package signer

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	. "gopkg.in/check.v1"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	"github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
)

type BTCSignerSuite struct {
	btcSigner *Signer
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
	tss := &mocks.TSS{
		PrivKey: privateKey,
	}
	s.btcSigner, err = NewSigner(
		chains.Chain{},
		tss,
		nil,
		base.DefaultLogger(),
		config.BTCConfig{},
	)
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
	pkScript, err := bitcoin.PayToAddrScript(addr)

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
	pkScript, err := bitcoin.PayToAddrScript(addr)
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
	pkScript, err = bitcoin.PayToAddrScript(addr)
	c.Assert(err, IsNil)

	{
		txWitness, err := txscript.WitnessSignature(
			redeemTx,
			txSigHashes,
			0,
			100000000,
			pkScript,
			txscript.SigHashAll,
			privKey,
			true,
		)
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
		witnessHash, err := txscript.CalcWitnessSigHash(
			pkScript,
			txSigHashes,
			txscript.SigHashAll,
			redeemTx,
			0,
			100000000,
		)
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

func TestAddWithdrawTxOutputs(t *testing.T) {
	// Create test signer and receiver address
	signer, err := NewSigner(
		chains.Chain{},
		mocks.NewTSSMainnet(),
		nil,
		base.DefaultLogger(),
		config.BTCConfig{},
	)
	require.NoError(t, err)

	// tss address and script
	tssAddr := signer.TSS().BTCAddressWitnessPubkeyHash()
	tssScript, err := bitcoin.PayToAddrScript(tssAddr)
	require.NoError(t, err)
	fmt.Printf("tss address: %s", tssAddr.EncodeAddress())

	// receiver addresses
	receiver := "bc1qaxf82vyzy8y80v000e7t64gpten7gawewzu42y"
	to, err := chains.DecodeBtcAddress(receiver, chains.BitcoinMainnet.ChainId)
	require.NoError(t, err)
	toScript, err := bitcoin.PayToAddrScript(to)
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
					require.ErrorContains(t, err, tt.message)
				}
				return
			} else {
				require.NoError(t, err)
				require.True(t, reflect.DeepEqual(tt.txout, tt.tx.TxOut))
			}
		})
	}
}

// Coverage doesn't seem to pick this up from the suite
func TestNewBTCSigner(t *testing.T) {
	// test private key with EVM address
	//// EVM: 0x236C7f53a90493Bb423411fe4117Cb4c2De71DfB
	// BTC testnet3: muGe9prUBjQwEnX19zG26fVRHNi8z7kSPo
	skHex := "7b8507ba117e069f4a3f456f505276084f8c92aee86ac78ae37b4d1801d35fa8"
	privateKey, err := crypto.HexToECDSA(skHex)
	require.NoError(t, err)
	tss := &mocks.TSS{
		PrivKey: privateKey,
	}
	btcSigner, err := NewSigner(
		chains.Chain{},
		tss,
		nil,
		base.DefaultLogger(),
		config.BTCConfig{})
	require.NoError(t, err)
	require.NotNil(t, btcSigner)
}
