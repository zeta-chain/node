package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

const (
	chunkSize = 500
)

func main() {
	cfg := MustGetConfig()

	connCfg := &rpcclient.ConnConfig{
		Host:         cfg.RPCHost,
		User:         cfg.RPCUser,
		Pass:         cfg.RPCPass,
		DisableTLS:   true,
		HTTPPostMode: true,
	}
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatalf("error creating rpc connection : %v\n", err.Error())
	}

	// get Wallet data
	wif, addr, err := getWallet(cfg.WalletPK)
	if err != nil {
		log.Fatalf("error getting wallet data: %v\n", err)
	}
	privKey := *wif.PrivKey

	// get the current block height.
	blockCount, err := client.GetBlockCount()
	if err != nil {
		log.Fatalf("error getting block height : %v\n", err)
	}
	log.Printf("Block height: %d", blockCount)

	// List unspent.
	address, err := btcutil.DecodeAddress(cfg.WalletAddress, &chaincfg.TestNet3Params)
	if err != nil {
		log.Fatalf("error decoding wallet address (%s) : %s\n", cfg.WalletAddress, err.Error())
	}
	addresses := []btcutil.Address{address}
	var utxos []btcjson.ListUnspentResult

	// first pass, populate utxos array
	for i := cfg.MinConf; i < cfg.MaxConf; i += chunkSize {
		unspents, err := client.ListUnspentMinMaxAddresses(i, i+chunkSize, addresses)
		if err != nil {
			log.Fatalf(" error pulling utxos: %v\n", err.Error())
		}
		var total float64
		for _, unspent := range unspents {
			total = total + unspent.Amount
		}
		log.Printf("fetched %d utxos, %v btc\n", len(unspents), total)
		utxos = append(utxos, unspents...)
	}

	// print totals
	var total float64
	for _, utxo := range utxos {
		total = total + utxo.Amount
	}
	fmt.Printf("totals: %d utxos, %v btc\n", len(utxos), total)

	// consolidate
	for i := 0; i < len(utxos); i += cfg.PrevCount {
		tx, err := consolidateTX(addr, privKey, utxos[i:i+cfg.PrevCount], wif.CompressPubKey)
		if err != nil {
			log.Printf("error consolidating tx : %v\n", err.Error())
			continue
		}
		printTX(tx)
	}
}

func consolidateTX(addr *btcutil.AddressWitnessPubKeyHash, privKey btcec.PrivateKey, prevOuts []btcjson.ListUnspentResult, compress bool) (*wire.MsgTx, error) {
	var total float64
	pkScript, err := payToWitnessPubKeyHashScript(addr.WitnessProgram())
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
		total = total + prevOut.Amount
	}

	totalSatoshis, err := getSatoshis(total)
	if err != nil {
		return nil, err
	}
	txOut := wire.NewTxOut(totalSatoshis, pkScript)
	tx.AddTxOut(txOut)

	return tx, nil
}

func getWallet(pk string) (*btcutil.WIF, *btcutil.AddressWitnessPubKeyHash, error) {
	wif, err := btcutil.DecodeWIF(pk)
	if err != nil {
		return nil, nil, err
	}
	addr, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(wif.SerializePubKey()), &chaincfg.TestNet3Params)
	if err != nil {
		return nil, nil, err
	}
	return wif, addr, nil
}

func printTX(tx *wire.MsgTx) {
	buf := new(bytes.Buffer)
	if err := tx.Serialize(buf); err != nil {
		log.Printf("error printing tx: %v\n", err.Error())
	}
	log.Printf("tx: %v\n", hex.EncodeToString(buf.Bytes()))
}

func payToWitnessPubKeyHashScript(pubKeyHash []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(pubKeyHash).Script()
}
