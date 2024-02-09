package local

import (
	"fmt"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/config"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/e2etests"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
)

// bitcoinTestRoutine runs Bitcoin related smoke tests
func bitcoinTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	initBitcoinNetwork bool,
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
				err = fmt.Errorf("bitcoin panic: %v, stack trace %s", r, stack[:n])
			}
		}()

		// initialize runner for bitcoin test
		bitcoinRunner, err := initTestRunner(
			"bitcoin",
			conf,
			deployerRunner,
			UserBitcoinAddress,
			UserBitcoinPrivateKey,
			runner.NewLogger(verbose, color.FgYellow, "bitcoin"),
		)
		if err != nil {
			return err
		}

		bitcoinRunner.Logger.Print("🏃 starting Bitcoin tests")
		startTime := time.Now()

		// funding the account
		txUSDTSend := deployerRunner.SendUSDTOnEvm(UserBitcoinAddress, 1000)
		bitcoinRunner.WaitForTxReceiptOnEvm(txUSDTSend)

		// depositing the necessary tokens on ZetaChain
		txEtherDeposit := bitcoinRunner.DepositEther(false)
		txERC20Deposit := bitcoinRunner.DepositERC20()

		bitcoinRunner.WaitForMinedCCTX(txEtherDeposit)
		bitcoinRunner.WaitForMinedCCTX(txERC20Deposit)

		bitcoinRunner.SetupBitcoinAccount(initBitcoinNetwork)
		bitcoinRunner.DepositBTC(true)

		// run bitcoin test
		// Note: due to the extensive block generation in Bitcoin localnet, block header test is run first
		// to make it faster to catch up with the latest block header
		if err := bitcoinRunner.RunE2ETestsFromNames(
			e2etests.AllE2ETests,
			e2etests.TestBitcoinWithdrawName,
			e2etests.TestZetaWithdrawBTCRevertName,
			e2etests.TestCrosschainSwapName,
		); err != nil {
			return fmt.Errorf("bitcoin tests failed: %v", err)
		}

		if err := bitcoinRunner.CheckBtcTSSBalance(); err != nil {
			return err
		}

		bitcoinRunner.Logger.Print("🍾 Bitcoin tests completed in %s", time.Since(startTime).String())

		return err
	}
}
