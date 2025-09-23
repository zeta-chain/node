package observer

import (
	"context"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

const (
	// maxNoAssetCallExcessAmount is the maximum excess amount beyond depositor fee for Bitcoin NoAssetCall
	// A static value 100K satoshis (0.001 BTC) is used to prevent accidental funds loss exceeding $110 at the time of writing
	maxNoAssetCallExcessAmount = 100_000
)

// ObserveInbound observes the Bitcoin chain for inbounds and post votes to zetacore
func (ob *Observer) ObserveInbound(ctx context.Context) error {
	logger := ob.Logger().Inbound

	// keep last block up-to-date
	if err := ob.updateLastBlock(ctx); err != nil {
		return err
	}

	// scan SAFE confirmed blocks
	startBlockSafe, endBlockSafe := ob.GetScanRangeInboundSafe(config.MaxBlocksPerScan)
	if startBlockSafe < endBlockSafe {
		// observe inbounds for the block range [startBlock, endBlock-1]
		lastScannedNew, err := ob.observeInboundInBlockRange(ctx, startBlockSafe, endBlockSafe-1)
		if err != nil {
			logger.Error().
				Err(err).
				Uint64("from", startBlockSafe).
				Uint64("to", endBlockSafe-1).
				Msg("error observing inbounds in block range")
		}

		// save last scanned block to both memory and db
		if lastScannedNew > ob.LastBlockScanned() {
			logger.Info().
				Uint64("from", startBlockSafe).
				Uint64("to", lastScannedNew).
				Msg("observed blocks for inbounds")
			if err := ob.SaveLastBlockScanned(lastScannedNew); err != nil {
				return errors.Wrapf(err, "unable to save last scanned Bitcoin block %d", lastScannedNew)
			}
		}
	}

	// scan FAST confirmed blocks if available
	_, endBlockFast := ob.GetScanRangeInboundFast(config.MaxBlocksPerScan)
	if endBlockSafe < endBlockFast {
		_, err := ob.observeInboundInBlockRange(ctx, endBlockSafe, endBlockFast-1)
		if err != nil {
			logger.Error().
				Err(err).
				Uint64("from", endBlockSafe).
				Uint64("to", endBlockFast-1).
				Msg("error observing inbounds in block range (fast)")
		}
	}

	return nil
}

// observeInboundInBlockRange observes inbounds for given block range [startBlock, toBlock (inclusive)]
// It returns the last successfully scanned block height, so the caller knows where to resume next time
func (ob *Observer) observeInboundInBlockRange(ctx context.Context, startBlock, toBlock uint64) (uint64, error) {
	for blockNumber := startBlock; blockNumber <= toBlock; blockNumber++ {
		// query incoming gas asset to TSS address
		// #nosec G115 always in range
		res, err := ob.GetBlockByNumberCached(ctx, int64(blockNumber))
		if err != nil {
			// we have to re-scan this block next time
			return blockNumber - 1, errors.Wrapf(err, "error getting bitcoin block %d", blockNumber)
		}

		if len(res.Block.Tx) <= 1 {
			continue
		}

		// filter incoming txs to TSS address
		tssAddress := ob.TSSAddressString()

		// #nosec G115 always positive
		events, err := FilterAndParseIncomingTx(
			ctx,
			ob.bitcoinClient,
			res.Block.Tx,
			uint64(res.Block.Height),
			tssAddress,
			ob.logger.Inbound,
			ob.netParams,
		)
		if err != nil {
			// we have to re-scan this block next time
			return blockNumber - 1, errors.Wrapf(err, "error filtering incoming txs for block %d", blockNumber)
		}

		// post inbound vote message to zetacore
		for _, event := range events {
			msg := ob.GetInboundVoteFromBtcEvent(event)
			if msg != nil {
				// skip early observed inbound that is not eligible for fast confirmation
				if msg.ConfirmationMode == crosschaintypes.ConfirmationMode_FAST {
					eligible, err := ob.IsInboundEligibleForFastConfirmation(ctx, msg)
					if err != nil {
						return blockNumber - 1, errors.Wrapf(
							err,
							"unable to determine inbound fast confirmation eligibility for tx %s",
							event.TxHash,
						)
					}
					if !eligible {
						continue
					}
				}

				_, err = ob.PostVoteInbound(ctx, msg, zetacore.PostVoteInboundExecutionGasLimit)
				if err != nil {
					// we have to re-scan this block next time
					return blockNumber - 1, errors.Wrapf(err, "error posting inbound vote for tx %s", event.TxHash)
				}
			}
		}
	}

	// successful processed all blocks in [startBlock, toBlock]
	return toBlock, nil
}

// FilterAndParseIncomingTx given txs list returned by the "getblock 2" RPC command, return the txs that are relevant to us
// relevant tx must have the following vouts as the first two vouts:
// vout0: p2wpkh to the TSS address (targetAddress)
// vout1: OP_RETURN memo, base64 encoded
func FilterAndParseIncomingTx(
	ctx context.Context,
	bitcoinClient BitcoinClient,
	txs []btcjson.TxRawResult,
	blockNumber uint64,
	tssAddress string,
	logger zerolog.Logger,
	netParams *chaincfg.Params,
) ([]*BTCInboundEvent, error) {
	events := make([]*BTCInboundEvent, 0)

	for idx, tx := range txs {
		if idx == 0 {
			// the first tx is coinbase; we do not process coinbase tx
			continue
		}

		event, err := GetBtcEventWithWitness(
			ctx,
			bitcoinClient,
			tx,
			tssAddress,
			blockNumber,
			logger,
			netParams,
			common.CalcDepositorFee,
		)
		if err != nil {
			// unable to parse the tx, the caller should retry
			return nil, errors.Wrapf(err, "error getting btc event for tx %s in block %d", tx.Txid, blockNumber)
		}

		if event != nil {
			events = append(events, event)
		}
	}

	return events, nil
}

// GetInboundVoteFromBtcEvent converts a BTCInboundEvent to a MsgVoteInbound to enable voting on the inbound on zetacore
//
// Returns:
//   - a valid MsgVoteInbound message, or
//   - nil if no valid message can be created for whatever reasons:
//     invalid data, not processable, invalid amount, etc.
func (ob *Observer) GetInboundVoteFromBtcEvent(event *BTCInboundEvent) *crosschaintypes.MsgVoteInbound {
	// prepare logger
	logger := ob.logger.Inbound.With().Str(logs.FieldTx, event.TxHash).Logger()

	// decode event memo bytes
	// if the memo is invalid, we set the status in the event, the inbound will be observed as invalid
	if err := event.DecodeMemoBytes(ob.Chain().ChainId); err != nil {
		logger.Info().
			Err(err).
			Str("memo", hex.EncodeToString(event.MemoBytes)).
			Msg("invalid memo")
		event.SetStatusAndErrMessage(crosschaintypes.InboundStatus_INVALID_MEMO, err.Error())
	}

	// check if the event is processable
	if !ob.IsEventProcessable(*event) {
		return nil
	}

	// resolve the amount to be used in inbound vote message
	if err := event.ResolveAmountForMsgVoteInbound(); err != nil {
		// should never happen, otherwise skip this tx
		logger.Error().Err(err).Msg("unable to resolve msg vote amount")
		return nil
	}

	return ob.NewInboundVoteFromEvent(event)
}

// NewInboundVoteFromEvent creates a MsgVoteInbound message from a inbound event
func (ob *Observer) NewInboundVoteFromEvent(event *BTCInboundEvent) *crosschaintypes.MsgVoteInbound {
	// determine confirmation mode
	confirmationMode := crosschaintypes.ConfirmationMode_FAST
	if ob.IsBlockConfirmedForInboundSafe(event.BlockNumber) {
		confirmationMode = crosschaintypes.ConfirmationMode_SAFE
	}

	return crosschaintypes.NewMsgVoteInbound(
		ob.ZetacoreClient().GetKeys().GetOperatorAddress().String(),
		event.FromAddress,
		ob.Chain().ChainId,
		event.FromAddress,
		event.ToAddress,
		ob.ZetacoreClient().Chain().ChainId,
		event.AmountForMsgVoteInbound,
		event.Message(),
		event.TxHash,
		event.BlockNumber,
		0,
		event.CoinType(),
		"",
		0,
		crosschaintypes.ProtocolContractVersion_V2,
		false, // no arbitrary call for deposit to ZetaChain
		event.Status,
		confirmationMode,
		crosschaintypes.WithRevertOptions(event.RevertOptions()),
		crosschaintypes.WithCrossChainCall(event.IsCrossChainCall()),
		crosschaintypes.WithErrorMessage(event.ErrorMessage),
	)
}
