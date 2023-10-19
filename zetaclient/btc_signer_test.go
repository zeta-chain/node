package zetaclient

import (
	"encoding/hex"
	"fmt"
	"sort"
	"sync"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	. "gopkg.in/check.v1"
)

type BTCSignerSuite struct {
	btcSigner *BTCSigner
}

var _ = Suite(&BTCSignerSuite{})

func (s *BTCSignerSuite) SetUpTest(c *C) {
	// test private key with EVM address
	//// EVM: 0x236C7f53a90493Bb423411fe4117Cb4c2De71DfB
	// BTC testnet3: muGe9prUBjQwEnX19zG26fVRHNi8z7kSPo
	skHex := "7b8507ba117e069f4a3f456f505276084f8c92aee86ac78ae37b4d1801d35fa8"
	privateKey, err := crypto.HexToECDSA(skHex)
	pkBytes := crypto.FromECDSAPub(&privateKey.PublicKey)
	c.Logf("pubkey: %d", len(pkBytes))
	// Uncomment the following code to generate new random private key pairs
	//privateKey, err := crypto.GenerateKey()
	//privkeyBytes := crypto.FromECDSA(privateKey)
	//c.Logf("privatekey %s", hex.EncodeToString(privkeyBytes))
	c.Assert(err, IsNil)
	tss := TestSigner{
		PrivKey: privateKey,
	}
	s.btcSigner, err = NewBTCSigner(config.BTCConfig{}, &tss, zerolog.Logger{}, &TelemetryServer{})
	c.Assert(err, IsNil)
}

func (s *BTCSignerSuite) TestP2PH(c *C) {
	// Ordinarily the private key would come from whatever storage mechanism
	// is being used, but for this example just hard code it.
	privKeyBytes, err := hex.DecodeString("22a47fa09a223f2aa079edf85a7c2" +
		"d4f8720ee63e502ee2869afab7de234b80c")
	c.Assert(err, IsNil)

	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), privKeyBytes)
	pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())
	addr, err := btcutil.NewAddressPubKeyHash(pubKeyHash, &chaincfg.RegressionNetParams)
	c.Assert(err, IsNil)

	// For this example, create a fake transaction that represents what
	// would ordinarily be the real transaction that is being spent.  It
	// contains a single output that pays to address in the amount of 1 BTC.
	originTx := wire.NewMsgTx(wire.TxVersion)
	prevOut := wire.NewOutPoint(&chainhash.Hash{}, ^uint32(0))
	txIn := wire.NewTxIn(prevOut, []byte{txscript.OP_0, txscript.OP_0}, nil)
	originTx.AddTxIn(txIn)
	pkScript, err := txscript.PayToAddrScript(addr)

	c.Assert(err, IsNil)

	txOut := wire.NewTxOut(100000000, pkScript)
	originTx.AddTxOut(txOut)
	originTxHash := originTx.TxHash()

	// Create the transaction to redeem the fake transaction.
	redeemTx := wire.NewMsgTx(wire.TxVersion)

	// Add the input(s) the redeeming transaction will spend.  There is no
	// signature script at this point since it hasn't been created or signed
	// yet, hence nil is provided for it.
	prevOut = wire.NewOutPoint(&originTxHash, 0)
	txIn = wire.NewTxIn(prevOut, nil, nil)
	redeemTx.AddTxIn(txIn)

	// Ordinarily this would contain that actual destination of the funds,
	// but for this example don't bother.
	txOut = wire.NewTxOut(0, nil)
	redeemTx.AddTxOut(txOut)

	// Sign the redeeming transaction.
	lookupKey := func(a btcutil.Address) (*btcec.PrivateKey, bool, error) {
		return privKey, true, nil
	}
	// Notice that the script database parameter is nil here since it isn't
	// used.  It must be specified when pay-to-script-hash transactions are
	// being signed.
	sigScript, err := txscript.SignTxOutput(&chaincfg.MainNetParams,
		redeemTx, 0, originTx.TxOut[0].PkScript, txscript.SigHashAll,
		txscript.KeyClosure(lookupKey), nil, nil)
	c.Assert(err, IsNil)

	redeemTx.TxIn[0].SignatureScript = sigScript

	// Prove that the transaction has been validly signed by executing the
	// script pair.
	flags := txscript.ScriptBip16 | txscript.ScriptVerifyDERSignatures |
		txscript.ScriptStrictMultiSig |
		txscript.ScriptDiscourageUpgradableNops
	vm, err := txscript.NewEngine(originTx.TxOut[0].PkScript, redeemTx, 0,
		flags, nil, nil, -1)
	c.Assert(err, IsNil)

	err = vm.Execute()
	c.Assert(err, IsNil)

	fmt.Println("Transaction successfully signed")
}

func (s *BTCSignerSuite) TestP2WPH(c *C) {
	// Ordinarily the private key would come from whatever storage mechanism
	// is being used, but for this example just hard code it.
	privKeyBytes, err := hex.DecodeString("22a47fa09a223f2aa079edf85a7c2" +
		"d4f8720ee63e502ee2869afab7de234b80c")
	c.Assert(err, IsNil)

	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), privKeyBytes)
	pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())
	//addr, err := btcutil.NewAddressPubKeyHash(pubKeyHash, &chaincfg.RegressionNetParams)
	addr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, &chaincfg.RegressionNetParams)
	c.Assert(err, IsNil)

	// For this example, create a fake transaction that represents what
	// would ordinarily be the real transaction that is being spent.  It
	// contains a single output that pays to address in the amount of 1 BTC.
	originTx := wire.NewMsgTx(wire.TxVersion)
	prevOut := wire.NewOutPoint(&chainhash.Hash{}, ^uint32(0))
	txIn := wire.NewTxIn(prevOut, []byte{txscript.OP_0, txscript.OP_0}, nil)
	originTx.AddTxIn(txIn)
	pkScript, err := txscript.PayToAddrScript(addr)
	c.Assert(err, IsNil)
	txOut := wire.NewTxOut(100000000, pkScript)
	originTx.AddTxOut(txOut)
	originTxHash := originTx.TxHash()

	// Create the transaction to redeem the fake transaction.
	redeemTx := wire.NewMsgTx(wire.TxVersion)

	// Add the input(s) the redeeming transaction will spend.  There is no
	// signature script at this point since it hasn't been created or signed
	// yet, hence nil is provided for it.
	prevOut = wire.NewOutPoint(&originTxHash, 0)
	txIn = wire.NewTxIn(prevOut, nil, nil)
	redeemTx.AddTxIn(txIn)

	// Ordinarily this would contain that actual destination of the funds,
	// but for this example don't bother.
	txOut = wire.NewTxOut(0, nil)
	redeemTx.AddTxOut(txOut)
	txSigHashes := txscript.NewTxSigHashes(redeemTx)
	pkScript, err = payToWitnessPubKeyHashScript(addr.WitnessProgram())

	{
		txWitness, err := txscript.WitnessSignature(redeemTx, txSigHashes, 0, 100000000, pkScript, txscript.SigHashAll, privKey, true)
		redeemTx.TxIn[0].Witness = txWitness
		// Prove that the transaction has been validly signed by executing the
		// script pair.
		flags := txscript.ScriptBip16 | txscript.ScriptVerifyDERSignatures |
			txscript.ScriptStrictMultiSig |
			txscript.ScriptDiscourageUpgradableNops
		vm, err := txscript.NewEngine(originTx.TxOut[0].PkScript, redeemTx, 0,
			flags, nil, nil, -1)
		c.Assert(err, IsNil)

		err = vm.Execute()
		c.Assert(err, IsNil)
	}

	{
		witnessHash, err := txscript.CalcWitnessSigHash(pkScript, txSigHashes, txscript.SigHashAll, redeemTx, 0, 100000000)
		c.Assert(err, IsNil)
		sig, err := privKey.Sign(witnessHash)
		txWitness := wire.TxWitness{append(sig.Serialize(), byte(txscript.SigHashAll)), pubKeyHash}
		redeemTx.TxIn[0].Witness = txWitness

		flags := txscript.ScriptBip16 | txscript.ScriptVerifyDERSignatures |
			txscript.ScriptStrictMultiSig |
			txscript.ScriptDiscourageUpgradableNops
		vm, err := txscript.NewEngine(originTx.TxOut[0].PkScript, redeemTx, 0,
			flags, nil, nil, -1)
		c.Assert(err, IsNil)

		err = vm.Execute()
		c.Assert(err, IsNil)
	}

	fmt.Println("Transaction successfully signed")
}

// helper function to create a new BitcoinChainClient
func createTestClient(t *testing.T) *BitcoinChainClient {
	skHex := "7b8507ba117e069f4a3f456f505276084f8c92aee86ac78ae37b4d1801d35fa8"
	privateKey, err := crypto.HexToECDSA(skHex)
	require.Nil(t, err)
	tss := TestSigner{
		PrivKey: privateKey,
	}
	tssAddress := tss.BTCAddressWitnessPubkeyHash().EncodeAddress()

	// Create BitcoinChainClient
	client := &BitcoinChainClient{
		Tss:               tss,
		Mu:                &sync.Mutex{},
		includedTxResults: make(map[string]btcjson.GetTransactionResult),
	}

	// Create 10 dummy UTXOs (22.44 BTC in total)
	client.utxos = make([]btcjson.ListUnspentResult, 0, 10)
	amounts := []float64{0.01, 0.12, 0.18, 0.24, 0.5, 1.26, 2.97, 3.28, 5.16, 8.72}
	for _, amount := range amounts {
		client.utxos = append(client.utxos, btcjson.ListUnspentResult{Address: tssAddress, Amount: amount})
	}
	return client
}

func mineTxNSetNonceMark(ob *BitcoinChainClient, nonce uint64, txid string, preMarkIndex int) {
	// Mine transaction
	outTxID := ob.GetTxID(nonce)
	ob.includedTxResults[outTxID] = btcjson.GetTransactionResult{TxID: txid}

	// Set nonce mark
	if preMarkIndex >= 0 {
		tssAddress := ob.Tss.BTCAddressWitnessPubkeyHash().EncodeAddress()
		nonceMark := btcjson.ListUnspentResult{TxID: txid, Address: tssAddress, Amount: float64(common.NonceMarkAmount(nonce)) * 1e-8}
		ob.utxos[preMarkIndex] = nonceMark
		sort.SliceStable(ob.utxos, func(i, j int) bool {
			return ob.utxos[i].Amount < ob.utxos[j].Amount
		})
	}
}

func TestSelectUTXOs(t *testing.T) {
	ob := createTestClient(t)
	tssAddress := ob.Tss.BTCAddressWitnessPubkeyHash().EncodeAddress()
	dummyTxID := "6e6f71d281146c1fc5c755b35908ee449f26786c84e2ae18f98b268de40b7ec4"

	// Case1: nonce = 0, bootstrap
	// 		input: utxoCap = 5, amount = 0.01, nonce = 0
	// 		output: [0.01], 0.01
	result, amount, err := ob.SelectUTXOs(0.01, 5, 0, true)
	require.Nil(t, err)
	require.Equal(t, 0.01, amount)
	require.Equal(t, ob.utxos[0:1], result)

	// Case2: nonce = 1, must FAIL and wait for previous transaction to be mined
	// 		input: utxoCap = 5, amount = 0.5, nonce = 1
	// 		output: error
	result, amount, err = ob.SelectUTXOs(0.5, 5, 1, true)
	require.NotNil(t, err)
	require.Nil(t, result)
	require.Zero(t, amount)
	require.Equal(t, "getOutTxidByNonce: cannot find outTx txid for nonce 0", err.Error())
	mineTxNSetNonceMark(ob, 0, dummyTxID, -1) // mine a transaction for nonce 0

	// Case3: nonce = 1, must FAIL without nonce mark utxo
	// 		input: utxoCap = 5, amount = 0.5, nonce = 1
	// 		output: error
	result, amount, err = ob.SelectUTXOs(0.5, 5, 1, true)
	require.NotNil(t, err)
	require.Nil(t, result)
	require.Zero(t, amount)
	require.Equal(t, "findNonceMarkUTXO: cannot find nonce-mark utxo with nonce 0", err.Error())

	// add nonce-mark utxo for nonce 0
	nonceMark0 := btcjson.ListUnspentResult{TxID: dummyTxID, Address: tssAddress, Amount: float64(common.NonceMarkAmount(0)) * 1e-8}
	ob.utxos = append([]btcjson.ListUnspentResult{nonceMark0}, ob.utxos...)

	// Case4: nonce = 1, should pass now
	// 		input: utxoCap = 5, amount = 0.5, nonce = 1
	// 		output: [0.00002, 0.01, 0.12, 0.18, 0.24], 0.55002
	result, amount, err = ob.SelectUTXOs(0.5, 5, 1, true)
	require.Nil(t, err)
	require.Equal(t, 0.55002, amount)
	require.Equal(t, ob.utxos[0:5], result)
	mineTxNSetNonceMark(ob, 1, dummyTxID, 0) // mine a transaction and set nonce-mark utxo for nonce 1

	// Case5:
	// 		input: utxoCap = 5, amount = 1.0, nonce = 2
	// 		output: [0.00002001, 0.01, 0.12, 0.18, 0.24, 0.5], 1.05002001
	result, amount, err = ob.SelectUTXOs(1.0, 5, 2, true)
	require.Nil(t, err)
	assert.InEpsilon(t, 1.05002001, amount, 1e-8)
	require.Equal(t, ob.utxos[0:6], result)
	mineTxNSetNonceMark(ob, 2, dummyTxID, 0) // mine a transaction and set nonce-mark utxo for nonce 2

	// Case6: should include nonce-mark utxo on the LEFT
	// 		input: utxoCap = 5, amount = 8.05, nonce = 3
	// 		output: [0.00002002, 0.24, 0.5, 1.26, 2.97, 3.28], 8.25002002
	result, amount, err = ob.SelectUTXOs(8.05, 5, 3, true)
	require.Nil(t, err)
	assert.InEpsilon(t, 8.25002002, amount, 1e-8)
	expected := append([]btcjson.ListUnspentResult{ob.utxos[0]}, ob.utxos[4:9]...)
	require.Equal(t, expected, result)
	mineTxNSetNonceMark(ob, 24105431, dummyTxID, 0) // mine a transaction and set nonce-mark utxo for nonce 24105431

	// Case7: should include nonce-mark utxo on the RIGHT
	// 		input: utxoCap = 5, amount = 0.503, nonce = 24105432
	// 		output: [0.24107432, 0.01, 0.12, 0.18, 0.24], 0.55002002
	result, amount, err = ob.SelectUTXOs(0.503, 5, 24105432, true)
	require.Nil(t, err)
	assert.InEpsilon(t, 0.79107431, amount, 1e-8)
	expected = append([]btcjson.ListUnspentResult{ob.utxos[4]}, ob.utxos[0:4]...)
	require.Equal(t, expected, result)
	mineTxNSetNonceMark(ob, 24105432, dummyTxID, 4) // mine a transaction and set nonce-mark utxo for nonce 24105432

	// Case8: should include nonce-mark utxo in the MIDDLE
	// 		input: utxoCap = 5, amount = 1.0, nonce = 24105433
	// 		output: [0.24107432, 0.12, 0.18, 0.24, 0.5], 1.28107432
	result, amount, err = ob.SelectUTXOs(1.0, 5, 24105433, true)
	require.Nil(t, err)
	assert.InEpsilon(t, 1.28107432, amount, 1e-8)
	expected = append([]btcjson.ListUnspentResult{ob.utxos[4]}, ob.utxos[1:4]...)
	expected = append(expected, ob.utxos[5])
	require.Equal(t, expected, result)

	// Case9: should work with maximum amount
	// 		input: utxoCap = 5, amount = 16.03
	// 		output: [0.24107432, 1.26, 2.97, 3.28, 5.16, 8.72], 21.63107432
	result, amount, err = ob.SelectUTXOs(16.03, 5, 24105433, true)
	require.Nil(t, err)
	assert.InEpsilon(t, 21.63107432, amount, 1e-8)
	expected = append([]btcjson.ListUnspentResult{ob.utxos[4]}, ob.utxos[6:11]...)
	require.Equal(t, expected, result)

	// Case10: must FAIL due to insufficient funds
	// 		input: utxoCap = 5, amount = 21.64
	// 		output: error
	result, amount, err = ob.SelectUTXOs(21.64, 5, 24105433, true)
	require.NotNil(t, err)
	require.Nil(t, result)
	require.Zero(t, amount)
	require.Equal(t, "SelectUTXOs: not enough btc in reserve - available : 21.63107432 , tx amount : 21.64", err.Error())
}
