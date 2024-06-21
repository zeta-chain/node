package local

import (
	"fmt"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/zeta-chain/zetacore/e2e/e2etests"

	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/runner"
)

// adminTestRoutine runs admin functions tests
func migrationTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				// print stack trace
				stack := make([]byte, 4096)
				n := runtime.Stack(stack, false)
				err = fmt.Errorf("admin panic: %v, stack trace %s", r, stack[:n])
			}
		}()

		// initialize runner for erc20 advanced test
		migrationTestRunner, err := initTestRunner(
			"admin",
			conf,
			deployerRunner,
			UserAdminAddress,
			UserAdminPrivateKey,
			runner.NewLogger(verbose, color.FgHiWhite, "migration"),
			runner.WithZetaTxServer(deployerRunner.ZetaTxServer),
		)
		if err != nil {
			return err
		}

		migrationTestRunner.Logger.Print("🏃 starting migration tests")
		startTime := time.Now()

		// funding the account
		// we transfer around the total supply of Zeta to the admin for the chain migration test
		txZetaSend := deployerRunner.SendZetaOnEvm(UserAdminAddress, 20_500_000_000)
		txERC20Send := deployerRunner.SendERC20OnEvm(UserAdminAddress, 1000)
		migrationTestRunner.WaitForTxReceiptOnEvm(txZetaSend)
		migrationTestRunner.WaitForTxReceiptOnEvm(txERC20Send)

		// depositing the necessary tokens on ZetaChain
		txZetaDeposit := migrationTestRunner.DepositZeta()
		txEtherDeposit := migrationTestRunner.DepositEther(false)
		txERC20Deposit := migrationTestRunner.DepositERC20()
		migrationTestRunner.WaitForMinedCCTX(txZetaDeposit)
		migrationTestRunner.WaitForMinedCCTX(txEtherDeposit)
		migrationTestRunner.WaitForMinedCCTX(txERC20Deposit)

		// run erc20 advanced test
		testsToRun, err := migrationTestRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("migration tests failed: %v", err)
		}

		if err := migrationTestRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("migration tests failed: %v", err)
		}

		migrationTestRunner.Logger.Print("🍾 migration tests completed in %s", time.Since(startTime).String())

		return err
	}

}
