package observer

import (
	"context"
	"math"
	"sort"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
)

// the relative path to the testdata directory
var TestDataDir = "../../../"

// MockBTCObserverMainnet creates a mock Bitcoin mainnet observer for testing
func MockBTCObserverMainnet(t *testing.T) *Observer {
	// setup mock arguments
	chain := chains.BitcoinMainnet
	btcClient := mocks.NewMockBTCRPCClient().WithBlockCount(100)
	params := mocks.MockChainParams(chain.ChainId, 10)
	tss := mocks.NewTSSMainnet()

	// create Bitcoin observer
	ob, err := NewObserver(chain, btcClient, params, nil, tss, testutils.SQLiteMemory, base.Logger{}, nil)
	require.NoError(t, err)

	return ob
}

// helper function to create a test Bitcoin observer
func createObserverWithPrivateKey(t *testing.T) *Observer {
	skHex := "7b8507ba117e069f4a3f456f505276084f8c92aee86ac78ae37b4d1801d35fa8"
	privateKey, err := crypto.HexToECDSA(skHex)
	require.NoError(t, err)
	tss := &mocks.TSS{
		PrivKey: privateKey,
	}

	// create Bitcoin observer with mock tss
	ob := MockBTCObserverMainnet(t)
	ob.WithTSS(tss)

	return ob
}

// helper function to create a test Bitcoin observer with UTXOs
func createObserverWithUTXOs(t *testing.T) *Observer {
	// Create Bitcoin observer
	ob := createObserverWithPrivateKey(t)
	tssAddress := ob.TSS().BTCAddressWitnessPubkeyHash().EncodeAddress()

	// Create 10 dummy UTXOs (22.44 BTC in total)
	ob.utxos = make([]btcjson.ListUnspentResult, 0, 10)
	amounts := []float64{0.01, 0.12, 0.18, 0.24, 0.5, 1.26, 2.97, 3.28, 5.16, 8.72}
	for _, amount := range amounts {
		ob.utxos = append(ob.utxos, btcjson.ListUnspentResult{Address: tssAddress, Amount: amount})
	}
	return ob
}

func mineTxNSetNonceMark(ob *Observer, nonce uint64, txid string, preMarkIndex int) {
	// Mine transaction
	outboundID := ob.GetTxID(nonce)
	ob.includedTxResults[outboundID] = &btcjson.GetTransactionResult{TxID: txid}

	// Set nonce mark
	tssAddress := ob.TSS().BTCAddressWitnessPubkeyHash().EncodeAddress()
	nonceMark := btcjson.ListUnspentResult{
		TxID:    txid,
		Address: tssAddress,
		Amount:  float64(chains.NonceMarkAmount(nonce)) * 1e-8,
	}
	if preMarkIndex >= 0 { // replace nonce-mark utxo
		ob.utxos[preMarkIndex] = nonceMark

	} else { // add nonce-mark utxo directly
		ob.utxos = append(ob.utxos, nonceMark)
	}
	sort.SliceStable(ob.utxos, func(i, j int) bool {
		return ob.utxos[i].Amount < ob.utxos[j].Amount
	})
}

func TestCheckTSSVout(t *testing.T) {
	// the archived outbound raw result file and cctx file
	// https://blockstream.info/tx/030cd813443f7b70cc6d8a544d320c6d8465e4528fc0f3410b599dc0b26753a0
	chain := chains.BitcoinMainnet
	chainID := chain.ChainId
	nonce := uint64(148)

	// create mainnet mock client
	ob := MockBTCObserverMainnet(t)

	t.Run("valid TSS vout should pass", func(t *testing.T) {
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()
		err := ob.checkTSSVout(params, rawResult.Vout)
		require.NoError(t, err)
	})
	t.Run("should fail if vout length < 2 or > 3", func(t *testing.T) {
		_, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()

		err := ob.checkTSSVout(params, []btcjson.Vout{{}})
		require.ErrorContains(t, err, "invalid number of vouts")

		err = ob.checkTSSVout(params, []btcjson.Vout{{}, {}, {}, {}})
		require.ErrorContains(t, err, "invalid number of vouts")
	})
	t.Run("should fail on invalid TSS vout", func(t *testing.T) {
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()

		// invalid TSS vout
		rawResult.Vout[0].ScriptPubKey.Hex = "invalid script"
		err := ob.checkTSSVout(params, rawResult.Vout)
		require.Error(t, err)
	})
	t.Run("should fail if vout 0 is not to the TSS address", func(t *testing.T) {
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()

		// not TSS address, bc1qh297vdt8xq6df5xae9z8gzd4jsu9a392mp0dus
		rawResult.Vout[0].ScriptPubKey.Hex = "0014ba8be635673034d4d0ddc9447409b594385ec4aa"
		err := ob.checkTSSVout(params, rawResult.Vout)
		require.ErrorContains(t, err, "not match TSS address")
	})
	t.Run("should fail if vout 0 not match nonce mark", func(t *testing.T) {
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()

		// not match nonce mark
		rawResult.Vout[0].Value = 0.00000147
		err := ob.checkTSSVout(params, rawResult.Vout)
		require.ErrorContains(t, err, "not match nonce-mark amount")
	})
	t.Run("should fail if vout 1 is not to the receiver address", func(t *testing.T) {
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()

		// not receiver address, bc1qh297vdt8xq6df5xae9z8gzd4jsu9a392mp0dus
		rawResult.Vout[1].ScriptPubKey.Hex = "0014ba8be635673034d4d0ddc9447409b594385ec4aa"
		err := ob.checkTSSVout(params, rawResult.Vout)
		require.ErrorContains(t, err, "not match params receiver")
	})
	t.Run("should fail if vout 1 not match payment amount", func(t *testing.T) {
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()

		// not match payment amount
		rawResult.Vout[1].Value = 0.00011000
		err := ob.checkTSSVout(params, rawResult.Vout)
		require.ErrorContains(t, err, "not match params amount")
	})
	t.Run("should fail if vout 2 is not to the TSS address", func(t *testing.T) {
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()

		// not TSS address, bc1qh297vdt8xq6df5xae9z8gzd4jsu9a392mp0dus
		rawResult.Vout[2].ScriptPubKey.Hex = "0014ba8be635673034d4d0ddc9447409b594385ec4aa"
		err := ob.checkTSSVout(params, rawResult.Vout)
		require.ErrorContains(t, err, "not match TSS address")
	})
}

func TestCheckTSSVoutCancelled(t *testing.T) {
	// the archived outbound raw result file and cctx file
	// https://blockstream.info/tx/030cd813443f7b70cc6d8a544d320c6d8465e4528fc0f3410b599dc0b26753a0
	chain := chains.BitcoinMainnet
	chainID := chain.ChainId
	nonce := uint64(148)

	// create mainnet mock client
	ob := MockBTCObserverMainnet(t)

	t.Run("valid TSS vout should pass", func(t *testing.T) {
		// remove change vout to simulate cancelled tx
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		rawResult.Vout[1] = rawResult.Vout[2]
		rawResult.Vout = rawResult.Vout[:2]
		params := cctx.GetCurrentOutboundParam()

		err := ob.checkTSSVoutCancelled(params, rawResult.Vout)
		require.NoError(t, err)
	})
	t.Run("should fail if vout length < 1 or > 2", func(t *testing.T) {
		_, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		params := cctx.GetCurrentOutboundParam()

		err := ob.checkTSSVoutCancelled(params, []btcjson.Vout{})
		require.ErrorContains(t, err, "invalid number of vouts")

		err = ob.checkTSSVoutCancelled(params, []btcjson.Vout{{}, {}, {}})
		require.ErrorContains(t, err, "invalid number of vouts")
	})
	t.Run("should fail if vout 0 is not to the TSS address", func(t *testing.T) {
		// remove change vout to simulate cancelled tx
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		rawResult.Vout[1] = rawResult.Vout[2]
		rawResult.Vout = rawResult.Vout[:2]
		params := cctx.GetCurrentOutboundParam()

		// not TSS address, bc1qh297vdt8xq6df5xae9z8gzd4jsu9a392mp0dus
		rawResult.Vout[0].ScriptPubKey.Hex = "0014ba8be635673034d4d0ddc9447409b594385ec4aa"
		err := ob.checkTSSVoutCancelled(params, rawResult.Vout)
		require.ErrorContains(t, err, "not match TSS address")
	})
	t.Run("should fail if vout 0 not match nonce mark", func(t *testing.T) {
		// remove change vout to simulate cancelled tx
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		rawResult.Vout[1] = rawResult.Vout[2]
		rawResult.Vout = rawResult.Vout[:2]
		params := cctx.GetCurrentOutboundParam()

		// not match nonce mark
		rawResult.Vout[0].Value = 0.00000147
		err := ob.checkTSSVoutCancelled(params, rawResult.Vout)
		require.ErrorContains(t, err, "not match nonce-mark amount")
	})
	t.Run("should fail if vout 1 is not to the TSS address", func(t *testing.T) {
		// remove change vout to simulate cancelled tx
		rawResult, cctx := testutils.LoadBTCTxRawResultNCctx(t, TestDataDir, chainID, nonce)
		rawResult.Vout[1] = rawResult.Vout[2]
		rawResult.Vout[1].N = 1 // swap vout index
		rawResult.Vout = rawResult.Vout[:2]
		params := cctx.GetCurrentOutboundParam()

		// not TSS address, bc1qh297vdt8xq6df5xae9z8gzd4jsu9a392mp0dus
		rawResult.Vout[1].ScriptPubKey.Hex = "0014ba8be635673034d4d0ddc9447409b594385ec4aa"
		err := ob.checkTSSVoutCancelled(params, rawResult.Vout)
		require.ErrorContains(t, err, "not match TSS address")
	})
}

func TestSelectUTXOs(t *testing.T) {
	ctx := context.Background()

	ob := createObserverWithUTXOs(t)
	dummyTxID := "6e6f71d281146c1fc5c755b35908ee449f26786c84e2ae18f98b268de40b7ec4"

	// Case1: nonce = 0, bootstrap
	// 		input: utxoCap = 5, amount = 0.01, nonce = 0
	// 		output: [0.01], 0.01
	result, amount, _, _, err := ob.SelectUTXOs(ctx, 0.01, 5, 0, math.MaxUint16, true)
	require.NoError(t, err)
	require.Equal(t, 0.01, amount)
	require.Equal(t, ob.utxos[0:1], result)

	// Case2: nonce = 1, must FAIL and wait for previous transaction to be mined
	// 		input: utxoCap = 5, amount = 0.5, nonce = 1
	// 		output: error
	result, amount, _, _, err = ob.SelectUTXOs(ctx, 0.5, 5, 1, math.MaxUint16, true)
	require.Error(t, err)
	require.Nil(t, result)
	require.Zero(t, amount)
	require.Equal(t, "getOutboundIDByNonce: cannot find outbound txid for nonce 0", err.Error())
	mineTxNSetNonceMark(ob, 0, dummyTxID, -1) // mine a transaction and set nonce-mark utxo for nonce 0

	// Case3: nonce = 1, should pass now
	// 		input: utxoCap = 5, amount = 0.5, nonce = 1
	// 		output: [0.00002, 0.01, 0.12, 0.18, 0.24], 0.55002
	result, amount, _, _, err = ob.SelectUTXOs(ctx, 0.5, 5, 1, math.MaxUint16, true)
	require.NoError(t, err)
	require.Equal(t, 0.55002, amount)
	require.Equal(t, ob.utxos[0:5], result)
	mineTxNSetNonceMark(ob, 1, dummyTxID, 0) // mine a transaction and set nonce-mark utxo for nonce 1

	// Case4:
	// 		input: utxoCap = 5, amount = 1.0, nonce = 2
	// 		output: [0.00002001, 0.01, 0.12, 0.18, 0.24, 0.5], 1.05002001
	result, amount, _, _, err = ob.SelectUTXOs(ctx, 1.0, 5, 2, math.MaxUint16, true)
	require.NoError(t, err)
	require.InEpsilon(t, 1.05002001, amount, 1e-8)
	require.Equal(t, ob.utxos[0:6], result)
	mineTxNSetNonceMark(ob, 2, dummyTxID, 0) // mine a transaction and set nonce-mark utxo for nonce 2

	// Case5: should include nonce-mark utxo on the LEFT
	// 		input: utxoCap = 5, amount = 8.05, nonce = 3
	// 		output: [0.00002002, 0.24, 0.5, 1.26, 2.97, 3.28], 8.25002002
	result, amount, _, _, err = ob.SelectUTXOs(ctx, 8.05, 5, 3, math.MaxUint16, true)
	require.NoError(t, err)
	require.InEpsilon(t, 8.25002002, amount, 1e-8)
	expected := append([]btcjson.ListUnspentResult{ob.utxos[0]}, ob.utxos[4:9]...)
	require.Equal(t, expected, result)
	mineTxNSetNonceMark(ob, 24105431, dummyTxID, 0) // mine a transaction and set nonce-mark utxo for nonce 24105431

	// Case6: should include nonce-mark utxo on the RIGHT
	// 		input: utxoCap = 5, amount = 0.503, nonce = 24105432
	// 		output: [0.24107432, 0.01, 0.12, 0.18, 0.24], 0.55002002
	result, amount, _, _, err = ob.SelectUTXOs(ctx, 0.503, 5, 24105432, math.MaxUint16, true)
	require.NoError(t, err)
	require.InEpsilon(t, 0.79107431, amount, 1e-8)
	expected = append([]btcjson.ListUnspentResult{ob.utxos[4]}, ob.utxos[0:4]...)
	require.Equal(t, expected, result)
	mineTxNSetNonceMark(ob, 24105432, dummyTxID, 4) // mine a transaction and set nonce-mark utxo for nonce 24105432

	// Case7: should include nonce-mark utxo in the MIDDLE
	// 		input: utxoCap = 5, amount = 1.0, nonce = 24105433
	// 		output: [0.24107432, 0.12, 0.18, 0.24, 0.5], 1.28107432
	result, amount, _, _, err = ob.SelectUTXOs(ctx, 1.0, 5, 24105433, math.MaxUint16, true)
	require.NoError(t, err)
	require.InEpsilon(t, 1.28107432, amount, 1e-8)
	expected = append([]btcjson.ListUnspentResult{ob.utxos[4]}, ob.utxos[1:4]...)
	expected = append(expected, ob.utxos[5])
	require.Equal(t, expected, result)

	// Case8: should work with maximum amount
	// 		input: utxoCap = 5, amount = 16.03
	// 		output: [0.24107432, 1.26, 2.97, 3.28, 5.16, 8.72], 21.63107432
	result, amount, _, _, err = ob.SelectUTXOs(ctx, 16.03, 5, 24105433, math.MaxUint16, true)
	require.NoError(t, err)
	require.InEpsilon(t, 21.63107432, amount, 1e-8)
	expected = append([]btcjson.ListUnspentResult{ob.utxos[4]}, ob.utxos[6:11]...)
	require.Equal(t, expected, result)

	// Case9: must FAIL due to insufficient funds
	// 		input: utxoCap = 5, amount = 21.64
	// 		output: error
	result, amount, _, _, err = ob.SelectUTXOs(ctx, 21.64, 5, 24105433, math.MaxUint16, true)
	require.Error(t, err)
	require.Nil(t, result)
	require.Zero(t, amount)
	require.Equal(
		t,
		"SelectUTXOs: not enough btc in reserve - available : 21.63107432 , tx amount : 21.64",
		err.Error(),
	)
}

func TestUTXOConsolidation(t *testing.T) {
	ctx := context.Background()

	dummyTxID := "6e6f71d281146c1fc5c755b35908ee449f26786c84e2ae18f98b268de40b7ec4"

	t.Run("should not consolidate", func(t *testing.T) {
		ob := createObserverWithUTXOs(t)
		mineTxNSetNonceMark(ob, 0, dummyTxID, -1) // mine a transaction and set nonce-mark utxo for nonce 0

		// input: utxoCap = 10, amount = 0.01, nonce = 1, rank = 10
		// output: [0.00002, 0.01], 0.01002
		result, amount, clsdtUtxo, clsdtValue, err := ob.SelectUTXOs(ctx, 0.01, 10, 1, 10, true)
		require.NoError(t, err)
		require.Equal(t, 0.01002, amount)
		require.Equal(t, ob.utxos[0:2], result)
		require.Equal(t, uint16(0), clsdtUtxo)
		require.Equal(t, 0.0, clsdtValue)
	})

	t.Run("should consolidate 1 utxo", func(t *testing.T) {
		ob := createObserverWithUTXOs(t)
		mineTxNSetNonceMark(ob, 0, dummyTxID, -1) // mine a transaction and set nonce-mark utxo for nonce 0

		// input: utxoCap = 9, amount = 0.01, nonce = 1, rank = 9
		// output: [0.00002, 0.01, 0.12], 0.13002
		result, amount, clsdtUtxo, clsdtValue, err := ob.SelectUTXOs(ctx, 0.01, 9, 1, 9, true)
		require.NoError(t, err)
		require.Equal(t, 0.13002, amount)
		require.Equal(t, ob.utxos[0:3], result)
		require.Equal(t, uint16(1), clsdtUtxo)
		require.Equal(t, 0.12, clsdtValue)
	})

	t.Run("should consolidate 3 utxos", func(t *testing.T) {
		ob := createObserverWithUTXOs(t)
		mineTxNSetNonceMark(ob, 0, dummyTxID, -1) // mine a transaction and set nonce-mark utxo for nonce 0

		// input: utxoCap = 5, amount = 0.01, nonce = 0, rank = 5
		// output: [0.00002, 0.014, 1.26, 0.5, 0.2], 2.01002
		result, amount, clsdtUtxo, clsdtValue, err := ob.SelectUTXOs(ctx, 0.01, 5, 1, 5, true)
		require.NoError(t, err)
		require.Equal(t, 2.01002, amount)
		expected := make([]btcjson.ListUnspentResult, 2)
		copy(expected, ob.utxos[0:2])
		for i := 6; i >= 4; i-- { // append consolidated utxos in descending order
			expected = append(expected, ob.utxos[i])
		}
		require.Equal(t, expected, result)
		require.Equal(t, uint16(3), clsdtUtxo)
		require.Equal(t, 2.0, clsdtValue)
	})

	t.Run("should consolidate all utxos using rank 1", func(t *testing.T) {
		ob := createObserverWithUTXOs(t)
		mineTxNSetNonceMark(ob, 0, dummyTxID, -1) // mine a transaction and set nonce-mark utxo for nonce 0

		// input: utxoCap = 12, amount = 0.01, nonce = 0, rank = 1
		// output: [0.00002, 0.01, 8.72, 5.16, 3.28, 2.97, 1.26, 0.5, 0.24, 0.18, 0.12], 22.44002
		result, amount, clsdtUtxo, clsdtValue, err := ob.SelectUTXOs(ctx, 0.01, 12, 1, 1, true)
		require.NoError(t, err)
		require.Equal(t, 22.44002, amount)
		expected := make([]btcjson.ListUnspentResult, 2)
		copy(expected, ob.utxos[0:2])
		for i := 10; i >= 2; i-- { // append consolidated utxos in descending order
			expected = append(expected, ob.utxos[i])
		}
		require.Equal(t, expected, result)
		require.Equal(t, uint16(9), clsdtUtxo)
		require.Equal(t, 22.43, clsdtValue)
	})

	t.Run("should consolidate 3 utxos sparse", func(t *testing.T) {
		ob := createObserverWithUTXOs(t)
		mineTxNSetNonceMark(
			ob,
			24105431,
			dummyTxID,
			-1,
		) // mine a transaction and set nonce-mark utxo for nonce 24105431

		// input: utxoCap = 5, amount = 0.13, nonce = 24105432, rank = 5
		// output: [0.24107431, 0.01, 0.12, 1.26, 0.5, 0.24], 2.37107431
		result, amount, clsdtUtxo, clsdtValue, err := ob.SelectUTXOs(ctx, 0.13, 5, 24105432, 5, true)
		require.NoError(t, err)
		require.InEpsilon(t, 2.37107431, amount, 1e-8)
		expected := append([]btcjson.ListUnspentResult{ob.utxos[4]}, ob.utxos[0:2]...)
		expected = append(expected, ob.utxos[6])
		expected = append(expected, ob.utxos[5])
		expected = append(expected, ob.utxos[3])
		require.Equal(t, expected, result)
		require.Equal(t, uint16(3), clsdtUtxo)
		require.Equal(t, 2.0, clsdtValue)
	})

	t.Run("should consolidate all utxos sparse", func(t *testing.T) {
		ob := createObserverWithUTXOs(t)
		mineTxNSetNonceMark(
			ob,
			24105431,
			dummyTxID,
			-1,
		) // mine a transaction and set nonce-mark utxo for nonce 24105431

		// input: utxoCap = 12, amount = 0.13, nonce = 24105432, rank = 1
		// output: [0.24107431, 0.01, 0.12, 8.72, 5.16, 3.28, 2.97, 1.26, 0.5, 0.24, 0.18], 22.68107431
		result, amount, clsdtUtxo, clsdtValue, err := ob.SelectUTXOs(ctx, 0.13, 12, 24105432, 1, true)
		require.NoError(t, err)
		require.InEpsilon(t, 22.68107431, amount, 1e-8)
		expected := append([]btcjson.ListUnspentResult{ob.utxos[4]}, ob.utxos[0:2]...)
		for i := 10; i >= 5; i-- { // append consolidated utxos in descending order
			expected = append(expected, ob.utxos[i])
		}
		expected = append(expected, ob.utxos[3])
		expected = append(expected, ob.utxos[2])
		require.Equal(t, expected, result)
		require.Equal(t, uint16(8), clsdtUtxo)
		require.Equal(t, 22.31, clsdtValue)
	})
}
