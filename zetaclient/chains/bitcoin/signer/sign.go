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
	"github.com/zeta-chain/node/zetaclient/logs"
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
	logger := signer.Logger().Std.With().Uint64(logs.FieldNonce, txData.nonce).Logger()

	nonceMark := chains.NonceMarkAmount(txData.nonce)

	// we don't know how many UTXOs will be used beforehand, so we do
	// a conservative estimation using the maximum size of the outbound tx:
	// estimateFee = feeRate * maxTxSize
	// #nosec G115 always in range
	estimateFee := float64(int64(txData.feeRate)*common.OutboundBytesMax) / 1e8
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
		return nil, errors.Wrapf(err, "FetchUTXOs failed for nonce %d", txData.nonce)
	}

	// select N UTXOs to cover the total expense
	selected, err := ob.SelectUTXOs(
		ctx,
		totalAmount,
		MaxNoOfInputsPerTx,
		txData.nonce,
		consolidationRank,
	)
	if err != nil {
		return nil, err
	}

	// build tx and add inputs
	tx := wire.NewMsgTx(wire.TxVersion)
	inAmounts, err := AddTxInputs(tx, selected.UTXOs)
	if err != nil {
		return nil, err
	}

	// size checking
	// #nosec G115 always positive
	txSize, err := common.EstimateOutboundSize(int64(len(selected.UTXOs)), []btcutil.Address{txData.to})
	if err != nil {
		return nil, err
	}
	if txSize > common.OutboundBytesMax {
		// in case of accident
		logger.Warn().
			Int64("tx_size", txSize).
			Int64("outbound_bytes_max", common.OutboundBytesMax).
			Msg("tx size is greater than outbound bytes max")
		txSize = common.OutboundBytesMax
	}

	// fee calculation
	// #nosec G115 always in range
	fees := txSize * int64(txData.feeRate)

	// add tx outputs
	inputValue := selected.Value
	err = signer.AddWithdrawTxOutputs(tx,
		txData.to, inputValue, txData.amountSats, nonceMark, fees, txData.cancelTx)
	if err != nil {
		return nil, err
	}
	logger.Info().
		Uint64("tx_rate", txData.feeRate).
		Int64("tx_fees", fees).
		Uint16("tx_consolidated_utxos", selected.ConsolidatedUTXOs).
		Float64("tx_consolidated_value", selected.ConsolidatedValue).
		Msg("signing bitcoin outbound")

	// sign the tx
	err = signer.SignTx(ctx, tx, inAmounts, txData.height, txData.nonce)
	if err != nil {
		return nil, errors.Wrap(err, "SignTx failed")
	}

	return tx, nil
}

// AddTxInputs adds the inputs to the tx and returns input amounts
func AddTxInputs(tx *wire.MsgTx, utxos []btcjson.ListUnspentResult) ([]int64, error) {
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
//
// Note: float64 is used for 'inputValue' because UTXOs struct uses float64.
// But we need to use 'int64' for the outputs because NewTxOut expects int64.
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
		signer.Logger().Std.Info().
			Int64("remaining_sats", remainingSats).
			Msg("adjust remainder value to avoid duplicate nonce-mark")
		remainingSats--
	}

	// 1st output: the nonce-mark btc to TSS self
	payToSelfScript, err := signer.TSS().PubKey().BTCPayToAddrScript(signer.Chain().ChainId)
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
	pkScript, err := signer.TSS().PubKey().BTCPayToAddrScript(signer.Chain().ChainId)
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
		return errors.Wrap(err, "SignBatch failed")
	}

	// add witnesses to the tx
	pkCompressed := signer.TSS().PubKey().Bytes(true)
	hashType := txscript.SigHashAll
	for ix := range tx.TxIn {
		sig65B := sig65Bs[ix]
		R := &btcec.ModNScalar{}
		R.SetBytes((*[32]byte)(sig65B[:32]))
		S := &btcec.ModNScalar{}
		S.SetBytes((*[32]byte)(sig65B[32:64]))
		sig := btcecdsa.NewSignature(R, S)

		txWitness := wire.TxWitness{append(sig.Serialize(), byte(hashType)), pkCompressed}
		tx.TxIn[ix].Witness = txWitness
	}

	return nil
}
