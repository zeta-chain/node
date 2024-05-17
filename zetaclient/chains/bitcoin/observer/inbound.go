package observer

import (
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
	"github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/types"
	"github.com/zeta-chain/zetacore/zetaclient/zetacore"
)

// WatchInbound watches Bitcoin chain for incoming txs and post votes to zetacore
func (ob *Observer) WatchInbound() {
	ticker, err := types.NewDynamicTicker("Bitcoin_WatchInbound", ob.GetChainParams().InboundTicker)
	if err != nil {
		ob.logger.Inbound.Error().Err(err).Msg("error creating ticker")
		return
	}
	defer ticker.Stop()

	ob.logger.Inbound.Info().Msgf("WatchInbound started for chain %d", ob.chain.ChainId)
	sampledLogger := ob.logger.Inbound.Sample(&zerolog.BasicSampler{N: 10})

	for {
		select {
		case <-ticker.C():
			if !context.IsInboundObservationEnabled(ob.coreContext, ob.GetChainParams()) {
				sampledLogger.Info().Msgf("WatchInbound: inbound observation is disabled for chain %d", ob.chain.ChainId)
				continue
			}
			err := ob.ObserveInbound()
			if err != nil {
				ob.logger.Inbound.Error().Err(err).Msg("WatchInbound error observing in tx")
			}
			ticker.UpdateInterval(ob.GetChainParams().InboundTicker, ob.logger.Inbound)
		case <-ob.stop:
			ob.logger.Inbound.Info().Msgf("WatchInbound stopped for chain %d", ob.chain.ChainId)
			return
		}
	}
}

func (ob *Observer) ObserveInbound() error {
	// get and update latest block height
	cnt, err := ob.rpcClient.GetBlockCount()
	if err != nil {
		return fmt.Errorf("observeInboundBTC: error getting block number: %s", err)
	}
	if cnt < 0 {
		return fmt.Errorf("observeInboundBTC: block number is negative: %d", cnt)
	}
	if cnt < ob.GetLastBlockHeight() {
		return fmt.Errorf("observeInboundBTC: block number should not decrease: current %d last %d", cnt, ob.GetLastBlockHeight())
	}
	ob.SetLastBlockHeight(cnt)

	// skip if current height is too low
	// #nosec G701 always in range
	confirmedBlockNum := cnt - int64(ob.GetChainParams().ConfirmationCount)
	if confirmedBlockNum < 0 {
		return fmt.Errorf("observeInboundBTC: skipping observer, current block number %d is too low", cnt)
	}

	// skip if no new block is confirmed
	lastScanned := ob.GetLastBlockHeightScanned()
	if lastScanned >= confirmedBlockNum {
		return nil
	}

	// query incoming gas asset to TSS address
	blockNumber := lastScanned + 1
	res, err := ob.GetBlockByNumberCached(blockNumber)
	if err != nil {
		ob.logger.Inbound.Error().Err(err).Msgf("observeInboundBTC: error getting bitcoin block %d", blockNumber)
		return err
	}
	ob.logger.Inbound.Info().Msgf("observeInboundBTC: block %d has %d txs, current block %d, last block %d",
		blockNumber, len(res.Block.Tx), cnt, lastScanned)

	// add block header to zetacore
	// TODO: consider having a separate ticker(from TSS scaning) for posting block headers
	// https://github.com/zeta-chain/node/issues/1847
	flags := ob.coreContext.GetCrossChainFlags()
	if flags.BlockHeaderVerificationFlags != nil && flags.BlockHeaderVerificationFlags.IsBtcTypeChainEnabled {
		err = ob.postBlockHeader(blockNumber)
		if err != nil {
			ob.logger.Inbound.Warn().Err(err).Msgf("observeInboundBTC: error posting block header %d", blockNumber)
		}
	}

	if len(res.Block.Tx) > 1 {
		// get depositor fee
		depositorFee := bitcoin.CalcDepositorFee(res.Block, ob.chain.ChainId, ob.netParams, ob.logger.Inbound)

		// filter incoming txs to TSS address
		tssAddress := ob.Tss.BTCAddress()

		// add block header to zetacore
		// TODO: consider having a separate ticker(from TSS scaning) for posting block headers
		// https://github.com/zeta-chain/node/issues/1847
		blockHeaderVerification, found := ob.coreContext.GetBlockHeaderEnabledChains(ob.chain.ChainId)
		if found && blockHeaderVerification.Enabled {
			err = ob.postBlockHeader(blockNumber)
			if err != nil {
				ob.logger.Inbound.Warn().Err(err).Msgf("observeInboundBTC: error posting block header %d", blockNumber)
			}
		}

		if len(res.Block.Tx) > 1 {
			// get depositor fee
			depositorFee := bitcoin.CalcDepositorFee(res.Block, ob.chain.ChainId, ob.netParams, ob.logger.Inbound)

			// filter incoming txs to TSS address
			tssAddress := ob.Tss.BTCAddress()
			// #nosec G701 always positive
			inbounds, err := FilterAndParseIncomingTx(
				ob.rpcClient,
				res.Block.Tx,
				uint64(res.Block.Height),
				tssAddress,
				ob.logger.Inbound,
				ob.netParams,
				depositorFee,
			)
			if err != nil {
				ob.logger.Inbound.Error().Err(err).Msgf("observeInboundBTC: error filtering incoming txs for block %d", blockNumber)
				return err // we have to re-scan this block next time
			}

			// post inbound vote message to zetacore
			for _, inbound := range inbounds {
				msg := ob.GetInboundVoteMessageFromBtcEvent(inbound)
				if msg != nil {
					zetaHash, ballot, err := ob.zetacoreClient.PostVoteInbound(zetacore.PostVoteInboundGasLimit, zetacore.PostVoteInboundExecutionGasLimit, msg)
					if err != nil {
						ob.logger.Inbound.Error().Err(err).Msgf("observeInboundBTC: error posting to zetacore for tx %s", inbound.TxHash)
						return err // we have to re-scan this block next time
					} else if zetaHash != "" {
						ob.logger.Inbound.Info().Msgf("observeInboundBTC: PostVoteInbound zeta tx hash: %s inbound %s ballot %s fee %v",
							zetaHash, inbound.TxHash, ballot, depositorFee)
					}
				}
			}
		}

		// Save LastBlockHeight
		ob.SetLastBlockHeightScanned(blockNumber)

		// #nosec G701 always positive
		inbounds, err := FilterAndParseIncomingTx(
			ob.rpcClient,
			res.Block.Tx,
			uint64(res.Block.Height),
			tssAddress,
			ob.logger.Inbound,
			ob.netParams,
			depositorFee,
		)
		if err != nil {
			ob.logger.Inbound.Error().Err(err).Msgf("observeInboundBTC: error filtering incoming txs for block %d", blockNumber)
			return err // we have to re-scan this block next time
		}

		// post inbound vote message to zetacore
		for _, inbound := range inbounds {
			msg := ob.GetInboundVoteMessageFromBtcEvent(inbound)
			if msg != nil {
				zetaHash, ballot, err := ob.zetacoreClient.PostVoteInbound(zetacore.PostVoteInboundGasLimit, zetacore.PostVoteInboundExecutionGasLimit, msg)
				if err != nil {
					ob.logger.Inbound.Error().Err(err).Msgf("observeInboundBTC: error posting to zetacore for tx %s", inbound.TxHash)
					return err // we have to re-scan this block next time
				} else if zetaHash != "" {
					ob.logger.Inbound.Info().Msgf("observeInboundBTC: PostVoteInbound zeta tx hash: %s inbound %s ballot %s fee %v",
						zetaHash, inbound.TxHash, ballot, depositorFee)
				}
			}
		}
	}

	// Save LastBlockHeight
	ob.SetLastBlockHeightScanned(blockNumber)

	// #nosec G701 always positive
	if err := ob.db.Save(types.ToLastBlockSQLType(uint64(blockNumber))).Error; err != nil {
		ob.logger.Inbound.Error().Err(err).Msgf("observeInboundBTC: error writing last scanned block %d to db", blockNumber)
	}

	return nil
}

// WatchInboundTracker watches zetacore for bitcoin inbound trackers
func (ob *Observer) WatchInboundTracker() {
	ticker, err := types.NewDynamicTicker("Bitcoin_WatchInboundTracker", ob.GetChainParams().InboundTicker)
	if err != nil {
		ob.logger.Inbound.Err(err).Msg("error creating ticker")
		return
	}

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			if !context.IsInboundObservationEnabled(ob.coreContext, ob.GetChainParams()) {
				continue
			}
			err := ob.ProcessInboundTrackers()
			if err != nil {
				ob.logger.Inbound.Error().Err(err).Msgf("error observing inbound tracker for chain %d", ob.chain.ChainId)
			}
			ticker.UpdateInterval(ob.GetChainParams().InboundTicker, ob.logger.Inbound)
		case <-ob.stop:
			ob.logger.Inbound.Info().Msgf("WatchInboundTracker stopped for chain %d", ob.chain.ChainId)
			return
		}
	}
}

// ProcessInboundTrackers processes inbound trackers
func (ob *Observer) ProcessInboundTrackers() error {
	trackers, err := ob.zetacoreClient.GetInboundTrackersForChain(ob.chain.ChainId)
	if err != nil {
		return err
	}

	for _, tracker := range trackers {
		ob.logger.Inbound.Info().Msgf("checking tracker with hash :%s and coin-type :%s ", tracker.TxHash, tracker.CoinType)
		ballotIdentifier, err := ob.CheckReceiptForBtcTxHash(tracker.TxHash, true)
		if err != nil {
			return err
		}
		ob.logger.Inbound.Info().Msgf("Vote submitted for inbound Tracker, Chain : %s,Ballot Identifier : %s, coin-type %s", ob.chain.ChainName, ballotIdentifier, coin.CoinType_Gas.String())
	}

	return nil
}

// CheckReceiptForBtcTxHash checks the receipt for a btc tx hash
func (ob *Observer) CheckReceiptForBtcTxHash(txHash string, vote bool) (string, error) {
	hash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return "", err
	}

	tx, err := ob.rpcClient.GetRawTransactionVerbose(hash)
	if err != nil {
		return "", err
	}

	blockHash, err := chainhash.NewHashFromStr(tx.BlockHash)
	if err != nil {
		return "", err
	}

	blockVb, err := ob.rpcClient.GetBlockVerboseTx(blockHash)
	if err != nil {
		return "", err
	}

	if len(blockVb.Tx) <= 1 {
		return "", fmt.Errorf("block %d has no transactions", blockVb.Height)
	}

	depositorFee := bitcoin.CalcDepositorFee(blockVb, ob.chain.ChainId, ob.netParams, ob.logger.Inbound)
	tss, err := ob.zetacoreClient.GetBtcTssAddress(ob.chain.ChainId)
	if err != nil {
		return "", err
	}

	// #nosec G701 always positive
	event, err := GetBtcEvent(ob.rpcClient, *tx, tss, uint64(blockVb.Height), ob.logger.Inbound, ob.netParams, depositorFee)
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

	zetaHash, ballot, err := ob.zetacoreClient.PostVoteInbound(zetacore.PostVoteInboundGasLimit, zetacore.PostVoteInboundExecutionGasLimit, msg)
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

func (ob *Observer) GetInboundVoteMessageFromBtcEvent(inbound *BTCInboundEvent) *crosschaintypes.MsgVoteInbound {
	ob.logger.Inbound.Debug().Msgf("Processing inbound: %s", inbound.TxHash)
	amount := big.NewFloat(inbound.Value)
	amount = amount.Mul(amount, big.NewFloat(1e8))
	amountInt, _ := amount.Int(nil)
	message := hex.EncodeToString(inbound.MemoBytes)

	// compliance check
	// if the inbound contains restricted addresses, return nil
	if ob.IsInboundRestricted(inbound) {
		return nil
	}

	return zetacore.GetInBoundVoteMessage(
		inbound.FromAddress,
		ob.chain.ChainId,
		inbound.FromAddress,
		inbound.FromAddress,
		ob.zetacoreClient.Chain().ChainId,
		cosmosmath.NewUintFromBigInt(amountInt),
		message,
		inbound.TxHash,
		inbound.BlockNumber,
		0,
		coin.CoinType_Gas,
		"",
		ob.zetacoreClient.GetKeys().GetOperatorAddress().String(),
		0,
	)
}

// IsInboundRestricted returns true if the inTx contains restricted addresses
func (ob *Observer) IsInboundRestricted(inTx *BTCInboundEvent) bool {
	receiver := ""
	parsedAddress, _, err := chains.ParseAddressAndData(hex.EncodeToString(inTx.MemoBytes))
	if err == nil && parsedAddress != (ethcommon.Address{}) {
		receiver = parsedAddress.Hex()
	}
	if config.ContainRestrictedAddress(inTx.FromAddress, receiver) {
		compliance.PrintComplianceLog(ob.logger.Inbound, ob.logger.Compliance,
			false, ob.chain.ChainId, inTx.TxHash, inTx.FromAddress, receiver, "BTC")
		return true
	}
	return false
}

// GetBtcEvent either returns a valid BTCInboundEvent or nil
// Note: the caller should retry the tx on error (e.g., GetSenderAddressByVin failed)
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
				logger.Info().Msgf("GetBtcEvent: btc deposit amount %v in txid %s is less than depositor fee %v", vout0.Value, tx.Txid, depositorFee)
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
