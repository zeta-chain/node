//go:build PRIVNET
// +build PRIVNET

package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/app"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/contracts/contextapp"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/contracts/zevmswap"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"google.golang.org/grpc"
)

var (
	ZetaChainID          = "athens_101-1"
	DeployerAddress      = ethcommon.HexToAddress("0xE5C5367B8224807Ac2207d350E60e1b6F27a7ecC")
	DeployerPrivateKey   = "d87baf7bf6dc560a252596678c12e41f7d1682837f05b29d411bc3f78ae2c263"
	TSSAddress           = ethcommon.HexToAddress("0x0Da38EA1B43758F55eB97590D41e244913A00b26")
	BTCTSSAddress, _     = btcutil.DecodeAddress("bcrt1q78nlhm7mr7t6z8a93z3y93k75ftppcukt5ayay", config.BitconNetParams)
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

	// FungibleAdminMnemonic is the mnemonic for the admin account of the fungible module
	//nolint:gosec - disable nosec because this is a test account
	FungibleAdminMnemonic = "snow grace federal cupboard arrive fancy gym lady uniform rotate exercise either leave alien grass"
	FungibleAdminName     = "fungibleadmin"
	FungibleAdminAddress  = "zeta1srsq755t654agc0grpxj4y3w0znktrpr9tcdgk"
)

var RootCmd = &cobra.Command{
	Use:   "smoketest",
	Short: "Smoke Test CLI",
}

var LocalCmd = &cobra.Command{
	Use:   "local",
	Short: "Run Local Smoketest",
	Run:   LocalSmokeTest,
}

type localArgs struct {
	contractsDeployed bool
	waitForHeight     int64
}

var localTestArgs = localArgs{}

func init() {
	RootCmd.AddCommand(LocalCmd)
	LocalCmd.Flags().BoolVar(&localTestArgs.contractsDeployed, "deployed", false, "set to to true if running smoketest again with existing state")
	LocalCmd.Flags().Int64Var(&localTestArgs.waitForHeight, "wait-for", 0, "block height for smoketest to begin, ex. --wait-for 100")
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
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

	if localTestArgs.waitForHeight != 0 {
		WaitForBlockHeight(localTestArgs.waitForHeight)
	}

	// set account prefix to zeta
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	cfg.Seal()

	// initialize clients
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

	cctxClient := crosschaintypes.NewQueryClient(grpcConn)
	fungibleClient := fungibletypes.NewQueryClient(grpcConn)
	authClient := authtypes.NewQueryClient(grpcConn)
	bankClient := banktypes.NewQueryClient(grpcConn)
	observerClient := observertypes.NewQueryClient(grpcConn)

	//Wait for Genesis
	time.Sleep(30 * time.Second)

	// initialize client to send messages to ZetaChain
	zetaTxServer, err := NewZetaTxServer(
		"http://zetacore0:26657",
		[]string{FungibleAdminName},
		[]string{FungibleAdminMnemonic},
	)
	if err != nil {
		panic(err)
	}

	//Wait for keygen to be completed. ~ height 30
	for {
		time.Sleep(5 * time.Second)
		response, err := cctxClient.LastZetaHeight(context.Background(), &crosschaintypes.QueryLastZetaHeightRequest{})
		if err != nil {
			fmt.Printf("cctxClient.LastZetaHeight error: %s", err)
			continue
		}
		if response.Height >= 60 {
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

	smokeTest := NewSmokeTest(
		goerliClient,
		zevmClient,
		cctxClient,
		zetaTxServer,
		fungibleClient,
		authClient,
		bankClient,
		observerClient,
		goerliAuth,
		zevmAuth,
		btcRPCClient,
	)

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

	//test system contract context upgrade
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

	smokeTest.TestPauseZRC20()
	smokeTest.CheckZRC20ReserveAndSupply()

	smokeTest.TestERC20DepositAndCallRefund()
	smokeTest.CheckZRC20ReserveAndSupply()

	smokeTest.TestUpdateBytecode()
	smokeTest.CheckZRC20ReserveAndSupply()

	smokeTest.TestEtherDepositAndCall()
	smokeTest.CheckZRC20ReserveAndSupply()

	smokeTest.TestDepositEtherLiquidityCap()
	smokeTest.CheckZRC20ReserveAndSupply()

	// add your dev test here
	smokeTest.TestMyTest()

	smokeTest.wg.Wait()
}
