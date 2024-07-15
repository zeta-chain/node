package observer

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	cosmosmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/compliance"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	zctx "github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/types"
	"github.com/zeta-chain/zetacore/zetaclient/zetacore"
)

// WatchInbound watches Bitcoin chain for inbounds on a ticker
// It starts a ticker and run ObserveInbound
// TODO(revamp): move all ticker related methods in the same file
func (ob *Observer) WatchInbound(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	ticker, err := types.NewDynamicTicker("Bitcoin_WatchInbound", ob.GetChainParams().InboundTicker)
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
			if !app.IsInboundObservationEnabled(ob.GetChainParams()) {
				sampledLogger.Info().
					Msgf("WatchInbound: inbound observation is disabled for chain %d", ob.Chain().ChainId)
				continue
			}
			err := ob.ObserveInbound(ctx)
			if err != nil {
				ob.logger.Inbound.Error().Err(err).Msg("WatchInbound error observing in tx")
			}
			ticker.UpdateInterval(ob.GetChainParams().InboundTicker, ob.logger.Inbound)
		case <-ob.StopChannel():
			ob.logger.Inbound.Info().Msgf("WatchInbound stopped for chain %d", ob.Chain().ChainId)
			return nil
		}
	}
}

// ObserveInbound observes the Bitcoin chain for inbounds and post votes to zetacore
// TODO(revamp): simplify this function into smaller functions
func (ob *Observer) ObserveInbound(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	zetaCoreClient := ob.ZetacoreClient()

	// get and update latest block height
	cnt, err := ob.btcClient.GetBlockCount()
	if err != nil {
		return fmt.Errorf("observeInboundBTC: error getting block number: %s", err)
	}
	if cnt < 0 {
		return fmt.Errorf("observeInboundBTC: block number is negative: %d", cnt)
	}
	// #nosec G115 checked positive
	lastBlock := uint64(cnt)
	if lastBlock < ob.LastBlock() {
		return fmt.Errorf(
			"observeInboundBTC: block number should not decrease: current %d last %d",
			cnt,
			ob.LastBlock(),
		)
	}
	ob.WithLastBlock(lastBlock)

	// skip if current height is too low
	if lastBlock < ob.GetChainParams().ConfirmationCount {
		return fmt.Errorf("observeInboundBTC: skipping observer, current block number %d is too low", cnt)
	}

	// skip if no new block is confirmed
	lastScanned := ob.LastBlockScanned()
	if lastScanned >= lastBlock-ob.GetChainParams().ConfirmationCount {
		return nil
	}

	// query incoming gas asset to TSS address
	blockNumber := lastScanned + 1
	// #nosec G115 always in range
	res, err := ob.GetBlockByNumberCached(int64(blockNumber))
	if err != nil {
		ob.logger.Inbound.Error().Err(err).Msgf("observeInboundBTC: error getting bitcoin block %d", blockNumber)
		return err
	}
	ob.logger.Inbound.Info().Msgf("observeInboundBTC: block %d has %d txs, current block %d, last block %d",
		blockNumber, len(res.Block.Tx), cnt, lastScanned)

	// add block header to zetacore
	// TODO: consider having a separate ticker(from TSS scaning) for posting block headers
	// https://github.com/zeta-chain/node/issues/1847
	// TODO: move this logic in its own routine
	// https://github.com/zeta-chain/node/issues/2204
	blockHeaderVerification, found := app.GetBlockHeaderEnabledChains(ob.Chain().ChainId)
	if found && blockHeaderVerification.Enabled {
		// #nosec G115 always in range
		err = ob.postBlockHeader(ctx, int64(blockNumber))
		if err != nil {
			ob.logger.Inbound.Warn().Err(err).Msgf("observeInboundBTC: error posting block header %d", blockNumber)
		}
	}

	if len(res.Block.Tx) > 1 {
		// get depositor fee
		depositorFee := bitcoin.CalcDepositorFee(res.Block, ob.Chain().ChainId, ob.netParams, ob.logger.Inbound)

		// filter incoming txs to TSS address
		tssAddress := ob.TSS().BTCAddress()

		// #nosec G115 always positive
		inbounds, err := FilterAndParseIncomingTx(
			ob.btcClient,
			res.Block.Tx,
			uint64(res.Block.Height),
			tssAddress,
			ob.logger.Inbound,
			ob.netParams,
			depositorFee,
		)
		if err != nil {
			ob.logger.Inbound.Error().
				Err(err).
				Msgf("observeInboundBTC: error filtering incoming txs for block %d", blockNumber)
			return err // we have to re-scan this block next time
		}

		// post inbound vote message to zetacore
		for _, inbound := range inbounds {
			msg := ob.GetInboundVoteMessageFromBtcEvent(inbound)
			if msg != nil {
				zetaHash, ballot, err := zetaCoreClient.PostVoteInbound(
					ctx,
					zetacore.PostVoteInboundGasLimit,
					zetacore.PostVoteInboundExecutionGasLimit,
					msg,
				)
				if err != nil {
					ob.logger.Inbound.Error().
						Err(err).
						Msgf("observeInboundBTC: error posting to zetacore for tx %s", inbound.TxHash)
					return err // we have to re-scan this block next time
				} else if zetaHash != "" {
					ob.logger.Inbound.Info().Msgf("observeInboundBTC: PostVoteInbound zeta tx hash: %s inbound %s ballot %s fee %v",
						zetaHash, inbound.TxHash, ballot, depositorFee)
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

	ticker, err := types.NewDynamicTicker("Bitcoin_WatchInboundTracker", ob.GetChainParams().InboundTicker)
	if err != nil {
		ob.logger.Inbound.Err(err).Msg("error creating ticker")
		return err
	}

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			if !app.IsInboundObservationEnabled(ob.GetChainParams()) {
				continue
			}
			err := ob.ProcessInboundTrackers(ctx)
			if err != nil {
				ob.logger.Inbound.Error().
					Err(err).
					Msgf("error observing inbound tracker for chain %d", ob.Chain().ChainId)
			}
			ticker.UpdateInterval(ob.GetChainParams().InboundTicker, ob.logger.Inbound)
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
			Msgf("checking tracker with hash :%s and coin-type :%s ", tracker.TxHash, tracker.CoinType)
		ballotIdentifier, err := ob.CheckReceiptForBtcTxHash(ctx, tracker.TxHash, true)
		if err != nil {
			return err
		}
		ob.logger.Inbound.Info().
			Msgf("Vote submitted for inbound Tracker, Chain : %s,Ballot Identifier : %s, coin-type %s", ob.Chain().ChainName, ballotIdentifier, coin.CoinType_Gas.String())
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

	depositorFee := bitcoin.CalcDepositorFee(blockVb, ob.Chain().ChainId, ob.netParams, ob.logger.Inbound)
	tss, err := ob.ZetacoreClient().GetBTCTSSAddress(ctx, ob.Chain().ChainId)
	if err != nil {
		return "", err
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

	msg := ob.GetInboundVoteMessageFromBtcEvent(event)
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
			zetaHash, txHash, ballot, depositorFee)
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
	depositorFee float64,
) ([]*BTCInboundEvent, error) {
	inbounds := make([]*BTCInboundEvent, 0)
	for idx, tx := range txs {
		if idx == 0 {
			continue // the first tx is coinbase; we do not process coinbase tx
		}

		inbound, err := GetBtcEvent(rpcClient, tx, tssAddress, blockNumber, logger, netParams, depositorFee)
		if err != nil {
			// unable to parse the tx, the caller should retry
			return nil, errors.Wrapf(err, "error getting btc event for tx %s in block %d", tx.Txid, blockNumber)
		}

		if inbound != nil {
			inbounds = append(inbounds, inbound)
			logger.Info().Msgf("FilterAndParseIncomingTx: found btc event for tx %s in block %d", tx.Txid, blockNumber)
		}
	}
	return inbounds, nil
}

// GetInboundVoteMessageFromBtcEvent converts a BTCInboundEvent to a MsgVoteInbound to enable voting on the inbound on zetacore
func (ob *Observer) GetInboundVoteMessageFromBtcEvent(inbound *BTCInboundEvent) *crosschaintypes.MsgVoteInbound {
	ob.logger.Inbound.Debug().Msgf("Processing inbound: %s", inbound.TxHash)
	amount := big.NewFloat(inbound.Value)
	amount = amount.Mul(amount, big.NewFloat(1e8))
	amountInt, _ := amount.Int(nil)
	message := hex.EncodeToString(inbound.MemoBytes)

	// compliance check
	// if the inbound contains restricted addresses, return nil
	if ob.DoesInboundContainsRestrictedAddress(inbound) {
		return nil
	}

	return zetacore.GetInboundVoteMessage(
		inbound.FromAddress,
		ob.Chain().ChainId,
		inbound.FromAddress,
		inbound.FromAddress,
		ob.ZetacoreClient().Chain().ChainId,
		cosmosmath.NewUintFromBigInt(amountInt),
		message,
		inbound.TxHash,
		inbound.BlockNumber,
		0,
		coin.CoinType_Gas,
		"",
		ob.ZetacoreClient().GetKeys().GetOperatorAddress().String(),
		0,
	)
}

// DoesInboundContainsRestrictedAddress returns true if the inbound contains restricted addresses
// TODO(revamp): move all compliance related functions in a specific file
func (ob *Observer) DoesInboundContainsRestrictedAddress(inTx *BTCInboundEvent) bool {
	receiver := ""
	parsedAddress, _, err := chains.ParseAddressAndData(hex.EncodeToString(inTx.MemoBytes))
	if err == nil && parsedAddress != (ethcommon.Address{}) {
		receiver = parsedAddress.Hex()
	}
	if config.ContainRestrictedAddress(inTx.FromAddress, receiver) {
		compliance.PrintComplianceLog(ob.logger.Inbound, ob.logger.Compliance,
			false, ob.Chain().ChainId, inTx.TxHash, inTx.FromAddress, receiver, "BTC")
		return true
	}
	return false
}

// GetBtcEvent either returns a valid BTCInboundEvent or nil
// Note: the caller should retry the tx on error (e.g., GetSenderAddressByVin failed)
// TODO(revamp): simplify this function
func GetBtcEvent(
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
			memo, found, err = bitcoin.DecodeOpReturnMemo(vout1.ScriptPubKey.Hex, tx.Txid)
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

		fromAddress, err := GetSenderAddressByVin(rpcClient, tx.Vin[0], netParams)
		if err != nil {
			return nil, errors.Wrapf(err, "error getting sender address for inbound: %s", tx.Txid)
		}

		return &BTCInboundEvent{
			FromAddress: fromAddress,
			ToAddress:   tssAddress,
			Value:       value,
			MemoBytes:   memo,
			BlockNumber: blockNumber,
			TxHash:      tx.Txid,
		}, nil
	}
	return nil, nil
}
