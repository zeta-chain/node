package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
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
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
)

type SmokeTest struct {
	zevmClient   *ethclient.Client
	goerliClient *ethclient.Client
	btcRPCClient *rpcclient.Client
	zetaTxServer ZetaTxServer

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

func NewSmokeTest(
	goerliClient *ethclient.Client,
	zevmClient *ethclient.Client,
	cctxClient crosschaintypes.QueryClient,
	zetaTxServer ZetaTxServer,
	fungibleClient fungibletypes.QueryClient,
	authClient authtypes.QueryClient,
	bankClient banktypes.QueryClient,
	observerClient observertypes.QueryClient,
	goerliAuth *bind.TransactOpts,
	zevmAuth *bind.TransactOpts,
	btcRPCClient *rpcclient.Client,
) *SmokeTest {
	// query system contract address
	systemContractRes, err := fungibleClient.SystemContract(context.Background(), &fungibletypes.QueryGetSystemContractRequest{})
	if err != nil {
		panic(err)
	}

	SystemContract, err := systemcontract.NewSystemContract(HexToAddress(systemContractRes.SystemContract.SystemContract), zevmClient)
	if err != nil {
		panic(err)
	}
	systemContractAddr := HexToAddress(systemContractRes.SystemContract.SystemContract)
	fmt.Printf("System contract address: %s\n", systemContractAddr)

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
		zetaTxServer:       zetaTxServer,
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
		SystemContractAddr: systemContractAddr,
	}
}
