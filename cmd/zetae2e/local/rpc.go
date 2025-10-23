package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
)

// rpcTestRoutine runs zevm json rpc tests
func rpcTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		account := conf.AdditionalAccounts.UserRPC
		// initialize runner for rpc test
		rpcTestRunner, err := initTestRunner(
			"rpc",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgMagenta, "rpc"),
			runner.WithZetaTxServer(deployerRunner.ZetaTxServer),
		)
		if err != nil {
			return err
		}

		rpcTestRunner.Logger.Print("🏃 starting RPC tests")
		startTime := time.Now()

		// run RPC test
		testsToRun, err := rpcTestRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("RPC tests failed: %v", err)
		}

		if err := rpcTestRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("RPC tests failed: %v", err)
		}
		rpcTestRunner.Logger.Print("🍾 RPC tests completed in %s", time.Since(startTime).String())

		return nil
	}
}
