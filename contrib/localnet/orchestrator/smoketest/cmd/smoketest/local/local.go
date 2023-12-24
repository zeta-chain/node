package local

import (
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	"golang.org/x/sync/errgroup"
)

const (
	flagContractsDeployed = "deployed"
	flagWaitForHeight     = "wait-for"
	flagConfigFile        = "config"
	flagVerbose           = "verbose"
	flagTestAdmin         = "test-admin"
	flagTestCustom        = "test-custom"
	flagSkipRegular       = "skip-regular"
)

var (
	SmokeTestTimeout = 30 * time.Minute
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
		false,
		"set to true to enable verbose logging",
	)
	cmd.Flags().Bool(
		flagTestAdmin,
		false,
		"set to true to run admin tests",
	)
	cmd.Flags().Bool(
		flagTestCustom,
		false,
		"set to true to run custom tests",
	)
	cmd.Flags().Bool(
		flagSkipRegular,
		false,
		"set to true to skip regular tests",
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
	logger := runner.NewLogger(verbose, color.FgWhite, "setup")
	testAdmin, err := cmd.Flags().GetBool(flagTestAdmin)
	if err != nil {
		panic(err)
	}
	testCustom, err := cmd.Flags().GetBool(flagTestCustom)
	if err != nil {
		panic(err)
	}
	skipRegular, err := cmd.Flags().GetBool(flagSkipRegular)
	if err != nil {
		panic(err)
	}

	testStartTime := time.Now()
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
	logger.Print("⏳ wait 60s for genesis")
	time.Sleep(60 * time.Second)

	// initialize deployer runner with config
	deployerRunner, err := runnerFromConfig(conf, DeployerAddress, DeployerPrivateKey, logger)
	if err != nil {
		panic(err)
	}

	// wait for keygen to be completed
	waitKeygenHeight(deployerRunner.CctxClient, logger)

	// setting up the networks
	logger.Print("⚙️ setting up networks")
	startTime := time.Now()
	deployerRunner.SetTSSAddresses()
	deployerRunner.SetupEVM(contractsDeployed)
	deployerRunner.SetZEVMContracts()
	deployerRunner.MintUSDTOnEvm(10000)
	logger.Print("✅ setup completed in %s", time.Since(startTime))

	// run tests
	var eg errgroup.Group
	if !skipRegular {
		eg.Go(erc20TestRoutine(conf, deployerRunner, verbose))
		eg.Go(zetaTestRoutine(conf, deployerRunner, verbose))
		eg.Go(bitcoinTestRoutine(conf, deployerRunner, verbose))
		eg.Go(ethereumTestRoutine(conf, deployerRunner, verbose))
	}
	if testAdmin {
		eg.Go(adminTestRoutine(conf, deployerRunner, verbose))
	}
	if testCustom {
		eg.Go(miscTestRoutine(conf, deployerRunner, verbose))
	}

	if err := eg.Wait(); err != nil {
		logger.Print("❌ %v", err)
		logger.Print("❌ smoke tests failed after %s", time.Since(testStartTime).String())
		os.Exit(1)
	}

	logger.Print("✅ smoke tests completed in %s", time.Since(testStartTime).String())
}
