package local

import (
	"time"

	"github.com/fatih/color"
	"github.com/pkg/errors"

	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/e2etests"
	"github.com/zeta-chain/zetacore/e2e/runner"
)

// postMigrationTestRoutine runs post TSS migration tests
func postTSSMigrationTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		account := conf.AdditionalAccounts.UserBitcoin
		// initialize runner for post migration test
		postMigrationRunner, err := initTestRunner(
			"postTSSMigration",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgMagenta, "postTSSMigration"),
		)
		if err != nil {
			return err
		}

		postMigrationRunner.Logger.Print("🏃 starting post TSS migration tests")
		startTime := time.Now()

		testsToRun, err := postMigrationRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return errors.Wrap(err, "post TSS migration tests failed")
		}

		if err := postMigrationRunner.RunE2ETests(testsToRun); err != nil {
			return errors.Wrap(err, "post TSS migration tests failed")
		}

		if err := postMigrationRunner.CheckBtcTSSBalance(); err != nil {
			return err
		}

		postMigrationRunner.Logger.Print("🍾 post TSS migration tests completed in %s", time.Since(startTime).String())

		return err
	}
}
