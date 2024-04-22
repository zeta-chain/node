package e2etests

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"golang.org/x/sync/errgroup"
)

// rateLimiterFlags are the rate limiter flags for the test
var rateLimiterFlags = crosschaintypes.RateLimiterFlags{
	Enabled: true,
	Rate:    sdk.NewUint(110000000000000000), // 0.1 ZETA, this value is used so rate is reached
	Window:  5,
}

func TestRateLimiter(r *runner.E2ERunner, _ []string) {
	r.Logger.Info("TestRateLimiter")

	// First test without rate limiter
	r.Logger.Print("rate limiter disabled")
	if err := createAndWaitWithdraws(r); err != nil {
		panic(err)
	}

	// Set the rate limiter to 0.11ZETA per 10 blocks
	// These rate limiter flags will only allow to process 1 withdraw per 10 blocks
	r.Logger.Info("setting up rate limiter flags")
	if err := setupRateLimiterFlags(r, rateLimiterFlags); err != nil {
		panic(err)
	}

	// Test with rate limiter
	r.Logger.Print("rate limiter enabled")
	if err := createAndWaitWithdraws(r); err != nil {
		panic(err)
	}

	// Disable rate limiter
	r.Logger.Info("disabling rate limiter")
	if err := setupRateLimiterFlags(r, crosschaintypes.RateLimiterFlags{Enabled: false}); err != nil {
		panic(err)
	}

	// Test without rate limiter again
	r.Logger.Print("rate limiter disabled")
	if err := createAndWaitWithdraws(r); err != nil {
		panic(err)
	}
}

// setupRateLimiterFlags sets up the rate limiter flags with flags defined in the test
func setupRateLimiterFlags(r *runner.E2ERunner, flags crosschaintypes.RateLimiterFlags) error {
	adminAddr, err := r.ZetaTxServer.GetAccountAddressFromName(utils.FungibleAdminName)
	if err != nil {
		return err
	}
	_, err = r.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, crosschaintypes.NewMsgUpdateRateLimiterFlags(
		adminAddr,
		flags,
	))
	if err != nil {
		return err
	}

	return nil
}

// createAndWaitWithdraws performs 10 withdraws
func createAndWaitWithdraws(r *runner.E2ERunner) error {
	startTime := time.Now()

	r.Logger.Print("starting 10 withdraws of 0.1 ZETA each")

	// Perform 10 withdraws to log time for completion
	txs := make([]*ethtypes.Transaction, 10)
	for i := 0; i < 10; i++ {
		txs[i] = r.WithdrawZeta("100000000000000000", false)
	}

	// start a error group to wait for all the withdraws to be mined
	g, ctx := errgroup.WithContext(r.Ctx)
	for i, tx := range txs {
		// capture the loop variables
		tx, i := tx, i

		// start a goroutine to wait for the withdraw to be mined
		g.Go(func() error {
			return waitForZetaWithdrawMined(ctx, r, tx, i, startTime)
		})
	}

	// wait for all the withdraws to be mined
	if err := g.Wait(); err != nil {
		return err
	}

	duration := time.Now().Sub(startTime).Seconds()
	block, err := r.ZEVMClient.BlockNumber(ctx)
	if err != nil {
		return err
	}
	r.Logger.Print("all 10 withdraws completed in %vs at block %d", duration, block)

	return nil
}

// waitForZetaWithdrawMined waits for a zeta withdraw to be mined
// we first wait to get the receipt
// NOTE: this could be a more general function but we define it here for this test because we emit in the function logs specific to this test
func waitForZetaWithdrawMined(ctx context.Context, r *runner.E2ERunner, tx *ethtypes.Transaction, index int, startTime time.Time) error {
	// wait for the tx receipt
	receipt := utils.MustWaitForTxReceipt(ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "zeta withdraw")
	if receipt.Status == 0 {
		return fmt.Errorf("zeta withdraw failed, logs: %+v", receipt.Logs)
	}

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInTxHash(ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "zeta withdraw")
	if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_OutboundMined {
		return fmt.Errorf(
			"expected cctx status to be %s; got %s, message %s",
			crosschaintypes.CctxStatus_OutboundMined,
			cctx.CctxStatus.Status.String(),
			cctx.CctxStatus.StatusMessage,
		)
	}

	// record the time for completion
	duration := time.Now().Sub(startTime).Seconds()
	block, err := r.ZEVMClient.BlockNumber(ctx)
	if err != nil {
		return err
	}
	r.Logger.Print("cctx %d mined in %vs at block %d", index, duration, block)

	return nil
}
