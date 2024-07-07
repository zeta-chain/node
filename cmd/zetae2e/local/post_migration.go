package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/zeta-chain/zetacore/e2e/config"
	"github.com/zeta-chain/zetacore/e2e/e2etests"
	"github.com/zeta-chain/zetacore/e2e/runner"
)

// postMigrationTestRoutine runs Bitcoin related e2e tests
func postMigrationTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		account := conf.AdditionalAccounts.UserBitcoin
		// initialize runner for bitcoin test
		postMigrationRunner, err := initTestRunner(
			"postMigration",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgMagenta, "postMigrationRunner"),
		)
		if err != nil {
			return err
		}

		postMigrationRunner.Logger.Print("🏃 starting postMigration tests")
		startTime := time.Now()

		// funding the account
		//txERC20Send := deployerRunner.SendERC20OnEvm(account.EVMAddress(), 1000)
		//postMigrationRunner.WaitForTxReceiptOnEvm(txERC20Send)
		//
		//// depositing the necessary tokens on ZetaChain
		//txEtherDeposit := postMigrationRunner.DepositEther(false)
		//txERC20Deposit := postMigrationRunner.DepositERC20()
		//
		//postMigrationRunner.WaitForMinedCCTX(txEtherDeposit)
		//postMigrationRunner.WaitForMinedCCTX(txERC20Deposit)
		//
		//postMigrationRunner.Name = "bitcoin"
		//postMigrationRunner.SetupBitcoinAccount(initBitcoinNetwork)
		//postMigrationRunner.Name = "postMigration"
		//postMigrationRunner.DepositBTC(testHeader)

		// run bitcoin test
		// Note: due to the extensive block generation in Bitcoin localnet, block header test is run first
		// to make it faster to catch up with the latest block header
		testsToRun, err := postMigrationRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("postMigrationRunner tests failed: %v", err)
		}

		if err := postMigrationRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("postMigrationRunner tests failed: %v", err)
		}

		if err := postMigrationRunner.CheckBtcTSSBalance(); err != nil {
			return err
		}

		postMigrationRunner.Logger.Print("🍾 PostMigration tests completed in %s", time.Since(startTime).String())

		return err
	}
}
