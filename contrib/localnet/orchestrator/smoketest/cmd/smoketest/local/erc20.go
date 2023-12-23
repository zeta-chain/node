package local

import (
	"fmt"
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
				err = fmt.Errorf("erc20 panic: %v", r)
			}
		}()

		// initialize runner for erc20 test
		erc20Runner, err := initERC20Runner(conf, deployerRunner, verbose)
		if err != nil {
			panic(err)
		}

		erc20Runner.Logger.Print("🏃 starting erc20 tests")
		startTime := time.Now()

		// funding the account
		deployerRunner.SendZetaOnEvm(UserERC20Address, 1000)
		deployerRunner.SendUSDTOnEvm(UserERC20Address, 10)

		// depositing the necessary tokens on ZetaChain
		erc20Runner.DepositZeta()
		erc20Runner.DepositEther()
		erc20Runner.DepositERC20()

		// run erc20 test
		if err := erc20Runner.RunSmokeTestsFromNames(
			smoketests.AllSmokeTests,
			smoketests.TestDepositAndCallRefundName,
			//smoketests.TestMultipleERC20DepositName,
			smoketests.TestWithdrawERC20Name,
			//smoketests.TestMultipleWithdrawsName,
			//smoketests.TestERC20DepositAndCallRefundName,
		); err != nil {
			return fmt.Errorf("erc20 tests failed: %v", err)
		}

		erc20Runner.Logger.Print("🍾 erc20 tests completed in %s", time.Since(startTime).String())

		return err
	}
}

// erc20AdvancedTestRoutine runs erc20 advanced related smoke tests
func erc20AdvancedTestRoutine(
	conf config.Config,
	deployerRunner *runner.SmokeTestRunner,
	verbose bool,
) func() error {
	return func() (err error) {
		// return an error on panic
		// https://github.com/zeta-chain/node/issues/1500
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("erc20 advanced panic: %v", r)
			}
		}()

		// initialize runner for erc20 advanced test
		erc20AdvancedRunner, err := initERC20AdvancedRunner(conf, deployerRunner, verbose)
		if err != nil {
			panic(err)
		}

		erc20AdvancedRunner.Logger.Print("🏃 starting erc20 advanced tests")
		startTime := time.Now()

		// funding the account
		deployerRunner.SendZetaOnEvm(UserERC20AdvancedAddress, 1000)
		deployerRunner.SendUSDTOnEvm(UserERC20AdvancedAddress, 1000)

		// depositing the necessary tokens on ZetaChain
		erc20AdvancedRunner.DepositZeta()
		erc20AdvancedRunner.DepositEther()
		erc20AdvancedRunner.DepositERC20()

		// run erc20 advanced test
		if err := erc20AdvancedRunner.RunSmokeTestsFromNames(
			smoketests.AllSmokeTests,
			smoketests.TestZRC20SwapName,
			//smoketests.TestPauseZRC20Name,
			//smoketests.TestUpdateBytecodeName,
			//smoketests.TestWhitelistERC20Name,
		); err != nil {
			return fmt.Errorf("erc20 advanced tests failed: %v", err)
		}

		erc20AdvancedRunner.Logger.Print("🍾 erc20 advanced tests completed in %s", time.Since(startTime).String())

		return err
	}
}

// initERC20Runner initializes a runner for erc20 tests
func initERC20Runner(
	conf config.Config,
	deployerRunner *runner.SmokeTestRunner,
	verbose bool,
) (*runner.SmokeTestRunner, error) {
	// initialize runner for erc20 test
	erc20Runner, err := runnerFromConfig(
		conf,
		UserERC20Address,
		UserERC20PrivateKey,
		runner.NewLogger(verbose, color.FgGreen, "erc20"),
	)
	if err != nil {
		return nil, err
	}
	if err := erc20Runner.CopyAddressesFrom(deployerRunner); err != nil {
		return nil, err
	}
	return erc20Runner, nil
}

// initERC20AdvancedRunner initializes a runner for erc20 advanced tests
func initERC20AdvancedRunner(
	conf config.Config,
	deployerRunner *runner.SmokeTestRunner,
	verbose bool,
) (*runner.SmokeTestRunner, error) {
	// initialize runner for erc20 advanced test
	erc20AdvancedRunner, err := runnerFromConfig(
		conf,
		UserERC20AdvancedAddress,
		UserERC20AdvancedPrivateKey,
		runner.NewLogger(verbose, color.FgHiGreen, "erc20advanced"),
	)
	if err != nil {
		return nil, err
	}
	if err := erc20AdvancedRunner.CopyAddressesFrom(deployerRunner); err != nil {
		return nil, err
	}
	return erc20AdvancedRunner, nil
}
