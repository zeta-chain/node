package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/config"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/smoketests"
)

// ethereumTestRoutine runs Ethereum related smoke tests
func ethereumTestRoutine(
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
				err = fmt.Errorf("ethereum panic: %v", r)
			}
		}()

		// initialize runner for ether test
		ethereumRunner, err := initEtherRunner(conf, deployerRunner, verbose)
		if err != nil {
			panic(err)
		}

		ethereumRunner.Logger.Print("🏃 starting Ethereum tests")
		startTime := time.Now()

		// funding the account
		deployerRunner.SendZetaOnEvm(UserEtherAddress, 1000)

		// depositing the necessary tokens on ZetaChain
		ethereumRunner.DepositZeta()
		ethereumRunner.DepositEther()
		ethereumRunner.SetupContextApp()

		// run ethereum test
		if err := ethereumRunner.RunSmokeTestsFromNames(
			smoketests.AllSmokeTests,
			smoketests.TestContextUpgradeName,
			smoketests.TestEtherDepositAndCallName,
			//smoketests.TestDepositEtherLiquidityCapName,
		); err != nil {
			return fmt.Errorf("ethereum tests failed: %v", err)
		}

		ethereumRunner.Logger.Print("🍾 Ethereum tests completed in %s", time.Since(startTime).String())

		return err
	}
}

// initEtherRunner initializes a runner for ether tests
func initEtherRunner(
	conf config.Config,
	deployerRunner *runner.SmokeTestRunner,
	verbose bool,
) (*runner.SmokeTestRunner, error) {
	// initialize runner for ether test
	etherRunner, err := runnerFromConfig(
		conf,
		UserEtherAddress,
		UserEtherPrivateKey,
		runner.NewLogger(verbose, color.FgMagenta, "ether"),
	)
	if err != nil {
		return nil, err
	}
	if err := etherRunner.CopyAddressesFrom(deployerRunner); err != nil {
		return nil, err
	}
	return etherRunner, nil
}
