package local

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	zetae2econfig "github.com/zeta-chain/zetacore/cmd/zetae2e/config"
	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/e2etests"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/txserver"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/testutil"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
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
	flagSkipHeaderProof   = "skip-header-proof"
)

var (
	TestTimeout = 15 * time.Minute
)

var noError = testutil.NoError

// NewLocalCmd returns the local command
// which runs the E2E tests locally on the machine with localnet for each blockchain
func NewLocalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "local",
		Short: "Run Local E2E tests",
		Run:   localE2ETest,
	}
	cmd.Flags().Bool(flagContractsDeployed, false, "set to to true if running tests again with existing state")
	cmd.Flags().Int64(flagWaitForHeight, 0, "block height for tests to begin, ex. --wait-for 100")
	cmd.Flags().String(FlagConfigFile, "", "config file to use for the tests")
	cmd.Flags().Bool(flagVerbose, false, "set to true to enable verbose logging")
	cmd.Flags().Bool(flagTestAdmin, false, "set to true to run admin tests")
	cmd.Flags().Bool(flagTestPerformance, false, "set to true to run performance tests")
	cmd.Flags().Bool(flagTestCustom, false, "set to true to run custom tests")
	cmd.Flags().Bool(flagSkipRegular, false, "set to true to skip regular tests")
	cmd.Flags().Bool(flagLight, false, "run the most basic regular tests, useful for quick checks")
	cmd.Flags().Bool(flagSetupOnly, false, "set to true to only setup the networks")
	cmd.Flags().String(flagConfigOut, "", "config file to write the deployed contracts from the setup")
	cmd.Flags().Bool(flagSkipSetup, false, "set to true to skip setup")
	cmd.Flags().Bool(flagSkipBitcoinSetup, false, "set to true to skip bitcoin wallet setup")
	cmd.Flags().Bool(flagSkipHeaderProof, false, "set to true to skip header proof tests")

	return cmd
}

func localE2ETest(cmd *cobra.Command, _ []string) {
	// fetch flags
	var (
		waitForHeight     = must(cmd.Flags().GetInt64(flagWaitForHeight))
		contractsDeployed = must(cmd.Flags().GetBool(flagContractsDeployed))
		verbose           = must(cmd.Flags().GetBool(flagVerbose))
		configOut         = must(cmd.Flags().GetString(flagConfigOut))
		testAdmin         = must(cmd.Flags().GetBool(flagTestAdmin))
		testPerformance   = must(cmd.Flags().GetBool(flagTestPerformance))
		testCustom        = must(cmd.Flags().GetBool(flagTestCustom))
		skipRegular       = must(cmd.Flags().GetBool(flagSkipRegular))
		light             = must(cmd.Flags().GetBool(flagLight))
		setupOnly         = must(cmd.Flags().GetBool(flagSetupOnly))
		skipSetup         = must(cmd.Flags().GetBool(flagSkipSetup))
		skipBitcoinSetup  = must(cmd.Flags().GetBool(flagSkipBitcoinSetup))
		skipHeaderProof   = must(cmd.Flags().GetBool(flagSkipHeaderProof))
	)

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
	noError(err)

	// initialize context
	ctx, cancel := context.WithCancel(context.Background())

	// wait for a specific height on ZetaChain
	if waitForHeight != 0 {
		noError(utils.WaitForBlockHeight(ctx, waitForHeight, conf.RPCs.ZetaCoreRPC, logger))
	}

	// set account prefix to zeta
	setCosmosConfig()

	zetaTxServer, err := txserver.NewZetaTxServer(
		conf.RPCs.ZetaCoreRPC,
		[]string{utils.EmergencyPolicyName, utils.OperationalPolicyName, utils.AdminPolicyName},
		[]string{
			conf.PolicyAccounts.EmergencyPolicyAccount.RawPrivateKey.String(),
			conf.PolicyAccounts.OperationalPolicyAccount.RawPrivateKey.String(),
			conf.PolicyAccounts.AdminPolicyAccount.RawPrivateKey.String(),
		},
		conf.ZetaChainID,
	)
	noError(err)

	// initialize deployer runner with config
	deployerRunner, err := zetae2econfig.RunnerFromConfig(
		ctx,
		"deployer",
		cancel,
		conf,
		conf.DefaultAccount,
		logger,
		runner.WithZetaTxServer(zetaTxServer),
	)
	noError(err)

	// wait for keygen to be completed
	// if setup is skipped, we assume that the keygen is already completed
	if !skipSetup {
		waitKeygenHeight(ctx, deployerRunner.CctxClient, logger)
	}

	// query and set the TSS
	noError(deployerRunner.SetTSSAddresses())

	if !skipHeaderProof {
		noError(deployerRunner.EnableHeaderVerification([]int64{
			chains.GoerliLocalnet.ChainId,
			chains.BitcoinRegtest.ChainId,
		}))
	}

	// setting up the networks
	if !skipSetup {
		logger.Print("⚙️ setting up networks")
		startTime := time.Now()

		deployerRunner.SetupEVM(contractsDeployed, true)
		deployerRunner.SetZEVMContracts()
		noError(deployerRunner.FundEmissionsPool())

		deployerRunner.MintERC20OnEvm(10000)

		logger.Print("✅ setup completed in %s", time.Since(startTime))
	}

	// if a config output is specified, write the config
	if configOut != "" {
		newConfig := zetae2econfig.ExportContractsFromRunner(deployerRunner, conf)

		// write config into stdout
		configOut, err := filepath.Abs(configOut)
		noError(err)

		noError(config.WriteConfig(configOut, newConfig))

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
		// defines all tests, if light is enabled, only the most basic tests are run and advanced are skipped
		erc20Tests := []string{
			e2etests.TestERC20WithdrawName,
			e2etests.TestMultipleERC20WithdrawsName,
			e2etests.TestERC20DepositAndCallRefundName,
			e2etests.TestZRC20SwapName,
		}
		erc20AdvancedTests := []string{
			e2etests.TestERC20DepositRestrictedName,
		}
		zetaTests := []string{
			e2etests.TestZetaWithdrawName,
			e2etests.TestMessagePassingExternalChainsName,
			e2etests.TestMessagePassingRevertFailExternalChainsName,
			e2etests.TestMessagePassingRevertSuccessExternalChainsName,
		}
		zetaAdvancedTests := []string{
			e2etests.TestZetaDepositRestrictedName,
			e2etests.TestZetaDepositName,
			e2etests.TestZetaDepositNewAddressName,
		}
		zevmMPTests := []string{}
		zevmMPAdvancedTests := []string{
			e2etests.TestMessagePassingZEVMToEVMName,
			e2etests.TestMessagePassingEVMtoZEVMName,
			e2etests.TestMessagePassingEVMtoZEVMRevertName,
			e2etests.TestMessagePassingZEVMtoEVMRevertName,
			e2etests.TestMessagePassingZEVMtoEVMRevertFailName,
			e2etests.TestMessagePassingEVMtoZEVMRevertFailName,
		}
		bitcoinTests := []string{
			e2etests.TestBitcoinDepositName,
			e2etests.TestBitcoinDepositRefundName,
			e2etests.TestBitcoinWithdrawSegWitName,
			e2etests.TestBitcoinWithdrawInvalidAddressName,
			e2etests.TestZetaWithdrawBTCRevertName,
			e2etests.TestCrosschainSwapName,
		}
		bitcoinAdvancedTests := []string{
			e2etests.TestBitcoinWithdrawTaprootName,
			e2etests.TestBitcoinWithdrawLegacyName,
			e2etests.TestBitcoinWithdrawMultipleName,
			e2etests.TestBitcoinWithdrawP2SHName,
			e2etests.TestBitcoinWithdrawP2WSHName,
			e2etests.TestBitcoinWithdrawRestrictedName,
		}
		ethereumTests := []string{
			e2etests.TestEtherWithdrawName,
			e2etests.TestContextUpgradeName,
			e2etests.TestEtherDepositAndCallName,
			e2etests.TestEtherDepositAndCallRefundName,
		}
		ethereumAdvancedTests := []string{
			e2etests.TestEtherWithdrawRestrictedName,
		}

		if !light {
			erc20Tests = append(erc20Tests, erc20AdvancedTests...)
			zetaTests = append(zetaTests, zetaAdvancedTests...)
			zevmMPTests = append(zevmMPTests, zevmMPAdvancedTests...)
			bitcoinTests = append(bitcoinTests, bitcoinAdvancedTests...)
			ethereumTests = append(ethereumTests, ethereumAdvancedTests...)
		}

		// skip the header proof test if we run light test or skipHeaderProof is enabled
		testHeader := !light && !skipHeaderProof

		eg.Go(erc20TestRoutine(conf, deployerRunner, verbose, erc20Tests...))
		eg.Go(zetaTestRoutine(conf, deployerRunner, verbose, zetaTests...))
		eg.Go(zevmMPTestRoutine(conf, deployerRunner, verbose, zevmMPTests...))
		eg.Go(bitcoinTestRoutine(conf, deployerRunner, verbose, !skipBitcoinSetup, testHeader, bitcoinTests...))
		eg.Go(ethereumTestRoutine(conf, deployerRunner, verbose, testHeader, ethereumTests...))
	}
	if testAdmin {
		eg.Go(adminTestRoutine(conf, deployerRunner, verbose,
			e2etests.TestRateLimiterName,
			e2etests.TestPauseZRC20Name,
			e2etests.TestUpdateBytecodeZRC20Name,
			e2etests.TestUpdateBytecodeConnectorName,
			e2etests.TestDepositEtherLiquidityCapName,

			// TestMigrateChainSupportName tests EVM chain migration. Currently this test doesn't work with Anvil because pre-EIP1559 txs are not supported
			// See issue below for details
			// TODO: renenable this test as per the issue below
			// https://github.com/zeta-chain/node/issues/1980
			// e2etests.TestMigrateChainSupportName,
		))
	}
	if testPerformance {
		eg.Go(ethereumDepositPerformanceRoutine(conf, deployerRunner, verbose, e2etests.TestStressEtherDepositName))
		eg.Go(ethereumWithdrawPerformanceRoutine(conf, deployerRunner, verbose, e2etests.TestStressEtherWithdrawName))
	}
	if testCustom {
		eg.Go(miscTestRoutine(conf, deployerRunner, verbose, e2etests.TestMyTestName))
	}

	// while tests are executed, monitor blocks in parallel to check if system txs are on top and they have biggest priority
	txPriorityErrCh := make(chan error, 1)
	ctx, monitorPriorityCancel := context.WithCancel(context.Background())
	go monitorTxPriorityInBlocks(ctx, conf, txPriorityErrCh)

	if err := eg.Wait(); err != nil {
		deployerRunner.CtxCancel()
		monitorPriorityCancel()
		logger.Print("❌ %v", err)
		logger.Print("❌ e2e tests failed after %s", time.Since(testStartTime).String())
		os.Exit(1)
	}

	// if all tests pass, cancel txs priority monitoring and check if tx priority is not correct in some blocks
	logger.Print("⏳ e2e tests passed, checking tx priority")
	monitorPriorityCancel()
	if err := <-txPriorityErrCh; err != nil {
		logger.Print("❌ %v", err)
		logger.Print("❌ e2e tests failed after %s", time.Since(testStartTime).String())
		os.Exit(1)
	}

	logger.Print("✅ e2e tests completed in %s", time.Since(testStartTime).String())

	// print and validate report
	networkReport, err := deployerRunner.GenerateNetworkReport()
	if err != nil {
		logger.Print("❌ failed to generate network report %v", err)
	}
	deployerRunner.PrintNetworkReport(networkReport)
	if err := networkReport.Validate(); err != nil {
		logger.Print("❌ network report validation failed %v", err)
		os.Exit(1)
	}

	os.Exit(0)
}

// waitKeygenHeight waits for keygen height
func waitKeygenHeight(
	ctx context.Context,
	cctxClient crosschaintypes.QueryClient,
	logger *runner.Logger,
) {
	// wait for keygen to be completed
	keygenHeight := int64(35)
	logger.Print("⏳ wait height %v for keygen to be completed", keygenHeight)
	for {
		time.Sleep(2 * time.Second)
		response, err := cctxClient.LastZetaHeight(ctx, &crosschaintypes.QueryLastZetaHeightRequest{})
		if err != nil {
			logger.Error("cctxClient.LastZetaHeight error: %s", err)
			continue
		}
		if response.Height >= keygenHeight {
			break
		}
		logger.Info("Last ZetaHeight: %d", response.Height)
	}
}

func must[T any](v T, err error) T {
	return testutil.Must(v, err)
}
