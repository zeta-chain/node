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

// ethereumTestRoutine runs Ethereum related smoke tests
func ethereumTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
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
				err = fmt.Errorf("ethereum panic: %v, stack trace %s", r, stack[:n])
			}
		}()

		// initialize runner for ether test
		ethereumRunner, err := initTestRunner(
			"ether",
			conf,
			deployerRunner,
			UserEtherAddress,
			UserEtherPrivateKey,
			runner.NewLogger(verbose, color.FgMagenta, "ether"),
		)
		if err != nil {
			return err
		}

		ethereumRunner.Logger.Print("🏃 starting Ethereum tests")
		startTime := time.Now()

		// depositing the necessary tokens on ZetaChain
		txEtherDeposit := ethereumRunner.DepositEther(true)
		ethereumRunner.WaitForMinedCCTX(txEtherDeposit)

		// run ethereum test
		// Note: due to the extensive block generation in Ethereum localnet, block header test is run first
		// to make it faster to catch up with the latest block header
		if err := ethereumRunner.RunE2ETestsFromNames(
			e2etests.AllE2ETests,
			e2etests.TestEtherWithdrawName,
			e2etests.TestContextUpgradeName,
			e2etests.TestEtherDepositAndCallName,
			e2etests.TestDepositAndCallRefundName,
		); err != nil {
			return fmt.Errorf("ethereum tests failed: %v", err)
		}

		ethereumRunner.Logger.Print("🍾 Ethereum tests completed in %s", time.Since(startTime).String())

		return err
	}
}
