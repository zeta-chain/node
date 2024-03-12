package local

import (
	"fmt"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/e2etests"
	"github.com/zeta-chain/zetacore/e2e/runner"
)

// erc20TestRoutine runs erc20 related e2e tests
func erc20TestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		// return an error on panic
		// TODO: remove and instead return errors in the tests
		// https://github.com/zeta-chain/node/issues/1500
		defer func() {
			if r := recover(); r != nil {
				// print stack trace
				stack := make([]byte, 4096)
				n := runtime.Stack(stack, false)
				err = fmt.Errorf("erc20 panic: %v, stack trace %s", r, stack[:n])
			}
		}()

		// initialize runner for erc20 test
		erc20Runner, err := initTestRunner(
			"erc20",
			conf,
			deployerRunner,
			UserERC20Address,
			UserERC20PrivateKey,
			runner.NewLogger(verbose, color.FgGreen, "erc20"),
		)
		if err != nil {
			return err
		}

		erc20Runner.Logger.Print("🏃 starting erc20 tests")
		startTime := time.Now()

		// funding the account
		txERC20Send := deployerRunner.SendERC20OnEvm(UserERC20Address, 10)
		erc20Runner.WaitForTxReceiptOnEvm(txERC20Send)

		// depositing the necessary tokens on ZetaChain
		txEtherDeposit := erc20Runner.DepositEther(false)
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
