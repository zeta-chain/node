package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/e2etests"
	"github.com/zeta-chain/zetacore/e2e/runner"
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
			"migration",
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

		return err
	}
}
