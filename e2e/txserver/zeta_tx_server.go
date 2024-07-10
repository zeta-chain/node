package txserver

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/crypto/hd"
	etherminttypes "github.com/evmos/ethermint/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"

	"github.com/zeta-chain/zetacore/app"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	emissionstypes "github.com/zeta-chain/zetacore/x/emissions/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// EmissionsPoolAddress is the address of the emissions pool
// This address is constant for all networks because it is derived from emissions name
const EmissionsPoolAddress = "zeta1w43fn2ze2wyhu5hfmegr6vp52c3dgn0srdgymy"

// ZetaTxServer is a ZetaChain tx server for E2E test
type ZetaTxServer struct {
	clientCtx    client.Context
	txFactory    tx.Factory
	name         []string
	mnemonic     []string
	address      []string
	blockTimeout time.Duration
}

// NewZetaTxServer returns a new TxServer with provided account
func NewZetaTxServer(rpcAddr string, names []string, privateKeys []string, chainID string) (*ZetaTxServer, error) {
	ctx := context.Background()

	if len(names) == 0 {
		return nil, errors.New("no account provided")
	}

	if len(names) != len(privateKeys) {
		return nil, errors.New("invalid names and privateKeys")
	}

	// initialize rpc and check status
	rpc, err := rpchttp.New(rpcAddr, "/websocket")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize rpc: %s", err.Error())
	}
	if _, err = rpc.Status(ctx); err != nil {
		return nil, fmt.Errorf("failed to query rpc: %s", err.Error())
	}

	// initialize codec
	cdc, reg := newCodec()

	// initialize keyring
	kr := keyring.NewInMemory(cdc, hd.EthSecp256k1Option())

	addresses := make([]string, 0, len(names))

	// create accounts
	for i := range names {
		err = kr.ImportPrivKeyHex(names[i], privateKeys[i], string(hd.EthSecp256k1Type))
		if err != nil {
			return nil, fmt.Errorf("failed to create account: %w", err)
		}
		r, err := kr.Key(names[i])
		if err != nil {
			return nil, fmt.Errorf("failed to get account key: %w", err)
		}
		accAddr, err := r.GetAddress()
		if err != nil {
			return nil, fmt.Errorf("failed to get account address: %w", err)
		}

		addresses = append(addresses, accAddr.String())
	}

	clientCtx := newContext(rpc, cdc, reg, kr, chainID)
	txf := newFactory(clientCtx)

	return &ZetaTxServer{
		clientCtx:    clientCtx,
		txFactory:    txf,
		name:         names,
		address:      addresses,
		blockTimeout: 2 * time.Minute,
	}, nil
}

// GetAccountName returns the account name from the given index
// returns empty string if index is out of bound, error should be handled by caller
func (zts ZetaTxServer) GetAccountName(index int) string {
	if index >= len(zts.name) {
		return ""
	}
	return zts.name[index]
}

// GetAccountAddress returns the account address from the given index
// returns empty string if index is out of bound, error should be handled by caller
func (zts ZetaTxServer) GetAccountAddress(index int) string {
	if index >= len(zts.address) {
		return ""
	}
	return zts.address[index]
}

// GetAccountAddressFromName returns the account address from the given name
func (zts ZetaTxServer) GetAccountAddressFromName(name string) (string, error) {
	acc, err := zts.clientCtx.Keyring.Key(name)
	if err != nil {
		return "", err
	}
	addr, err := acc.GetAddress()
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

// MustGetAccountAddressFromName returns the account address from the given name.It panics on error
// and should be used in tests only
func (zts ZetaTxServer) MustGetAccountAddressFromName(name string) string {
	acc, err := zts.clientCtx.Keyring.Key(name)
	if err != nil {
		panic(err)
	}
	addr, err := acc.GetAddress()
	if err != nil {
		panic(err)
	}
	return addr.String()
}

// GetAllAccountAddress returns all account addresses
func (zts ZetaTxServer) GetAllAccountAddress() []string {
	return zts.address
}

// GetAccountMnemonic returns the account name from the given index
// returns empty string if index is out of bound, error should be handled by caller
func (zts ZetaTxServer) GetAccountMnemonic(index int) string {
	if index >= len(zts.mnemonic) {
		return ""
	}
	return zts.mnemonic[index]
}

// BroadcastTx broadcasts a tx to ZetaChain with the provided msg from the account
// and waiting for blockTime for tx to be included in the block
func (zts ZetaTxServer) BroadcastTx(account string, msg sdktypes.Msg) (*sdktypes.TxResponse, error) {
	// Find number and sequence and set it
	acc, err := zts.clientCtx.Keyring.Key(account)
	if err != nil {
		return nil, err
	}
	addr, err := acc.GetAddress()
	if err != nil {
		return nil, err
	}
	accountNumber, accountSeq, err := zts.clientCtx.AccountRetriever.GetAccountNumberSequence(zts.clientCtx, addr)
	if err != nil {
		return nil, err
	}
	zts.txFactory = zts.txFactory.WithAccountNumber(accountNumber).WithSequence(accountSeq)

	txBuilder, err := zts.txFactory.BuildUnsignedTx(msg)
	if err != nil {
		return nil, err
	}

	// Sign tx
	err = tx.Sign(zts.txFactory, account, txBuilder, true)
	if err != nil {
		return nil, err
	}
	txBytes, err := zts.clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}
	return broadcastWithBlockTimeout(zts, txBytes)
}

func broadcastWithBlockTimeout(zts ZetaTxServer, txBytes []byte) (*sdktypes.TxResponse, error) {
	res, err := zts.clientCtx.BroadcastTx(txBytes)
	if err != nil {
		if res == nil {
			return nil, err
		}
		return &sdktypes.TxResponse{
			Code:      res.Code,
			Codespace: res.Codespace,
			TxHash:    res.TxHash,
		}, err
	}

	exitAfter := time.After(zts.blockTimeout)
	hash, err := hex.DecodeString(res.TxHash)
	if err != nil {
		return nil, err
	}
	for {
		select {
		case <-exitAfter:
			return nil, fmt.Errorf("timed out after waiting for tx to get included in the block: %d", zts.blockTimeout)
		case <-time.After(time.Millisecond * 100):
			resTx, err := zts.clientCtx.Client.Tx(context.TODO(), hash, false)
			if err == nil {
				return mkTxResult(zts.clientCtx, resTx)
			}
		}
	}
}

func mkTxResult(clientCtx client.Context, resTx *coretypes.ResultTx) (*sdktypes.TxResponse, error) {
	txb, err := clientCtx.TxConfig.TxDecoder()(resTx.Tx)
	if err != nil {
		return nil, err
	}
	p, ok := txb.(intoAny)
	if !ok {
		return nil, fmt.Errorf("expecting a type implementing intoAny, got: %T", txb)
	}
	resBlock, err := clientCtx.Client.Block(context.TODO(), &resTx.Height)
	if err != nil {
		return nil, err
	}
	return sdktypes.NewResponseResultTx(resTx, p.AsAny(), resBlock.Block.Time.Format(time.RFC3339)), nil
}

type intoAny interface {
	AsAny() *codectypes.Any
}

// EnableHeaderVerification enables the header verification for the given chain IDs
func (zts ZetaTxServer) EnableHeaderVerification(account string, chainIDList []int64) error {
	// retrieve account
	acc, err := zts.clientCtx.Keyring.Key(account)
	if err != nil {
		return err
	}
	addr, err := acc.GetAddress()
	if err != nil {
		return err
	}

	_, err = zts.BroadcastTx(account, lightclienttypes.NewMsgEnableHeaderVerification(
		addr.String(),
		chainIDList,
	))
	return err
}

// DeploySystemContractsAndZRC20 deploys the system contracts and ZRC20 contracts
// returns the addresses of uniswap factory, router and erc20 zrc20
func (zts ZetaTxServer) DeploySystemContractsAndZRC20(
	account, erc20Addr string,
) (string, string, string, string, string, error) {
	// retrieve account
	acc, err := zts.clientCtx.Keyring.Key(account)
	if err != nil {
		return "", "", "", "", "", err
	}
	addr, err := acc.GetAddress()
	if err != nil {
		return "", "", "", "", "", err
	}

	// deploy new system contracts
	res, err := zts.BroadcastTx(account, fungibletypes.NewMsgDeploySystemContracts(addr.String()))
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("failed to deploy system contracts: %s", err.Error())
	}

	systemContractAddress, err := FetchAttributeFromTxResponse(res, "system_contract")
	if err != nil {
		return "", "", "", "", "", fmt.Errorf(
			"failed to fetch system contract address: %s; rawlog %s",
			err.Error(),
			res.RawLog,
		)
	}

	// get system contract
	_, err = zts.BroadcastTx(account, fungibletypes.NewMsgUpdateSystemContract(addr.String(), systemContractAddress))
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("failed to set system contract: %s", err.Error())
	}

	// get uniswap contract addresses
	uniswapV2FactoryAddr, err := FetchAttributeFromTxResponse(res, "uniswap_v2_factory")
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("failed to fetch uniswap v2 factory address: %s", err.Error())
	}
	uniswapV2RouterAddr, err := FetchAttributeFromTxResponse(res, "uniswap_v2_router")
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("failed to fetch uniswap v2 router address: %s", err.Error())
	}

	// get zevm connector address
	zevmConnectorAddr, err := FetchAttributeFromTxResponse(res, "connector_zevm")
	if err != nil {
		return "", "", "", "", "", fmt.Errorf(
			"failed to fetch zevm connector address: %s, txResponse: %s",
			err.Error(),
			res.String(),
		)
	}

	// get wzeta address
	wzetaAddr, err := FetchAttributeFromTxResponse(res, "wzeta")
	if err != nil {
		return "", "", "", "", "", fmt.Errorf(
			"failed to fetch wzeta address: %s, txResponse: %s",
			err.Error(),
			res.String(),
		)
	}

	// deploy eth zrc20
	_, err = zts.BroadcastTx(account, fungibletypes.NewMsgDeployFungibleCoinZRC20(
		addr.String(),
		"",
		chains.GoerliLocalnet.ChainId,
		18,
		"ETH",
		"gETH",
		coin.CoinType_Gas,
		100000,
	))
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("failed to deploy eth zrc20: %s", err.Error())
	}

	// deploy btc zrc20
	_, err = zts.BroadcastTx(account, fungibletypes.NewMsgDeployFungibleCoinZRC20(
		addr.String(),
		"",
		chains.BitcoinRegtest.ChainId,
		8,
		"BTC",
		"tBTC",
		coin.CoinType_Gas,
		100000,
	))
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("failed to deploy btc zrc20: %s", err.Error())
	}

	// deploy erc20 zrc20
	res, err = zts.BroadcastTx(account, fungibletypes.NewMsgDeployFungibleCoinZRC20(
		addr.String(),
		erc20Addr,
		chains.GoerliLocalnet.ChainId,
		6,
		"USDT",
		"USDT",
		coin.CoinType_ERC20,
		100000,
	))
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("failed to deploy erc20 zrc20: %s", err.Error())
	}

	// fetch the erc20 zrc20 contract address and remove the quotes
	erc20zrc20Addr, err := FetchAttributeFromTxResponse(res, "Contract")
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("failed to fetch erc20 zrc20 contract address: %s", err.Error())
	}
	if !ethcommon.IsHexAddress(erc20zrc20Addr) {
		return "", "", "", "", "", fmt.Errorf("invalid address in event: %s", erc20zrc20Addr)
	}
	return uniswapV2FactoryAddr, uniswapV2RouterAddr, zevmConnectorAddr, wzetaAddr, erc20zrc20Addr, nil
}

// FundEmissionsPool funds the emissions pool with the given amount
func (zts ZetaTxServer) FundEmissionsPool(account string, amount *big.Int) error {
	// retrieve account
	acc, err := zts.clientCtx.Keyring.Key(account)
	if err != nil {
		return err
	}
	addr, err := acc.GetAddress()
	if err != nil {
		return err
	}

	// retrieve account address
	emissionPoolAccAddr, err := sdktypes.AccAddressFromBech32(EmissionsPoolAddress)
	if err != nil {
		return err
	}

	// convert amount
	amountInt := sdktypes.NewIntFromBigInt(amount)

	// fund emissions pool
	_, err = zts.BroadcastTx(account, banktypes.NewMsgSend(
		addr,
		emissionPoolAccAddr,
		sdktypes.NewCoins(sdktypes.NewCoin(config.BaseDenom, amountInt)),
	))
	return err
}

// newCodec returns the codec for msg server
func newCodec() (*codec.ProtoCodec, codectypes.InterfaceRegistry) {
	encodingConfig := app.MakeEncodingConfig()
	interfaceRegistry := encodingConfig.InterfaceRegistry
	cdc := codec.NewProtoCodec(interfaceRegistry)

	sdktypes.RegisterInterfaces(interfaceRegistry)
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	authz.RegisterInterfaces(interfaceRegistry)
	banktypes.RegisterInterfaces(interfaceRegistry)
	stakingtypes.RegisterInterfaces(interfaceRegistry)
	slashingtypes.RegisterInterfaces(interfaceRegistry)
	upgradetypes.RegisterInterfaces(interfaceRegistry)
	distrtypes.RegisterInterfaces(interfaceRegistry)
	evidencetypes.RegisterInterfaces(interfaceRegistry)
	crisistypes.RegisterInterfaces(interfaceRegistry)
	evmtypes.RegisterInterfaces(interfaceRegistry)
	etherminttypes.RegisterInterfaces(interfaceRegistry)
	crosschaintypes.RegisterInterfaces(interfaceRegistry)
	emissionstypes.RegisterInterfaces(interfaceRegistry)
	fungibletypes.RegisterInterfaces(interfaceRegistry)
	observertypes.RegisterInterfaces(interfaceRegistry)
	lightclienttypes.RegisterInterfaces(interfaceRegistry)
	authoritytypes.RegisterInterfaces(interfaceRegistry)

	return cdc, interfaceRegistry
}

// newContext returns the client context for msg server
func newContext(
	rpc *rpchttp.HTTP,
	cdc *codec.ProtoCodec,
	reg codectypes.InterfaceRegistry,
	kr keyring.Keyring,
	chainID string,
) client.Context {
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)
	return client.Context{}.
		WithChainID(chainID).
		WithInterfaceRegistry(reg).
		WithCodec(cdc).
		WithTxConfig(txConfig).
		WithLegacyAmino(codec.NewLegacyAmino()).
		WithInput(os.Stdin).
		WithOutput(os.Stdout).
		WithBroadcastMode(flags.BroadcastSync).
		WithClient(rpc).
		WithSkipConfirmation(true).
		WithFromName("creator").
		WithFromAddress(sdktypes.AccAddress{}).
		WithKeyring(kr).
		WithAccountRetriever(authtypes.AccountRetriever{})
}

// newFactory returns the tx factory for msg server
func newFactory(clientCtx client.Context) tx.Factory {
	return tx.Factory{}.
		WithChainID(clientCtx.ChainID).
		WithKeybase(clientCtx.Keyring).
		WithGas(10000000).
		WithGasAdjustment(1).
		WithSignMode(signing.SignMode_SIGN_MODE_UNSPECIFIED).
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithTxConfig(clientCtx.TxConfig).
		WithFees("100000000000000000azeta")
}

type messageLog struct {
	Events []event `json:"events"`
}

type event struct {
	Type       string      `json:"type"`
	Attributes []attribute `json:"attributes"`
}

type attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// FetchAttributeFromTxResponse fetches the attribute from the tx response
func FetchAttributeFromTxResponse(res *sdktypes.TxResponse, key string) (string, error) {
	var logs []messageLog
	err := json.Unmarshal([]byte(res.RawLog), &logs)
	if err != nil {
		return "", err
	}

	var attributes []string
	for _, log := range logs {
		for _, event := range log.Events {
			for _, attr := range event.Attributes {
				attributes = append(attributes, attr.Key)
				if strings.EqualFold(attr.Key, key) {
					address := attr.Value

					if len(address) < 2 {
						return "", fmt.Errorf("invalid address: %s", address)
					}

					// trim the quotes
					address = address[1 : len(address)-1]

					return address, nil
				}
			}
		}
	}

	return "", fmt.Errorf("attribute %s not found, attributes:  %+v", key, attributes)
}
