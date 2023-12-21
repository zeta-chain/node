package main

import (
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/config"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
)

const (
	flagContractsDeployed = "deployed"
	flagWaitForHeight     = "wait-for"
	flagConfigFile        = "config"
	flagVerbose           = "verbose"
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
	cmd.Flags().Bool(
		flagVerbose,
		true,
		"set to true to enable verbose logging",
	)
	return cmd
}

func localSmokeTest(cmd *cobra.Command, _ []string) {
	// fetch flags
	waitForHeight, err := cmd.Flags().GetInt64(flagWaitForHeight)
	if err != nil {
		panic(err)
	}
	contractsDeployed, err := cmd.Flags().GetBool(flagContractsDeployed)
	if err != nil {
		panic(err)
	}
	verbose, err := cmd.Flags().GetBool(flagVerbose)
	if err != nil {
		panic(err)
	}
	logger := runner.NewLogger(verbose)

	testStartTime := time.Now()
	defer func() {
		logger.Print("✅ smoke tests completed in %s", time.Since(testStartTime).String())
	}()

	logger.Print("starting smoke tests")

	// start timer
	go func() {
		time.Sleep(SmokeTestTimeout)
		logger.Error("Smoke test timed out after %s", SmokeTestTimeout.String())
		os.Exit(1)
	}()

	// initialize smoke tests config
	conf, err := getConfig(cmd)
	if err != nil {
		panic(err)
	}

	// wait for a specific height on ZetaChain
	if waitForHeight != 0 {
		utils.WaitForBlockHeight(waitForHeight, conf.RPCs.ZetaCoreRPC, logger)
	}

	// set account prefix to zeta
	setCosmosConfig()

	// wait for Genesis
	logger.Print("⏳ wait 40s for genesis")
	time.Sleep(40 * time.Second)

	// initialize runner with config
	sm, err := runnerFromConfig(conf, DeployerAddress, DeployerPrivateKey, logger)
	if err != nil {
		panic(err)
	}

	// wait for keygen to be completed
	waitKeygenHeight(sm.CctxClient, logger)

	// setting up the networks
	startTime := time.Now()

	// setup TSS addresses
	logger.Print("⚙️ setting up TSS address")
	sm.SetTSSAddresses()

	// setup the external network
	logger.Print("⚙️ setting up Bitcoin network")
	sm.SetupBitcoin()
	logger.Print("⚙️ setting up Goerli network")
	sm.SetupEVM(contractsDeployed)

	// deploy and set zevm contract
	logger.Print("⚙️ deploying system contracts and ZRC20s on ZEVM")
	sm.SetZEVMContracts()

	// deposits on ZetaChain
	//sm.DepositEther()
	sm.DepositZeta()
	//sm.DepositBTC()
	//sm.DepositERC20()

	// deploy zevm swap and context apps
	//logger.Print("⚙️ setting up ZEVM swap and context apps")
	//sm.SetupZEVMSwapApp()
	//sm.SetupContextApp()

	logger.Print("✅ setup completed in %s", time.Since(startTime))

	// run all smoke tests
	//sm.RunSmokeTests(smoketests.AllSmokeTests)

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
