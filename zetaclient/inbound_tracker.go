package zetaclient

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"golang.org/x/net/context"
)

// ExternalChainWatcherForNewInboundTrackerSuggestions At each tick, gets a list of Inbound tracker suggestions from zeta-core and tries to check if the in-tx was confirmed.
// If it was, it tries to broadcast the confirmation vote. If this zeta client has previously broadcast the vote, the tx would be rejected
func (ob *EVMChainClient) ExternalChainWatcherForNewInboundTrackerSuggestions() {
	ticker, err := NewDynamicTicker(
		fmt.Sprintf("EVM_ExternalChainWatcher_InboundTrackerSuggestions_%d", ob.chain.ChainId),
		ob.GetChainParams().InTxTicker,
	)
	if err != nil {
		ob.logger.ExternalChainWatcher.Err(err).Msg("error creating ticker")
		return
	}

	defer ticker.Stop()
	ob.logger.ExternalChainWatcher.Info().Msg("ExternalChainWatcher for inboundTrackerSuggestions started")
	for {
		select {
		case <-ticker.C():
			err := ob.ObserveTrackerSuggestions()
			if err != nil {
				ob.logger.ExternalChainWatcher.Err(err).Msg("ObserveTrackerSuggestions error")
			}
			ticker.UpdateInterval(ob.GetChainParams().InTxTicker, ob.logger.ExternalChainWatcher)
		case <-ob.stop:
			ob.logger.ExternalChainWatcher.Info().Msg("ExternalChainWatcher for inboundTrackerSuggestions stopped")
			return
		}
	}
}

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
	zetaHash, ballot, err := ob.zetaClient.PostSend(PostSendEVMGasLimit, msg)
	if err != nil {
		ob.logger.WatchInTx.Error().Err(err).Msg("error posting to zeta core")
		return "", err
	} else if ballot == "" {
		ob.logger.WatchInTx.Info().Msgf("BTC deposit detected and reported: PostSend zeta tx: %s ballot %s", zetaHash, ballot)
	}
	return msg.Digest(), nil
}

func (ob *EVMChainClient) ObserveTrackerSuggestions() error {
	trackers, err := ob.zetaClient.GetInboundTrackersForChain(ob.chain.ChainId)
	if err != nil {
		return err
	}
	for _, tracker := range trackers {
		ob.logger.ExternalChainWatcher.Info().Msgf("checking tracker with hash :%s and coin-type :%s ", tracker.TxHash, tracker.CoinType)
		switch tracker.CoinType {
		case common.CoinType_Zeta:
			ballotIdentifier, err := ob.CheckReceiptForCoinTypeZeta(tracker.TxHash, true)
			if err != nil {
				return err
			}
			ob.logger.ExternalChainWatcher.Info().Msgf("Vote submitted for inbound Tracker,Chain : %s,Ballot Identifier : %s, coin-type %s", ob.chain.ChainName, ballotIdentifier, common.CoinType_Zeta.String())
		case common.CoinType_ERC20:
			ballotIdentifier, err := ob.CheckReceiptForCoinTypeERC20(tracker.TxHash, true)
			if err != nil {
				return err
			}
			ob.logger.ExternalChainWatcher.Info().Msgf("Vote submitted for inbound Tracker,Chain : %s,Ballot Identifier : %s, coin-type %s", ob.chain.ChainName, ballotIdentifier, common.CoinType_ERC20.String())
		case common.CoinType_Gas:
			ballotIdentifier, err := ob.CheckReceiptForCoinTypeGas(tracker.TxHash, true)
			if err != nil {
				return err
			}
			ob.logger.ExternalChainWatcher.Info().Msgf("Vote submitted for inbound Tracker,Chain : %s,Ballot Identifier : %s, coin-type %s", ob.chain.ChainName, ballotIdentifier, common.CoinType_Gas.String())
		}
	}
	return nil
}

func (ob *EVMChainClient) CheckReceiptForCoinTypeZeta(txHash string, vote bool) (string, error) {
	addrConnector, connector, err := ob.GetConnectorContract()
	if err != nil {
		return "", err
	}
	hash := ethcommon.HexToHash(txHash)
	receipt, err := ob.evmClient.TransactionReceipt(context.Background(), hash)
	if err != nil {
		return "", err
	}

	// check if the tx is confirmed
	lastHeight := ob.GetLastBlockHeight()
	if !ob.HasEnoughConfirmations(receipt, lastHeight) {
		return "", fmt.Errorf("txHash %s has not been confirmed yet: receipt block %d current block %d", txHash, receipt.BlockNumber, lastHeight)
	}

	var msg types.MsgVoteOnObservedInboundTx
	for _, log := range receipt.Logs {
		event, err := connector.ParseZetaSent(*log)
		if err == nil && event != nil {
			// sanity check tx event
			err = ob.CheckEvmTxLog(&event.Raw, addrConnector, txHash, TopicsZetaSent)
			if err == nil {
				msg, err = ob.GetInboundVoteMsgForZetaSentEvent(event)
				if err == nil {
					break
				}
			} else {
				ob.logger.ExternalChainWatcher.Error().Err(err).Msg("CheckEvmTxLog error on ZetaSent event")
			}
		}
	}
	if !vote {
		return msg.Digest(), nil
	}

	zetaHash, ballot, err := ob.zetaClient.PostSend(PostSendNonEVMGasLimit, &msg)
	if err != nil {
		ob.logger.ExternalChainWatcher.Error().Err(err).Msg("error posting to zeta core")
		return "", err
	} else if zetaHash != "" {
		ob.logger.ExternalChainWatcher.Info().Msgf("ZetaSent event detected and reported: PostSend zeta tx: %s ballot %s", zetaHash, ballot)
	}

	return msg.Digest(), nil
}

func (ob *EVMChainClient) CheckReceiptForCoinTypeERC20(txHash string, vote bool) (string, error) {
	addrCustory, custody, err := ob.GetERC20CustodyContract()
	if err != nil {
		return "", err
	}
	hash := ethcommon.HexToHash(txHash)
	receipt, err := ob.evmClient.TransactionReceipt(context.Background(), hash)
	if err != nil {
		return "", err
	}

	// check if the tx is confirmed
	lastHeight := ob.GetLastBlockHeight()
	if !ob.HasEnoughConfirmations(receipt, lastHeight) {
		return "", fmt.Errorf("txHash %s has not been confirmed yet: receipt block %d current block %d", txHash, receipt.BlockNumber, lastHeight)
	}

	var msg types.MsgVoteOnObservedInboundTx
	for _, log := range receipt.Logs {
		zetaDeposited, err := custody.ParseDeposited(*log)
		if err == nil && zetaDeposited != nil {
			// sanity check tx event
			err = ob.CheckEvmTxLog(&zetaDeposited.Raw, addrCustory, txHash, TopicsDeposited)
			if err == nil {
				msg, err = ob.GetInboundVoteMsgForDepositedEvent(zetaDeposited)
				if err == nil {
					break
				}
			} else {
				ob.logger.ExternalChainWatcher.Error().Err(err).Msg("CheckEvmTxLog error on Deposited event")
			}
		}
	}
	if !vote {
		return msg.Digest(), nil
	}

	zetaHash, ballot, err := ob.zetaClient.PostSend(PostSendEVMGasLimit, &msg)
	if err != nil {
		ob.logger.ExternalChainWatcher.Error().Err(err).Msg("error posting to zeta core")
		return "", err
	} else if zetaHash != "" {
		ob.logger.ExternalChainWatcher.Info().Msgf("ERC20 Deposit event detected and reported: PostSend zeta tx: %s ballot %s", zetaHash, ballot)
	}

	return msg.Digest(), nil
}

func (ob *EVMChainClient) CheckReceiptForCoinTypeGas(txHash string, vote bool) (string, error) {
	hash := ethcommon.HexToHash(txHash)
	tx, isPending, err := ob.evmClient.TransactionByHash(context.Background(), hash)
	if err != nil {
		return "", err
	}
	if isPending {
		return "", errors.New("tx is still pending")
	}

	receipt, err := ob.evmClient.TransactionReceipt(context.Background(), hash)
	if err != nil {
		ob.logger.ExternalChainWatcher.Err(err).Msg("TransactionReceipt error")
		return "", err
	}
	if receipt.Status != 1 { // 1: successful, 0: failed
		ob.logger.ExternalChainWatcher.Info().Msgf("tx %s failed; don't act", tx.Hash().Hex())
		return "", errors.New("tx not successful yet")
	}

	// check if the tx is confirmed
	lastHeight := ob.GetLastBlockHeight()
	if !ob.HasEnoughConfirmations(receipt, lastHeight) {
		return "", fmt.Errorf("txHash %s has not been confirmed yet: receipt block %d current block %d", txHash, receipt.BlockNumber, lastHeight)
	}

	block, err := ob.evmClient.BlockByNumber(context.Background(), receipt.BlockNumber)
	if err != nil {
		ob.logger.ExternalChainWatcher.Err(err).Msg("BlockByNumber error")
		return "", err
	}
	from, err := ob.evmClient.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
	if err != nil {
		ob.logger.ExternalChainWatcher.Err(err).Msg("TransactionSender error; trying local recovery (assuming LondonSigner dynamic fee tx type) of sender address")
		signer := ethtypes.NewLondonSigner(big.NewInt(ob.chain.ChainId))
		from, err = signer.Sender(tx)
		if err != nil {
			ob.logger.ExternalChainWatcher.Err(err).Msg("local recovery of sender address failed")
			return "", err
		}
	}
	msg := ob.GetInboundVoteMsgForTokenSentToTSS(tx.Hash(), tx.Value(), receipt, from, tx.Data())
	if !vote {
		return msg.Digest(), nil
	}

	zetaHash, ballot, err := ob.zetaClient.PostSend(PostSendEVMGasLimit, msg)
	if err != nil {
		ob.logger.ExternalChainWatcher.Error().Err(err).Msg("error posting to zeta core")
		return "", err
	} else if zetaHash != "" {
		ob.logger.ExternalChainWatcher.Info().Msgf("Gas deposit detected and reported: PostSend zeta tx: %s ballot %s", zetaHash, ballot)
	}

	return msg.Digest(), nil
}
