package observer

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

// ObserveInbound observes the Bitcoin chain for inbounds and post votes to zetacore
func (ob *Observer) ObserveInbound(ctx context.Context) error {
	logger := ob.Logger().Inbound.With().Str(logs.FieldMethod, "observe_inbound").Logger()

	// keep last block up-to-date
	if err := ob.updateLastBlock(ctx); err != nil {
		return err
	}

	// get the block range to scan
	startBlock, endBlock := ob.GetScanRangeInboundSafe(config.MaxBlocksPerScan)
	if startBlock >= endBlock {
		return nil
	}

	// observe inbounds for the block range [startBlock, endBlock-1]
	lastScannedNew, err := ob.observeInboundInBlockRange(ctx, startBlock, endBlock-1)
	if err != nil {
		logger.Error().
			Err(err).
			Uint64("from", startBlock).
			Uint64("to", endBlock-1).
			Msg("error observing inbounds in block range")
	}

	// save last scanned block to both memory and db
	if lastScannedNew > ob.LastBlockScanned() {
		logger.Info().Uint64("from", startBlock).Uint64("to", lastScannedNew).Msg("observed blocks for inbounds")
		if err := ob.SaveLastBlockScanned(lastScannedNew); err != nil {
			return errors.Wrapf(err, "unable to save last scanned Bitcoin block %d", lastScannedNew)
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
			ob.rpc,
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
	rpc RPC,
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

		event, err := GetBtcEvent(ctx, rpc, tx, tssAddress, blockNumber, logger, netParams, common.CalcDepositorFee)
		if err != nil {
			// unable to parse the tx, the caller should retry
			return nil, errors.Wrapf(err, "error getting btc event for tx %s in block %d", tx.Txid, blockNumber)
		}

		if event != nil {
			events = append(events, event)
			logger.Info().Msgf("FilterAndParseIncomingTx: found btc event for tx %s in block %d", tx.Txid, blockNumber)
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
func (ob *Observer) GetInboundVoteFromBtcEvent(event *BTCInboundEvent) *types.MsgVoteInbound {
	// prepare logger fields
	lf := map[string]any{
		logs.FieldMethod: "GetInboundVoteFromBtcEvent",
		logs.FieldTx:     event.TxHash,
	}

	// decode event memo bytes
	err := event.DecodeMemoBytes(ob.Chain().ChainId)
	if err != nil {
		ob.Logger().Inbound.Info().Fields(lf).Msgf("invalid memo bytes: %s", hex.EncodeToString(event.MemoBytes))
		return nil
	}

	// check if the event is processable
	if !ob.IsEventProcessable(*event) {
		return nil
	}

	// convert the amount to integer (satoshis)
	amountSats, err := common.GetSatoshis(event.Value)
	if err != nil {
		ob.Logger().Inbound.Error().Err(err).Fields(lf).Msgf("can't convert value %f to satoshis", event.Value)
		return nil
	}
	amountInt := big.NewInt(amountSats)

	// create inbound vote message contract V1 for legacy memo
	if event.MemoStd == nil {
		return ob.NewInboundVoteFromLegacyMemo(event, amountInt)
	}

	// create inbound vote message for standard memo
	return ob.NewInboundVoteFromStdMemo(event, amountInt)
}

// GetBtcEvent returns a valid BTCInboundEvent or nil
// it uses witness data to extract the sender address
func GetBtcEvent(
	ctx context.Context,
	rpc RPC,
	tx btcjson.TxRawResult,
	tssAddress string,
	blockNumber uint64,
	logger zerolog.Logger,
	netParams *chaincfg.Params,
	feeCalculator common.DepositorFeeCalculator,
) (*BTCInboundEvent, error) {
	return GetBtcEventWithWitness(ctx, rpc, tx, tssAddress, blockNumber, logger, netParams, feeCalculator)
}

// GetSenderAddressByVin get the sender address from the transaction input (vin)
func GetSenderAddressByVin(
	ctx context.Context,
	rpc RPC,
	vin btcjson.Vin,
	net *chaincfg.Params,
) (string, error) {
	// query previous raw transaction by txid
	hash, err := chainhash.NewHashFromStr(vin.Txid)
	if err != nil {
		return "", err
	}

	// this requires running bitcoin node with 'txindex=1'
	tx, err := rpc.GetRawTransaction(ctx, hash)
	if err != nil {
		return "", errors.Wrapf(err, "error getting raw transaction %s", vin.Txid)
	}

	// #nosec G115 - always in range
	if len(tx.MsgTx().TxOut) <= int(vin.Vout) {
		return "", fmt.Errorf("vout index %d out of range for tx %s", vin.Vout, vin.Txid)
	}

	// decode sender address from previous pkScript
	pkScript := tx.MsgTx().TxOut[vin.Vout].PkScript

	return common.DecodeSenderFromScript(pkScript, net)
}
