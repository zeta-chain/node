package e2etests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/testutil/sample"
	zetabitcoin "github.com/zeta-chain/node/zetaclient/chains/bitcoin"
)

func TestBitcoinDepositAndCallRevert(r *runner.E2ERunner, args []string) {
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

	r.Logger.Info("Sent %f BTC to TSS with invalid memo, got refund of %d satoshis", amount, value)
}
