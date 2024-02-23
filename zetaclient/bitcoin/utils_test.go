package bitcoin

import (
	"path"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

func TestDecodeP2WPKHVout(t *testing.T) {
	// load archived outtx raw result
	// https://blockstream.info/tx/030cd813443f7b70cc6d8a544d320c6d8465e4528fc0f3410b599dc0b26753a0
	var rawResult btcjson.TxRawResult
	err := testutils.LoadObjectFromJSONFile(&rawResult, path.Join("../", testutils.TestDataPathBTC, "outtx_8332_148_raw_result.json"))
	require.NoError(t, err)
	require.Len(t, rawResult.Vout, 3)

	// it's a mainnet outtx
	chain := common.BtcMainnetChain()
	nonce := uint64(148)

	// decode vout 0, nonce mark 148
	receiver, amount, err := DecodeP2WPKHVout(rawResult.Vout[0], chain)
	require.NoError(t, err)
	require.Equal(t, testutils.TSSAddressBTCMainnet, receiver)
	require.Equal(t, common.NonceMarkAmount(nonce), amount)

	// decode vout 1, payment 0.00012000 BTC
	receiver, amount, err = DecodeP2WPKHVout(rawResult.Vout[1], chain)
	require.NoError(t, err)
	require.Equal(t, "bc1qpsdlklfcmlcfgm77c43x65ddtrt7n0z57hsyjp", receiver)
	require.Equal(t, int64(12000), amount)

	// decode vout 2, change 0.39041489 BTC
	receiver, amount, err = DecodeP2WPKHVout(rawResult.Vout[2], chain)
	require.NoError(t, err)
	require.Equal(t, testutils.TSSAddressBTCMainnet, receiver)
	require.Equal(t, int64(39041489), amount)
}

func TestDecodeP2WPKHVoutErrors(t *testing.T) {
	// load archived outtx raw result
	// https://blockstream.info/tx/030cd813443f7b70cc6d8a544d320c6d8465e4528fc0f3410b599dc0b26753a0
	var rawResult btcjson.TxRawResult
	err := testutils.LoadObjectFromJSONFile(&rawResult, path.Join("../", testutils.TestDataPathBTC, "outtx_8332_148_raw_result.json"))
	require.NoError(t, err)

	chain := common.BtcMainnetChain()

	t.Run("should return error on invalid amount", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.Value = -0.5 // negative amount, should not happen
		_, _, err := DecodeP2WPKHVout(invalidVout, chain)
		require.Error(t, err)
		require.ErrorContains(t, err, "error getting satoshis")
	})
	t.Run("should return error on invalid script", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Hex = "invalid script"
		_, _, err := DecodeP2WPKHVout(invalidVout, chain)
		require.Error(t, err)
		require.ErrorContains(t, err, "error decoding scriptPubKey")
	})
	t.Run("should return error on unsupported script", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// can use any invalid script, https://blockstream.info/tx/e95c6ff206103716129c8e3aa8def1427782af3490589d1ea35ccf0122adbc25 (P2SH)
		invalidVout.ScriptPubKey.Hex = "a91413b2388e6532653a4b369b7e4ed130f7b81626cc87"
		_, _, err := DecodeP2WPKHVout(invalidVout, chain)
		require.Error(t, err)
		require.ErrorContains(t, err, "unsupported scriptPubKey")
	})
	t.Run("should return error on unsupported witness version", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// use a fake witness version 1, even if version 0 is the only witness version defined in BIP141
		invalidVout.ScriptPubKey.Hex = "01140c1bfb7d38dff0946fdec5626d51ad58d7e9bc54"
		_, _, err := DecodeP2WPKHVout(invalidVout, chain)
		require.Error(t, err)
		require.ErrorContains(t, err, "unsupported witness in scriptPubKey")
	})
}
