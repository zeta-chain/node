package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
)

// zevmMPTestRoutine runs ZEVM message passing related e2e tests
func zevmMPTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		account := conf.AdditionalAccounts.UserZEVMMPTest
		// initialize runner for zevm mp test
		zevmMPRunner, err := initTestRunner(
			"zevm_mp",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgHiRed, "zevm_mp"),
		)
		if err != nil {
			return err
		}

		zevmMPRunner.Logger.Print("🏃 starting ZEVM Message Passing tests")
		startTime := time.Now()

		// funding the account
		txZetaSend := deployerRunner.SendZetaOnEvm(account.EVMAddress(), 1000)
		zevmMPRunner.WaitForTxReceiptOnEvm(txZetaSend)

		// depositing the necessary tokens on ZetaChain
		txZetaDeposit := zevmMPRunner.DepositZeta()
		txEtherDeposit := zevmMPRunner.DepositEther()
		zevmMPRunner.WaitForMinedCCTX(txZetaDeposit)
		zevmMPRunner.WaitForMinedCCTX(txEtherDeposit)

		// run zevm message passing test
		testsToRun, err := zevmMPRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("zevm message passing tests failed: %v", err)
		}

		if err := zevmMPRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("zevm message passing tests failed: %v", err)
		}

		zevmMPRunner.Logger.Print("🍾 ZEVM message passing tests completed in %s", time.Since(startTime).String())

		return err
	}
}
