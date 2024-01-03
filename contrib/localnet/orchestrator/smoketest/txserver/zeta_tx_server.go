package txserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
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
	etherminttypes "github.com/evmos/ethermint/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"github.com/zeta-chain/zetacore/common"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	emissionstypes "github.com/zeta-chain/zetacore/x/emissions/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// ZetaTxServer is a ZetaChain tx server for smoke test
type ZetaTxServer struct {
	clientCtx client.Context
	txFactory tx.Factory
	name      []string
	mnemonic  []string
	address   []string
}

// NewZetaTxServer returns a new TxServer with provided account
func NewZetaTxServer(rpcAddr string, names []string, mnemonics []string, chainID string) (ZetaTxServer, error) {
	ctx := context.Background()

	if len(names) == 0 {
		return ZetaTxServer{}, errors.New("no account provided")
	}

	if len(names) != len(mnemonics) {
		return ZetaTxServer{}, errors.New("invalid names and mnemonics")
	}

	// initialize rpc and check status
	rpc, err := rpchttp.New(rpcAddr, "/websocket")
	if err != nil {
		return ZetaTxServer{}, fmt.Errorf("failed to initialize rpc: %s", err.Error())
	}
	if _, err = rpc.Status(ctx); err != nil {
		return ZetaTxServer{}, fmt.Errorf("failed to query rpc: %s", err.Error())
	}

	// initialize codec
	cdc, reg := newCodec()

	// initialize keyring
	kr := keyring.NewInMemory(cdc)

	addresses := make([]string, 0, len(names))

	// create accounts
	for i := range names {
		r, err := kr.NewAccount(names[i], mnemonics[i], "", sdktypes.FullFundraiserPath, hd.Secp256k1)
		if err != nil {
			return ZetaTxServer{}, fmt.Errorf("failed to create account: %s", err.Error())
		}
		accAddr, err := r.GetAddress()
		if err != nil {
			return ZetaTxServer{}, fmt.Errorf("failed to get account address: %s", err.Error())
		}

		addresses = append(addresses, accAddr.String())
	}

	clientCtx := newContext(rpc, cdc, reg, kr, chainID)
	txf := newFactory(clientCtx)

	return ZetaTxServer{
		clientCtx: clientCtx,
		txFactory: txf,
		name:      names,
		mnemonic:  mnemonics,
		address:   addresses,
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

// GetAccountMnemonic returns the account name from the given index
// returns empty string if index is out of bound, error should be handled by caller
func (zts ZetaTxServer) GetAccountMnemonic(index int) string {
	if index >= len(zts.mnemonic) {
		return ""
	}
	return zts.mnemonic[index]
}

// BroadcastTx broadcasts a tx to ZetaChain with the provided msg from the account
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

	// Broadcast tx
	return zts.clientCtx.BroadcastTx(txBytes)
}

// DeploySystemContractsAndZRC20 deploys the system contracts and ZRC20 contracts
// returns the addresses of uniswap factory, router and usdt zrc20
func (zts ZetaTxServer) DeploySystemContractsAndZRC20(account, usdtERC20Addr string) (string, string, string, error) {
	// retrieve account
	acc, err := zts.clientCtx.Keyring.Key(account)
	if err != nil {
		return "", "", "", err
	}
	addr, err := acc.GetAddress()
	if err != nil {
		return "", "", "", err
	}

	// deploy new system contracts
	res, err := zts.BroadcastTx(account, fungibletypes.NewMsgDeploySystemContracts(addr.String()))
	if err != nil {
		return "", "", "", fmt.Errorf("failed to deploy system contracts: %s", err.Error())
	}

	systemContractAddress, err := fetchAttribute(res, "system_contract")
	if err != nil {
		return "", "", "", fmt.Errorf("failed to fetch system contract address: %s; rawlog %s", err.Error(), res.RawLog)
	}

	// set system contract
	_, err = zts.BroadcastTx(account, fungibletypes.NewMsgUpdateSystemContract(addr.String(), systemContractAddress))
	if err != nil {
		return "", "", "", fmt.Errorf("failed to set system contract: %s", err.Error())
	}

	// set uniswap contract addresses
	uniswapV2FactoryAddr, err := fetchAttribute(res, "uniswap_v2_factory")
	if err != nil {
		return "", "", "", fmt.Errorf("failed to fetch uniswap v2 factory address: %s", err.Error())
	}
	uniswapV2RouterAddr, err := fetchAttribute(res, "uniswap_v2_router")
	if err != nil {
		return "", "", "", fmt.Errorf("failed to fetch uniswap v2 router address: %s", err.Error())
	}

	// deploy eth zrc20
	_, err = zts.BroadcastTx(account, fungibletypes.NewMsgDeployFungibleCoinZRC20(
		addr.String(),
		"",
		common.GoerliLocalnetChain().ChainId,
		18,
		"ETH",
		"gETH",
		common.CoinType_Gas,
		100000,
	))
	if err != nil {
		return "", "", "", fmt.Errorf("failed to deploy eth zrc20: %s", err.Error())
	}

	// deploy btc zrc20
	_, err = zts.BroadcastTx(account, fungibletypes.NewMsgDeployFungibleCoinZRC20(
		addr.String(),
		"",
		common.BtcRegtestChain().ChainId,
		8,
		"BTC",
		"tBTC",
		common.CoinType_Gas,
		100000,
	))
	if err != nil {
		return "", "", "", fmt.Errorf("failed to deploy btc zrc20: %s", err.Error())
	}

	// deploy usdt zrc20
	res, err = zts.BroadcastTx(account, fungibletypes.NewMsgDeployFungibleCoinZRC20(
		addr.String(),
		usdtERC20Addr,
		common.GoerliLocalnetChain().ChainId,
		6,
		"USDT",
		"USDT",
		common.CoinType_ERC20,
		100000,
	))
	if err != nil {
		return "", "", "", fmt.Errorf("failed to deploy usdt zrc20: %s", err.Error())
	}

	// fetch the usdt zrc20 contract address and remove the quotes
	usdtZRC20Addr, err := fetchAttribute(res, "Contract")
	if err != nil {
		return "", "", "", fmt.Errorf("failed to fetch usdt zrc20 contract address: %s", err.Error())
	}
	if !ethcommon.IsHexAddress(usdtZRC20Addr) {
		return "", "", "", fmt.Errorf("invalid address in event: %s", usdtZRC20Addr)
	}
	return uniswapV2FactoryAddr, uniswapV2RouterAddr, usdtZRC20Addr, nil
}

// InitializeCoreParams sets the core params with local Goerli and BtcRegtest chains enabled
func (zts ZetaTxServer) InitializeCoreParams(account string) error {
	// set btc regtest  core params
	btcCoreParams := observertypes.GetDefaultBtcRegtestCoreParams()
	btcCoreParams.IsSupported = true
	if err := zts.UpdateCoreParams(account, btcCoreParams); err != nil {
		return fmt.Errorf("failed to set core params for bitcoin: %s", err.Error())
	}

	// set goerli localnet core params
	goerliCoreParams := observertypes.GetDefaultGoerliLocalnetCoreParams()
	goerliCoreParams.IsSupported = true
	if err := zts.UpdateCoreParams(account, goerliCoreParams); err != nil {
		return fmt.Errorf("failed to set core params for bitcoin: %s", err.Error())
	}

	return nil
}

// UpdateCoreParams updates the core params
func (zts ZetaTxServer) UpdateCoreParams(account string, cp *observertypes.CoreParams) error {
	// retrieve account
	acc, err := zts.clientCtx.Keyring.Key(account)
	if err != nil {
		return err
	}
	addr, err := acc.GetAddress()
	if err != nil {
		return err
	}

	_, err = zts.BroadcastTx(account, observertypes.NewMsgUpdateCoreParams(addr.String(), cp))
	if err != nil {
		return fmt.Errorf("failed to set core params for bitcoin: %s", err.Error())
	}

	return nil
}

// newCodec returns the codec for msg server
func newCodec() (*codec.ProtoCodec, codectypes.InterfaceRegistry) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
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
		WithBroadcastMode(flags.BroadcastBlock).
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

// fetchAttribute fetches the attribute from the tx response
func fetchAttribute(res *sdktypes.TxResponse, key string) (string, error) {
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
