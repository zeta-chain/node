package bitcoin

import (
	"path"
	"strings"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
)

func TestDecodeVoutP2TR(t *testing.T) {
	// load archived tx raw result
	// https://mempool.space/tx/259fc21e63e138136c8f19270a0f7ca10039a66a474f91d23a17896f46e677a7
	chain := common.BtcMainnetChain()
	txHash := "259fc21e63e138136c8f19270a0f7ca10039a66a474f91d23a17896f46e677a7"
	net := &chaincfg.MainNetParams
	nameTx := path.Join("../", testutils.TestDataPathBTC, testutils.FileNameBTCTxByType(chain.ChainId, "P2TR", txHash))

	var rawResult btcjson.TxRawResult
	testutils.LoadObjectFromJSONFile(&rawResult, nameTx)
	require.Len(t, rawResult.Vout, 2)

	// decode vout 0, P2TR
	receiver, err := DecodeVoutP2TR(rawResult.Vout[0], net)
	require.NoError(t, err)
	require.Equal(t, "bc1p4scddlkkuw9486579autxumxmkvuphm5pz4jvf7f6pdh50p2uzqstawjt9", receiver)
}

func TestDecodeVoutP2TRErrors(t *testing.T) {
	// load archived tx raw result
	// https://mempool.space/tx/259fc21e63e138136c8f19270a0f7ca10039a66a474f91d23a17896f46e677a7
	chain := common.BtcMainnetChain()
	txHash := "259fc21e63e138136c8f19270a0f7ca10039a66a474f91d23a17896f46e677a7"
	net := &chaincfg.MainNetParams
	nameTx := path.Join("../", testutils.TestDataPathBTC, testutils.FileNameBTCTxByType(chain.ChainId, "P2TR", txHash))

	var rawResult btcjson.TxRawResult
	testutils.LoadObjectFromJSONFile(&rawResult, nameTx)

	t.Run("should return error on wrong script type", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Type = "witness_v0_keyhash" // use non-P2TR script type
		_, err := DecodeVoutP2TR(invalidVout, net)
		require.ErrorContains(t, err, "want scriptPubKey type witness_v1_taproot")
	})
	t.Run("should return error on invalid script", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Hex = "invalid script"
		_, err := DecodeVoutP2TR(invalidVout, net)
		require.ErrorContains(t, err, "error decoding scriptPubKey")
	})
	t.Run("should return error on wrong script length", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Hex = "0020" // 2 bytes, should be 34
		_, err := DecodeVoutP2TR(invalidVout, net)
		require.ErrorContains(t, err, "invalid P2TR scriptPubKey")
	})
	t.Run("should return error on invalid OP_1", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// modify the OP_1 '51' to OP_2 '52'
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "51", "52", 1)
		_, err := DecodeVoutP2TR(invalidVout, net)
		require.ErrorContains(t, err, "invalid P2TR scriptPubKey")
	})
	t.Run("should return error on wrong hash length", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// modify the length '20' to '19'
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "5120", "5119", 1)
		_, err := DecodeVoutP2TR(invalidVout, net)
		require.ErrorContains(t, err, "invalid P2TR scriptPubKey")
	})
}

func TestDecodeVoutP2WSH(t *testing.T) {
	// load archived tx raw result
	// https://mempool.space/tx/791bb9d16f7ab05f70a116d18eaf3552faf77b9d5688699a480261424b4f7e53
	chain := common.BtcMainnetChain()
	txHash := "791bb9d16f7ab05f70a116d18eaf3552faf77b9d5688699a480261424b4f7e53"
	net := &chaincfg.MainNetParams
	nameTx := path.Join("../", testutils.TestDataPathBTC, testutils.FileNameBTCTxByType(chain.ChainId, "P2WSH", txHash))

	var rawResult btcjson.TxRawResult
	testutils.LoadObjectFromJSONFile(&rawResult, nameTx)
	require.Len(t, rawResult.Vout, 1)

	// decode vout 0, P2WSH
	receiver, err := DecodeVoutP2WSH(rawResult.Vout[0], net)
	require.NoError(t, err)
	require.Equal(t, "bc1qqv6pwn470vu0tssdfha4zdk89v3c8ch5lsnyy855k9hcrcv3evequdmjmc", receiver)
}

func TestDecodeVoutP2WSHErrors(t *testing.T) {
	// load archived tx raw result
	// https://mempool.space/tx/791bb9d16f7ab05f70a116d18eaf3552faf77b9d5688699a480261424b4f7e53
	chain := common.BtcMainnetChain()
	txHash := "791bb9d16f7ab05f70a116d18eaf3552faf77b9d5688699a480261424b4f7e53"
	net := &chaincfg.MainNetParams
	nameTx := path.Join("../", testutils.TestDataPathBTC, testutils.FileNameBTCTxByType(chain.ChainId, "P2WSH", txHash))

	var rawResult btcjson.TxRawResult
	testutils.LoadObjectFromJSONFile(&rawResult, nameTx)

	t.Run("should return error on wrong script type", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Type = "witness_v0_keyhash" // use non-P2WSH script type
		_, err := DecodeVoutP2WSH(invalidVout, net)
		require.ErrorContains(t, err, "want scriptPubKey type witness_v0_scripthash")
	})
	t.Run("should return error on invalid script", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Hex = "invalid script"
		_, err := DecodeVoutP2WSH(invalidVout, net)
		require.ErrorContains(t, err, "error decoding scriptPubKey")
	})
	t.Run("should return error on wrong script length", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Hex = "0020" // 2 bytes, should be 34
		_, err := DecodeVoutP2WSH(invalidVout, net)
		require.ErrorContains(t, err, "invalid P2WSH scriptPubKey")
	})
	t.Run("should return error on invalid OP_0", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// modify the OP_0 '00' to OP_1 '51'
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "00", "51", 1)
		_, err := DecodeVoutP2WSH(invalidVout, net)
		require.ErrorContains(t, err, "invalid P2WSH scriptPubKey")
	})
	t.Run("should return error on wrong hash length", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// modify the length '20' to '19'
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "0020", "0019", 1)
		_, err := DecodeVoutP2WSH(invalidVout, net)
		require.ErrorContains(t, err, "invalid P2WSH scriptPubKey")
	})
}

func TestDecodeP2WPKHVout(t *testing.T) {
	// load archived outtx raw result
	// https://mempool.space/tx/030cd813443f7b70cc6d8a544d320c6d8465e4528fc0f3410b599dc0b26753a0
	chain := common.BtcMainnetChain()
	nonce := uint64(148)
	net := &chaincfg.MainNetParams
	nameTx := path.Join("../", testutils.TestDataPathBTC, testutils.FileNameBTCOuttx(chain.ChainId, nonce))

	var rawResult btcjson.TxRawResult
	testutils.LoadObjectFromJSONFile(&rawResult, nameTx)
	require.Len(t, rawResult.Vout, 3)

	// decode vout 0, nonce mark 148
	receiver, err := DecodeVoutP2WPKH(rawResult.Vout[0], net)
	require.NoError(t, err)
	require.Equal(t, testutils.TSSAddressBTCMainnet, receiver)

	// decode vout 1, payment 0.00012000 BTC
	receiver, err = DecodeVoutP2WPKH(rawResult.Vout[1], net)
	require.NoError(t, err)
	require.Equal(t, "bc1qpsdlklfcmlcfgm77c43x65ddtrt7n0z57hsyjp", receiver)

	// decode vout 2, change 0.39041489 BTC
	receiver, err = DecodeVoutP2WPKH(rawResult.Vout[2], net)
	require.NoError(t, err)
	require.Equal(t, testutils.TSSAddressBTCMainnet, receiver)
}

func TestDecodeP2WPKHVoutErrors(t *testing.T) {
	// load archived outtx raw result
	// https://mempool.space/tx/030cd813443f7b70cc6d8a544d320c6d8465e4528fc0f3410b599dc0b26753a0
	chain := common.BtcMainnetChain()
	nonce := uint64(148)
	net := &chaincfg.MainNetParams
	nameTx := path.Join("../", testutils.TestDataPathBTC, testutils.FileNameBTCOuttx(chain.ChainId, nonce))

	var rawResult btcjson.TxRawResult
	testutils.LoadObjectFromJSONFile(&rawResult, nameTx)

	t.Run("should return error on wrong script type", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Type = "scripthash" // use non-P2WPKH script type
		_, err := DecodeVoutP2WPKH(invalidVout, net)
		require.ErrorContains(t, err, "want scriptPubKey type witness_v0_keyhash")
	})
	t.Run("should return error on invalid script", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Hex = "invalid script"
		_, err := DecodeVoutP2WPKH(invalidVout, net)
		require.ErrorContains(t, err, "error decoding scriptPubKey")
	})
	t.Run("should return error on wrong script length", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Hex = "0014" // 2 bytes, should be 22
		_, err := DecodeVoutP2WPKH(invalidVout, net)
		require.ErrorContains(t, err, "invalid P2WPKH scriptPubKey")
	})
	t.Run("should return error on wrong hash length", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// modify the length '14' to '13'
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "0014", "0013", 1)
		_, err := DecodeVoutP2WPKH(invalidVout, net)
		require.ErrorContains(t, err, "invalid P2WPKH scriptPubKey")
	})
}

func TestDecodeVoutP2SH(t *testing.T) {
	// load archived tx raw result
	// https://mempool.space/tx/fd68c8b4478686ca6f5ae4c28eaab055490650dbdaa6c2c8e380a7e075958a21
	chain := common.BtcMainnetChain()
	txHash := "fd68c8b4478686ca6f5ae4c28eaab055490650dbdaa6c2c8e380a7e075958a21"
	net := &chaincfg.MainNetParams
	nameTx := path.Join("../", testutils.TestDataPathBTC, testutils.FileNameBTCTxByType(chain.ChainId, "P2SH", txHash))

	var rawResult btcjson.TxRawResult
	testutils.LoadObjectFromJSONFile(&rawResult, nameTx)
	require.Len(t, rawResult.Vout, 2)

	// decode vout 0, P2SH
	receiver, err := DecodeVoutP2SH(rawResult.Vout[0], net)
	require.NoError(t, err)
	require.Equal(t, "327z4GyFM8Y8DiYfasGKQWhRK4MvyMSEgE", receiver)
}

func TestDecodeVoutP2SHErrors(t *testing.T) {
	// load archived tx raw result
	// https://mempool.space/tx/fd68c8b4478686ca6f5ae4c28eaab055490650dbdaa6c2c8e380a7e075958a21
	chain := common.BtcMainnetChain()
	txHash := "fd68c8b4478686ca6f5ae4c28eaab055490650dbdaa6c2c8e380a7e075958a21"
	net := &chaincfg.MainNetParams
	nameTx := path.Join("../", testutils.TestDataPathBTC, testutils.FileNameBTCTxByType(chain.ChainId, "P2SH", txHash))

	var rawResult btcjson.TxRawResult
	testutils.LoadObjectFromJSONFile(&rawResult, nameTx)

	t.Run("should return error on wrong script type", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Type = "witness_v0_keyhash" // use non-P2SH script type
		_, err := DecodeVoutP2SH(invalidVout, net)
		require.ErrorContains(t, err, "want scriptPubKey type scripthash")
	})
	t.Run("should return error on invalid script", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Hex = "invalid script"
		_, err := DecodeVoutP2SH(invalidVout, net)
		require.ErrorContains(t, err, "error decoding scriptPubKey")
	})
	t.Run("should return error on wrong script length", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Hex = "0014" // 2 bytes, should be 23
		_, err := DecodeVoutP2SH(invalidVout, net)
		require.ErrorContains(t, err, "invalid P2SH scriptPubKey")
	})
	t.Run("should return error on invalid OP_HASH160", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// modify the OP_HASH160 'a9' to OP_HASH256 'aa'
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "a9", "aa", 1)
		_, err := DecodeVoutP2SH(invalidVout, net)
		require.ErrorContains(t, err, "invalid P2SH scriptPubKey")
	})
	t.Run("should return error on wrong data length", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// modify the length '14' to '13'
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "a914", "a913", 1)
		_, err := DecodeVoutP2SH(invalidVout, net)
		require.ErrorContains(t, err, "invalid P2SH scriptPubKey")
	})
	t.Run("should return error on invalid OP_EQUAL", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "87", "88", 1)
		_, err := DecodeVoutP2SH(invalidVout, net)
		require.ErrorContains(t, err, "invalid P2SH scriptPubKey")
	})
}

func TestDecodeVoutP2PKH(t *testing.T) {
	// load archived tx raw result
	// https://mempool.space/tx/9c741de6e17382b7a9113fc811e3558981a35a360e3d1262a6675892c91322ca
	chain := common.BtcMainnetChain()
	txHash := "9c741de6e17382b7a9113fc811e3558981a35a360e3d1262a6675892c91322ca"
	net := &chaincfg.MainNetParams
	nameTx := path.Join("../", testutils.TestDataPathBTC, testutils.FileNameBTCTxByType(chain.ChainId, "P2PKH", txHash))

	var rawResult btcjson.TxRawResult
	testutils.LoadObjectFromJSONFile(&rawResult, nameTx)
	require.Len(t, rawResult.Vout, 2)

	// decode vout 0, P2PKH
	receiver, err := DecodeVoutP2PKH(rawResult.Vout[0], net)
	require.NoError(t, err)
	require.Equal(t, "1FueivsE338W2LgifJ25HhTcVJ7CRT8kte", receiver)
}

func TestDecodeVoutP2PKHErrors(t *testing.T) {
	// load archived tx raw result
	// https://mempool.space/tx/9c741de6e17382b7a9113fc811e3558981a35a360e3d1262a6675892c91322ca
	chain := common.BtcMainnetChain()
	txHash := "9c741de6e17382b7a9113fc811e3558981a35a360e3d1262a6675892c91322ca"
	net := &chaincfg.MainNetParams
	nameTx := path.Join("../", testutils.TestDataPathBTC, testutils.FileNameBTCTxByType(chain.ChainId, "P2PKH", txHash))

	var rawResult btcjson.TxRawResult
	testutils.LoadObjectFromJSONFile(&rawResult, nameTx)

	t.Run("should return error on wrong script type", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Type = "scripthash" // use non-P2PKH script type
		_, err := DecodeVoutP2PKH(invalidVout, net)
		require.ErrorContains(t, err, "want scriptPubKey type pubkeyhash")
	})
	t.Run("should return error on invalid script", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Hex = "invalid script"
		_, err := DecodeVoutP2PKH(invalidVout, net)
		require.ErrorContains(t, err, "error decoding scriptPubKey")
	})
	t.Run("should return error on wrong script length", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		invalidVout.ScriptPubKey.Hex = "76a914" // 3 bytes, should be 25
		_, err := DecodeVoutP2PKH(invalidVout, net)
		require.ErrorContains(t, err, "invalid P2PKH scriptPubKey")
	})
	t.Run("should return error on invalid OP_DUP", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// modify the OP_DUP '76' to OP_NIP '77'
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "76", "77", 1)
		_, err := DecodeVoutP2PKH(invalidVout, net)
		require.ErrorContains(t, err, "invalid P2PKH scriptPubKey")
	})
	t.Run("should return error on invalid OP_HASH160", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// modify the OP_HASH160 'a9' to OP_HASH256 'aa'
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "76a9", "76aa", 1)
		_, err := DecodeVoutP2PKH(invalidVout, net)
		require.ErrorContains(t, err, "invalid P2PKH scriptPubKey")
	})
	t.Run("should return error on wrong data length", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// modify the length '14' to '13'
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "76a914", "76a913", 1)
		_, err := DecodeVoutP2PKH(invalidVout, net)
		require.ErrorContains(t, err, "invalid P2PKH scriptPubKey")
	})
	t.Run("should return error on invalid OP_EQUALVERIFY", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// modify the OP_EQUALVERIFY '88' to OP_RESERVED1 '89'
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "88ac", "89ac", 1)
		_, err := DecodeVoutP2PKH(invalidVout, net)
		require.ErrorContains(t, err, "invalid P2PKH scriptPubKey")
	})
	t.Run("should return error on invalid OP_CHECKSIG", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// modify the OP_CHECKSIG 'ac' to OP_CHECKSIGVERIFY 'ad'
		invalidVout.ScriptPubKey.Hex = strings.Replace(invalidVout.ScriptPubKey.Hex, "88ac", "88ad", 1)
		_, err := DecodeVoutP2PKH(invalidVout, net)
		require.ErrorContains(t, err, "invalid P2PKH scriptPubKey")
	})
}

func TestDecodeTSSVout(t *testing.T) {
	chain := common.BtcMainnetChain()

	t.Run("should decode P2TR vout", func(t *testing.T) {
		// https://mempool.space/tx/259fc21e63e138136c8f19270a0f7ca10039a66a474f91d23a17896f46e677a7
		txHash := "259fc21e63e138136c8f19270a0f7ca10039a66a474f91d23a17896f46e677a7"
		nameTx := path.Join("../", testutils.TestDataPathBTC, testutils.FileNameBTCTxByType(chain.ChainId, "P2TR", txHash))
		var rawResult btcjson.TxRawResult
		testutils.LoadObjectFromJSONFile(&rawResult, nameTx)

		receiverExpected := "bc1p4scddlkkuw9486579autxumxmkvuphm5pz4jvf7f6pdh50p2uzqstawjt9"
		receiver, amount, err := DecodeTSSVout(rawResult.Vout[0], receiverExpected, chain)
		require.NoError(t, err)
		require.Equal(t, receiverExpected, receiver)
		require.Equal(t, int64(45000), amount)
	})
	t.Run("should decode P2WSH vout", func(t *testing.T) {
		// https://mempool.space/tx/791bb9d16f7ab05f70a116d18eaf3552faf77b9d5688699a480261424b4f7e53
		txHash := "791bb9d16f7ab05f70a116d18eaf3552faf77b9d5688699a480261424b4f7e53"
		nameTx := path.Join("../", testutils.TestDataPathBTC, testutils.FileNameBTCTxByType(chain.ChainId, "P2WSH", txHash))
		var rawResult btcjson.TxRawResult
		testutils.LoadObjectFromJSONFile(&rawResult, nameTx)

		receiverExpected := "bc1qqv6pwn470vu0tssdfha4zdk89v3c8ch5lsnyy855k9hcrcv3evequdmjmc"
		receiver, amount, err := DecodeTSSVout(rawResult.Vout[0], receiverExpected, chain)
		require.NoError(t, err)
		require.Equal(t, receiverExpected, receiver)
		require.Equal(t, int64(36557203), amount)
	})
	t.Run("should decode P2WPKH vout", func(t *testing.T) {
		// https://mempool.space/tx/5d09d232bfe41c7cb831bf53fc2e4029ab33a99087fd5328a2331b52ff2ebe5b
		txHash := "5d09d232bfe41c7cb831bf53fc2e4029ab33a99087fd5328a2331b52ff2ebe5b"
		nameTx := path.Join("../", testutils.TestDataPathBTC, testutils.FileNameBTCTxByType(chain.ChainId, "P2WPKH", txHash))
		var rawResult btcjson.TxRawResult
		testutils.LoadObjectFromJSONFile(&rawResult, nameTx)

		receiverExpected := "bc1qaxf82vyzy8y80v000e7t64gpten7gawewzu42y"
		receiver, amount, err := DecodeTSSVout(rawResult.Vout[0], receiverExpected, chain)
		require.NoError(t, err)
		require.Equal(t, receiverExpected, receiver)
		require.Equal(t, int64(79938), amount)
	})
	t.Run("should decode P2SH vout", func(t *testing.T) {
		// https://mempool.space/tx/fd68c8b4478686ca6f5ae4c28eaab055490650dbdaa6c2c8e380a7e075958a21
		txHash := "fd68c8b4478686ca6f5ae4c28eaab055490650dbdaa6c2c8e380a7e075958a21"
		nameTx := path.Join("../", testutils.TestDataPathBTC, testutils.FileNameBTCTxByType(chain.ChainId, "P2SH", txHash))
		var rawResult btcjson.TxRawResult
		testutils.LoadObjectFromJSONFile(&rawResult, nameTx)

		receiverExpected := "327z4GyFM8Y8DiYfasGKQWhRK4MvyMSEgE"
		receiver, amount, err := DecodeTSSVout(rawResult.Vout[0], receiverExpected, chain)
		require.NoError(t, err)
		require.Equal(t, receiverExpected, receiver)
		require.Equal(t, int64(1003881), amount)
	})
	t.Run("should decode P2PKH vout", func(t *testing.T) {
		// https://mempool.space/tx/9c741de6e17382b7a9113fc811e3558981a35a360e3d1262a6675892c91322ca
		txHash := "9c741de6e17382b7a9113fc811e3558981a35a360e3d1262a6675892c91322ca"
		nameTx := path.Join("../", testutils.TestDataPathBTC, testutils.FileNameBTCTxByType(chain.ChainId, "P2PKH", txHash))
		var rawResult btcjson.TxRawResult
		testutils.LoadObjectFromJSONFile(&rawResult, nameTx)

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
	chain := common.BtcMainnetChain()
	txHash := "259fc21e63e138136c8f19270a0f7ca10039a66a474f91d23a17896f46e677a7"
	nameTx := path.Join("../", testutils.TestDataPathBTC, testutils.FileNameBTCTxByType(chain.ChainId, "P2TR", txHash))

	var rawResult btcjson.TxRawResult
	testutils.LoadObjectFromJSONFile(&rawResult, nameTx)
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
		invalidChain := common.Chain{ChainId: 123}
		receiver, amount, err := DecodeTSSVout(invalidVout, receiverExpected, invalidChain)
		require.ErrorContains(t, err, "error GetBTCChainParams")
		require.Empty(t, receiver)
		require.Zero(t, amount)
	})
	t.Run("should return error when invalid receiver passed", func(t *testing.T) {
		invalidVout := rawResult.Vout[0]
		// use testnet params to decode mainnet receiver
		wrongChain := common.BtcTestNetChain()
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
