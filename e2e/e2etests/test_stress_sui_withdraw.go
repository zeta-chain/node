package e2etests

import (
	"math/big"
	"sync"

	"cosmossdk.io/errors"
	"github.com/montanaflynn/stats"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"
	"golang.org/x/sync/errgroup"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

// TestStressSuiWithdraw tests the stressing withdrawal of SUI
func TestStressSuiWithdraw(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	// ARRANGE
	// Given amount
	amount := utils.ParseBigInt(r, args[0])
	numWithdrawals := utils.ParseInt(r, args[1])

	// Given signer
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err, "unable to get deployer signer")

	r.Logger.Print("starting stress test of %d SUI withdrawals", numWithdrawals)

	// approve enough ZRC20 SUI token for the withdrawals
	r.ApproveSUIZRC20(r.GatewayZEVMAddr)

	// create a wait group to wait for all the withdrawals to complete
	var eg errgroup.Group

	// store durations as float64 seconds like prometheus
	withdrawDurations := make([]float64, 0, numWithdrawals)
	mu := sync.Mutex{}

	// ACT
	// send the withdrawals SUI
	revertOptions := gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)}
	for i := range numWithdrawals {
		// each goroutine captures its own copy of i
		i := i
		tx := r.SuiWithdraw(signer.Address(), amount, r.SUIZRC20Addr, revertOptions)

		// wait for receipt before next withdrawal to avoid race condition
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
		utils.RequireTxSuccessful(r, receipt)

		r.Logger.Print("index %d: started with tx hash: %s", i, tx.Hash().Hex())

		eg.Go(func() error {
			totalTime, err := r.SuiMonitorCCTXByInboundHash(tx.Hash().Hex(), i)
			if err != nil {
				return errors.Wrap(err, "failed to monitor withdraw")
			}

			mu.Lock()
			withdrawDurations = append(withdrawDurations, totalTime.Seconds())
			mu.Unlock()

			return nil
		})
	}

	err = eg.Wait()

	// Print statistics
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
	r.Logger.Print("All SUI withdrawals completed")
}
