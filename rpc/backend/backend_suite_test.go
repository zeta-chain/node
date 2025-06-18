package backend

import (
	"bufio"
	"math/big"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/suite"

	cmtrpctypes "github.com/cometbft/cometbft/rpc/core/types"

	dbm "github.com/cosmos/cosmos-db"
	evmdconfig "github.com/cosmos/evm/cmd/evmd/config"
	"github.com/cosmos/evm/crypto/hd"
	"github.com/cosmos/evm/encoding"
	"github.com/cosmos/evm/indexer"
	"github.com/cosmos/evm/testutil/constants"
	"github.com/cosmos/evm/testutil/integration/evm/network"
	utiltx "github.com/cosmos/evm/testutil/tx"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/zeta-chain/node/rpc/backend/mocks"
	rpctypes "github.com/zeta-chain/node/rpc/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type TestSuite struct {
	suite.Suite

	create  network.CreateEvmApp
	options []network.ConfigOption
	backend *Backend
	from    common.Address
	acc     sdk.AccAddress
	signer  keyring.Signer
}

func NewTestSuite(create network.CreateEvmApp, options ...network.ConfigOption) *TestSuite {
	return &TestSuite{
		create:  create,
		options: options,
	}
}

var ChainID = constants.ExampleChainID

// SetupTest is executed before every TestSuite test
func (s *TestSuite) SetupTest() {
	ctx := server.NewDefaultContext()
	ctx.Viper.Set("telemetry.global-labels", []interface{}{})
	ctx.Viper.Set("evm.evm-chain-id", evmdconfig.EVMChainID)

	baseDir := s.T().TempDir()
	nodeDirName := "node"
	clientDir := filepath.Join(baseDir, nodeDirName, "evmoscli")
	keyRing, err := s.generateTestKeyring(clientDir)
	if err != nil {
		panic(err)
	}

	// Create Account with set sequence
	s.acc = sdk.AccAddress(utiltx.GenerateAddress().Bytes())
	accounts := map[string]client.TestAccount{}
	accounts[s.acc.String()] = client.TestAccount{
		Address: s.acc,
		Num:     uint64(1),
		Seq:     uint64(1),
	}

	from, priv := utiltx.NewAddrKey()
	s.from = from
	s.signer = utiltx.NewSigner(priv)
	s.Require().NoError(err)

	nw := network.New(s.create, s.options...)
	encodingConfig := nw.GetEncodingConfig()
	clientCtx := client.Context{}.WithChainID(ChainID.ChainID).
		WithHeight(1).
		WithTxConfig(encodingConfig.TxConfig).
		WithKeyringDir(clientDir).
		WithKeyring(keyRing).
		WithAccountRetriever(client.TestAccountRetriever{Accounts: accounts}).
		WithClient(mocks.NewClient(s.T()))

	allowUnprotectedTxs := false
	idxer := indexer.NewKVIndexer(dbm.NewMemDB(), ctx.Logger, clientCtx)

	s.backend = NewBackend(ctx, ctx.Logger, clientCtx, allowUnprotectedTxs, idxer)
	s.backend.Cfg.JSONRPC.GasCap = 0
	s.backend.Cfg.JSONRPC.EVMTimeout = 0
	s.backend.Cfg.JSONRPC.AllowInsecureUnlock = true
	s.backend.Cfg.EVM.EVMChainID = 262144
	s.backend.QueryClient.QueryClient = mocks.NewEVMQueryClient(s.T())
	s.backend.QueryClient.FeeMarket = mocks.NewFeeMarketQueryClient(s.T())
	s.backend.Ctx = rpctypes.ContextWithHeight(1)

	// Add codec
	s.backend.ClientCtx.Codec = encodingConfig.Codec
}

// buildEthereumTx returns an example legacy Ethereum transaction
func (s *TestSuite) buildEthereumTx() (*evmtypes.MsgEthereumTx, []byte) {
	ethTxParams := evmtypes.EvmTxArgs{
		ChainID:  s.backend.EvmChainID,
		Nonce:    uint64(0),
		To:       &common.Address{},
		Amount:   big.NewInt(0),
		GasLimit: 100000,
		GasPrice: big.NewInt(1),
	}
	msgEthereumTx := evmtypes.NewTx(&ethTxParams)

	// A valid msg should have empty `From`
	msgEthereumTx.From = s.from.Hex()

	txBuilder := s.backend.ClientCtx.TxConfig.NewTxBuilder()
	err := txBuilder.SetMsgs(msgEthereumTx)
	s.Require().NoError(err)

	bz, err := s.backend.ClientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	s.Require().NoError(err)
	return msgEthereumTx, bz
}

// buildEthereumTx returns an example legacy Ethereum transaction
func (s *TestSuite) buildEthereumTxWithChainID(eip155ChainID *big.Int) *evmtypes.MsgEthereumTx {
	ethTxParams := evmtypes.EvmTxArgs{
		ChainID:  eip155ChainID,
		Nonce:    uint64(0),
		To:       &common.Address{},
		Amount:   big.NewInt(0),
		GasLimit: 100000,
		GasPrice: big.NewInt(1),
	}
	msgEthereumTx := evmtypes.NewTx(&ethTxParams)

	// A valid msg should have empty `From`
	msgEthereumTx.From = s.from.Hex()

	txBuilder := s.backend.ClientCtx.TxConfig.NewTxBuilder()
	err := txBuilder.SetMsgs(msgEthereumTx)
	s.Require().NoError(err)

	return msgEthereumTx
}

// buildFormattedBlock returns a formatted block for testing
func (s *TestSuite) buildFormattedBlock(
	blockRes *cmtrpctypes.ResultBlockResults,
	resBlock *cmtrpctypes.ResultBlock,
	fullTx bool,
	tx *evmtypes.MsgEthereumTx,
	validator sdk.AccAddress,
	baseFee *big.Int,
) map[string]interface{} {
	header := resBlock.Block.Header
	gasLimit := int64(^uint32(0))                                             // for `MaxGas = -1` (DefaultConsensusParams)
	gasUsed := new(big.Int).SetUint64(uint64(blockRes.TxsResults[0].GasUsed)) //nolint:gosec // G115 // won't exceed uint64

	root := common.Hash{}.Bytes()
	receipt := ethtypes.NewReceipt(root, false, gasUsed.Uint64())
	bloom := ethtypes.CreateBloom(receipt)

	ethRPCTxs := []interface{}{}
	if tx != nil {
		if fullTx {
			rpcTx, err := rpctypes.NewRPCTransaction(
				tx.AsTransaction(),
				common.BytesToHash(header.Hash()),
				uint64(header.Height), //nolint:gosec // G115 // won't exceed uint64
				uint64(0),
				baseFee,
				s.backend.EvmChainID,
			)
			s.Require().NoError(err)
			ethRPCTxs = []interface{}{rpcTx}
		} else {
			ethRPCTxs = []interface{}{common.HexToHash(tx.Hash)}
		}
	}

	return rpctypes.FormatBlock(
		header,
		resBlock.Block.Size(),
		gasLimit,
		gasUsed,
		ethRPCTxs,
		bloom,
		common.BytesToAddress(validator.Bytes()),
		baseFee,
	)
}

func (s *TestSuite) generateTestKeyring(clientDir string) (keyring.Keyring, error) {
	buf := bufio.NewReader(os.Stdin)
	encCfg := encoding.MakeConfig(evmdconfig.EVMChainID)
	return keyring.New(sdk.KeyringServiceName(), keyring.BackendTest, clientDir, buf, encCfg.Codec, []keyring.Option{hd.EthSecp256k1Option()}...)
}

func (s *TestSuite) signAndEncodeEthTx(msgEthereumTx *evmtypes.MsgEthereumTx) []byte {
	from, priv := utiltx.NewAddrKey()
	signer := utiltx.NewSigner(priv)

	ethSigner := ethtypes.LatestSigner(s.backend.ChainConfig())
	msgEthereumTx.From = from.String()
	err := msgEthereumTx.Sign(ethSigner, signer)
	s.Require().NoError(err)

	evmDenom := evmtypes.GetEVMCoinDenom()
	tx, err := msgEthereumTx.BuildTx(s.backend.ClientCtx.TxConfig.NewTxBuilder(), evmDenom)
	s.Require().NoError(err)

	txEncoder := s.backend.ClientCtx.TxConfig.TxEncoder()
	txBz, err := txEncoder(tx)
	s.Require().NoError(err)

	return txBz
}
