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
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	emissionstypes "github.com/zeta-chain/zetacore/x/emissions/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// EmissionsPoolAddress is the address of the emissions pool
// This address is constant for all networks because it is derived from emissions name
const EmissionsPoolAddress = "zeta1w43fn2ze2wyhu5hfmegr6vp52c3dgn0srdgymy"

// ZetaTxServer is a ZetaChain tx server for E2E test
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

	// Broadcast tx and wait 10 mins to be included in block
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

	var blockTimeout time.Duration = 10 * time.Minute
	exitAfter := time.After(blockTimeout)
	hash, err := hex.DecodeString(res.TxHash)
	if err != nil {
		return nil, err
	}
	for {
		select {
		case <-exitAfter:
			return nil, fmt.Errorf("timed out after waiting for tx to get included in the block: %d", blockTimeout)
		case <-time.After(time.Millisecond * 100):
			resTx, err := zts.clientCtx.Client.Tx(context.TODO(), hash, false)
			if err == nil {
				resBlock, err := zts.clientCtx.Client.Block(context.TODO(), &resTx.Height)
				if err == nil {
					return mkTxResult(zts.clientCtx.TxConfig, resTx, resBlock)
				}
			}
		}
	}
}

func mkTxResult(txConfig client.TxConfig, resTx *coretypes.ResultTx, resBlock *coretypes.ResultBlock) (*sdktypes.TxResponse, error) {
	txb, err := txConfig.TxDecoder()(resTx.Tx)
	if err != nil {
		return nil, err
	}
	p, ok := txb.(intoAny)
	if !ok {
		return nil, fmt.Errorf("expecting a type implementing intoAny, got: %T", txb)
	}
	any := p.AsAny()
	return sdktypes.NewResponseResultTx(resTx, any, resBlock.Block.Time.Format(time.RFC3339)), nil
}

type intoAny interface {
	AsAny() *codectypes.Any
}

// DeploySystemContractsAndZRC20 deploys the system contracts and ZRC20 contracts
// returns the addresses of uniswap factory, router and erc20 zrc20
func (zts ZetaTxServer) DeploySystemContractsAndZRC20(account, erc20Addr string) (string, string, string, string, string, error) {
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
		return "", "", "", "", "", fmt.Errorf("failed to fetch system contract address: %s; rawlog %s", err.Error(), res.RawLog)
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
		return "", "", "", "", "", fmt.Errorf("failed to fetch zevm connector address: %s, txResponse: %s", err.Error(), res.String())
	}

	// get wzeta address
	wzetaAddr, err := FetchAttributeFromTxResponse(res, "wzeta")
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("failed to fetch wzeta address: %s, txResponse: %s", err.Error(), res.String())
	}

	// deploy eth zrc20
	_, err = zts.BroadcastTx(account, fungibletypes.NewMsgDeployFungibleCoinZRC20(
		addr.String(),
		"",
		chains.GoerliLocalnetChain().ChainId,
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
		chains.BtcRegtestChain().ChainId,
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
		chains.GoerliLocalnetChain().ChainId,
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
