// Package signer implements the ChainSigner interface for BTC
package signer

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	btcecdsa "github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/compliance"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/outboundprocessor"
)

const (
	// the maximum number of inputs per outbound
	MaxNoOfInputsPerTx = 20

	// the rank below (or equal to) which we consolidate UTXOs
	consolidationRank = 10

	// broadcastBackoff is the initial backoff duration for retrying broadcast
	broadcastBackoff = 1000 * time.Millisecond

	// broadcastRetries is the maximum number of retries for broadcasting a transaction
	broadcastRetries = 5
)

// Signer deals with signing BTC transactions and implements the ChainSigner interface
type Signer struct {
	*base.Signer

	// client is the RPC client to interact with the Bitcoin chain
	client interfaces.BTCRPCClient
}

// NewSigner creates a new Bitcoin signer
func NewSigner(
	chain chains.Chain,
	tss interfaces.TSSSigner,
	logger base.Logger,
	cfg config.BTCConfig,
) (*Signer, error) {
	// create base signer
	baseSigner := base.NewSigner(chain, tss, logger)

	// create the bitcoin rpc client using the provided config
	connCfg := &rpcclient.ConnConfig{
		Host:         cfg.RPCHost,
		User:         cfg.RPCUsername,
		Pass:         cfg.RPCPassword,
		HTTPPostMode: true,
		DisableTLS:   true,
		Params:       cfg.RPCParams,
	}
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create bitcoin rpc client")
	}

	return &Signer{
		Signer: baseSigner,
		client: client,
	}, nil
}

// AddWithdrawTxOutputs adds the 3 outputs to the withdraw tx
// 1st output: the nonce-mark btc to TSS itself
// 2nd output: the payment to the recipient
// 3rd output: the remaining btc to TSS itself
func (signer *Signer) AddWithdrawTxOutputs(
	tx *wire.MsgTx,
	to btcutil.Address,
	total float64,
	amount float64,
	nonceMark int64,
	fees *big.Int,
	cancelTx bool,
) error {
	// convert withdraw amount to satoshis
	amountSatoshis, err := common.GetSatoshis(amount)
	if err != nil {
		return err
	}

	// calculate remaining btc (the change) to TSS self
	remaining := total - amount
	remainingSats, err := common.GetSatoshis(remaining)
	if err != nil {
		return err
	}
	remainingSats -= fees.Int64()
	remainingSats -= nonceMark
	if remainingSats < 0 {
		return fmt.Errorf("remainder value is negative: %d", remainingSats)
	} else if remainingSats == nonceMark {
		signer.Logger().Std.Info().Msgf("adjust remainder value to avoid duplicate nonce-mark: %d", remainingSats)
		remainingSats--
	}

	// 1st output: the nonce-mark btc to TSS self
	tssAddrP2WPKH, err := signer.TSS().PubKey().AddressBTC(signer.Chain().ChainId)
	if err != nil {
		return err
	}
	payToSelfScript, err := txscript.PayToAddrScript(tssAddrP2WPKH)
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
		txOut2 := wire.NewTxOut(amountSatoshis, pkScript)
		tx.AddTxOut(txOut2)
	} else {
		// send the amount to TSS self if tx is cancelled
		remainingSats += amountSatoshis
	}

	// 3rd output: the remaining btc to TSS self
	if remainingSats > 0 {
		txOut3 := wire.NewTxOut(remainingSats, payToSelfScript)
		tx.AddTxOut(txOut3)
	}
	return nil
}

// SignWithdrawTx receives utxos sorted by value, amount in BTC, feeRate in BTC per Kb
// TODO(revamp): simplify the function
func (signer *Signer) SignWithdrawTx(
	ctx context.Context,
	to btcutil.Address,
	amount float64,
	gasPrice *big.Int,
	sizeLimit uint64,
	observer *observer.Observer,
	height uint64,
	nonce uint64,
	chain chains.Chain,
	cancelTx bool,
) (*wire.MsgTx, error) {
	estimateFee := float64(gasPrice.Uint64()*common.OutboundBytesMax) / 1e8
	nonceMark := chains.NonceMarkAmount(nonce)

	// refresh unspent UTXOs and continue with keysign regardless of error
	err := observer.FetchUTXOs(ctx)
	if err != nil {
		signer.Logger().
			Std.Error().
			Err(err).
			Msgf("SignGasWithdraw: FetchUTXOs error: nonce %d chain %d", nonce, chain.ChainId)
	}

	// select N UTXOs to cover the total expense
	prevOuts, total, consolidatedUtxo, consolidatedValue, err := observer.SelectUTXOs(
		ctx,
		amount+estimateFee+float64(nonceMark)*1e-8,
		MaxNoOfInputsPerTx,
		nonce,
		consolidationRank,
		false,
	)
	if err != nil {
		return nil, err
	}

	// build tx with selected unspents
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

	// size checking
	// #nosec G115 always positive
	txSize, err := common.EstimateOutboundSize(uint64(len(prevOuts)), []btcutil.Address{to})
	if err != nil {
		return nil, err
	}
	if sizeLimit < common.BtcOutboundBytesWithdrawer { // ZRC20 'withdraw' charged less fee from end user
		signer.Logger().Std.Info().
			Msgf("sizeLimit %d is less than BtcOutboundBytesWithdrawer %d for nonce %d", sizeLimit, txSize, nonce)
	}
	if txSize < common.OutboundBytesMin { // outbound shouldn't be blocked a low sizeLimit
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
	fees := new(big.Int).Mul(big.NewInt(int64(txSize)), gasPrice)
	signer.Logger().
		Std.Info().
		Msgf("bitcoin outbound nonce %d gasPrice %s size %d fees %s consolidated %d utxos of value %v",
			nonce, gasPrice.String(), txSize, fees.String(), consolidatedUtxo, consolidatedValue)

	// add tx outputs
	err = signer.AddWithdrawTxOutputs(tx, to, total, amount, nonceMark, fees, cancelTx)
	if err != nil {
		return nil, err
	}

	// sign the tx
	sigHashes := txscript.NewTxSigHashes(tx, txscript.NewCannedPrevOutputFetcher([]byte{}, 0))
	witnessHashes := make([][]byte, len(tx.TxIn))
	for ix := range tx.TxIn {
		amt, err := common.GetSatoshis(prevOuts[ix].Amount)
		if err != nil {
			return nil, err
		}
		pkScript, err := hex.DecodeString(prevOuts[ix].ScriptPubKey)
		if err != nil {
			return nil, err
		}
		witnessHashes[ix], err = txscript.CalcWitnessSigHash(pkScript, sigHashes, txscript.SigHashAll, tx, ix, amt)
		if err != nil {
			return nil, err
		}
	}

	sig65Bs, err := signer.TSS().SignBatch(ctx, witnessHashes, height, nonce, chain.ChainId)
	if err != nil {
		return nil, fmt.Errorf("SignBatch error: %v", err)
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

	return tx, nil
}

// Broadcast sends the signed transaction to the network
func (signer *Signer) Broadcast(signedTx *wire.MsgTx) error {
	fmt.Printf("BTCSigner: Broadcasting: %s\n", signedTx.TxHash().String())

	var outBuff bytes.Buffer
	err := signedTx.Serialize(&outBuff)
	if err != nil {
		return err
	}
	str := hex.EncodeToString(outBuff.Bytes())
	fmt.Printf("BTCSigner: Transaction Data: %s\n", str)

	hash, err := signer.client.SendRawTransaction(signedTx, true)
	if err != nil {
		return err
	}

	signer.Logger().Std.Info().Msgf("Broadcasting BTC tx , hash %s ", hash)
	return nil
}

// TryProcessOutbound signs and broadcasts a BTC transaction from a new outbound
// TODO(revamp): simplify the function
func (signer *Signer) TryProcessOutbound(
	ctx context.Context,
	cctx *types.CrossChainTx,
	outboundProcessor *outboundprocessor.Processor,
	outboundID string,
	observer *observer.Observer,
	zetacoreClient interfaces.ZetacoreClient,
	height uint64,
) {
	// end outbound process on panic
	defer func() {
		outboundProcessor.EndTryProcess(outboundID)
		if err := recover(); err != nil {
			signer.Logger().Std.Error().Msgf("BTC TryProcessOutbound: %s, caught panic error: %v", cctx.Index, err)
		}
	}()

	// prepare logger
	params := cctx.GetCurrentOutboundParam()
	// prepare logger fields
	lf := map[string]any{
		logs.FieldMethod: "TryProcessOutbound",
		logs.FieldCctx:   cctx.Index,
		logs.FieldNonce:  params.TssNonce,
	}
	logger := signer.Logger().Std.With().Fields(lf).Logger()

	// support gas token only for Bitcoin outbound
	coinType := cctx.InboundParams.CoinType
	if coinType == coin.CoinType_Zeta || coinType == coin.CoinType_ERC20 {
		logger.Error().Msg("can only send BTC to a BTC network")
		return
	}

	chain := observer.Chain()
	outboundTssNonce := params.TssNonce
	signerAddress, err := zetacoreClient.GetKeys().GetAddress()
	if err != nil {
		logger.Error().Err(err).Msg("cannot get signer address")
		return
	}
	lf["signer"] = signerAddress.String()

	// get size limit and gas price
	sizelimit := params.CallOptions.GasLimit
	gasprice, ok := new(big.Int).SetString(params.GasPrice, 10)
	if !ok || gasprice.Cmp(big.NewInt(0)) < 0 {
		logger.Error().Msgf("cannot convert gas price  %s ", params.GasPrice)
		return
	}

	// Check receiver P2WPKH address
	to, err := chains.DecodeBtcAddress(params.Receiver, params.ReceiverChainId)
	if err != nil {
		logger.Error().Err(err).Msgf("cannot decode address %s ", params.Receiver)
		return
	}
	if !chains.IsBtcAddressSupported(to) {
		logger.Error().Msgf("unsupported address %s", params.Receiver)
		return
	}
	amount := float64(params.Amount.Uint64()) / 1e8

	// Add 1 satoshi/byte to gasPrice to avoid minRelayTxFee issue
	networkInfo, err := signer.client.GetNetworkInfo()
	if err != nil {
		logger.Error().Err(err).Msgf("cannot get bitcoin network info")
		return
	}
	satPerByte := common.FeeRateToSatPerByte(networkInfo.RelayFee)
	gasprice.Add(gasprice, satPerByte)

	// compliance check
	restrictedCCTX := compliance.IsCctxRestricted(cctx)
	if restrictedCCTX {
		compliance.PrintComplianceLog(logger, signer.Logger().Compliance,
			true, chain.ChainId, cctx.Index, cctx.InboundParams.Sender, params.Receiver, "BTC")
	}

	// check dust amount
	dustAmount := params.Amount.Uint64() < constant.BTCWithdrawalDustAmount
	if dustAmount {
		logger.Warn().Msgf("dust amount %d sats, canceling tx", params.Amount.Uint64())
	}

	// set the amount to 0 when the tx should be cancelled
	cancelTx := restrictedCCTX || dustAmount
	if cancelTx {
		amount = 0.0
	}

	// sign withdraw tx
	tx, err := signer.SignWithdrawTx(
		ctx,
		to,
		amount,
		gasprice,
		sizelimit,
		observer,
		height,
		outboundTssNonce,
		chain,
		cancelTx,
	)
	if err != nil {
		logger.Warn().Err(err).Msg("SignWithdrawTx failed")
		return
	}
	logger.Info().Msg("Key-sign success")

	// FIXME: add prometheus metrics
	_, err = zetacoreClient.GetObserverList(ctx)
	if err != nil {
		logger.Warn().
			Err(err).Stringer("observation_type", observertypes.ObservationType_OutboundTx).
			Msg("unable to get observer list, observation")
	}
	if tx != nil {
		outboundHash := tx.TxHash().String()
		lf[logs.FieldTx] = outboundHash

		// try broacasting tx with increasing backoff (1s, 2s, 4s, 8s, 16s) in case of RPC error
		backOff := broadcastBackoff
		for i := 0; i < broadcastRetries; i++ {
			time.Sleep(backOff)
			err := signer.Broadcast(tx)
			if err != nil {
				logger.Warn().Err(err).Fields(lf).Msgf("Broadcasting Bitcoin tx, retry %d", i)
				backOff *= 2
				continue
			}
			logger.Info().Fields(lf).Msgf("Broadcast Bitcoin tx successfully")
			zetaHash, err := zetacoreClient.PostOutboundTracker(
				ctx,
				chain.ChainId,
				outboundTssNonce,
				outboundHash,
			)
			if err != nil {
				logger.Err(err).Fields(lf).Msgf("Unable to add Bitcoin outbound tracker")
			}
			lf[logs.FieldZetaTx] = zetaHash
			logger.Info().Fields(lf).Msgf("Add Bitcoin outbound tracker successfully")

			// Save successfully broadcasted transaction to btc chain observer
			observer.SaveBroadcastedTx(outboundHash, outboundTssNonce)

			break // successful broadcast; no need to retry
		}
	}
}
