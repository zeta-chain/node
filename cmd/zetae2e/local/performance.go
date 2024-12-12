package local

// performance.go provides routines that run the stress tests for different actions (deposit, withdraw) to measure network performance
// Note: the routine provided here should not be used concurrently with other routines as these reuse the accounts of other routines

import (
	"fmt"
	"math/big"
	"time"

	"github.com/fatih/color"

	"github.com/zeta-chain/node/e2e/config"
	"github.com/zeta-chain/node/e2e/e2etests"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// ethereumDepositPerformanceRoutine runs performance tests for Ether deposit
func ethereumDepositPerformanceRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		// initialize runner for ether test
		r, err := initTestRunner(
			"ether",
			conf,
			deployerRunner,
			conf.AdditionalAccounts.UserLegacyERC20,
			runner.NewLogger(verbose, color.FgHiMagenta, "perf_eth_deposit"),
		)
		if err != nil {
			return err
		}

		r.Logger.Print("🏃 starting Ethereum deposit performance tests")
		startTime := time.Now()

		tests, err := r.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("ethereum deposit performance test failed: %v", err)
		}

		if err := r.RunE2ETests(tests); err != nil {
			return fmt.Errorf("ethereum deposit performance test failed: %v", err)
		}

		r.Logger.Print("🍾 Ethereum deposit performance test completed in %s", time.Since(startTime).String())

		return err
	}
}

// ethereumWithdrawPerformanceRoutine runs performance tests for Ether withdraw
func ethereumWithdrawPerformanceRoutine(
	conf config.Config,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	testNames ...string,
) func() error {
	return func() (err error) {
		// initialize runner for ether test
		r, err := initTestRunner(
			"ether",
			conf,
			deployerRunner,
			conf.AdditionalAccounts.UserLegacyEther,
			runner.NewLogger(verbose, color.FgHiBlue, "perf_eth_withdraw"),
		)
		if err != nil {
			return err
		}

		if r.ReceiptTimeout == 0 {
			r.ReceiptTimeout = 15 * time.Minute
		}
		if r.CctxTimeout == 0 {
			r.CctxTimeout = 15 * time.Minute
		}

		r.Logger.Print("🏃 starting Ethereum withdraw performance tests")
		startTime := time.Now()

		// depositing the necessary tokens on ZetaChain
		txEtherDeposit := r.LegacyDepositEther()
		r.WaitForMinedCCTX(txEtherDeposit)

		tests, err := r.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("ethereum withdraw performance test failed: %v", err)
		}

		if err := r.RunE2ETests(tests); err != nil {
			return fmt.Errorf("ethereum withdraw performance test failed: %v", err)
		}

		r.Logger.Print("🍾 Ethereum withdraw performance test completed in %s", time.Since(startTime).String())

		return err
	}
}

// solanaDepositPerformanceRoutine runs performance tests for solana deposits
func solanaDepositPerformanceRoutine(
	conf config.Config,
	name string,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	account config.Account,
	testNames ...string,
) func() error {
	return func() (err error) {
		// initialize runner for solana test
		r, err := initTestRunner(
			"solana",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgHiMagenta, name),
			runner.WithZetaTxServer(deployerRunner.ZetaTxServer),
		)
		if err != nil {
			return err
		}

		if r.ReceiptTimeout == 0 {
			r.ReceiptTimeout = 15 * time.Minute
		}
		if r.CctxTimeout == 0 {
			r.CctxTimeout = 30 * time.Minute
		}

		r.Logger.Print("🏃 starting solana deposit performance tests")
		startTime := time.Now()

		tests, err := r.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("solana deposit performance test failed: %v", err)
		}

		if err := r.RunE2ETests(tests); err != nil {
			return fmt.Errorf("solana deposit performance test failed: %v", err)
		}

		r.Logger.Print("🍾 solana deposit performance test completed in %s", time.Since(startTime).String())

		return err
	}
}

// solanaWithdrawPerformanceRoutine runs performance tests for solana withdrawals
func solanaWithdrawPerformanceRoutine(
	conf config.Config,
	name string,
	deployerRunner *runner.E2ERunner,
	verbose bool,
	account config.Account,
	testNames ...string,
) func() error {
	return func() (err error) {
		// initialize runner for solana test
		r, err := initTestRunner(
			"solana",
			conf,
			deployerRunner,
			account,
			runner.NewLogger(verbose, color.FgHiGreen, name),
			runner.WithZetaTxServer(deployerRunner.ZetaTxServer),
		)
		if err != nil {
			return err
		}

		if r.ReceiptTimeout == 0 {
			r.ReceiptTimeout = 15 * time.Minute
		}
		if r.CctxTimeout == 0 {
			r.CctxTimeout = 30 * time.Minute
		}

		r.Logger.Print("🏃 starting solana withdraw performance tests")
		startTime := time.Now()

		// load deployer private key
		privKey := r.GetSolanaPrivKey()

		// execute the deposit sol transaction
		amount := big.NewInt(0).Mul(big.NewInt(1e9), big.NewInt(100)) // 100 sol in lamports
		sig := r.SOLDepositAndCall(nil, r.EVMAddress(), amount, nil)

		// wait for the cctx to be mined
		cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient, r.Logger, r.CctxTimeout)
		r.Logger.CCTX(*cctx, "solana_deposit")
		utils.RequireCCTXStatus(r, cctx, types.CctxStatus_OutboundMined)

		// same amount for spl
		sig = r.SPLDepositAndCall(&privKey, amount.Uint64(), r.SPLAddr, r.EVMAddress(), nil)

		// wait for the cctx to be mined
		cctx = utils.WaitCctxMinedByInboundHash(r.Ctx, sig.String(), r.CctxClient, r.Logger, r.CctxTimeout)
		r.Logger.CCTX(*cctx, "solana_deposit_spl")
		utils.RequireCCTXStatus(r, cctx, types.CctxStatus_OutboundMined)

		tests, err := r.GetE2ETestsToRunByName(
			e2etests.AllE2ETests,
			testNames...,
		)
		if err != nil {
			return fmt.Errorf("solana withdraw performance test failed: %v", err)
		}

		if err := r.RunE2ETests(tests); err != nil {
			return fmt.Errorf("solana withdraw performance test failed: %v", err)
		}

		r.Logger.Print("🍾 solana withdraw performance test completed in %s", time.Since(startTime).String())

		return err
	}
}
