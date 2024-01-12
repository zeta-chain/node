package local

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/config"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/smoketests"
)

// bitcoinTestRoutine runs Bitcoin related smoke tests
func bitcoinTestRoutine(
	conf config.Config,
	deployerRunner *runner.SmokeTestRunner,
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
		printBtcSupply(bitcoinRunner)
		txUSDTSend := deployerRunner.SendUSDTOnEvm(UserBitcoinAddress, 1000)
		bitcoinRunner.WaitForTxReceiptOnEvm(txUSDTSend)
		printBtcSupply(bitcoinRunner)
		if !initBitcoinNetwork {
			bitcoinRunner.Logger.Print("sleeping 10sec...")
			time.Sleep(10 * time.Second)
			printBtcSupply(bitcoinRunner)
		}
		// depositing the necessary tokens on ZetaChain
		txEtherDeposit := bitcoinRunner.DepositEther(false)
		bitcoinRunner.WaitForMinedCCTX(txEtherDeposit)
		printBtcSupply(bitcoinRunner)
		txERC20Deposit := bitcoinRunner.DepositERC20()
		bitcoinRunner.WaitForMinedCCTX(txERC20Deposit)
		printBtcSupply(bitcoinRunner)
		bitcoinRunner.SetupBitcoinAccount(initBitcoinNetwork)
		if initBitcoinNetwork {
			printBtcSupply(bitcoinRunner)
			bitcoinRunner.DepositBTC(false)
			printBtcSupply(bitcoinRunner)
		}
		printBtcSupply(bitcoinRunner)
		if err := bitcoinRunner.CheckBtcTSSBalance(); err != nil {
			return err
		}

		// run bitcoin test
		// Note: due to the extensive block generation in Bitcoin localnet, block header test is run first
		// to make it faster to catch up with the latest block header
		if err := bitcoinRunner.RunSmokeTestsFromNames(
			smoketests.AllSmokeTests,
			//smoketests.TestBitcoinWithdrawName,
			//smoketests.TestSendZetaOutBTCRevertName,
			//smoketests.TestCrosschainSwapName,
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

func printBtcSupply(runner *runner.SmokeTestRunner) {
	zrc20Supply, err := runner.BTCZRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		runner.Logger.Print("failed to get btc supply: %v", err)
		return
	}
	runner.Logger.Print("btc supply: %v", zrc20Supply.String())
}
