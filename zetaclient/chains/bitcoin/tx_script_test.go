package bitcoin

import (
	"encoding/hex"
	"path"
	"strings"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/constant"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

// the relative path to the testdata directory
var TestDataDir = "../../"

func TestDecodeVoutP2TR(t *testing.T) {
	// load archived tx raw result
	// https://mempool.space/tx/259fc21e63e138136c8f19270a0f7ca10039a66a474f91d23a17896f46e677a7
	chain := chains.BtcMainnetChain
	txHash := "259fc21e63e138136c8f19270a0f7ca10039a66a474f91d23a17896f46e677a7"
	net := &chaincfg.MainNetParams

	rawResult := testutils.LoadBTCTxRawResult(t, TestDataDir, chain.ChainId, "P2TR", txHash)
	require.Len(t, rawResult.Vout, 2)

	// decode vout 0, P2TR
	receiver, err := DecodeScriptP2TR(rawResult.Vout[0].ScriptPubKey.Hex, net)
	require.NoError(t, err)
	require.Equal(t, "bc1p4scddlkkuw9486579autxumxmkvuphm5pz4jvf7f6pdh50p2uzqstawjt9", receiver)
}

func TestDecodeVoutP2TRErrors(t *testing.T) {
	// load archived tx raw result
	// https://mempool.space/tx/259fc21e63e138136c8f19270a0f7ca10039a66a474f91d23a17896f46e677a7
	chain := chains.BtcMainnetChain
	txHash := "259fc21e63e138136c8f19270a0f7ca10039a66a474f91d23a17896f46e677a7"
	net := &chaincfg.MainNetParams
	rawResult := testutils.LoadBTCTxRawResult(t, TestDataDir, chain.ChainId, "P2TR", txHash)

	t.Run("should return error on invalid script", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Hex = "invalid script"
		_, err := DecodeScriptP2TR(invalidVout.ScriptPubKey.Hex, net)
		require.ErrorContains(t, err, "error decoding script")
	})
	t.Run("should return error on wrong script length", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Hex = "0020" // 2 bytes, should be 34
		_, err := DecodeScriptP2TR(invalidVout.ScriptPubKey.Hex, net)
		require.ErrorContains(t, err, "invalid P2TR script")
	})
	t.Run("should return error on invalid OP_1", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// modify the OP_1 '51' to OP_2 '52'
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "51", "52", 1)
		_, err := DecodeScriptP2TR(invalidVout.ScriptPubKey.Hex, net)
		require.ErrorContains(t, err, "invalid P2TR script")
	})
	t.Run("should return error on wrong hash length", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// modify the length '20' to '19'
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "5120", "5119", 1)
		_, err := DecodeScriptP2TR(invalidVout.ScriptPubKey.Hex, net)
		require.ErrorContains(t, err, "invalid P2TR script")
	})
}

func TestDecodeVoutP2WSH(t *testing.T) {
	// load archived tx raw result
	// https://mempool.space/tx/791bb9d16f7ab05f70a116d18eaf3552faf77b9d5688699a480261424b4f7e53
	chain := chains.BtcMainnetChain
	txHash := "791bb9d16f7ab05f70a116d18eaf3552faf77b9d5688699a480261424b4f7e53"
	net := &chaincfg.MainNetParams

	rawResult := testutils.LoadBTCTxRawResult(t, TestDataDir, chain.ChainId, "P2WSH", txHash)
	require.Len(t, rawResult.Vout, 1)

	// decode vout 0, P2WSH
	receiver, err := DecodeScriptP2WSH(rawResult.Vout[0].ScriptPubKey.Hex, net)
	require.NoError(t, err)
	require.Equal(t, "bc1qqv6pwn470vu0tssdfha4zdk89v3c8ch5lsnyy855k9hcrcv3evequdmjmc", receiver)
}

func TestDecodeVoutP2WSHErrors(t *testing.T) {
	// load archived tx raw result
	// https://mempool.space/tx/791bb9d16f7ab05f70a116d18eaf3552faf77b9d5688699a480261424b4f7e53
	chain := chains.BtcMainnetChain
	txHash := "791bb9d16f7ab05f70a116d18eaf3552faf77b9d5688699a480261424b4f7e53"
	net := &chaincfg.MainNetParams
	rawResult := testutils.LoadBTCTxRawResult(t, TestDataDir, chain.ChainId, "P2WSH", txHash)

	t.Run("should return error on invalid script", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Hex = "invalid script"
		_, err := DecodeScriptP2WSH(invalidVout.ScriptPubKey.Hex, net)
		require.ErrorContains(t, err, "error decoding script")
	})
	t.Run("should return error on wrong script length", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Hex = "0020" // 2 bytes, should be 34
		_, err := DecodeScriptP2WSH(invalidVout.ScriptPubKey.Hex, net)
		require.ErrorContains(t, err, "invalid P2WSH script")
	})
	t.Run("should return error on invalid OP_0", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// modify the OP_0 '00' to OP_1 '51'
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "00", "51", 1)
		_, err := DecodeScriptP2WSH(invalidVout.ScriptPubKey.Hex, net)
		require.ErrorContains(t, err, "invalid P2WSH script")
	})
	t.Run("should return error on wrong hash length", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// modify the length '20' to '19'
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "0020", "0019", 1)
		_, err := DecodeScriptP2WSH(invalidVout.ScriptPubKey.Hex, net)
		require.ErrorContains(t, err, "invalid P2WSH script")
	})
}

func TestDecodeP2WPKHVout(t *testing.T) {
	// load archived outtx raw result
	// https://mempool.space/tx/030cd813443f7b70cc6d8a544d320c6d8465e4528fc0f3410b599dc0b26753a0
	chain := chains.BtcMainnetChain
	nonce := uint64(148)
	net := &chaincfg.MainNetParams
	nameTx := path.Join(TestDataDir, testutils.TestDataPathBTC, testutils.FileNameBTCOuttx(chain.ChainId, nonce))

	var rawResult btcjson.TxRawResult
	testutils.LoadObjectFromJSONFile(t, &rawResult, nameTx)
	require.Len(t, rawResult.Vout, 3)

	// decode vout 0, nonce mark 148
	receiver, err := DecodeScriptP2WPKH(rawResult.Vout[0].ScriptPubKey.Hex, net)
	require.NoError(t, err)
	require.Equal(t, testutils.TSSAddressBTCMainnet, receiver)

	// decode vout 1, payment 0.00012000 BTC
	receiver, err = DecodeScriptP2WPKH(rawResult.Vout[1].ScriptPubKey.Hex, net)
	require.NoError(t, err)
	require.Equal(t, "bc1qpsdlklfcmlcfgm77c43x65ddtrt7n0z57hsyjp", receiver)

	// decode vout 2, change 0.39041489 BTC
	receiver, err = DecodeScriptP2WPKH(rawResult.Vout[2].ScriptPubKey.Hex, net)
	require.NoError(t, err)
	require.Equal(t, testutils.TSSAddressBTCMainnet, receiver)
}

func TestDecodeP2WPKHVoutErrors(t *testing.T) {
	// load archived outtx raw result
	// https://mempool.space/tx/030cd813443f7b70cc6d8a544d320c6d8465e4528fc0f3410b599dc0b26753a0
	chain := chains.BtcMainnetChain
	nonce := uint64(148)
	net := &chaincfg.MainNetParams
	nameTx := path.Join(TestDataDir, testutils.TestDataPathBTC, testutils.FileNameBTCOuttx(chain.ChainId, nonce))

	var rawResult btcjson.TxRawResult
	testutils.LoadObjectFromJSONFile(t, &rawResult, nameTx)

	t.Run("should return error on invalid script", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Hex = "invalid script"
		_, err := DecodeScriptP2WPKH(invalidVout.ScriptPubKey.Hex, net)
		require.ErrorContains(t, err, "error decoding script")
	})
	t.Run("should return error on wrong script length", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Hex = "0014" // 2 bytes, should be 22
		_, err := DecodeScriptP2WPKH(invalidVout.ScriptPubKey.Hex, net)
		require.ErrorContains(t, err, "invalid P2WPKH script")
	})
	t.Run("should return error on wrong hash length", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// modify the length '14' to '13'
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "0014", "0013", 1)
		_, err := DecodeScriptP2WPKH(invalidVout.ScriptPubKey.Hex, net)
		require.ErrorContains(t, err, "invalid P2WPKH script")
	})
}

func TestDecodeVoutP2SH(t *testing.T) {
	// load archived tx raw result
	// https://mempool.space/tx/fd68c8b4478686ca6f5ae4c28eaab055490650dbdaa6c2c8e380a7e075958a21
	chain := chains.BtcMainnetChain
	txHash := "fd68c8b4478686ca6f5ae4c28eaab055490650dbdaa6c2c8e380a7e075958a21"
	net := &chaincfg.MainNetParams

	rawResult := testutils.LoadBTCTxRawResult(t, TestDataDir, chain.ChainId, "P2SH", txHash)
	require.Len(t, rawResult.Vout, 2)

	// decode vout 0, P2SH
	receiver, err := DecodeScriptP2SH(rawResult.Vout[0].ScriptPubKey.Hex, net)
	require.NoError(t, err)
	require.Equal(t, "327z4GyFM8Y8DiYfasGKQWhRK4MvyMSEgE", receiver)
}

func TestDecodeVoutP2SHErrors(t *testing.T) {
	// load archived tx raw result
	// https://mempool.space/tx/fd68c8b4478686ca6f5ae4c28eaab055490650dbdaa6c2c8e380a7e075958a21
	chain := chains.BtcMainnetChain
	txHash := "fd68c8b4478686ca6f5ae4c28eaab055490650dbdaa6c2c8e380a7e075958a21"
	net := &chaincfg.MainNetParams
	rawResult := testutils.LoadBTCTxRawResult(t, TestDataDir, chain.ChainId, "P2SH", txHash)

	t.Run("should return error on invalid script", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Hex = "invalid script"
		_, err := DecodeScriptP2SH(invalidVout.ScriptPubKey.Hex, net)
		require.ErrorContains(t, err, "error decoding script")
	})
	t.Run("should return error on wrong script length", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Hex = "0014" // 2 bytes, should be 23
		_, err := DecodeScriptP2SH(invalidVout.ScriptPubKey.Hex, net)
		require.ErrorContains(t, err, "invalid P2SH script")
	})
	t.Run("should return error on invalid OP_HASH160", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// modify the OP_HASH160 'a9' to OP_HASH256 'aa'
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "a9", "aa", 1)
		_, err := DecodeScriptP2SH(invalidVout.ScriptPubKey.Hex, net)
		require.ErrorContains(t, err, "invalid P2SH script")
	})
	t.Run("should return error on wrong data length", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// modify the length '14' to '13'
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "a914", "a913", 1)
		_, err := DecodeScriptP2SH(invalidVout.ScriptPubKey.Hex, net)
		require.ErrorContains(t, err, "invalid P2SH script")
	})
	t.Run("should return error on invalid OP_EQUAL", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "87", "88", 1)
		_, err := DecodeScriptP2SH(invalidVout.ScriptPubKey.Hex, net)
		require.ErrorContains(t, err, "invalid P2SH script")
	})
}

func TestDecodeVoutP2PKH(t *testing.T) {
	// load archived tx raw result
	// https://mempool.space/tx/9c741de6e17382b7a9113fc811e3558981a35a360e3d1262a6675892c91322ca
	chain := chains.BtcMainnetChain
	txHash := "9c741de6e17382b7a9113fc811e3558981a35a360e3d1262a6675892c91322ca"
	net := &chaincfg.MainNetParams

	rawResult := testutils.LoadBTCTxRawResult(t, TestDataDir, chain.ChainId, "P2PKH", txHash)
	require.Len(t, rawResult.Vout, 2)

	// decode vout 0, P2PKH
	receiver, err := DecodeScriptP2PKH(rawResult.Vout[0].ScriptPubKey.Hex, net)
	require.NoError(t, err)
	require.Equal(t, "1FueivsE338W2LgifJ25HhTcVJ7CRT8kte", receiver)
}

func TestDecodeVoutP2PKHErrors(t *testing.T) {
	// load archived tx raw result
	// https://mempool.space/tx/9c741de6e17382b7a9113fc811e3558981a35a360e3d1262a6675892c91322ca
	chain := chains.BtcMainnetChain
	txHash := "9c741de6e17382b7a9113fc811e3558981a35a360e3d1262a6675892c91322ca"
	net := &chaincfg.MainNetParams
	rawResult := testutils.LoadBTCTxRawResult(t, TestDataDir, chain.ChainId, "P2PKH", txHash)

	t.Run("should return error on invalid script", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Hex = "invalid script"
		_, err := DecodeScriptP2PKH(invalidVout.ScriptPubKey.Hex, net)
		require.ErrorContains(t, err, "error decoding script")
	})
	t.Run("should return error on wrong script length", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Hex = "76a914" // 3 bytes, should be 25
		_, err := DecodeScriptP2PKH(invalidVout.ScriptPubKey.Hex, net)
		require.ErrorContains(t, err, "invalid P2PKH script")
	})
	t.Run("should return error on invalid OP_DUP", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// modify the OP_DUP '76' to OP_NIP '77'
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "76", "77", 1)
		_, err := DecodeScriptP2PKH(invalidVout.ScriptPubKey.Hex, net)
		require.ErrorContains(t, err, "invalid P2PKH script")
	})
	t.Run("should return error on invalid OP_HASH160", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// modify the OP_HASH160 'a9' to OP_HASH256 'aa'
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "76a9", "76aa", 1)
		_, err := DecodeScriptP2PKH(invalidVout.ScriptPubKey.Hex, net)
		require.ErrorContains(t, err, "invalid P2PKH script")
	})
	t.Run("should return error on wrong data length", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// modify the length '14' to '13'
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "76a914", "76a913", 1)
		_, err := DecodeScriptP2PKH(invalidVout.ScriptPubKey.Hex, net)
		require.ErrorContains(t, err, "invalid P2PKH script")
	})
	t.Run("should return error on invalid OP_EQUALVERIFY", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// modify the OP_EQUALVERIFY '88' to OP_RESERVED1 '89'
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "88ac", "89ac", 1)
		_, err := DecodeScriptP2PKH(invalidVout.ScriptPubKey.Hex, net)
		require.ErrorContains(t, err, "invalid P2PKH script")
	})
	t.Run("should return error on invalid OP_CHECKSIG", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// modify the OP_CHECKSIG 'ac' to OP_CHECKSIGVERIFY 'ad'
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "88ac", "88ad", 1)
		_, err := DecodeScriptP2PKH(invalidVout.ScriptPubKey.Hex, net)
		require.ErrorContains(t, err, "invalid P2PKH script")
	})
}

func TestDecodeOpReturnMemo(t *testing.T) {
	// load archived intx raw result
	// https://mempool.space/tx/847139aa65aa4a5ee896375951cbf7417cfc8a4d6f277ec11f40cd87319f04aa
	chain := chains.BtcMainnetChain
	txHash := "847139aa65aa4a5ee896375951cbf7417cfc8a4d6f277ec11f40cd87319f04aa"
	scriptHex := "6a1467ed0bcc4e1256bc2ce87d22e190d63a120114bf"
	rawResult := testutils.LoadBTCIntxRawResult(t, TestDataDir, chain.ChainId, txHash, false)
	require.True(t, len(rawResult.Vout) >= 2)
	require.Equal(t, scriptHex, rawResult.Vout[1].ScriptPubKey.Hex)

	t.Run("should decode memo from OP_RETURN output", func(t *testing.T) {
		memo, found, err := DecodeOpReturnMemo(rawResult.Vout[1].ScriptPubKey.Hex, txHash)
		require.NoError(t, err)
		require.True(t, found)
		// [OP_RETURN, 0x14,<20-byte-hash>]
		require.Equal(t, scriptHex[4:], hex.EncodeToString(memo))
	})
	t.Run("should return nil memo non-OP_RETURN output", func(t *testing.T) {
		// modify the OP_RETURN to OP_1
		scriptInvalid := strings.Replace(scriptHex, "6a", "51", 1)
		memo, found, err := DecodeOpReturnMemo(scriptInvalid, txHash)
		require.NoError(t, err)
		require.False(t, found)
		require.Nil(t, memo)
	})
	t.Run("should return nil memo on invalid script", func(t *testing.T) {
		// use known short script
		scriptInvalid := "00"
		memo, found, err := DecodeOpReturnMemo(scriptInvalid, txHash)
		require.NoError(t, err)
		require.False(t, found)
		require.Nil(t, memo)
	})
}

func TestDecodeOpReturnMemoErrors(t *testing.T) {
	// https://mempool.space/tx/847139aa65aa4a5ee896375951cbf7417cfc8a4d6f277ec11f40cd87319f04aa
	txHash := "847139aa65aa4a5ee896375951cbf7417cfc8a4d6f277ec11f40cd87319f04aa"
	scriptHex := "6a1467ed0bcc4e1256bc2ce87d22e190d63a120114bf"

	t.Run("should return error on invalid memo size", func(t *testing.T) {
		// use invalid memo size
		scriptInvalid := strings.Replace(scriptHex, "6a14", "6axy", 1)
		memo, found, err := DecodeOpReturnMemo(scriptInvalid, txHash)
		require.ErrorContains(t, err, "error decoding memo size")
		require.False(t, found)
		require.Nil(t, memo)
	})
	t.Run("should return error on memo size mismatch", func(t *testing.T) {
		// use wrong memo size
		scriptInvalid := strings.Replace(scriptHex, "6a14", "6a13", 1)
		memo, found, err := DecodeOpReturnMemo(scriptInvalid, txHash)
		require.ErrorContains(t, err, "memo size mismatch")
		require.False(t, found)
		require.Nil(t, memo)
	})
	t.Run("should return error on invalid hex", func(t *testing.T) {
		// use invalid hex
		scriptInvalid := strings.Replace(scriptHex, "6a1467", "6a14xy", 1)
		memo, found, err := DecodeOpReturnMemo(scriptInvalid, txHash)
		require.ErrorContains(t, err, "error hex decoding memo")
		require.False(t, found)
		require.Nil(t, memo)
	})
	t.Run("should return nil memo on donation tx", func(t *testing.T) {
		// use donation sctipt "6a0a4920616d207269636821"
		scriptDonation := "6a0a" + hex.EncodeToString([]byte(constant.DonationMessage))
		memo, found, err := DecodeOpReturnMemo(scriptDonation, txHash)
		require.ErrorContains(t, err, "donation tx")
		require.False(t, found)
		require.Nil(t, memo)
	})
}

func TestDecodeTSSVout(t *testing.T) {
	chain := chains.BtcMainnetChain

	t.Run("should decode P2TR vout", func(t *testing.T) {
		// https://mempool.space/tx/259fc21e63e138136c8f19270a0f7ca10039a66a474f91d23a17896f46e677a7
		txHash := "259fc21e63e138136c8f19270a0f7ca10039a66a474f91d23a17896f46e677a7"
		rawResult := testutils.LoadBTCTxRawResult(t, TestDataDir, chain.ChainId, "P2TR", txHash)

		receiverExpected := "bc1p4scddlkkuw9486579autxumxmkvuphm5pz4jvf7f6pdh50p2uzqstawjt9"
		receiver, amount, err := DecodeTSSVout(rawResult.Vout[0], receiverExpected, chain)
		require.NoError(t, err)
		require.Equal(t, receiverExpected, receiver)
		require.Equal(t, int64(45000), amount)
	})
	t.Run("should decode P2WSH vout", func(t *testing.T) {
		// https://mempool.space/tx/791bb9d16f7ab05f70a116d18eaf3552faf77b9d5688699a480261424b4f7e53
		txHash := "791bb9d16f7ab05f70a116d18eaf3552faf77b9d5688699a480261424b4f7e53"
		rawResult := testutils.LoadBTCTxRawResult(t, TestDataDir, chain.ChainId, "P2WSH", txHash)

		receiverExpected := "bc1qqv6pwn470vu0tssdfha4zdk89v3c8ch5lsnyy855k9hcrcv3evequdmjmc"
		receiver, amount, err := DecodeTSSVout(rawResult.Vout[0], receiverExpected, chain)
		require.NoError(t, err)
		require.Equal(t, receiverExpected, receiver)
		require.Equal(t, int64(36557203), amount)
	})
	t.Run("should decode P2WPKH vout", func(t *testing.T) {
		// https://mempool.space/tx/5d09d232bfe41c7cb831bf53fc2e4029ab33a99087fd5328a2331b52ff2ebe5b
		txHash := "5d09d232bfe41c7cb831bf53fc2e4029ab33a99087fd5328a2331b52ff2ebe5b"
		rawResult := testutils.LoadBTCTxRawResult(t, TestDataDir, chain.ChainId, "P2WPKH", txHash)

		receiverExpected := "bc1qaxf82vyzy8y80v000e7t64gpten7gawewzu42y"
		receiver, amount, err := DecodeTSSVout(rawResult.Vout[0], receiverExpected, chain)
		require.NoError(t, err)
		require.Equal(t, receiverExpected, receiver)
		require.Equal(t, int64(79938), amount)
	})
	t.Run("should decode P2SH vout", func(t *testing.T) {
		// https://mempool.space/tx/fd68c8b4478686ca6f5ae4c28eaab055490650dbdaa6c2c8e380a7e075958a21
		txHash := "fd68c8b4478686ca6f5ae4c28eaab055490650dbdaa6c2c8e380a7e075958a21"
		rawResult := testutils.LoadBTCTxRawResult(t, TestDataDir, chain.ChainId, "P2SH", txHash)

		receiverExpected := "327z4GyFM8Y8DiYfasGKQWhRK4MvyMSEgE"
		receiver, amount, err := DecodeTSSVout(rawResult.Vout[0], receiverExpected, chain)
		require.NoError(t, err)
		require.Equal(t, receiverExpected, receiver)
		require.Equal(t, int64(1003881), amount)
	})
	t.Run("should decode P2PKH vout", func(t *testing.T) {
		// https://mempool.space/tx/9c741de6e17382b7a9113fc811e3558981a35a360e3d1262a6675892c91322ca
		txHash := "9c741de6e17382b7a9113fc811e3558981a35a360e3d1262a6675892c91322ca"
		rawResult := testutils.LoadBTCTxRawResult(t, TestDataDir, chain.ChainId, "P2PKH", txHash)

		receiverExpected := "1FueivsE338W2LgifJ25HhTcVJ7CRT8kte"
		receiver, amount, err := DecodeTSSVout(rawResult.Vout[0], receiverExpected, chain)
		require.NoError(t, err)
		require.Equal(t, receiverExpected, receiver)
		require.Equal(t, int64(1140000), amount)
	})
}

func TestDecodeTSSVoutErrors(t *testing.T) {
	// load archived tx raw result
	// https://mempool.space/tx/259fc21e63e138136c8f19270a0f7ca10039a66a474f91d23a17896f46e677a7
	chain := chains.BtcMainnetChain
	txHash := "259fc21e63e138136c8f19270a0f7ca10039a66a474f91d23a17896f46e677a7"

	rawResult := testutils.LoadBTCTxRawResult(t, TestDataDir, chain.ChainId, "P2TR", txHash)
	receiverExpected := "bc1p4scddlkkuw9486579autxumxmkvuphm5pz4jvf7f6pdh50p2uzqstawjt9"

	t.Run("should return error on invalid amount", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.Value = -0.05 // use negative amount
		receiver, amount, err := DecodeTSSVout(invalidVout, receiverExpected, chain)
		require.ErrorContains(t, err, "error getting satoshis")
		require.Empty(t, receiver)
		require.Zero(t, amount)
	})
	t.Run("should return error on invalid btc chain", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// use invalid chain
		invalidChain := chains.Chain{ChainId: 123}
		receiver, amount, err := DecodeTSSVout(invalidVout, receiverExpected, invalidChain)
		require.ErrorContains(t, err, "error GetBTCChainParams")
		require.Empty(t, receiver)
		require.Zero(t, amount)
	})
	t.Run("should return error when invalid receiver passed", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// use testnet params to decode mainnet receiver
		wrongChain := chains.BtcTestNetChain
		receiver, amount, err := DecodeTSSVout(invalidVout, "bc1qulmx8ej27cj0xe20953cztr2excnmsqvuh0s5c", wrongChain)
		require.ErrorContains(t, err, "error decoding receiver")
		require.Empty(t, receiver)
		require.Zero(t, amount)
	})
	t.Run("should return error on decoding failure", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Hex = "invalid script"
		receiver, amount, err := DecodeTSSVout(invalidVout, receiverExpected, chain)
		require.ErrorContains(t, err, "error decoding TSS vout")
		require.Empty(t, receiver)
		require.Zero(t, amount)
	})
}
