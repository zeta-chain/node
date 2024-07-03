package local

import (
	"fmt"
	"runtime"
	"time"

	"github.com/fatih/color"

	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/e2etests"
	"github.com/zeta-chain/zetacore/e2e/runner"
)

// bitcoinTestRoutine runs Bitcoin related e2e tests
func solanaTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	initBitcoinNetwork bool,
	testHeader bool,
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
				err = fmt.Errorf("solana panic: %v, stack trace %s", r, stack[:n])
			}
		}()
		account := conf.AdditionalAccounts.UserBitcoin

		// initialize runner for bitcoin test
		solanaRunner, err := initTestRunner(
			"solana",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgCyan, "solana"),
		)

		if err != nil {
			return err
		}

		solanaRunner.Logger.Print("🏃 starting Solana tests")
		startTime := time.Now()

		// run bitcoin test
		// Note: due to the extensive block generation in Bitcoin localnet, block header test is run first
		// to make it faster to catch up with the latest block header
		testsToRun, err := solanaRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("solana tests failed: %v", err)
		}

		if err := solanaRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("solana tests failed: %v", err)
		}

		if err := solanaRunner.CheckBtcTSSBalance(); err != nil {
			return err
		}

		solanaRunner.Logger.Print("🍾 Solana tests completed in %s", time.Since(startTime).String())

		return err
	}
}
