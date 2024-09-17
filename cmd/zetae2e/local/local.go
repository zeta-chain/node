package local

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/zeta-chain/node/app"
	zetae2econfig "github.com/zeta-chain/node/cmd/zetae2e/config"
	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/txserver"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
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
	flagTestSolana        = "test-solana"
	flagTestTON           = "test-ton"
	flagSkipRegular       = "skip-regular"
	flagLight             = "light"
	flagSetupOnly         = "setup-only"
	flagSkipSetup         = "skip-setup"
	flagTestTSSMigration  = "test-tss-migration"
	flagSkipBitcoinSetup  = "skip-bitcoin-setup"
	flagSkipHeaderProof   = "skip-header-proof"
	flagTestV2            = "test-v2"
	flagTestV2Migration   = "test-v2-migration"
	flagSkipTrackerCheck  = "skip-tracker-check"
	flagSkipPrecompiles   = "skip-precompiles"
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
	cmd.Flags().Bool(flagTestSolana, false, "set to true to run solana tests")
	cmd.Flags().Bool(flagTestTON, false, "set to true to run TON tests")
	cmd.Flags().Bool(flagSkipRegular, false, "set to true to skip regular tests")
	cmd.Flags().Bool(flagLight, false, "run the most basic regular tests, useful for quick checks")
	cmd.Flags().Bool(flagSetupOnly, false, "set to true to only setup the networks")
	cmd.Flags().String(flagConfigOut, "", "config file to write the deployed contracts from the setup")
	cmd.Flags().Bool(flagSkipSetup, false, "set to true to skip setup")
	cmd.Flags().Bool(flagSkipBitcoinSetup, false, "set to true to skip bitcoin wallet setup")
	cmd.Flags().Bool(flagSkipHeaderProof, false, "set to true to skip header proof tests")
	cmd.Flags().Bool(flagTestTSSMigration, false, "set to true to include a migration test at the end")
	cmd.Flags().Bool(flagTestV2, false, "set to true to run tests for v2 contracts")
	cmd.Flags().Bool(flagTestV2Migration, false, "set to true to run tests for v2 contracts migration test")
	cmd.Flags().Bool(flagSkipTrackerCheck, false, "set to true to skip tracker check at the end of the tests")
	cmd.Flags().Bool(flagSkipPrecompiles, false, "set to true to skip stateful precompiled contracts test")

	return cmd
}

// TODO: simplify this file: put the different type of tests in separate files
// https://github.com/zeta-chain/node/issues/2762
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
		testSolana        = must(cmd.Flags().GetBool(flagTestSolana))
		testTON           = must(cmd.Flags().GetBool(flagTestTON))
		skipRegular       = must(cmd.Flags().GetBool(flagSkipRegular))
		light             = must(cmd.Flags().GetBool(flagLight))
		setupOnly         = must(cmd.Flags().GetBool(flagSetupOnly))
		skipSetup         = must(cmd.Flags().GetBool(flagSkipSetup))
		skipBitcoinSetup  = must(cmd.Flags().GetBool(flagSkipBitcoinSetup))
		skipHeaderProof   = must(cmd.Flags().GetBool(flagSkipHeaderProof))
		skipTrackerCheck  = must(cmd.Flags().GetBool(flagSkipTrackerCheck))
		testTSSMigration  = must(cmd.Flags().GetBool(flagTestTSSMigration))
		testV2            = must(cmd.Flags().GetBool(flagTestV2))
		testV2Migration   = must(cmd.Flags().GetBool(flagTestV2Migration))
		skipPrecompiles   = must(cmd.Flags().GetBool(flagSkipPrecompiles))
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

	// temporary spaghetti to overcome e2e flags limitations
	if !testTON {
		conf.RPCs.TONSidecarURL = ""
	}

	// initialize context
	ctx, cancel := context.WithCancel(context.Background())

	// wait for a specific height on ZetaChain
	if waitForHeight != 0 {
		noError(utils.WaitForBlockHeight(ctx, waitForHeight, conf.RPCs.ZetaCoreRPC, logger))
	}

	app.SetConfig()

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

	// set the authority client to the zeta tx server to be able to query message permissions
	deployerRunner.ZetaTxServer.SetAuthorityClient(deployerRunner.AuthorityClient)

	// wait for keygen to be completed
	// if setup is skipped, we assume that the keygen is already completed
	if !skipSetup {
		waitKeygenHeight(ctx, deployerRunner.CctxClient, deployerRunner.ObserverClient, logger, 10)
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

		if testV2 {
			deployerRunner.SetupEVMV2()
		}

		deployerRunner.SetZEVMSystemContracts()

		if testV2 {
			// NOTE: v2 (gateway) setup called here because system contract needs to be set first, then gateway, then zrc20
			deployerRunner.SetZEVMContractsV2()
		}

		deployerRunner.SetZEVMZRC20s()

		if testSolana {
			deployerRunner.SetSolanaContracts(conf.AdditionalAccounts.UserSolana.SolanaPrivateKey.String())
		}
		noError(deployerRunner.FundEmissionsPool())

		deployerRunner.MintERC20OnEvm(1000000)

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

	// run the v2 migration
	if testV2Migration {
		deployerRunner.RunV2Migration()
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
		precompiledContractTests := []string{}

		if !skipPrecompiles {
			precompiledContractTests = []string{
				e2etests.TestPrecompilesPrototypeName,
				e2etests.TestPrecompilesPrototypeThroughContractName,
				e2etests.TestPrecompilesStakingName,
				e2etests.TestPrecompilesStakingThroughContractName,
			}
		}

		if !light {
			erc20Tests = append(erc20Tests, erc20AdvancedTests...)
			zetaTests = append(zetaTests, zetaAdvancedTests...)
			zevmMPTests = append(zevmMPTests, zevmMPAdvancedTests...)
			bitcoinTests = append(bitcoinTests, bitcoinAdvancedTests...)
			ethereumTests = append(ethereumTests, ethereumAdvancedTests...)
		}

		eg.Go(statefulPrecompilesTestRoutine(conf, deployerRunner, verbose, precompiledContractTests...))
		eg.Go(erc20TestRoutine(conf, deployerRunner, verbose, erc20Tests...))
		eg.Go(zetaTestRoutine(conf, deployerRunner, verbose, zetaTests...))
		eg.Go(zevmMPTestRoutine(conf, deployerRunner, verbose, zevmMPTests...))
		eg.Go(bitcoinTestRoutine(conf, deployerRunner, verbose, !skipBitcoinSetup, bitcoinTests...))
		eg.Go(ethereumTestRoutine(conf, deployerRunner, verbose, ethereumTests...))
	}

	if testAdmin {
		eg.Go(adminTestRoutine(conf, deployerRunner, verbose,
			e2etests.TestWhitelistERC20Name,
			e2etests.TestPauseZRC20Name,
			e2etests.TestUpdateBytecodeZRC20Name,
			e2etests.TestUpdateBytecodeConnectorName,
			e2etests.TestDepositEtherLiquidityCapName,
			e2etests.TestCriticalAdminTransactionsName,
			e2etests.TestPauseERC20CustodyName,
			e2etests.TestMigrateERC20CustodyFundsName,

			// Test the rate limiter functionalities
			// this test is currently incomplete and takes 10m to run
			// TODO: define assertion, and make more optimized
			// https://github.com/zeta-chain/node/issues/2090
			//e2etests.TestRateLimiterName,

			// TestMigrateChainSupportName tests EVM chain migration. Currently this test doesn't work with Anvil because pre-EIP1559 txs are not supported
			// See issue below for details
			// TODO: reenable this test as per the issue below
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
	if testSolana {
		if deployerRunner.SolanaClient == nil {
			logger.Print("❌ solana client is nil, maybe solana rpc is not set")
			os.Exit(1)
		}
		solanaTests := []string{
			e2etests.TestSolanaDepositName,
			e2etests.TestSolanaWithdrawName,
			e2etests.TestSolanaDepositAndCallName,
			e2etests.TestSolanaDepositAndCallRefundName,
			e2etests.TestSolanaDepositRestrictedName,
			e2etests.TestSolanaWithdrawRestrictedName,
		}
		eg.Go(solanaTestRoutine(conf, deployerRunner, verbose, solanaTests...))
	}

	if testTON {
		if deployerRunner.Clients.TON == nil {
			logger.Print("❌ TON client is nil, maybe TON lite-server config is not set")
			os.Exit(1)
		}

		tonTests := []string{
			e2etests.TestTONDepositName,
		}

		eg.Go(tonTestRoutine(conf, deployerRunner, verbose, tonTests...))
	}

	if testV2 {
		// update the ERC20 custody contract for v2 tests
		// note: not run in testV2Migration because it is already run in the migration process
		deployerRunner.UpdateChainParamsV2Contracts()
	}

	if testV2 || testV2Migration {
		startV2Tests(&eg, conf, deployerRunner, verbose)
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
	logger.Print("⏳ e2e tests passed,checking tx priority")
	monitorPriorityCancel()
	if err := <-txPriorityErrCh; err != nil && errors.Is(err, errWrongTxPriority) {
		logger.Print("❌ %v", err)
		logger.Print("❌ e2e tests failed after %s", time.Since(testStartTime).String())
		os.Exit(1)
	}

	logger.Print("✅ e2e tests completed in %s", time.Since(testStartTime).String())

	if testTSSMigration {
		TSSMigration(deployerRunner, logger, verbose, conf)
	}
	// Verify that there are no trackers left over after tests complete
	if !skipTrackerCheck {
		deployerRunner.EnsureNoTrackers()
	}
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
	observerClient observertypes.QueryClient,
	logger *runner.Logger,
	bufferBlocks int64,
) {
	// wait for keygen to be completed
	resp, err := observerClient.Keygen(ctx, &observertypes.QueryGetKeygenRequest{})
	if err != nil {
		logger.Error("observerClient.Keygen error: %s", err)
		return
	}
	if resp.Keygen == nil {
		logger.Error("observerClient.Keygen keygen is nil")
		return
	}
	if resp.Keygen.Status != observertypes.KeygenStatus_PendingKeygen {
		return
	}
	keygenHeight := resp.Keygen.BlockNumber
	logger.Print("⏳ wait height %v for keygen to be completed", keygenHeight)
	for {
		time.Sleep(2 * time.Second)
		response, err := cctxClient.LastZetaHeight(ctx, &crosschaintypes.QueryLastZetaHeightRequest{})
		if err != nil {
			logger.Error("cctxClient.LastZetaHeight error: %s", err)
			continue
		}
		if response.Height >= keygenHeight+bufferBlocks {
			break
		}
		logger.Info("Last ZetaHeight: %d", response.Height)
	}
}

func must[T any](v T, err error) T {
	return testutil.Must(v, err)
}
