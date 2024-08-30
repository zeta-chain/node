package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
)

// zetaTestRoutine runs Zeta transfer and message passing related e2e tests
func zetaTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		account := conf.AdditionalAccounts.UserZetaTest
		// initialize runner for zeta test
		zetaRunner, err := initTestRunner(
			"zeta",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgBlue, "zeta"),
		)
		if err != nil {
			return err
		}

		zetaRunner.Logger.Print("🏃 starting Zeta tests")
		startTime := time.Now()

		// funding the account
		txZetaSend := deployerRunner.SendZetaOnEvm(account.EVMAddress(), 1000)
		zetaRunner.WaitForTxReceiptOnEvm(txZetaSend)

		// depositing the necessary tokens on ZetaChain
		txZetaDeposit := zetaRunner.DepositZeta()
		txEtherDeposit := zetaRunner.DepositEther()
		zetaRunner.WaitForMinedCCTX(txZetaDeposit)
		zetaRunner.WaitForMinedCCTX(txEtherDeposit)

		// run zeta test
		testsToRun, err := zetaRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("zeta tests failed: %v", err)
		}

		if err := zetaRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("zeta tests failed: %v", err)
		}

		zetaRunner.Logger.Print("🍾 Zeta tests completed in %s", time.Since(startTime).String())

		return err
	}
}
