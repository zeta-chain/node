package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
)

// legacyERC20TestRoutine runs erc20 related e2e tests
func legacyERC20TestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		account := conf.AdditionalAccounts.UserLegacyERC20
		// initialize runner for erc20 test
		erc20Runner, err := initTestRunner(
			"erc20",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgGreen, "erc20"),
			runner.WithZetaTxServer(deployerRunner.ZetaTxServer),
		)
		if err != nil {
			return err
		}

		erc20Runner.Logger.Print("🏃 starting erc20 tests")
		startTime := time.Now()

		// funding the account
		txERC20Send := deployerRunner.SendERC20OnEVM(account.EVMAddress(), 10000)
		erc20Runner.WaitForTxReceiptOnEVM(txERC20Send)

		// depositing the necessary tokens on ZetaChain
		txEtherDeposit := erc20Runner.LegacyDepositEther()
		txERC20Deposit := erc20Runner.LegacyDepositERC20()
		erc20Runner.WaitForMinedCCTX(txEtherDeposit)
		erc20Runner.WaitForMinedCCTX(txERC20Deposit)

		// run erc20 test
		testsToRun, err := erc20Runner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("erc20 tests failed: %v", err)
		}

		if err := erc20Runner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("erc20 tests failed: %v", err)
		}

		erc20Runner.Logger.Print("🍾 erc20 tests completed in %s", time.Since(startTime).String())

		return err
	}
}

// legacyEthereumTestRoutine runs Ethereum related e2e tests
func legacyEthereumTestRoutine(
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
			conf.AdditionalAccounts.UserLegacyEther,
			runner.NewLogger(verbose, color.FgMagenta, "ether"),
		)
		if err != nil {
			return err
		}

		ethereumRunner.Logger.Print("🏃 starting Ethereum tests")
		startTime := time.Now()

		// depositing the necessary tokens on ZetaChain
		txEtherDeposit := ethereumRunner.LegacyDepositEther()
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

// legacyZEVMMPTestRoutine runs ZEVM message passing related e2e tests
func legacyZEVMMPTestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		account := conf.AdditionalAccounts.UserLegacyZEVMMP
		// initialize runner for zevm mp test
		zevmMPRunner, err := initTestRunner(
			"zevm_mp",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgHiRed, "zevm_mp"),
		)
		if err != nil {
			return err
		}

		zevmMPRunner.Logger.Print("🏃 starting ZEVM Message Passing tests")
		startTime := time.Now()

		// funding the account
		txZetaSend := deployerRunner.TransferZETAOnEvm(account.EVMAddress(), 1000)
		zevmMPRunner.WaitForTxReceiptOnEVM(txZetaSend)

		// depositing the necessary tokens on ZetaChain
		txZetaDeposit := zevmMPRunner.LegacyDepositZeta()
		txEtherDeposit := zevmMPRunner.LegacyDepositEther()
		zevmMPRunner.WaitForMinedCCTX(txZetaDeposit)
		zevmMPRunner.WaitForMinedCCTX(txEtherDeposit)

		// run zevm message passing test
		testsToRun, err := zevmMPRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("zevm message passing tests failed: %v", err)
		}

		if err := zevmMPRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("zevm message passing tests failed: %v", err)
		}

		zevmMPRunner.Logger.Print("🍾 ZEVM message passing tests completed in %s", time.Since(startTime).String())

		return err
	}
}

// legacyZETATestRoutine runs Zeta transfer and message passing related e2e tests
func legacyZETATestRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		account := conf.AdditionalAccounts.UserLegacyZeta
		// initialize runner for zeta test
		zetaRunner, err := initTestRunner(
			"zeta",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgBlue, "zeta"),
		)
		if err != nil {
			return err
		}

		zetaRunner.Logger.Print("🏃 starting Zeta tests")
		startTime := time.Now()

		// funding the account
		txZetaSend := deployerRunner.TransferZETAOnEvm(account.EVMAddress(), 1000)
		zetaRunner.WaitForTxReceiptOnEVM(txZetaSend)

		// depositing the necessary tokens on ZetaChain
		txZetaDeposit := zetaRunner.LegacyDepositZeta()
		txEtherDeposit := zetaRunner.LegacyDepositEther()
		zetaRunner.WaitForMinedCCTX(txZetaDeposit)
		zetaRunner.WaitForMinedCCTX(txEtherDeposit)

		// run zeta test
		testsToRun, err := zetaRunner.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("zeta tests failed: %v", err)
		}

		if err := zetaRunner.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("zeta tests failed: %v", err)
		}

		zetaRunner.Logger.Print("🍾 Zeta tests completed in %s", time.Since(startTime).String())

		return err
	}
}
