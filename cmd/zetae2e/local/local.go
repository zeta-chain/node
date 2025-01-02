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
	flagTestSolana        = "test-solana"
	flagTestTON           = "test-ton"
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
)

var (
	TestTimeout = 20 * time.Minute
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
	cmd.Flags().Bool(flagTestLegacy, false, "set to true to run legacy EVM tests")
	cmd.Flags().Bool(flagSkipTrackerCheck, false, "set to true to skip tracker check at the end of the tests")
	cmd.Flags().Bool(flagSkipPrecompiles, false, "set to true to skip stateful precompiled contracts test")
	cmd.Flags().
		Bool(flagUpgradeContracts, false, "set to true to upgrade Gateways and ERC20Custody contracts during setup for ZEVM and EVM")

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
		testLegacy        = must(cmd.Flags().GetBool(flagTestLegacy))
		skipPrecompiles   = must(cmd.Flags().GetBool(flagSkipPrecompiles))
		upgradeContracts  = must(cmd.Flags().GetBool(flagUpgradeContracts))
		setupSolana       = testSolana || testPerformance
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
		skipPrecompiles = true
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

	// monitor block production to ensure we fail fast if there are consensus failures
	// this is not run in an errgroup since only returning an error will not exit immedately
	// this needs to be early to quickly detect consensus failure during genesis
	go monitorBlockProductionExit(ctx, conf)

	// set the authority client to the zeta tx server to be able to query message permissions
	deployerRunner.ZetaTxServer.SetAuthorityClient(deployerRunner.AuthorityClient)

	// run setup steps that do not require tss
	if !skipSetup {
		noError(deployerRunner.FundEmissionsPool())
	}

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

	e2eStartHeight, err := deployerRunner.Clients.Zetacore.GetBlockHeight(ctx)
	noError(err)

	// setting up the networks
	if !skipSetup {
		logger.Print("⚙️ setting up networks")
		startTime := time.Now()

		//setup protocol contracts v1 as they are still supported for now
		deployerRunner.LegacySetupEVM(contractsDeployed)

		// setup protocol contracts on the connected EVM chain
		deployerRunner.SetupEVM()

		if setupSolana {
			deployerRunner.SetupSolana(
				conf.Contracts.Solana.GatewayProgramID.String(),
				conf.AdditionalAccounts.UserSolana.SolanaPrivateKey.String(),
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

	// run tests
	var eg errgroup.Group

	if !skipRegular {
		// start the EVM tests
		startEVMTests(&eg, conf, deployerRunner, verbose)

		// start the bitcoin tests
		// btc withdraw tests are those that need a Bitcoin node wallet to send UTXOs
		bitcoinDepositTests := []string{
			e2etests.TestBitcoinDonationName,
			e2etests.TestBitcoinDepositName,
			e2etests.TestBitcoinDepositAndCallName,
			e2etests.TestBitcoinDepositAndCallRevertName,
			e2etests.TestBitcoinStdMemoDepositName,
			e2etests.TestBitcoinStdMemoDepositAndCallName,
			e2etests.TestBitcoinStdMemoDepositAndCallRevertName,
			e2etests.TestBitcoinStdMemoInscribedDepositAndCallName,
			e2etests.TestCrosschainSwapName,
		}
		bitcoinDepositTestsAdvanced := []string{
			e2etests.TestBitcoinDepositAndCallRevertWithDustName,
			e2etests.TestBitcoinStdMemoDepositAndCallRevertOtherAddressName,
		}
		bitcoinWithdrawTests := []string{
			e2etests.TestBitcoinWithdrawSegWitName,
			e2etests.TestBitcoinWithdrawInvalidAddressName,
			e2etests.TestLegacyZetaWithdrawBTCRevertName,
		}
		bitcoinWithdrawTestsAdvanced := []string{
			e2etests.TestBitcoinWithdrawTaprootName,
			e2etests.TestBitcoinWithdrawLegacyName,
			e2etests.TestBitcoinWithdrawP2SHName,
			e2etests.TestBitcoinWithdrawP2WSHName,
			e2etests.TestBitcoinWithdrawMultipleName,
			e2etests.TestBitcoinWithdrawRestrictedName,
		}

		if !light {
			// if light is enabled, only the most basic tests are run and advanced are skipped
			bitcoinDepositTests = append(bitcoinDepositTests, bitcoinDepositTestsAdvanced...)
			bitcoinWithdrawTests = append(bitcoinWithdrawTests, bitcoinWithdrawTestsAdvanced...)
		}
		bitcoinDepositTestRoutine, bitcoinWithdrawTestRoutine := bitcoinTestRoutines(
			conf,
			deployerRunner,
			verbose,
			!skipBitcoinSetup,
			bitcoinDepositTests,
			bitcoinWithdrawTests,
		)
		eg.Go(bitcoinDepositTestRoutine)
		eg.Go(bitcoinWithdrawTestRoutine)
	}

	if !skipPrecompiles {
		precompiledContractTests := []string{
			e2etests.TestPrecompilesPrototypeName,
			e2etests.TestPrecompilesPrototypeThroughContractName,
			// Disabled until further notice, check https://github.com/zeta-chain/node/issues/3005.
			// e2etests.TestPrecompilesStakingThroughContractName,
			e2etests.TestPrecompilesBankName,
			e2etests.TestPrecompilesBankFailName,
			e2etests.TestPrecompilesBankThroughContractName,
		}
		if e2eStartHeight < 100 {
			// these tests require a clean system
			// since unstaking has an unbonding period
			precompiledContractTests = append(precompiledContractTests,
				e2etests.TestPrecompilesStakingName,
				e2etests.TestPrecompilesDistributeName,
				e2etests.TestPrecompilesDistributeNonZRC20Name,
				e2etests.TestPrecompilesDistributeThroughContractName,
			)
		} else {
			logger.Print("⚠️ partial precompiled run (unclean state)")
		}
		eg.Go(statefulPrecompilesTestRoutine(conf, deployerRunner, verbose, precompiledContractTests...))
	}

	if testAdmin {
		eg.Go(adminTestRoutine(conf, deployerRunner, verbose,
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

			// Currently this test doesn't work with Anvil because pre-EIP1559 txs are not supported
			// See issue below for details
			// TODO: reenable this test as per the issue below
			// https://github.com/zeta-chain/node/issues/1980
			// e2etests.TestMigrateChainSupportName,
		))
	}

	if testPerformance {
		eg.Go(ethereumDepositPerformanceRoutine(conf, deployerRunner, verbose, e2etests.TestStressEtherDepositName))
		eg.Go(ethereumWithdrawPerformanceRoutine(conf, deployerRunner, verbose, e2etests.TestStressEtherWithdrawName))
		eg.Go(
			solanaDepositPerformanceRoutine(
				conf,
				"perf_sol_deposit",
				deployerRunner,
				verbose,
				conf.AdditionalAccounts.UserSolana,
				e2etests.TestStressSolanaDepositName,
			),
		)
		eg.Go(
			solanaDepositPerformanceRoutine(
				conf,
				"perf_spl_deposit",
				deployerRunner,
				verbose,
				conf.AdditionalAccounts.UserSPL,
				e2etests.TestStressSPLDepositName,
			),
		)
		eg.Go(
			solanaWithdrawPerformanceRoutine(
				conf,
				"perf_sol_withdraw",
				deployerRunner,
				verbose,
				conf.AdditionalAccounts.UserSolana,
				e2etests.TestStressSolanaWithdrawName,
			),
		)
		eg.Go(
			solanaWithdrawPerformanceRoutine(
				conf,
				"perf_spl_withdraw",
				deployerRunner,
				verbose,
				conf.AdditionalAccounts.UserSPL,
				e2etests.TestStressSPLWithdrawName,
			),
		)
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
			e2etests.TestSolanaDepositAndCallRevertName,
			e2etests.TestSolanaDepositAndCallRevertWithDustName,
			e2etests.TestSolanaDepositRestrictedName,
			e2etests.TestSolanaWithdrawRestrictedName,
			// TODO move under admin tests
			// https://github.com/zeta-chain/node/issues/3085
			e2etests.TestSPLDepositName,
			e2etests.TestSPLDepositAndCallName,
			e2etests.TestSPLWithdrawName,
			e2etests.TestSPLWithdrawAndCreateReceiverAtaName,
			e2etests.TestSolanaWhitelistSPLName,
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
			e2etests.TestTONDepositAndCallName,
			e2etests.TestTONDepositAndCallRefundName,
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
		deployerRunner.CtxCancel()
		monitorPriorityCancel()
		logger.Print("❌ %v", err)
		logger.Print("❌ e2e tests failed after %s", time.Since(testStartTime).String())
		os.Exit(1)
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

	logger.Print("✅ e2e tests completed in %s", time.Since(testStartTime).String())

	if testTSSMigration {
		TSSMigration(deployerRunner, logger, verbose, conf)
	}

	// Verify that there are no trackers left over after tests complete
	if !skipTrackerCheck {
		deployerRunner.EnsureNoTrackers()
	}

	// Verify that the balance of restricted address is zero
	deployerRunner.EnsureZeroBalanceOnRestrictedAddressZEVM()

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
