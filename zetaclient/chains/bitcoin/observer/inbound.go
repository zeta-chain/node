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

// WatchInTx watches Bitcoin chain for incoming txs and post votes to zetacore
func (ob *Observer) WatchInTx() {
	ticker, err := types.NewDynamicTicker("Bitcoin_WatchInTx", ob.GetChainParams().InTxTicker)
	if err != nil {
		ob.logger.InTx.Error().Err(err).Msg("error creating ticker")
		return
	}
	defer ticker.Stop()

	ob.logger.InTx.Info().Msgf("WatchInTx started for chain %d", ob.chain.ChainId)
	sampledLogger := ob.logger.InTx.Sample(&zerolog.BasicSampler{N: 10})

	for {
		select {
		case <-ticker.C():
			if !context.IsInboundObservationEnabled(ob.coreContext, ob.GetChainParams()) {
				sampledLogger.Info().Msgf("WatchInTx: inbound observation is disabled for chain %d", ob.chain.ChainId)
				continue
			}
			err := ob.ObserveInTx()
			if err != nil {
				ob.logger.InTx.Error().Err(err).Msg("WatchInTx error observing in tx")
			}
			ticker.UpdateInterval(ob.GetChainParams().InTxTicker, ob.logger.InTx)
		case <-ob.stop:
			ob.logger.InTx.Info().Msgf("WatchInTx stopped for chain %d", ob.chain.ChainId)
			return
		}
	}
}

func (ob *Observer) ObserveInTx() error {
	// get and update latest block height
	cnt, err := ob.rpcClient.GetBlockCount()
	if err != nil {
		return fmt.Errorf("observeInTxBTC: error getting block number: %s", err)
	}
	if cnt < 0 {
		return fmt.Errorf("observeInTxBTC: block number is negative: %d", cnt)
	}
	if cnt < ob.GetLastBlockHeight() {
		return fmt.Errorf("observeInTxBTC: block number should not decrease: current %d last %d", cnt, ob.GetLastBlockHeight())
	}
	ob.SetLastBlockHeight(cnt)

	// skip if current height is too low
	// #nosec G701 always in range
	confirmedBlockNum := cnt - int64(ob.GetChainParams().ConfirmationCount)
	if confirmedBlockNum < 0 {
		return fmt.Errorf("observeInTxBTC: skipping observer, current block number %d is too low", cnt)
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
		ob.logger.InTx.Error().Err(err).Msgf("observeInTxBTC: error getting bitcoin block %d", blockNumber)
		return err
	}
	ob.logger.InTx.Info().Msgf("observeInTxBTC: block %d has %d txs, current block %d, last block %d",
		blockNumber, len(res.Block.Tx), cnt, lastScanned)

	// add block header to zetacore
	// TODO: consider having a separate ticker(from TSS scaning) for posting block headers
	// https://github.com/zeta-chain/node/issues/1847
	flags := ob.coreContext.GetCrossChainFlags()
	if flags.BlockHeaderVerificationFlags != nil && flags.BlockHeaderVerificationFlags.IsBtcTypeChainEnabled {
		err = ob.postBlockHeader(blockNumber)
		if err != nil {
			ob.logger.InTx.Warn().Err(err).Msgf("observeInTxBTC: error posting block header %d", blockNumber)
		}
	}

	if len(res.Block.Tx) > 1 {
		// get depositor fee
		depositorFee := bitcoin.CalcDepositorFee(res.Block, ob.chain.ChainId, ob.netParams, ob.logger.InTx)

		// filter incoming txs to TSS address
		tssAddress := ob.Tss.BTCAddress()

		// add block header to zetacore
		// TODO: consider having a separate ticker(from TSS scaning) for posting block headers
		// https://github.com/zeta-chain/node/issues/1847
		blockHeaderVerification, found := ob.coreContext.GetBlockHeaderEnabledChains(ob.chain.ChainId)
		if found && blockHeaderVerification.Enabled {
			err = ob.postBlockHeader(blockNumber)
			if err != nil {
				ob.logger.InTx.Warn().Err(err).Msgf("observeInTxBTC: error posting block header %d", blockNumber)
			}
		}

		if len(res.Block.Tx) > 1 {
			// get depositor fee
			depositorFee := bitcoin.CalcDepositorFee(res.Block, ob.chain.ChainId, ob.netParams, ob.logger.InTx)

			// filter incoming txs to TSS address
			tssAddress := ob.Tss.BTCAddress()
			// #nosec G701 always positive
			inTxs, err := FilterAndParseIncomingTx(
				ob.rpcClient,
				res.Block.Tx,
				uint64(res.Block.Height),
				tssAddress,
				ob.logger.InTx,
				ob.netParams,
				depositorFee,
			)
			if err != nil {
				ob.logger.InTx.Error().Err(err).Msgf("observeInTxBTC: error filtering incoming txs for block %d", blockNumber)
				return err // we have to re-scan this block next time
			}

			// post inbound vote message to zetacore
			for _, inTx := range inTxs {
				msg := ob.GetInboundVoteMessageFromBtcEvent(inTx)
				if msg != nil {
					zetaHash, ballot, err := ob.zetacoreClient.PostVoteInbound(zetacore.PostVoteInboundGasLimit, zetacore.PostVoteInboundExecutionGasLimit, msg)
					if err != nil {
						ob.logger.InTx.Error().Err(err).Msgf("observeInTxBTC: error posting to zetacore for tx %s", inTx.TxHash)
						return err // we have to re-scan this block next time
					} else if zetaHash != "" {
						ob.logger.InTx.Info().Msgf("observeInTxBTC: PostVoteInbound zeta tx hash: %s inTx %s ballot %s fee %v",
							zetaHash, inTx.TxHash, ballot, depositorFee)
					}
				}
			}
		}

		// Save LastBlockHeight
		ob.SetLastBlockHeightScanned(blockNumber)

		// #nosec G701 always positive
		inTxs, err := FilterAndParseIncomingTx(
			ob.rpcClient,
			res.Block.Tx,
			uint64(res.Block.Height),
			tssAddress,
			ob.logger.InTx,
			ob.netParams,
			depositorFee,
		)
		if err != nil {
			ob.logger.InTx.Error().Err(err).Msgf("observeInTxBTC: error filtering incoming txs for block %d", blockNumber)
			return err // we have to re-scan this block next time
		}

		// post inbound vote message to zetacore
		for _, inTx := range inTxs {
			msg := ob.GetInboundVoteMessageFromBtcEvent(inTx)
			if msg != nil {
				zetaHash, ballot, err := ob.zetacoreClient.PostVoteInbound(zetacore.PostVoteInboundGasLimit, zetacore.PostVoteInboundExecutionGasLimit, msg)
				if err != nil {
					ob.logger.InTx.Error().Err(err).Msgf("observeInTxBTC: error posting to zetacore for tx %s", inTx.TxHash)
					return err // we have to re-scan this block next time
				} else if zetaHash != "" {
					ob.logger.InTx.Info().Msgf("observeInTxBTC: PostVoteInbound zeta tx hash: %s inTx %s ballot %s fee %v",
						zetaHash, inTx.TxHash, ballot, depositorFee)
				}
			}
		}
	}

	// Save LastBlockHeight
	ob.SetLastBlockHeightScanned(blockNumber)

	// #nosec G701 always positive
	if err := ob.db.Save(types.ToLastBlockSQLType(uint64(blockNumber))).Error; err != nil {
		ob.logger.InTx.Error().Err(err).Msgf("observeInTxBTC: error writing last scanned block %d to db", blockNumber)
	}

	return nil
}

// WatchIntxTracker watches zetacore for bitcoin intx trackers
func (ob *Observer) WatchIntxTracker() {
	ticker, err := types.NewDynamicTicker("Bitcoin_WatchIntxTracker", ob.GetChainParams().InTxTicker)
	if err != nil {
		ob.logger.InTx.Err(err).Msg("error creating ticker")
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
				ob.logger.InTx.Error().Err(err).Msgf("error observing intx tracker for chain %d", ob.chain.ChainId)
			}
			ticker.UpdateInterval(ob.GetChainParams().InTxTicker, ob.logger.InTx)
		case <-ob.stop:
			ob.logger.InTx.Info().Msgf("WatchIntxTracker stopped for chain %d", ob.chain.ChainId)
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
		ob.logger.InTx.Info().Msgf("checking tracker with hash :%s and coin-type :%s ", tracker.TxHash, tracker.CoinType)
		ballotIdentifier, err := ob.CheckReceiptForBtcTxHash(tracker.TxHash, true)
		if err != nil {
			return err
		}
		ob.logger.InTx.Info().Msgf("Vote submitted for inbound Tracker, Chain : %s,Ballot Identifier : %s, coin-type %s", ob.chain.ChainName, ballotIdentifier, coin.CoinType_Gas.String())
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

	depositorFee := bitcoin.CalcDepositorFee(blockVb, ob.chain.ChainId, ob.netParams, ob.logger.InTx)
	tss, err := ob.zetacoreClient.GetBtcTssAddress(ob.chain.ChainId)
	if err != nil {
		return "", err
	}

	// #nosec G701 always positive
	event, err := GetBtcEvent(ob.rpcClient, *tx, tss, uint64(blockVb.Height), ob.logger.InTx, ob.netParams, depositorFee)
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
		ob.logger.InTx.Error().Err(err).Msg("error posting to zetacore")
		return "", err
	} else if zetaHash != "" {
		ob.logger.InTx.Info().Msgf("BTC deposit detected and reported: PostVoteInbound zeta tx hash: %s inTx %s ballot %s fee %v",
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
) ([]*BTCInTxEvent, error) {
	inTxs := make([]*BTCInTxEvent, 0)
	for idx, tx := range txs {
		if idx == 0 {
			continue // the first tx is coinbase; we do not process coinbase tx
		}

		inTx, err := GetBtcEvent(rpcClient, tx, tssAddress, blockNumber, logger, netParams, depositorFee)
		if err != nil {
			// unable to parse the tx, the caller should retry
			return nil, errors.Wrapf(err, "error getting btc event for tx %s in block %d", tx.Txid, blockNumber)
		}

		if inTx != nil {
			inTxs = append(inTxs, inTx)
			logger.Info().Msgf("FilterAndParseIncomingTx: found btc event for tx %s in block %d", tx.Txid, blockNumber)
		}
	}
	return inTxs, nil
}

func (ob *Observer) GetInboundVoteMessageFromBtcEvent(inTx *BTCInTxEvent) *crosschaintypes.MsgVoteOnObservedInboundTx {
	ob.logger.InTx.Debug().Msgf("Processing inTx: %s", inTx.TxHash)
	amount := big.NewFloat(inTx.Value)
	amount = amount.Mul(amount, big.NewFloat(1e8))
	amountInt, _ := amount.Int(nil)
	message := hex.EncodeToString(inTx.MemoBytes)

	// compliance check
	// if the inbound contains restricted addresses, return nil
	if ob.IsInTxRestricted(inTx) {
		return nil
	}

	return zetacore.GetInBoundVoteMessage(
		inTx.FromAddress,
		ob.chain.ChainId,
		inTx.FromAddress,
		inTx.FromAddress,
		ob.zetacoreClient.Chain().ChainId,
		cosmosmath.NewUintFromBigInt(amountInt),
		message,
		inTx.TxHash,
		inTx.BlockNumber,
		0,
		coin.CoinType_Gas,
		"",
		ob.zetacoreClient.GetKeys().GetOperatorAddress().String(),
		0,
	)
}

// IsInTxRestricted returns true if the inTx contains restricted addresses
func (ob *Observer) IsInTxRestricted(inTx *BTCInTxEvent) bool {
	receiver := ""
	parsedAddress, _, err := chains.ParseAddressAndData(hex.EncodeToString(inTx.MemoBytes))
	if err == nil && parsedAddress != (ethcommon.Address{}) {
		receiver = parsedAddress.Hex()
	}
	if config.ContainRestrictedAddress(inTx.FromAddress, receiver) {
		compliance.PrintComplianceLog(ob.logger.InTx, ob.logger.Compliance,
			false, ob.chain.ChainId, inTx.TxHash, inTx.FromAddress, receiver, "BTC")
		return true
	}
	return false
}

// GetBtcEvent either returns a valid BTCInTxEvent or nil
// Note: the caller should retry the tx on error (e.g., GetSenderAddressByVin failed)
func GetBtcEvent(
	rpcClient interfaces.BTCRPCClient,
	tx btcjson.TxRawResult,
	tssAddress string,
	blockNumber uint64,
	logger zerolog.Logger,
	netParams *chaincfg.Params,
	depositorFee float64,
) (*BTCInTxEvent, error) {
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
			return nil, fmt.Errorf("GetBtcEvent: no input found for intx: %s", tx.Txid)
		}

		fromAddress, err := GetSenderAddressByVin(rpcClient, tx.Vin[0], netParams)
		if err != nil {
			return nil, errors.Wrapf(err, "error getting sender address for intx: %s", tx.Txid)
		}

		return &BTCInTxEvent{
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
