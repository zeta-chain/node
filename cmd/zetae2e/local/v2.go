package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/e2etests"
	"github.com/zeta-chain/zetacore/e2e/runner"
)

// erc20TestRoutine runs v2 related e2e tests
// TODO: this routine will be broken down in the future and will replace most current tests
// we keep a single routine for v2 for simplicity
// https://github.com/zeta-chain/node/issues/2554
func v2TestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		account := conf.AdditionalAccounts.UserERC20
		// initialize runner for erc20 test
		v2Runner, err := initTestRunner(
			"v2",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgHiYellow, "v2"),
			runner.WithZetaTxServer(deployerRunner.ZetaTxServer),
		)
		if err != nil {
			return err
		}

		v2Runner.Logger.Print("🏃 starting v2 tests")
		startTime := time.Now()

		// funding the account
		txERC20Send := deployerRunner.SendERC20OnEvm(account.EVMAddress(), 10)
		v2Runner.WaitForTxReceiptOnEvm(txERC20Send)

		// depositing the necessary tokens on ZetaChain
		// TODO: update with v2 deposits
		// https://github.com/zeta-chain/node/issues/2554
		txEtherDeposit := v2Runner.DepositEther(false)
		txERC20Deposit := v2Runner.DepositERC20()
		v2Runner.WaitForMinedCCTX(txEtherDeposit)
		v2Runner.WaitForMinedCCTX(txERC20Deposit)

		// run erc20 test
		testsToRun, err := v2Runner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("v2 tests failed: %v", err)
		}

		if err := v2Runner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("v2 tests failed: %v", err)
		}

		v2Runner.Logger.Print("🍾 v2 tests completed in %s", time.Since(startTime).String())

		return err
	}
}
