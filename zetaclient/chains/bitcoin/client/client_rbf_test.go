package client_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/client"
	btccommon "github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
	"github.com/zeta-chain/node/zetaclient/common"
	"github.com/zeta-chain/node/zetaclient/config"
)

// Test_BitcoinRBFLive runs RBF tests on a live Bitcoin network.
func Test_BitcoinRBFLive(t *testing.T) {
	if !common.LiveTestEnabled() {
		t.Skip("skipping live test")
	}

	// RBF tests requires private key to be set
	if common.IsEnvVariableSet(common.EnvBTCTestPK) {
		t.Run("RBFTransaction", Run_RBFTransaction)
		t.Run("RBFTransaction_Chained_CPFP", Run_RBFTransaction_Chained_CPFP)
	}

	t.Run("PendingMempoolTx", Run_PendingMempoolTx)
}

// setupRBFTest initializes the test suite, privateKey, sender, receiver
func setupRBFTest(t *testing.T) (*testSuite, *secp256k1.PrivateKey, btcutil.Address, btcutil.Address) {
	// network to use
	chain := chains.BitcoinTestnet4
	net, err := chains.GetBTCChainParams(chain.ChainId)
	require.NoError(t, err)
	config := config.BTCConfig{
		RPCHost:   os.Getenv(common.EnvBtcRPCMainnet),
		RPCParams: "testnet3",
	}

	// load test private key
	privKeyHex := os.Getenv(common.EnvBTCTestPK)
	privKeyBytes, err := hex.DecodeString(privKeyHex)
	require.NoError(t, err)

	// construct a secp256k1 private key object
	privKey := secp256k1.PrivKeyFromBytes(privKeyBytes)
	pubKeyHash := btcutil.Hash160(privKey.PubKey().SerializeCompressed())
	sender, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, net)
	require.NoError(t, err)
	fmt.Printf("sender  : %s\n", sender.EncodeAddress())

	// receiver address
	to, err := btcutil.DecodeAddress("tb1qxr8zcffrkmqwvtkzjz8nxs05p2vs6pt9rzq27a", net)
	require.NoError(t, err)
	fmt.Printf("receiver: %s\n", to.EncodeAddress())

	// create a new test suite
	ts := newTestSuite(t, config)

	return ts, privKey, sender, to
}

func Run_RBFTransaction(t *testing.T) {
	// setup test
	ts, privKey, sender, to := setupRBFTest(t)

	// define amount, fee rate and bump fee reserved
	amount := 0.00001
	nonceMark := chains.NonceMarkAmount(1)
	feeRate := int64(2)
	bumpFeeReserved := int64(10000)

	// STEP 1
	// build and send tx1
	nonceMark += 1
	txHash1 := buildAndSendRBFTx(
		t,
		ts.ctx,
		ts.Client,
		privKey,
		nil,
		sender,
		to,
		amount,
		nonceMark,
		feeRate,
		bumpFeeReserved,
	)
	fmt.Printf("sent tx1: %s\n", txHash1)

	// STEP 2
	// build and send tx2 (child of tx1)
	nonceMark += 1
	txHash2 := buildAndSendRBFTx(
		t,
		ts.ctx,
		ts.Client,
		privKey,
		txHash1,
		sender,
		to,
		amount,
		nonceMark,
		feeRate,
		bumpFeeReserved,
	)
	fmt.Printf("sent tx2: %s\n", txHash2)

	// STEP 3
	// wait for a short time before bumping fee
	rawTx1, confirmed := waitForTxConfirmation(ts.ctx, ts.Client, sender, txHash1, 10*time.Second)
	if confirmed {
		fmt.Println("Opps: tx1 confirmed, no chance to bump fee; please start over")
		return
	}

	// STEP 4
	// bump gas fee for tx1 (the parent of tx2)
	// we assume that tx1, tx2 and tx3 have same vBytes for simplicity
	// two rules to satisfy:
	//   - feeTx3 >= feeTx1 + feeTx2
	//   - additionalFees >= vSizeTx3 * minRelayFeeRate
	// see: https://github.com/bitcoin/bitcoin/blob/5b8046a6e893b7fad5a93631e6d1e70db31878af/src/policy/rbf.cpp#L166-L183
	minRelayFeeRate := int64(1)
	feeRateIncrease := minRelayFeeRate
	sizeTx3 := mempool.GetTxVirtualSize(rawTx1)
	additionalFees := (sizeTx3 + 1) * (feeRate + feeRateIncrease) // +1 in case Bitcoin Core rounds up the vSize
	fmt.Printf("additional fee: %d sats\n", additionalFees)
	tx3, err := bumpRBFTxFee(rawTx1.MsgTx(), additionalFees)
	require.NoError(t, err)

	// STEP 5
	// sign and send tx3, which replaces tx1
	signTx(t, ts.ctx, ts.Client, privKey, tx3)
	txHash3, err := ts.Client.SendRawTransaction(ts.ctx, tx3, true)
	require.NoError(t, err)
	fmt.Printf("sent tx3: %s\n", txHash3)

	// STEP 6
	// wait for tx3 confirmation
	rawTx3, confirmed := waitForTxConfirmation(ts.ctx, ts.Client, sender, txHash3, 30*time.Minute)
	require.True(t, confirmed)
	printTx(rawTx3.MsgTx())
	fmt.Println("tx3 confirmed")

	// STEP 7
	// tx1 and tx2 must be dropped
	ensureTxDropped(t, ts.ctx, ts.Client, txHash1)
	fmt.Println("tx1 dropped")
	ensureTxDropped(t, ts.ctx, ts.Client, txHash2)
	fmt.Println("tx2 dropped")
}

// Run_RBFTransaction_Chained_CPFP tests Child-Pays-For-Parent (CPFP) fee bumping strategy for chained RBF transactions
func Run_RBFTransaction_Chained_CPFP(t *testing.T) {
	// setup test
	ts, privKey, sender, to := setupRBFTest(t)

	// define amount, fee rate and bump fee reserved
	amount := 0.00001
	nonceMark := chains.NonceMarkAmount(0)
	feeRate := int64(2)
	bumpFeeReserved := int64(10000)

	// STEP 1
	// build and send tx1
	nonceMark += 1
	txHash1 := buildAndSendRBFTx(
		t,
		ts.ctx,
		ts.Client,
		privKey,
		nil,
		sender,
		to,
		amount,
		nonceMark,
		feeRate,
		bumpFeeReserved,
	)
	fmt.Printf("sent tx1: %s\n", txHash1)

	// STEP 2
	// build and send tx2 (child of tx1)
	nonceMark += 1
	txHash2 := buildAndSendRBFTx(
		t,
		ts.ctx,
		ts.Client,
		privKey,
		txHash1,
		sender,
		to,
		amount,
		nonceMark,
		feeRate,
		bumpFeeReserved,
	)
	fmt.Printf("sent tx2: %s\n", txHash2)

	// STEP 3
	// build and send tx3 (child of tx2)
	nonceMark += 1
	txHash3 := buildAndSendRBFTx(
		t,
		ts.ctx,
		ts.Client,
		privKey,
		txHash2,
		sender,
		to,
		amount,
		nonceMark,
		feeRate,
		bumpFeeReserved,
	)
	fmt.Printf("sent tx3: %s\n", txHash3)

	// STEP 4
	// wait for a short time before bumping fee
	rawTx3, confirmed := waitForTxConfirmation(ts.ctx, ts.Client, sender, txHash3, 10*time.Second)
	if confirmed {
		fmt.Println("Opps: tx3 confirmed, no chance to bump fee; please start over")
		return
	}

	// STEP 5
	// bump gas fee for tx3 (the child/grandchild of tx1/tx2)
	// we assume that tx3 has same vBytes as the fee-bump tx (tx4) for simplicity
	// two rules to satisfy:
	//   - feeTx4 >= feeTx3
	//   - additionalFees >= vSizeTx4 * minRelayFeeRate
	// see: https://github.com/bitcoin/bitcoin/blob/5b8046a6e893b7fad5a93631e6d1e70db31878af/src/policy/rbf.cpp#L166-L183
	minRelayFeeRate := int64(1)
	feeRateIncrease := minRelayFeeRate
	additionalFees := (mempool.GetTxVirtualSize(rawTx3) + 1) * feeRateIncrease
	fmt.Printf("additional fee: %d sats\n", additionalFees)
	tx4, err := bumpRBFTxFee(rawTx3.MsgTx(), additionalFees)
	require.NoError(t, err)

	// STEP 6
	// sign and send tx4, which replaces tx3
	signTx(t, ts.ctx, ts.Client, privKey, tx4)
	txHash4, err := ts.Client.SendRawTransaction(ts.ctx, tx4, true)
	require.NoError(t, err)
	fmt.Printf("sent tx4: %s\n", txHash4)

	// STEP 7
	// wait for tx4 confirmation
	rawTx4, confirmed := waitForTxConfirmation(ts.ctx, ts.Client, sender, txHash4, 30*time.Minute)
	require.True(t, confirmed)
	printTx(rawTx4.MsgTx())
	fmt.Println("tx4 confirmed")

	// STEP 8
	// tx3 must be dropped
	ensureTxDropped(t, ts.ctx, ts.Client, txHash3)
	fmt.Println("tx1 dropped")
}

func Run_PendingMempoolTx(t *testing.T) {
	// network to use
	config := config.BTCConfig{
		RPCHost:   os.Getenv(common.EnvBtcRPCMainnet),
		RPCParams: "mainnet",
	}

	// create a new test suite
	ts := newTestSuite(t, config)

	// get mempool transactions
	mempoolTxs, err := ts.Client.GetRawMempool(ts.ctx)
	require.NoError(t, err)
	fmt.Printf("mempool txs: %d\n", len(mempoolTxs))

	// get last block height
	lastHeight, err := ts.Client.GetBlockCount(ts.ctx)
	require.NoError(t, err)
	fmt.Printf("block height: %d\n", lastHeight)

	const (
		// average minutes per block is about 10 minutes
		minutesPerBlockAverage = 10.0

		// maxBlockTimeDiffPercentage is the maximum error percentage between the estimated and actual block time
		// note: 25% is a percentage to make sure the test is not too strict
		maxBlockTimeDiffPercentage = 0.25
	)

	// the goal of the test is to ensure the 'Time' and 'Height' provided by the mempool are correct.
	// otherwise, zetaclient should not rely on these information to schedule RBF/CPFP transactions.
	// loop through the mempool to sample N pending txs that are pending for more than 2 hours
	N := 10
	for i := len(mempoolTxs) - 1; i >= 0; i-- {
		txHash := mempoolTxs[i]
		entry, err := ts.Client.GetMempoolEntry(ts.ctx, txHash.String())
		if err == nil {
			require.Positive(t, entry.Fee)
			txTime := time.Unix(entry.Time, 0)
			txTimeStr := txTime.Format(time.DateTime)
			elapsed := time.Since(txTime)
			if elapsed > 30*time.Minute {
				// calculate average block time
				elapsedBlocks := lastHeight - entry.Height
				minutesPerBlockCalculated := elapsed.Minutes() / float64(elapsedBlocks)
				blockTimeDiff := minutesPerBlockAverage - minutesPerBlockCalculated
				if blockTimeDiff < 0 {
					blockTimeDiff = -blockTimeDiff
				}

				// the block time difference should fall within 25% of the average block time
				require.Less(t, blockTimeDiff, minutesPerBlockAverage*maxBlockTimeDiffPercentage)
				fmt.Printf(
					"txid: %s, height: %d, time: %s, pending: %f minutes, block time: %f minutes, diff: %f%%\n",
					txHash,
					entry.Height,
					txTimeStr,
					elapsed.Minutes(),
					minutesPerBlockCalculated,
					blockTimeDiff/minutesPerBlockAverage*100,
				)

				// break if we have enough samples
				if N -= 1; N == 0 {
					break
				}
			}
		}
	}
}

// buildAndSendRBFTx builds, signs and sends an RBF transaction
func buildAndSendRBFTx(
	t *testing.T,
	ctx context.Context,
	client *client.Client,
	privKey *secp256k1.PrivateKey,
	parent *chainhash.Hash,
	sender, to btcutil.Address,
	amount float64,
	nonceMark int64,
	feeRate int64,
	bumpFeeReserved int64,
) *chainhash.Hash {
	// list outputs
	utxos := listUTXOs(ctx, client, sender)
	require.NotEmpty(t, utxos)

	// ensure all inputs are from the parent tx
	if parent != nil {
		for _, out := range utxos {
			require.Equal(t, parent.String(), out.TxID)
		}
	}

	// build tx opt-in RBF
	tx := buildRBFTx(t, utxos, sender, to, amount, nonceMark, feeRate, bumpFeeReserved)

	// sign tx
	signTx(t, ctx, client, privKey, tx)

	// broadcast tx
	txHash, err := client.SendRawTransaction(ctx, tx, true)
	require.NoError(t, err)

	return txHash
}

func listUTXOs(ctx context.Context, client *client.Client, address btcutil.Address) []btcjson.ListUnspentResult {
	utxos, err := client.ListUnspentMinMaxAddresses(ctx, 0, 9999999, []btcutil.Address{address})
	if err != nil {
		fmt.Printf("ListUnspent failed: %s\n", err)
		return nil
	}

	// sort utxos by amount, txid, vout
	sort.SliceStable(utxos, func(i, j int) bool {
		if utxos[i].Amount == utxos[j].Amount {
			if utxos[i].TxID == utxos[j].TxID {
				return utxos[i].Vout < utxos[j].Vout
			}
			return utxos[i].TxID < utxos[j].TxID
		}
		return utxos[i].Amount < utxos[j].Amount
	})

	// print utxos
	fmt.Println("utxos:")
	for _, out := range utxos {
		fmt.Printf(
			"  txid: %s, vout: %d, amount: %f, confirmation: %d\n",
			out.TxID,
			out.Vout,
			out.Amount,
			out.Confirmations,
		)
	}

	return utxos
}

func buildRBFTx(
	t *testing.T,
	utxos []btcjson.ListUnspentResult,
	sender, to btcutil.Address,
	amount float64,
	nonceMark int64,
	feeRate int64,
	bumpFeeReserved int64,
) *wire.MsgTx {
	// build tx with all unspents
	total := 0.0
	tx := wire.NewMsgTx(wire.TxVersion)
	for _, output := range utxos {
		hash, err := chainhash.NewHashFromStr(output.TxID)
		require.NoError(t, err)

		// add input
		outpoint := wire.NewOutPoint(hash, output.Vout)
		txIn := wire.NewTxIn(outpoint, nil, nil)
		txIn.Sequence = 1 // opt-in for RBF
		tx.AddTxIn(txIn)
		total += output.Amount
	}
	totalSats, err := btccommon.GetSatoshis(total)
	require.NoError(t, err)

	// amount to send in satoshis
	amountSats, err := btccommon.GetSatoshis(amount)
	require.NoError(t, err)

	// calculate tx fee
	txSize, err := btccommon.EstimateOutboundSize(int64(len(utxos)), []btcutil.Address{to})
	require.NoError(t, err)
	fees := int64(txSize) * feeRate

	// make sure total is greater than amount + fees
	require.GreaterOrEqual(t, totalSats, nonceMark+amountSats+fees+bumpFeeReserved)

	// 1st output: simulated nonce-mark amount to self
	pkScriptSender, err := txscript.PayToAddrScript(sender)
	require.NoError(t, err)
	txOut0 := wire.NewTxOut(nonceMark, pkScriptSender)
	tx.AddTxOut(txOut0)

	// 2nd output: payment to receiver
	pkScriptReceiver, err := txscript.PayToAddrScript(to)
	require.NoError(t, err)
	txOut1 := wire.NewTxOut(amountSats, pkScriptReceiver)
	tx.AddTxOut(txOut1)

	// 3rd output: change to self
	changeSats := totalSats - nonceMark - amountSats - fees
	require.GreaterOrEqual(t, changeSats, bumpFeeReserved)
	txOut2 := wire.NewTxOut(changeSats, pkScriptSender)
	tx.AddTxOut(txOut2)

	return tx
}

func signTx(t *testing.T, ctx context.Context, client *client.Client, privKey *secp256k1.PrivateKey, tx *wire.MsgTx) {
	// we know that the first output is the nonce-mark amount, so it contains the sender pkScript
	pkScriptSender := tx.TxOut[0].PkScript

	sigHashes := txscript.NewTxSigHashes(tx, txscript.NewCannedPrevOutputFetcher([]byte{}, 0))
	for idx, input := range tx.TxIn {
		// get input amount from previous tx outpoint via RPC
		rawTx, err := client.GetRawTransaction(ctx, &input.PreviousOutPoint.Hash)
		require.NoError(t, err)
		amount := rawTx.MsgTx().TxOut[input.PreviousOutPoint.Index].Value

		// calculate witness signature hash for signing
		witnessHash, err := txscript.CalcWitnessSigHash(pkScriptSender, sigHashes, txscript.SigHashAll, tx, idx, amount)
		require.NoError(t, err)

		// sign the witness hash
		sig := ecdsa.Sign(privKey, witnessHash)
		tx.TxIn[idx].Witness = wire.TxWitness{
			append(sig.Serialize(), byte(txscript.SigHashAll)),
			privKey.PubKey().SerializeCompressed(),
		}
	}

	printTx(tx)
}

func printTx(tx *wire.MsgTx) {
	fmt.Printf("\n==============================================================\n")
	fmt.Printf("tx version: %d\n", tx.Version)
	fmt.Printf("tx locktime: %d\n", tx.LockTime)
	fmt.Println("tx inputs:")
	for i, vin := range tx.TxIn {
		fmt.Printf("  input[%d]:\n", i)
		fmt.Printf("    prevout hash: %s\n", vin.PreviousOutPoint.Hash)
		fmt.Printf("    prevout index: %d\n", vin.PreviousOutPoint.Index)
		fmt.Printf("    sig script: %s\n", hex.EncodeToString(vin.SignatureScript))
		fmt.Printf("    sequence: %d\n", vin.Sequence)
		fmt.Printf("    witness: \n")
		for j, w := range vin.Witness {
			fmt.Printf("      witness[%d]: %s\n", j, hex.EncodeToString(w))
		}
	}
	fmt.Println("tx outputs:")
	for i, vout := range tx.TxOut {
		fmt.Printf("  output[%d]:\n", i)
		fmt.Printf("    value: %d\n", vout.Value)
		fmt.Printf("    script: %s\n", hex.EncodeToString(vout.PkScript))
	}
	fmt.Printf("==============================================================\n\n")
}

func peekUnconfirmedTx(ctx context.Context, client *client.Client, txHash *chainhash.Hash) (*btcutil.Tx, bool) {
	confirmed := false

	// try querying tx result
	getTxResult, err := client.GetTransaction(ctx, txHash)
	if err == nil {
		confirmed = getTxResult.Confirmations > 0
		fmt.Printf("tx confirmations: %d\n", getTxResult.Confirmations)
	} else {
		fmt.Printf("GetTxResultByHash failed: %s\n", err)
	}

	// query tx from mempool
	entry, err := client.GetMempoolEntry(ctx, txHash.String())
	switch {
	case err != nil:
		fmt.Println("tx in mempool: NO")
	default:
		txTime := time.Unix(entry.Time, 0)
		txTimeStr := txTime.Format(time.DateTime)
		elapsed := int64(time.Since(txTime).Seconds())
		fmt.Printf(
			"tx in mempool: YES, VSize: %d, height: %d, time: %s, elapsed: %d\n",
			entry.VSize,
			entry.Height,
			txTimeStr,
			elapsed,
		)
	}

	// query the raw tx
	rawTx, err := client.GetRawTransaction(ctx, txHash)
	if err != nil {
		fmt.Printf("GetRawTransaction failed: %s\n", err)
	}

	return rawTx, confirmed
}

func waitForTxConfirmation(
	ctx context.Context,
	client *client.Client,
	sender btcutil.Address,
	txHash *chainhash.Hash,
	timeOut time.Duration,
) (*btcutil.Tx, bool) {
	start := time.Now()
	for {
		rawTx, confirmed := peekUnconfirmedTx(ctx, client, txHash)
		listUTXOs(ctx, client, sender)
		fmt.Println()

		if confirmed {
			return rawTx, true
		}
		if time.Since(start) > timeOut {
			return rawTx, false
		}

		time.Sleep(5 * time.Second)
	}
}

func bumpRBFTxFee(oldTx *wire.MsgTx, additionalFee int64) (*wire.MsgTx, error) {
	// copy the old tx and reset
	newTx := oldTx.Copy()
	for idx := range newTx.TxIn {
		newTx.TxIn[idx].Witness = wire.TxWitness{}
		newTx.TxIn[idx].Sequence = 1
	}

	// original change needs to be enough to cover the additional fee
	if newTx.TxOut[2].Value <= additionalFee {
		return nil, errors.New("change amount is not enough to cover the additional fee")
	}

	// bump fee by reducing the change amount
	newTx.TxOut[2].Value = newTx.TxOut[2].Value - additionalFee

	return newTx, nil
}

func ensureTxDropped(t *testing.T, ctx context.Context, client *client.Client, txHash *chainhash.Hash) {
	// dropped tx must has negative confirmations (if returned)
	getTxResult, err := client.GetTransaction(ctx, txHash)
	if err == nil {
		require.Negative(t, getTxResult.Confirmations)
	}

	// dropped tx should be removed from mempool
	entry, err := client.GetMempoolEntry(ctx, txHash.String())
	require.Error(t, err)
	require.Nil(t, entry)

	// dropped tx should not be found
	// -5: No such mempool or blockchain transaction
	rawTx, err := client.GetRawTransaction(ctx, txHash)
	require.Error(t, err)
	require.Nil(t, rawTx)
}
