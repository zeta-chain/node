//go:build PRIVNET
// +build PRIVNET

package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/btcsuite/btcutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/contracts/contextapp"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/contracts/zevmswap"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	erc20custody "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	zetaeth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zeta.eth.sol"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"
	systemcontract "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/systemcontract.sol"
	zrc20 "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	uniswapv2factory "github.com/zeta-chain/protocol-contracts/pkg/uniswap/v2-core/contracts/uniswapv2factory.sol"
	uniswapv2router "github.com/zeta-chain/protocol-contracts/pkg/uniswap/v2-periphery/contracts/uniswapv2router02.sol"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/contracts/erc20"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	"google.golang.org/grpc"
)

var (
	DeployerAddress      = ethcommon.HexToAddress("0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC")
	DeployerPrivateKey   = "d87baf7bf6dc560a252596678c12e41f7d1682837f05b29d411bc3f78ae2c263"
	TSSAddress           = ethcommon.HexToAddress("0x0Da38EA1B43758F55eB97590D41e244913A00b26")
	BTCTSSAddress, _     = btcutil.DecodeAddress("bcrt1q78nlhm7mr7t6z8a93z3y93k75ftppcukt5ayay", config.BitconNetParams)
	BLOCK                = 2 * time.Second // should be 2x block time
	BigZero              = big.NewInt(0)
	SmokeTestTimeout     = 24 * time.Hour // smoke test fails if timeout is reached
	USDTZRC20Addr        = "0x48f80608B672DC30DC7e3dbBd0343c5F02C738Eb"
	USDTERC20Addr        = "0xff3135df4F2775f4091b81f4c7B6359CfA07862a"
	ERC20CustodyAddr     = "0xD28D6A0b8189305551a0A8bd247a6ECa9CE781Ca"
	UniswapV2FactoryAddr = "0x9fd96203f7b22bCF72d9DCb40ff98302376cE09c"
	UniswapV2RouterAddr  = "0x2ca7d64A7EFE2D62A725E2B35Cf7230D6677FfEe"
	//SystemContractAddr   = "0x91d18e54DAf4F677cB28167158d6dd21F6aB3921"
	//ZEVMSwapAppAddr      = "0x65a45c57636f9BcCeD4fe193A602008578BcA90b"
	HexToAddress = ethcommon.HexToAddress
)

type SmokeTest struct {
	zevmClient   *ethclient.Client
	goerliClient *ethclient.Client
	btcRPCClient *rpcclient.Client

	cctxClient     types.QueryClient
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

var RootCmd = &cobra.Command{
	Use:   "smoketest",
	Short: "Smoke Test CLI",
}

var LocalCmd = &cobra.Command{
	Use:   "local",
	Short: "Run Local Smoketest",
	Run:   LocalSmokeTest,
}

func init() {
	RootCmd.AddCommand(LocalCmd)
}

func NewSmokeTest(goerliClient *ethclient.Client, zevmClient *ethclient.Client,
	cctxClient types.QueryClient, fungibleClient fungibletypes.QueryClient,
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

	response := &types.QueryGetTssAddressResponse{}
	for {
		response, err = cctxClient.GetTssAddress(context.Background(), &types.QueryGetTssAddressRequest{})
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

func LocalSmokeTest(_ *cobra.Command, _ []string) {
	testStartTime := time.Now()
	defer func() {
		fmt.Println("Smoke test took", time.Since(testStartTime))
	}()
	go func() {
		time.Sleep(SmokeTestTimeout)
		fmt.Println("Smoke test timed out after", SmokeTestTimeout)
		os.Exit(1)
	}()

	connCfg := &rpcclient.ConnConfig{
		Host:         "bitcoin:18443",
		User:         "smoketest",
		Pass:         "123",
		HTTPPostMode: true,
		DisableTLS:   true,
		Params:       "testnet3",
	}
	btcRPCClient, err := rpcclient.New(connCfg, nil)
	if err != nil {
		panic(err)
	}

	goerliClient, err := ethclient.Dial("http://eth:8545")
	if err != nil {
		panic(err)
	}

	bal, err := goerliClient.BalanceAt(context.TODO(), DeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deployer address: %s, balance: %d Ether\n", DeployerAddress.Hex(), bal.Div(bal, big.NewInt(1e18)))

	chainid, err := goerliClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	deployerPrivkey, err := crypto.HexToECDSA(DeployerPrivateKey)
	if err != nil {
		panic(err)
	}
	goerliAuth, err := bind.NewKeyedTransactorWithChainID(deployerPrivkey, chainid)
	if err != nil {
		panic(err)
	}

	grpcConn, err := grpc.Dial("zetacore0:9090", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	cctxClient := types.NewQueryClient(grpcConn)
	fungibleClient := fungibletypes.NewQueryClient(grpcConn)
	authClient := authtypes.NewQueryClient(grpcConn)
	bankClient := banktypes.NewQueryClient(grpcConn)
	observerClient := observertypes.NewQueryClient(grpcConn)

	// Wait for Genesis and keygen to be completed. ~ height 30
	time.Sleep(20 * time.Second)
	for {
		time.Sleep(5 * time.Second)
		response, err := cctxClient.LastZetaHeight(context.Background(), &types.QueryLastZetaHeightRequest{})
		if err != nil {
			fmt.Printf("cctxClient.LastZetaHeight error: %s", err)
			continue
		}
		if response.Height >= 30 {
			break
		}
		fmt.Printf("Last ZetaHeight: %d\n", response.Height)
	}

	// get the clients for tests
	var zevmClient *ethclient.Client
	for {
		time.Sleep(5 * time.Second)
		fmt.Printf("dialing zevm client: http://zetacore0:8545\n")
		zevmClient, err = ethclient.Dial("http://zetacore0:8545")
		if err != nil {
			continue
		}
		break
	}
	chainid, err = zevmClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	zevmAuth, err := bind.NewKeyedTransactorWithChainID(deployerPrivkey, chainid)
	if err != nil {
		panic(err)
	}

	smokeTest := NewSmokeTest(goerliClient, zevmClient, cctxClient, fungibleClient, authClient, bankClient, observerClient, goerliAuth, zevmAuth, btcRPCClient)

	// The following deployment must happen here and in this order, please do not change
	// ==================== Deploying contracts ====================
	startTime := time.Now()

	// test the block of increaseAllowance & decreaseAllowance

	smokeTest.TestBitcoinSetup()
	smokeTest.TestSetupZetaTokenAndConnectorAndZEVMContracts()
	smokeTest.TestDepositEtherIntoZRC20()

	smokeTest.TestSendZetaIn()

	zevmSwapAppAddr, tx, _, err := zevmswap.DeployZEVMSwapApp(smokeTest.zevmAuth, smokeTest.zevmClient, smokeTest.UniswapV2RouterAddr, smokeTest.SystemContractAddr)
	if err != nil {
		panic(err)
	}
	receipt := MustWaitForTxReceipt(zevmClient, tx)
	if receipt.Status != 1 {
		panic("ZEVMSwapApp deployment failed")
	}
	zevmSwapApp, err := zevmswap.NewZEVMSwapApp(zevmSwapAppAddr, zevmClient)
	fmt.Printf("ZEVMSwapApp contract address: %s, tx hash: %s\n", zevmSwapAppAddr.Hex(), tx.Hash().Hex())
	smokeTest.ZEVMSwapAppAddr = zevmSwapAppAddr
	smokeTest.ZEVMSwapApp = zevmSwapApp

	// test system contract context upgrade
	contextAppAddr, tx, _, err := contextapp.DeployContextApp(smokeTest.zevmAuth, smokeTest.zevmClient)
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(zevmClient, tx)
	if receipt.Status != 1 {
		panic("ContextApp deployment failed")
	}
	contextApp, err := contextapp.NewContextApp(contextAppAddr, zevmClient)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ContextApp contract address: %s, tx hash: %s\n", contextAppAddr.Hex(), tx.Hash().Hex())
	smokeTest.ContextAppAddr = contextAppAddr
	smokeTest.ContextApp = contextApp

	fmt.Printf("## Essential tests takes %s\n", time.Since(startTime))
	fmt.Printf("## The DeployerAddress %s is funded on the following networks:\n", DeployerAddress.Hex())
	fmt.Printf("##   Ether on Ethereum private net\n")
	fmt.Printf("##   ZETA on ZetaChain EVM\n")
	fmt.Printf("##   ETH ZRC20 on ZetaChain\n")
	// The following tests are optional tests; comment out the ones you don't want to run
	// temporarily to reduce dev/test cycle turnaround time
	smokeTest.CheckZRC20ReserveAndSupply()

	smokeTest.TestContextUpgrade()

	smokeTest.TestDepositAndCallRefund()
	smokeTest.CheckZRC20ReserveAndSupply()

	smokeTest.TestERC20Deposit()
	smokeTest.CheckZRC20ReserveAndSupply()

	smokeTest.TestERC20Withdraw()
	//smokeTest.WithdrawBitcoinMultipleTimes(5)
	smokeTest.CheckZRC20ReserveAndSupply()

	smokeTest.TestSendZetaOut()
	smokeTest.CheckZRC20ReserveAndSupply()

	smokeTest.TestSendZetaOutBTCRevert()
	smokeTest.CheckZRC20ReserveAndSupply()

	smokeTest.TestMessagePassing()
	smokeTest.CheckZRC20ReserveAndSupply()

	smokeTest.TestZRC20Swap()
	smokeTest.CheckZRC20ReserveAndSupply()

	smokeTest.TestBitcoinWithdraw()
	smokeTest.CheckZRC20ReserveAndSupply()

	smokeTest.TestCrosschainSwap()
	smokeTest.CheckZRC20ReserveAndSupply()

	smokeTest.TestMessagePassingRevertFail()
	smokeTest.CheckZRC20ReserveAndSupply()

	smokeTest.TestMessagePassingRevertSuccess()
	smokeTest.CheckZRC20ReserveAndSupply()

	// add your dev test here
	smokeTest.TestMyTest()

	//{
	//	LoudPrintf("Test ZRC20 blocked methods\n")
	//	sampleZRC20, err := zrc20.NewZRC20(smokeTest.USDTZRC20Addr, zevmClient)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	_, err = sampleZRC20.Approve(smokeTest.zevmAuth, smokeTest.ERC20CustodyAddr, big.NewInt(1000000))
	//	if err != nil {
	//		panic(err)
	//	}
	//	tx, err := sampleZRC20.IncreaseAllowance(smokeTest.zevmAuth, smokeTest.ERC20CustodyAddr, big.NewInt(1000000))
	//	if err != nil {
	//		panic(err)
	//	}
	//	receipt := MustWaitForTxReceipt(zevmClient, tx)
	//	if receipt.Status != 0 {
	//		panic("IncreaseAllowance should be blocked!")
	//	} else {
	//		fmt.Printf("IncreaseAllowance is blocked as expected\n")
	//	}
	//	tx, err = sampleZRC20.DecreaseAllowance(smokeTest.zevmAuth, smokeTest.ERC20CustodyAddr, big.NewInt(1000000))
	//	if err != nil {
	//		panic(err)
	//	}
	//	receipt = MustWaitForTxReceipt(zevmClient, tx)
	//	if receipt.Status != 0 {
	//		panic("IncreaseAllowance should be blocked!")
	//	} else {
	//		fmt.Printf("IncreaseAllowance is blocked as expected\n")
	//	}
	//}

	smokeTest.wg.Wait()
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
