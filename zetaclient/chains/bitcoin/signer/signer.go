// Package signer implements the ChainSigner interface for BTC
package signer

import (
	"bytes"
	"context"
	"encoding/hex"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/outboundprocessor"
)

const (
	// broadcastBackoff is the initial backoff duration for retrying broadcast
	broadcastBackoff = 1000 * time.Millisecond

	// broadcastRetries is the maximum number of retries for broadcasting a transaction
	broadcastRetries = 5
)

type RPC interface {
	GetNetworkInfo(ctx context.Context) (*btcjson.GetNetworkInfoResult, error)
	GetRawTransaction(ctx context.Context, hash *chainhash.Hash) (*btcutil.Tx, error)
	GetEstimatedFeeRate(ctx context.Context, confTarget int64, regnet bool) (int64, error)
	SendRawTransaction(ctx context.Context, tx *wire.MsgTx, allowHighFees bool) (*chainhash.Hash, error)
	GetTotalMempoolParentsSizeNFees(
		ctx context.Context,
		childHash string,
		timeout time.Duration,
	) (int64, float64, int64, int64, error)
}

// Signer deals with signing & broadcasting BTC transactions.
type Signer struct {
	*base.Signer
	rpc RPC
}

// New creates a new Bitcoin signer
func New(chain chains.Chain, tss interfaces.TSSSigner, rpc RPC, logger base.Logger) *Signer {
	return &Signer{
		Signer: base.NewSigner(chain, tss, logger),
		rpc:    rpc,
	}
}

// Broadcast sends the signed transaction to the network
func (signer *Signer) Broadcast(ctx context.Context, signedTx *wire.MsgTx) error {
	var outBuff bytes.Buffer
	if err := signedTx.Serialize(&outBuff); err != nil {
		return errors.Wrap(err, "unable to serialize tx")
	}

	signer.Logger().Std.Info().
		Str(logs.FieldTx, signedTx.TxHash().String()).
		Str("signer.tx_payload", hex.EncodeToString(outBuff.Bytes())).
		Msg("Broadcasting transaction")

	_, err := signer.rpc.SendRawTransaction(ctx, signedTx, true)
	if err != nil {
		return errors.Wrap(err, "unable to broadcast raw tx")
	}

	return nil
}

// PkScriptTSS returns the TSS pkScript
func (signer *Signer) PkScriptTSS() ([]byte, error) {
	tssAddrP2WPKH, err := signer.TSS().PubKey().AddressBTC(signer.Chain().ChainId)
	if err != nil {
		return nil, err
	}
	return txscript.PayToAddrScript(tssAddrP2WPKH)
}

// TryProcessOutbound signs and broadcasts a BTC transaction from a new outbound
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
	chain := signer.Chain()
	params := cctx.GetCurrentOutboundParam()
	lf := map[string]any{
		logs.FieldMethod: "TryProcessOutbound",
		logs.FieldCctx:   cctx.Index,
		logs.FieldNonce:  params.TssNonce,
	}
	signerAddress, err := zetacoreClient.GetKeys().GetAddress()
	if err == nil {
		lf["signer"] = signerAddress.String()
	}
	logger := signer.Logger().Std.With().Fields(lf).Logger()

	// query network info to get minRelayFee (typically 1000 satoshis)
	networkInfo, err := signer.rpc.GetNetworkInfo(ctx)
	if err != nil {
		logger.Error().Err(err).Msgf("failed get bitcoin network info")
		return
	}
	minRelayFee := networkInfo.RelayFee
	if minRelayFee <= 0 {
		logger.Error().Msgf("invalid minimum relay fee: %f", minRelayFee)
		return
	}

	var (
		rbfTx    = false
		signedTx *wire.MsgTx
		stuckTx  = observer.GetLastStuckOutbound()
	)

	// sign outbound
	if stuckTx != nil && params.TssNonce == stuckTx.Nonce {
		// sign RBF tx
		rbfTx = true
		signedTx, err = signer.SignRBFTx(ctx, cctx, height, stuckTx.Tx, minRelayFee)
		if err != nil {
			logger.Error().Err(err).Msg("SignRBFTx failed")
			return
		}
		logger.Info().Str(logs.FieldTx, signedTx.TxID()).Msg("SignRBFTx succeed")
	} else {
		// setup outbound data
		txData, err := NewOutboundData(cctx, chain.ChainId, height, minRelayFee, logger, signer.Logger().Compliance)
		if err != nil {
			logger.Error().Err(err).Msg("failed to setup Bitcoin outbound data")
			return
		}

		// sign withdraw tx
		signedTx, err = signer.SignWithdrawTx(ctx, txData, observer)
		if err != nil {
			logger.Error().Err(err).Msg("SignWithdrawTx failed")
			return
		}
		logger.Info().Str(logs.FieldTx, signedTx.TxID()).Msg("SignWithdrawTx succeed")
	}

	// broadcast signed outbound
	signer.BroadcastOutbound(ctx, signedTx, params.TssNonce, rbfTx, cctx, observer, zetacoreClient)
}

// BroadcastOutbound sends the signed transaction to the Bitcoin network
func (signer *Signer) BroadcastOutbound(
	ctx context.Context,
	tx *wire.MsgTx,
	nonce uint64,
	rbfTx bool,
	cctx *types.CrossChainTx,
	ob *observer.Observer,
	zetacoreClient interfaces.ZetacoreClient,
) {
	txHash := tx.TxID()

	// prepare logger fields
	lf := map[string]any{
		logs.FieldMethod: "broadcastOutbound",
		logs.FieldNonce:  nonce,
		logs.FieldTx:     txHash,
		logs.FieldCctx:   cctx.Index,
	}
	logger := signer.Logger().Std

	// double check to ensure the tx is still the last outbound
	if rbfTx {
		if ob.GetPendingNonce() > nonce+1 {
			logger.Warn().Fields(lf).Msgf("RBF tx nonce is outdated, skipping broadcast")
			return
		}
	}

	// try broacasting tx with increasing backoff (1s, 2s, 4s, 8s, 16s) in case of RPC error
	backOff := broadcastBackoff
	for i := 0; i < broadcastRetries; i++ {
		time.Sleep(backOff)

		// broadcast tx
		err := signer.Broadcast(ctx, tx)
		if err != nil {
			logger.Warn().Err(err).Fields(lf).Msgf("broadcasting Bitcoin outbound, retry %d", i)
			backOff *= 2
			continue
		}
		logger.Info().Fields(lf).Msg("broadcasted Bitcoin outbound successfully")

		// save tx local db
		ob.SaveBroadcastedTx(txHash, nonce)

		// add tx to outbound tracker so that all observers know about it
		zetaHash, err := zetacoreClient.PostOutboundTracker(ctx, ob.Chain().ChainId, nonce, txHash)
		if err != nil {
			logger.Err(err).Fields(lf).Msg("unable to add Bitcoin outbound tracker")
		} else {
			lf[logs.FieldZetaTx] = zetaHash
			logger.Info().Fields(lf).Msg("add Bitcoin outbound tracker successfully")
		}

		// try including this outbound as early as possible
		_, included := ob.TryIncludeOutbound(ctx, cctx, txHash)
		if included {
			logger.Info().Fields(lf).Msg("included newly broadcasted Bitcoin outbound")
		}

		// successful broadcast; no need to retry
		break
	}
}
