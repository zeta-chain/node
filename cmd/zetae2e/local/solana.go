package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
)

// solanaTestRoutine runs Solana related e2e tests
func solanaTestRoutine(
	conf config.Config,
	name string,
	account config.Account,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		// initialize runner for solana test
		solanaRunner, err := initTestRunner(
			name,
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgCyan, name),
			runner.WithZetaTxServer(deployerRunner.ZetaTxServer),
		)
		if err != nil {
			return err
		}

		solanaRunner.Logger.Print("🏃 starting %s tests", name)
		startTime := time.Now()
		solanaRunner.SetupSolanaAccount()

		// run solana test
		testsToRun, err := solanaRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("%s tests failed: %v", name, err)
		}

		if err := solanaRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("%s tests failed: %v", name, err)
		}

		// check gateway SOL balance against ZRC20 total supply
		solanaRunner.CheckSolanaTSSBalance()

		solanaRunner.Logger.Print("🍾 %s tests completed in %s", name, time.Since(startTime).String())

		return err
	}
}
