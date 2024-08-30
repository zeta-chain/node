package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
)

// erc20TestRoutine runs erc20 related e2e tests
func erc20TestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		account := conf.AdditionalAccounts.UserERC20
		// initialize runner for erc20 test
		erc20Runner, err := initTestRunner(
			"erc20",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgGreen, "erc20"),
			runner.WithZetaTxServer(deployerRunner.ZetaTxServer),
		)
		if err != nil {
			return err
		}

		erc20Runner.Logger.Print("🏃 starting erc20 tests")
		startTime := time.Now()

		// funding the account
		txERC20Send := deployerRunner.SendERC20OnEvm(account.EVMAddress(), 10)
		erc20Runner.WaitForTxReceiptOnEvm(txERC20Send)

		// depositing the necessary tokens on ZetaChain
		txEtherDeposit := erc20Runner.DepositEther()
		txERC20Deposit := erc20Runner.DepositERC20()
		erc20Runner.WaitForMinedCCTX(txEtherDeposit)
		erc20Runner.WaitForMinedCCTX(txERC20Deposit)

		// run erc20 test
		testsToRun, err := erc20Runner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("erc20 tests failed: %v", err)
		}

		if err := erc20Runner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("erc20 tests failed: %v", err)
		}

		erc20Runner.Logger.Print("🍾 erc20 tests completed in %s", time.Since(startTime).String())

		return err
	}
}
