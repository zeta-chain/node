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

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

// ObserveInbound observes the Bitcoin chain for inbounds and post votes to zetacore
// TODO(revamp): simplify this function into smaller functions
func (ob *Observer) ObserveInbound(ctx context.Context) error {
	// get and update latest block height
	currentBlock, err := ob.rpc.GetBlockCount(ctx)
	if err != nil {
		return fmt.Errorf("observeInboundBTC: error getting block number: %s", err)
	}
	if currentBlock < 0 {
		return fmt.Errorf("observeInboundBTC: block number is negative: %d", currentBlock)
	}

	// 0 will be returned if the node is not synced
	if currentBlock == 0 {
		ob.nodeEnabled.Store(false)
		ob.logger.Inbound.Debug().Err(err).Msg("WatchInbound: Bitcoin node is not enabled")
		return nil
	}

	ob.nodeEnabled.Store(true)

	// #nosec G115 checked positive
	lastBlock := uint64(currentBlock)
	if lastBlock < ob.LastBlock() {
		return fmt.Errorf(
			"observeInboundBTC: block number should not decrease: current %d last %d",
			currentBlock,
			ob.LastBlock(),
		)
	}
	ob.WithLastBlock(lastBlock)

	// check confirmation
	blockNumber := ob.LastBlockScanned() + 1
	if !ob.IsBlockConfirmed(blockNumber) {
		return nil
	}

	// query incoming gas asset to TSS address
	// #nosec G115 always in range
	res, err := ob.GetBlockByNumberCached(ctx, int64(blockNumber))
	if err != nil {
		ob.logger.Inbound.Error().Err(err).Msgf("observeInboundBTC: error getting bitcoin block %d", blockNumber)
		return err
	}
	ob.logger.Inbound.Info().Msgf("observeInboundBTC: block %d has %d txs, current block %d, last scanned %d",
		blockNumber, len(res.Block.Tx), currentBlock, blockNumber-1)

	// add block header to zetacore
	if len(res.Block.Tx) > 1 {
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
			ob.logger.Inbound.Error().
				Err(err).
				Msgf("observeInboundBTC: error filtering incoming txs for block %d", blockNumber)
			return err // we have to re-scan this block next time
		}

		// post inbound vote message to zetacore
		for _, event := range events {
			msg := ob.GetInboundVoteFromBtcEvent(event)
			if msg != nil {
				_, err = ob.PostVoteInbound(ctx, msg, zetacore.PostVoteInboundExecutionGasLimit)
				if err != nil {
					return errors.Wrapf(err, "error PostVoteInbound") // we have to re-scan this block next time
				}
			}
		}
	}

	// save last scanned block to both memory and db
	if err := ob.SaveLastBlockScanned(blockNumber); err != nil {
		ob.logger.Inbound.Error().
			Err(err).
			Msgf("observeInboundBTC: error writing last scanned block %d to db", blockNumber)
	}

	return nil
}

// ProcessInboundTrackers processes inbound trackers
// TODO(revamp): move inbound tracker logic in a specific file
func (ob *Observer) ProcessInboundTrackers(ctx context.Context) error {
	trackers, err := ob.ZetacoreClient().GetInboundTrackersForChain(ctx, ob.Chain().ChainId)
	if err != nil {
		return err
	}

	for _, tracker := range trackers {
		ob.logger.Inbound.Info().
			Str("tracker.hash", tracker.TxHash).
			Str("tracker.coin-type", tracker.CoinType.String()).
			Msgf("checking tracker")
		ballotIdentifier, err := ob.CheckReceiptForBtcTxHash(ctx, tracker.TxHash, true)
		if err != nil {
			return err
		}
		ob.logger.Inbound.Info().
			Str("inbound.chain", ob.Chain().Name).
			Str("inbound.ballot", ballotIdentifier).
			Str("inbound.coin-type", coin.CoinType_Gas.String()).
			Msgf("Vote submitted for inbound Tracker")
	}

	return nil
}

// CheckReceiptForBtcTxHash checks the receipt for a btc tx hash
func (ob *Observer) CheckReceiptForBtcTxHash(ctx context.Context, txHash string, vote bool) (string, error) {
	hash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return "", err
	}

	tx, err := ob.rpc.GetRawTransactionVerbose(ctx, hash)
	if err != nil {
		return "", err
	}

	blockHash, err := chainhash.NewHashFromStr(tx.BlockHash)
	if err != nil {
		return "", err
	}

	blockVb, err := ob.rpc.GetBlockVerbose(ctx, blockHash)
	if err != nil {
		return "", err
	}

	if len(blockVb.Tx) <= 1 {
		return "", fmt.Errorf("block %d has no transactions", blockVb.Height)
	}

	tss, err := ob.ZetacoreClient().GetBTCTSSAddress(ctx, ob.Chain().ChainId)
	if err != nil {
		return "", err
	}

	// check confirmation
	// #nosec G115 block height always positive
	if !ob.IsBlockConfirmed(uint64(blockVb.Height)) {
		return "", fmt.Errorf("block %d is not confirmed yet", blockVb.Height)
	}

	// #nosec G115 always positive
	event, err := GetBtcEvent(
		ctx,
		ob.rpc,
		*tx,
		tss,
		uint64(blockVb.Height),
		ob.logger.Inbound,
		ob.netParams,
		common.CalcDepositorFee,
	)
	if err != nil {
		return "", err
	}

	if event == nil {
		return "", errors.New("no btc deposit event found")
	}

	msg := ob.GetInboundVoteFromBtcEvent(event)
	if msg == nil {
		return "", errors.New("no message built for btc sent to TSS")
	}

	if !vote {
		return msg.Digest(), nil
	}

	return ob.PostVoteInbound(ctx, msg, zetacore.PostVoteInboundExecutionGasLimit)
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
