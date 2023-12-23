package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/config"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/smoketests"
)

// miscTestRoutine runs miscellaneous smoke tests
func miscTestRoutine(
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
				err = fmt.Errorf("misc test panic: %v", r)
			}
		}()

		// initialize runner for misc test
		miscRunner, err := initMiscRunner(conf, deployerRunner, verbose)
		if err != nil {
			panic(err)
		}

		miscRunner.Logger.Print("🏃 starting miscellaneous tests")
		startTime := time.Now()

		// funding the account
		deployerRunner.SendZetaOnEvm(UserMiscAddress, 1000)

		// depositing the necessary tokens on ZetaChain
		miscRunner.DepositZeta()

		// run misc test
		if err := miscRunner.RunSmokeTestsFromNames(
			smoketests.AllSmokeTests,
			//smoketests.TestBlockHeadersName,
			smoketests.TestMyTestName,
		); err != nil {
			return fmt.Errorf("misc tests failed: %v", err)
		}

		miscRunner.Logger.Print("🍾 miscellaneous tests completed in %s", time.Since(startTime).String())

		return err
	}
}

// initMiscRunner initializes a runner for miscellaneous tests
func initMiscRunner(
	conf config.Config,
	deployerRunner *runner.SmokeTestRunner,
	verbose bool,
) (*runner.SmokeTestRunner, error) {
	// initialize runner for misc test
	miscRunner, err := runnerFromConfig(
		conf,
		UserMiscAddress,
		UserMiscPrivateKey,
		runner.NewLogger(verbose, color.FgCyan, "misc"),
	)
	if err != nil {
		return nil, err
	}
	if err := miscRunner.CopyAddressesFrom(deployerRunner); err != nil {
		return nil, err
	}
	return miscRunner, nil
}
