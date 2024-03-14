package local

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	zetae2econfig "github.com/zeta-chain/zetacore/cmd/zetae2e/config"
	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/e2etests"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"golang.org/x/sync/errgroup"
)

const (
	flagContractsDeployed = "deployed"
	flagWaitForHeight     = "wait-for"
	FlagConfigFile        = "config"
	flagConfigOut         = "config-out"
	flagVerbose           = "verbose"
	flagTestAdmin         = "test-admin"
	flagTestPerformance   = "test-performance"
	flagTestCustom        = "test-custom"
	flagSkipRegular       = "skip-regular"
	flagLight             = "light"
	flagSetupOnly         = "setup-only"
	flagSkipSetup         = "skip-setup"
	flagSkipBitcoinSetup  = "skip-bitcoin-setup"
)

var (
	TestTimeout = 15 * time.Minute
)

// NewLocalCmd returns the local command
// which runs the E2E tests locally on the machine with localnet for each blockchain
func NewLocalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "local",
		Short: "Run Local E2E tests",
		Run:   localE2ETest,
	}
	cmd.Flags().Bool(
		flagContractsDeployed,
		false,
		"set to to true if running tests again with existing state",
	)
	cmd.Flags().Int64(
		flagWaitForHeight,
		0,
		"block height for tests to begin, ex. --wait-for 100",
	)
	cmd.Flags().String(
		FlagConfigFile,
		"",
		"config file to use for the tests",
	)
	cmd.Flags().String(
		flagConfigOut,
		"",
		"config file to write the deployed contracts from the setup",
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
		flagTestPerformance,
		false,
		"set to true to run performance tests",
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
	cmd.Flags().Bool(
		flagLight,
		false,
		"run the most basic regular tests, useful for quick checks",
	)
	cmd.Flags().Bool(
		flagSetupOnly,
		false,
		"set to true to only setup the networks",
	)
	cmd.Flags().Bool(
		flagSkipSetup,
		false,
		"set to true to skip setup",
	)
	cmd.Flags().Bool(
		flagSkipBitcoinSetup,
		false,
		"set to true to skip bitcoin wallet setup",
	)

	return cmd
}

func localE2ETest(cmd *cobra.Command, _ []string) {
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
	configOut, err := cmd.Flags().GetString(flagConfigOut)
	if err != nil {
		panic(err)
	}
	testAdmin, err := cmd.Flags().GetBool(flagTestAdmin)
	if err != nil {
		panic(err)
	}
	testPerformance, err := cmd.Flags().GetBool(flagTestPerformance)
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
	light, err := cmd.Flags().GetBool(flagLight)
	if err != nil {
		panic(err)
	}
	setupOnly, err := cmd.Flags().GetBool(flagSetupOnly)
	if err != nil {
		panic(err)
	}
	skipSetup, err := cmd.Flags().GetBool(flagSkipSetup)
	if err != nil {
		panic(err)
	}
	skipBitcoinSetup, err := cmd.Flags().GetBool(flagSkipBitcoinSetup)
	if err != nil {
		panic(err)
	}

	logger := runner.NewLogger(verbose, color.FgWhite, "setup")

	testStartTime := time.Now()
	logger.Print("starting E2E tests")

	if testAdmin {
		logger.Print("⚠️ admin tests enabled")
	}

	if testPerformance {
		logger.Print("⚠️ performance tests enabled, regular tests will be skipped")
		skipRegular = true
	}

	// start timer
	go func() {
		time.Sleep(TestTimeout)
		logger.Error("Test timed out after %s", TestTimeout.String())
		os.Exit(1)
	}()

	// initialize tests config
	conf, err := GetConfig(cmd)
	if err != nil {
		panic(err)
	}

	// initialize context
	ctx, cancel := context.WithCancel(context.Background())

	// wait for a specific height on ZetaChain
	if waitForHeight != 0 {
		utils.WaitForBlockHeight(ctx, waitForHeight, conf.RPCs.ZetaCoreRPC, logger)
	}

	// set account prefix to zeta
	setCosmosConfig()

	// wait for Genesis
	// if setup is skip, we assume that the genesis is already created
	if !skipSetup {
		logger.Print("⏳ wait 70s for genesis")
		time.Sleep(70 * time.Second)
	}

	// initialize deployer runner with config
	deployerRunner, err := zetae2econfig.RunnerFromConfig(
		ctx,
		"deployer",
		cancel,
		conf,
		DeployerAddress,
		DeployerPrivateKey,
		utils.FungibleAdminName,
		FungibleAdminMnemonic,
		logger,
	)
	if err != nil {
		panic(err)
	}

	// wait for keygen to be completed
	// if setup is skipped, we assume that the keygen is already completed
	if !skipSetup {
		waitKeygenHeight(ctx, deployerRunner.CctxClient, logger)
	}

	// query and set the TSS
	if err := deployerRunner.SetTSSAddresses(); err != nil {
		panic(err)
	}

	// setting up the networks
	if !skipSetup {
		logger.Print("⚙️ setting up networks")
		startTime := time.Now()
		deployerRunner.SetupEVM(contractsDeployed)
		deployerRunner.SetZEVMContracts()
		deployerRunner.MintERC20OnEvm(10000)
		logger.Print("✅ setup completed in %s", time.Since(startTime))
	}

	// if a config output is specified, write the config
	if configOut != "" {
		newConfig := zetae2econfig.ExportContractsFromRunner(deployerRunner, conf)
		configOut, err := filepath.Abs(configOut)
		if err != nil {
			panic(err)
		}

		// write config into stdout
		if err := config.WriteConfig(configOut, newConfig); err != nil {
			panic(err)
		}

		logger.Print("✅ config file written in %s", configOut)
	}

	deployerRunner.PrintContractAddresses()

	// if setup only, quit
	if setupOnly {
		logger.Print("✅ the localnet has been setup")
		os.Exit(0)
	}

	// run tests
	var eg errgroup.Group
	if !skipRegular {
		// defines all tests, if light is enabled, only the most basic tests are run
		erc20Tests := []string{
			e2etests.TestERC20WithdrawName,
			e2etests.TestMultipleWithdrawsName,
			e2etests.TestERC20DepositAndCallRefundName,
			e2etests.TestZRC20SwapName,
		}
		erc20AdvancedTests := []string{
			e2etests.TestERC20DepositRestrictedName,
		}
		zetaTests := []string{
			e2etests.TestZetaWithdrawName,
			e2etests.TestMessagePassingName,
			e2etests.TestMessagePassingRevertFailName,
			e2etests.TestMessagePassingRevertSuccessName,
		}
		zetaAdvancedTests := []string{
			e2etests.TestZetaDepositRestrictedName,
		}
		bitcoinTests := []string{
			e2etests.TestBitcoinWithdrawName,
			e2etests.TestBitcoinWithdrawInvalidAddressName,
			e2etests.TestZetaWithdrawBTCRevertName,
			e2etests.TestCrosschainSwapName,
		}
		bitcoinAdvancedTests := []string{
			e2etests.TestBitcoinWithdrawRestrictedName,
		}
		ethereumTests := []string{
			e2etests.TestEtherWithdrawName,
			e2etests.TestContextUpgradeName,
			e2etests.TestEtherDepositAndCallName,
			e2etests.TestDepositAndCallRefundName,
		}
		ethereumAdvancedTests := []string{
			e2etests.TestEtherWithdrawRestrictedName,
		}

		if !light {
			erc20Tests = append(erc20Tests, erc20AdvancedTests...)
			zetaTests = append(zetaTests, zetaAdvancedTests...)
			bitcoinTests = append(bitcoinTests, bitcoinAdvancedTests...)
			ethereumTests = append(ethereumTests, ethereumAdvancedTests...)
		}

		eg.Go(erc20TestRoutine(conf, deployerRunner, verbose, erc20Tests...))
		eg.Go(zetaTestRoutine(conf, deployerRunner, verbose, zetaTests...))
		eg.Go(bitcoinTestRoutine(conf, deployerRunner, verbose, !skipBitcoinSetup, !light, bitcoinTests...))
		eg.Go(ethereumTestRoutine(conf, deployerRunner, verbose, !light, ethereumTests...))
	}
	if testAdmin {
		eg.Go(adminTestRoutine(conf, deployerRunner, verbose,
			e2etests.TestPauseZRC20Name,
			e2etests.TestUpdateBytecodeName,
			e2etests.TestDepositEtherLiquidityCapName,
		))
	}
	if testPerformance {
		eg.Go(ethereumDepositPerformanceRoutine(conf, deployerRunner, verbose, e2etests.TestStressEtherDepositName))
		eg.Go(ethereumWithdrawPerformanceRoutine(conf, deployerRunner, verbose, e2etests.TestStressEtherWithdrawName))
	}
	if testCustom {
		eg.Go(miscTestRoutine(conf, deployerRunner, verbose, e2etests.TestMyTestName))
	}

	if err := eg.Wait(); err != nil {
		deployerRunner.CtxCancel()
		logger.Print("❌ %v", err)
		logger.Print("❌ e2e tests failed after %s", time.Since(testStartTime).String())
		os.Exit(1)
	}

	logger.Print("✅ e2e tests completed in %s", time.Since(testStartTime).String())
}
