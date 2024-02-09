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

// miscTestRoutine runs miscellaneous smoke tests
func miscTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
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
				err = fmt.Errorf("misc panic: %v, stack trace %s", r, stack[:n])
			}
		}()

		// initialize runner for misc test
		miscRunner, err := initTestRunner(
			"misc",
			conf,
			deployerRunner,
			UserMiscAddress,
			UserMiscPrivateKey,
			runner.NewLogger(verbose, color.FgCyan, "misc"),
		)
		if err != nil {
			return err
		}

		miscRunner.Logger.Print("🏃 starting miscellaneous tests")
		startTime := time.Now()

		// funding the account
		txZetaSend := deployerRunner.SendZetaOnEvm(UserMiscAddress, 1000)
		miscRunner.WaitForTxReceiptOnEvm(txZetaSend)

		// depositing the necessary tokens on ZetaChain
		txZetaDeposit := miscRunner.DepositZeta()
		miscRunner.WaitForMinedCCTX(txZetaDeposit)

		// run misc test
		if err := miscRunner.RunE2ETestsFromNames(
			e2etests.AllE2ETests,
			//e2etests.TestBlockHeadersName,
			e2etests.TestMyTestName,
		); err != nil {
			return fmt.Errorf("misc tests failed: %v", err)
		}

		miscRunner.Logger.Print("🍾 miscellaneous tests completed in %s", time.Since(startTime).String())

		return err
	}
}
