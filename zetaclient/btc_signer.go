package zetaclient

import (
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
)

type BTCSigner struct {
	tssSigner *TestSigner
	logger    zerolog.Logger
}

func NewBTCSigner(tssSigner *TestSigner) (*BTCSigner, error) {
	return &BTCSigner{
		tssSigner: tssSigner,
		logger:    log.With().Str("module", "BTCSigner").Logger(),
	}, nil
}

// SignWithdrawTx receives utxos sorted by value
func (signer *BTCSigner) SignWithdrawTx(to *btcutil.AddressWitnessPubKeyHash, amountBTC float64, feeRateBTCPerKB float64, utxos []btcjson.ListUnspentResult, pendingUtxos *leveldb.DB) (*wire.MsgTx, error) {
	var total float64
	var prevOuts []btcjson.ListUnspentResult
	// select N utxo sufficient to cover the amount
	for _, utxo := range utxos {
		total = total + utxo.Amount
		unspents = append(prevOuts, utxo)
		if total >= amount {
			break
		}
	}
	if total < amount {
		return nil, fmt.Errorf("not enough BTC in reserve - available : %v , tx amount : %v\n", total, amount)
	}
	// build tx with selected unspents
	pkScript, err := payToWitnessPubKeyHashScript(to.WitnessProgram())
	if err != nil {
		return nil, err
	}
	tx := wire.NewMsgTx(wire.TxVersion)
	for ix, prevOut := range prevOuts {
		hash, err := chainhash.NewHashFromStr(prevOut.TxID)
		if err != nil {
			return nil, err
		}
		outpoint := wire.NewOutPoint(hash, prevOut.Vout)
		txIn := wire.NewTxIn(outpoint, nil, nil)
		tx.AddTxIn(txIn)

		sigHashes := txscript.NewTxSigHashes(tx)
		satoshis, err := getSatoshis(prevOut.Amount)
		if err != nil {
			return nil, err
		}
		txWitness, err := txscript.WitnessSignature(tx, sigHashes, ix, satoshis, pkScript, txscript.SigHashAll, &privKey, compress)
		if err != nil {
			return nil, err
		}
		tx.TxIn[ix].Witness = txWitness
	}

	amountSatoshis, err := getSatoshis(Amount)
	if err != nil {
		return nil, err
	}
	// add txout
	txOut := wire.NewTxOut(amountSatoshis, pkScript)
	// get fees
	fees, err := getSatoshis(float64(tx.SerializeSize()) * feeRate / 1024)
	if err != nil {
		return nil, err
	}
	txOut.Value = amountSatoshis - fees
	tx.AddTxOut(txOut)

	// sign the tx
	for ix := range tx.TxIn {
		amt, err := getSatoshis(prevOuts[ix].Amount)
		if err != nil {
			return nil, err
		}
		witnessHash, err := txscript.CalcWitnessSigHash(subscript, sigHashes, txscript.SigHashAll, tx, ix, amt)
		if err != nil {
			return nil, err
		}

		sig65B, err := tss.Sign(witnessHash)
		R := big.NewInt(0).SetBytes(sig65B[:32])
		S := big.NewInt(0).SetBytes(sig65B[32:64])
		sig := btcec.Signature{
			R: R,
			S: S,
		}
		if err != nil {
			return nil, err
		}

		pkCompressed := tss.PubKeyCompressedBytes()
		txWitness := wire.TxWitness{append(sig.Serialize(), byte(hashType)), pkCompressed}
		tx.TxIn[ix].Witness = txWitness
	}

	return tx, nil
}
