package observer

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"golang.org/x/exp/rand"

	"github.com/zeta-chain/node/pkg/chains"
)

func Test_FetchUTXOs(t *testing.T) {
	// create test suite
	ob, utxos := newTestSuitWithUTXOs(t)

	// check number of UTXOs again
	require.Equal(t, len(utxos), ob.TelemetryServer().GetNumberOfUTXOs())
}

func Test_SelectUTXOs(t *testing.T) {
	ctx := context.Background()
	dummyTxID := "6e6f71d281146c1fc5c755b35908ee449f26786c84e2ae18f98b268de40b7ec4"

	t.Run("noce = 0, should bootstrap", func(t *testing.T) {
		// 		input: utxoCap = 5, amount = 0.01, nonce = 0
		// 		output: [0.01], 0.01
		ob, utxos := newTestSuitWithUTXOs(t)
		selected, err := ob.SelectUTXOs(ctx, 0.01, 5, 0, math.MaxUint16)
		require.NoError(t, err)
		require.Equal(t, 0.01, selected.Value)
		require.Equal(t, utxos[0:1], selected.UTXOs)
	})

	t.Run("nonce = 1, must FAIL and wait for previous transaction to be mined", func(t *testing.T) {
		// 		input: utxoCap = 5, amount = 0.5, nonce = 1
		// 		output: error
		ob, _ := newTestSuitWithUTXOs(t)
		selected, err := ob.SelectUTXOs(ctx, 0.5, 5, 1, math.MaxUint16)
		require.Error(t, err)
		require.Nil(t, selected.UTXOs)
		require.Zero(t, selected.Value)
		require.ErrorContains(t, err, "error getting cctx for nonce 0")
	})

	t.Run("nonce = 1, should pass when nonce mark 0 is set", func(t *testing.T) {
		// 		input: utxoCap = 5, amount = 0.5, nonce = 1
		// 		output: [0.00002, 0.01, 0.12, 0.18, 0.24], 0.55002
		ob, utxos := createTestSuitWithUTXOsAndNonceMark(t, 0, dummyTxID)
		selected, err := ob.SelectUTXOs(ctx, 0.5, 5, 1, math.MaxUint16)
		require.NoError(t, err)
		require.Equal(t, 0.55002, selected.Value)
		require.Equal(t, utxos[0:5], selected.UTXOs)
	})

	t.Run("nonce = 2, should pass when nonce mark 1 is set", func(t *testing.T) {
		// 		input: utxoCap = 5, amount = 1.0, nonce = 2
		// 		output: [0.00002001, 0.01, 0.12, 0.18, 0.24, 0.5], 1.05002001
		ob, utxos := createTestSuitWithUTXOsAndNonceMark(t, 1, dummyTxID)
		selected, err := ob.SelectUTXOs(ctx, 1.0, 5, 2, math.MaxUint16)
		require.NoError(t, err)
		require.InEpsilon(t, 1.05002001, selected.Value, 1e-8)
		require.Equal(t, utxos[0:6], selected.UTXOs)
	})

	t.Run("nonce = 3, should select nonce-mark utxo on the LEFT", func(t *testing.T) {
		// 		input: utxoCap = 5, amount = 8.05, nonce = 3
		// 		output: [0.00002002, 0.24, 0.5, 1.26, 2.97, 3.28], 8.25002002
		ob, utxos := createTestSuitWithUTXOsAndNonceMark(t, 2, dummyTxID)
		selected, err := ob.SelectUTXOs(ctx, 8.05, 5, 3, math.MaxUint16)
		require.NoError(t, err)
		require.InEpsilon(t, 8.25002002, selected.Value, 1e-8)
		expected := append([]btcjson.ListUnspentResult{utxos[0]}, utxos[4:9]...)
		require.Equal(t, expected, selected.UTXOs)
	})

	t.Run("nonce = 24105432, should select nonce-mark utxo on the RIGHT", func(t *testing.T) {
		// 		input: utxoCap = 5, amount = 0.503, nonce = 24105432
		// 		output: [0.24107432, 0.01, 0.12, 0.18, 0.24], 0.7910731
		ob, utxos := createTestSuitWithUTXOsAndNonceMark(t, 24105431, dummyTxID)
		selected, err := ob.SelectUTXOs(ctx, 0.503, 5, 24105432, math.MaxUint16)
		require.NoError(t, err)
		require.InEpsilon(t, 0.79107431, selected.Value, 1e-8)
		expected := append([]btcjson.ListUnspentResult{utxos[4]}, utxos[0:4]...)
		require.Equal(t, expected, selected.UTXOs)
	})

	t.Run("nonce = 24105433, should select nonce-mark utxo in the MIDDLE", func(t *testing.T) {
		// 		input: utxoCap = 5, amount = 1.0, nonce = 24105433
		// 		output: [0.24107432, 0.12, 0.18, 0.24, 0.5], 1.28107432
		ob, utxos := createTestSuitWithUTXOsAndNonceMark(t, 24105432, dummyTxID)
		selected, err := ob.SelectUTXOs(ctx, 1.0, 5, 24105433, math.MaxUint16)
		require.NoError(t, err)
		require.InEpsilon(t, 1.28107432, selected.Value, 1e-8)
		expected := append([]btcjson.ListUnspentResult{utxos[4]}, utxos[1:4]...)
		expected = append(expected, utxos[5])
		require.Equal(t, expected, selected.UTXOs)
	})

	t.Run("nonce = 24105433, should select biggest utxos to maximize amount", func(t *testing.T) {
		// 		input: utxoCap = 5, amount = 16.03
		// 		output: [0.24107432, 1.26, 2.97, 3.28, 5.16, 8.72], 21.63107432
		ob, utxos := createTestSuitWithUTXOsAndNonceMark(t, 24105432, dummyTxID)
		selected, err := ob.SelectUTXOs(ctx, 16.03, 5, 24105433, math.MaxUint16)
		require.NoError(t, err)
		require.InEpsilon(t, 21.63107432, selected.Value, 1e-8)
		expected := append([]btcjson.ListUnspentResult{utxos[4]}, utxos[6:11]...)
		require.Equal(t, expected, selected.UTXOs)
	})

	t.Run("nonce = 24105433, should fail due to insufficient funds", func(t *testing.T) {
		// 		input: utxoCap = 5, amount = 21.64
		// 		output: error
		ob, _ := createTestSuitWithUTXOsAndNonceMark(t, 24105432, dummyTxID)
		selected, err := ob.SelectUTXOs(ctx, 21.64, 5, 24105433, math.MaxUint16)
		require.Error(t, err)
		require.Nil(t, selected.UTXOs)
		require.Zero(t, selected.Value)
		require.ErrorContains(t, err, "not enough btc in reserve - available : 21.63107432 , tx amount : 21.64")
	})
}

func Test_SelectUTXOs_Consolidation(t *testing.T) {
	ctx := context.Background()
	dummyTxID := "6e6f71d281146c1fc5c755b35908ee449f26786c84e2ae18f98b268de40b7ec4"

	t.Run("should not consolidate", func(t *testing.T) {
		// create test suite and set nonce-mark utxo for nonce 0
		ob, utxos := createTestSuitWithUTXOsAndNonceMark(t, 0, dummyTxID)

		// input: utxoCap = 10, amount = 0.01, nonce = 1, rank = 10
		// output: [0.00002, 0.01], 0.01002
		res, err := ob.SelectUTXOs(ctx, 0.01, 10, 1, 10)
		require.NoError(t, err)
		require.Equal(t, 0.01002, res.Value)
		require.Equal(t, utxos[0:2], res.UTXOs)
		require.Zero(t, res.ConsolidatedUTXOs)
		require.Zero(t, res.ConsolidatedValue)
	})

	t.Run("should consolidate 1 utxo", func(t *testing.T) {
		// create test suite and set nonce-mark utxo for nonce 0
		ob, utxos := createTestSuitWithUTXOsAndNonceMark(t, 0, dummyTxID)

		// input: utxoCap = 9, amount = 0.01, nonce = 1, rank = 9
		// output: [0.00002, 0.01, 0.12], 0.13002
		selected, err := ob.SelectUTXOs(ctx, 0.01, 9, 1, 9)
		require.NoError(t, err)
		require.Equal(t, 0.13002, selected.Value)
		require.Equal(t, utxos[0:3], selected.UTXOs)
		require.Equal(t, uint16(1), selected.ConsolidatedUTXOs)
		require.Equal(t, 0.12, selected.ConsolidatedValue)
	})

	t.Run("should consolidate 3 utxos", func(t *testing.T) {
		// create test suite and set nonce-mark utxo for nonce 0
		ob, utxos := createTestSuitWithUTXOsAndNonceMark(t, 0, dummyTxID)

		// input: utxoCap = 5, amount = 0.01, nonce = 0, rank = 5
		// output: [0.00002, 0.014, 1.26, 0.5, 0.2], 2.01002
		selected, err := ob.SelectUTXOs(ctx, 0.01, 5, 1, 5)
		require.NoError(t, err)
		require.Equal(t, 2.01002, selected.Value)
		expected := make([]btcjson.ListUnspentResult, 2)
		copy(expected, utxos[0:2])
		for i := 6; i >= 4; i-- { // append consolidated utxos in descending order
			expected = append(expected, utxos[i])
		}
		require.Equal(t, expected, selected.UTXOs)
		require.Equal(t, uint16(3), selected.ConsolidatedUTXOs)
		require.Equal(t, 2.0, selected.ConsolidatedValue)
	})

	t.Run("should consolidate all utxos using rank 1", func(t *testing.T) {
		// create test suite and set nonce-mark utxo for nonce 0
		ob, utxos := createTestSuitWithUTXOsAndNonceMark(t, 0, dummyTxID)

		// input: utxoCap = 12, amount = 0.01, nonce = 0, rank = 1
		// output: [0.00002, 0.01, 8.72, 5.16, 3.28, 2.97, 1.26, 0.5, 0.24, 0.18, 0.12], 22.44002
		selected, err := ob.SelectUTXOs(ctx, 0.01, 12, 1, 1)
		require.NoError(t, err)
		require.Equal(t, 22.44002, selected.Value)
		expected := make([]btcjson.ListUnspentResult, 2)
		copy(expected, utxos[0:2])
		for i := 10; i >= 2; i-- { // append consolidated utxos in descending order
			expected = append(expected, utxos[i])
		}
		require.Equal(t, expected, selected.UTXOs)
		require.Equal(t, uint16(9), selected.ConsolidatedUTXOs)
		require.Equal(t, 22.43, selected.ConsolidatedValue)
	})

	t.Run("should consolidate 3 utxos sparse", func(t *testing.T) {
		// create test suite and set nonce-mark utxo for nonce 24105431
		ob, utxos := createTestSuitWithUTXOsAndNonceMark(t, 24105431, dummyTxID)

		// input: utxoCap = 5, amount = 0.13, nonce = 24105432, rank = 5
		// output: [0.24107431, 0.01, 0.12, 1.26, 0.5, 0.24], 2.37107431
		selected, err := ob.SelectUTXOs(ctx, 0.13, 5, 24105432, 5)
		require.NoError(t, err)
		require.InEpsilon(t, 2.37107431, selected.Value, 1e-8)
		expected := append([]btcjson.ListUnspentResult{utxos[4]}, utxos[0:2]...)
		expected = append(expected, utxos[6])
		expected = append(expected, utxos[5])
		expected = append(expected, utxos[3])
		require.Equal(t, expected, selected.UTXOs)
		require.Equal(t, uint16(3), selected.ConsolidatedUTXOs)
		require.Equal(t, 2.0, selected.ConsolidatedValue)
	})

	t.Run("should consolidate all utxos sparse", func(t *testing.T) {
		// create test suite and set nonce-mark utxo for nonce 24105431
		ob, utxos := createTestSuitWithUTXOsAndNonceMark(t, 24105431, dummyTxID)

		// input: utxoCap = 12, amount = 0.13, nonce = 24105432, rank = 1
		// output: [0.24107431, 0.01, 0.12, 8.72, 5.16, 3.28, 2.97, 1.26, 0.5, 0.24, 0.18], 22.68107431
		selected, err := ob.SelectUTXOs(ctx, 0.13, 12, 24105432, 1)
		require.NoError(t, err)
		require.InEpsilon(t, 22.68107431, selected.Value, 1e-8)
		expected := append([]btcjson.ListUnspentResult{utxos[4]}, utxos[0:2]...)
		for i := 10; i >= 5; i-- { // append consolidated utxos in descending order
			expected = append(expected, utxos[i])
		}
		expected = append(expected, utxos[3])
		expected = append(expected, utxos[2])
		require.Equal(t, expected, selected.UTXOs)
		require.Equal(t, uint16(8), selected.ConsolidatedUTXOs)
		require.Equal(t, 22.31, selected.ConsolidatedValue)
	})
}

// helper function to create a test suite with UTXOs
func newTestSuitWithUTXOs(t *testing.T) (*testSuite, []btcjson.ListUnspentResult) {
	// create test observer
	ob := newTestSuite(t, chains.BitcoinMainnet)

	// get test UTXOs
	tssAddress, err := ob.TSS().PubKey().AddressBTC(ob.Chain().ChainId)
	require.NoError(t, err)
	utxos := getTestUTXOs(tssAddress.EncodeAddress())

	// mock up pending nonces and UTXOs
	pendingNonces := observertypes.PendingNonces{}
	ob.zetacore.On("GetPendingNoncesByChain", mock.Anything, mock.Anything).Maybe().Return(pendingNonces, nil)
	ob.client.On("ListUnspentMinMaxAddresses", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Maybe().
		Return(utxos, nil)

	// update UTXOs
	err = ob.FetchUTXOs(context.Background())
	require.NoError(t, err)

	return ob, utxos
}

// helper function to create a test suite with UTXOs and nonce mark
func createTestSuitWithUTXOsAndNonceMark(
	t *testing.T,
	nonce uint64,
	txid string,
) (*testSuite, []btcjson.ListUnspentResult) {
	// create test observer
	ob := newTestSuite(t, chains.BitcoinMainnet)

	// make a nonce mark UTXO
	tssAddress, err := ob.TSS().PubKey().AddressBTC(ob.Chain().ChainId)
	require.NoError(t, err)
	nonceMark := btcjson.ListUnspentResult{
		TxID:          txid,
		Address:       tssAddress.EncodeAddress(),
		Amount:        float64(chains.NonceMarkAmount(nonce)) * 1e-8,
		Confirmations: 1,
	}

	// get test UTXOs and append nonce-mark UTXO
	utxos := getTestUTXOs(tssAddress.EncodeAddress())
	utxos = append(utxos, nonceMark)

	// mock up pending nonces and UTXOs
	pendingNonces := observertypes.PendingNonces{}
	ob.zetacore.On("GetPendingNoncesByChain", mock.Anything, mock.Anything).Maybe().Return(pendingNonces, nil)
	ob.client.On("ListUnspentMinMaxAddresses", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Maybe().
		Return(utxos, nil)

	// update UTXOs
	err = ob.FetchUTXOs(context.Background())
	require.NoError(t, err)

	// set nonce-mark
	ob.Observer.SetIncludedTx(nonce, &btcjson.GetTransactionResult{TxID: txid})

	return ob, utxos
}

// getTestUTXOs returns a list of constant UTXOs for testing
func getTestUTXOs(owner string) []btcjson.ListUnspentResult {
	// create 10 constant dummy UTXOs (22.44 BTC in total)
	utxos := make([]btcjson.ListUnspentResult, 0, 10)
	amounts := []float64{0.01, 0.12, 0.18, 0.24, 0.5, 1.26, 2.97, 3.28, 5.16, 8.72}
	for _, amount := range amounts {
		utxos = append(utxos, btcjson.ListUnspentResult{
			Address:       owner,
			Amount:        amount,
			Confirmations: 1,
		})
	}

	// shuffle the UTXOs, zetaclient will always sort them
	rand.Seed(uint64(time.Now().Second()))
	rand.Shuffle(len(utxos), func(i, j int) {
		utxos[i], utxos[j] = utxos[j], utxos[i]
	})

	return utxos
}
