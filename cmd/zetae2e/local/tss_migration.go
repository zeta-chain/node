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
)

// tssMigrationTestRoutine runs TSS migration related e2e tests
func tssMigrationTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		account := conf.AdditionalAccounts.UserMigration
		// initialize runner for migration test
		tssMigrationTestRunner, err := initTestRunner(
			"tssMigration",
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

		if err := tssMigrationTestRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("TSS migration tests failed: %v", err)
		}
		if err := tssMigrationTestRunner.CheckBtcTSSBalance(); err != nil {
			return err
		}

		tssMigrationTestRunner.Logger.Print("🍾 TSS migration tests completed in %s", time.Since(startTime).String())

		return nil
	}
}

func TSSMigration(deployerRunner *runner.E2ERunner, logger *runner.Logger, verbose bool, conf config.Config) {
	migrationStartTime := time.Now()
	logger.Print("🏁 starting tss migration")

	response, err := deployerRunner.CctxClient.LastZetaHeight(
		deployerRunner.Ctx,
		&crosschaintypes.QueryLastZetaHeightRequest{},
	)
	require.NoError(deployerRunner, err)
	err = deployerRunner.ZetaTxServer.UpdateKeygen(response.Height)
	require.NoError(deployerRunner, err)

	// Generate new TSS
	waitKeygenHeight(deployerRunner.Ctx, deployerRunner.CctxClient, deployerRunner.ObserverClient, logger, 0)

	// Run migration
	// migrationRoutine runs migration e2e test , which migrates funds from the older TSS to the new one
	// The zetaclient restarts required for this process are managed by the background workers in zetaclient (TSSListener)
	fn := tssMigrationTestRoutine(conf, deployerRunner, verbose, e2etests.TestMigrateTSSName)

	if err := fn(); err != nil {
		logger.Print("❌ %v", err)
		logger.Print("❌ tss migration failed")
		os.Exit(1)
	}
	deployerRunner.UpdateTssAddressForConnector()
	deployerRunner.UpdateTssAddressForErc20custody()
	logger.Print("✅ migration completed in %s ", time.Since(migrationStartTime).String())
}
