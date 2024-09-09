package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
)

// adminTestRoutine runs admin functions tests
func adminTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		account := conf.AdditionalAccounts.UserAdmin
		// initialize runner for erc20 advanced test
		adminRunner, err := initTestRunner(
			"admin",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgHiGreen, "admin"),
			runner.WithZetaTxServer(deployerRunner.ZetaTxServer),
		)
		if err != nil {
			return err
		}

		adminRunner.Logger.Print("🏃 starting admin tests")
		startTime := time.Now()

		// funding the account
		// we transfer around the total supply of Zeta to the admin for the chain migration test
		txZetaSend := deployerRunner.SendZetaOnEvm(account.EVMAddress(), 20_500_000_000)
		txERC20Send := deployerRunner.SendERC20OnEvm(account.EVMAddress(), 1000)
		adminRunner.WaitForTxReceiptOnEvm(txZetaSend)
		adminRunner.WaitForTxReceiptOnEvm(txERC20Send)

		// depositing the necessary tokens on ZetaChain
		txZetaDeposit := adminRunner.DepositZeta()
		txEtherDeposit := adminRunner.DepositEther()
		txERC20Deposit := adminRunner.DepositERC20()
		adminRunner.WaitForMinedCCTX(txZetaDeposit)
		adminRunner.WaitForMinedCCTX(txEtherDeposit)
		adminRunner.WaitForMinedCCTX(txERC20Deposit)

		// run erc20 advanced test
		testsToRun, err := adminRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("admin tests failed: %v", err)
		}

		if err := adminRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("admin tests failed: %v", err)
		}

		adminRunner.Logger.Print("🍾 admin tests completed in %s", time.Since(startTime).String())

		return err
	}
}
