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
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/bitcoin"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/utils"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient"
)

// DepositBTC deposits BTC on ZetaChain
func (sm *SmokeTestRunner) DepositBTC() {
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
	fmt.Printf("ListUnspent:\n")
	fmt.Printf("  spendableAmount: %f\n", spendableAmount)
	fmt.Printf("  spendableUTXOs: %d\n", spendableUTXOs)
	fmt.Printf("Now sending two txs to TSS address...\n")
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

	fmt.Printf("testing if the deposit into BTC ZRC20 is successful...\n")

	// check if the deposit is successful
	BTCZRC20Addr, err := sm.SystemContract.GasCoinZRC20ByChainId(&bind.CallOpts{}, big.NewInt(common.BtcRegtestChain().ChainId))
	if err != nil {
		panic(err)
	}
	sm.BTCZRC20Addr = BTCZRC20Addr
	fmt.Printf("BTCZRC20Addr: %s\n", BTCZRC20Addr.Hex())
	BTCZRC20, err := zrc20.NewZRC20(BTCZRC20Addr, sm.ZevmClient)
	if err != nil {
		panic(err)
	}
	sm.BTCZRC20 = BTCZRC20
	initialBalance, err := BTCZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
	if err != nil {
		panic(err)
	}
	for {
		time.Sleep(5 * time.Second)
		balance, err := BTCZRC20.BalanceOf(&bind.CallOpts{}, sm.DeployerAddress)
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

	// prove the two transactions of the deposit
	utils.LoudPrintf("Bitcoin Merkle Proof\n")

	sm.ProveBTCTransaction(txHash1)
	sm.ProveBTCTransaction(txHash2)
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
			fmt.Printf("waiting for block header to show up in observer... current hash %s; err %s\n", hash.String(), err.Error())
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
	fmt.Printf("OK: txProof verified for inTx: %s\n", txHash.String())
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
	fmt.Printf("txid: %+v\n", txid)
	_, err = btc.GenerateToAddress(6, btcDeployerAddress, nil)
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

	events := zetaclient.FilterAndParseIncomingTx(
		[]btcjson.TxRawResult{*rawtx},
		0,
		sm.BTCTSSAddress.EncodeAddress(),
		&log.Logger,
		common.BtcRegtestChain().ChainId,
	)
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
