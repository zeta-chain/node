package local

import (
	"fmt"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/config"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/smoketests"
)

// zetaTestRoutine runs Zeta transfer and message passing related smoke tests
func zetaTestRoutine(
	conf config.Config,
	deployerRunner *runner.SmokeTestRunner,
	verbose bool,
) func() error {
	return func() (err error) {
		// return an error on panic
		// TODO: remove and instead return errors in the smoke tests
		// https://github.com/zeta-chain/node/issues/1500
		defer func() {
			if r := recover(); r != nil {
				// print stack trace
				stack := make([]byte, 4096)
				n := runtime.Stack(stack, false)
				err = fmt.Errorf("zeta panic: %v, stack trace %s", r, stack[:n])
			}
		}()

		// initialize runner for zeta test
		zetaRunner, err := initTestRunner(
			"zeta",
			conf,
			deployerRunner,
			UserZetaTestAddress,
			UserZetaTestPrivateKey,
			runner.NewLogger(verbose, color.FgBlue, "zeta"),
		)
		if err != nil {
			return err
		}

		zetaRunner.Logger.Print("🏃 starting Zeta tests")
		startTime := time.Now()

		// funding the account
		txZetaSend := deployerRunner.SendZetaOnEvm(UserZetaTestAddress, 1000)
		zetaRunner.WaitForTxReceiptOnEvm(txZetaSend)

		// depositing the necessary tokens on ZetaChain
		txZetaDeposit := zetaRunner.DepositZeta()
		txEtherDeposit := zetaRunner.DepositEther()
		zetaRunner.WaitForMinedCCTX(txZetaDeposit)
		zetaRunner.WaitForMinedCCTX(txEtherDeposit)

		// run zeta test
		if err := zetaRunner.RunSmokeTestsFromNames(
			smoketests.AllSmokeTests,
			smoketests.TestSendZetaOutName,
			smoketests.TestMessagePassingName,
			smoketests.TestMessagePassingRevertFailName,
			smoketests.TestMessagePassingRevertSuccessName,
		); err != nil {
			return fmt.Errorf("zeta tests failed: %v", err)
		}

		zetaRunner.Logger.Print("🍾 Zeta tests completed in %s", time.Since(startTime).String())

		return err
	}
}
