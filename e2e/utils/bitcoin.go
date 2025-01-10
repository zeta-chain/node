package utils

import (
	"context"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/stretchr/testify/require"
)

// MustHaveDroppedTx ensures the given tx has been dropped
func MustHaveDroppedTx(ctx context.Context, client *rpcclient.Client, txHash *chainhash.Hash) {
	t := TestingFromContext(ctx)

	// dropped tx has negative confirmations
	txResult, err := client.GetTransaction(txHash)
	if err == nil {
		require.Negative(t, txResult.Confirmations)
	}

	// dropped tx should be removed from mempool
	entry, err := client.GetMempoolEntry(txHash.String())
	require.Error(t, err)
	require.Nil(t, entry)

	// dropped tx won't exist in blockchain
	// -5: No such mempool or blockchain transaction
	rawTx, err := client.GetRawTransaction(txHash)
	require.Error(t, err)
	require.Nil(t, rawTx)
}

// MustHaveMinedTx ensures the given tx has been mined
func MustHaveMinedTx(ctx context.Context, client *rpcclient.Client, txHash *chainhash.Hash) *btcjson.TxRawResult {
	t := TestingFromContext(ctx)

	// positive confirmations
	txResult, err := client.GetTransaction(txHash)
	require.NoError(t, err)
	require.Positive(t, txResult.Confirmations)

	// tx exists in blockchain
	rawResult, err := client.GetRawTransactionVerbose(txHash)
	require.NoError(t, err)

	return rawResult
}
