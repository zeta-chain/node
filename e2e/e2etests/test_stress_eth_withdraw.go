package e2etests

import (
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestStressEtherWithdraw tests the stressing withdraw of ether
func TestStressEtherWithdraw(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	// parse withdraw amount and number of withdraws
	withdrawalAmount := utils.ParseBigInt(r, args[0])

	numWithdraws, err := strconv.Atoi(args[1])
	require.NoError(r, err)
	require.GreaterOrEqual(r, numWithdraws, 1)

	tx, err := r.ETHZRC20.Approve(r.ZEVMAuth, r.ETHZRC20Addr, big.NewInt(1e18))
	require.NoError(r, err)

	r.WaitForTxReceiptOnZEVM(tx)

	r.Logger.Print("starting stress test of %d withdraws", numWithdraws)

	// create a wait group to wait for all the withdraws to complete
	var eg errgroup.Group

	// store durations as float64 seconds like prometheus
	withdrawDurations := []float64{}
	withdrawDurationsLock := sync.Mutex{}

	// send the withdraws
	for i := 0; i < numWithdraws; i++ {
		i := i

		tx, err := r.ETHZRC20.Withdraw(r.ZEVMAuth, r.EVMAddress().Bytes(), withdrawalAmount)
		require.NoError(r, err)

		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		utils.RequireTxSuccessful(r, receipt)

		r.Logger.Print("index %d: starting withdraw, tx hash: %s", i, tx.Hash().Hex())

		eg.Go(func() error {
			startTime := time.Now()
			cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.ReceiptTimeout)
			if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
				return fmt.Errorf(
					"index %d: withdraw cctx failed with status %s, message %s, cctx index %s",
					i,
					cctx.CctxStatus.Status,
					cctx.CctxStatus.StatusMessage,
					cctx.Index,
				)
			}
			timeToComplete := time.Since(startTime)
			r.Logger.Print("index %d: withdraw cctx success in %s", i, timeToComplete.String())

			withdrawDurationsLock.Lock()
			withdrawDurations = append(withdrawDurations, timeToComplete.Seconds())
			withdrawDurationsLock.Unlock()

			return nil
		})
	}

	err = eg.Wait()

	desc, descErr := stats.Describe(withdrawDurations, false, &[]float64{50.0, 75.0, 90.0, 95.0})
	if descErr != nil {
		r.Logger.Print("❌ failed to calculate latency report: %v", descErr)
	}

	r.Logger.Print("Latency report:")
	r.Logger.Print("min:  %.2f", desc.Min)
	r.Logger.Print("max:  %.2f", desc.Max)
	r.Logger.Print("mean: %.2f", desc.Mean)
	r.Logger.Print("std:  %.2f", desc.Std)
	for _, p := range desc.DescriptionPercentiles {
		r.Logger.Print("p%.0f:  %.2f", p.Percentile, p.Value)
	}

	require.NoError(r, err)

	r.Logger.Print("all withdraws completed")
}
