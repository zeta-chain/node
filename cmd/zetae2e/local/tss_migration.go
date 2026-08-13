package local

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// tssMigrationTestRoutine runs TSS migration related e2e tests
func tssMigrationTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	expectedTssCount int,
	testNames ...string,
) func() error {
	return func() (err error) {
		account := conf.AdditionalAccounts.UserMigration
		// initialize runner for migration test
		tssMigrationTestRunner, err := initTestRunner(
			"triggerTSSMigration",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgHiGreen, "migration"),
			runner.WithZetaTxServer(deployerRunner.ZetaTxServer),
		)
		if err != nil {
			return err
		}

		tssMigrationTestRunner.Logger.Print("🏃 starting TSS migration tests")
		startTime := time.Now()

		if len(testNames) == 0 {
			tssMigrationTestRunner.Logger.Print("🍾 TSS migration tests completed in %s", time.Since(startTime).String())
			return nil
		}
		// run TSS migration test
		testsToRun, err := tssMigrationTestRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("TSS migration tests failed: %v", err)
		}
		tssMigrationTestRunner.WaitForTSSGeneration(int64(expectedTssCount))

		if err := tssMigrationTestRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("TSS migration tests failed: %v", err)
		}
		tssMigrationTestRunner.CheckBTCTSSBalance()
		tssMigrationTestRunner.Logger.Print("🍾 TSS migration tests completed in %s", time.Since(startTime).String())
		return nil
	}
}

// triggerTSSMigration generates a new TSS and migrates funds to it.
//
// This does not complete today. zetaclient no longer runs a keygen ceremony once a TSS key
// exists, so the wait below never sees KeyGenSuccess and blocks forever. Both CI suites that
// call it are disabled for the same reason (.github/workflows/e2e.yml), and the why is in the
// step 5 comment in zetaclient/tss/setup.go. Re-enable rotation before re-enabling this.
func triggerTSSMigration(
	deployerRunner *runner.E2ERunner,
	logger *runner.Logger,
	verbose bool,
	conf config.Config,
	testSolana bool,
	testSui bool,
	testTon bool,
) {
	migrationStartTime := time.Now()
	logger.Print("🏁 starting tss migration")

	tssList, err := deployerRunner.ObserverClient.TssHistory(
		deployerRunner.Ctx,
		&observertypes.QueryTssHistoryRequest{},
	)
	require.NoError(deployerRunner, err)
	// Increase this number to generate more than 1 TSS.
	// The migration always happens to the latest one, this is set on zetacore directly
	numberOfTssToGenerate := 1
	expectedTssCount := numberOfTssToGenerate + len(tssList.TssList)

	// Generate new TSS address(es)
	for i := 0; i < numberOfTssToGenerate; i++ {
		logger.Print("🔑 generating TSS %d/%d", i+1, numberOfTssToGenerate)

		response, err := deployerRunner.CctxClient.LastZetaHeight(
			deployerRunner.Ctx,
			&crosschaintypes.QueryLastZetaHeightRequest{},
		)
		require.NoError(deployerRunner, err)
		err = deployerRunner.ZetaTxServer.UpdateKeygen(response.Height)
		require.NoError(deployerRunner, err)

		// Generate new TSS
		noError(
			waitKeygenHeight(deployerRunner.Ctx, deployerRunner.CctxClient, deployerRunner.ObserverClient, logger, 0),
		)
	}

	// Run migration
	// migrationRoutine runs migration e2e test , which migrates funds from the older TSS to the new one
	// The TSSListener still restarts zetaclient when the TSS address changes or a new key lands in
	// history, but it no longer watches the keygen record, so nothing restarts the clients to run
	// the ceremony that would produce that new key in the first place.
	fn := tssMigrationTestRoutine(conf, deployerRunner, verbose, expectedTssCount, e2etests.TestMigrateTSSName)

	if err := fn(); err != nil {
		logger.Print("❌ %v", err)
		logger.Print("❌ tss migration failed")
		os.Exit(1)
	}

	// Update TSS address for contracts in connected chains
	// TODO : Update TSS address for other chains if necessary
	// https://github.com/zeta-chain/node/issues/3599
	deployerRunner.UpdateTSSAddressForConnectorNative()
	deployerRunner.UpdateTSSAddressForERC20custody()
	deployerRunner.UpdateTSSAddressForGateway()
	if testSolana {
		deployerRunner.UpdateTSSAddressSolana(
			conf.Contracts.Solana.GatewayProgramID.String(),
			conf.AdditionalAccounts.UserSolana.SolanaPrivateKey.String())
	}
	if testSui {
		deployerRunner.UpdateTSSAddressSui(conf.RPCs.SuiFaucet)
	}

	if testTon {
		deployerRunner.UpdateTSSAddressTON(
			conf.Contracts.TON.GatewayAccountID.String(),
			conf.RPCs.TONFaucet,
		)
	}
	logger.Print("✅ migration completed in %s ", time.Since(migrationStartTime).String())
}
