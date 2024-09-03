package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"golang.org/x/sync/errgroup"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
)

// startV2Tests starts v2 related tests in parallel
func startV2Tests(eg *errgroup.Group, conf config.Config, deployerRunner *runner.E2ERunner, verbose bool) {
	// Test happy paths for gas token workflow
	eg.Go(v2TestRoutine(conf, "eth", conf.AdditionalAccounts.UserV2Ether, color.FgHiGreen, deployerRunner, verbose,
		e2etests.TestV2ETHDepositName,
		e2etests.TestV2ETHDepositAndCallName,
		e2etests.TestV2ETHWithdrawName,
		e2etests.TestV2ETHWithdrawAndCallName,
		e2etests.TestV2ZEVMToEVMCallName,
		e2etests.TestV2EVMToZEVMCallName,
	))

	// Test happy paths for erc20 token workflow
	eg.Go(v2TestRoutine(conf, "erc20", conf.AdditionalAccounts.UserV2ERC20, color.FgHiBlue, deployerRunner, verbose,
		e2etests.TestV2ETHDepositName, // necessary to pay fees on ZEVM
		e2etests.TestV2ERC20DepositName,
		e2etests.TestV2ERC20DepositAndCallName,
		e2etests.TestV2ERC20WithdrawName,
		e2etests.TestV2ERC20WithdrawAndCallName,
	))

	// Test revert cases for gas token workflow
	eg.Go(
		v2TestRoutine(
			conf,
			"eth-revert",
			conf.AdditionalAccounts.UserV2EtherRevert,
			color.FgHiYellow,
			deployerRunner,
			verbose,
			e2etests.TestV2ETHDepositName, // necessary to pay fees on ZEVM and withdraw
			e2etests.TestV2ETHDepositAndCallRevertName,
			e2etests.TestV2ETHDepositAndCallRevertWithCallName,
			e2etests.TestV2ETHWithdrawAndCallRevertName,
			e2etests.TestV2ETHWithdrawAndCallRevertWithCallName,
		),
	)

	// Test revert cases for erc20 token workflow
	eg.Go(
		v2TestRoutine(
			conf,
			"erc20-revert",
			conf.AdditionalAccounts.UserV2ERC20Revert,
			color.FgHiRed,
			deployerRunner,
			verbose,
			e2etests.TestV2ETHDepositName,   // necessary to pay fees on ZEVM
			e2etests.TestV2ERC20DepositName, // necessary to have assets to withdraw
			e2etests.TestOperationAddLiquidityETHName, // liquidity with gas and ERC20 are necessary for reverts
			e2etests.TestOperationAddLiquidityERC20Name,
			e2etests.TestV2ERC20DepositAndCallRevertName,
			e2etests.TestV2ERC20DepositAndCallRevertWithCallName,
			e2etests.TestV2ERC20WithdrawAndCallRevertName,
			e2etests.TestV2ERC20WithdrawAndCallRevertWithCallName,
		),
	)
}

// v2TestRoutine runs v2 related e2e tests
// TODO: this routine will be broken down in the future and will replace most current tests
// we keep a single routine for v2 for simplicity
// https://github.com/zeta-chain/node/issues/2554
func v2TestRoutine(
	conf config.Config,
	name string,
	account config.Account,
	color color.Attribute,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		name = "v2-" + name

		// initialize runner for erc20 test
		v2Runner, err := initTestRunner(
			name,
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color, name),
			runner.WithZetaTxServer(deployerRunner.ZetaTxServer),
		)
		if err != nil {
			return err
		}

		v2Runner.Logger.Print("🏃 starting %s tests", name)
		startTime := time.Now()

		// funding the account
		txERC20Send := deployerRunner.SendERC20OnEvm(account.EVMAddress(), 10000)
		v2Runner.WaitForTxReceiptOnEvm(txERC20Send)

		// run erc20 test
		testsToRun, err := v2Runner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("%s tests failed: %v", name, err)
		}

		if err := v2Runner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("%s tests failed: %v", name, err)
		}

		v2Runner.Logger.Print("🍾 %s tests completed in %s", name, time.Since(startTime).String())

		return err
	}
}
