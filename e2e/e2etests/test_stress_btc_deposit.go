package e2etests

import (
	"fmt"
	"time"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestStressBTCDeposit tests the stressing deposit of BTC
func TestStressBTCDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	depositAmount := parseFloat(r, args[0])
	numDeposits := parseInt(r, args[1])

	r.SetBtcAddress(r.Name, false)

	r.Logger.Print("starting stress test of %d deposits", numDeposits)

	// create a wait group to wait for all the deposits to complete
	var eg errgroup.Group

	// send the deposits
	for i := 0; i < numDeposits; i++ {
		i := i
		txHash := r.DepositBTCWithAmount(depositAmount)
		r.Logger.Print("index %d: starting deposit, tx hash: %s", i, txHash.String())

		eg.Go(func() error { return monitorBTCDeposit(r, txHash, i, time.Now()) })
	}

	require.NoError(r, eg.Wait())

	r.Logger.Print("all deposits completed")
}

// monitorBTCDeposit monitors the deposit of BTC, returns once the deposit is complete
func monitorBTCDeposit(r *runner.E2ERunner, hash *chainhash.Hash, index int, startTime time.Time) error {
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, hash.String(), r.CctxClient, r.Logger, r.ReceiptTimeout)
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
