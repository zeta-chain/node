package e2etests

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/e2e/contracts/testgasconsumer"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

const (
	// ModeFailFast fails the test on first transaction failure
	ModeFailFast = false
	// ModeBestEffort tolerates transaction failures and reports statistics at the end
	ModeBestEffort = true
)

type txResult struct {
	index   uint64
	hash    common.Hash
	receipt *types.Receipt
	err     error
}

// TestStressZEVM tests stressing direct interactions with the zEVM using calls that consume a lot of gas
func TestStressZEVM(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse number of transactions
	txNumbers := utils.ParseBigInt(r, args[0])
	totalTxs := txNumbers.Uint64()

	// configure test mode
	bestEffortMode := ModeBestEffort

	r.Logger.Print("starting stress test of %d calls (mode: best-effort=%v)", totalTxs, bestEffortMode)

	// Deploy the GasConsumer contract
	gasConsumerAddress, txDeploy, gasConsumer, err := testgasconsumer.DeployTestGasConsumer(
		r.ZEVMAuth,
		r.ZEVMClient,
		big.NewInt(1000000),
	)
	require.NoError(r, err)
	r.WaitForTxReceiptOnZEVM(txDeploy)

	// Get initial nonce
	initialNonce, err := r.ZEVMClient.PendingNonceAt(context.Background(), r.ZEVMAuth.From)
	require.NoError(r, err)

	r.Logger.Print("starting nonce: %d", initialNonce)

	// Channels for transaction tracking
	txHashes := make(chan txResult, totalTxs)
	results := make(chan txResult, totalTxs)

	// Statistics tracking
	var (
		sentCount    atomic.Uint64
		successCount atomic.Uint64
		failedCount  atomic.Uint64
	)

	// Start receipt monitoring goroutines
	var wg sync.WaitGroup
	numMonitors := 10 // number of parallel receipt monitors
	for i := 0; i < numMonitors; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			monitorReceipts(r, txHashes, results, bestEffortMode, workerID)
		}(i)
	}

	// Send all transactions with pre-calculated nonces
	sendStart := time.Now()
	for i := uint64(0); i < totalTxs; i++ {
		nonce := initialNonce + i

		// Create a new transactor with specific nonce
		auth := *r.ZEVMAuth // copy the auth
		auth.Nonce = big.NewInt(int64(nonce))

		// Send transaction
		tx, err := gasConsumer.OnCall(
			&auth,
			testgasconsumer.TestGasConsumerzContext{
				Origin:  []byte{},
				Sender:  gasConsumerAddress,
				ChainID: big.NewInt(0),
			},
			gasConsumerAddress,
			big.NewInt(0),
			[]byte{},
		)

		if err != nil {
			r.Logger.Print("index %d (nonce %d): failed to send tx: %v", i, nonce, err)
			if !bestEffortMode {
				require.FailNow(r, fmt.Sprintf("failed to send transaction %d: %v", i, err))
			}
			failedCount.Add(1)
			continue
		}

		sentCount.Add(1)
		txHashes <- txResult{index: i, hash: tx.Hash()}

		if i%100 == 0 {
			r.Logger.Print("sent %d/%d transactions", i+1, totalTxs)
		}

		// Small delay to avoid overwhelming the node
		time.Sleep(time.Millisecond * 10)
	}

	close(txHashes)
	sendDuration := time.Since(sendStart)

	r.Logger.Print("all %d transactions sent in %v", sentCount.Load(), sendDuration)
	r.Logger.Print("waiting for receipts...")

	// Wait for all monitors to finish
	wg.Wait()
	close(results)

	// Collect and analyze results
	receiptStart := time.Now()
	var failedTxs []txResult
	for result := range results {
		if result.err != nil || result.receipt == nil {
			failedCount.Add(1)
			failedTxs = append(failedTxs, result)
			if !bestEffortMode {
				require.FailNow(r, fmt.Sprintf("transaction %d failed: %v", result.index, result.err))
			}
		} else if result.receipt.Status == types.ReceiptStatusFailed {
			failedCount.Add(1)
			failedTxs = append(failedTxs, result)
			if !bestEffortMode {
				require.FailNow(r, fmt.Sprintf("transaction %d reverted", result.index))
			}
		} else {
			successCount.Add(1)
		}
	}
	receiptDuration := time.Since(receiptStart)

	// Print final statistics
	r.Logger.Print("═══════════════════════════════════════")
	r.Logger.Print("Stress Test Results:")
	r.Logger.Print("═══════════════════════════════════════")
	r.Logger.Print("Total transactions: %d", totalTxs)
	r.Logger.Print("Sent successfully:  %d", sentCount.Load())
	r.Logger.Print("Succeeded:          %d (%.2f%%)", successCount.Load(), float64(successCount.Load())/float64(totalTxs)*100)
	r.Logger.Print("Failed:             %d (%.2f%%)", failedCount.Load(), float64(failedCount.Load())/float64(totalTxs)*100)
	r.Logger.Print("Send duration:      %v", sendDuration)
	r.Logger.Print("Receipt duration:   %v", receiptDuration)
	r.Logger.Print("Total duration:     %v", sendDuration+receiptDuration)
	r.Logger.Print("TPS (send):         %.2f", float64(sentCount.Load())/sendDuration.Seconds())
	r.Logger.Print("═══════════════════════════════════════")

	if len(failedTxs) > 0 && len(failedTxs) <= 10 {
		r.Logger.Print("Failed transaction details:")
		for _, failed := range failedTxs {
			r.Logger.Print("  - Index %d, Hash %s: %v", failed.index, failed.hash.Hex(), failed.err)
		}
	} else if len(failedTxs) > 10 {
		r.Logger.Print("First 10 failed transactions:")
		for i := 0; i < 10; i++ {
			r.Logger.Print("  - Index %d, Hash %s: %v", failedTxs[i].index, failedTxs[i].hash.Hex(), failedTxs[i].err)
		}
		r.Logger.Print("  ... and %d more", len(failedTxs)-10)
	}

	if !bestEffortMode && failedCount.Load() > 0 {
		require.FailNow(r, fmt.Sprintf("%d transactions failed", failedCount.Load()))
	}

	r.Logger.Print("stress test completed")
}

// monitorReceipts monitors transaction receipts in a goroutine
func monitorReceipts(
	r *runner.E2ERunner,
	txHashes <-chan txResult,
	results chan<- txResult,
	bestEffortMode bool,
	workerID int,
) {
	for txRes := range txHashes {
		result := txRes

		// Poll for receipt with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		receipt, err := waitForReceipt(ctx, r.ZEVMClient, txRes.hash)
		cancel()

		if err != nil {
			result.err = fmt.Errorf("failed to get receipt: %w", err)
			if !bestEffortMode {
				r.Logger.Print("worker %d: index %d: failed to get receipt: %v", workerID, txRes.index, err)
			}
		} else {
			result.receipt = receipt
			if receipt.Status == types.ReceiptStatusSuccessful {
				if txRes.index%100 == 0 {
					r.Logger.Print("worker %d: index %d: ✓ confirmed (gas used: %d)", workerID, txRes.index, receipt.GasUsed)
				}
			} else {
				r.Logger.Print("worker %d: index %d: ✗ reverted", workerID, txRes.index)
			}
		}

		results <- result
	}
}

// waitForReceipt polls for a transaction receipt until it's available or timeout
func waitForReceipt(ctx context.Context, client bind.DeployBackend, hash common.Hash) (*types.Receipt, error) {
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
