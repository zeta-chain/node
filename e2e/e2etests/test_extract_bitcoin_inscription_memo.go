package e2etests

import (
	"github.com/stretchr/testify/require"
	zetabitcoin "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin"
	btcobserver "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/observer"

	"github.com/zeta-chain/zetacore/e2e/runner"
)

func TestExtractBitcoinInscriptionMemo(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	r.SetBtcAddress(r.Name, false)

	r.InscribeToTSSFromDeployerWithMemo()
	_, err := r.GenerateToAddressIfLocalBitcoin(6, r.BTCDeployerAddress)
	require.NoError(r, err)
	gtx, err := r.BtcRPCClient.GetTransaction(txid)
	require.NoError(r, err)
	r.Logger.Info("rawtx confirmation: %d", gtx.BlockIndex)
	rawtx, err := r.BtcRPCClient.GetRawTransactionVerbose(txid)
	require.NoError(r, err)

	depositorFee := zetabitcoin.DefaultDepositorFee
	events, err := btcobserver.FilterAndParseIncomingTx(
		r.BtcRPCClient,
		[]btcjson.TxRawResult{*rawtx},
		0,
		r.BTCTSSAddress.EncodeAddress(),
		log.Logger,
		r.BitcoinParams,
		depositorFee,
	)
	require.NoError(r, err)
	r.Logger.Info("bitcoin inbound events:")
	for _, event := range events {
		r.Logger.Info("  TxHash: %s", event.TxHash)
		r.Logger.Info("  From: %s", event.FromAddress)
		r.Logger.Info("  To: %s", event.ToAddress)
		r.Logger.Info("  Amount: %f", event.Value)
		r.Logger.Info("  Memo: %x", event.MemoBytes)
	}
}
