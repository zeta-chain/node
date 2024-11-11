package e2etests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

// TestBitcoinDepositAndCallRevertWithDust sends a Bitcoin deposit that reverts with a dust amount in the revert outbound.
// Given the dust is too smart, the CCTX should revert
func TestBitcoinDepositAndCallRevertWithDust(r *runner.E2ERunner, args []string) {
	// ARRANGE
	// Given BTC address
	r.SetBtcAddress(r.Name, false)

	require.Len(r, args, 0)

	// Given "Live" BTC network
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	// 0.002 BTC is consumed in a deposit and revert, the dust is set to 1000 satoshis in the protocol
	// Therefore the deposit amount should be 0.002 + 0.000005 = 0.00200500 should trigger the condition
	// As only 500 satoshis are left after the deposit

	amount := 0.00200500

	// Given a list of UTXOs
	utxos, err := r.ListDeployerUTXOs()
	require.NoError(r, err)
	require.NotEmpty(r, utxos)

	// ACT
	// Send BTC to TSS address with a dummy memo
	// zetacore should revert cctx if call is made on a non-existing address
	nonExistReceiver := sample.EthAddress()
	badMemo := append(nonExistReceiver.Bytes(), []byte("gibberish-memo")...)
	txHash, err := r.SendToTSSFromDeployerWithMemo(amount, utxos, badMemo)
	require.NoError(r, err)
	require.NotEmpty(r, txHash)

	// ASSERT
	// Now we want to make sure refund TX is completed.
	cctx := utils.WaitCctxRevertedByInboundHash(r.Ctx, r, txHash.String(), r.CctxClient)

	// Check revert tx receiver address and amount
	receiver, value := r.QueryOutboundReceiverAndAmount(cctx.OutboundParams[1].Hash)
	assert.Equal(r, r.BTCDeployerAddress.EncodeAddress(), receiver)
	assert.Positive(r, value)

	r.Logger.Print("BITCOIN: Amount received: %d", value)
	r.Logger.Info("Sent %f BTC to TSS with invalid memo, got refund of %d satoshis", amount, value)
}
