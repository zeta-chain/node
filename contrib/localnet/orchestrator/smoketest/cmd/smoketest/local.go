package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/btcsuite/btcd/rpcclient"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/app"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/config"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/smoketests"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/txserver"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc"
)

const (
	flagContractsDeployed = "deployed"
	flagWaitForHeight     = "wait-for"
	flagConfigFile        = "config"
)

func NewLocalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "local",
		Short: "Run Local Smoketest",
		Run:   localSmokeTest,
	}
	cmd.Flags().Bool(
		flagContractsDeployed,
		false,
		"set to to true if running smoketest again with existing state",
	)
	cmd.Flags().Int64(
		flagWaitForHeight,
		0,
		"block height for smoketest to begin, ex. --wait-for 100",
	)
	cmd.Flags().String(
		flagConfigFile,
		"",
		"config file to use for the smoketest",
	)
	return cmd
}

func localSmokeTest(cmd *cobra.Command, _ []string) {
	testStartTime := time.Now()
	defer func() {
		fmt.Println("Smoke test took", time.Since(testStartTime))
	}()

	// start timer
	go func() {
		time.Sleep(SmokeTestTimeout)
		fmt.Println("Smoke test timed out after", SmokeTestTimeout)
		os.Exit(1)
	}()

	// fetch flags
	waitForHeight, err := cmd.Flags().GetInt64(flagWaitForHeight)
	if err != nil {
		panic(err)
	}
	contractsDeployed, err := cmd.Flags().GetBool(flagContractsDeployed)
	if err != nil {
		panic(err)
	}

	// initialize smoke tests config
	conf, err := getConfig(cmd)
	if err != nil {
		panic(err)
	}

	// wait for a specific height on ZetaChain
	if waitForHeight != 0 {
		utils.WaitForBlockHeight(waitForHeight, conf.RPCs.ZetaCoreRPC)
	}

	// set account prefix to zeta
	cosmosConf := sdk.GetConfig()
	cosmosConf.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	cosmosConf.Seal()

	// initialize clients
	// TODO: add connection values to config
	// https://github.com/zeta-chain/node-private/issues/41
	connCfg := &rpcclient.ConnConfig{
		Host:         conf.RPCs.Bitcoin,
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

	goerliClient, err := ethclient.Dial(conf.RPCs.EVM)
	if err != nil {
		panic(err)
	}

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

	grpcConn, err := grpc.Dial(conf.RPCs.ZetaCoreGRPC, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	cctxClient := crosschaintypes.NewQueryClient(grpcConn)
	fungibleClient := fungibletypes.NewQueryClient(grpcConn)
	authClient := authtypes.NewQueryClient(grpcConn)
	bankClient := banktypes.NewQueryClient(grpcConn)
	observerClient := observertypes.NewQueryClient(grpcConn)

	// wait for Genesis
	time.Sleep(30 * time.Second)

	// initialize client to send messages to ZetaChain
	zetaTxServer, err := txserver.NewZetaTxServer(
		conf.RPCs.ZetaCoreRPC,
		[]string{utils.FungibleAdminName},
		[]string{FungibleAdminMnemonic},
		conf.ZetaChainID,
	)
	if err != nil {
		panic(err)
	}

	// wait for keygen to be completed. ~ height 30
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

	// setup client and auth for zevm
	var zevmClient *ethclient.Client
	for {
		time.Sleep(5 * time.Second)
		fmt.Printf("dialing zevm client: %s\n", conf.RPCs.Zevm)
		zevmClient, err = ethclient.Dial(conf.RPCs.Zevm)
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

	// initialize smoke test runner
	sm := runner.NewSmokeTestRunner(
		DeployerAddress,
		DeployerPrivateKey,
		FungibleAdminMnemonic,
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

	// setting up the networks
	startTime := time.Now()

	// setup TSS addresses
	sm.SetTSSAddresses()

	// setup the external network
	sm.SetupBitcoin()
	sm.SetupEVM(contractsDeployed)

	// deploy and set zevm contract
	sm.SetZEVMContracts()

	// deposit on ZetaChain
	sm.DepositEtherIntoZRC20()
	sm.SendZetaIn()
	sm.DepositBTC()

	// deploy zevm swap and context apps
	sm.SetupZEVMSwapApp()
	sm.SetupContextApp()

	fmt.Printf("## Setup takes %s\n", time.Since(startTime))

	// run all smoke tests
	sm.RunSmokeTests(smoketests.AllSmokeTests)

	sm.WG.Wait()
}

func getConfig(cmd *cobra.Command) (config.Config, error) {
	configFile, err := cmd.Flags().GetString(flagConfigFile)
	if err != nil {
		return config.Config{}, err
	}

	// use default config if no config file is specified
	if configFile == "" {
		return config.DefaultConfig(), nil
	}

	return config.ReadConfig(configFile)
}
