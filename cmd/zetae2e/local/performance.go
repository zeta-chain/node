package local

// performance.go provides routines that run the stress tests for different actions (deposit, withdraw) to measure network performance
// Note: the routine provided here should not be used concurrently with other routines as these reuse the accounts of other routines

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/e2etests"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"runtime"
	"time"
)

// ethereumDepositPerformanceRoutine runs Ethereum withdraw stress tests
func ethereumDepositPerformanceRoutine(
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
				err = fmt.Errorf("ethereum deposit perf panic: %v, stack trace %s", r, stack[:n])
			}
		}()

		// initialize runner for ether test
		ethereumRunner, err := initTestRunner(
			"ether",
			conf,
			deployerRunner,
			UserEtherAddress,
			UserEtherPrivateKey,
			runner.NewLogger(verbose, color.FgMagenta, "perf_eth_deposit"),
		)
		if err != nil {
			return err
		}

		ethereumRunner.Logger.Print("🏃 starting Ethereum deposit performance tests")
		startTime := time.Now()

		if err := ethereumRunner.RunE2ETestsFromNames(
			e2etests.AllE2ETests,
			e2etests.TestStressEtherDepositName,
		); err != nil {
			return fmt.Errorf("thereum deposit performance test failed: %v", err)
		}

		ethereumRunner.Logger.Print("🍾 Ethereum deposit performance test completed in %s", time.Since(startTime).String())

		return err
	}
}

// ethereumWithdrawPerformanceRoutine runs Ethereum withdraw stress tests
func ethereumWithdrawPerformanceRoutine(
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
				err = fmt.Errorf("ethereum withdraw perf panic: %v, stack trace %s", r, stack[:n])
			}
		}()

		// initialize runner for ether test
		ethereumRunner, err := initTestRunner(
			"ether",
			conf,
			deployerRunner,
			UserEtherAddress,
			UserEtherPrivateKey,
			runner.NewLogger(verbose, color.FgMagenta, "perf_eth_withdraw"),
		)
		if err != nil {
			return err
		}

		ethereumRunner.Logger.Print("🏃 starting Ethereum withdraw performance tests")
		startTime := time.Now()

		// depositing the necessary tokens on ZetaChain
		txEtherDeposit := ethereumRunner.DepositEther(true)
		ethereumRunner.WaitForMinedCCTX(txEtherDeposit)

		if err := ethereumRunner.RunE2ETestsFromNames(
			e2etests.AllE2ETests,
			e2etests.TestStressEtherWithdrawName,
		); err != nil {
			return fmt.Errorf("thereum withdraw performance test failed: %v", err)
		}

		ethereumRunner.Logger.Print("🍾 Ethereum withdraw performance test completed in %s", time.Since(startTime).String())

		return err
	}
}
