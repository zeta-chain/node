package e2etests

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

type withdrawResult struct {
	index       int
	txHash      common.Hash
	cctx        *crosschaintypes.CrossChainTx
	duration    time.Duration
	err         error
	startTime   time.Time
	submittedAt time.Time
}

// TestStressEtherWithdraw tests the stressing withdraw of ether
func TestStressEtherWithdraw(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 4)

	// parse withdraw amount and number of withdraws
	withdrawalAmount := utils.ParseBigInt(r, args[0])

	numWithdraws, err := strconv.Atoi(args[1])
	require.NoError(r, err)
	require.GreaterOrEqual(r, numWithdraws, 1)

	batchSize := utils.ParseInt(r, args[2])
	batchInterval := utils.ParseInt(r, args[3])

	// configure test mode
	bestEffortMode := true

	r.Logger.Print("starting stress test of %d withdrawals (mode: best-effort=%v)", numWithdraws, bestEffortMode)

	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	// Get initial nonce
	initialNonce, err := r.ZEVMClient.PendingNonceAt(context.Background(), r.ZEVMAuth.From)
	require.NoError(r, err)
	r.Logger.Print("starting nonce: %d", initialNonce)

	// Channels for tracking
	withdrawTxs := make(chan withdrawResult, numWithdraws)
	results := make(chan withdrawResult, numWithdraws)

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
			monitorWithdrawals(r, withdrawTxs, results, bestEffortMode, workerID)
		}(i)
	}

	// Send all withdrawals with batching
	sendStart := time.Now()
	batchCount := (numWithdraws + batchSize - 1) / batchSize
	r.Logger.Print("sending %d withdrawals in %d batches (batch size: %d, interval: %dms)",
		numWithdraws, batchCount, batchSize, batchInterval)

	currentNonce := initialNonce

	for batchIdx := 0; batchIdx < batchCount; batchIdx++ {
		batchStart := batchIdx * batchSize
		batchEnd := minInt((batchIdx+1)*batchSize, numWithdraws)
		actualBatchSize := batchEnd - batchStart

		r.Logger.Print("sending batch %d/%d (%d withdrawals)", batchIdx+1, batchCount, actualBatchSize)

		// Send all withdrawals in this batch
		for i := batchStart; i < batchEnd; i++ {
			txSent := false
			retryCount := 0

			// Retry loop - infinite retries for mempool full, otherwise fail
			for !txSent {
				// Create a new transactor with specific nonce
				auth := *r.ZEVMAuth
				// #nosec G115 e2eTest - always in range
				auth.Nonce = big.NewInt(int64(currentNonce))

				// Send withdrawal transaction
				r.Logger.Print("gatewayzevm %s", r.GatewayZEVMAddr.Hex())
				tx, err := r.GatewayZEVM.Withdraw0(
					&auth,
					r.EVMAddress().Bytes(),
					withdrawalAmount,
					r.ETHZRC20Addr,
					gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
				)

				if err != nil {
					// Check if mempool is full
					if isErrMempoolFull(err) {
						retryCount++
						r.Logger.Print("index %d (nonce %d): mempool is full (retry %d), waiting 5 seconds...",
							i, currentNonce, retryCount)
						time.Sleep(5 * time.Second)
						continue // retry with same nonce
					}

					// Other errors - fail or skip
					r.Logger.Print("index %d (nonce %d): failed to send withdrawal: %v", i, currentNonce, err)
					if !bestEffortMode {
						require.FailNow(r, fmt.Sprintf("failed to send withdrawal %d: %v", i, err))
					}
					failedCount.Add(1)
					break // exit retry loop and move to next withdrawal
				}

				sentCount.Add(1)
				currentNonce++
				txSent = true

				r.Logger.Print("index %d: withdrawal broadcast, tx hash: %s", i, tx.Hash().Hex())

				withdrawTxs <- withdrawResult{
					index:       i,
					txHash:      tx.Hash(),
					startTime:   time.Now(),
					submittedAt: time.Now(),
				}

				// Small delay within batch
				time.Sleep(time.Millisecond * 20)
			}
		}

		r.Logger.Print("batch %d/%d sent (%d withdrawals)", batchIdx+1, batchCount, actualBatchSize)

		// Wait before sending next batch
		if batchIdx < batchCount-1 {
			r.Logger.Print("waiting %dms before next batch...", batchInterval)
			time.Sleep(time.Duration(batchInterval) * time.Millisecond)
		}
	}

	close(withdrawTxs)
	sendDuration := time.Since(sendStart)

	r.Logger.Print("all %d withdrawals sent in %v", sentCount.Load(), sendDuration)
	r.Logger.Print("waiting for CCTXs to complete...")

	// Wait for all monitors to finish
	wg.Wait()
	close(results)

	// Collect and analyze results
	var failedWithdraws []withdrawResult

	for result := range results {
		if result.err != nil {
			failedCount.Add(1)
			failedWithdraws = append(failedWithdraws, result)
			if !bestEffortMode {
				require.FailNow(r, fmt.Sprintf("withdrawal %d failed: %v", result.index, result.err))
			}
		} else if result.cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
			failedCount.Add(1)
			failedWithdraws = append(failedWithdraws, result)
			if !bestEffortMode {
				require.FailNow(r, fmt.Sprintf("withdrawal %d cctx failed with status %s",
					result.index, result.cctx.CctxStatus.Status))
			}
		} else {
			successCount.Add(1)
		}
	}

	// Print final statistics
	r.Logger.Print("═══════════════════════════════════════")
	r.Logger.Print("Stress Test Results:")
	r.Logger.Print("═══════════════════════════════════════")
	r.Logger.Print("Configuration:")
	r.Logger.Print("  Batch size:          %d withdrawals", batchSize)
	r.Logger.Print("  Batch interval:      %dms", batchInterval)
	r.Logger.Print("  Total batches:       %d", batchCount)
	r.Logger.Print("Results:")
	r.Logger.Print("  Total withdrawals:   %d", numWithdraws)
	r.Logger.Print("  Sent successfully:   %d", sentCount.Load())
	r.Logger.Print("  Succeeded:           %d (%.2f%%)", successCount.Load(),
		float64(successCount.Load())/float64(numWithdraws)*100)
	r.Logger.Print("  Failed:              %d (%.2f%%)", failedCount.Load(),
		float64(failedCount.Load())/float64(numWithdraws)*100)
	r.Logger.Print("═══════════════════════════════════════")

	if len(failedWithdraws) > 0 && len(failedWithdraws) <= 10 {
		r.Logger.Print("Failed withdrawal details:")
		for _, failed := range failedWithdraws {
			if failed.cctx != nil {
				r.Logger.Print("  - Index %d, Hash %s, Status: %s, Message: %s",
					failed.index, failed.txHash.Hex(),
					failed.cctx.CctxStatus.Status, failed.cctx.CctxStatus.StatusMessage)
			} else {
				r.Logger.Print("  - Index %d, Hash %s: %v", failed.index, failed.txHash.Hex(), failed.err)
			}
		}
	} else if len(failedWithdraws) > 10 {
		r.Logger.Print("First 10 failed withdrawals:")
		for i := 0; i < 10; i++ {
			failed := failedWithdraws[i]
			if failed.cctx != nil {
				r.Logger.Print("  - Index %d, Hash %s, Status: %s",
					failed.index, failed.txHash.Hex(), failed.cctx.CctxStatus.Status)
			} else {
				r.Logger.Print("  - Index %d, Hash %s: %v", failed.index, failed.txHash.Hex(), failed.err)
			}
		}
		r.Logger.Print("  ... and %d more", len(failedWithdraws)-10)
	}

	if !bestEffortMode && failedCount.Load() > 0 {
		require.FailNow(r, fmt.Sprintf("%d withdrawals failed", failedCount.Load()))
	}

	r.Logger.Print("stress test completed")
}

// monitorWithdrawals monitors withdrawal CCTXs in a goroutine
func monitorWithdrawals(
	r *runner.E2ERunner,
	withdrawTxs <-chan withdrawResult,
	results chan<- withdrawResult,
	bestEffortMode bool,
	workerID int,
) {
	for withdraw := range withdrawTxs {
		result := withdraw

		// First, wait for the transaction receipt on zEVM
		ctx, cancel := context.WithTimeout(context.Background(), r.ReceiptTimeout)
		receipt, err := waitForZEVMReceipt(ctx, r.ZEVMClient, withdraw.txHash)
		cancel()

		if err != nil {
			result.err = fmt.Errorf("failed to get zEVM receipt: %w", err)
			r.Logger.Print("worker %d: index %d: ✗ failed to get receipt: %v",
				workerID, withdraw.index, err)
			results <- result
			continue
		}

		if receipt.Status == 0 {
			result.err = fmt.Errorf("transaction reverted on zEVM")
			r.Logger.Print("worker %d: index %d: ✗ transaction reverted",
				workerID, withdraw.index)
			results <- result
			continue
		}

		// Now wait for CCTX to be mined
		cctx := utils.WaitCctxMinedByInboundHash(
			r.Ctx,
			withdraw.txHash.Hex(),
			r.CctxClient,
			r.Logger,
			r.CctxTimeout,
		)

		result.cctx = cctx
		result.duration = time.Since(withdraw.startTime)

		if cctx.CctxStatus.Status == crosschaintypes.CctxStatus_OutboundMined {
			if withdraw.index%10 == 0 || !bestEffortMode {
				r.Logger.Print("worker %d: index %d: ✓ completed in %v",
					workerID, withdraw.index, result.duration)
			}
		} else {
			result.err = fmt.Errorf("cctx failed with status %s: %s",
				cctx.CctxStatus.Status, cctx.CctxStatus.StatusMessage)
			r.Logger.Print("worker %d: index %d: ✗ failed with status %s",
				workerID, withdraw.index, cctx.CctxStatus.Status)
		}

		results <- result
	}
}

// waitForZEVMReceipt polls for a transaction receipt on zEVM until it's available or timeout
func waitForZEVMReceipt(ctx context.Context, client bind.DeployBackend, hash common.Hash) (*types.Receipt, error) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			receipt, err := client.TransactionReceipt(ctx, hash)
			if err == nil {
				return receipt, nil
			}
			// If error is not "not found", return it
			if err.Error() != "not found" && err.Error() != "transaction not found" {
				return nil, err
			}
			// Otherwise continue polling
		}
	}
}

// minInt returns the minimum of two integers
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// isErrMempoolFull checks if the error indicates mempool is full
// there are two types of error messages reported when it happens
func isErrMempoolFull(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, "mempool is full") ||
		strings.Contains(errMsg, "pool reached max tx capacity")
}
