package evm

import (
	"errors"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/bitcoin"
	"github.com/zeta-chain/zetacore/zetaclient/zetabridge"
	"golang.org/x/net/context"
)

// ExternalChainWatcherForNewInboundTrackerSuggestions At each tick, gets a list of Inbound tracker suggestions from zeta-core and tries to check if the in-tx was confirmed.
// If it was, it tries to broadcast the confirmation vote. If this zeta client has previously broadcast the vote, the tx would be rejected
func (ob *EVMChainClient) ExternalChainWatcherForNewInboundTrackerSuggestions() {
	ticker, err := bitcoin.NewDynamicTicker(
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

	var msg *types.MsgVoteOnObservedInboundTx
	for _, log := range receipt.Logs {
		event, err := connector.ParseZetaSent(*log)
		if err == nil && event != nil {
			// sanity check tx event
			err = ob.CheckEvmTxLog(&event.Raw, addrConnector, txHash, TopicsZetaSent)
			if err == nil {
				msg = ob.GetInboundVoteMsgForZetaSentEvent(event)
				if msg != nil {
					break
				}
			} else {
				ob.logger.ExternalChainWatcher.Error().Err(err).Msg("CheckEvmTxLog error on ZetaSent event")
			}
		}
	}
	if msg == nil {
		return "", errors.New("no ZetaSent event found")
	}
	if !vote {
		return msg.Digest(), nil
	}

	zetaHash, ballot, err := ob.zetaClient.PostVoteInbound(zetabridge.PostVoteInboundGasLimit, zetabridge.PostVoteInboundMessagePassingExecutionGasLimit, msg)
	if err != nil {
		ob.logger.ExternalChainWatcher.Error().Err(err).Msg("error posting to zeta core")
		return "", err
	} else if zetaHash != "" {
		ob.logger.ExternalChainWatcher.Info().Msgf("ZetaSent event detected and reported: PostVoteInbound zeta tx: %s ballot %s", zetaHash, ballot)
	}

	return msg.Digest(), nil
}

func (ob *EVMChainClient) CheckReceiptForCoinTypeERC20(txHash string, vote bool) (string, error) {
	addrCustory, custody, err := ob.GetERC20CustodyContract()
	if err != nil {
		return "", err
	}
	// get transaction, receipt and sender
	hash := ethcommon.HexToHash(txHash)
	tx, _, err := ob.evmClient.TransactionByHash(context.Background(), hash)
	if err != nil {
		return "", err
	}
	receipt, err := ob.evmClient.TransactionReceipt(context.Background(), hash)
	if err != nil {
		return "", err
	}
	sender, err := ob.evmClient.TransactionSender(context.Background(), tx, receipt.BlockHash, receipt.TransactionIndex)
	if err != nil {
		return "", err
	}

	// check if the tx is confirmed
	lastHeight := ob.GetLastBlockHeight()
	if !ob.HasEnoughConfirmations(receipt, lastHeight) {
		return "", fmt.Errorf("txHash %s has not been confirmed yet: receipt block %d current block %d", txHash, receipt.BlockNumber, lastHeight)
	}

	var msg *types.MsgVoteOnObservedInboundTx
	for _, log := range receipt.Logs {
		zetaDeposited, err := custody.ParseDeposited(*log)
		if err == nil && zetaDeposited != nil {
			// sanity check tx event
			err = ob.CheckEvmTxLog(&zetaDeposited.Raw, addrCustory, txHash, TopicsDeposited)
			if err == nil {
				msg = ob.GetInboundVoteMsgForDepositedEvent(zetaDeposited, sender)
				if err == nil {
					break
				}
			} else {
				ob.logger.ExternalChainWatcher.Error().Err(err).Msg("CheckEvmTxLog error on ERC20CustodyDeposited event")
			}
		}
	}
	if msg == nil {
		return "", errors.New("no ERC20CustodyDeposited event found")
	}
	if !vote {
		return msg.Digest(), nil
	}

	zetaHash, ballot, err := ob.zetaClient.PostVoteInbound(zetabridge.PostVoteInboundGasLimit, zetabridge.PostVoteInboundExecutionGasLimit, msg)
	if err != nil {
		ob.logger.ExternalChainWatcher.Error().Err(err).Msg("error posting to zeta core")
		return "", err
	} else if zetaHash != "" {
		ob.logger.ExternalChainWatcher.Info().Msgf("ERC20 Deposit event detected and reported: PostVoteInbound zeta tx: %s ballot %s", zetaHash, ballot)
	}

	return msg.Digest(), nil
}

func (ob *EVMChainClient) CheckReceiptForCoinTypeGas(txHash string, vote bool) (string, error) {
	// TSS address should be set
	tssAddress := ob.Tss.EVMAddress()
	if tssAddress == (ethcommon.Address{}) {
		return "", errors.New("TSS address not set")
	}

	// check transaction and receipt
	hash := ethcommon.HexToHash(txHash)
	tx, isPending, err := ob.evmClient.TransactionByHash(context.Background(), hash)
	if err != nil {
		return "", err
	}
	if isPending {
		return "", errors.New("tx is still pending")
	}
	if tx.To() == nil {
		return "", errors.New("tx.To() is nil")
	}
	if *tx.To() != tssAddress {
		return "", fmt.Errorf("tx.To() %s is not TSS address", tssAddress.Hex())
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
	sender, err := ob.evmClient.TransactionSender(context.Background(), tx, receipt.BlockHash, receipt.TransactionIndex)
	if err != nil {
		return "", err
	}

	// check if the tx is confirmed
	lastHeight := ob.GetLastBlockHeight()
	if !ob.HasEnoughConfirmations(receipt, lastHeight) {
		return "", fmt.Errorf("txHash %s has not been confirmed yet: receipt block %d current block %d", txHash, receipt.BlockNumber, lastHeight)
	}
	msg := ob.GetInboundVoteMsgForTokenSentToTSS(tx, sender, receipt.BlockNumber.Uint64())
	if msg == nil {
		return "", errors.New("no message built for token sent to TSS")
	}
	if !vote {
		return msg.Digest(), nil
	}

	zetaHash, ballot, err := ob.zetaClient.PostVoteInbound(zetabridge.PostVoteInboundGasLimit, zetabridge.PostVoteInboundExecutionGasLimit, msg)
	if err != nil {
		ob.logger.ExternalChainWatcher.Error().Err(err).Msg("error posting to zeta core")
		return "", err
	} else if zetaHash != "" {
		ob.logger.ExternalChainWatcher.Info().Msgf("Gas deposit detected and reported: PostVoteInbound zeta tx: %s ballot %s", zetaHash, ballot)
	}

	return msg.Digest(), nil
}
