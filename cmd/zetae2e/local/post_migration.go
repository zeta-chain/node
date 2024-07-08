package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/e2etests"
	"github.com/zeta-chain/zetacore/e2e/runner"
)

// postMigrationTestRoutine runs post migration tests
func postMigrationTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		account := conf.AdditionalAccounts.UserBitcoin
		// initialize runner for post migration test
		postMigrationRunner, err := initTestRunner(
			"postMigration",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgMagenta, "postMigrationRunner"),
		)
		if err != nil {
			return err
		}

		postMigrationRunner.Logger.Print("🏃 starting postMigration tests")
		startTime := time.Now()

		testsToRun, err := postMigrationRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("postMigrationRunner tests failed: %v", err)
		}

		if err := postMigrationRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("postMigrationRunner tests failed: %v", err)
		}

		if err := postMigrationRunner.CheckBtcTSSBalance(); err != nil {
			return err
		}

		postMigrationRunner.Logger.Print("🍾 PostMigration tests completed in %s", time.Since(startTime).String())

		return err
	}
}
