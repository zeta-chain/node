package main

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
)

func (sm *SmokeTest) TestBitcoinSetup() {
	{
		btc := sm.btcRpcClient
		_, err := btc.CreateWallet("smoketest", rpcclient.WithCreateWalletBlank())
		if err != nil {
			panic(err)
		}
		skBytes, err := hex.DecodeString(DeployerPrivateKey)
		sk, _ := btcec.PrivKeyFromBytes(btcec.S256(), skBytes)
		privkeyWIF, err := btcutil.NewWIF(sk, &chaincfg.RegressionNetParams, true)
		if err != nil {
			panic(err)
		}
		err = btc.ImportPrivKeyRescan(privkeyWIF, "deployer", true)
		if err != nil {
			panic(err)
		}
		BTCDeployerAddress, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(privkeyWIF.PrivKey.PubKey().SerializeCompressed()), &chaincfg.RegressionNetParams)
		if err != nil {
			panic(err)
		}
		fmt.Printf("BTCDeployerAddress: %s\n", BTCDeployerAddress.EncodeAddress())

		err = btc.ImportAddress(BTCTSSAddress.EncodeAddress())
		if err != nil {
			panic(err)
		}
		_, err = btc.GenerateToAddress(101, BTCDeployerAddress, nil)
		if err != nil {
			panic(err)
		}
		bal, err := btc.GetBalance("*")
		if err != nil {
			panic(err)
		}
		_, err = btc.GenerateToAddress(1, BTCTSSAddress, nil)
		if err != nil {
			panic(err)
		}
		bal, err = btc.GetBalance("*")
		if err != nil {
			panic(err)
		}
		fmt.Printf("balance: %f\n", bal.ToBTC())

		bals, err := btc.GetBalances()
		if err != nil {
			panic(err)
		}
		fmt.Printf("balances: \n")
		fmt.Printf("  mine: %+v\n", bals.Mine)
		if bals.WatchOnly != nil {
			fmt.Printf("  watchonly: %+v\n", bals.WatchOnly)
		}
		fmt.Printf("TSS Address: %s\n", BTCTSSAddress.EncodeAddress())
		utxos, err := btc.ListUnspent()
		if err != nil {
			panic(err)
		}
		for _, utxo := range utxos {
			fmt.Printf("utxo: %+v\n", utxo)
		}

		// send 1 BTC to TSS address
		input0 := btcjson.TransactionInput{utxos[0].TxID, utxos[0].Vout}
		input1 := btcjson.TransactionInput{utxos[1].TxID, utxos[1].Vout}
		inputs := []btcjson.TransactionInput{input0, input1}
		fee := btcutil.Amount(0.0001 * btcutil.SatoshiPerBitcoin)
		change := btcutil.Amount((utxos[0].Amount+utxos[1].Amount)*(btcutil.SatoshiPerBitcoin)) - fee - btcutil.Amount(1*btcutil.SatoshiPerBitcoin)
		amounts := map[btcutil.Address]btcutil.Amount{
			BTCTSSAddress:      btcutil.Amount(1 * btcutil.SatoshiPerBitcoin),
			BTCDeployerAddress: change,
		}
		tx, err := btc.CreateRawTransaction(inputs, amounts, nil)
		if err != nil {
			panic(err)
		}
		fmt.Printf("raw transaction: \n")
		for idx, txout := range tx.TxOut {
			fmt.Printf("txout %d\n", idx)
			fmt.Printf("  value: %d\n", txout.Value)
			fmt.Printf("  PkScript: %x\n", txout.PkScript)
		}
		stx, signed, err := btc.SignRawTransactionWithWallet(tx)
		if err != nil {
			panic(err)
		}
		fmt.Printf("signed tx: %+v, all inputs signed?: %+v\n", stx, signed)
		txid, err := btc.SendRawTransaction(stx, true)
		if err != nil {
			panic(err)
		}
		fmt.Printf("txid: %+v\n", txid)
		_, err = btc.GenerateToAddress(6, BTCDeployerAddress, nil)
		if err != nil {
			panic(err)
		}
		gtx, err := btc.GetTransaction(txid)
		if err != nil {
			panic(err)
		}
		fmt.Printf("rawtx: %+v\n", gtx)
	}
}
