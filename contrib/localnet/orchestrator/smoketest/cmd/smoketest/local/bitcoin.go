package local

import (
	"fmt"
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
) func() error {
	return func() (err error) {
		// return an error on panic
		// TODO: remove and instead return errors in the smoke tests
		// https://github.com/zeta-chain/node/issues/1500
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("bitcoin panic: %v", r)
			}
		}()

		// initialize runner for bitcoin test
		bitcoinRunner, err := initBitcoinRunner(conf, deployerRunner, verbose)
		if err != nil {
			panic(err)
		}

		bitcoinRunner.Logger.Print("🏃 starting Bitcoin tests")
		startTime := time.Now()

		// funding the account
		deployerRunner.SendZetaOnEvm(UserBitcoinAddress, 1000)

		// depositing the necessary tokens on ZetaChain
		bitcoinRunner.DepositZeta()
		bitcoinRunner.DepositEther()
		bitcoinRunner.SetupBitcoinAccount()
		bitcoinRunner.DepositBTC()
		bitcoinRunner.SetupZEVMSwapApp()

		// run bitcoin test
		if err := bitcoinRunner.RunSmokeTestsFromNames(
			smoketests.AllSmokeTests,
			smoketests.TestBitcoinWithdrawName,
			smoketests.TestSendZetaOutBTCRevertName,
			//smoketests.TestCrosschainSwapName,
		); err != nil {
			return fmt.Errorf("bitcoin tests failed: %v", err)
		}

		bitcoinRunner.Logger.Print("🍾 Bitcoin tests completed in %s", time.Since(startTime).String())

		return err
	}
}

// initBitcoinRunner initializes a runner for bitcoin tests
func initBitcoinRunner(
	conf config.Config,
	deployerRunner *runner.SmokeTestRunner,
	verbose bool,
) (*runner.SmokeTestRunner, error) {
	// initialize runner for bitcoin test
	bitcoinRunner, err := runnerFromConfig(
		conf,
		UserBitcoinAddress,
		UserBitcoinPrivateKey,
		runner.NewLogger(verbose, color.FgYellow, "bitcoin"),
	)
	if err != nil {
		return nil, err
	}
	if err := bitcoinRunner.CopyAddressesFrom(deployerRunner); err != nil {
		return nil, err
	}
	return bitcoinRunner, nil
}
