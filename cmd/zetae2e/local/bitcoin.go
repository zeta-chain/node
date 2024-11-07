package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/testutil"
)

// initBitcoinTestRunners initializes Bitcoin deposit and withdraw test runners
func initBitcoinTestRunners(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	initNetwork bool,
	depositTests []string,
	withdrawTests []string,
) (func() error, func() error) {
	// initialize runner for deposit tests
	// deposit tests need Bitcoin node wallet to handle UTXOs
	account := conf.AdditionalAccounts.UserBitcoinDeposit
	runnerDeposit := initRunner(
		"btc_deposit",
		account,
		conf,
		deployerRunner,
		color.FgYellow,
		verbose,
		initNetwork,
		true,
	)

	// initialize runner for withdraw tests
	// withdraw tests DON'T use Bitcoin node wallet
	account = conf.AdditionalAccounts.UserBitcoinWithdraw
	runnerWithdraw := initRunner(
		"btc_withdraw",
		account,
		conf,
		deployerRunner,
		color.FgHiYellow,
		verbose,
		initNetwork,
		false,
	)

	// initialize funds
	// send BTC to TSS for gas fees and to tester ZEVM address
	if initNetwork {
		// mine 101 blocks to ensure the BTC rewards are spendable
		// Note: the block rewards can be sent to any address in here
		_, err := runnerDeposit.GenerateToAddressIfLocalBitcoin(101, runnerDeposit.BTCDeployerAddress)
		require.NoError(runnerDeposit, err)

		// send BTC to ZEVM addresses
		runnerDeposit.DepositBTC(runnerDeposit.EVMAddress())
		runnerDeposit.DepositBTC(runnerWithdraw.EVMAddress())
	}

	// create test routines
	routineDeposit := createTestRoutine(runnerDeposit, depositTests)
	routineWithdraw := createTestRoutine(runnerWithdraw, withdrawTests)

	return routineDeposit, routineWithdraw
}

// initRunner initializes the runner for given test name and account
func initRunner(
	name string,
	account config.Account,
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	printColor color.Attribute,
	verbose, initNetwork, createWallet bool,
) *runner.E2ERunner {
	// initialize runner for bitcoin test
	runner, err := initTestRunner(name, conf, deployerRunner, account, runner.NewLogger(verbose, printColor, name))
	testutil.NoError(err)

	// setup TSS address and setup deployer wallet
	runner.SetupBitcoinAccounts(createWallet)

	// initialize funds
	if initNetwork {
		// send some BTC block rewards to the deployer address
		_, err = runner.GenerateToAddressIfLocalBitcoin(4, runner.BTCDeployerAddress)
		require.NoError(runner, err)

		// send ERC20 token on EVM
		txERC20Send := deployerRunner.SendERC20OnEvm(account.EVMAddress(), 1000)
		runner.WaitForTxReceiptOnEvm(txERC20Send)

		// deposit ETH and ERC20 tokens on ZetaChain
		txEtherDeposit := runner.DepositEther()
		txERC20Deposit := runner.DepositERC20()

		runner.WaitForMinedCCTX(txEtherDeposit)
		runner.WaitForMinedCCTX(txERC20Deposit)
	}

	return runner
}

// createTestRoutine creates a test routine for given test names
func createTestRoutine(r *runner.E2ERunner, testNames []string) func() error {
	return func() (err error) {
		r.Logger.Print("🏃 starting bitcoin tests")
		startTime := time.Now()

		// run bitcoin tests
		testsToRun, err := r.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("bitcoin tests failed: %v", err)
		}

		if err := r.RunE2ETests(testsToRun); err != nil {
			return fmt.Errorf("bitcoin tests failed: %v", err)
		}

		r.Logger.Print("🍾 bitcoin tests completed in %s", time.Since(startTime).String())

		return err
	}
}
