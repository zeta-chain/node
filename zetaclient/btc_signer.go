package zetaclient

import (
	"encoding/json"
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

// SignWithdrawTx receives utxos sorted by value, amount in BTC, feeRate in BTC per Kb
func (signer *BTCSigner) SignWithdrawTx(to *btcutil.AddressWitnessPubKeyHash, amount float64, feeRate float64, utxos []btcjson.ListUnspentResult, pendingUTXOs *leveldb.DB) (*wire.MsgTx, error) {
	var total float64
	var prevOuts []btcjson.ListUnspentResult
	// select N utxo sufficient to cover the amount
	for _, utxo := range utxos {
		// check for pending utxBos
		if _, err := pendingUTXOs.Get([]byte(utxoKey(utxo)), nil); err != nil {
			if err == leveldb.ErrNotFound {
				total = total + utxo.Amount
				prevOuts = append(prevOuts, utxo)
				if total >= amount {
					break
				}
			} else {
				return nil, err
			}
		}
	}
	if total < amount {
		return nil, fmt.Errorf("not enough btc in reserve - available : %v , tx amount : %v", total, amount)
	}
	// build tx with selected unspents
	pkScript, err := payToWitnessPubKeyHashScript(to.WitnessProgram())
	if err != nil {
		return nil, err
	}
	tx := wire.NewMsgTx(wire.TxVersion)
	for _, prevOut := range prevOuts {
		hash, err := chainhash.NewHashFromStr(prevOut.TxID)
		if err != nil {
			return nil, err
		}
		outpoint := wire.NewOutPoint(hash, prevOut.Vout)
		txIn := wire.NewTxIn(outpoint, nil, nil)
		tx.AddTxIn(txIn)
	}

	amountSatoshis, err := getSatoshis(amount)
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
	sigHashes := txscript.NewTxSigHashes(tx)
	for ix := range tx.TxIn {
		amt, err := getSatoshis(prevOuts[ix].Amount)
		if err != nil {
			return nil, err
		}
		witnessHash, err := txscript.CalcWitnessSigHash(pkScript, sigHashes, txscript.SigHashAll, tx, ix, amt)
		if err != nil {
			return nil, err
		}

		sig65B, err := signer.tssSigner.Sign(witnessHash)
		R := big.NewInt(0).SetBytes(sig65B[:32])
		S := big.NewInt(0).SetBytes(sig65B[32:64])
		sig := btcec.Signature{
			R: R,
			S: S,
		}
		if err != nil {
			return nil, err
		}

		pkCompressed := signer.tssSigner.PubKeyCompressedBytes()
		hashType := txscript.SigHashAll
		txWitness := wire.TxWitness{append(sig.Serialize(), byte(hashType)), pkCompressed}
		tx.TxIn[ix].Witness = txWitness
	}

	// update pending utxos db
	err = signer.updatePendingUTXOs(pendingUTXOs, prevOuts)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (signer *BTCSigner) updatePendingUTXOs(pendingDB *leveldb.DB, utxos []btcjson.ListUnspentResult) error {
	for _, utxo := range utxos {
		bytes, err := json.Marshal(utxo)
		if err != nil {
			return err
		}
		err = pendingDB.Put([]byte(utxoKey(utxo)), bytes, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
