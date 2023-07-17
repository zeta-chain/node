package zetaclient

import (
	"encoding/hex"
	"fmt"
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
	"github.com/stretchr/testify/require"
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
	s.btcSigner, err = NewBTCSigner(&tss, nil, zerolog.Logger{}, &TelemetryServer{})
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

func TestSelectUTXOs(t *testing.T) {
	// Create 10 dummy UTXOs (22.44 BTC in total)
	utxos := make([]btcjson.ListUnspentResult, 0, 10)
	amounts := []float64{0.01, 0.12, 0.18, 0.24, 0.5, 1.26, 2.97, 3.28, 5.16, 8.72}
	for _, amount := range amounts {
		utxos = append(utxos, btcjson.ListUnspentResult{Amount: amount})
	}

	// Case1:
	// 		input: utxoCap = 5, amount = 0.01,
	// 		output: [0.01], 0.01
	result, amount, err := selectUTXOs(utxos, 0.01, 5, 0, "")
	require.Nil(t, err)
	require.Equal(t, 0.01, amount)
	require.Equal(t, utxos[0:1], result)

	// Case2:
	// 		input: utxoCap = 5, amount = 0.5
	// 		output: [0.01, 0.12, 0.18, 0.24], 0.55
	result, amount, err = selectUTXOs(utxos, 0.5, 5, 0, "")
	require.Nil(t, err)
	require.Equal(t, 0.55, amount)
	require.Equal(t, utxos[0:4], result)

	// Case3:
	// 		input: utxoCap = 5, amount = 1.0
	// 		output: [0.01, 0.12, 0.18, 0.24, 0.5], 1.05
	result, amount, err = selectUTXOs(utxos, 1.0, 5, 0, "")
	require.Nil(t, err)
	require.Equal(t, 1.05, amount)
	require.Equal(t, utxos[0:5], result)

	// Case4:
	// 		input: utxoCap = 5, amount = 8.05
	// 		output: [0.24, 0.5, 1.26, 2.97, 3.28], 8.25
	result, amount, err = selectUTXOs(utxos, 8.05, 5, 0, "")
	require.Nil(t, err)
	require.Equal(t, 8.25, amount)
	require.Equal(t, utxos[3:8], result)

	// Case5:
	// 		input: utxoCap = 5, amount = 16.03
	// 		output: [1.26, 2.97, 3.28, 5.16, 8.72], 21.39
	result, amount, err = selectUTXOs(utxos, 16.03, 5, 0, "")
	require.Nil(t, err)
	require.Equal(t, 21.39, amount)
	require.Equal(t, utxos[5:10], result)

	// Case6:
	// 		input: utxoCap = 5, amount = 21.4
	// 		output: error
	result, amount, err = selectUTXOs(utxos, 21.4, 5, 0, "")
	require.NotNil(t, err)
	require.Nil(t, result)
	require.Equal(t, 0.0, amount)
	require.Equal(t, "not enough btc in reserve - available : 21.39 , tx amount : 21.4", err.Error())

	// TODO: add a case with nonce > 0 so that a utxo with value 2000+nonce needs to be selected
}
