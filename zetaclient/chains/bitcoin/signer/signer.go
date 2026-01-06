// Package signer implements the ChainSigner interface for BTC
package signer

import (
	"bytes"
	"context"
	"encoding/hex"
	"runtime/debug"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/pkg/retry"
	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/client"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/mode"
)

const (
	// broadcastBackoff is the backoff duration for retrying broadcast
	broadcastBackoff = time.Second * 6

	// broadcastRetries is the maximum number of retries for broadcasting a transaction
	broadcastRetries = 10
)

type BitcoinClient interface {
	GetNetworkInfo(context.Context) (*btcjson.GetNetworkInfoResult, error)
	GetRawTransaction(context.Context, *chainhash.Hash) (*btcutil.Tx, error)
	GetEstimatedFeeRate(_ context.Context, confTarget int64) (uint64, error)
	GetMempoolTxsAndFees(_ context.Context, childHash string) (client.MempoolTxsAndFees, error)

	// This is a mutating function that does not get called when zetaclient is in dry-mode.
	SendRawTransaction(_ context.Context, _ *wire.MsgTx, allowHighFees bool) (*chainhash.Hash, error)
}

// Signer deals with signing & broadcasting BTC transactions.
type Signer struct {
	*base.Signer
	bitcoinClient BitcoinClient
	isRegnet      bool
}

// New creates a new Bitcoin signer
func New(baseSigner *base.Signer, bitcoinClient BitcoinClient) *Signer {
	chainID := baseSigner.Chain().ChainId
	return &Signer{
		Signer:        baseSigner,
		bitcoinClient: bitcoinClient,
		isRegnet:      chains.IsBitcoinRegnet(chainID),
	}
}

// Broadcast sends the signed transaction to the network
func (signer *Signer) Broadcast(ctx context.Context, signedTx *wire.MsgTx) error {
	var outBuff bytes.Buffer
	if err := signedTx.Serialize(&outBuff); err != nil {
		return errors.Wrap(err, "unable to serialize tx")
	}

	signer.Logger().Std.Info().
		Stringer(logs.FieldTx, signedTx.TxHash()).
		Str("signer_tx_payload", hex.EncodeToString(outBuff.Bytes())).
		Msg("broadcasting transaction")

	_, err := signer.bitcoinClient.SendRawTransaction(ctx, signedTx, true)
	if err != nil {
		return errors.Wrap(err, "unable to broadcast raw tx")
	}

	return nil
}

// TryProcessOutbound signs and broadcasts a BTC transaction from a new outbound
func (signer *Signer) TryProcessOutbound(
	ctx context.Context,
	cctx *types.CrossChainTx,
	observer *observer.Observer,
	height uint64,
) {
	outboundID := base.OutboundIDFromCCTX(cctx)
	signer.MarkOutbound(outboundID, true)

	// end outbound process on panic
	defer func() {
		if err := recover(); err != nil {
			signer.Logger().Std.Error().
				Str(logs.FieldCctxIndex, cctx.Index).
				Interface("panic", err).
				Str("stack_trace", string(debug.Stack())).
				Msg("caught panic error")
		}

		signer.MarkOutbound(outboundID, false)
	}()

	// prepare logger
	params := cctx.GetCurrentOutboundParam()
	lf := map[string]any{
		logs.FieldCctxIndex: cctx.Index,
		logs.FieldNonce:     params.TssNonce,
	}
	signerAddress, err := observer.ZetaRepo().GetKeysAddress()
	if err != nil {
		return
	}
	lf["signer"] = signerAddress
	logger := signer.Logger().Std.With().Fields(lf).Logger()

	// query network info to get minRelayFee (typically 1000 satoshis)
	networkInfo, err := signer.bitcoinClient.GetNetworkInfo(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed get bitcoin network info")
		return
	}
	minRelayFee := networkInfo.RelayFee
	if minRelayFee <= 0 {
		logger.Error().Float64("min_relay_fee", minRelayFee).Msg("invalid minimum relay fee")
		return
	}

	// compliance check
	isRestricted := !signer.PassesCompliance(cctx)

	// setup outbound data
	txData, err := NewOutboundData(cctx, height, minRelayFee, isRestricted, logger)
	if err != nil {
		logger.Error().Err(err).Msg("failed to setup Bitcoin outbound data")
		return
	}

	var (
		signedTx       *wire.MsgTx
		stuckTx, found = observer.LastStuckOutbound()
		rbfTx          = found && stuckTx.Nonce == txData.nonce
	)

	if signer.ClientMode.IsDryMode() {
		logger.Info().Stringer(logs.FieldMode, mode.DryMode).Msg("skipping outbound processing")
		return
	}

	// sign outbound
	if !rbfTx {
		// sign withdraw tx
		signedTx, err = signer.SignWithdrawTx(ctx, txData, observer)
		if err != nil {
			logger.Error().Err(err).Msg("call to SignWithdrawTx failed")
			return
		}
		logger.Info().Str(logs.FieldTx, signedTx.TxID()).Msg("call to SignWithdrawTx succeed")
	} else {
		// sign RBF tx
		signedTx, err = signer.SignRBFTx(ctx, txData, stuckTx.Tx)
		if err != nil {
			logger.Error().Err(err).Msg("call to SignRBFTx failed")
			return
		}
		logger.Info().Str(logs.FieldTx, signedTx.TxID()).Msg("call to SignRBFTx succeed")
	}

	// broadcast signed outbound
	signer.BroadcastOutbound(ctx, signedTx, params.TssNonce, rbfTx, cctx, observer)
}

// BroadcastOutbound sends the signed transaction to the Bitcoin network
func (signer *Signer) BroadcastOutbound(
	ctx context.Context,
	tx *wire.MsgTx,
	nonce uint64,
	rbfTx bool,
	cctx *types.CrossChainTx,
	ob *observer.Observer,
) {
	txHash := tx.TxID()

	// prepare logger fields
	logger := signer.Logger().Std.With().
		Uint64(logs.FieldNonce, nonce).
		Str(logs.FieldTx, txHash).
		Str(logs.FieldCctxIndex, cctx.Index).
		Logger()

	// double check to ensure the tx being replaced is still the last outbound.
	// when CCTX gets stuck at nonce 'N', the pending nonce will stop incrementing
	// and stay at 'N' or 'N+1' (at most).
	if rbfTx && ob.GetPendingNonce() > nonce+1 {
		logger.Warn().Msg("RBF tx nonce is outdated; skipping broadcasting")
		return
	}

	// try broacasting tx with backoff in case of RPC error
	broadcast := func() error {
		return retry.Retry(signer.Broadcast(ctx, tx))
	}

	bo := backoff.NewConstantBackOff(broadcastBackoff)
	boWithMaxRetries := backoff.WithMaxRetries(bo, broadcastRetries)
	err := retry.DoWithBackoff(broadcast, boWithMaxRetries)
	if err != nil {
		logger.Error().Err(err).Msg("unable to broadcast Bitcoin outbound")
		return
	}
	logger.Info().Msg("broadcasted Bitcoin outbound successfully")

	// save tx local db and ignore db error.
	// db error is not critical and should not block outbound tracker.
	if err := ob.SaveBroadcastedTx(txHash, nonce); err != nil {
		logger.Error().Err(err).Msg("unable to save broadcasted Bitcoin outbound")
	}

	deadline := time.Now().Add(4 * constant.ZetaBlockTime)
	trackerBo := backoff.NewConstantBackOff(constant.ZetaBlockTime)

	call := func() error {
		_, err := ob.ZetaRepo().PostOutboundTracker(ctx, logger, nonce, txHash)
		return err
	}

	// Block the post outbound tracker call to make sure it gets posted before proceeding.
	if err := retry.DoWithDeadline(call, trackerBo, deadline); err != nil {
		logger.Error().Err(err).Msg("failed to post outbound tracker")
	}

	_, included := ob.TryIncludeOutbound(ctx, cctx, txHash)
	if included {
		logger.Info().Msg("included newly broadcasted Bitcoin outbound")
	}
}
