package bitcoin

import (
	"math/rand"
	"testing"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/pkg/chains"
)

const (
	// btc address script types
	ScriptTypeP2TR   = "witness_v1_taproot"
	ScriptTypeP2WSH  = "witness_v0_scripthash"
	ScriptTypeP2WPKH = "witness_v0_keyhash"
	ScriptTypeP2SH   = "scripthash"
	ScriptTypeP2PKH  = "pubkeyhash"
)

var testAddressMap = map[string]string{
	ScriptTypeP2TR:   "bc1p4scddlkkuw9486579autxumxmkvuphm5pz4jvf7f6pdh50p2uzqstawjt9",
	ScriptTypeP2WSH:  "bc1qqv6pwn470vu0tssdfha4zdk89v3c8ch5lsnyy855k9hcrcv3evequdmjmc",
	ScriptTypeP2WPKH: "bc1qaxf82vyzy8y80v000e7t64gpten7gawewzu42y",
	ScriptTypeP2SH:   "327z4GyFM8Y8DiYfasGKQWhRK4MvyMSEgE",
	ScriptTypeP2PKH:  "1FueivsE338W2LgifJ25HhTcVJ7CRT8kte",
}

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

func generateKeyPair(t *testing.T, net *chaincfg.Params) (*btcec.PrivateKey, btcutil.Address, []byte) {
	privateKey, err := btcec.NewPrivateKey(btcec.S256())
	require.NoError(t, err)
	pubKeyHash := btcutil.Hash160(privateKey.PubKey().SerializeCompressed())
	addr, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, net)
	require.NoError(t, err)
	//fmt.Printf("New address: %s\n", addr.EncodeAddress())
	pkScript, err := PayToAddrScript(addr)
	require.NoError(t, err)
	return privateKey, addr, pkScript
}

// getTestAddrScript returns hard coded test address scripts by script type
func getTestAddrScript(t *testing.T, scriptType string) btcutil.Address {
	chain := chains.BitcoinMainnet
	inputAddress, found := testAddressMap[scriptType]
	require.True(t, found)
	address, err := chains.DecodeBtcAddress(inputAddress, chain.ChainId)
	require.NoError(t, err)
	return address
}

// createPkScripts creates 10 random amount of scripts to the given address 'to'
func createPkScripts(t *testing.T, to btcutil.Address, repeat int) ([]btcutil.Address, [][]byte) {
	pkScript, err := PayToAddrScript(to)
	require.NoError(t, err)

	addrs := []btcutil.Address{}
	pkScripts := [][]byte{}
	for i := 0; i < repeat; i++ {
		addrs = append(addrs, to)
		pkScripts = append(pkScripts, pkScript)
	}
	return addrs, pkScripts
}

func addTxInputs(t *testing.T, tx *wire.MsgTx, txids []string) {
	preTxSize := tx.SerializeSize()
	for _, txid := range txids {
		hash, err := chainhash.NewHashFromStr(txid)
		require.NoError(t, err)
		outpoint := wire.NewOutPoint(hash, uint32(rand.Intn(100)))
		txIn := wire.NewTxIn(outpoint, nil, nil)
		tx.AddTxIn(txIn)
		require.Equal(t, bytesPerInput, tx.SerializeSize()-preTxSize)
		//fmt.Printf("tx size: %d, input %d size: %d\n", tx.SerializeSize(), i, tx.SerializeSize()-preTxSize)
		preTxSize = tx.SerializeSize()
	}
}

func addTxOutputs(t *testing.T, tx *wire.MsgTx, payerScript []byte, payeeScripts [][]byte) {
	preTxSize := tx.SerializeSize()

	// 1st output to payer
	value1 := int64(1 + rand.Intn(100000000))
	txOut1 := wire.NewTxOut(value1, payerScript)
	tx.AddTxOut(txOut1)
	require.Equal(t, bytesPerOutputP2WPKH, tx.SerializeSize()-preTxSize)
	//fmt.Printf("tx size: %d, output 1: %d\n", tx.SerializeSize(), tx.SerializeSize()-preTxSize)
	preTxSize = tx.SerializeSize()

	// output to payee list
	for _, payeeScript := range payeeScripts {
		value := int64(1 + rand.Intn(100000000))
		txOut := wire.NewTxOut(value, payeeScript)
		tx.AddTxOut(txOut)
		//fmt.Printf("tx size: %d, output %d: %d\n", tx.SerializeSize(), i+1, tx.SerializeSize()-preTxSize)
		preTxSize = tx.SerializeSize()
	}

	// 3rd output to payee
	value3 := int64(1 + rand.Intn(100000000))
	txOut3 := wire.NewTxOut(value3, payerScript)
	tx.AddTxOut(txOut3)
	require.Equal(t, bytesPerOutputP2WPKH, tx.SerializeSize()-preTxSize)
	//fmt.Printf("tx size: %d, last output: %d\n", tx.SerializeSize(), tx.SerializeSize()-preTxSize)
}

func addTxInputsOutputsAndSignTx(
	t *testing.T, tx *wire.MsgTx,
	privateKey *btcec.PrivateKey,
	payerScript []byte,
	txids []string,
	payeeScripts [][]byte) {
	// Add inputs
	addTxInputs(t, tx, txids)

	// Add outputs
	addTxOutputs(t, tx, payerScript, payeeScripts)

	// Payer sign the redeeming transaction.
	signTx(t, tx, payerScript, privateKey)
}

func signTx(t *testing.T, tx *wire.MsgTx, payerScript []byte, privateKey *btcec.PrivateKey) {
	preTxSize := tx.SerializeSize()
	sigHashes := txscript.NewTxSigHashes(tx)
	for ix := range tx.TxIn {
		amount := int64(1 + rand.Intn(100000000))
		witnessHash, err := txscript.CalcWitnessSigHash(payerScript, sigHashes, txscript.SigHashAll, tx, ix, amount)
		require.NoError(t, err)
		sig, err := privateKey.Sign(witnessHash)
		require.NoError(t, err)

		pkCompressed := privateKey.PubKey().SerializeCompressed()
		txWitness := wire.TxWitness{append(sig.Serialize(), byte(txscript.SigHashAll)), pkCompressed}
		tx.TxIn[ix].Witness = txWitness

		//fmt.Printf("tx size: %d, witness %d: %d\n", tx.SerializeSize(), ix+1, tx.SerializeSize()-preTxSize)
		if ix == 0 {
			bytesIncur := bytes1stWitness + len(tx.TxIn) - 1 // e.g., 130 bytes for a 21-input tx
			require.True(t, tx.SerializeSize()-preTxSize >= bytesIncur-5)
			require.True(t, tx.SerializeSize()-preTxSize <= bytesIncur+5)
		} else {
			require.True(t, tx.SerializeSize()-preTxSize >= bytesPerWitness-5)
			require.True(t, tx.SerializeSize()-preTxSize <= bytesPerWitness+5)
		}
		preTxSize = tx.SerializeSize()
	}
}

func TestOutboundSize2In3Out(t *testing.T) {
	// Generate payer/payee private keys and P2WPKH addresss
	privateKey, _, payerScript := generateKeyPair(t, &chaincfg.TestNet3Params)
	_, payee, payeeScript := generateKeyPair(t, &chaincfg.TestNet3Params)

	// 2 example UTXO txids to use in the test.
	utxosTxids := exampleTxids[:2]

	// Create a new transaction
	tx := wire.NewMsgTx(wire.TxVersion)

	// Add inputs and outputs and sign the transaction
	addTxInputsOutputsAndSignTx(t, tx, privateKey, payerScript, utxosTxids, [][]byte{payeeScript})

	// Estimate the tx size in vByte
	// #nosec G115 always positive
	vError := uint64(1) // 1 vByte error tolerance
	vBytes := uint64(blockchain.GetTransactionWeight(btcutil.NewTx(tx)) / blockchain.WitnessScaleFactor)
	vBytesEstimated, err := EstimateOutboundSize(uint64(len(utxosTxids)), []btcutil.Address{payee})
	require.NoError(t, err)
	if vBytes > vBytesEstimated {
		require.True(t, vBytes-vBytesEstimated <= vError)
	} else {
		require.True(t, vBytesEstimated-vBytes <= vError)
	}
}

func TestOutboundSize21In3Out(t *testing.T) {
	// Generate payer/payee private keys and P2WPKH addresss
	privateKey, _, payerScript := generateKeyPair(t, &chaincfg.TestNet3Params)
	_, payee, payeeScript := generateKeyPair(t, &chaincfg.TestNet3Params)

	// Create a new transaction
	tx := wire.NewMsgTx(wire.TxVersion)

	// Add inputs and outputs and sign the transaction
	addTxInputsOutputsAndSignTx(t, tx, privateKey, payerScript, exampleTxids, [][]byte{payeeScript})

	// Estimate the tx size in vByte
	// #nosec G115 always positive
	vError := uint64(21 / 4) // 5 vBytes error tolerance
	vBytes := uint64(blockchain.GetTransactionWeight(btcutil.NewTx(tx)) / blockchain.WitnessScaleFactor)
	vBytesEstimated, err := EstimateOutboundSize(uint64(len(exampleTxids)), []btcutil.Address{payee})
	require.NoError(t, err)
	if vBytes > vBytesEstimated {
		require.True(t, vBytes-vBytesEstimated <= vError)
	} else {
		require.True(t, vBytesEstimated-vBytes <= vError)
	}
}

func TestOutboundSizeXIn3Out(t *testing.T) {
	// Generate payer/payee private keys and P2WPKH addresss
	privateKey, _, payerScript := generateKeyPair(t, &chaincfg.TestNet3Params)
	_, payee, payeeScript := generateKeyPair(t, &chaincfg.TestNet3Params)

	// Create new transactions with X (2 <= X <= 21) inputs and 3 outputs respectively
	for x := 2; x <= 21; x++ {
		// Create transaction. Add inputs and outputs and sign the transaction
		tx := wire.NewMsgTx(wire.TxVersion)
		addTxInputsOutputsAndSignTx(t, tx, privateKey, payerScript, exampleTxids[:x], [][]byte{payeeScript})

		// Estimate the tx size
		// #nosec G115 always positive
		vError := uint64(
			0.25 + float64(x)/4,
		) // 1st witness incurs 0.25 more vByte error than others (which incurs 1/4 vByte per witness)
		vBytes := uint64(blockchain.GetTransactionWeight(btcutil.NewTx(tx)) / blockchain.WitnessScaleFactor)
		vBytesEstimated, err := EstimateOutboundSize(uint64(len(exampleTxids[:x])), []btcutil.Address{payee})
		require.NoError(t, err)
		if vBytes > vBytesEstimated {
			require.True(t, vBytes-vBytesEstimated <= vError)
			//fmt.Printf("%d error percentage: %.2f%%\n", float64(vBytes-vBytesEstimated)/float64(vBytes)*100)
		} else {
			require.True(t, vBytesEstimated-vBytes <= vError)
			//fmt.Printf("error percentage: %.2f%%\n", float64(vBytesEstimated-vBytes)/float64(vBytes)*100)
		}
	}
}

func TestGetOutputSizeByAddress(t *testing.T) {
	// test nil P2TR address and non-nil P2TR address
	nilP2TR := (*chains.AddressTaproot)(nil)
	sizeNilP2TR, err := GetOutputSizeByAddress(nilP2TR)
	require.NoError(t, err)
	require.Equal(t, uint64(0), sizeNilP2TR)

	addrP2TR := getTestAddrScript(t, ScriptTypeP2TR)
	sizeP2TR, err := GetOutputSizeByAddress(addrP2TR)
	require.NoError(t, err)
	require.Equal(t, uint64(bytesPerOutputP2TR), sizeP2TR)

	// test nil P2WSH address and non-nil P2WSH address
	nilP2WSH := (*btcutil.AddressWitnessScriptHash)(nil)
	sizeNilP2WSH, err := GetOutputSizeByAddress(nilP2WSH)
	require.NoError(t, err)
	require.Equal(t, uint64(0), sizeNilP2WSH)

	addrP2WSH := getTestAddrScript(t, ScriptTypeP2WSH)
	sizeP2WSH, err := GetOutputSizeByAddress(addrP2WSH)
	require.NoError(t, err)
	require.Equal(t, uint64(bytesPerOutputP2WSH), sizeP2WSH)

	// test nil P2WPKH address and non-nil P2WPKH address
	nilP2WPKH := (*btcutil.AddressWitnessPubKeyHash)(nil)
	sizeNilP2WPKH, err := GetOutputSizeByAddress(nilP2WPKH)
	require.NoError(t, err)
	require.Equal(t, uint64(0), sizeNilP2WPKH)

	addrP2WPKH := getTestAddrScript(t, ScriptTypeP2WPKH)
	sizeP2WPKH, err := GetOutputSizeByAddress(addrP2WPKH)
	require.NoError(t, err)
	require.Equal(t, uint64(bytesPerOutputP2WPKH), sizeP2WPKH)

	// test nil P2SH address and non-nil P2SH address
	nilP2SH := (*btcutil.AddressScriptHash)(nil)
	sizeNilP2SH, err := GetOutputSizeByAddress(nilP2SH)
	require.NoError(t, err)
	require.Equal(t, uint64(0), sizeNilP2SH)

	addrP2SH := getTestAddrScript(t, ScriptTypeP2SH)
	sizeP2SH, err := GetOutputSizeByAddress(addrP2SH)
	require.NoError(t, err)
	require.Equal(t, uint64(bytesPerOutputP2SH), sizeP2SH)

	// test nil P2PKH address and non-nil P2PKH address
	nilP2PKH := (*btcutil.AddressPubKeyHash)(nil)
	sizeNilP2PKH, err := GetOutputSizeByAddress(nilP2PKH)
	require.NoError(t, err)
	require.Equal(t, uint64(0), sizeNilP2PKH)

	addrP2PKH := getTestAddrScript(t, ScriptTypeP2PKH)
	sizeP2PKH, err := GetOutputSizeByAddress(addrP2PKH)
	require.NoError(t, err)
	require.Equal(t, uint64(bytesPerOutputP2PKH), sizeP2PKH)

	// test unsupported address type
	nilP2PK := (*btcutil.AddressPubKey)(nil)
	sizeP2PK, err := GetOutputSizeByAddress(nilP2PK)
	require.ErrorContains(t, err, "cannot get output size for address type")
	require.Equal(t, uint64(0), sizeP2PK)
}

func TestOutputSizeP2TR(t *testing.T) {
	// Generate payer/payee private keys and P2WPKH addresss
	privateKey, _, payerScript := generateKeyPair(t, &chaincfg.TestNet3Params)
	payee := getTestAddrScript(t, ScriptTypeP2TR)

	// Create a new transaction and 10 random amount of payee scripts
	tx := wire.NewMsgTx(wire.TxVersion)
	payees, payeeScripts := createPkScripts(t, payee, 10)

	// Add inputs and outputs and sign the transaction
	addTxInputsOutputsAndSignTx(t, tx, privateKey, payerScript, exampleTxids[:2], payeeScripts)

	// Estimate the tx size in vByte
	// #nosec G115 always positive
	vBytes := uint64(blockchain.GetTransactionWeight(btcutil.NewTx(tx)) / blockchain.WitnessScaleFactor)
	vBytesEstimated, err := EstimateOutboundSize(2, payees)
	require.NoError(t, err)
	require.Equal(t, vBytes, vBytesEstimated)
}

func TestOutputSizeP2WSH(t *testing.T) {
	// Generate payer/payee private keys and P2WPKH addresss
	privateKey, _, payerScript := generateKeyPair(t, &chaincfg.TestNet3Params)
	payee := getTestAddrScript(t, ScriptTypeP2WSH)

	// Create a new transaction and 10 random amount of payee scripts
	tx := wire.NewMsgTx(wire.TxVersion)
	payees, payeeScripts := createPkScripts(t, payee, 10)

	// Add inputs and outputs and sign the transaction
	addTxInputsOutputsAndSignTx(t, tx, privateKey, payerScript, exampleTxids[:2], payeeScripts)

	// Estimate the tx size in vByte
	// #nosec G115 always positive
	vBytes := uint64(blockchain.GetTransactionWeight(btcutil.NewTx(tx)) / blockchain.WitnessScaleFactor)
	vBytesEstimated, err := EstimateOutboundSize(2, payees)
	require.NoError(t, err)
	require.Equal(t, vBytes, vBytesEstimated)
}

func TestOutputSizeP2SH(t *testing.T) {
	// Generate payer/payee private keys and P2SH addresss
	privateKey, _, payerScript := generateKeyPair(t, &chaincfg.TestNet3Params)
	payee := getTestAddrScript(t, ScriptTypeP2SH)

	// Create a new transaction and 10 random amount of payee scripts
	tx := wire.NewMsgTx(wire.TxVersion)
	payees, payeeScripts := createPkScripts(t, payee, 10)

	// Add inputs and outputs and sign the transaction
	addTxInputsOutputsAndSignTx(t, tx, privateKey, payerScript, exampleTxids[:2], payeeScripts)

	// Estimate the tx size in vByte
	// #nosec G115 always positive
	vBytes := uint64(blockchain.GetTransactionWeight(btcutil.NewTx(tx)) / blockchain.WitnessScaleFactor)
	vBytesEstimated, err := EstimateOutboundSize(2, payees)
	require.NoError(t, err)
	require.Equal(t, vBytes, vBytesEstimated)
}

func TestOutputSizeP2PKH(t *testing.T) {
	// Generate payer/payee private keys and P2PKH addresss
	privateKey, _, payerScript := generateKeyPair(t, &chaincfg.TestNet3Params)
	payee := getTestAddrScript(t, ScriptTypeP2PKH)

	// Create a new transaction and 10 random amount of payee scripts
	tx := wire.NewMsgTx(wire.TxVersion)
	payees, payeeScripts := createPkScripts(t, payee, 10)

	// Add inputs and outputs and sign the transaction
	addTxInputsOutputsAndSignTx(t, tx, privateKey, payerScript, exampleTxids[:2], payeeScripts)

	// Estimate the tx size in vByte
	// #nosec G115 always positive
	vBytes := uint64(blockchain.GetTransactionWeight(btcutil.NewTx(tx)) / blockchain.WitnessScaleFactor)
	vBytesEstimated, err := EstimateOutboundSize(2, payees)
	require.NoError(t, err)
	require.Equal(t, vBytes, vBytesEstimated)
}

func TestOutboundSizeBreakdown(t *testing.T) {
	// a list of all types of addresses
	payees := []btcutil.Address{
		getTestAddrScript(t, ScriptTypeP2TR),
		getTestAddrScript(t, ScriptTypeP2WSH),
		getTestAddrScript(t, ScriptTypeP2WPKH),
		getTestAddrScript(t, ScriptTypeP2SH),
		getTestAddrScript(t, ScriptTypeP2PKH),
	}

	// add all outbound sizes paying to each address
	txSizeTotal := uint64(0)
	for _, payee := range payees {
		sizeOutput, err := EstimateOutboundSize(2, []btcutil.Address{payee})
		require.NoError(t, err)
		txSizeTotal += sizeOutput
	}

	// calculate the average outbound size
	// #nosec G115 always in range
	txSizeAverage := uint64((float64(txSizeTotal))/float64(len(payees)) + 0.5)

	// get deposit fee
	txSizeDepositor := OutboundSizeDepositor()
	require.Equal(t, uint64(68), txSizeDepositor)

	// get withdrawer fee
	txSizeWithdrawer := OutboundSizeWithdrawer()
	require.Equal(t, uint64(177), txSizeWithdrawer)

	// total outbound size == (deposit fee + withdrawer fee), 245 = 68 + 177
	require.Equal(t, OutboundBytesAvg, txSizeAverage)
	require.Equal(t, txSizeAverage, txSizeDepositor+txSizeWithdrawer)

	// check default depositor fee
	depositFee := DepositorFee(defaultDepositorFeeRate)
	require.Equal(t, depositFee, 0.00001360)
}

func TestOutboundSizeMinMaxError(t *testing.T) {
	// P2TR output is the largest in size; P2WPKH is the smallest
	toP2TR := getTestAddrScript(t, ScriptTypeP2TR)
	toP2WPKH := getTestAddrScript(t, ScriptTypeP2WPKH)

	// Estimate the largest outbound size in vByte
	sizeMax, err := EstimateOutboundSize(21, []btcutil.Address{toP2TR})
	require.NoError(t, err)
	require.Equal(t, OutboundBytesMax, sizeMax)

	// Estimate the smallest outbound size in vByte
	sizeMin, err := EstimateOutboundSize(2, []btcutil.Address{toP2WPKH})
	require.NoError(t, err)
	require.Equal(t, OutboundBytesMin, sizeMin)

	// Estimate unknown address type
	nilP2PK := (*btcutil.AddressPubKey)(nil)
	size, err := EstimateOutboundSize(1, []btcutil.Address{nilP2PK})
	require.Error(t, err)
	require.Equal(t, uint64(0), size)
}
