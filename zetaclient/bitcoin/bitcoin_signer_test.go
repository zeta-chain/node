package bitcoin

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"reflect"
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
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	clientcommon "github.com/zeta-chain/zetacore/zetaclient/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	corecontext "github.com/zeta-chain/zetacore/zetaclient/core_context"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/stub"
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
	tss := interfaces.TestSigner{
		PrivKey: privateKey,
	}
	cfg := config.NewConfig()
	s.btcSigner, err = NewBTCSigner(
		config.BTCConfig{},
		&tss,
		clientcommon.DefaultLoggers(),
		&metrics.TelemetryServer{},
		corecontext.NewZetaCoreContext(cfg))
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
	pkScript, err := PayToAddrScript(addr)

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
	pkScript, err := PayToAddrScript(addr)
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
	pkScript, err = PayToAddrScript(addr)
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

// helper function to create a test BitcoinChainClient
func createTestClient(t *testing.T) *BTCChainClient {
	skHex := "7b8507ba117e069f4a3f456f505276084f8c92aee86ac78ae37b4d1801d35fa8"
	privateKey, err := crypto.HexToECDSA(skHex)
	require.Nil(t, err)
	tss := interfaces.TestSigner{
		PrivKey: privateKey,
	}
	return &BTCChainClient{
		Tss:               tss,
		Mu:                &sync.Mutex{},
		includedTxResults: make(map[string]*btcjson.GetTransactionResult),
	}
}

// helper function to create a test BitcoinChainClient with UTXOs
func createTestClientWithUTXOs(t *testing.T) *BTCChainClient {
	// Create BitcoinChainClient
	client := createTestClient(t)
	tssAddress := client.Tss.BTCAddressWitnessPubkeyHash().EncodeAddress()

	// Create 10 dummy UTXOs (22.44 BTC in total)
	client.utxos = make([]btcjson.ListUnspentResult, 0, 10)
	amounts := []float64{0.01, 0.12, 0.18, 0.24, 0.5, 1.26, 2.97, 3.28, 5.16, 8.72}
	for _, amount := range amounts {
		client.utxos = append(client.utxos, btcjson.ListUnspentResult{Address: tssAddress, Amount: amount})
	}
	return client
}

func TestAddWithdrawTxOutputs(t *testing.T) {
	// Create test signer and receiver address
	signer, err := NewBTCSigner(config.BTCConfig{}, stub.NewTSSMainnet(), clientcommon.DefaultLoggers(), &metrics.TelemetryServer{}, nil)
	require.NoError(t, err)

	// tss address and script
	tssAddr := signer.tssSigner.BTCAddressWitnessPubkeyHash()
	tssScript, err := PayToAddrScript(tssAddr)
	require.NoError(t, err)
	fmt.Printf("tss address: %s", tssAddr.EncodeAddress())

	// receiver addresses
	receiver := "bc1qaxf82vyzy8y80v000e7t64gpten7gawewzu42y"
	to, err := chains.DecodeBtcAddress(receiver, chains.BtcMainnetChain.ChainId)
	require.NoError(t, err)
	toScript, err := PayToAddrScript(to)
	require.NoError(t, err)

	// test cases
	tests := []struct {
		name     string
		tx       *wire.MsgTx
		to       btcutil.Address
		total    float64
		amount   float64
		nonce    int64
		fees     *big.Int
		cancelTx bool
		fail     bool
		message  string
		txout    []*wire.TxOut
	}{
		{
			name:   "should add outputs successfully",
			tx:     wire.NewMsgTx(wire.TxVersion),
			to:     to,
			total:  1.00012000,
			amount: 0.2,
			nonce:  10000,
			fees:   big.NewInt(2000),
			fail:   false,
			txout: []*wire.TxOut{
				{Value: 10000, PkScript: tssScript},
				{Value: 20000000, PkScript: toScript},
				{Value: 80000000, PkScript: tssScript},
			},
		},
		{
			name:   "should add outputs without change successfully",
			tx:     wire.NewMsgTx(wire.TxVersion),
			to:     to,
			total:  0.20012000,
			amount: 0.2,
			nonce:  10000,
			fees:   big.NewInt(2000),
			fail:   false,
			txout: []*wire.TxOut{
				{Value: 10000, PkScript: tssScript},
				{Value: 20000000, PkScript: toScript},
			},
		},
		{
			name:     "should cancel tx successfully",
			tx:       wire.NewMsgTx(wire.TxVersion),
			to:       to,
			total:    1.00012000,
			amount:   0.2,
			nonce:    10000,
			fees:     big.NewInt(2000),
			cancelTx: true,
			fail:     false,
			txout: []*wire.TxOut{
				{Value: 10000, PkScript: tssScript},
				{Value: 100000000, PkScript: tssScript},
			},
		},
		{
			name:   "should fail on invalid amount",
			tx:     wire.NewMsgTx(wire.TxVersion),
			to:     to,
			total:  1.00012000,
			amount: -0.5,
			fail:   true,
		},
		{
			name:   "should fail when total < amount",
			tx:     wire.NewMsgTx(wire.TxVersion),
			to:     to,
			total:  0.00012000,
			amount: 0.2,
			fail:   true,
		},
		{
			name:    "should fail when total < fees + amount + nonce",
			tx:      wire.NewMsgTx(wire.TxVersion),
			to:      to,
			total:   0.20011000,
			amount:  0.2,
			nonce:   10000,
			fees:    big.NewInt(2000),
			fail:    true,
			message: "remainder value is negative",
		},
		{
			name:   "should not produce duplicate nonce mark",
			tx:     wire.NewMsgTx(wire.TxVersion),
			to:     to,
			total:  0.20022000, //  0.2 + fee + nonceMark * 2
			amount: 0.2,
			nonce:  10000,
			fees:   big.NewInt(2000),
			fail:   false,
			txout: []*wire.TxOut{
				{Value: 10000, PkScript: tssScript},
				{Value: 20000000, PkScript: toScript},
				{Value: 9999, PkScript: tssScript}, // nonceMark - 1
			},
		},
		{
			name:   "should fail on invalid to address",
			tx:     wire.NewMsgTx(wire.TxVersion),
			to:     nil,
			total:  1.00012000,
			amount: 0.2,
			nonce:  10000,
			fees:   big.NewInt(2000),
			fail:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := signer.AddWithdrawTxOutputs(tt.tx, tt.to, tt.total, tt.amount, tt.nonce, tt.fees, tt.cancelTx)
			if tt.fail {
				require.Error(t, err)
				if tt.message != "" {
					require.Contains(t, err.Error(), tt.message)
				}
				return
			} else {
				require.NoError(t, err)
				require.True(t, reflect.DeepEqual(tt.txout, tt.tx.TxOut))
			}
		})
	}
}

func mineTxNSetNonceMark(ob *BTCChainClient, nonce uint64, txid string, preMarkIndex int) {
	// Mine transaction
	outboundID := ob.GetTxID(nonce)
	ob.includedTxResults[outboundID] = &btcjson.GetTransactionResult{TxID: txid}

	// Set nonce mark
	tssAddress := ob.Tss.BTCAddressWitnessPubkeyHash().EncodeAddress()
	nonceMark := btcjson.ListUnspentResult{TxID: txid, Address: tssAddress, Amount: float64(chains.NonceMarkAmount(nonce)) * 1e-8}
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
	ob := createTestClientWithUTXOs(t)
	dummyTxID := "6e6f71d281146c1fc5c755b35908ee449f26786c84e2ae18f98b268de40b7ec4"

	// Case1: nonce = 0, bootstrap
	// 		input: utxoCap = 5, amount = 0.01, nonce = 0
	// 		output: [0.01], 0.01
	result, amount, _, _, err := ob.SelectUTXOs(0.01, 5, 0, math.MaxUint16, true)
	require.Nil(t, err)
	require.Equal(t, 0.01, amount)
	require.Equal(t, ob.utxos[0:1], result)

	// Case2: nonce = 1, must FAIL and wait for previous transaction to be mined
	// 		input: utxoCap = 5, amount = 0.5, nonce = 1
	// 		output: error
	result, amount, _, _, err = ob.SelectUTXOs(0.5, 5, 1, math.MaxUint16, true)
	require.NotNil(t, err)
	require.Nil(t, result)
	require.Zero(t, amount)
	require.Equal(t, "getOutboundIDByNonce: cannot find outTx txid for nonce 0", err.Error())
	mineTxNSetNonceMark(ob, 0, dummyTxID, -1) // mine a transaction and set nonce-mark utxo for nonce 0

	// Case3: nonce = 1, should pass now
	// 		input: utxoCap = 5, amount = 0.5, nonce = 1
	// 		output: [0.00002, 0.01, 0.12, 0.18, 0.24], 0.55002
	result, amount, _, _, err = ob.SelectUTXOs(0.5, 5, 1, math.MaxUint16, true)
	require.Nil(t, err)
	require.Equal(t, 0.55002, amount)
	require.Equal(t, ob.utxos[0:5], result)
	mineTxNSetNonceMark(ob, 1, dummyTxID, 0) // mine a transaction and set nonce-mark utxo for nonce 1

	// Case4:
	// 		input: utxoCap = 5, amount = 1.0, nonce = 2
	// 		output: [0.00002001, 0.01, 0.12, 0.18, 0.24, 0.5], 1.05002001
	result, amount, _, _, err = ob.SelectUTXOs(1.0, 5, 2, math.MaxUint16, true)
	require.Nil(t, err)
	require.InEpsilon(t, 1.05002001, amount, 1e-8)
	require.Equal(t, ob.utxos[0:6], result)
	mineTxNSetNonceMark(ob, 2, dummyTxID, 0) // mine a transaction and set nonce-mark utxo for nonce 2

	// Case5: should include nonce-mark utxo on the LEFT
	// 		input: utxoCap = 5, amount = 8.05, nonce = 3
	// 		output: [0.00002002, 0.24, 0.5, 1.26, 2.97, 3.28], 8.25002002
	result, amount, _, _, err = ob.SelectUTXOs(8.05, 5, 3, math.MaxUint16, true)
	require.Nil(t, err)
	require.InEpsilon(t, 8.25002002, amount, 1e-8)
	expected := append([]btcjson.ListUnspentResult{ob.utxos[0]}, ob.utxos[4:9]...)
	require.Equal(t, expected, result)
	mineTxNSetNonceMark(ob, 24105431, dummyTxID, 0) // mine a transaction and set nonce-mark utxo for nonce 24105431

	// Case6: should include nonce-mark utxo on the RIGHT
	// 		input: utxoCap = 5, amount = 0.503, nonce = 24105432
	// 		output: [0.24107432, 0.01, 0.12, 0.18, 0.24], 0.55002002
	result, amount, _, _, err = ob.SelectUTXOs(0.503, 5, 24105432, math.MaxUint16, true)
	require.Nil(t, err)
	require.InEpsilon(t, 0.79107431, amount, 1e-8)
	expected = append([]btcjson.ListUnspentResult{ob.utxos[4]}, ob.utxos[0:4]...)
	require.Equal(t, expected, result)
	mineTxNSetNonceMark(ob, 24105432, dummyTxID, 4) // mine a transaction and set nonce-mark utxo for nonce 24105432

	// Case7: should include nonce-mark utxo in the MIDDLE
	// 		input: utxoCap = 5, amount = 1.0, nonce = 24105433
	// 		output: [0.24107432, 0.12, 0.18, 0.24, 0.5], 1.28107432
	result, amount, _, _, err = ob.SelectUTXOs(1.0, 5, 24105433, math.MaxUint16, true)
	require.Nil(t, err)
	require.InEpsilon(t, 1.28107432, amount, 1e-8)
	expected = append([]btcjson.ListUnspentResult{ob.utxos[4]}, ob.utxos[1:4]...)
	expected = append(expected, ob.utxos[5])
	require.Equal(t, expected, result)

	// Case8: should work with maximum amount
	// 		input: utxoCap = 5, amount = 16.03
	// 		output: [0.24107432, 1.26, 2.97, 3.28, 5.16, 8.72], 21.63107432
	result, amount, _, _, err = ob.SelectUTXOs(16.03, 5, 24105433, math.MaxUint16, true)
	require.Nil(t, err)
	require.InEpsilon(t, 21.63107432, amount, 1e-8)
	expected = append([]btcjson.ListUnspentResult{ob.utxos[4]}, ob.utxos[6:11]...)
	require.Equal(t, expected, result)

	// Case9: must FAIL due to insufficient funds
	// 		input: utxoCap = 5, amount = 21.64
	// 		output: error
	result, amount, _, _, err = ob.SelectUTXOs(21.64, 5, 24105433, math.MaxUint16, true)
	require.NotNil(t, err)
	require.Nil(t, result)
	require.Zero(t, amount)
	require.Equal(t, "SelectUTXOs: not enough btc in reserve - available : 21.63107432 , tx amount : 21.64", err.Error())
}

func TestUTXOConsolidation(t *testing.T) {
	dummyTxID := "6e6f71d281146c1fc5c755b35908ee449f26786c84e2ae18f98b268de40b7ec4"

	t.Run("should not consolidate", func(t *testing.T) {
		ob := createTestClientWithUTXOs(t)
		mineTxNSetNonceMark(ob, 0, dummyTxID, -1) // mine a transaction and set nonce-mark utxo for nonce 0

		// input: utxoCap = 10, amount = 0.01, nonce = 1, rank = 10
		// output: [0.00002, 0.01], 0.01002
		result, amount, clsdtUtxo, clsdtValue, err := ob.SelectUTXOs(0.01, 10, 1, 10, true)
		require.Nil(t, err)
		require.Equal(t, 0.01002, amount)
		require.Equal(t, ob.utxos[0:2], result)
		require.Equal(t, uint16(0), clsdtUtxo)
		require.Equal(t, 0.0, clsdtValue)
	})

	t.Run("should consolidate 1 utxo", func(t *testing.T) {
		ob := createTestClientWithUTXOs(t)
		mineTxNSetNonceMark(ob, 0, dummyTxID, -1) // mine a transaction and set nonce-mark utxo for nonce 0

		// input: utxoCap = 9, amount = 0.01, nonce = 1, rank = 9
		// output: [0.00002, 0.01, 0.12], 0.13002
		result, amount, clsdtUtxo, clsdtValue, err := ob.SelectUTXOs(0.01, 9, 1, 9, true)
		require.Nil(t, err)
		require.Equal(t, 0.13002, amount)
		require.Equal(t, ob.utxos[0:3], result)
		require.Equal(t, uint16(1), clsdtUtxo)
		require.Equal(t, 0.12, clsdtValue)
	})

	t.Run("should consolidate 3 utxos", func(t *testing.T) {
		ob := createTestClientWithUTXOs(t)
		mineTxNSetNonceMark(ob, 0, dummyTxID, -1) // mine a transaction and set nonce-mark utxo for nonce 0

		// input: utxoCap = 5, amount = 0.01, nonce = 0, rank = 5
		// output: [0.00002, 0.014, 1.26, 0.5, 0.2], 2.01002
		result, amount, clsdtUtxo, clsdtValue, err := ob.SelectUTXOs(0.01, 5, 1, 5, true)
		require.Nil(t, err)
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
		ob := createTestClientWithUTXOs(t)
		mineTxNSetNonceMark(ob, 0, dummyTxID, -1) // mine a transaction and set nonce-mark utxo for nonce 0

		// input: utxoCap = 12, amount = 0.01, nonce = 0, rank = 1
		// output: [0.00002, 0.01, 8.72, 5.16, 3.28, 2.97, 1.26, 0.5, 0.24, 0.18, 0.12], 22.44002
		result, amount, clsdtUtxo, clsdtValue, err := ob.SelectUTXOs(0.01, 12, 1, 1, true)
		require.Nil(t, err)
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
		ob := createTestClientWithUTXOs(t)
		mineTxNSetNonceMark(ob, 24105431, dummyTxID, -1) // mine a transaction and set nonce-mark utxo for nonce 24105431

		// input: utxoCap = 5, amount = 0.13, nonce = 24105432, rank = 5
		// output: [0.24107431, 0.01, 0.12, 1.26, 0.5, 0.24], 2.37107431
		result, amount, clsdtUtxo, clsdtValue, err := ob.SelectUTXOs(0.13, 5, 24105432, 5, true)
		require.Nil(t, err)
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
		ob := createTestClientWithUTXOs(t)
		mineTxNSetNonceMark(ob, 24105431, dummyTxID, -1) // mine a transaction and set nonce-mark utxo for nonce 24105431

		// input: utxoCap = 12, amount = 0.13, nonce = 24105432, rank = 1
		// output: [0.24107431, 0.01, 0.12, 8.72, 5.16, 3.28, 2.97, 1.26, 0.5, 0.24, 0.18], 22.68107431
		result, amount, clsdtUtxo, clsdtValue, err := ob.SelectUTXOs(0.13, 12, 24105432, 1, true)
		require.Nil(t, err)
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

// Coverage doesn't seem to pick this up from the suite
func TestNewBTCSigner(t *testing.T) {
	// test private key with EVM address
	//// EVM: 0x236C7f53a90493Bb423411fe4117Cb4c2De71DfB
	// BTC testnet3: muGe9prUBjQwEnX19zG26fVRHNi8z7kSPo
	skHex := "7b8507ba117e069f4a3f456f505276084f8c92aee86ac78ae37b4d1801d35fa8"
	privateKey, err := crypto.HexToECDSA(skHex)
	require.NoError(t, err)
	tss := interfaces.TestSigner{
		PrivKey: privateKey,
	}
	cfg := config.NewConfig()
	btcSigner, err := NewBTCSigner(
		config.BTCConfig{},
		&tss,
		clientcommon.DefaultLoggers(),
		&metrics.TelemetryServer{},
		corecontext.NewZetaCoreContext(cfg))
	require.NoError(t, err)
	require.NotNil(t, btcSigner)
}
