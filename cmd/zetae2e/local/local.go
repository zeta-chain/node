package local

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"

	zetae2econfig "github.com/zeta-chain/node/cmd/zetae2e/config"
	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/txserver"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/errgroup"
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
	flagIterations        = "iterations"
	flagTestSolana        = "test-solana"
	flagTestTON           = "test-ton"
	flagTestSui           = "test-sui"
	flagSkipRegular       = "skip-regular"
	flagLight             = "light"
	flagSetupOnly         = "setup-only"
	flagSkipSetup         = "skip-setup"
	flagTestTSSMigration  = "test-tss-migration"
	flagSkipBitcoinSetup  = "skip-bitcoin-setup"
	flagSkipHeaderProof   = "skip-header-proof"
	flagTestLegacy        = "test-legacy"
	flagSkipTrackerCheck  = "skip-tracker-check"
	flagSkipPrecompiles   = "skip-precompiles"
	flagUpgradeContracts  = "upgrade-contracts"
	flagTestFilter        = "test-filter"
	flagTestStaking       = "test-staking"
)

var (
	TestTimeout        = 20 * time.Minute
	ErrTopLevelTimeout = errors.New("top level test timeout")
	noError            = testutil.NoError
)

// NewLocalCmd returns the local command
// which runs the E2E tests locally on the machine with localnet for each blockchain
func NewLocalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "local",
		Short: "Run Local E2E tests",
		Run:   localE2ETest,
	}
	cmd.Flags().Bool(flagContractsDeployed, false, "set to to true if running tests again with existing state")
	cmd.Flags().Int64(flagWaitForHeight, 1, "block height for tests to begin, ex. --wait-for 100")
	cmd.Flags().String(FlagConfigFile, "", "config file to use for the tests")
	cmd.Flags().Bool(flagVerbose, false, "set to true to enable verbose logging")
	cmd.Flags().Bool(flagTestAdmin, false, "set to true to run admin tests")
	cmd.Flags().Bool(flagTestPerformance, false, "set to true to run performance tests")
	cmd.Flags().Int(flagIterations, 100, "number of iterations to run each performance test")
	cmd.Flags().Bool(flagTestSolana, false, "set to true to run solana tests")
	cmd.Flags().Bool(flagTestTON, false, "set to true to run TON tests")
	cmd.Flags().Bool(flagTestSui, false, "set to true to run Sui tests")
	cmd.Flags().Bool(flagSkipRegular, false, "set to true to skip regular tests")
	cmd.Flags().Bool(flagLight, false, "run the most basic regular tests, useful for quick checks")
	cmd.Flags().Bool(flagSetupOnly, false, "set to true to only setup the networks")
	cmd.Flags().String(flagConfigOut, "", "config file to write the deployed contracts from the setup")
	cmd.Flags().Bool(flagSkipSetup, false, "set to true to skip setup")
	cmd.Flags().Bool(flagSkipBitcoinSetup, false, "set to true to skip bitcoin wallet setup")
	cmd.Flags().Bool(flagSkipHeaderProof, false, "set to true to skip header proof tests")
	cmd.Flags().Bool(flagTestTSSMigration, false, "set to true to include a migration test at the end")
	cmd.Flags().Bool(flagTestLegacy, false, "set to true to run legacy EVM tests")
	cmd.Flags().Bool(flagSkipTrackerCheck, false, "set to true to skip tracker check at the end of the tests")
	cmd.Flags().Bool(flagSkipPrecompiles, true, "set to true to skip stateful precompiled contracts test")
	cmd.Flags().
		Bool(flagUpgradeContracts, false, "set to true to upgrade Gateways and ERC20Custody contracts during setup for ZEVM and EVM")
	cmd.Flags().String(flagTestFilter, "", "regexp filter to limit which test to run")
	cmd.Flags().Bool(flagTestStaking, false, "set to true to run staking tests")

	cmd.AddCommand(NewGetZetaclientBootstrap())

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
		iterations        = must(cmd.Flags().GetInt(flagIterations))
		testSolana        = must(cmd.Flags().GetBool(flagTestSolana))
		testTON           = must(cmd.Flags().GetBool(flagTestTON))
		testSui           = must(cmd.Flags().GetBool(flagTestSui))
		skipRegular       = must(cmd.Flags().GetBool(flagSkipRegular))
		light             = must(cmd.Flags().GetBool(flagLight))
		setupOnly         = must(cmd.Flags().GetBool(flagSetupOnly))
		skipSetup         = must(cmd.Flags().GetBool(flagSkipSetup))
		skipBitcoinSetup  = must(cmd.Flags().GetBool(flagSkipBitcoinSetup))
		skipHeaderProof   = must(cmd.Flags().GetBool(flagSkipHeaderProof))
		skipTrackerCheck  = must(cmd.Flags().GetBool(flagSkipTrackerCheck))
		testTSSMigration  = must(cmd.Flags().GetBool(flagTestTSSMigration))
		testLegacy        = must(cmd.Flags().GetBool(flagTestLegacy))
		skipPrecompiles   = must(cmd.Flags().GetBool(flagSkipPrecompiles))
		upgradeContracts  = must(cmd.Flags().GetBool(flagUpgradeContracts))
		setupSolana       = testSolana || testPerformance
		testFilterStr     = must(cmd.Flags().GetString(flagTestFilter))
		testStaking       = must(cmd.Flags().GetBool(flagTestStaking))
	)

	testFilter := regexp.MustCompile(testFilterStr)

	logger := runner.NewLogger(verbose, color.FgWhite, "setup")

	testStartTime := time.Now()

	logger.Print("starting E2E tests")

	if verbose {
		logger.Info("Flags")
		cmd.Flags().VisitAll(func(flag *pflag.Flag) {
			logger.Info(`--%s="%s"`, flag.Name, flag.Value.String())
		})
	}

	if testAdmin {
		logger.Print("⚠️ admin tests enabled")
	}

	if testPerformance {
		logger.Print("⚠️ performance tests enabled, regular tests will be skipped")
		skipRegular = true
		skipPrecompiles = true
	}

	// initialize tests config
	conf, err := GetConfig(cmd)
	noError(err)

	if testPerformance && iterations > 100 {
		TestTimeout = time.Hour
	}

	// initialize context
	ctx, timeoutCancel := context.WithTimeoutCause(context.Background(), TestTimeout, ErrTopLevelTimeout)
	defer timeoutCancel()
	ctx, cancel := context.WithCancelCause(ctx)

	// route os signals to context cancellation.
	// using NotifyContext will ensure that the second signal
	// will not be handled and should kill the process.
	go func() {
		notifyCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
		<-notifyCtx.Done()
		cancel(fmt.Errorf("notify context: %w", notifyCtx.Err()))
		stop()
	}()

	// wait for a specific height on ZetaChain
	noError(utils.WaitForBlockHeight(ctx, waitForHeight, conf.RPCs.ZetaCoreRPC, logger))

	zetaTxServer, err := txserver.NewZetaTxServer(
		conf.RPCs.ZetaCoreRPC,
		[]string{
			utils.EmergencyPolicyName,
			utils.OperationalPolicyName,
			utils.AdminPolicyName,
			utils.UserEmissionsWithdrawName,
		},
		[]string{
			conf.PolicyAccounts.EmergencyPolicyAccount.RawPrivateKey.String(),
			conf.PolicyAccounts.OperationalPolicyAccount.RawPrivateKey.String(),
			conf.PolicyAccounts.AdminPolicyAccount.RawPrivateKey.String(),
			conf.AdditionalAccounts.UserEmissionsWithdraw.RawPrivateKey.String(),
		},
		conf.ZetaChainID,
	)
	noError(err)

	// Drop this cond after TON e2e is included in the default suite
	if !testTON {
		conf.RPCs.TON = ""
	}

	// initialize deployer runner with config
	deployerRunner, err := zetae2econfig.RunnerFromConfig(
		ctx,
		"deployer",
		cancel,
		conf,
		conf.DefaultAccount,
		logger,
		runner.WithZetaTxServer(zetaTxServer),
		runner.WithTestFilter(testFilter),
	)
	noError(err)

	// monitor block production to ensure we fail fast if there are consensus failures
	go monitorBlockProductionCancel(ctx, cancel, conf)

	// set the authority client to the zeta tx server to be able to query message permissions
	deployerRunner.ZetaTxServer.SetAuthorityClient(deployerRunner.AuthorityClient)

	if !skipSetup {
		// run setup steps that do not require tss
		noError(deployerRunner.FundEmissionsPool())

		// wait for keygen to be completed
		// if setup is skipped, we assume that the keygen is already completed
		noError(waitKeygenHeight(ctx, deployerRunner.CctxClient, deployerRunner.ObserverClient, logger, 10))
	}

	// query and set the TSS
	noError(deployerRunner.SetTSSAddresses())

	if !skipHeaderProof {
		noError(deployerRunner.EnableHeaderVerification([]int64{
			chains.GoerliLocalnet.ChainId,
			chains.BitcoinRegtest.ChainId,
		}))
	}

	e2eStartHeight, err := deployerRunner.Clients.Zetacore.GetBlockHeight(ctx)
	noError(err)

	// setting up the networks
	if !skipSetup {
		logger.Print("⚙️ setting up networks")
		startTime := time.Now()

		// setup TSS address and setup deployer wallet
		deployerRunner.SetupBitcoinAccounts(true)

		//setup protocol contracts v1 as they are still supported for now
		deployerRunner.LegacySetupEVM(contractsDeployed, testLegacy)

		// setup protocol contracts on the connected EVM chain
		deployerRunner.SetupEVM()

		if setupSolana {
			deployerRunner.SetupSolana(
				conf.Contracts.Solana.GatewayProgramID.String(),
				conf.AdditionalAccounts.UserSolana.SolanaPrivateKey.String(),
				conf.AdditionalAccounts.UserSPL.SolanaPrivateKey.String(),
			)
		}

		deployerRunner.SetupZEVMProtocolContracts()
		deployerRunner.SetupLegacyZEVMContracts()

		zrc20Deployment := txserver.ZRC20Deployment{
			ERC20Addr: deployerRunner.ERC20Addr,
			SPLAddr:   nil,
		}
		if setupSolana {
			zrc20Deployment.SPLAddr = deployerRunner.SPLAddr.ToPointer()
		}
		deployerRunner.SetupZEVMZRC20s(zrc20Deployment)

		// Update the chain params to contains protocol contract addresses
		deployerRunner.UpdateProtocolContractsInChainParams()

		if testTON {
			deployerRunner.SetupTON(
				conf.RPCs.TONFaucet,
				conf.AdditionalAccounts.UserTON,
			)
		}

		if testSui {
			deployerRunner.SetupSui(conf.RPCs.SuiFaucet)
		}

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

	if upgradeContracts {
		deployerRunner.UpgradeGatewaysAndERC20Custody()
	}
	// always mint ERC20 before every test execution
	deployerRunner.MintERC20OnEVM(1e10)

	// Run the proposals under the start sequence(proposals_e2e_start folder)
	if !skipRegular {
		noError(deployerRunner.CreateGovProposals(runner.StartOfE2E))
	}

	// run tests
	var eg errgroup.Group

	if !skipRegular {
		// start the EVM tests
		startEVMTests(&eg, conf, deployerRunner, verbose)
		startBitcoinTests(&eg, conf, deployerRunner, verbose, light, skipBitcoinSetup)
	}

	if !skipPrecompiles {
		precompiledContractTests := []string{
			//e2etests.TestPrecompilesPrototypeName,
			//e2etests.TestPrecompilesPrototypeThroughContractName,
			//// Disabled until further notice, check https://github.com/zeta-chain/node/issues/3005.
			//// e2etests.TestPrecompilesStakingThroughContractName,
			//e2etests.TestPrecompilesBankName,
			//e2etests.TestPrecompilesBankFailName,
			//e2etests.TestPrecompilesBankThroughContractName,
		}
		if e2eStartHeight < 100 {
			// these tests require a clean system
			// since unstaking has an unbonding period
			//precompiledContractTests = append(precompiledContractTests,
			//	e2etests.TestPrecompilesStakingName,
			//	e2etests.TestPrecompilesDistributeName,
			//	e2etests.TestPrecompilesDistributeNonZRC20Name,
			//	e2etests.TestPrecompilesDistributeThroughContractName,
			//)
			// prevent lint error
			_ = precompiledContractTests
		} else {
			logger.Print("⚠️ partial precompiled run (unclean state)")
		}
		eg.Go(statefulPrecompilesTestRoutine(conf, deployerRunner, verbose, precompiledContractTests...))
	}

	if testAdmin {
		eg.Go(adminTestRoutine(conf, deployerRunner, verbose,
			e2etests.TestUpdateZRC20NameName,
			e2etests.TestZetaclientSignerOffsetName,
			e2etests.TestZetaclientRestartHeightName,
			e2etests.TestWhitelistERC20Name,
			e2etests.TestPauseZRC20Name,
			e2etests.TestUpdateBytecodeZRC20Name,
			e2etests.TestUpdateBytecodeConnectorName,
			e2etests.TestDepositEtherLiquidityCapName,
			e2etests.TestCriticalAdminTransactionsName,
			e2etests.TestPauseERC20CustodyName,
			e2etests.TestMigrateERC20CustodyFundsName,
			e2etests.TestUpdateOperationalChainParamsName,

			// Currently this test doesn't work with Anvil because pre-EIP1559 txs are not supported
			// See issue below for details
			// TODO: reenable this test as per the issue below
			// https://github.com/zeta-chain/node/issues/1980
			// e2etests.TestMigrateChainSupportName,
		))
	}

	if testPerformance {
		eg.Go(
			ethereumDepositPerformanceRoutine(
				conf,
				deployerRunner,
				verbose,
				[]string{e2etests.TestStressEtherDepositName},
				iterations,
			),
		)
		eg.Go(
			ethereumWithdrawPerformanceRoutine(
				conf,
				deployerRunner,
				verbose,
				[]string{e2etests.TestStressEtherWithdrawName},
				iterations,
			),
		)
		eg.Go(
			solanaDepositPerformanceRoutine(
				conf,
				"perf_sol_deposit",
				deployerRunner,
				verbose,
				conf.AdditionalAccounts.UserSolana,
				[]string{e2etests.TestStressSolanaDepositName},
			),
		)
		eg.Go(
			solanaDepositPerformanceRoutine(
				conf,
				"perf_spl_deposit",
				deployerRunner,
				verbose,
				conf.AdditionalAccounts.UserSPL,
				[]string{e2etests.TestStressSPLDepositName},
			),
		)
		eg.Go(
			solanaWithdrawPerformanceRoutine(
				conf,
				"perf_sol_withdraw",
				deployerRunner,
				verbose,
				conf.AdditionalAccounts.UserSolana,
				[]string{e2etests.TestStressSolanaWithdrawName},
			),
		)
		eg.Go(
			solanaWithdrawPerformanceRoutine(
				conf,
				"perf_spl_withdraw",
				deployerRunner,
				verbose,
				conf.AdditionalAccounts.UserSPL,
				[]string{e2etests.TestStressSPLWithdrawName},
			),
		)
	}

	if testSolana {
		if deployerRunner.SolanaClient == nil {
			logger.Print("❌ solana client is nil, maybe solana rpc is not set")
			os.Exit(1)
		}
		// Run only basic solana tests if an upgrade is in progress
		// This is done to avoid running the tests that take too long to complete
		// Related : https://github.com/zeta-chain/node/issues/3666
		solanaTests := []string{
			e2etests.TestSolanaDepositName,
			e2etests.TestSolanaWithdrawName,
		}

		splTests := []string{
			e2etests.TestSolanaDepositName,
			e2etests.TestSPLDepositName,
		}

		if !deployerRunner.IsRunningUpgrade() {
			solanaTests = append(solanaTests, []string{
				e2etests.TestSolanaDepositThroughProgramName,
				e2etests.TestSolanaDepositAndCallName,
				e2etests.TestSolanaWithdrawAndCallName,
				e2etests.TestZEVMToSolanaCallName,
				e2etests.TestSolanaWithdrawAndCallRevertWithCallName,
				e2etests.TestSolanaDepositAndCallRevertName,
				e2etests.TestSolanaDepositAndCallRevertWithCallName,
				e2etests.TestSolanaDepositAndCallRevertWithCallThatRevertsName,
				e2etests.TestSolanaDepositAndCallRevertWithDustName,
				e2etests.TestSolanaDepositRestrictedName,
				e2etests.TestSolanaToZEVMCallName,
				e2etests.TestSolanaWithdrawRestrictedName,
			}...)

			splTests = append(splTests, []string{
				e2etests.TestSPLDepositAndCallName,
				e2etests.TestSPLDepositAndCallRevertName,
				e2etests.TestSPLDepositAndCallRevertWithCallName,
				e2etests.TestSPLDepositAndCallRevertWithCallThatRevertsName,
				e2etests.TestSPLWithdrawName,
				e2etests.TestSPLWithdrawAndCallName,
				e2etests.TestSPLWithdrawAndCallRevertName,
				e2etests.TestSPLWithdrawAndCreateReceiverAtaName,
				// TODO move under admin tests
				// https://github.com/zeta-chain/node/issues/3085
				e2etests.TestSolanaWhitelistSPLName,
			}...)
		}

		eg.Go(
			solanaTestRoutine(
				conf,
				"solana",
				conf.AdditionalAccounts.UserSolana,
				deployerRunner,
				verbose,
				solanaTests...),
		)
		eg.Go(
			solanaTestRoutine(
				conf,
				"spl",
				conf.AdditionalAccounts.UserSPL,
				deployerRunner,
				verbose,
				splTests...),
		)
	}

	if testSui {
		suiTests := []string{
			e2etests.TestSuiDepositName,
			e2etests.TestSuiDepositAndCallRevertName,
			e2etests.TestSuiDepositAndCallName,
			e2etests.TestSuiTokenDepositName,
			e2etests.TestSuiTokenDepositAndCallName,
			e2etests.TestSuiTokenDepositAndCallRevertName,
			e2etests.TestSuiWithdrawName,
			e2etests.TestSuiWithdrawAndCallName,
			e2etests.TestSuiWithdrawRevertWithCallName,
			e2etests.TestSuiWithdrawAndCallRevertWithCallName,
			e2etests.TestSuiTokenWithdrawName,
			e2etests.TestSuiTokenWithdrawAndCallName,
			e2etests.TestSuiTokenWithdrawAndCallRevertWithCallName,
			e2etests.TestSuiDepositRestrictedName,
			e2etests.TestSuiWithdrawRestrictedName,
		}
		eg.Go(suiTestRoutine(conf, deployerRunner, verbose, suiTests...))
	}

	if testTON {
		if deployerRunner.Clients.TON == nil {
			logger.Print("❌ TON client is nil, maybe TON lite-server config is not set")
			os.Exit(1)
		}

		tonTests := []string{
			e2etests.TestTONDepositName,
			e2etests.TestTONDepositAndCallName,
			e2etests.TestTONDepositAndCallRefundName,
			e2etests.TestTONDepositRestrictedName,
			e2etests.TestTONWithdrawName,
			e2etests.TestTONWithdrawConcurrentName,
		}

		eg.Go(tonTestRoutine(conf, deployerRunner, verbose, tonTests...))
	}

	if testLegacy {
		eg.Go(legacyERC20TestRoutine(conf, deployerRunner, verbose,
			e2etests.TestLegacyERC20WithdrawName,
			e2etests.TestLegacyMultipleERC20WithdrawsName,
			e2etests.TestLegacyERC20DepositAndCallRefundName))
		eg.Go(legacyZETATestRoutine(conf, deployerRunner, verbose,
			e2etests.TestLegacyZetaWithdrawName,
			e2etests.TestLegacyMessagePassingExternalChainsName,
			e2etests.TestLegacyMessagePassingRevertFailExternalChainsName,
			e2etests.TestLegacyMessagePassingRevertSuccessExternalChainsName,
			e2etests.TestLegacyZetaDepositRestrictedName,
			e2etests.TestLegacyZetaDepositName,
			e2etests.TestLegacyZetaDepositNewAddressName,
		))
		eg.Go(legacyZEVMMPTestRoutine(conf, deployerRunner, verbose,
			e2etests.TestLegacyMessagePassingZEVMToEVMName,
			e2etests.TestLegacyMessagePassingEVMtoZEVMName,
			e2etests.TestLegacyMessagePassingEVMtoZEVMRevertName,
			e2etests.TestLegacyMessagePassingZEVMtoEVMRevertName,
			e2etests.TestLegacyMessagePassingZEVMtoEVMRevertFailName,
			e2etests.TestLegacyMessagePassingEVMtoZEVMRevertFailName,
		))
		eg.Go(legacyEthereumTestRoutine(conf, deployerRunner, verbose,
			e2etests.TestLegacyEtherWithdrawName,
			e2etests.TestLegacyEtherDepositAndCallName,
			e2etests.TestLegacyEtherDepositAndCallRefundName,
		))
	}

	// while tests are executed, monitor blocks in parallel to check if system txs are on top and they have biggest priority
	txPriorityErrCh := make(chan error, 1)
	ctx, monitorPriorityCancel := context.WithCancel(context.Background())
	go monitorTxPriorityInBlocks(ctx, conf, txPriorityErrCh)

	if err := eg.Wait(); err != nil {
		deployerRunner.CtxCancel(err)
		monitorPriorityCancel()
		logger.Print("❌ %v", err)
		logger.Print("❌ e2e tests failed after %s", time.Since(testStartTime).String())
		os.Exit(1)
	}

	// Default ballot maturity is set to 30 blocks.
	// We can wait for 31 blocks to ensure that all ballots created during the test are matured, as emission rewards may be slashed for some of the observers based on their vote.
	// This seems to be a problem only in performance tests where we are creating a lot of ballots in a short time. We do not need to slow down regular tests for this check as we expect all observers to vote correctly.
	if testPerformance {
		deployerRunner.WaitForBlocks(31)
	}

	noError(deployerRunner.WithdrawEmissions())

	// if all tests pass, cancel txs priority monitoring and check if tx priority is not correct in some blocks
	logger.Print("⏳ e2e tests passed, checking tx priority")
	monitorPriorityCancel()
	if err := <-txPriorityErrCh; err != nil && errors.Is(err, errWrongTxPriority) {
		logger.Print("❌ %v", err)
		logger.Print("❌ e2e tests failed after %s", time.Since(testStartTime).String())
		os.Exit(1)
	}
	if !skipRegular {
		noError(deployerRunner.CreateGovProposals(runner.EndOfE2E))
	}

	logger.Print("✅ e2e tests completed in %s", time.Since(testStartTime).String())

	if testSolana {
		require.True(
			deployerRunner,
			deployerRunner.VerifySolanaContractsUpgrade(conf.AdditionalAccounts.UserSolana.SolanaPrivateKey.String()),
		)
	}

	if testTSSMigration {
		TSSMigration(deployerRunner, logger, verbose, conf)
	}

	// Verify that there are no trackers left over after tests complete
	if !skipTrackerCheck {
		deployerRunner.EnsureNoTrackers()
	}

	// Verify that the balance of restricted address is zero
	deployerRunner.EnsureZeroBalanceOnRestrictedAddressZEVM()

	if !deployerRunner.IsRunningUpgrade() {
		// Verify that there are no stale ballots left over after tests complete
		deployerRunner.EnsureNoStaleBallots()
	}

	// This should only be run at the end to the test as it would remove the observer.
	if testStaking {
		e2etests.UndelegateToBelowMinimumObserverDelegation(deployerRunner, []string{})
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
) error {
	// wait for keygen to be completed
	resp, err := observerClient.Keygen(ctx, &observertypes.QueryGetKeygenRequest{})

	switch {
	case err != nil:
		return errors.Wrap(err, "observerClient.Keygen error")
	case resp.Keygen == nil:
		return errors.New("keygen is nil")
	case resp.Keygen.Status != observertypes.KeygenStatus_PendingKeygen:
		return errors.Errorf("keygen is not pending (status: %s)", resp.Keygen.Status.String())
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

		logger.Info("Last ZetaHeight: %d", response.Height)

		if response.Height >= keygenHeight+bufferBlocks {
			return nil
		}
	}
}

func must[T any](v T, err error) T {
	return testutil.Must(v, err)
}
