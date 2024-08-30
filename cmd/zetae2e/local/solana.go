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
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		// initialize runner for solana test
		solanaRunner, err := initTestRunner(
			"solana",
			conf,
			deployerRunner,
			conf.AdditionalAccounts.UserSolana,
			runner.NewLogger(verbose, color.FgCyan, "solana"),
		)
		if err != nil {
			return err
		}

		solanaRunner.Logger.Print("🏃 starting Solana tests")
		startTime := time.Now()
		solanaRunner.SetupSolanaAccount()

		// run solana test
		testsToRun, err := solanaRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("solana tests failed: %v", err)
		}

		if err := solanaRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("solana tests failed: %v", err)
		}

		// check gateway SOL balance against ZRC20 total supply
		if err := solanaRunner.CheckSolanaTSSBalance(); err != nil {
			return err
		}

		solanaRunner.Logger.Print("🍾 solana tests completed in %s", time.Since(startTime).String())

		return err
	}
}
