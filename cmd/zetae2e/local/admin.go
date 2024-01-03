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

// adminTestRoutine runs admin functions tests
func adminTestRoutine(
	conf config.Config,
	deployerRunner *runner.SmokeTestRunner,
	verbose bool,
) func() error {
	return func() (err error) {
		// return an error on panic
		// https://github.com/zeta-chain/node/issues/1500
		defer func() {
			if r := recover(); r != nil {
				// print stack trace
				stack := make([]byte, 4096)
				n := runtime.Stack(stack, false)
				err = fmt.Errorf("admin panic: %v, stack trace %s", r, stack[:n])
			}
		}()

		// initialize runner for erc20 advanced test
		adminRunner, err := initTestRunner(
			"admin",
			conf,
			deployerRunner,
			UserAdminAddress,
			UserAdminPrivateKey,
			runner.NewLogger(verbose, color.FgGreen, "admin"),
		)
		if err != nil {
			return err
		}

		adminRunner.Logger.Print("🏃 starting admin tests")
		startTime := time.Now()

		// funding the account
		txZetaSend := deployerRunner.SendZetaOnEvm(UserAdminAddress, 1000)
		txUSDTSend := deployerRunner.SendUSDTOnEvm(UserAdminAddress, 1000)
		adminRunner.WaitForTxReceiptOnEvm(txZetaSend)
		adminRunner.WaitForTxReceiptOnEvm(txUSDTSend)

		// depositing the necessary tokens on ZetaChain
		txZetaDeposit := adminRunner.DepositZeta()
		txEtherDeposit := adminRunner.DepositEther(false)
		txERC20Deposit := adminRunner.DepositERC20()
		adminRunner.WaitForMinedCCTX(txZetaDeposit)
		adminRunner.WaitForMinedCCTX(txEtherDeposit)
		adminRunner.WaitForMinedCCTX(txERC20Deposit)

		// run erc20 advanced test
		if err := adminRunner.RunSmokeTestsFromNames(
			smoketests.AllSmokeTests,
			smoketests.TestPauseZRC20Name,
			smoketests.TestUpdateBytecodeName,
			smoketests.TestDepositEtherLiquidityCapName,
		); err != nil {
			return fmt.Errorf("admim tests failed: %v", err)
		}

		adminRunner.Logger.Print("🍾 admin tests completed in %s", time.Since(startTime).String())

		return err
	}
}
