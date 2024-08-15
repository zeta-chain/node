package e2etests

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	zetabitcoin "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin"
	btcobserver "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin/observer"

	"github.com/zeta-chain/zetacore/e2e/runner"
)

func TestExtractBitcoinInscriptionMemo(r *runner.E2ERunner, args []string) {
	r.Logger.Info("Testing extract memo from btc inscription")

	r.SetBtcAddress(r.Name, false)

	// obtain some initial fund
	stop := r.MineBlocksIfLocalBitcoin()
	defer stop()

	// list deployer utxos
	utxos, err := r.ListDeployerUTXOs()
	require.NoError(r, err)

	amount := parseFloat(r, args[0])
	memo, _ := hex.DecodeString(
		"72f080c854647755d0d9e6f6821f6931f855b9acffd53d87433395672756d58822fd143360762109ab898626556b1c3b8d3096d2361f1297df4a41c1b429471a9aa2fc9be5f27c13b3863d6ac269e4b587d8389f8fd9649859935b0d48dea88cdb40f20c",
	)

	txid, err := r.InscribeToTSSFromDeployerWithMemo(amount, utxos, memo)
	require.NoError(r, err)

	_, err = r.GenerateToAddressIfLocalBitcoin(6, r.BTCDeployerAddress)
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
