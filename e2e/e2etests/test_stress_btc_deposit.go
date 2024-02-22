package e2etests

import (
	"fmt"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"golang.org/x/sync/errgroup"
)

// TestStressBTCDeposit tests the stressing deposit of BTC
func TestStressBTCDeposit(r *runner.E2ERunner) {
	// number of deposits to perform
	numDeposits := 100

	r.SetBtcAddress(r.Name, false)

	r.Logger.Print("starting stress test of %d deposits", numDeposits)

	// create a wait group to wait for all the deposits to complete
	var eg errgroup.Group

	// send the deposits
	for i := 0; i < numDeposits; i++ {
		i := i
		txHash := r.DepositBTCWithAmount(0.001)
		r.Logger.Print("index %d: starting deposit, tx hash: %s", i, txHash.String())

		eg.Go(func() error {
			return MonitorBTCDeposit(r, txHash, i, time.Now())
		})
	}

	// wait for all the deposits to complete
	if err := eg.Wait(); err != nil {
		panic(err)
	}

	r.Logger.Print("all deposits completed")
}

// MonitorBTCDeposit monitors the deposit of BTC, returns once the deposit is complete
func MonitorBTCDeposit(r *runner.E2ERunner, hash *chainhash.Hash, index int, startTime time.Time) error {
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, hash.String(), r.CctxClient, r.Logger, r.ReceiptTimeout)
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		return fmt.Errorf(
			"index %d: deposit cctx failed with status %s, message %s, cctx index %s",
			index,
			cctx.CctxStatus.Status,
			cctx.CctxStatus.StatusMessage,
			cctx.Index,
		)
	}
	timeToComplete := time.Now().Sub(startTime)
	r.Logger.Print("index %d: deposit cctx success in %s", index, timeToComplete.String())

	return nil
}
