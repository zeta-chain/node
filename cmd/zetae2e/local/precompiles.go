package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
)

// statefulPrecompilesTestRoutine runs steateful precompiles related e2e tests
func statefulPrecompilesTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		account := conf.AdditionalAccounts.UserPrecompile

		precompileRunner, err := initTestRunner(
			"precompiles",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgRed, "precompiles"),
		)
		if err != nil {
			return err
		}

		precompileRunner.Logger.Print("🏃 starting stateful precompiled contracts tests")
		startTime := time.Now()

		testsToRun, err := precompileRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("precompiled contracts tests failed: %v", err)
		}

		if err := precompileRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("precompiled contracts tests failed: %v", err)
		}

		precompileRunner.Logger.Print("🍾 precompiled contracts tests completed in %s", time.Since(startTime).String())

		return err
	}
}
