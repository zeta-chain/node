package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
)

// suiTestRoutine runs Sui related e2e tests
func suiTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		// initialize runner for sui test
		suiRunner, err := initTestRunner(
			"sui",
			conf,
			deployerRunner,
			conf.AdditionalAccounts.UserSui,
			runner.NewLogger(verbose, color.FgHiCyan, "sui"),
			runner.WithZetaTxServer(deployerRunner.ZetaTxServer),
		)
		if err != nil {
			return err
		}

		suiRunner.Logger.Print("🏃 starting Sui tests")
		startTime := time.Now()

		suiRunnerSigner, err := suiRunner.Account.SuiSigner()
		if err != nil {
			return err
		}

		// get tokens for the account
		suiRunner.RequestSuiFromFaucet(conf.RPCs.SuiFaucet, suiRunnerSigner.Address())

		// mint fungible tokens to the account
		txRes := deployerRunner.SuiMintUSDC("100000000000", suiRunnerSigner.Address())

		deployerRunner.Logger.Info("Sui USDC mint tx: %s", txRes.Digest)

		// run sui test
		testsToRun, err := suiRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("sui tests failed: %v", err)
		}

		// update Sui gateway package information after upgrade
		// because the old gateway package ID is deprecated after upgrade
		if suiRunner.IsRunningUpgrade() {
			suiRunner.SuiUpdateGatewayInfo()
		}

		if err := suiRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("sui tests failed: %v", err)
		}

		// check gateway SUI balance against ZRC20 total supply
		// TODO: https://github.com/zeta-chain/node/issues/4139 remove after v35
		// skip if its running upgrade as gas limit improvement is not on release/v32
		if !suiRunner.IsRunningUpgrade() {
			suiRunner.CheckSUITSSBalance()
		}

		suiRunner.Logger.Print("🍾 sui tests completed in %s", time.Since(startTime).String())

		return err
	}
}
