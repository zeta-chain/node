package bitcoin

import (
	"encoding/hex"
	"path"
	"strings"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/zetaclient/testutils"
)

// the relative path to the testdata directory
var TestDataDir = "../../"

func TestDecodeVoutP2TR(t *testing.T) {
	// load archived tx raw result
	// https://mempool.space/tx/259fc21e63e138136c8f19270a0f7ca10039a66a474f91d23a17896f46e677a7
	chain := chains.BitcoinMainnet
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
	chain := chains.BitcoinMainnet
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
	chain := chains.BitcoinMainnet
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
	chain := chains.BitcoinMainnet
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
	// load archived outbound raw result
	// https://mempool.space/tx/030cd813443f7b70cc6d8a544d320c6d8465e4528fc0f3410b599dc0b26753a0
	chain := chains.BitcoinMainnet
	nonce := uint64(148)
	net := &chaincfg.MainNetParams
	nameTx := path.Join(TestDataDir, testutils.TestDataPathBTC, testutils.FileNameBTCOutbound(chain.ChainId, nonce))

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
	// load archived outbound raw result
	// https://mempool.space/tx/030cd813443f7b70cc6d8a544d320c6d8465e4528fc0f3410b599dc0b26753a0
	chain := chains.BitcoinMainnet
	nonce := uint64(148)
	net := &chaincfg.MainNetParams
	nameTx := path.Join(TestDataDir, testutils.TestDataPathBTC, testutils.FileNameBTCOutbound(chain.ChainId, nonce))

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
	chain := chains.BitcoinMainnet
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
	chain := chains.BitcoinMainnet
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
	chain := chains.BitcoinMainnet
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
	chain := chains.BitcoinMainnet
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
	// load archived inbound raw result
	// https://mempool.space/tx/847139aa65aa4a5ee896375951cbf7417cfc8a4d6f277ec11f40cd87319f04aa
	chain := chains.BitcoinMainnet
	txHash := "847139aa65aa4a5ee896375951cbf7417cfc8a4d6f277ec11f40cd87319f04aa"
	scriptHex := "6a1467ed0bcc4e1256bc2ce87d22e190d63a120114bf"
	rawResult := testutils.LoadBTCInboundRawResult(t, TestDataDir, chain.ChainId, txHash, false)
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
	chain := chains.BitcoinMainnet

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
	chain := chains.BitcoinMainnet
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
		wrongChain := chains.BitcoinTestnet
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

func TestDecodeScript(t *testing.T) {
	t.Run("should decode longer data ok", func(t *testing.T) {
		// 600 bytes of random data generated offline
		data := "2001a7bae79bd61c2368fe41a565061d6cf22b4f509fbc1652caea06d98b8fd0c7ac00634d0802c7faa771dd05f27993d22c42988758882d20080241074462884c8774e1cdf4b04e5b3b74b6568bd1769722708306c66270b6b2a7f68baced83627eeeb2d494e8a1749277b92a4c5a90b1b4f6038e5f704405515109d4d0021612ad298b8dad6e12245f8f0020e11a7a319652ba6abe261958201ce5e83131cd81302c0ecec60d4afa9f72540fc84b6b9c1f3d903ab25686df263b192a403a4aa22b799ba24369c49ff4042012589a07d4211e05f80f18a1262de5a1577ce0ec9e1fa9283cfa25d98d7d0b4217951dfcb8868570318c63f1e1424cfdb7d7a33c6b9e3ced4b2ffa0178b3a5fac8bace2991e382a402f56a2c6a9191463740910056483e4fd0f5ac729ffac66bf1b3ec4570c4e75c116f7d9fd65718ec3ed6c7647bf335b77e7d6a4e2011276dc8031b78403a1ad82c92fb339ec916c263b6dd0f003ba4381ad5410e90e88effbfa7f961b8e8a6011c525643a434f7abe2c1928a892cc57d6291831216c4e70cb80a39a79a3889211070e767c23db396af9b4c2093c3743d8cbcbfcb73d29361ecd3857e94ab3c800be1299fd36a5685ec60607a60d8c2e0f99ff0b8b9e86354d39a43041f7d552e95fe2d33b6fc0f540715da0e7e1b344c778afe73f82d00881352207b719f67dcb00b4ff645974d4fd7711363d26400e2852890cb6ea9cbfe63ac43080870049b1023be984331560c6350bb64da52b4b81bc8910934915f0a96701f4c50646d5386146596443bee9b2d116706e1687697fb42542196c1d764419c23a914896f9212946518ac59e1ba5d1fc37e503313133ebdf2ced5785e0eaa9738fe3f9ad73646e733931ebb7cff26e96106fe68"
		script, _ := hex.DecodeString(data)

		memo, isFound, err := DecodeScript(script)
		require.Nil(t, err)
		require.True(t, isFound)

		// the expected memo
		expected := "c7faa771dd05f27993d22c42988758882d20080241074462884c8774e1cdf4b04e5b3b74b6568bd1769722708306c66270b6b2a7f68baced83627eeeb2d494e8a1749277b92a4c5a90b1b4f6038e5f704405515109d4d0021612ad298b8dad6e12245f8f0020e11a7a319652ba6abe261958201ce5e83131cd81302c0ecec60d4afa9f72540fc84b6b9c1f3d903ab25686df263b192a403a4aa22b799ba24369c49ff4042012589a07d4211e05f80f18a1262de5a1577ce0ec9e1fa9283cfa25d98d7d0b4217951dfcb8868570318c63f1e1424cfdb7d7a33c6b9e3ced4b2ffa0178b3a5fac8bace2991e382a402f56a2c6a9191463740910056483e4fd0f5ac729ffac66bf1b3ec4570c4e75c116f7d9fd65718ec3ed6c7647bf335b77e7d6a4e2011276dc8031b78403a1ad82c92fb339ec916c263b6dd0f003ba4381ad5410e90e88effbfa7f961b8e8a6011c525643a434f7abe2c1928a892cc57d6291831216c4e70cb80a39a79a3889211070e767c23db396af9b4c2093c3743d8cbcbfcb73d29361ecd3857e94ab3c800be1299fd36a5685ec60607a60d8c2e0f99ff0b8b9e86354d39a43041f7d552e95fe2d33b6fc0f540715da0e7e1b344c778afe73f82d00881352207b719f67dcb00b4ff645974d4fd7711363d26400e2852890cb6ea9cbfe63ac43080870049b1023be984331560c6350bb64da52b4b81bc8910934915f0a96701f646d5386146596443bee9b2d116706e1687697fb42542196c1d764419c23a914896f9212946518ac59e1ba5d1fc37e503313133ebdf2ced5785e0eaa9738fe3f9ad73646e733931ebb7cff26e96106fe"
		require.Equal(t, hex.EncodeToString(memo), expected)
	})

	t.Run("should decode shorter data ok", func(t *testing.T) {
		// 81 bytes of random data generated offline
		data := "20d6f59371037bf30115d9fd6016f0e3ef552cdfc0367ee20aa9df3158f74aaeb4ac00634c51bdd33073d76f6b4ae6510d69218100575eafabadd16e5faf9f42bd2fbbae402078bdcaa4c0413ce96d053e3c0bbd4d5944d6857107d640c248bdaaa7de959d9c1e6b9962b51428e5a554c28c397160881668"
		script, _ := hex.DecodeString(data)

		memo, isFound, err := DecodeScript(script)
		require.Nil(t, err)
		require.True(t, isFound)

		// the expected memo
		expected := "bdd33073d76f6b4ae6510d69218100575eafabadd16e5faf9f42bd2fbbae402078bdcaa4c0413ce96d053e3c0bbd4d5944d6857107d640c248bdaaa7de959d9c1e6b9962b51428e5a554c28c3971608816"
		require.Equal(t, hex.EncodeToString(memo), expected)
	})

	t.Run("decode error due to missing data byte", func(t *testing.T) {
		// missing OP_ENDIF at the end
		data := "20cabd6ecc0245c40f27ca6299dcd3732287c317f3946734f04e27568fc5334218ac00634d0802000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004c500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000068"
		script, _ := hex.DecodeString(data)

		memo, isFound, err := DecodeScript(script)
		require.ErrorContains(t, err, "should contain more data, but script ended")
		require.False(t, isFound)
		require.Nil(t, memo)
	})

	t.Run("decode error due to missing data for public key", func(t *testing.T) {
		// missing OP_ENDIF at the end
		data := "2001a7bae79bd61c2368fe41a565061d6cf22b4f509fbc1652caea06d98b8fd0"
		script, _ := hex.DecodeString(data)

		memo, isFound, err := DecodeScript(script)
		require.ErrorContains(t, err, "cannot obtain public key bytes")
		require.False(t, isFound)
		require.Nil(t, memo)
	})

	t.Run("decode error due to missing OP_CHECKSIG", func(t *testing.T) {
		// missing OP_ENDIF at the end
		data := "2001a7bae79bd61c2368fe41a565061d6cf22b4f509fbc1652caea06d98b8fd0c7ab"
		script, _ := hex.DecodeString(data)

		memo, isFound, err := DecodeScript(script)
		require.ErrorContains(t, err, "cannot parse OP_CHECKSIG")
		require.False(t, isFound)
		require.Nil(t, memo)
	})
}
