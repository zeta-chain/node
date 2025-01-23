package signer

import (
	"context"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	btcecdsa "github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
)

const (
	// the maximum number of inputs per outbound
	MaxNoOfInputsPerTx = 20

	// the rank below (or equal to) which we consolidate UTXOs
	consolidationRank = 10

	// reservedRBFFees is the amount of BTC reserved for RBF fee bumping.
	// the TSS keysign stops automatically when transactions get stuck in the mempool
	// 0.01 BTC can bump 10 transactions (1KB each) by 100 sat/vB
	reservedRBFFees = 0.01

	// rbfTxInSequenceNum is the sequence number used to signal an opt-in full-RBF (Replace-By-Fee) transaction
	// Setting sequenceNum to "1" effectively makes the transaction timelocks irrelevant.
	// See: https://github.com/bitcoin/bips/blob/master/bip-0125.mediawiki
	// See: https://github.com/BlockchainCommons/Learning-Bitcoin-from-the-Command-Line/blob/master/05_2_Resending_a_Transaction_with_RBF.md
	rbfTxInSequenceNum uint32 = 1
)

// SignWithdrawTx signs a BTC withdrawal tx and returns the signed tx
func (signer *Signer) SignWithdrawTx(
	ctx context.Context,
	txData *OutboundData,
	ob *observer.Observer,
) (*wire.MsgTx, error) {
	nonceMark := chains.NonceMarkAmount(txData.nonce)
	estimateFee := float64(txData.feeRate*common.OutboundBytesMax) / 1e8
	totalAmount := txData.amount + estimateFee + reservedRBFFees + float64(nonceMark)*1e-8

	// refreshing UTXO list before TSS keysign is important:
	// 1. all TSS outbounds have opted-in for RBF to be replaceable
	// 2. using old UTXOs may lead to accidental double-spending, which may trigger unwanted RBF
	//
	// Note: unwanted RBF is very unlikely to happen for two reasons:
	// 1. it requires 2/3 TSS signers to accidentally sign the same tx using same outdated UTXOs.
	// 2. RBF requires a higher fee rate than the original tx, otherwise it will fail.
	err := ob.FetchUTXOs(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "FetchUTXOs failed")
	}

	// select N UTXOs to cover the total expense
	selected, err := ob.SelectUTXOs(
		ctx,
		totalAmount,
		MaxNoOfInputsPerTx,
		txData.nonce,
		consolidationRank,
		false,
	)
	if err != nil {
		return nil, err
	}

	// build tx and add inputs
	tx := wire.NewMsgTx(wire.TxVersion)
	inAmounts, err := signer.AddTxInputs(tx, selected.UTXOs)
	if err != nil {
		return nil, err
	}

	// size checking
	// #nosec G115 always positive
	txSize, err := common.EstimateOutboundSize(int64(len(selected.UTXOs)), []btcutil.Address{txData.to})
	if err != nil {
		return nil, err
	}
	if txData.txSize < common.BtcOutboundBytesWithdrawer { // ZRC20 'withdraw' charged less fee from end user
		signer.Logger().Std.Info().
			Msgf("txSize %d is less than BtcOutboundBytesWithdrawer %d for nonce %d", txData.txSize, txSize, txData.nonce)
	}
	if txSize < common.OutboundBytesMin { // outbound shouldn't be blocked by low sizeLimit
		signer.Logger().Std.Warn().
			Msgf("txSize %d is less than outboundBytesMin %d; use outboundBytesMin", txSize, common.OutboundBytesMin)
		txSize = common.OutboundBytesMin
	}
	if txSize > common.OutboundBytesMax { // in case of accident
		signer.Logger().Std.Warn().
			Msgf("txSize %d is greater than outboundBytesMax %d; use outboundBytesMax", txSize, common.OutboundBytesMax)
		txSize = common.OutboundBytesMax
	}

	// fee calculation
	// #nosec G115 always in range (checked above)
	fees := txSize * txData.feeRate
	signer.Logger().
		Std.Info().
		Msgf("bitcoin outbound nonce %d feeRate %d size %d fees %d consolidated %d utxos of value %v",
			txData.nonce, txData.feeRate, txSize, fees, selected.ConsolidatedUTXOs, selected.ConsolidatedValue)

	// add tx outputs
	inputValue := selected.Value
	err = signer.AddWithdrawTxOutputs(tx, txData.to, inputValue, txData.amountSats, nonceMark, fees, txData.cancelTx)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	// sign the tx
	err = signer.SignTx(ctx, tx, inAmounts, txData.height, txData.nonce)
	if err != nil {
		return nil, errors.Wrap(err, "SignTx failed")
	}

	return tx, nil
}

// AddTxInputs adds the inputs to the tx and returns input amounts
func (signer *Signer) AddTxInputs(tx *wire.MsgTx, utxos []btcjson.ListUnspentResult) ([]int64, error) {
	amounts := make([]int64, len(utxos))
	for i, utxo := range utxos {
		hash, err := chainhash.NewHashFromStr(utxo.TxID)
		if err != nil {
			return nil, err
		}

		// add input and set 'nSequence' to opt-in for RBF
		// it doesn't matter on which input we set the RBF sequence
		outpoint := wire.NewOutPoint(hash, utxo.Vout)
		txIn := wire.NewTxIn(outpoint, nil, nil)
		if i == 0 {
			txIn.Sequence = rbfTxInSequenceNum
		}
		tx.AddTxIn(txIn)

		// store the amount for later signing use
		amount, err := common.GetSatoshis(utxos[i].Amount)
		if err != nil {
			return nil, err
		}
		amounts[i] = amount
	}

	return amounts, nil
}

// AddWithdrawTxOutputs adds the 3 outputs to the withdraw tx
// 1st output: the nonce-mark btc to TSS itself
// 2nd output: the payment to the recipient
// 3rd output: the remaining btc to TSS itself
func (signer *Signer) AddWithdrawTxOutputs(
	tx *wire.MsgTx,
	to btcutil.Address,
	inputValue float64,
	amountSats int64,
	nonceMark int64,
	fees int64,
	cancelTx bool,
) error {
	// convert withdraw amount to BTC
	amount := float64(amountSats) / 1e8

	// calculate remaining btc (the change) to TSS self
	remaining := inputValue - amount
	remainingSats, err := common.GetSatoshis(remaining)
	if err != nil {
		return err
	}
	remainingSats -= fees
	remainingSats -= nonceMark
	if remainingSats < 0 {
		return fmt.Errorf("remainder value is negative: %d", remainingSats)
	} else if remainingSats == nonceMark {
		signer.Logger().Std.Info().Msgf("adjust remainder value to avoid duplicate nonce-mark: %d", remainingSats)
		remainingSats--
	}

	// 1st output: the nonce-mark btc to TSS self
	payToSelfScript, err := signer.TSSToPkScript()
	if err != nil {
		return err
	}
	txOut1 := wire.NewTxOut(nonceMark, payToSelfScript)
	tx.AddTxOut(txOut1)

	// 2nd output: the payment to the recipient
	if !cancelTx {
		pkScript, err := txscript.PayToAddrScript(to)
		if err != nil {
			return err
		}
		txOut2 := wire.NewTxOut(amountSats, pkScript)
		tx.AddTxOut(txOut2)
	} else {
		// send the amount to TSS self if tx is cancelled
		remainingSats += amountSats
	}

	// 3rd output: the remaining btc to TSS self
	if remainingSats >= constant.BTCWithdrawalDustAmount {
		txOut3 := wire.NewTxOut(remainingSats, payToSelfScript)
		tx.AddTxOut(txOut3)
	}
	return nil
}

// SignTx signs the given tx with TSS
func (signer *Signer) SignTx(
	ctx context.Context,
	tx *wire.MsgTx,
	inputAmounts []int64,
	height uint64,
	nonce uint64,
) error {
	pkScript, err := signer.TSSToPkScript()
	if err != nil {
		return err
	}

	// calculate sighashes to sign
	sigHashes := txscript.NewTxSigHashes(tx, txscript.NewCannedPrevOutputFetcher([]byte{}, 0))
	witnessHashes := make([][]byte, len(tx.TxIn))
	for ix := range tx.TxIn {
		amount := inputAmounts[ix]
		witnessHashes[ix], err = txscript.CalcWitnessSigHash(pkScript, sigHashes, txscript.SigHashAll, tx, ix, amount)
		if err != nil {
			return err
		}
	}

	// sign the tx with TSS
	sig65Bs, err := signer.TSS().SignBatch(ctx, witnessHashes, height, nonce, signer.Chain().ChainId)
	if err != nil {
		return fmt.Errorf("SignBatch failed: %v", err)
	}

	for ix := range tx.TxIn {
		sig65B := sig65Bs[ix]
		R := &btcec.ModNScalar{}
		R.SetBytes((*[32]byte)(sig65B[:32]))
		S := &btcec.ModNScalar{}
		S.SetBytes((*[32]byte)(sig65B[32:64]))
		sig := btcecdsa.NewSignature(R, S)

		pkCompressed := signer.TSS().PubKey().Bytes(true)
		hashType := txscript.SigHashAll
		txWitness := wire.TxWitness{append(sig.Serialize(), byte(hashType)), pkCompressed}
		tx.TxIn[ix].Witness = txWitness
	}

	return nil
}
