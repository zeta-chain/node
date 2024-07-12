package signer

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/stretchr/testify/suite"

	"github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
)

type BTCSignTestSuite struct {
	suite.Suite
	testSigner interfaces.TSSSigner
}

const (
	prevOut = "07a84f4bd45a633e93871be5c98d958afd13a37f3cf5010f40eec0840d19f5fa"
	pk      = "cQkjdfeMU8vHvE6jErnFVqZYYZnGGYy64jH6zovbSXdfTjte6QgY"
)

func (suite *BTCSignTestSuite) SetupTest() {
	wif, _ := btcutil.DecodeWIF(pk)
	privateKey := wif.PrivKey

	suite.testSigner = &mocks.TSS{ // fake TSS
		PrivKey: privateKey.ToECDSA(),
	}
	addr := suite.testSigner.BTCAddressWitnessPubkeyHash()
	suite.T().Logf("segwit addr: %s", addr)
}

func (suite *BTCSignTestSuite) TearDownSuite() {
}

func (suite *BTCSignTestSuite) TestSign() {

	// build a tx used for both signatures
	tx, txSigHashes, idx, amt, subscript, privKey, compress, err := buildTX()
	suite.Require().NoError(err)

	// sign tx using wallet signature
	walletSignedTX, err := getWalletTX(tx, txSigHashes, idx, amt, subscript, txscript.SigHashAll, privKey, compress)
	suite.Require().NoError(err)
	suite.T().Logf("wallet signed tx : %v\n", walletSignedTX)

	// sign tx using tss signature
	tssSignedTX, err := getTSSTX(suite.testSigner, tx, txSigHashes, idx, amt, subscript, txscript.SigHashAll)
	suite.Require().NoError(err)
	suite.T().Logf("tss signed tx :    %v\n", tssSignedTX)
}

func (suite *BTCSignTestSuite) TestSubmittedTx() {

}

func TestBTCSign(t *testing.T) {
	suite.Run(t, new(BTCSignTestSuite))
}

func buildTX() (*wire.MsgTx, *txscript.TxSigHashes, int, int64, []byte, *btcec.PrivateKey, bool, error) {
	wif, err := btcutil.DecodeWIF(pk)
	if err != nil {
		return nil, nil, 0, 0, nil, nil, false, err
	}
	fmt.Printf("is wif.netid testnet3 %v\n", wif.IsForNet(&chaincfg.TestNet3Params))
	fmt.Printf("is wif pubkey compressed %v\n", wif.CompressPubKey)

	addr, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &chaincfg.TestNet3Params)
	if err != nil {
		return nil, nil, 0, 0, nil, nil, false, err
	}

	hash, err := chainhash.NewHashFromStr(prevOut)
	if err != nil {
		return nil, nil, 0, 0, nil, nil, false, err
	}
	outpoint := wire.NewOutPoint(hash, 0)

	// build tx
	tx := wire.NewMsgTx(wire.TxVersion)
	txIn := wire.NewTxIn(outpoint, nil, nil)
	tx.AddTxIn(txIn)

	pkScript, err := bitcoin.PayToAddrScript(addr)
	if err != nil {
		return nil, nil, 0, 0, nil, nil, false, err
	}
	txOut := wire.NewTxOut(47000, pkScript)
	tx.AddTxOut(txOut)

	txSigHashes := txscript.NewTxSigHashes(tx)

	privKey := btcec.PrivateKey(*wif.PrivKey.ToECDSA())

	return tx, txSigHashes, int(0), int64(65236), pkScript, &privKey, wif.CompressPubKey, nil
}

func getWalletTX(
	tx *wire.MsgTx,
	sigHashes *txscript.TxSigHashes,
	idx int,
	amt int64,
	subscript []byte,
	hashType txscript.SigHashType,
	privKey *btcec.PrivateKey,
	compress bool,
) (string, error) {
	txWitness, err := txscript.WitnessSignature(tx, sigHashes, idx, amt, subscript, hashType, privKey, compress)
	if err != nil {
		return "", err
	}

	tx.TxIn[0].Witness = txWitness

	buf := new(bytes.Buffer)
	if err := tx.Serialize(buf); err != nil {
		return "", err
	}
	walletTx := hex.EncodeToString(buf.Bytes())
	return walletTx, nil
}

func getTSSTX(
	tss interfaces.TSSSigner,
	tx *wire.MsgTx,
	sigHashes *txscript.TxSigHashes,
	idx int,
	amt int64,
	subscript []byte,
	hashType txscript.SigHashType,
) (string, error) {
	ctx := context.Background()

	witnessHash, err := txscript.CalcWitnessSigHash(subscript, sigHashes, txscript.SigHashAll, tx, idx, amt)
	if err != nil {
		return "", err
	}

	sig65B, err := tss.Sign(ctx, witnessHash, 10, 10, 0, "")
	R := big.NewInt(0).SetBytes(sig65B[:32])
	S := big.NewInt(0).SetBytes(sig65B[32:64])
	sig := btcec.Signature{
		R: R,
		S: S,
	}
	if err != nil {
		return "", err
	}

	pkCompressed := tss.PubKeyCompressedBytes()
	txWitness := wire.TxWitness{append(sig.Serialize(), byte(hashType)), pkCompressed}
	tx.TxIn[0].Witness = txWitness

	buf := new(bytes.Buffer)
	err = tx.Serialize(buf)
	if err != nil {
		return "", err
	}

	tssTX := hex.EncodeToString(buf.Bytes())
	return tssTX, nil
}
