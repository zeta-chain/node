// Package signer implements the ChainSigner interface for BTC
package signer

import (
	"bytes"
	"context"
	"encoding/hex"
	"time"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	ethcommon "github.com/ethereum/go-ethereum/common"

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

var _ interfaces.ChainSigner = (*Signer)(nil)

// Signer deals with signing BTC transactions and implements the ChainSigner interface
type Signer struct {
	*base.Signer

	// client is the RPC client to interact with the Bitcoin chain
	client interfaces.BTCRPCClient
}

// NewSigner creates a new Bitcoin signer
func NewSigner(
	chain chains.Chain,
	client interfaces.BTCRPCClient,
	tss interfaces.TSSSigner,
	logger base.Logger,
) *Signer {
	// create base signer
	baseSigner := base.NewSigner(chain, tss, logger)

	return &Signer{
		Signer: baseSigner,
		client: client,
	}
}

// TODO: get rid of below four get/set functions for Bitcoin, as they are not needed in future
// https://github.com/zeta-chain/node/issues/2532
// SetZetaConnectorAddress does nothing for BTC
func (signer *Signer) SetZetaConnectorAddress(_ ethcommon.Address) {
}

// SetERC20CustodyAddress does nothing for BTC
func (signer *Signer) SetERC20CustodyAddress(_ ethcommon.Address) {
}

// GetZetaConnectorAddress returns dummy address
func (signer *Signer) GetZetaConnectorAddress() ethcommon.Address {
	return ethcommon.Address{}
}

// GetERC20CustodyAddress returns dummy address
func (signer *Signer) GetERC20CustodyAddress() ethcommon.Address {
	return ethcommon.Address{}
}

// SetGatewayAddress does nothing for BTC
// Note: TSS address will be used as gateway address for Bitcoin
func (signer *Signer) SetGatewayAddress(_ string) {
}

// GetGatewayAddress returns empty address
// Note: same as SetGatewayAddress
func (signer *Signer) GetGatewayAddress() string {
	return ""
}

// PkScriptTSS returns the TSS pkScript
func (signer *Signer) PkScriptTSS() ([]byte, error) {
	tssAddrP2WPKH, err := signer.TSS().PubKey().AddressBTC(signer.Chain().ChainId)
	if err != nil {
		return nil, err
	}
	return txscript.PayToAddrScript(tssAddrP2WPKH)
}

// Broadcast sends the signed transaction to the network
func (signer *Signer) Broadcast(signedTx *wire.MsgTx) error {
	var outBuff bytes.Buffer
	err := signedTx.Serialize(&outBuff)
	if err != nil {
		return err
	}
	str := hex.EncodeToString(outBuff.Bytes())

	_, err = signer.client.SendRawTransaction(signedTx, true)
	if err != nil {
		return err
	}

	signer.Logger().Std.Info().Msgf("Broadcasted transaction data: %s ", str)
	return nil
}

// TryProcessOutbound signs and broadcasts a BTC transaction from a new outbound
func (signer *Signer) TryProcessOutbound(
	ctx context.Context,
	cctx *types.CrossChainTx,
	outboundProcessor *outboundprocessor.Processor,
	outboundID string,
	chainObserver interfaces.ChainObserver,
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

	// convert chain observer to BTC observer
	btcObserver, ok := chainObserver.(*observer.Observer)
	if !ok {
		logger.Error().Msg("chain observer is not a bitcoin observer")
		return
	}

	// query network info to get minRelayFee (typically 1000 satoshis)
	networkInfo, err := signer.client.GetNetworkInfo()
	if err != nil {
		logger.Error().Err(err).Msgf("failed get bitcoin network info")
		return
	}
	minRelayFee := networkInfo.RelayFee

	var (
		signedTx *wire.MsgTx
		stuckTx  = btcObserver.GetLastStuckOutbound()
	)

	// sign outbound
	if stuckTx != nil && params.TssNonce == stuckTx.Nonce {
		signedTx, err = signer.SignRBFTx(ctx, cctx, stuckTx.Tx, minRelayFee)
		if err != nil {
			logger.Error().Err(err).Msg("SignRBFTx failed")
			return
		}
		logger.Info().Msg("SignRBFTx success")
	} else {
		// setup outbound data
		txData, err := NewOutboundData(cctx, chain.ChainId, height, minRelayFee, logger, signer.Logger().Compliance)
		if err != nil {
			logger.Error().Err(err).Msg("failed to setup Bitcoin outbound data")
			return
		}

		// sign withdraw tx
		signedTx, err = signer.SignWithdrawTx(ctx, txData, btcObserver)
		if err != nil {
			logger.Error().Err(err).Msg("SignWithdrawTx failed")
			return
		}
		logger.Info().Msg("SignWithdrawTx success")
	}

	// broadcast signed outbound
	signer.BroadcastOutbound(ctx, signedTx, params.TssNonce, cctx, btcObserver, zetacoreClient)
}

// BroadcastOutbound sends the signed transaction to the Bitcoin network
func (signer *Signer) BroadcastOutbound(
	ctx context.Context,
	tx *wire.MsgTx,
	nonce uint64,
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

	// try broacasting tx with increasing backoff (1s, 2s, 4s, 8s, 16s) in case of RPC error
	backOff := broadcastBackoff
	for i := 0; i < broadcastRetries; i++ {
		time.Sleep(backOff)

		// broadcast tx
		err := signer.Broadcast(tx)
		if err != nil {
			logger.Warn().Err(err).Fields(lf).Msgf("broadcasting Bitcoin outbound, retry %d", i)
			backOff *= 2
			continue
		}
		logger.Info().Fields(lf).Msgf("broadcasted Bitcoin outbound successfully")

		// save tx local db
		ob.SaveBroadcastedTx(txHash, nonce)

		// add tx to outbound tracker so that all observers know about it
		zetaHash, err := zetacoreClient.PostOutboundTracker(ctx, ob.Chain().ChainId, nonce, txHash)
		if err != nil {
			logger.Err(err).Fields(lf).Msgf("unable to add Bitcoin outbound tracker")
		} else {
			lf[logs.FieldZetaTx] = zetaHash
			logger.Info().Fields(lf).Msgf("add Bitcoin outbound tracker successfully")
		}

		// try including this outbound as early as possible
		_, included := ob.TryIncludeOutbound(ctx, cctx, txHash)
		if included {
			logger.Info().Fields(lf).Msgf("included newly broadcasted Bitcoin outbound")
		}

		// successful broadcast; no need to retry
		break
	}
}
