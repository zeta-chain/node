package e2etests

import (
	"context"
	"strconv"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	zetabitcoin "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin"
)

func TestBitcoinDepositRefund(r *runner.E2ERunner, args []string) {
	// ARRANGE
	// Given amount to send
	require.Len(r, args, 1)
	amount := parseFloat(r, args[0])
	amount += zetabitcoin.DefaultDepositorFee

	// Given BTC address
	r.SetBtcAddress(r.Name, false)

	// Given a list of UTXOs
	utxos, err := r.ListDeployerUTXOs()
	require.NoError(r, err)
	require.NotEmpty(r, utxos)

	// ACT
	// Send BTC to TSS address with a dummy memo
	txHash, err := r.SendToTSSFromDeployerWithMemo(amount, utxos, []byte("gibberish-memo"))
	require.NoError(r, err)
	require.NotEmpty(r, txHash)

	// Wait for processing in zetaclient
	cctxs := utils.WaitCctxByInboundHash(r.Ctx, r, txHash.String(), r.CctxClient)

	// ASSERT
	require.Len(r, cctxs, 1)

	// Check that it's status is related to tx reversal
	expectedStatuses := []string{types.CctxStatus_PendingRevert.String(), types.CctxStatus_Reverted.String()}
	actualStatus := cctxs[0].CctxStatus.Status.String()

	require.Contains(r, expectedStatuses, actualStatus)

	r.Logger.Info("CCTX revert status: %s", actualStatus)

	// Now we want to make sure refund TX is completed. Let's check that zetaclient issued a refund on BTC
	ctx, cancel := context.WithTimeout(r.Ctx, time.Minute*10)
	defer cancel()

	searchForCrossChainWithBtcRefund := utils.Matches(func(tx types.CrossChainTx) bool {
		if len(tx.OutboundParams) != 2 {
			return false
		}

		btcRefundTx := tx.OutboundParams[1]

		return btcRefundTx.Hash != ""
	})

	cctxs = utils.WaitCctxByInboundHash(ctx, r, txHash.String(), r.CctxClient, searchForCrossChainWithBtcRefund)
	require.Len(r, cctxs, 1)

	// todo check that BTC refund is completed
}

func parseFloat(t require.TestingT, s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	require.NoError(t, err, "unable to parse float %q", s)
	return f
}
