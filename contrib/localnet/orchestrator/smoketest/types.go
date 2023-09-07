//go:build PRIVNET
// +build PRIVNET

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
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
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	zetaeth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zeta.eth.sol"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/systemcontract.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/uniswap/v2-core/contracts/uniswapv2factory.sol"
	uniswapv2router "github.com/zeta-chain/protocol-contracts/pkg/uniswap/v2-periphery/contracts/uniswapv2router02.sol"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/contracts/contextapp"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/contracts/erc20"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/contracts/zevmswap"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	emissionstypes "github.com/zeta-chain/zetacore/x/emissions/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
)

type SmokeTest struct {
	zevmClient   *ethclient.Client
	goerliClient *ethclient.Client
	btcRPCClient *rpcclient.Client

	cctxClient     crosschaintypes.QueryClient
	fungibleClient fungibletypes.QueryClient
	authClient     authtypes.QueryClient
	bankClient     banktypes.QueryClient
	observerClient observertypes.QueryClient

	wg               sync.WaitGroup
	ZetaEth          *zetaeth.ZetaEth
	ZetaEthAddr      ethcommon.Address
	ConnectorEth     *zetaconnectoreth.ZetaConnectorEth
	ConnectorEthAddr ethcommon.Address
	goerliAuth       *bind.TransactOpts
	zevmAuth         *bind.TransactOpts

	ERC20CustodyAddr     ethcommon.Address
	ERC20Custody         *erc20custody.ERC20Custody
	USDTERC20Addr        ethcommon.Address
	USDTERC20            *erc20.USDT
	USDTZRC20Addr        ethcommon.Address
	USDTZRC20            *zrc20.ZRC20
	ETHZRC20Addr         ethcommon.Address
	ETHZRC20             *zrc20.ZRC20
	BTCZRC20Addr         ethcommon.Address
	BTCZRC20             *zrc20.ZRC20
	UniswapV2FactoryAddr ethcommon.Address
	UniswapV2Factory     *uniswapv2factory.UniswapV2Factory
	UniswapV2RouterAddr  ethcommon.Address
	UniswapV2Router      *uniswapv2router.UniswapV2Router02
	TestDAppAddr         ethcommon.Address
	ZEVMSwapAppAddr      ethcommon.Address
	ZEVMSwapApp          *zevmswap.ZEVMSwapApp
	ContextAppAddr       ethcommon.Address
	ContextApp           *contextapp.ContextApp

	SystemContract     *systemcontract.SystemContract
	SystemContractAddr ethcommon.Address
}

func NewSmokeTest(goerliClient *ethclient.Client, zevmClient *ethclient.Client,
	cctxClient crosschaintypes.QueryClient, fungibleClient fungibletypes.QueryClient,
	authClient authtypes.QueryClient, bankClient banktypes.QueryClient, observerClient observertypes.QueryClient,
	goerliAuth *bind.TransactOpts, zevmAuth *bind.TransactOpts,
	btcRPCClient *rpcclient.Client) *SmokeTest {
	// query system contract address
	systemContractAddr, err := fungibleClient.SystemContract(context.Background(), &fungibletypes.QueryGetSystemContractRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("System contract address: %s\n", systemContractAddr)

	SystemContract, err := systemcontract.NewSystemContract(HexToAddress(systemContractAddr.SystemContract.SystemContract), zevmClient)
	if err != nil {
		panic(err)
	}
	SystemContractAddr := HexToAddress(systemContractAddr.SystemContract.SystemContract)

	response := &crosschaintypes.QueryGetTssAddressResponse{}
	for {
		response, err = cctxClient.GetTssAddress(context.Background(), &crosschaintypes.QueryGetTssAddressRequest{})
		if err != nil {
			fmt.Printf("cctxClient.TSS error %s\n", err.Error())
			fmt.Printf("TSS not ready yet, waiting for TSS to be appear in zetacore netowrk...\n")
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	TSSAddress = ethcommon.HexToAddress(response.Eth)
	BTCTSSAddress, _ = btcutil.DecodeAddress(response.Btc, config.BitconNetParams)
	fmt.Printf("TSS EthAddress: %s\n TSS BTC address %s\n", response.GetEth(), response.GetBtc())

	return &SmokeTest{
		zevmClient:         zevmClient,
		goerliClient:       goerliClient,
		cctxClient:         cctxClient,
		fungibleClient:     fungibleClient,
		authClient:         authClient,
		bankClient:         bankClient,
		observerClient:     observerClient,
		wg:                 sync.WaitGroup{},
		goerliAuth:         goerliAuth,
		zevmAuth:           zevmAuth,
		btcRPCClient:       btcRPCClient,
		SystemContract:     SystemContract,
		SystemContractAddr: SystemContractAddr,
	}
}

// ZetaTxServer is a ZetaChain tx server for smoke test
type ZetaTxServer struct {
	clientCtx client.Context
	txFactory tx.Factory
}

// NewTxServer returns a new TxServer with provided account
func NewTxServer(names []string, mnemonics []string) (ZetaTxServer, error) {
	ctx := context.Background()
	rpcAddr := "zetacore0:26657"

	if len(names) != len(mnemonics) {
		return ZetaTxServer{}, errors.New("invalid names and mnemonics")
	}

	// initialize rpc and check status
	rpc, err := rpchttp.New(rpcAddr, "/websocket")
	if err != nil {
		return ZetaTxServer{}, err
	}
	if _, err = rpc.Status(ctx); err != nil {
		return ZetaTxServer{}, err
	}

	// initialize codec
	cdc, reg := newCodec()

	// initialize keyring
	kr := keyring.NewInMemory(cdc)

	// create accounts
	for i := range names {
		_, err = kr.NewAccount(names[i], mnemonics[i], "", sdktypes.FullFundraiserPath, hd.Secp256k1)
		if err != nil {
			return ZetaTxServer{}, err
		}
	}

	clientCtx := newContext(rpc, cdc, reg, kr)
	txf := newFactory(clientCtx)

	return ZetaTxServer{
		clientCtx: clientCtx,
		txFactory: txf,
	}, nil
}

// BroadcastTx broadcasts a tx to ZetaChain with the provided msg
func (zts ZetaTxServer) BroadcastTx(account string, msg sdktypes.Msg) (*sdktypes.TxResponse, error) {
	// Set account number and sequence number
	// txf := txf.WithAccountNumber(n).WithSequence(m)

	// Set the fees
	// txf = txf.WithFees(fees)

	// Set the gas prices
	// txf = txf.WithGasPrices(gasPrices)

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

// newCodec returns the codec for msg server
func newCodec() (*codec.ProtoCodec, codectypes.InterfaceRegistry) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)

	authtypes.RegisterInterfaces(interfaceRegistry)
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	sdktypes.RegisterInterfaces(interfaceRegistry)
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	crosschaintypes.RegisterInterfaces(interfaceRegistry)
	emissionstypes.RegisterInterfaces(interfaceRegistry)
	fungibletypes.RegisterInterfaces(interfaceRegistry)
	observertypes.RegisterInterfaces(interfaceRegistry)

	return cdc, interfaceRegistry
}

// newContext returns the client context for msg server
func newContext(rpc *rpchttp.HTTP, cdc *codec.ProtoCodec, reg codectypes.InterfaceRegistry, kr keyring.Keyring) client.Context {
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)
	return client.Context{}.
		WithChainID(ZetaChainID).
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
		WithKeyring(kr)
	//WithAccountRetriever(accountRetriever)
	//WithGenerateOnly(false)
}

// newFactory returns the tx factory for msg server
func newFactory(clientCtx client.Context) tx.Factory {
	return tx.Factory{}.
		WithChainID(clientCtx.ChainID).
		WithKeybase(clientCtx.Keyring).
		WithGas(300000).
		WithGasAdjustment(1.0).
		WithSignMode(signing.SignMode_SIGN_MODE_UNSPECIFIED).
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithTxConfig(clientCtx.TxConfig)
}
