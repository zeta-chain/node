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
		txZetaSend := deployerRunner.TransferZETAOnEvm(account.EVMAddress(), 20_500_000_000)
		txERC20Send := deployerRunner.SendERC20OnEVM(account.EVMAddress(), 1000)
		adminRunner.WaitForTxReceiptOnEVM(txZetaSend)
		adminRunner.WaitForTxReceiptOnEVM(txERC20Send)

		// depositing the necessary tokens on ZetaChain to the deployer account
		// only deposit ZETA if V2 ZETA flows are enabled (gateway deposits don't work otherwise)
		if adminRunner.IsV2ZETAEnabled() {
			txZetaDeposit := adminRunner.DepositZETAToDeployer()
			adminRunner.WaitForMinedCCTX(txZetaDeposit.Hash())
		}
		txEtherDeposit := adminRunner.DepositEtherToDeployer()
		txERC20Deposit := adminRunner.DepositERC20ToDeployer()
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
		deployerRunner.ERC20CustodyAddr = adminRunner.ERC20CustodyAddr // update the address as ERC20MigrateFunds migrates funds to the new address

		return err
	}
}
