package local

import (
	"fmt"
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
				err = fmt.Errorf("zeta panic: %v", r)
			}
		}()

		// initialize runner for zeta test
		zetaRunner, err := initZetaRunner(conf, deployerRunner, verbose)
		if err != nil {
			panic(err)
		}

		zetaRunner.Logger.Print("🏃 starting Zeta tests")
		startTime := time.Now()

		// funding the account
		deployerRunner.SendZetaOnEvm(UserZetaTestAddress, 1000)

		// depositing the necessary tokens on ZetaChain
		zetaRunner.DepositZeta()
		zetaRunner.DepositEther()

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

// initZetaRunner initializes a runner for zeta tests
func initZetaRunner(
	conf config.Config,
	deployerRunner *runner.SmokeTestRunner,
	verbose bool,
) (*runner.SmokeTestRunner, error) {
	// initialize runner for zeta test
	zetaRunner, err := runnerFromConfig(
		conf,
		UserZetaTestAddress,
		UserZetaTestPrivateKey,
		runner.NewLogger(verbose, color.FgBlue, "zeta"),
	)
	if err != nil {
		return nil, err
	}
	if err := zetaRunner.CopyAddressesFrom(deployerRunner); err != nil {
		return nil, err
	}
	return zetaRunner, nil
}
