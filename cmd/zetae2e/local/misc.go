package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
)

// miscTestRoutine runs miscellaneous e2e tests
func miscTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		account := conf.AdditionalAccounts.UserMisc
		// initialize runner for misc test
		miscRunner, err := initTestRunner(
			"misc",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgCyan, "misc"),
		)
		if err != nil {
			return err
		}

		miscRunner.Logger.Print("🏃 starting miscellaneous tests")
		startTime := time.Now()

		// funding the account
		txZetaSend := deployerRunner.SendZetaOnEvm(account.EVMAddress(), 1000)
		miscRunner.WaitForTxReceiptOnEvm(txZetaSend)

		// depositing the necessary tokens on ZetaChain
		txZetaDeposit := miscRunner.DepositZeta()
		miscRunner.WaitForMinedCCTX(txZetaDeposit)

		// run misc test
		testsToRun, err := miscRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("misc tests failed: %v", err)
		}

		if err := miscRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("misc tests failed: %v", err)
		}

		miscRunner.Logger.Print("🍾 miscellaneous tests completed in %s", time.Since(startTime).String())

		return err
	}
}
