package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/e2etests"
	"github.com/zeta-chain/zetacore/e2e/runner"
)

// erc20TestRoutine runs erc20 related e2e tests
func statefulPrecompilesTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		account := conf.AdditionalAccounts.UserERC20

		// initialize runner for erc20 test
		erc20Runner, err := initTestRunner(
			"precompiles",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgRed, "precompiles"),
		)
		if err != nil {
			return err
		}

		erc20Runner.Logger.Print("🏃 starting stateful precompiled contracts tests")
		startTime := time.Now()

		// run erc20 test
		testsToRun, err := erc20Runner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("precompiled contracts tests failed: %v", err)
		}

		if err := erc20Runner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("precompiled contracts tests failed: %v", err)
		}

		erc20Runner.Logger.Print("🍾 precompiled contracts tests completed in %s", time.Since(startTime).String())

		return err
	}
}
