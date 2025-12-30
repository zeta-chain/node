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

	// Generate 2 new TSS addresses
	for i := 0; i < numberOfTssToGenerate; i++ {
		logger.Print("🔑 generating TSS %d/2", i+1)

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
	// The zetaclient restarts required for this process are managed by the background workers in zetaclient (TSSListener)
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
