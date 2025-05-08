package utils

import (
	"context"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/client"
)

// MustHaveDroppedBitcoinTx ensures the given Bitcoin tx has been dropped
func MustHaveDroppedBitcoinTx(ctx context.Context, client *client.Client, txHash *chainhash.Hash) {
	t := TestingFromContext(ctx)

	// dropped tx has negative confirmations
	txResult, err := client.GetTransaction(ctx, txHash)
	if err == nil {
		require.Negative(t, txResult.Confirmations)
	}

	// dropped tx should be removed from mempool
	entry, err := client.GetMempoolEntry(ctx, txHash.String())
	require.Error(t, err)
	require.Nil(t, entry)

	// dropped tx won't exist in blockchain
	// -5: No such mempool or blockchain transaction
	rawTx, err := client.GetRawTransaction(ctx, txHash)
	require.Error(t, err)
	require.Nil(t, rawTx)
}

// MustHaveMinedBitcoinTx ensures the given Bitcoin tx has been mined
func MustHaveMinedBitcoinTx(ctx context.Context, client *client.Client, txHash *chainhash.Hash) *btcjson.TxRawResult {
	t := TestingFromContext(ctx)

	// positive confirmations
	txResult, err := client.GetTransaction(ctx, txHash)
	require.NoError(t, err)
	require.Positive(t, txResult.Confirmations)

	// tx exists in blockchain
	rawResult, err := client.GetRawTransactionVerbose(ctx, txHash)
	require.NoError(t, err)

	return rawResult
}
