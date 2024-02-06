package zetaclient

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"
	"testing"

	"github.com/btcsuite/btcd/blockchain"
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
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	. "gopkg.in/check.v1"
)

type BTCSignerSuite struct {
	btcSigner *BTCSigner
}

var _ = Suite(&BTCSignerSuite{})

// 21 example UTXO txids to use in the test.
var exampleTxids = []string{
	"c1729638e1c9b6bfca57d11bf93047d98b65594b0bf75d7ee68bf7dc80dc164e",
	"54f9ebbd9e3ad39a297da54bf34a609b6831acbea0361cb5b7b5c8374f5046aa",
	"b18a55a34319cfbedebfcfe1a80fef2b92ad8894d06caf8293a0344824c2cfbc",
	"969fb309a4df7c299972700da788b5d601c0c04bab4ab46fff79d0335a7d75de",
	"6c71913061246ffc20e268c1b0e65895055c36bfbf1f8faf92dcad6f8242121e",
	"ba6d6e88cb5a97556684a1232719a3ffe409c5c9501061e1f59741bc412b3585",
	"69b56c3c8c5d1851f9eaec256cd49f290b477a5d43e2aef42ef25d3c1d9f4b33",
	"b87effd4cb46fe1a575b5b1ba0289313dc9b4bc9e615a3c6cbc0a14186921fdf",
	"3135433054523f5e220621c9e3d48efbbb34a6a2df65635c2a3e7d462d3e1cda",
	"8495c22a9ce6359ab53aa048c13b41c64fdf5fe141f516ba2573cc3f9313f06e",
	"f31583544b475370d7b9187c9a01b92e44fb31ac5fcfa7fc55565ac64043aa9a",
	"c03d55f9f717c1df978623e2e6b397b720999242f9ead7db9b5988fee3fb3933",
	"ee55688439b47a5410cdc05bac46be0094f3af54d307456fdfe6ba8caf336e0b",
	"61895f86c70f0bc3eef55d9a00347b509fa90f7a344606a9774be98a3ee9e02a",
	"ffabb401a19d04327bd4a076671d48467dbcde95459beeab23df21686fd01525",
	"b7e1c03b9b73e4e90fc06da893072c5604203c49e66699acbb2f61485d822981",
	"185614d21973990138e478ce10e0a4014352df58044276d4e4c0093aa140f482",
	"4a2800f13d15dc0c82308761d6fe8f6d13b65e42d7ca96a42a3a7048830e8c55",
	"fb98f52e91db500735b185797cebb5848afbfe1289922d87e03b98c3da5b85ef",
	"7901c5e36d9e8456ac61b29b82048650672a889596cbd30a9f8910a589ffc5b3",
	"6bcd0850fd2fa1404290ed04d78d4ae718414f16d4fbfd344951add8dcf60326",
}

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
	c.Assert(err, IsNil)

	{
		txWitness, err := txscript.WitnessSignature(redeemTx, txSigHashes, 0, 100000000, pkScript, txscript.SigHashAll, privKey, true)
		c.Assert(err, IsNil)
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
		c.Assert(err, IsNil)
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

func generateKeyPair(t *testing.T, net *chaincfg.Params) (*btcec.PrivateKey, []byte) {
	privateKey, err := btcec.NewPrivateKey(btcec.S256())
	assert.Nil(t, err)
	pubKeyHash := btcutil.Hash160(privateKey.PubKey().SerializeCompressed())
	addr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, net)
	assert.Nil(t, err)
	//fmt.Printf("New address: %s\n", addr.EncodeAddress())
	pkScript, err := payToWitnessPubKeyHashScript(addr.WitnessProgram())
	assert.Nil(t, err)
	return privateKey, pkScript
}

func addTxInputs(t *testing.T, tx *wire.MsgTx, txids []string) {
	preTxSize := tx.SerializeSize()
	for _, txid := range txids {
		hash, err := chainhash.NewHashFromStr(txid)
		assert.Nil(t, err)
		outpoint := wire.NewOutPoint(hash, uint32(rand.Intn(100)))
		txIn := wire.NewTxIn(outpoint, nil, nil)
		tx.AddTxIn(txIn)
		assert.Equal(t, bytesPerInput, tx.SerializeSize()-preTxSize)
		//fmt.Printf("tx size: %d, input %d size: %d\n", tx.SerializeSize(), i, tx.SerializeSize()-preTxSize)
		preTxSize = tx.SerializeSize()
	}
}

func addTxOutputs(t *testing.T, tx *wire.MsgTx, payerScript, payeeScript []byte) {
	preTxSize := tx.SerializeSize()

	// 1st output to payer
	value1 := int64(1 + rand.Intn(100000000))
	txOut1 := wire.NewTxOut(value1, payerScript)
	tx.AddTxOut(txOut1)
	assert.Equal(t, bytesPerOutput, tx.SerializeSize()-preTxSize)
	//fmt.Printf("tx size: %d, output 1: %d\n", tx.SerializeSize(), tx.SerializeSize()-preTxSize)
	preTxSize = tx.SerializeSize()

	// 2nd output to payee
	value2 := int64(1 + rand.Intn(100000000))
	txOut2 := wire.NewTxOut(value2, payeeScript)
	tx.AddTxOut(txOut2)
	assert.Equal(t, bytesPerOutput, tx.SerializeSize()-preTxSize)
	//fmt.Printf("tx size: %d, output 2: %d\n", tx.SerializeSize(), tx.SerializeSize()-preTxSize)
	preTxSize = tx.SerializeSize()

	// 3rd output to payee
	value3 := int64(1 + rand.Intn(100000000))
	txOut3 := wire.NewTxOut(value3, payeeScript)
	tx.AddTxOut(txOut3)
	assert.Equal(t, bytesPerOutput, tx.SerializeSize()-preTxSize)
	//fmt.Printf("tx size: %d, output 3: %d\n", tx.SerializeSize(), tx.SerializeSize()-preTxSize)
}

func signTx(t *testing.T, tx *wire.MsgTx, payerScript []byte, privateKey *btcec.PrivateKey) {
	preTxSize := tx.SerializeSize()
	sigHashes := txscript.NewTxSigHashes(tx)
	for ix := range tx.TxIn {
		amount := int64(1 + rand.Intn(100000000))
		witnessHash, err := txscript.CalcWitnessSigHash(payerScript, sigHashes, txscript.SigHashAll, tx, ix, amount)
		assert.Nil(t, err)
		sig, err := privateKey.Sign(witnessHash)
		assert.Nil(t, err)

		pkCompressed := privateKey.PubKey().SerializeCompressed()
		txWitness := wire.TxWitness{append(sig.Serialize(), byte(txscript.SigHashAll)), pkCompressed}
		tx.TxIn[ix].Witness = txWitness

		//fmt.Printf("tx size: %d, witness %d: %d\n", tx.SerializeSize(), ix+1, tx.SerializeSize()-preTxSize)
		if ix == 0 {
			bytesIncur := bytes1stWitness + len(tx.TxIn) - 1 // e.g., 130 bytes for a 21-input tx
			assert.True(t, tx.SerializeSize()-preTxSize >= bytesIncur-5)
			assert.True(t, tx.SerializeSize()-preTxSize <= bytesIncur+5)
		} else {
			assert.True(t, tx.SerializeSize()-preTxSize >= bytesPerWitness-5)
			assert.True(t, tx.SerializeSize()-preTxSize <= bytesPerWitness+5)
		}
		preTxSize = tx.SerializeSize()
	}
}

func TestP2WPHSize2In3Out(t *testing.T) {
	// Generate payer/payee private keys and P2WPKH addresss
	privateKey, payerScript := generateKeyPair(t, &chaincfg.TestNet3Params)
	_, payeeScript := generateKeyPair(t, &chaincfg.TestNet3Params)

	// 2 example UTXO txids to use in the test.
	utxosTxids := []string{
		"c1729638e1c9b6bfca57d11bf93047d98b65594b0bf75d7ee68bf7dc80dc164e",
		"54f9ebbd9e3ad39a297da54bf34a609b6831acbea0361cb5b7b5c8374f5046aa",
	}

	// Create a new transaction and add inputs
	tx := wire.NewMsgTx(wire.TxVersion)
	addTxInputs(t, tx, utxosTxids)

	// Add P2WPKH outputs
	addTxOutputs(t, tx, payerScript, payeeScript)

	// Payer sign the redeeming transaction.
	signTx(t, tx, payerScript, privateKey)

	// Estimate the tx size in vByte
	// #nosec G701 always positive
	vBytes := uint64(blockchain.GetTransactionWeight(btcutil.NewTx(tx)) / blockchain.WitnessScaleFactor)
	vBytesEstimated := EstimateSegWitTxSize(uint64(len(utxosTxids)), 3)
	assert.Equal(t, vBytes, vBytesEstimated)
	assert.Equal(t, vBytes, outTxBytesMin)
}

func TestP2WPHSize21In3Out(t *testing.T) {
	// Generate payer/payee private keys and P2WPKH addresss
	privateKey, payerScript := generateKeyPair(t, &chaincfg.TestNet3Params)
	_, payeeScript := generateKeyPair(t, &chaincfg.TestNet3Params)

	// Create a new transaction and add inputs
	tx := wire.NewMsgTx(wire.TxVersion)
	addTxInputs(t, tx, exampleTxids)

	// Add P2WPKH outputs
	addTxOutputs(t, tx, payerScript, payeeScript)

	// Payer sign the redeeming transaction.
	signTx(t, tx, payerScript, privateKey)

	// Estimate the tx size in vByte
	// #nosec G701 always positive
	vError := uint64(21 / 4) // 5 vBytes error tolerance
	vBytes := uint64(blockchain.GetTransactionWeight(btcutil.NewTx(tx)) / blockchain.WitnessScaleFactor)
	vBytesEstimated := EstimateSegWitTxSize(uint64(len(exampleTxids)), 3)
	assert.Equal(t, vBytesEstimated, outTxBytesMax)
	if vBytes > vBytesEstimated {
		assert.True(t, vBytes-vBytesEstimated <= vError)
	} else {
		assert.True(t, vBytesEstimated-vBytes <= vError)
	}
}

func TestP2WPHSizeXIn3Out(t *testing.T) {
	// Generate payer/payee private keys and P2WPKH addresss
	privateKey, payerScript := generateKeyPair(t, &chaincfg.TestNet3Params)
	_, payeeScript := generateKeyPair(t, &chaincfg.TestNet3Params)

	// Create new transactions with X (2 <= X <= 21) inputs and 3 outputs respectively
	for x := 2; x <= 21; x++ {
		tx := wire.NewMsgTx(wire.TxVersion)
		addTxInputs(t, tx, exampleTxids[:x])

		// Add P2WPKH outputs
		addTxOutputs(t, tx, payerScript, payeeScript)

		// Payer sign the redeeming transaction.
		signTx(t, tx, payerScript, privateKey)

		// Estimate the tx size
		// #nosec G701 always positive
		vError := uint64(0.25 + float64(x)/4) // 1st witness incur 0.25 vByte error, other witness incur 1/4 vByte error tolerance,
		vBytes := uint64(blockchain.GetTransactionWeight(btcutil.NewTx(tx)) / blockchain.WitnessScaleFactor)
		vBytesEstimated := EstimateSegWitTxSize(uint64(len(exampleTxids[:x])), 3)
		if vBytes > vBytesEstimated {
			assert.True(t, vBytes-vBytesEstimated <= vError)
			//fmt.Printf("%d error percentage: %.2f%%\n", float64(vBytes-vBytesEstimated)/float64(vBytes)*100)
		} else {
			assert.True(t, vBytesEstimated-vBytes <= vError)
			//fmt.Printf("error percentage: %.2f%%\n", float64(vBytesEstimated-vBytes)/float64(vBytes)*100)
		}
	}
}

func TestP2WPHSizeBreakdown(t *testing.T) {
	txSize2In3Out := EstimateSegWitTxSize(2, 3)
	assert.Equal(t, outTxBytesMin, txSize2In3Out)

	sz := EstimateSegWitTxSize(1, 1)
	fmt.Printf("1 input, 1 output: %d\n", sz)

	txSizeDepositor := SegWitTxSizeDepositor()
	assert.Equal(t, uint64(68), txSizeDepositor)

	txSizeWithdrawer := SegWitTxSizeWithdrawer()
	assert.Equal(t, uint64(171), txSizeWithdrawer)
	assert.Equal(t, txSize2In3Out, txSizeDepositor+txSizeWithdrawer) // 239 = 68 + 171

	depositFee := DepositorFee(20)
	assert.Equal(t, depositFee, 0.00001360)
}

// helper function to create a new BitcoinChainClient
func createTestClient(t *testing.T) *BitcoinChainClient {
	skHex := "7b8507ba117e069f4a3f456f505276084f8c92aee86ac78ae37b4d1801d35fa8"
	privateKey, err := crypto.HexToECDSA(skHex)
	assert.Nil(t, err)
	tss := TestSigner{
		PrivKey: privateKey,
	}
	tssAddress := tss.BTCAddressWitnessPubkeyHash().EncodeAddress()

	// Create BitcoinChainClient
	client := &BitcoinChainClient{
		Tss:               tss,
		Mu:                &sync.Mutex{},
		includedTxResults: make(map[string]*btcjson.GetTransactionResult),
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
	ob.includedTxResults[outTxID] = &btcjson.GetTransactionResult{TxID: txid}

	// Set nonce mark
	tssAddress := ob.Tss.BTCAddressWitnessPubkeyHash().EncodeAddress()
	nonceMark := btcjson.ListUnspentResult{TxID: txid, Address: tssAddress, Amount: float64(common.NonceMarkAmount(nonce)) * 1e-8}
	if preMarkIndex >= 0 { // replace nonce-mark utxo
		ob.utxos[preMarkIndex] = nonceMark

	} else { // add nonce-mark utxo directly
		ob.utxos = append(ob.utxos, nonceMark)
	}
	sort.SliceStable(ob.utxos, func(i, j int) bool {
		return ob.utxos[i].Amount < ob.utxos[j].Amount
	})
}

func TestSelectUTXOs(t *testing.T) {
	ob := createTestClient(t)
	dummyTxID := "6e6f71d281146c1fc5c755b35908ee449f26786c84e2ae18f98b268de40b7ec4"

	// Case1: nonce = 0, bootstrap
	// 		input: utxoCap = 5, amount = 0.01, nonce = 0
	// 		output: [0.01], 0.01
	result, amount, _, _, err := ob.SelectUTXOs(0.01, 5, 0, math.MaxUint16, true)
	assert.Nil(t, err)
	assert.Equal(t, 0.01, amount)
	assert.Equal(t, ob.utxos[0:1], result)

	// Case2: nonce = 1, must FAIL and wait for previous transaction to be mined
	// 		input: utxoCap = 5, amount = 0.5, nonce = 1
	// 		output: error
	result, amount, _, _, err = ob.SelectUTXOs(0.5, 5, 1, math.MaxUint16, true)
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Zero(t, amount)
	assert.Equal(t, "getOutTxidByNonce: cannot find outTx txid for nonce 0", err.Error())
	mineTxNSetNonceMark(ob, 0, dummyTxID, -1) // mine a transaction and set nonce-mark utxo for nonce 0

	// Case3: nonce = 1, should pass now
	// 		input: utxoCap = 5, amount = 0.5, nonce = 1
	// 		output: [0.00002, 0.01, 0.12, 0.18, 0.24], 0.55002
	result, amount, _, _, err = ob.SelectUTXOs(0.5, 5, 1, math.MaxUint16, true)
	assert.Nil(t, err)
	assert.Equal(t, 0.55002, amount)
	assert.Equal(t, ob.utxos[0:5], result)
	mineTxNSetNonceMark(ob, 1, dummyTxID, 0) // mine a transaction and set nonce-mark utxo for nonce 1

	// Case4:
	// 		input: utxoCap = 5, amount = 1.0, nonce = 2
	// 		output: [0.00002001, 0.01, 0.12, 0.18, 0.24, 0.5], 1.05002001
	result, amount, _, _, err = ob.SelectUTXOs(1.0, 5, 2, math.MaxUint16, true)
	assert.Nil(t, err)
	assert.InEpsilon(t, 1.05002001, amount, 1e-8)
	assert.Equal(t, ob.utxos[0:6], result)
	mineTxNSetNonceMark(ob, 2, dummyTxID, 0) // mine a transaction and set nonce-mark utxo for nonce 2

	// Case5: should include nonce-mark utxo on the LEFT
	// 		input: utxoCap = 5, amount = 8.05, nonce = 3
	// 		output: [0.00002002, 0.24, 0.5, 1.26, 2.97, 3.28], 8.25002002
	result, amount, _, _, err = ob.SelectUTXOs(8.05, 5, 3, math.MaxUint16, true)
	assert.Nil(t, err)
	assert.InEpsilon(t, 8.25002002, amount, 1e-8)
	expected := append([]btcjson.ListUnspentResult{ob.utxos[0]}, ob.utxos[4:9]...)
	assert.Equal(t, expected, result)
	mineTxNSetNonceMark(ob, 24105431, dummyTxID, 0) // mine a transaction and set nonce-mark utxo for nonce 24105431

	// Case6: should include nonce-mark utxo on the RIGHT
	// 		input: utxoCap = 5, amount = 0.503, nonce = 24105432
	// 		output: [0.24107432, 0.01, 0.12, 0.18, 0.24], 0.55002002
	result, amount, _, _, err = ob.SelectUTXOs(0.503, 5, 24105432, math.MaxUint16, true)
	assert.Nil(t, err)
	assert.InEpsilon(t, 0.79107431, amount, 1e-8)
	expected = append([]btcjson.ListUnspentResult{ob.utxos[4]}, ob.utxos[0:4]...)
	assert.Equal(t, expected, result)
	mineTxNSetNonceMark(ob, 24105432, dummyTxID, 4) // mine a transaction and set nonce-mark utxo for nonce 24105432

	// Case7: should include nonce-mark utxo in the MIDDLE
	// 		input: utxoCap = 5, amount = 1.0, nonce = 24105433
	// 		output: [0.24107432, 0.12, 0.18, 0.24, 0.5], 1.28107432
	result, amount, _, _, err = ob.SelectUTXOs(1.0, 5, 24105433, math.MaxUint16, true)
	assert.Nil(t, err)
	assert.InEpsilon(t, 1.28107432, amount, 1e-8)
	expected = append([]btcjson.ListUnspentResult{ob.utxos[4]}, ob.utxos[1:4]...)
	expected = append(expected, ob.utxos[5])
	assert.Equal(t, expected, result)

	// Case8: should work with maximum amount
	// 		input: utxoCap = 5, amount = 16.03
	// 		output: [0.24107432, 1.26, 2.97, 3.28, 5.16, 8.72], 21.63107432
	result, amount, _, _, err = ob.SelectUTXOs(16.03, 5, 24105433, math.MaxUint16, true)
	assert.Nil(t, err)
	assert.InEpsilon(t, 21.63107432, amount, 1e-8)
	expected = append([]btcjson.ListUnspentResult{ob.utxos[4]}, ob.utxos[6:11]...)
	assert.Equal(t, expected, result)

	// Case9: must FAIL due to insufficient funds
	// 		input: utxoCap = 5, amount = 21.64
	// 		output: error
	result, amount, _, _, err = ob.SelectUTXOs(21.64, 5, 24105433, math.MaxUint16, true)
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Zero(t, amount)
	assert.Equal(t, "SelectUTXOs: not enough btc in reserve - available : 21.63107432 , tx amount : 21.64", err.Error())
}

func TestUTXOConsolidation(t *testing.T) {
	dummyTxID := "6e6f71d281146c1fc5c755b35908ee449f26786c84e2ae18f98b268de40b7ec4"

	t.Run("should not consolidate", func(t *testing.T) {
		ob := createTestClient(t)
		mineTxNSetNonceMark(ob, 0, dummyTxID, -1) // mine a transaction and set nonce-mark utxo for nonce 0

		// input: utxoCap = 10, amount = 0.01, nonce = 1, rank = 10
		// output: [0.00002, 0.01], 0.01002
		result, amount, clsdtUtxo, clsdtValue, err := ob.SelectUTXOs(0.01, 10, 1, 10, true)
		assert.Nil(t, err)
		assert.Equal(t, 0.01002, amount)
		assert.Equal(t, ob.utxos[0:2], result)
		assert.Equal(t, uint16(0), clsdtUtxo)
		assert.Equal(t, 0.0, clsdtValue)
	})

	t.Run("should consolidate 1 utxo", func(t *testing.T) {
		ob := createTestClient(t)
		mineTxNSetNonceMark(ob, 0, dummyTxID, -1) // mine a transaction and set nonce-mark utxo for nonce 0

		// input: utxoCap = 9, amount = 0.01, nonce = 1, rank = 9
		// output: [0.00002, 0.01, 0.12], 0.13002
		result, amount, clsdtUtxo, clsdtValue, err := ob.SelectUTXOs(0.01, 9, 1, 9, true)
		assert.Nil(t, err)
		assert.Equal(t, 0.13002, amount)
		assert.Equal(t, ob.utxos[0:3], result)
		assert.Equal(t, uint16(1), clsdtUtxo)
		assert.Equal(t, 0.12, clsdtValue)
	})

	t.Run("should consolidate 3 utxos", func(t *testing.T) {
		ob := createTestClient(t)
		mineTxNSetNonceMark(ob, 0, dummyTxID, -1) // mine a transaction and set nonce-mark utxo for nonce 0

		// input: utxoCap = 5, amount = 0.01, nonce = 0, rank = 5
		// output: [0.00002, 0.014, 1.26, 0.5, 0.2], 2.01002
		result, amount, clsdtUtxo, clsdtValue, err := ob.SelectUTXOs(0.01, 5, 1, 5, true)
		assert.Nil(t, err)
		assert.Equal(t, 2.01002, amount)
		expected := make([]btcjson.ListUnspentResult, 2)
		copy(expected, ob.utxos[0:2])
		for i := 6; i >= 4; i-- { // append consolidated utxos in descending order
			expected = append(expected, ob.utxos[i])
		}
		assert.Equal(t, expected, result)
		assert.Equal(t, uint16(3), clsdtUtxo)
		assert.Equal(t, 2.0, clsdtValue)
	})

	t.Run("should consolidate all utxos using rank 1", func(t *testing.T) {
		ob := createTestClient(t)
		mineTxNSetNonceMark(ob, 0, dummyTxID, -1) // mine a transaction and set nonce-mark utxo for nonce 0

		// input: utxoCap = 12, amount = 0.01, nonce = 0, rank = 1
		// output: [0.00002, 0.01, 8.72, 5.16, 3.28, 2.97, 1.26, 0.5, 0.24, 0.18, 0.12], 22.44002
		result, amount, clsdtUtxo, clsdtValue, err := ob.SelectUTXOs(0.01, 12, 1, 1, true)
		assert.Nil(t, err)
		assert.Equal(t, 22.44002, amount)
		expected := make([]btcjson.ListUnspentResult, 2)
		copy(expected, ob.utxos[0:2])
		for i := 10; i >= 2; i-- { // append consolidated utxos in descending order
			expected = append(expected, ob.utxos[i])
		}
		assert.Equal(t, expected, result)
		assert.Equal(t, uint16(9), clsdtUtxo)
		assert.Equal(t, 22.43, clsdtValue)
	})

	t.Run("should consolidate 3 utxos sparse", func(t *testing.T) {
		ob := createTestClient(t)
		mineTxNSetNonceMark(ob, 24105431, dummyTxID, -1) // mine a transaction and set nonce-mark utxo for nonce 24105431

		// input: utxoCap = 5, amount = 0.13, nonce = 24105432, rank = 5
		// output: [0.24107431, 0.01, 0.12, 1.26, 0.5, 0.24], 2.37107431
		result, amount, clsdtUtxo, clsdtValue, err := ob.SelectUTXOs(0.13, 5, 24105432, 5, true)
		assert.Nil(t, err)
		assert.InEpsilon(t, 2.37107431, amount, 1e-8)
		expected := append([]btcjson.ListUnspentResult{ob.utxos[4]}, ob.utxos[0:2]...)
		expected = append(expected, ob.utxos[6])
		expected = append(expected, ob.utxos[5])
		expected = append(expected, ob.utxos[3])
		assert.Equal(t, expected, result)
		assert.Equal(t, uint16(3), clsdtUtxo)
		assert.Equal(t, 2.0, clsdtValue)
	})

	t.Run("should consolidate all utxos sparse", func(t *testing.T) {
		ob := createTestClient(t)
		mineTxNSetNonceMark(ob, 24105431, dummyTxID, -1) // mine a transaction and set nonce-mark utxo for nonce 24105431

		// input: utxoCap = 12, amount = 0.13, nonce = 24105432, rank = 1
		// output: [0.24107431, 0.01, 0.12, 8.72, 5.16, 3.28, 2.97, 1.26, 0.5, 0.24, 0.18], 22.68107431
		result, amount, clsdtUtxo, clsdtValue, err := ob.SelectUTXOs(0.13, 12, 24105432, 1, true)
		assert.Nil(t, err)
		assert.InEpsilon(t, 22.68107431, amount, 1e-8)
		expected := append([]btcjson.ListUnspentResult{ob.utxos[4]}, ob.utxos[0:2]...)
		for i := 10; i >= 5; i-- { // append consolidated utxos in descending order
			expected = append(expected, ob.utxos[i])
		}
		expected = append(expected, ob.utxos[3])
		expected = append(expected, ob.utxos[2])
		assert.Equal(t, expected, result)
		assert.Equal(t, uint16(8), clsdtUtxo)
		assert.Equal(t, 22.31, clsdtValue)
	})
}
