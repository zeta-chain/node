package bitcoin

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/zetaclient/types"
	"github.com/zeta-chain/zetacore/zetaclient/zetabridge"
)

// WatchIntxTracker watches zetacore for bitcoin intx trackers
func (ob *BTCChainClient) WatchIntxTracker() {
	ticker, err := types.NewDynamicTicker("Bitcoin_WatchIntxTracker", ob.GetChainParams().InTxTicker)
	if err != nil {
		ob.logger.InTx.Err(err).Msg("error creating ticker")
		return
	}

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			if !ob.GetChainParams().IsSupported {
				continue
			}
			err := ob.ObserveTrackerSuggestions()
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

func (ob *BTCChainClient) ObserveTrackerSuggestions() error {
	trackers, err := ob.zetaClient.GetInboundTrackersForChain(ob.chain.ChainId)
	if err != nil {
		return err
	}
	for _, tracker := range trackers {
		ob.logger.InTx.Info().Msgf("checking tracker with hash :%s and coin-type :%s ", tracker.TxHash, tracker.CoinType)
		ballotIdentifier, err := ob.CheckReceiptForBtcTxHash(tracker.TxHash, true)
		if err != nil {
			return err
		}
		ob.logger.InTx.Info().Msgf("Vote submitted for inbound Tracker,Chain : %s,Ballot Identifier : %s, coin-type %s", ob.chain.ChainName, ballotIdentifier, coin.CoinType_Gas.String())
	}
	return nil
}

func (ob *BTCChainClient) CheckReceiptForBtcTxHash(txHash string, vote bool) (string, error) {
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
	depositorFee := CalcDepositorFee(blockVb, ob.chain.ChainId, ob.netParams, ob.logger.InTx)
	tss, err := ob.zetaClient.GetBtcTssAddress(ob.chain.ChainId)
	if err != nil {
		return "", err
	}
	// #nosec G701 always positive
	event, err := GetBtcEvent(*tx, tss, uint64(blockVb.Height), &ob.logger.InTx, ob.netParams, depositorFee)
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
	zetaHash, ballot, err := ob.zetaClient.PostVoteInbound(zetabridge.PostVoteInboundGasLimit, zetabridge.PostVoteInboundExecutionGasLimit, msg)
	if err != nil {
		ob.logger.InTx.Error().Err(err).Msg("error posting to zeta core")
		return "", err
	} else if zetaHash != "" {
		ob.logger.InTx.Info().Msgf("BTC deposit detected and reported: PostVoteInbound zeta tx hash: %s inTx %s ballot %s fee %v",
			zetaHash, txHash, ballot, depositorFee)
	}
	return msg.Digest(), nil
}
