package backend

// import (
// 	"bufio"
// 	"math/big"
// 	"os"
// 	"path/filepath"
// 	"testing"

// 	"github.com/ethereum/go-ethereum/common"
// 	ethtypes "github.com/ethereum/go-ethereum/core/types"
// 	"github.com/stretchr/testify/suite"

// 	abci "github.com/cometbft/cometbft/abci/types"
// 	cmtrpctypes "github.com/cometbft/cometbft/rpc/core/types"

// 	dbm "github.com/cosmos/cosmos-db"
// 	"github.com/cosmos/evm/crypto/hd"
// 	"github.com/cosmos/evm/encoding"
// 	"github.com/cosmos/evm/indexer"
// 	"github.com/cosmos/evm/testutil/constants"
// 	utiltx "github.com/cosmos/evm/testutil/tx"
// 	evmtypes "github.com/cosmos/evm/x/vm/types"
// 	"github.com/zeta-chain/node/app"
// 	"github.com/zeta-chain/node/rpc/backend/mocks"
// 	rpctypes "github.com/zeta-chain/node/rpc/types"

// 	"github.com/cosmos/cosmos-sdk/client"
// 	"github.com/cosmos/cosmos-sdk/crypto/keyring"
// 	"github.com/cosmos/cosmos-sdk/server"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	zetacoredconfig "github.com/zeta-chain/node/cmd/zetacored/config"
// 	protov2 "google.golang.org/protobuf/proto"
// )

// type TestSuite struct {
// 	suite.Suite

// 	backend *Backend
// 	from    common.Address
// 	acc     sdk.AccAddress
// 	signer  keyring.Signer
// }

// var ChainID = constants.ExampleChainID

// // testTx is a dummy implementation of cosmos Tx used for testing.
// type testTx struct {
// }

// func (tx testTx) GetMsgs() []sdk.Msg                         { return nil }
// func (tx testTx) GetMsgsV2() ([]protov2.Message, error)      { return nil, nil }
// func (tx testTx) GetSigners() []sdk.AccAddress               { return nil }
// func (tx testTx) ValidateBasic() error                       { return nil }
// func (t testTx) ProtoMessage()                               { panic("not implemented") }
// func (t testTx) Reset()                                      { panic("not implemented") }
// func (t testTx) String() string                              { panic("not implemented") }
// func (t testTx) Bytes() []byte                               { panic("not implemented") }
// func (t testTx) VerifySignature(msg []byte, sig []byte) bool { panic("not implemented") }
// func (t testTx) Type() string                                { panic("not implemented") }

// var (
// 	_ sdk.Tx  = (*testTx)(nil)
// 	_ sdk.Msg = (*testTx)(nil)
// )

// func TestBackendTestSuite(t *testing.T) {
// 	suite.Run(t, new(TestSuite))
// }

// // SetupTest is executed before every TestSuite test
// func (s *TestSuite) SetupTest() {
// 	ctx := server.NewDefaultContext()
// 	ctx.Viper.Set("telemetry.global-labels", []interface{}{})
// 	ctx.Viper.Set("evm.evm-chain-id", ChainID.EVMChainID)

// 	baseDir := s.T().TempDir()
// 	nodeDirName := "node"
// 	clientDir := filepath.Join(baseDir, nodeDirName, "evmoscli")
// 	keyRing, err := s.generateTestKeyring(clientDir)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Create Account with set sequence
// 	s.acc = sdk.AccAddress(utiltx.GenerateAddress().Bytes())
// 	accounts := map[string]client.TestAccount{}
// 	accounts[s.acc.String()] = client.TestAccount{
// 		Address: s.acc,
// 		Num:     uint64(1),
// 		Seq:     uint64(1),
// 	}

// 	from, priv := utiltx.NewAddrKey()
// 	s.from = from
// 	s.signer = utiltx.NewSigner(priv)
// 	s.Require().NoError(err)

// 	encodingConfig := app.MakeEncodingConfig(ChainID.EVMChainID)
// 	clientCtx := client.Context{}.WithChainID(ChainID.ChainID).
// 		WithHeight(1).
// 		WithTxConfig(encodingConfig.TxConfig).
// 		WithKeyringDir(clientDir).
// 		WithKeyring(keyRing).
// 		WithAccountRetriever(client.TestAccountRetriever{Accounts: accounts}).
// 		WithClient(mocks.NewClient(s.T()))

// 	// TODO https://github.com/zeta-chain/node/issues/2157
// 	// rpc tests in evm module are integration tests, this workaround is needed to setup eth config
// 	// consider maintaining all synthetic txs changes and tests in evm module fork, so there is no need to have
// 	// only rpc forked in node repo
// 	ethCfg := evmtypes.DefaultChainConfig(ChainID.EVMChainID)
// 	configurator := evmtypes.NewEVMConfigurator()
// 	configurator.ResetTestConfig()
// 	err = configurator.
// 		WithChainConfig(ethCfg).
// 		WithEVMCoinInfo(evmtypes.EvmCoinInfo{
// 			Denom:         zetacoredconfig.BaseDenom,
// 			ExtendedDenom: zetacoredconfig.BaseDenom,
// 			DisplayDenom:  zetacoredconfig.BaseDenom,
// 			Decimals:      18,
// 		}).
// 		Configure()

// 	allowUnprotectedTxs := false
// 	idxer := indexer.NewKVIndexer(dbm.NewMemDB(), ctx.Logger, clientCtx)

// 	s.backend = NewBackend(ctx, ctx.Logger, clientCtx, allowUnprotectedTxs, idxer)
// 	s.backend.Cfg.EVM.EVMChainID = ChainID.EVMChainID
// 	s.backend.EvmChainID = big.NewInt(int64(ChainID.EVMChainID))
// 	s.backend.Cfg.JSONRPC.GasCap = 0
// 	s.backend.Cfg.JSONRPC.EVMTimeout = 0
// 	s.backend.Cfg.JSONRPC.AllowInsecureUnlock = true
// 	s.backend.QueryClient.QueryClient = mocks.NewEVMQueryClient(s.T())
// 	s.backend.QueryClient.FeeMarket = mocks.NewFeeMarketQueryClient(s.T())
// 	s.backend.Ctx = rpctypes.ContextWithHeight(1)

// 	// Add codec
// 	s.backend.ClientCtx.Codec = encodingConfig.Codec
// }

// // buildEthereumTx returns an example legacy Ethereum transaction
// func (s *TestSuite) buildEthereumTx() (*evmtypes.MsgEthereumTx, []byte) {
// 	ethTxParams := evmtypes.EvmTxArgs{
// 		ChainID:  s.backend.EvmChainID,
// 		Nonce:    uint64(0),
// 		To:       &common.Address{},
// 		Amount:   big.NewInt(0),
// 		GasLimit: 100000,
// 		GasPrice: big.NewInt(1),
// 	}
// 	msgEthereumTx := evmtypes.NewTx(&ethTxParams)

// 	s.signAndEncodeEthTx(msgEthereumTx)

// 	txBuilder := s.backend.ClientCtx.TxConfig.NewTxBuilder()
// 	err := txBuilder.SetMsgs(msgEthereumTx)
// 	s.Require().NoError(err)

// 	bz, err := s.backend.ClientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
// 	s.Require().NoError(err)
// 	return msgEthereumTx, bz
// }

// func (suite *TestSuite) buildSyntheticTxResult(txHash string) ([]byte, abci.ExecTxResult) {
// 	testTx := &testTx{}
// 	txBuilder := suite.backend.ClientCtx.TxConfig.NewTxBuilder()
// 	txBuilder.SetSignatures()
// 	txBuilder.SetMsgs(testTx)

// 	bz, _ := suite.backend.ClientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
// 	return bz, abci.ExecTxResult{
// 		Code: 0,
// 		Events: []abci.Event{
// 			{Type: evmtypes.EventTypeEthereumTx, Attributes: []abci.EventAttribute{
// 				{Key: "ethereumTxHash", Value: txHash},
// 				{Key: "txIndex", Value: "8888"},
// 				{Key: "txData", Value: "0x1234"},
// 				{Key: "amount", Value: "1000"},
// 				{Key: "txGasUsed", Value: "21000"},
// 				{Key: "txGasLimit", Value: "21000"},
// 				{Key: "txHash", Value: ""},
// 				{Key: "recipient", Value: "0x775b87ef5D82ca211811C1a02CE0fE0CA3a455d7"},
// 			}},
// 			{
// 				Type: "message", Attributes: []abci.EventAttribute{
// 					{Key: "sender", Value: "0x735b14BB79463307AAcBED86DAf3322B1e6226aB"},
// 					{Key: "txType", Value: "88"},
// 					{Key: "txNonce", Value: "1"},
// 				},
// 			},
// 		},
// 	}
// }

// // buildEthereumTx returns an example legacy Ethereum transaction
// func (s *TestSuite) buildEthereumTxWithChainID(eip155ChainID *big.Int) *evmtypes.MsgEthereumTx {
// 	ethTxParams := evmtypes.EvmTxArgs{
// 		ChainID:   eip155ChainID,
// 		Nonce:     uint64(0),
// 		To:        &common.Address{},
// 		Amount:    big.NewInt(0),
// 		GasLimit:  100000,
// 		GasPrice:  big.NewInt(1),
// 		GasFeeCap: big.NewInt(1),
// 	}
// 	msgEthereumTx := evmtypes.NewTx(&ethTxParams)

// 	// A valid msg should have empty `From`
// 	msgEthereumTx.From = s.from.Bytes()

// 	txBuilder := s.backend.ClientCtx.TxConfig.NewTxBuilder()
// 	err := txBuilder.SetMsgs(msgEthereumTx)
// 	s.Require().NoError(err)

// 	return msgEthereumTx
// }

// // buildFormattedBlock returns a formatted block for testing
// func (s *TestSuite) buildFormattedBlock(
// 	blockRes *cmtrpctypes.ResultBlockResults,
// 	resBlock *cmtrpctypes.ResultBlock,
// 	fullTx bool,
// 	tx *evmtypes.MsgEthereumTx,
// 	validator sdk.AccAddress,
// 	baseFee *big.Int,
// ) map[string]interface{} {
// 	header := resBlock.Block.Header
// 	gasLimit := int64(
// 		^uint32(0),
// 	) // for `MaxGas = -1` (DefaultConsensusParams)
// 	gasUsed := new(
// 		big.Int,
// 	).SetUint64(uint64(blockRes.TxsResults[0].GasUsed)) //#nosec G115 won't exceed uint64

// 	root := common.Hash{}.Bytes()
// 	receipt := ethtypes.NewReceipt(root, false, gasUsed.Uint64())
// 	bloom := ethtypes.CreateBloom(receipt)

// 	ethRPCTxs := []interface{}{}
// 	if tx != nil {
// 		if fullTx {
// 			rpcTx, err := rpctypes.NewRPCTransaction(
// 				tx,
// 				common.BytesToHash(header.Hash()),
// 				uint64(header.Height), //#nosec G115 won't exceed uint64
// 				uint64(0),
// 				baseFee,
// 				s.backend.EvmChainID,
// 			)
// 			s.Require().NoError(err)
// 			ethRPCTxs = []interface{}{rpcTx}
// 		} else {
// 			ethRPCTxs = []interface{}{common.HexToHash(tx.Hash)}
// 		}
// 	}

// 	return rpctypes.FormatBlock(
// 		header,
// 		resBlock.Block.Size(),
// 		gasLimit,
// 		gasUsed,
// 		ethRPCTxs,
// 		bloom,
// 		common.BytesToAddress(validator.Bytes()),
// 		baseFee,
// 	)
// }

// func (s *TestSuite) generateTestKeyring(clientDir string) (keyring.Keyring, error) {
// 	buf := bufio.NewReader(os.Stdin)
// 	encCfg := encoding.MakeConfig(ChainID.EVMChainID)
// 	return keyring.New(
// 		sdk.KeyringServiceName(),
// 		keyring.BackendTest,
// 		clientDir,
// 		buf,
// 		encCfg.Codec,
// 		[]keyring.Option{hd.EthSecp256k1Option()}...)
// }

// func (s *TestSuite) signAndEncodeEthTx(msgEthereumTx *evmtypes.MsgEthereumTx) []byte {
// 	from, priv := utiltx.NewAddrKey()
// 	signer := utiltx.NewSigner(priv)

// 	ethSigner := ethtypes.LatestSigner(s.backend.ChainConfig())
// 	msgEthereumTx.From = from.Bytes()
// 	err := msgEthereumTx.Sign(ethSigner, signer)
// 	s.Require().NoError(err)

// 	evmDenom := evmtypes.GetEVMCoinDenom()
// 	tx, err := msgEthereumTx.BuildTx(s.backend.ClientCtx.TxConfig.NewTxBuilder(), evmDenom)
// 	s.Require().NoError(err)

// 	txEncoder := s.backend.ClientCtx.TxConfig.TxEncoder()
// 	txBz, err := txEncoder(tx)
// 	s.Require().NoError(err)

// 	return txBz
// }
