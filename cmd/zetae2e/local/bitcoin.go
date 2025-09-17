package local

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/pkg/errgroup"
	"github.com/zeta-chain/node/testutil"
)

const (
	testGroupDepositName  = "btc_deposit"
	testGroupWithdrawName = "btc_withdraw"
)

// startBitcoinTests starts Bitcoin related tests
func startBitcoinTests(
	eg *errgroup.Group,
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	light, skipBitcoinSetup bool,
) {
	// start the bitcoin tests
	// btc withdraw tests are those that need a Bitcoin node wallet to send UTXOs
	bitcoinDepositTests := []string{
		//e2etests.TestBitcoinDonationName,
		e2etests.TestBitcoinDepositName,
		// e2etests.TestBitcoinDepositFastConfirmationName,
		// e2etests.TestBitcoinDepositAndCallName,
		// e2etests.TestBitcoinDepositAndCallRevertName,
		// e2etests.TestBitcoinStdMemoDepositName,
		// e2etests.TestBitcoinStdMemoDepositAndCallName,
		// e2etests.TestBitcoinStdMemoDepositAndCallRevertName,
		// e2etests.TestBitcoinStdMemoDepositAndCallRevertAndAbortName,
		e2etests.TestBitcoinStdMemoInscribedDepositAndCallName,
		// e2etests.TestBitcoinDepositAndAbortWithLowDepositFeeName,
		// e2etests.TestBitcoinDepositInvalidMemoRevertName,
		// e2etests.TestCrosschainSwapName,
	}
	bitcoinDepositTestsAdvanced := []string{
		// e2etests.TestBitcoinDepositAndCallRevertWithDustName,
		// e2etests.TestBitcoinStdMemoDepositAndCallRevertOtherAddressName,
		// e2etests.TestBitcoinDepositAndWithdrawWithDustName,
	}
	bitcoinWithdrawTests := []string{
		// need initial deposit to fund the withdraws
		// e2etests.TestBitcoinDepositName,
		// e2etests.TestBitcoinWithdrawSegWitName,
		// e2etests.TestBitcoinWithdrawInvalidAddressName,
		// e2etests.TestLegacyZetaWithdrawBTCRevertName,
	}
	bitcoinWithdrawTestsAdvanced := []string{
		// e2etests.TestBitcoinWithdrawTaprootName,
		// e2etests.TestBitcoinWithdrawLegacyName,
		// e2etests.TestBitcoinWithdrawP2SHName,
		// e2etests.TestBitcoinWithdrawP2WSHName,
		// e2etests.TestBitcoinWithdrawMultipleName,
		// e2etests.TestBitcoinWithdrawRestrictedName,
		// to run RBF test, change the constant 'minTxConfirmations' to 1 in the Bitcoin observer
		// https://github.com/zeta-chain/node/blob/5c2a8ffbc702130fd9538b1cd7640d0e04d3e4f6/zetaclient/chains/bitcoin/observer/outbound.go#L27
		//e2etests.TestBitcoinWithdrawRBFName,
	}

	if !light {
		// if light is enabled, only the most basic tests are run and advanced are skipped
		bitcoinDepositTests = append(bitcoinDepositTests, bitcoinDepositTestsAdvanced...)
		bitcoinWithdrawTests = append(bitcoinWithdrawTests, bitcoinWithdrawTestsAdvanced...)
	}
	bitcoinDepositTestRoutine, bitcoinWithdrawTestRoutine := bitcoinTestRoutines(
		conf,
		deployerRunner,
		verbose,
		!skipBitcoinSetup,
		bitcoinDepositTests,
		bitcoinWithdrawTests,
	)
	eg.Go(bitcoinDepositTestRoutine)
	eg.Go(bitcoinWithdrawTestRoutine)
}

// bitcoinTestRoutines returns test routines for deposit and withdraw tests
func bitcoinTestRoutines(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	initNetwork bool,
	depositTests []string,
	withdrawTests []string,
) (func() error, func() error) {
	// initialize runner for deposit tests
	account := conf.AdditionalAccounts.UserBitcoinDeposit
	runnerDeposit := initBitcoinRunner(
		testGroupDepositName,
		account,
		conf,
		deployerRunner,
		color.FgYellow,
		verbose,
		initNetwork,
	)

	// initialize runner for withdraw tests
	account = conf.AdditionalAccounts.UserBitcoinWithdraw
	runnerWithdraw := initBitcoinRunner(
		testGroupWithdrawName,
		account,
		conf,
		deployerRunner,
		color.FgHiYellow,
		verbose,
		initNetwork,
	)

	// initialize funds
	// send BTC to TSS for gas fees and to tester ZEVM address
	if initNetwork {
		// mine 101 blocks to ensure the BTC rewards are spendable
		// Note: the block rewards can be sent to any address in here
		_, err := deployerRunner.GenerateToAddressIfLocalBitcoin(101, deployerRunner.GetBtcAddress())
		require.NoError(runnerDeposit, err)

		// donate BTC to TSS and send BTC to ZEVM addresses
		deployerRunner.DonateBTC()
	}

	// create test routines
	routineDeposit := createBitcoinTestRoutine(runnerDeposit, depositTests, testGroupDepositName)
	routineWithdraw := createBitcoinTestRoutine(runnerWithdraw, withdrawTests, testGroupWithdrawName)

	return routineDeposit, routineWithdraw
}

// initBitcoinRunner initializes the Bitcoin runner for given test name and account
func initBitcoinRunner(
	name string,
	account config.Account,
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	printColor color.Attribute,
	verbose, initNetwork bool,
) *runner.E2ERunner {
	// initialize runner for bitcoin test
	runner, err := initTestRunner(
		name,
		conf,
		deployerRunner,
		account,
		runner.NewLogger(verbose, printColor, name),
		runner.WithZetaTxServer(deployerRunner.ZetaTxServer),
	)
	testutil.NoError(err)

	// initialize funds
	if initNetwork {
		address, _ := runner.GetBtcKeypair()

		// ensure address is imported
		err := runner.BtcRPCClient.ImportAddress(runner.Ctx, address.EncodeAddress())
		require.NoError(runner, err)

		// send some BTC block rewards
		_, err = runner.GenerateToAddressIfLocalBitcoin(101, address)
		require.NoError(runner, err)

		// send ERC20 token on EVM
		txERC20Send := deployerRunner.SendERC20OnEVM(account.EVMAddress(), 1000)
		runner.WaitForTxReceiptOnEVM(txERC20Send)

		// deposit ETH and ERC20 tokens on ZetaChain
		txEtherDeposit := runner.DepositEtherToDeployer()
		txERC20Deposit := runner.DepositERC20ToDeployer()

		runner.WaitForMinedCCTX(txEtherDeposit)
		runner.WaitForMinedCCTX(txERC20Deposit)
	}

	return runner
}

// createBitcoinTestRoutine creates a test routine for given test names
// The 'wgDependency' argument is used to wait for dependent routine to complete
func createBitcoinTestRoutine(r *runner.E2ERunner, testNames []string, name string) func() error {
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

		// mark deposit test group as done
		if name == testGroupDepositName {
			e2etests.DepdencyAllBitcoinDeposits.Done()
		}

		return err
	}
}
