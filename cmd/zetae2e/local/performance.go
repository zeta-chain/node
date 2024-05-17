package local

// performance.go provides routines that run the stress tests for different actions (deposit, withdraw) to measure network performance
// Note: the routine provided here should not be used concurrently with other routines as these reuse the accounts of other routines

import (
	"fmt"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/e2etests"
	"github.com/zeta-chain/zetacore/e2e/runner"
)

// ethereumDepositPerformanceRoutine runs performance tests for Ether deposit
func ethereumDepositPerformanceRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
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
		r, err := initTestRunner(
			"ether",
			conf,
			deployerRunner,
			UserERC20Address,
			UserERC20PrivateKey,
			runner.NewLogger(verbose, color.FgHiMagenta, "perf_eth_deposit"),
		)
		if err != nil {
			return err
		}

		r.Logger.Print("🏃 starting Ethereum deposit performance tests")
		startTime := time.Now()

		tests, err := r.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("ethereum deposit performance test failed: %v", err)
		}

		if err := r.RunE2ETests(tests); err != nil {
			return fmt.Errorf("ethereum deposit performance test failed: %v", err)
		}

		r.Logger.Print("🍾 Ethereum deposit performance test completed in %s", time.Since(startTime).String())

		return err
	}
}

// ethereumWithdrawPerformanceRoutine runs performance tests for Ether withdraw
func ethereumWithdrawPerformanceRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
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
		r, err := initTestRunner(
			"ether",
			conf,
			deployerRunner,
			UserEtherAddress,
			UserEtherPrivateKey,
			runner.NewLogger(verbose, color.FgHiBlue, "perf_eth_withdraw"),
		)
		if err != nil {
			return err
		}

		r.Logger.Print("🏃 starting Ethereum withdraw performance tests")
		startTime := time.Now()

		// depositing the necessary tokens on ZetaChain
		txEtherDeposit := r.DepositEther(false)
		r.WaitForMinedCCTX(txEtherDeposit)

		tests, err := r.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("ethereum withdraw performance test failed: %v", err)
		}

		if err := r.RunE2ETests(tests); err != nil {
			return fmt.Errorf("ethereum withdraw performance test failed: %v", err)
		}

		r.Logger.Print("🍾 Ethereum withdraw performance test completed in %s", time.Since(startTime).String())

		return err
	}
}
