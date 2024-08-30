package e2etests

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/x/crosschain/types"
	zetabitcoin "github.com/zeta-chain/node/zetaclient/chains/bitcoin"
)

func TestBitcoinDepositRefund(r *runner.E2ERunner, args []string) {
	// ARRANGE
	// Given BTC address
	r.SetBtcAddress(r.Name, false)

	// Given "Live" BTC network
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	// Given amount to send
	require.Len(r, args, 1)
	amount := parseFloat(r, args[0])
	amount += zetabitcoin.DefaultDepositorFee

	// Given a list of UTXOs
	utxos, err := r.ListDeployerUTXOs()
	require.NoError(r, err)
	require.NotEmpty(r, utxos)

	// ACT
	// Send BTC to TSS address with a dummy memo
	txHash, err := r.SendToTSSFromDeployerWithMemo(amount, utxos, []byte("gibberish-memo"))
	require.NoError(r, err)
	require.NotEmpty(r, txHash)

	// ASSERT
	// Now we want to make sure refund TX is completed.
	// Let's check that zetaclient issued a refund on BTC
	searchForCrossChainWithBtcRefund := utils.Matches(func(tx types.CrossChainTx) bool {
		return tx.GetCctxStatus().Status == types.CctxStatus_Reverted &&
			len(tx.OutboundParams) == 2 &&
			tx.OutboundParams[1].Hash != ""
	})

	cctxs := utils.WaitCctxByInboundHash(r.Ctx, r, txHash.String(), r.CctxClient, searchForCrossChainWithBtcRefund)
	require.Len(r, cctxs, 1)

	// Pick btc tx hash from the cctx
	btcTxHash, err := chainhash.NewHashFromStr(cctxs[0].OutboundParams[1].Hash)
	require.NoError(r, err)

	// Query the BTC network to check the refund transaction
	refundTx, err := r.BtcRPCClient.GetTransaction(btcTxHash)
	require.NoError(r, err, refundTx)

	// Finally, check the refund transaction details
	refundTxDetails := refundTx.Details[0]
	assert.Equal(r, "receive", refundTxDetails.Category)
	assert.Equal(r, r.BTCDeployerAddress.EncodeAddress(), refundTxDetails.Address)
	assert.NotEmpty(r, refundTxDetails.Amount)

	r.Logger.Info("Sent %f BTC to TSS with invalid memo, got refund of %f BTC", amount, refundTxDetails.Amount)
}
