//go:build PRIVNET
// +build PRIVNET

package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rs/zerolog/log"
	zrc20 "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient"
)

var (
	BTCDeployerAddress *btcutil.AddressWitnessPubKeyHash
)

func (sm *SmokeTest) TestBitcoinSetup() {
	LoudPrintf("Setup Bitcoin\n")
	startTime := time.Now()
	defer func() {
		fmt.Printf("Bitcoin setup took %s\n", time.Since(startTime))
	}()

	btc := sm.btcRPCClient
	_, err := btc.CreateWallet("smoketest", rpcclient.WithCreateWalletBlank())
	if err != nil {
		if !strings.Contains(err.Error(), "Database already exists") {
			panic(err)
		}
	}

	skBytes, err := hex.DecodeString(DeployerPrivateKey)
	if err != nil {
		panic(err)
	}
	sk, _ := btcec.PrivKeyFromBytes(btcec.S256(), skBytes)
	privkeyWIF, err := btcutil.NewWIF(sk, &chaincfg.RegressionNetParams, true)
	if err != nil {
		panic(err)
	}
	err = btc.ImportPrivKeyRescan(privkeyWIF, "deployer", true)
	if err != nil {
		panic(err)
	}
	BTCDeployerAddress, err = btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(privkeyWIF.PrivKey.PubKey().SerializeCompressed()), &chaincfg.RegressionNetParams)
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
	_, err = btc.GenerateToAddress(4, BTCDeployerAddress, nil)
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
	fmt.Printf("  mine (Deployer): %+v\n", bals.Mine)
	if bals.WatchOnly != nil {
		fmt.Printf("  watchonly (TSSAddress): %+v\n", bals.WatchOnly)
	}
	fmt.Printf("  TSS Address: %s\n", BTCTSSAddress.EncodeAddress())

	sm.DepositBTC()
}

func (sm *SmokeTest) DepositBTC() {
	btc := sm.btcRPCClient
	utxos, err := sm.btcRPCClient.ListUnspent()
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
	fmt.Printf("ListUnspent:\n")
	fmt.Printf("  spendableAmount: %f\n", spendableAmount)
	fmt.Printf("  spendableUTXOs: %d\n", spendableUTXOs)
	fmt.Printf("Now sending two txs to TSS address...\n")
	_, err = SendToTSSFromDeployerToDeposit(BTCTSSAddress, 1.1, utxos[:2], btc)
	if err != nil {
		panic(err)
	}
	_, err = SendToTSSFromDeployerToDeposit(BTCTSSAddress, 0.05, utxos[2:4], btc)
	if err != nil {
		panic(err)
	}
	_, err = SendToTSSFromDeployerWithMemo(BTCTSSAddress, 0.11, utxos[4:5], btc, []byte(zetaclient.DonationMessage))
	if err != nil {
		panic(err)
	}

	fmt.Printf("testing if the deposit into BTC ZRC20 is successful...\n")

	// check if the deposit is successful
	BTCZRC20Addr, err := sm.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(common.BtcRegtestChain().ChainId))
	if err != nil {
		panic(err)
	}
	sm.BTCZRC20Addr = BTCZRC20Addr
	fmt.Printf("BTCZRC20Addr: %s\n", BTCZRC20Addr.Hex())
	BTCZRC20, err := zrc20.NewZRC20(BTCZRC20Addr, sm.zevmClient)
	if err != nil {
		panic(err)
	}
	sm.BTCZRC20 = BTCZRC20
	initialBalance, err := BTCZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	if err != nil {
		panic(err)
	}
	for {
		time.Sleep(5 * time.Second)
		balance, err := BTCZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
		if err != nil {
			panic(err)
		}
		diff := big.NewInt(0)
		diff.Sub(balance, initialBalance)
		fmt.Printf("BTC Difference in balance: %d", diff.Uint64())
		if diff.Cmp(big.NewInt(1.15*btcutil.SatoshiPerBitcoin)) != 0 {
			fmt.Printf("waiting for BTC balance to show up in ZRC contract... current bal %d\n", balance)
		} else {
			fmt.Printf("BTC balance is in ZRC contract! Success\n")
			break
		}
	}
}

func (sm *SmokeTest) DepositBTCRefund() {
	LoudPrintf("Deposit BTC with invalid memo; should be refunded\n")
	btc := sm.btcRPCClient
	utxos, err := sm.btcRPCClient.ListUnspent()
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
	fmt.Printf("ListUnspent:\n")
	fmt.Printf("  spendableAmount: %f\n", spendableAmount)
	fmt.Printf("  spendableUTXOs: %d\n", spendableUTXOs)
	fmt.Printf("Now sending two txs to TSS address...\n")
	_, err = SendToTSSFromDeployerToDeposit(BTCTSSAddress, 1.1, utxos[:2], btc)
	if err != nil {
		panic(err)
	}
	_, err = SendToTSSFromDeployerToDeposit(BTCTSSAddress, 0.05, utxos[2:4], btc)
	if err != nil {
		panic(err)
	}

	fmt.Printf("testing if the deposit into BTC ZRC20 is successful...\n")

	// check if the deposit is successful
	BTCZRC20Addr, err := sm.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(common.BtcRegtestChain().ChainId))
	if err != nil {
		panic(err)
	}
	sm.BTCZRC20Addr = BTCZRC20Addr
	fmt.Printf("BTCZRC20Addr: %s\n", BTCZRC20Addr.Hex())
	BTCZRC20, err := zrc20.NewZRC20(BTCZRC20Addr, sm.zevmClient)
	if err != nil {
		panic(err)
	}
	sm.BTCZRC20 = BTCZRC20
	initialBalance, err := BTCZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	if err != nil {
		panic(err)
	}
	for {
		time.Sleep(5 * time.Second)
		balance, err := BTCZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
		if err != nil {
			panic(err)
		}
		diff := big.NewInt(0)
		diff.Sub(balance, initialBalance)
		if diff.Cmp(big.NewInt(1.15*btcutil.SatoshiPerBitcoin)) != 0 {
			fmt.Printf("waiting for BTC balance to show up in ZRC contract... current bal %d\n", balance)
		} else {
			fmt.Printf("BTC balance is in ZRC contract! Success\n")
			break
		}
	}
}

func (sm *SmokeTest) TestBitcoinWithdraw() {
	startTime := time.Now()
	defer func() {
		fmt.Printf("Bitcoin withdraw took %s\n", time.Since(startTime))
	}()
	LoudPrintf("Testing Bitcoin ZRC20 Withdraw...\n")
	// withdraw 0.1 BTC from ZRC20 to BTC address
	// first, approve the ZRC20 contract to spend 1 BTC from the deployer address
	sm.WithdrawBitcoin()
}

func (sm *SmokeTest) WithdrawBitcoin() {
	amount := big.NewInt(0.1 * btcutil.SatoshiPerBitcoin)

	// check if the deposit is successful
	BTCZRC20Addr, err := sm.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(common.BtcRegtestChain().ChainId))
	if err != nil {
		panic(err)
	}
	fmt.Printf("BTCZRC20Addr: %s\n", BTCZRC20Addr.Hex())
	BTCZRC20, err := zrc20.NewZRC20(BTCZRC20Addr, sm.zevmClient)
	if err != nil {
		panic(err)
	}
	balance, err := BTCZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	if err != nil {
		panic(err)
	}
	if balance.Cmp(amount) < 0 {
		panic(fmt.Errorf("not enough balance in ZRC20 contract"))
	}
	// approve the ZRC20 contract to spend 1 BTC from the deployer address
	{
		tx, err := BTCZRC20.Approve(sm.zevmAuth, BTCZRC20Addr, big.NewInt(amount.Int64()*2)) // approve more to cover withdraw fee
		if err != nil {
			panic(err)
		}
		receipt := MustWaitForTxReceipt(sm.zevmClient, tx)
		fmt.Printf("approve receipt: status %d\n", receipt.Status)
		if receipt.Status != 1 {
			panic(fmt.Errorf("approve receipt status is not 1"))
		}
	}
	go func() {
		for {
			time.Sleep(5 * time.Second)
			_, err = sm.btcRPCClient.GenerateToAddress(1, BTCDeployerAddress, nil)
			if err != nil {
				panic(err)
			}
		}
	}()
	// withdraw 0.1 BTC from ZRC20 to BTC address
	{
		_, gasFee, err := BTCZRC20.WithdrawGasFee(&bind.CallOpts{})
		if err != nil {
			panic(err)
		}
		fmt.Printf("withdraw gas fee: %d\n", gasFee)
		tx, err := BTCZRC20.Withdraw(sm.zevmAuth, []byte(BTCDeployerAddress.EncodeAddress()), amount)
		if err != nil {
			panic(err)
		}
		receipt := MustWaitForTxReceipt(sm.zevmClient, tx)
		fmt.Printf("withdraw receipt: status %d\n", receipt.Status)
		if receipt.Status != 1 {
			panic(fmt.Errorf("withdraw receipt status is not 1"))
		}
		_, err = sm.btcRPCClient.GenerateToAddress(10, BTCDeployerAddress, nil)
		if err != nil {
			panic(err)
		}
		cctx := WaitCctxMinedByInTxHash(receipt.TxHash.Hex(), sm.cctxClient)
		outTxHash := cctx.GetCurrentOutTxParam().OutboundTxHash
		hash, err := chainhash.NewHashFromStr(outTxHash)
		if err != nil {
			panic(err)
		}

		rawTx, err := sm.btcRPCClient.GetRawTransactionVerbose(hash)
		if err != nil {
			panic(err)
		}
		fmt.Printf("raw tx:\n")
		fmt.Printf("  TxIn: %d\n", len(rawTx.Vin))
		for idx, txIn := range rawTx.Vin {
			fmt.Printf("  TxIn %d:\n", idx)
			fmt.Printf("    TxID:Vout:  %s:%d\n", txIn.Txid, txIn.Vout)
			fmt.Printf("    ScriptSig: %s\n", txIn.ScriptSig.Hex)
		}
		fmt.Printf("  TxOut: %d\n", len(rawTx.Vout))
		for _, txOut := range rawTx.Vout {
			fmt.Printf("  TxOut %d:\n", txOut.N)
			fmt.Printf("    Value: %.8f\n", txOut.Value)
			fmt.Printf("    ScriptPubKey: %s\n", txOut.ScriptPubKey.Hex)
		}
	}
}

func (sm *SmokeTest) WithdrawBitcoinMultipleTimes(repeat int64) {
	totalAmount := big.NewInt(int64(0.1 * 1e8))
	amount := big.NewInt(int64(0.1 * 1e8 / float64(repeat)))

	// check if the deposit is successful
	BTCZRC20Addr, err := sm.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(common.BtcRegtestChain().ChainId))
	if err != nil {
		panic(err)
	}
	fmt.Printf("BTCZRC20Addr: %s\n", BTCZRC20Addr.Hex())
	BTCZRC20, err := zrc20.NewZRC20(BTCZRC20Addr, sm.zevmClient)
	if err != nil {
		panic(err)
	}
	balance, err := BTCZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	if err != nil {
		panic(err)
	}
	if balance.Cmp(totalAmount) < 0 {
		panic(fmt.Errorf("not enough balance in ZRC20 contract"))
	}
	// approve the ZRC20 contract to spend 1 BTC from the deployer address
	{
		tx, err := BTCZRC20.Approve(sm.zevmAuth, BTCZRC20Addr, totalAmount.Mul(totalAmount, big.NewInt(100))) // approve more to cover withdraw fee
		if err != nil {
			panic(err)
		}
		receipt := MustWaitForTxReceipt(sm.zevmClient, tx)
		fmt.Printf("approve receipt: status %d\n", receipt.Status)
		if receipt.Status != 1 {
			panic(fmt.Errorf("approve receipt status is not 1"))
		}
	}
	go func() {
		for {
			time.Sleep(3 * time.Second)
			_, err = sm.btcRPCClient.GenerateToAddress(1, BTCDeployerAddress, nil)
			if err != nil {
				panic(err)
			}
		}
	}()
	// withdraw 0.1 BTC from ZRC20 to BTC address
	for i := int64(0); i < repeat; i++ {
		_, gasFee, err := BTCZRC20.WithdrawGasFee(&bind.CallOpts{})
		if err != nil {
			panic(err)
		}
		fmt.Printf("withdraw gas fee: %d\n", gasFee)
		tx, err := BTCZRC20.Withdraw(sm.zevmAuth, []byte(BTCDeployerAddress.EncodeAddress()), amount)
		if err != nil {
			panic(err)
		}
		receipt := MustWaitForTxReceipt(sm.zevmClient, tx)
		fmt.Printf("withdraw receipt: status %d\n", receipt.Status)
		if receipt.Status != 1 {
			panic(fmt.Errorf("withdraw receipt status is not 1"))
		}
		_, err = sm.btcRPCClient.GenerateToAddress(10, BTCDeployerAddress, nil)
		if err != nil {
			panic(err)
		}
		cctx := WaitCctxMinedByInTxHash(receipt.TxHash.Hex(), sm.cctxClient)
		outTxHash := cctx.GetCurrentOutTxParam().OutboundTxHash
		hash, err := chainhash.NewHashFromStr(outTxHash)
		if err != nil {
			panic(err)
		}

		rawTx, err := sm.btcRPCClient.GetRawTransactionVerbose(hash)
		if err != nil {
			panic(err)
		}
		fmt.Printf("raw tx:\n")
		fmt.Printf("  TxIn: %d\n", len(rawTx.Vin))
		for idx, txIn := range rawTx.Vin {
			fmt.Printf("  TxIn %d:\n", idx)
			fmt.Printf("    TxID:Vout:  %s:%d\n", txIn.Txid, txIn.Vout)
			fmt.Printf("    ScriptSig: %s\n", txIn.ScriptSig.Hex)
		}
		fmt.Printf("  TxOut: %d\n", len(rawTx.Vout))
		for _, txOut := range rawTx.Vout {
			fmt.Printf("  TxOut %d:\n", txOut.N)
			fmt.Printf("    Value: %.8f\n", txOut.Value)
			fmt.Printf("    ScriptPubKey: %s\n", txOut.ScriptPubKey.Hex)
		}
	}
}

func SendToTSSFromDeployerToDeposit(to btcutil.Address, amount float64, inputUTXOs []btcjson.ListUnspentResult, btc *rpcclient.Client) (*chainhash.Hash, error) {
	return SendToTSSFromDeployerWithMemo(to, amount, inputUTXOs, btc, DeployerAddress.Bytes())
}

func SendToTSSFromDeployerWithMemo(to btcutil.Address, amount float64, inputUTXOs []btcjson.ListUnspentResult, btc *rpcclient.Client, memo []byte) (*chainhash.Hash, error) {
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
		BTCDeployerAddress: change,
	}
	tx, err := btc.CreateRawTransaction(inputs, amountMap, nil)
	if err != nil {
		panic(err)
	}

	nulldata, err := txscript.NullDataScript(memo) // this adds a OP_RETURN + single BYTE len prefix to the data
	if err != nil {
		panic(err)
	}
	fmt.Printf("nulldata (len %d): %x\n", len(nulldata), nulldata)
	if err != nil {
		panic(err)
	}
	memoOutput := wire.TxOut{Value: 0, PkScript: nulldata}
	tx.TxOut = append(tx.TxOut, &memoOutput)
	tx.TxOut[1], tx.TxOut[2] = tx.TxOut[2], tx.TxOut[1]

	// make sure that TxOut[0] is sent to "to" address; TxOut[2] is change to oneself. TxOut[1] is memo.
	if bytes.Compare(tx.TxOut[0].PkScript[2:], to.ScriptAddress()) != 0 {
		fmt.Printf("tx.TxOut[0].PkScript: %x\n", tx.TxOut[0].PkScript)
		fmt.Printf("to.ScriptAddress():   %x\n", to.ScriptAddress())
		fmt.Printf("swapping txout[0] with txout[2]\n")
		tx.TxOut[0], tx.TxOut[2] = tx.TxOut[2], tx.TxOut[0]
	}

	fmt.Printf("raw transaction: \n")
	for idx, txout := range tx.TxOut {
		fmt.Printf("txout %d\n", idx)
		fmt.Printf("  value: %d\n", txout.Value)
		fmt.Printf("  PkScript: %x\n", txout.PkScript)
	}
	var inputsForSign []btcjson.RawTxWitnessInput
	for i, input := range inputs {
		inputsForSign = append(inputsForSign, btcjson.RawTxWitnessInput{
			Txid: input.Txid, Vout: input.Vout, Amount: &amounts[i], ScriptPubKey: scriptPubkeys[i]})
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
	fmt.Printf("txid: %+v\n", txid)
	_, err = btc.GenerateToAddress(6, BTCDeployerAddress, nil)
	if err != nil {
		panic(err)
	}
	gtx, err := btc.GetTransaction(txid)
	if err != nil {
		panic(err)
	}
	fmt.Printf("rawtx confirmation: %d\n", gtx.BlockIndex)
	rawtx, err := btc.GetRawTransactionVerbose(txid)
	if err != nil {
		panic(err)
	}

	events := zetaclient.FilterAndParseIncomingTx([]btcjson.TxRawResult{*rawtx}, 0, BTCTSSAddress.EncodeAddress(), &log.Logger)
	fmt.Printf("bitcoin intx events:\n")
	for _, event := range events {
		fmt.Printf("  TxHash: %s\n", event.TxHash)
		fmt.Printf("  From: %s\n", event.FromAddress)
		fmt.Printf("  To: %s\n", event.ToAddress)
		fmt.Printf("  Amount: %f\n", event.Value)
		fmt.Printf("  Memo: %x\n", event.MemoBytes)
	}
	return txid, nil
}
