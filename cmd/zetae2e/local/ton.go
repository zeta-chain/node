package local

import (
	"time"

	"github.com/fatih/color"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
)

// tonTestRoutine runs TON related e2e tests
func tonTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		tonRunner, err := initTestRunner(
			"ton",
			conf,
			deployerRunner,
			conf.DefaultAccount,
			runner.NewLogger(verbose, color.FgCyan, "ton"),
		)
		if err != nil {
			return errors.Wrap(err, "unable to init ton test runner")
		}

		tonRunner.Logger.Print("🏃 starting TON tests")
		startTime := time.Now()

		tests, err := tonRunner.GetE2ETestsToRunByName(e2etests.AllE2ETests, testNames...)
		if err != nil {
			return errors.Wrap(err, "unable to get ton tests to run")
		}

		if err := tonRunner.SetupTON(); err != nil {
			return errors.Wrap(err, "unable to setup TON account")
		}

		if err := tonRunner.RunE2ETests(tests); err != nil {
			return errors.Wrap(err, "ton tests failed")
		}

		tonRunner.Logger.Print("🍾 ton tests completed in %s", time.Since(startTime).String())

		return nil
	}
}
