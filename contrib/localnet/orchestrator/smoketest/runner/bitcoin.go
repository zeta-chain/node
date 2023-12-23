package runner

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/bitcoin"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient"
)

// DepositBTC deposits BTC on ZetaChain
func (sm *SmokeTestRunner) DepositBTC() {
	sm.Logger.Print("⏳ depositing BTC into ZEVM")
	startTime := time.Now()
	defer func() {
		sm.Logger.Print("✅ BTC deposited in %s", time.Since(startTime))
	}()

	btc := sm.BtcRPCClient
	utxos, err := sm.BtcRPCClient.ListUnspent()
	if err != nil {
		panic(err)
	}
	spendableAmount := 0.0
	spendableUTXOs := 0
	for _, utxo := range utxos {
		if utxo.Spendable {
			spendableAmount += utxo.Amount
			spendableUTXOs++
		}
	}
	sm.Logger.Info("ListUnspent:")
	sm.Logger.Info("  spendableAmount: %f", spendableAmount)
	sm.Logger.Info("  spendableUTXOs: %d", spendableUTXOs)
	sm.Logger.Info("Now sending two txs to TSS address...")
	amount1 := 1.1 + zetaclient.BtcDepositorFeeMin
	txHash1, err := sm.SendToTSSFromDeployerToDeposit(sm.BTCTSSAddress, amount1, utxos[:2], btc, sm.BTCDeployerAddress)
	if err != nil {
		panic(err)
	}
	amount2 := 0.05 + zetaclient.BtcDepositorFeeMin
	txHash2, err := sm.SendToTSSFromDeployerToDeposit(sm.BTCTSSAddress, amount2, utxos[2:4], btc, sm.BTCDeployerAddress)
	if err != nil {
		panic(err)
	}
	_, err = sm.SendToTSSFromDeployerWithMemo(sm.BTCTSSAddress, 0.11, utxos[4:5], btc, []byte(zetaclient.DonationMessage), sm.BTCDeployerAddress)
	if err != nil {
		panic(err)
	}

	sm.Logger.Info("testing if the deposit into BTC ZRC20 is successful...")

	initialBalance, err := sm.BTCZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	for {
		time.Sleep(2 * time.Second)
		balance, err := sm.BTCZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
		if err != nil {
			panic(err)
		}
		diff := big.NewInt(0)
		diff.Sub(balance, initialBalance)
		sm.Logger.Info("BTC Difference in balance: %d", diff.Uint64())
		if diff.Cmp(big.NewInt(1.15*btcutil.SatoshiPerBitcoin)) != 0 {
			sm.Logger.Info("waiting for BTC balance to show up in ZRC contract... current bal %d", balance)
		} else {
			sm.Logger.Info("BTC balance is in ZRC contract! Success")
			break
		}
	}

	// prove the two transactions of the deposit
	sm.Logger.InfoLoud("Bitcoin Merkle Proof\n")

	_ = txHash1
	_ = txHash2
	//sm.ProveBTCTransaction(txHash1)
	//sm.ProveBTCTransaction(txHash2)
}

func (sm *SmokeTestRunner) ProveBTCTransaction(txHash *chainhash.Hash) {
	// get tx result
	btc := sm.BtcRPCClient
	txResult, err := btc.GetTransaction(txHash)
	if err != nil {
		panic("should get outTx result")
	}
	if txResult.Confirmations <= 0 {
		panic("outTx should have already confirmed")
	}
	txBytes, err := hex.DecodeString(txResult.Hex)
	if err != nil {
		panic(err)
	}

	// get the block with verbose transactions
	blockHash, err := chainhash.NewHashFromStr(txResult.BlockHash)
	if err != nil {
		panic(err)
	}
	blockVerbose, err := btc.GetBlockVerboseTx(blockHash)
	if err != nil {
		panic("should get block verbose tx")
	}

	// get the block header
	header, err := btc.GetBlockHeader(blockHash)
	if err != nil {
		panic("should get block header")
	}

	// collect all the txs in the block
	txns := []*btcutil.Tx{}
	for _, res := range blockVerbose.Tx {
		txBytes, err := hex.DecodeString(res.Hex)
		if err != nil {
			panic(err)
		}
		tx, err := btcutil.NewTxFromBytes(txBytes)
		if err != nil {
			panic(err)
		}
		txns = append(txns, tx)
	}

	// build merkle proof
	mk := bitcoin.NewMerkle(txns)
	path, index, err := mk.BuildMerkleProof(int(txResult.BlockIndex))
	if err != nil {
		panic("should build merkle proof")
	}

	// verify merkle proof statically
	pass := bitcoin.Prove(*txHash, header.MerkleRoot, path, index)
	if !pass {
		panic("should verify merkle proof")
	}

	hash := header.BlockHash()
	for {
		_, err := sm.ObserverClient.GetBlockHeaderByHash(context.Background(), &observertypes.QueryGetBlockHeaderByHashRequest{
			BlockHash: hash.CloneBytes(),
		})
		if err != nil {
			sm.Logger.Info("waiting for block header to show up in observer... current hash %s; err %s", hash.String(), err.Error())
		}
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	// verify merkle proof through RPC
	res, err := sm.ObserverClient.Prove(context.Background(), &observertypes.QueryProveRequest{
		ChainId:   common.BtcRegtestChain().ChainId,
		TxHash:    txHash.String(),
		BlockHash: blockHash.String(),
		Proof:     common.NewBitcoinProof(txBytes, path, index),
		TxIndex:   0, // bitcoin doesn't use txIndex
	})
	if err != nil {
		panic(err)
	}
	if !res.Valid {
		panic("txProof should be valid")
	}
	sm.Logger.Info("OK: txProof verified for inTx: %s", txHash.String())
}

func (sm *SmokeTestRunner) SendToTSSFromDeployerToDeposit(
	to btcutil.Address,
	amount float64,
	inputUTXOs []btcjson.ListUnspentResult,
	btc *rpcclient.Client,
	btcDeployerAddress *btcutil.AddressWitnessPubKeyHash,
) (*chainhash.Hash, error) {
	return sm.SendToTSSFromDeployerWithMemo(to, amount, inputUTXOs, btc, sm.DeployerAddress.Bytes(), btcDeployerAddress)
}

func (sm *SmokeTestRunner) SendToTSSFromDeployerWithMemo(
	to btcutil.Address,
	amount float64,
	inputUTXOs []btcjson.ListUnspentResult,
	btc *rpcclient.Client,
	memo []byte,
	btcDeployerAddress *btcutil.AddressWitnessPubKeyHash,
) (*chainhash.Hash, error) {
	utxos := inputUTXOs

	inputs := make([]btcjson.TransactionInput, len(utxos))
	inputSats := btcutil.Amount(0)
	amounts := make([]float64, len(utxos))
	scriptPubkeys := make([]string, len(utxos))
	for i, utxo := range utxos {
		inputs[i] = btcjson.TransactionInput{utxo.TxID, utxo.Vout}
		inputSats += btcutil.Amount(utxo.Amount * btcutil.SatoshiPerBitcoin)
		amounts[i] = utxo.Amount
		scriptPubkeys[i] = utxo.ScriptPubKey
	}
	feeSats := btcutil.Amount(0.0001 * btcutil.SatoshiPerBitcoin)
	amountSats := btcutil.Amount(amount * btcutil.SatoshiPerBitcoin)
	change := inputSats - feeSats - amountSats
	if change < 0 {
		return nil, fmt.Errorf("not enough input amount in sats; wanted %d, got %d", amountSats+feeSats, inputSats)
	}
	amountMap := map[btcutil.Address]btcutil.Amount{
		to:                 amountSats,
		btcDeployerAddress: change,
	}
	tx, err := btc.CreateRawTransaction(inputs, amountMap, nil)
	if err != nil {
		panic(err)
	}

	nulldata, err := txscript.NullDataScript(memo) // this adds a OP_RETURN + single BYTE len prefix to the data
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("nulldata (len %d): %x", len(nulldata), nulldata)
	if err != nil {
		panic(err)
	}
	memoOutput := wire.TxOut{Value: 0, PkScript: nulldata}
	tx.TxOut = append(tx.TxOut, &memoOutput)
	tx.TxOut[1], tx.TxOut[2] = tx.TxOut[2], tx.TxOut[1]

	// make sure that TxOut[0] is sent to "to" address; TxOut[2] is change to oneself. TxOut[1] is memo.
	if bytes.Compare(tx.TxOut[0].PkScript[2:], to.ScriptAddress()) != 0 {
		sm.Logger.Info("tx.TxOut[0].PkScript: %x", tx.TxOut[0].PkScript)
		sm.Logger.Info("to.ScriptAddress():   %x", to.ScriptAddress())
		sm.Logger.Info("swapping txout[0] with txout[2]")
		tx.TxOut[0], tx.TxOut[2] = tx.TxOut[2], tx.TxOut[0]
	}

	sm.Logger.Info("raw transaction: \n")
	for idx, txout := range tx.TxOut {
		sm.Logger.Info("txout %d", idx)
		sm.Logger.Info("  value: %d", txout.Value)
		sm.Logger.Info("  PkScript: %x", txout.PkScript)
	}

	inputsForSign := make([]btcjson.RawTxWitnessInput, len(inputs))
	for i, input := range inputs {
		inputsForSign[i] = btcjson.RawTxWitnessInput{
			Txid:         input.Txid,
			Vout:         input.Vout,
			Amount:       &amounts[i],
			ScriptPubKey: scriptPubkeys[i],
		}
	}

	//stx, signed, err := btc.SignRawTransactionWithWallet(tx)
	stx, signed, err := btc.SignRawTransactionWithWallet2(tx, inputsForSign)
	if err != nil {
		panic(err)
	}
	if !signed {
		panic("btc transaction not signed")
	}
	txid, err := btc.SendRawTransaction(stx, true)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("txid: %+v", txid)
	_, err = btc.GenerateToAddress(6, btcDeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	gtx, err := btc.GetTransaction(txid)
	if err != nil {
		panic(err)
	}
	sm.Logger.Info("rawtx confirmation: %d", gtx.BlockIndex)
	rawtx, err := btc.GetRawTransactionVerbose(txid)
	if err != nil {
		panic(err)
	}

	events := zetaclient.FilterAndParseIncomingTx(
		[]btcjson.TxRawResult{*rawtx},
		0,
		sm.BTCTSSAddress.EncodeAddress(),
		&log.Logger,
		common.BtcRegtestChain().ChainId,
	)
	sm.Logger.Info("bitcoin intx events:")
	for _, event := range events {
		sm.Logger.Info("  TxHash: %s", event.TxHash)
		sm.Logger.Info("  From: %s", event.FromAddress)
		sm.Logger.Info("  To: %s", event.ToAddress)
		sm.Logger.Info("  Amount: %f", event.Value)
		sm.Logger.Info("  Memo: %x", event.MemoBytes)
	}
	return txid, nil
}
