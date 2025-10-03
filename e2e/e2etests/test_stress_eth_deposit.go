package e2etests

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/montanaflynn/stats"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

type depositResult struct {
	index     int
	txHash    ethcommon.Hash
	cctx      *crosschaintypes.CrossChainTx
	duration  time.Duration
	err       error
	startTime time.Time
}

// TestStressEtherDeposit tests the stressing deposit of ether
func TestStressEtherDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 4)

	// parse deposit amount and number of deposits
	depositAmount := utils.ParseBigInt(r, args[0])
	numDeposits := utils.ParseInt(r, args[1])
	batchSize := utils.ParseInt(r, args[2])
	batchInterval := utils.ParseInt(r, args[3])

	// configure test mode
	// true: tolerates transaction failures and reports statistics at the end
	bestEffortMode := true

	r.Logger.Print("starting stress test of %d deposits (mode: best-effort=%v)", numDeposits, bestEffortMode)

	// Get initial nonce
	initialNonce, err := r.EVMClient.PendingNonceAt(context.Background(), r.EVMAddress())
	require.NoError(r, err)
	r.Logger.Print("starting nonce: %d", initialNonce)

	// Channels for tracking
	depositTxs := make(chan depositResult, numDeposits)
	results := make(chan depositResult, numDeposits)

	// Statistics tracking
	var (
		sentCount    atomic.Uint64
		successCount atomic.Uint64
		failedCount  atomic.Uint64
	)

	// Start CCTX monitoring goroutines
	var wg sync.WaitGroup
	numMonitors := 10
	for i := 0; i < numMonitors; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			monitorDeposits(r, depositTxs, results, bestEffortMode, workerID)
		}(i)
	}

	// Send all deposits with batching
	sendStart := time.Now()
	batchCount := (numDeposits + batchSize - 1) / batchSize
	r.Logger.Print("sending %d deposits in %d batches (batch size: %d, interval: %dms)",
		numDeposits, batchCount, batchSize, batchInterval)

	var currentNonce = initialNonce

	for batchIdx := 0; batchIdx < batchCount; batchIdx++ {
		batchStart := batchIdx * batchSize
		batchEnd := minInt((batchIdx+1)*batchSize, numDeposits)
		batchSize := batchEnd - batchStart

		r.Logger.Print("sending batch %d/%d (%d deposits)", batchIdx+1, batchCount, batchSize)

		// Send all deposits in this batch
		for i := batchStart; i < batchEnd; i++ {
			// Create a new transactor with specific nonce
			auth := *r.EVMAuth
			auth.Nonce = big.NewInt(int64(currentNonce))
			auth.Value = depositAmount

			// Send deposit transaction
			tx, err := r.GatewayEVM.Deposit0(
				&auth,
				r.EVMAddress(),
				gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
			)

			if err != nil {
				r.Logger.Print("index %d (nonce %d): failed to send deposit: %v", i, currentNonce, err)
				if !bestEffortMode {
					require.FailNow(r, fmt.Sprintf("failed to send deposit %d: %v", i, err))
				}
				failedCount.Add(1)
				currentNonce++
				continue
			}

			// Success - don't wait for receipt, just send to monitor
			sentCount.Add(1)
			currentNonce++

			r.Logger.Print("index %d: deposit broadcast, tx hash: %s", i, tx.Hash().Hex())

			depositTxs <- depositResult{
				index:     i,
				txHash:    tx.Hash(),
				startTime: time.Now(),
			}

			// Small delay within batch
			time.Sleep(time.Millisecond * 20)
		}

		r.Logger.Print("batch %d/%d sent (%d deposits)", batchIdx+1, batchCount, batchSize)

		// Wait before sending next batch
		if batchIdx < batchCount-1 {
			r.Logger.Print("waiting %dms before next batch...", batchInterval)
			time.Sleep(time.Duration(batchInterval) * time.Millisecond)
		}
	}

	close(depositTxs)
	sendDuration := time.Since(sendStart)

	r.Logger.Print("all %d deposits sent in %v", sentCount.Load(), sendDuration)
	r.Logger.Print("waiting for CCTXs to complete...")

	// Wait for all monitors to finish
	wg.Wait()
	close(results)

	// Collect and analyze results
	monitorStart := time.Now()
	var failedDeposits []depositResult
	var durations []float64

	for result := range results {
		if result.err != nil {
			failedCount.Add(1)
			failedDeposits = append(failedDeposits, result)
			if !bestEffortMode {
				require.FailNow(r, fmt.Sprintf("deposit %d failed: %v", result.index, result.err))
			}
		} else if result.cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
			failedCount.Add(1)
			failedDeposits = append(failedDeposits, result)
			if !bestEffortMode {
				require.FailNow(r, fmt.Sprintf("deposit %d cctx failed with status %s",
					result.index, result.cctx.CctxStatus.Status))
			}
		} else {
			successCount.Add(1)
			durations = append(durations, result.duration.Seconds())
		}
	}
	monitorDuration := time.Since(monitorStart)

	// Calculate statistics
	var desc *stats.Description
	var descErr error
	if len(durations) > 0 {
		desc, descErr = stats.Describe(durations, false, &[]float64{50.0, 75.0, 90.0, 95.0, 99.0})
		if descErr != nil {
			r.Logger.Print("warning: failed to calculate latency statistics: %v", descErr)
		}
	}

	// Print final statistics
	r.Logger.Print("═══════════════════════════════════════")
	r.Logger.Print("Stress Test Results:")
	r.Logger.Print("═══════════════════════════════════════")
	r.Logger.Print("Configuration:")
	r.Logger.Print("  Batch size:          %d deposits", batchSize)
	r.Logger.Print("  Batch interval:      %dms", batchInterval)
	r.Logger.Print("  Total batches:       %d", batchCount)
	r.Logger.Print("Results:")
	r.Logger.Print("  Total deposits:      %d", numDeposits)
	r.Logger.Print("  Sent successfully:   %d", sentCount.Load())
	r.Logger.Print("  Succeeded:           %d (%.2f%%)", successCount.Load(),
		float64(successCount.Load())/float64(numDeposits)*100)
	r.Logger.Print("  Failed:              %d (%.2f%%)", failedCount.Load(),
		float64(failedCount.Load())/float64(numDeposits)*100)
	r.Logger.Print("Timing:")
	r.Logger.Print("  Send duration:       %v", sendDuration)
	r.Logger.Print("  Monitor duration:    %v", monitorDuration)
	r.Logger.Print("  Total duration:      %v", sendDuration+monitorDuration)
	r.Logger.Print("  TPS (send):          %.2f", float64(sentCount.Load())/sendDuration.Seconds())

	if len(durations) > 0 && descErr == nil {
		r.Logger.Print("Latency Statistics (seconds):")
		r.Logger.Print("  Min:                 %.2f", desc.Min)
		r.Logger.Print("  Max:                 %.2f", desc.Max)
		r.Logger.Print("  Mean:                %.2f", desc.Mean)
		r.Logger.Print("  Std Dev:             %.2f", desc.Std)
		for _, p := range desc.DescriptionPercentiles {
			r.Logger.Print("  P%.0f:                 %.2f", p.Percentile, p.Value)
		}
	}
	r.Logger.Print("═══════════════════════════════════════")

	if len(failedDeposits) > 0 && len(failedDeposits) <= 10 {
		r.Logger.Print("Failed deposit details:")
		for _, failed := range failedDeposits {
			if failed.cctx != nil {
				r.Logger.Print("  - Index %d, Hash %s, Status: %s, Message: %s",
					failed.index, failed.txHash.Hex(),
					failed.cctx.CctxStatus.Status, failed.cctx.CctxStatus.StatusMessage)
			} else {
				r.Logger.Print("  - Index %d, Hash %s: %v", failed.index, failed.txHash.Hex(), failed.err)
			}
		}
	} else if len(failedDeposits) > 10 {
		r.Logger.Print("First 10 failed deposits:")
		for i := 0; i < 10; i++ {
			failed := failedDeposits[i]
			if failed.cctx != nil {
				r.Logger.Print("  - Index %d, Hash %s, Status: %s",
					failed.index, failed.txHash.Hex(), failed.cctx.CctxStatus.Status)
			} else {
				r.Logger.Print("  - Index %d, Hash %s: %v", failed.index, failed.txHash.Hex(), failed.err)
			}
		}
		r.Logger.Print("  ... and %d more", len(failedDeposits)-10)
	}

	if !bestEffortMode && failedCount.Load() > 0 {
		require.FailNow(r, fmt.Sprintf("%d deposits failed", failedCount.Load()))
	}

	r.Logger.Print("stress test completed")
}

// monitorDeposits monitors deposit CCTXs in a goroutine
func monitorDeposits(
	r *runner.E2ERunner,
	depositTxs <-chan depositResult,
	results chan<- depositResult,
	bestEffortMode bool,
	workerID int,
) {
	for deposit := range depositTxs {
		result := deposit

		// Wait for CCTX to be mined
		cctx := utils.WaitCctxMinedByInboundHash(
			r.Ctx,
			deposit.txHash.Hex(),
			r.CctxClient,
			r.Logger,
			r.CctxTimeout,
		)

		result.cctx = cctx
		result.duration = time.Since(deposit.startTime)

		if cctx.CctxStatus.Status == crosschaintypes.CctxStatus_OutboundMined {
			if deposit.index%10 == 0 || !bestEffortMode {
				r.Logger.Print("worker %d: index %d: ✓ completed in %v",
					workerID, deposit.index, result.duration)
			}
		} else {
			result.err = fmt.Errorf("cctx failed with status %s: %s",
				cctx.CctxStatus.Status, cctx.CctxStatus.StatusMessage)
			r.Logger.Print("worker %d: index %d: ✗ failed with status %s",
				workerID, deposit.index, cctx.CctxStatus.Status)
		}

		results <- result
	}
}
