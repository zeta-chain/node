package bitcoin

import (
	"errors"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/zetabridge"
)

func (ob *BitcoinChainClient) ExternalChainWatcherForNewInboundTrackerSuggestions() {
	ticker, err := NewDynamicTicker("Bitcoin_WatchInTx_InboundTrackerSuggestions", ob.GetChainParams().InTxTicker)
	if err != nil {
		ob.logger.WatchInTx.Err(err).Msg("error creating ticker")
		return
	}

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			err := ob.ObserveTrackerSuggestions()
			if err != nil {
				ob.logger.WatchInTx.Error().Err(err).Msg("error observing in tx")
			}
			ticker.UpdateInterval(ob.GetChainParams().InTxTicker, ob.logger.WatchInTx)
		case <-ob.stop:
			ob.logger.WatchInTx.Info().Msg("ExternalChainWatcher for BTC inboundTrackerSuggestions stopped")
			return
		}
	}
}

func (ob *BitcoinChainClient) ObserveTrackerSuggestions() error {
	trackers, err := ob.zetaClient.GetInboundTrackersForChain(ob.chain.ChainId)
	if err != nil {
		return err
	}
	for _, tracker := range trackers {
		ob.logger.WatchInTx.Info().Msgf("checking tracker with hash :%s and coin-type :%s ", tracker.TxHash, tracker.CoinType)
		ballotIdentifier, err := ob.CheckReceiptForBtcTxHash(tracker.TxHash, true)
		if err != nil {
			return err
		}
		ob.logger.WatchInTx.Info().Msgf("Vote submitted for inbound Tracker,Chain : %s,Ballot Identifier : %s, coin-type %s", ob.chain.ChainName, ballotIdentifier, common.CoinType_Gas.String())
	}
	return nil
}

func (ob *BitcoinChainClient) CheckReceiptForBtcTxHash(txHash string, vote bool) (string, error) {
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
	block, err := ob.rpcClient.GetBlockVerbose(blockHash)
	if err != nil {
		return "", err
	}
	tss, err := ob.zetaClient.GetBtcTssAddress()
	if err != nil {
		return "", err
	}
	// #nosec G701 always positive
	event, err := GetBtcEvent(*tx, tss, uint64(block.Height), &ob.logger.WatchInTx, ob.chain.ChainId)
	if err != nil {
		return "", err
	}
	if event == nil {
		return "", errors.New("no btc deposit event found")
	}
	msg := ob.GetInboundVoteMessageFromBtcEvent(event)
	if !vote {
		return msg.Digest(), nil
	}
	zetaHash, ballot, err := ob.zetaClient.PostVoteInbound(zetabridge.PostVoteInboundGasLimit, zetabridge.PostVoteInboundExecutionGasLimit, msg)
	if err != nil {
		ob.logger.WatchInTx.Error().Err(err).Msg("error posting to zeta core")
		return "", err
	} else if ballot == "" {
		ob.logger.WatchInTx.Info().Msgf("BTC deposit detected and reported: PostVoteInbound zeta tx: %s ballot %s", zetaHash, ballot)
	}
	return msg.Digest(), nil
}
