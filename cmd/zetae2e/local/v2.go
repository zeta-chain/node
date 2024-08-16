package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/e2etests"
	"github.com/zeta-chain/zetacore/e2e/runner"
)

// erc20TestRoutine runs v2 related e2e tests
// TODO: this routine will be broken down in the future and will replace most current tests
// we keep a single routine for v2 for simplicity
// https://github.com/zeta-chain/node/issues/2554
func v2TestRoutine(
	conf config.Config,
	name string,
	account config.Account,
	color color.Attribute,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		name = "v2-" + name

		// initialize runner for erc20 test
		v2Runner, err := initTestRunner(
			name,
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color, name),
			runner.WithZetaTxServer(deployerRunner.ZetaTxServer),
		)
		if err != nil {
			return err
		}

		v2Runner.Logger.Print("🏃 starting %s tests", name)
		startTime := time.Now()

		// funding the account
		txERC20Send := deployerRunner.SendERC20OnEvm(account.EVMAddress(), 10000)
		v2Runner.WaitForTxReceiptOnEvm(txERC20Send)

		// run erc20 test
		testsToRun, err := v2Runner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("%s tests failed: %v", name, err)
		}

		if err := v2Runner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("%s tests failed: %v", name, err)
		}

		v2Runner.Logger.Print("🍾 %s tests completed in %s", name, time.Since(startTime).String())

		return err
	}
}
