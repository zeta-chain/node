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

// erc20TestRoutine runs erc20 related smoke tests
func erc20TestRoutine(
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
		txZetaSend := deployerRunner.SendZetaOnEvm(UserERC20Address, 1000)
		txUSDTSend := deployerRunner.SendUSDTOnEvm(UserERC20Address, 10)

		erc20Runner.WaitForTxReceiptOnEvm(txZetaSend)
		erc20Runner.WaitForTxReceiptOnEvm(txUSDTSend)

		// depositing the necessary tokens on ZetaChain
		txZetaDeposit := erc20Runner.DepositZeta()
		txEtherDeposit := erc20Runner.DepositEther(false)
		txERC20Deposit := erc20Runner.DepositERC20()
		erc20Runner.WaitForMinedCCTX(txZetaDeposit)
		erc20Runner.WaitForMinedCCTX(txEtherDeposit)
		erc20Runner.WaitForMinedCCTX(txERC20Deposit)

		//erc20Runner.SetupBitcoinAccount()
		//erc20Runner.DepositBTC()

		// run erc20 test
		if err := erc20Runner.RunSmokeTestsFromNames(
			smoketests.AllSmokeTests,
			smoketests.TestMultipleERC20DepositName,
			smoketests.TestWithdrawERC20Name,
			smoketests.TestMultipleWithdrawsName,
			smoketests.TestERC20DepositAndCallRefundName,
			smoketests.TestZRC20SwapName,
		); err != nil {
			return fmt.Errorf("erc20 tests failed: %v", err)
		}

		erc20Runner.Logger.Print("🍾 erc20 tests completed in %s", time.Since(startTime).String())

		return err
	}
}
