package observer

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/coin"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/types"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

// WatchInbound watches Bitcoin chain for inbounds on a ticker
// It starts a ticker and run ObserveInbound
// TODO(revamp): move all ticker related methods in the same file
func (ob *Observer) WatchInbound(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	ticker, err := types.NewDynamicTicker("Bitcoin_WatchInbound", ob.ChainParams().InboundTicker)
	if err != nil {
		ob.logger.Inbound.Error().Err(err).Msg("error creating ticker")
		return err
	}
	defer ticker.Stop()

	ob.logger.Inbound.Info().Msgf("WatchInbound started for chain %d", ob.Chain().ChainId)
	sampledLogger := ob.logger.Inbound.Sample(&zerolog.BasicSampler{N: 10})

	// ticker loop
	for {
		select {
		case <-ticker.C():
			if !app.IsInboundObservationEnabled() {
				sampledLogger.Info().
					Msgf("WatchInbound: inbound observation is disabled for chain %d", ob.Chain().ChainId)
				continue
			}
			err := ob.ObserveInbound(ctx)
			if err != nil {
				// skip showing log for block number 0 as it means Bitcoin node is not enabled
				// TODO: prevent this routine from running if Bitcoin node is not enabled
				// https://github.com/zeta-chain/node/issues/2790
				if !errors.Is(err, bitcoin.ErrBitcoinNotEnabled) {
					ob.logger.Inbound.Error().Err(err).Msg("WatchInbound error observing in tx")
				} else {
					ob.logger.Inbound.Debug().Err(err).Msg("WatchInbound: Bitcoin node is not enabled")
				}
			}
			ticker.UpdateInterval(ob.ChainParams().InboundTicker, ob.logger.Inbound)
		case <-ob.StopChannel():
			ob.logger.Inbound.Info().Msgf("WatchInbound stopped for chain %d", ob.Chain().ChainId)
			return nil
		}
	}
}

// ObserveInbound observes the Bitcoin chain for inbounds and post votes to zetacore
// TODO(revamp): simplify this function into smaller functions
func (ob *Observer) ObserveInbound(ctx context.Context) error {
	// get and update latest block height
	currentBlock, err := ob.btcClient.GetBlockCount()
	if err != nil {
		return fmt.Errorf("observeInboundBTC: error getting block number: %s", err)
	}
	if currentBlock < 0 {
		return fmt.Errorf("observeInboundBTC: block number is negative: %d", currentBlock)
	}

	// 0 will be returned if the node is not synced
	if currentBlock == 0 {
		return errors.Wrap(bitcoin.ErrBitcoinNotEnabled, "observeInboundBTC: current block number 0 is too low")
	}

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
	res, err := ob.GetBlockByNumberCached(int64(blockNumber))
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
			ob.btcClient,
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

// WatchInboundTracker watches zetacore for bitcoin inbound trackers
// TODO(revamp): move all ticker related methods in the same file
func (ob *Observer) WatchInboundTracker(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	ticker, err := types.NewDynamicTicker("Bitcoin_WatchInboundTracker", ob.ChainParams().InboundTicker)
	if err != nil {
		ob.logger.Inbound.Err(err).Msg("error creating ticker")
		return err
	}

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			if !app.IsInboundObservationEnabled() {
				continue
			}
			err := ob.ProcessInboundTrackers(ctx)
			if err != nil {
				ob.logger.Inbound.Error().
					Err(err).
					Msgf("error observing inbound tracker for chain %d", ob.Chain().ChainId)
			}
			ticker.UpdateInterval(ob.ChainParams().InboundTicker, ob.logger.Inbound)
		case <-ob.StopChannel():
			ob.logger.Inbound.Info().Msgf("WatchInboundTracker stopped for chain %d", ob.Chain().ChainId)
			return nil
		}
	}
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

	tx, err := ob.btcClient.GetRawTransactionVerbose(hash)
	if err != nil {
		return "", err
	}

	blockHash, err := chainhash.NewHashFromStr(tx.BlockHash)
	if err != nil {
		return "", err
	}

	blockVb, err := ob.btcClient.GetBlockVerboseTx(blockHash)
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
	if !ob.IsBlockConfirmed(uint64(blockVb.Height)) {
		return "", fmt.Errorf("block %d is not confirmed yet", blockVb.Height)
	}

	// calculate depositor fee
	depositorFee, err := bitcoin.CalcDepositorFee(ob.btcClient, tx, ob.netParams)
	if err != nil {
		return "", errors.Wrapf(err, "error calculating depositor fee for inbound %s", tx.Txid)
	}

	// #nosec G115 always positive
	event, err := GetBtcEvent(
		ob.btcClient,
		*tx,
		tss,
		uint64(blockVb.Height),
		ob.logger.Inbound,
		ob.netParams,
		depositorFee,
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

	zetaHash, ballot, err := ob.ZetacoreClient().PostVoteInbound(
		ctx,
		zetacore.PostVoteInboundGasLimit,
		zetacore.PostVoteInboundExecutionGasLimit,
		msg,
	)
	if err != nil {
		ob.logger.Inbound.Error().Err(err).Msg("error posting to zetacore")
		return "", err
	} else if zetaHash != "" {
		ob.logger.Inbound.Info().Msgf("BTC deposit detected and reported: PostVoteInbound zeta tx hash: %s inbound %s ballot %s fee %v",
			zetaHash, txHash, ballot, event.DepositorFee)
	}

	return msg.Digest(), nil
}

// FilterAndParseIncomingTx given txs list returned by the "getblock 2" RPC command, return the txs that are relevant to us
// relevant tx must have the following vouts as the first two vouts:
// vout0: p2wpkh to the TSS address (targetAddress)
// vout1: OP_RETURN memo, base64 encoded
func FilterAndParseIncomingTx(
	rpcClient interfaces.BTCRPCClient,
	txs []btcjson.TxRawResult,
	blockNumber uint64,
	tssAddress string,
	logger zerolog.Logger,
	netParams *chaincfg.Params,
) ([]*BTCInboundEvent, error) {
	events := make([]*BTCInboundEvent, 0)
	for idx, tx := range txs {
		if idx == 0 {
			continue // the first tx is coinbase; we do not process coinbase tx
		}

		// calculate depositor fee
		depositorFee, err := bitcoin.CalcDepositorFee(rpcClient, &txs[idx], netParams)
		if err != nil {
			return nil, errors.Wrapf(err, "error calculating depositor fee for inbound %s", tx.Txid)
		}

		event, err := GetBtcEvent(rpcClient, tx, tssAddress, blockNumber, logger, netParams, depositorFee)
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
func (ob *Observer) GetInboundVoteFromBtcEvent(event *BTCInboundEvent) *crosschaintypes.MsgVoteInbound {
	// prepare logger fields
	lf := map[string]any{
		logs.FieldModule: logs.ModNameInbound,
		logs.FieldMethod: "GetInboundVoteMessageFromBtcEvent",
		logs.FieldChain:  ob.Chain().ChainId,
		logs.FieldTx:     event.TxHash,
	}

	// decode event memo bytes
	err := ob.DecodeEventMemoBytes(event)
	if err != nil {
		ob.Logger().Inbound.Info().Fields(lf).Msgf("invalid memo bytes: %s", hex.EncodeToString(event.MemoBytes))
		return nil
	}

	// check if the event is processable
	if !ob.CheckEventProcessability(event) {
		return nil
	}

	// convert the amount to integer (satoshis)
	amount := big.NewFloat(event.Value)
	amount = amount.Mul(amount, big.NewFloat(btcutil.SatoshiPerBitcoin))
	amountInt, _ := amount.Int(nil)

	// create inbound vote message contract V1 for legacy memo
	if event.MemoStd == nil {
		return ob.NewInboundVoteV1(event, amountInt)
	}

	// create inbound vote message for standard memo
	return ob.NewInboundVoteMemoStd(event, amountInt)
}

// GetBtcEvent returns a valid BTCInboundEvent or nil
// it uses witness data to extract the sender address, except for mainnet
func GetBtcEvent(
	rpcClient interfaces.BTCRPCClient,
	tx btcjson.TxRawResult,
	tssAddress string,
	blockNumber uint64,
	logger zerolog.Logger,
	netParams *chaincfg.Params,
	depositorFee float64,
) (*BTCInboundEvent, error) {
	if netParams.Name == chaincfg.MainNetParams.Name {
		return GetBtcEventWithoutWitness(rpcClient, tx, tssAddress, blockNumber, logger, netParams, depositorFee)
	}
	return GetBtcEventWithWitness(rpcClient, tx, tssAddress, blockNumber, logger, netParams, depositorFee)
}

// GetBtcEventWithoutWitness either returns a valid BTCInboundEvent or nil
// Note: the caller should retry the tx on error (e.g., GetSenderAddressByVin failed)
// TODO(revamp): simplify this function
func GetBtcEventWithoutWitness(
	rpcClient interfaces.BTCRPCClient,
	tx btcjson.TxRawResult,
	tssAddress string,
	blockNumber uint64,
	logger zerolog.Logger,
	netParams *chaincfg.Params,
	depositorFee float64,
) (*BTCInboundEvent, error) {
	found := false
	var value float64
	var memo []byte
	if len(tx.Vout) >= 2 {
		// 1st vout must have tss address as receiver with p2wpkh scriptPubKey
		vout0 := tx.Vout[0]
		script := vout0.ScriptPubKey.Hex
		if len(script) == 44 && script[:4] == "0014" {
			// P2WPKH output: 0x00 + 20 bytes of pubkey hash
			receiver, err := bitcoin.DecodeScriptP2WPKH(vout0.ScriptPubKey.Hex, netParams)
			if err != nil { // should never happen
				return nil, err
			}

			// skip irrelevant tx to us
			if receiver != tssAddress {
				return nil, nil
			}

			// deposit amount has to be no less than the minimum depositor fee
			if vout0.Value < depositorFee {
				logger.Info().
					Msgf("GetBtcEvent: btc deposit amount %v in txid %s is less than depositor fee %v", vout0.Value, tx.Txid, depositorFee)
				return nil, nil
			}
			value = vout0.Value - depositorFee

			// 2nd vout must be a valid OP_RETURN memo
			vout1 := tx.Vout[1]
			memo, found, err = bitcoin.DecodeOpReturnMemo(vout1.ScriptPubKey.Hex)
			if err != nil {
				logger.Error().Err(err).Msgf("GetBtcEvent: error decoding OP_RETURN memo: %s", vout1.ScriptPubKey.Hex)
				return nil, nil
			}
		}
	}
	// event found, get sender address
	if found {
		if len(tx.Vin) == 0 { // should never happen
			return nil, fmt.Errorf("GetBtcEvent: no input found for inbound: %s", tx.Txid)
		}

		// get sender address by input (vin)
		fromAddress, err := GetSenderAddressByVin(rpcClient, tx.Vin[0], netParams)
		if err != nil {
			return nil, errors.Wrapf(err, "error getting sender address for inbound: %s", tx.Txid)
		}

		// skip this tx and move on (e.g., due to unknown script type)
		// we don't know whom to refund if this tx gets reverted in zetacore
		if fromAddress == "" {
			return nil, nil
		}

		return &BTCInboundEvent{
			FromAddress:  fromAddress,
			ToAddress:    tssAddress,
			Value:        value,
			DepositorFee: depositorFee,
			MemoBytes:    memo,
			BlockNumber:  blockNumber,
			TxHash:       tx.Txid,
		}, nil
	}
	return nil, nil
}

// GetSenderAddressByVin get the sender address from the transaction input (vin)
func GetSenderAddressByVin(rpcClient interfaces.BTCRPCClient, vin btcjson.Vin, net *chaincfg.Params) (string, error) {
	// query previous raw transaction by txid
	hash, err := chainhash.NewHashFromStr(vin.Txid)
	if err != nil {
		return "", err
	}

	// this requires running bitcoin node with 'txindex=1'
	tx, err := rpcClient.GetRawTransaction(hash)
	if err != nil {
		return "", errors.Wrapf(err, "error getting raw transaction %s", vin.Txid)
	}

	// #nosec G115 - always in range
	if len(tx.MsgTx().TxOut) <= int(vin.Vout) {
		return "", fmt.Errorf("vout index %d out of range for tx %s", vin.Vout, vin.Txid)
	}

	// decode sender address from previous pkScript
	pkScript := tx.MsgTx().TxOut[vin.Vout].PkScript

	return bitcoin.DecodeSenderFromScript(pkScript, net)
}
