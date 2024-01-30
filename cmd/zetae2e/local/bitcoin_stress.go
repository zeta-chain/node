package local

import (
	"fmt"
	"math/big"
	"runtime"
	"time"

	"github.com/btcsuite/btcutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"

	"github.com/fatih/color"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/config"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/runner"
)

// bitcoinTestRoutine runs Bitcoin related smoke tests
func bitcoinStressTestRoutine(
	conf config.Config,
	deployerRunner *runner.SmokeTestRunner,
	verbose bool,
	initBitcoinNetwork bool,
) func() error {
	return func() (err error) {
		// return an error on panic
		// TODO: remove and instead return errors in the tests
		// https://github.com/zeta-chain/node/issues/1500
		defer func() {
			if r := recover(); r != nil {
				// print stack trace
				stack := make([]byte, 4096)
				n := runtime.Stack(stack, false)
				err = fmt.Errorf("bitcoin panic: %v, stack trace %s", r, stack[:n])
			}
		}()

		// initialize runner for bitcoin test
		bitcoinRunner, err := initTestRunner(
			"bitcoin",
			conf,
			deployerRunner,
			UserBitcoinAddress,
			UserBitcoinPrivateKey,
			runner.NewLogger(verbose, color.FgYellow, "bitcoin"),
		)
		if err != nil {
			return err
		}

		bitcoinRunner.Logger.Print("🏃 starting Bitcoin tests")
		//startTime :=  time.Now()

		// funding the account
		txUSDTSend := deployerRunner.SendUSDTOnEvm(UserBitcoinAddress, 1000)
		bitcoinRunner.WaitForTxReceiptOnEvm(txUSDTSend)

		// depositing the necessary tokens on ZetaChain
		txEtherDeposit := bitcoinRunner.DepositEther(false)
		txERC20Deposit := bitcoinRunner.DepositERC20()

		bitcoinRunner.WaitForMinedCCTX(txEtherDeposit)
		bitcoinRunner.WaitForMinedCCTX(txERC20Deposit)

		bitcoinRunner.SetupBitcoinAccount(initBitcoinNetwork)
		bitcoinRunner.DepositBTC(true)

		tx, err := bitcoinRunner.BTCZRC20.Approve(bitcoinRunner.ZevmAuth, bitcoinRunner.BTCZRC20Addr, big.NewInt(1e18))
		if err != nil {
			panic(err)
		}
		receipt := utils.MustWaitForTxReceipt(bitcoinRunner.Ctx, bitcoinRunner.ZevmClient, tx, bitcoinRunner.Logger, bitcoinRunner.ReceiptTimeout)
		if receipt.Status != 1 {
			panic(fmt.Errorf("approve receipt status is not 1"))
		}

		bitcoinRunner.MineBlocks()

		go WithdrawCCtx(bitcoinRunner)

		time.Sleep(360 * time.Second)

		return nil
	}
}

var (
	zevmNonce = big.NewInt(1)
)

// WithdrawCCtx withdraw USDT from ZEVM to EVM
func WithdrawCCtx(sm *runner.SmokeTestRunner) {
	ticker := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-ticker.C:
			WithdrawBTCZRC20(sm)
		}
	}
}

func WithdrawBTCZRC20(sm *runner.SmokeTestRunner) {
	defer func() {
		zevmNonce.Add(zevmNonce, big.NewInt(1))
	}()

	sm.ZevmAuth.Nonce = zevmNonce

	sm.Logger.Print("nonce %d: starting withdraw", zevmNonce)
	tx, err := sm.BTCZRC20.Withdraw(sm.ZevmAuth, []byte(sm.BTCDeployerAddress.EncodeAddress()), big.NewInt(0.01*btcutil.SatoshiPerBitcoin))
	if err != nil {
		panic(err)
	}

	go MonitorCCTXFromTxHash(sm, tx, zevmNonce.Int64())
}

func MonitorCCTXFromTxHash(sm *runner.SmokeTestRunner, tx *ethtypes.Transaction, nonce int64) {
	receipt := utils.MustWaitForTxReceipt(sm.Ctx, sm.ZevmClient, tx, sm.Logger, sm.ReceiptTimeout)
	if receipt.Status == 0 {
		sm.Logger.Print("nonce %d: withdraw evm tx failed", nonce)
		return
	}
	// mine 10 blocks to confirm the withdraw tx
	_, err := sm.BtcRPCClient.GenerateToAddress(10, sm.BTCDeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	cctx := utils.WaitCctxMinedByInTxHash(sm.Ctx, tx.Hash().Hex(), sm.CctxClient, sm.Logger, sm.ReceiptTimeout)
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		sm.Logger.Print(
			"nonce %d: withdraw cctx failed with status %s, message %s",
			nonce,
			cctx.CctxStatus.Status,
			cctx.CctxStatus.StatusMessage,
		)
		return
	}
	sm.Logger.Print("nonce %d: withdraw cctx success", nonce)
}
