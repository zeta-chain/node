package e2etests

import (
	"fmt"
	"math/big"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestStressBTCWithdraw tests the stressing withdraw of btc
func TestStressBTCWithdraw(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	withdrawalAmount := utils.ParseFloat(r, args[0])
	numWithdraws := utils.ParseInt(r, args[1])

	r.Logger.Print("starting stress test of %d withdraws", numWithdraws)

	// create a wait group to wait for all the withdraws to complete
	var eg errgroup.Group

	satAmount, err := btcutil.NewAmount(withdrawalAmount)
	require.NoError(r, err)

	// send the withdraws
	for i := 0; i < numWithdraws; i++ {
		i := i
		tx, err := r.BTCZRC20.Withdraw(
			r.ZEVMAuth,
			[]byte(r.GetBtcAddress().EncodeAddress()),
			big.NewInt(int64(satAmount)),
		)
		require.NoError(r, err)

		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		utils.RequireTxSuccessful(r, receipt)
		r.Logger.Print("index %d: starting withdraw, tx hash: %s", i, tx.Hash().Hex())

		eg.Go(func() error { return monitorBTCWithdraw(r, tx, i, time.Now()) })
	}

	require.NoError(r, eg.Wait())

	r.Logger.Print("all withdraws completed")
}

// monitorBTCWithdraw monitors the withdraw of BTC, returns once the withdraw is complete
func monitorBTCWithdraw(r *runner.E2ERunner, tx *ethtypes.Transaction, index int, startTime time.Time) error {
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.ReceiptTimeout)
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		return fmt.Errorf(
			"index %d: withdraw cctx failed with status %s, message %s, cctx index %s",
			index,
			cctx.CctxStatus.Status,
			cctx.CctxStatus.StatusMessage,
			cctx.Index,
		)
	}
	timeToComplete := time.Since(startTime)
	r.Logger.Print("index %d: withdraw cctx success in %s", index, timeToComplete.String())

	return nil
}
