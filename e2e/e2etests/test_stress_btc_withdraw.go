package e2etests

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/btcsuite/btcutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"golang.org/x/sync/errgroup"
)

// TestStressBTCWithdraw tests the stressing withdraw of btc
func TestStressBTCWithdraw(r *runner.E2ERunner, args []string) {
	if len(args) != 2 {
		panic("TestStressBTCWithdraw requires exactly two arguments: the withdrawal amount and the number of withdrawals.")
	}

	withdrawalAmount, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		panic("Invalid withdrawal amount specified for TestStressBTCWithdraw.")
	}

	numWithdraws, err := strconv.Atoi(args[1])
	if err != nil || numWithdraws < 1 {
		panic("Invalid number of withdrawals specified for TestStressBTCWithdraw.")
	}

	r.SetBtcAddress(r.Name, false)

	r.Logger.Print("starting stress test of %d withdraws", numWithdraws)

	// create a wait group to wait for all the withdraws to complete
	var eg errgroup.Group

	satAmount, err := btcutil.NewAmount(withdrawalAmount)
	if err != nil {
		panic(err)
	}

	// send the withdraws
	for i := 0; i < numWithdraws; i++ {
		i := i
		tx, err := r.BTCZRC20.Withdraw(
			r.ZEVMAuth,
			[]byte(r.BTCDeployerAddress.EncodeAddress()),
			big.NewInt(int64(satAmount)),
		)
		if err != nil {
			panic(err)
		}
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		if receipt.Status == 0 {
			//r.Logger.Info("index %d: withdraw evm tx failed", index)
			panic(fmt.Sprintf("index %d: withdraw btc tx %s failed", i, tx.Hash().Hex()))
		}
		r.Logger.Print("index %d: starting withdraw, tx hash: %s", i, tx.Hash().Hex())

		eg.Go(func() error {
			return monitorBTCWithdraw(r, tx, i, time.Now())
		})
	}

	// wait for all the withdraws to complete
	if err := eg.Wait(); err != nil {
		panic(err)
	}

	r.Logger.Print("all withdraws completed")
}

// monitorBTCWithdraw monitors the withdraw of BTC, returns once the withdraw is complete
func monitorBTCWithdraw(r *runner.E2ERunner, tx *ethtypes.Transaction, index int, startTime time.Time) error {
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.ReceiptTimeout)
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		return fmt.Errorf(
			"index %d: withdraw cctx failed with status %s, message %s, cctx index %s",
			index,
			cctx.CctxStatus.Status,
			cctx.CctxStatus.StatusMessage,
			cctx.Index,
		)
	}
	timeToComplete := time.Now().Sub(startTime)
	r.Logger.Print("index %d: withdraw cctx success in %s", index, timeToComplete.String())

	return nil
}
