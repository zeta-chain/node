package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
)

// ethereumTestRoutine runs Ethereum related e2e tests
func ethereumTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		// initialize runner for ether test
		ethereumRunner, err := initTestRunner(
			"ether",
			conf,
			deployerRunner,
			conf.AdditionalAccounts.UserEther,
			runner.NewLogger(verbose, color.FgMagenta, "ether"),
		)
		if err != nil {
			return err
		}

		ethereumRunner.Logger.Print("🏃 starting Ethereum tests")
		startTime := time.Now()

		// depositing the necessary tokens on ZetaChain
		txEtherDeposit := ethereumRunner.DepositEther()
		ethereumRunner.WaitForMinedCCTX(txEtherDeposit)

		// run ethereum test
		// Note: due to the extensive block generation in Ethereum localnet, block header test is run first
		// to make it faster to catch up with the latest block header
		testsToRun, err := ethereumRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("ethereum tests failed: %v", err)
		}

		if err := ethereumRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("ethereum tests failed: %v", err)
		}

		ethereumRunner.Logger.Print("🍾 Ethereum tests completed in %s", time.Since(startTime).String())

		return err
	}
}
