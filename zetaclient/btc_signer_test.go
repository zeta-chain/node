package zetaclient

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/crypto"
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
	s.btcSigner, err = NewBTCSigner(&tss)
	c.Assert(err, IsNil)
}

func (s *BTCSignerSuite) Test1(c *C) {
	addr := s.btcSigner.tssSigner.BTCAddressPubkey()
	originTx, err := createFakeOriginTx(addr.AddressPubKeyHash())
	c.Assert(err, IsNil)
	_ = originTx
}

// createFakeOriginTx creates a fake coinbase transaction that is used in the
// example as a stand-in for what ordinarily be the real transaction that is
// being spent.
func createFakeOriginTx(addr btcutil.Address) (*wire.MsgTx, error) {
	tx := wire.NewMsgTx(wire.TxVersion)
	prevOut := wire.NewOutPoint(&chainhash.Hash{}, ^uint32(0))

	txIn := wire.NewTxIn(prevOut, []byte{txscript.OP_0, txscript.OP_0}, [][]byte{})
	tx.AddTxIn(txIn)

	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return nil, err
	}
	txOut := wire.NewTxOut(100000000, pkScript)
	tx.AddTxOut(txOut)
	return tx, nil
}
