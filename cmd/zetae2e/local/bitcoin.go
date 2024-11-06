package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/testutil"
)

// startBitcoinTestRoutines starts Bitcoin deposit and withdraw tests in parallel
func startBitcoinTestRoutines(
	eg *errgroup.Group,
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	initNetwork bool,
	depositTests []string,
	withdrawTests []string,
) {
	// initialize runner for deposit tests
	runnerDeposit := initRunnerDeposit(conf, deployerRunner, verbose, initNetwork)

	// initialize runner for withdraw tests
	runnerWithdraw := initRunnerWithdraw(conf, deployerRunner, verbose, initNetwork)

	// initialize funds
	// send BTC to TSS for gas fees and to tester ZEVM address
	if initNetwork {
		// mine 101 blocks to ensure the BTC rewards are spendable
		// Note: the rewards can be sent to any address in here
		_, err := runnerDeposit.GenerateToAddressIfLocalBitcoin(101, runnerDeposit.BTCDeployerAddress)
		require.NoError(runnerDeposit, err)

		// send BTC to ZEVM addresses
		runnerDeposit.DepositBTC(runnerDeposit.EVMAddress())
		runnerDeposit.DepositBTC(runnerWithdraw.EVMAddress())
	}

	// create test routines
	routineDeposit := createTestRoutine(runnerDeposit, depositTests)
	routineWithdraw := createTestRoutine(runnerWithdraw, withdrawTests)

	// start test routines
	eg.Go(routineDeposit)
	eg.Go(routineWithdraw)
}

// initRunnerDeposit initializes the runner for deposit tests
func initRunnerDeposit(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose, initNetwork bool,
) *runner.E2ERunner {
	var (
		name         = "btc_deposit"
		account      = conf.AdditionalAccounts.UserBitcoin1
		createWallet = true // deposit tests need Bitcoin node wallet
	)

	return initRunner(name, account, conf, deployerRunner, verbose, initNetwork, createWallet)
}

// initRunnerWithdraw initializes the runner for withdraw tests
func initRunnerWithdraw(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose, initNetwork bool,
) *runner.E2ERunner {
	var (
		name         = "btc_withdraw"
		account      = conf.AdditionalAccounts.UserBitcoin2
		createWallet = false // withdraw tests do not use Bitcoin node wallet
	)

	return initRunner(name, account, conf, deployerRunner, verbose, initNetwork, createWallet)
}

// initRunner initializes the runner for given test name and account
func initRunner(
	name string,
	account config.Account,
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose, initNetwork, createWallet bool,
) *runner.E2ERunner {
	// initialize r for bitcoin test
	r, err := initTestRunner(
		name,
		conf,
		deployerRunner,
		account,
		runner.NewLogger(verbose, color.FgYellow, name),
	)
	testutil.NoError(err)

	// setup TSS address and setup deployer wallet
	r.SetupBitcoinAccounts(createWallet)

	// initialize funds
	if initNetwork {
		// send some BTC block rewards to the deployer address
		_, err = r.GenerateToAddressIfLocalBitcoin(4, r.BTCDeployerAddress)
		require.NoError(r, err)

		// send ERC20 token on EVM
		txERC20Send := deployerRunner.SendERC20OnEvm(account.EVMAddress(), 1000)
		r.WaitForTxReceiptOnEvm(txERC20Send)

		// deposit ETH and ERC20 tokens on ZetaChain
		txEtherDeposit := r.DepositEther()
		txERC20Deposit := r.DepositERC20()

		r.WaitForMinedCCTX(txEtherDeposit)
		r.WaitForMinedCCTX(txERC20Deposit)
	}

	return r
}

// createTestRoutine creates a test routine for given test names
// Note: due to the extensive block generation in Bitcoin localnet, block header test is run first
// to make it faster to catch up with the latest block header
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

		// if err := bitcoinRunner.CheckBtcTSSBalance(); err != nil {
		// 	return err
		// }

		r.Logger.Print("🍾 bitcoin tests completed in %s", time.Since(startTime).String())

		return err
	}
}
