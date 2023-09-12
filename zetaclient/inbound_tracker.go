package zetaclient

import (
	"errors"
	"fmt"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"golang.org/x/net/context"
)

func (ob *EVMChainClient) ExternalChainWatcherForNewInboundTrackerSuggestions() {
	// At each tick, query the Connector contract
	ticker := NewDynamicTicker(fmt.Sprintf("EVM_ExternalChainWatcher_InboundTrackerSuggestions_%d", ob.chain.ChainId), ob.GetCoreParams().InTxTicker)
	defer ticker.Stop()
	ob.logger.ExternalChainWatcher.Info().Msg("ExternalChainWatcher for inboundTrackerSuggestions started")
	for {
		select {
		case <-ticker.C():
			err := ob.ObserveTrackerSuggestions()
			if err != nil {
				ob.logger.ExternalChainWatcher.Err(err).Msg("ObserveTrackerSuggestions error")
			}
			ticker.UpdateInterval(ob.GetCoreParams().InTxTicker, ob.logger.ExternalChainWatcher)
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
		switch tracker.CoinType {
		case common.CoinType_Zeta:
			err := ob.CheckReceiptForCoinTypeZeta(tracker.TxHash)
			if err != nil {
				return err
			}
		case common.CoinType_ERC20:
			err := ob.CheckReceiptForCoinTypeERC20(tracker.TxHash)
			if err != nil {
				return err
			}
		case common.CoinType_Gas:
			err := ob.CheckReceiptForCoinTypeGas(tracker.TxHash)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (ob *EVMChainClient) CheckReceiptForCoinTypeZeta(txHash string) error {
	connector, err := ob.GetConnectorContract()
	if err != nil {
		return err
	}
	hash := ethcommon.HexToHash(txHash)
	receipt, err := ob.EvmClient.TransactionReceipt(context.Background(), hash)
	if err != nil {
		return err
	}
	var msg *types.MsgVoteOnObservedInboundTx
	for _, log := range receipt.Logs {
		event, err := connector.ParseZetaSent(*log)
		if err != nil {
			msg = ob.GetInboundVoteForZetaSentEvent(event)
			if msg != nil {
				break
			}
		}
	}
	zetaHash, err := ob.zetaClient.PostSend(PostSendNonEVMGasLimit, msg)
	if err != nil {
		ob.logger.ExternalChainWatcher.Error().Err(err).Msg("error posting to zeta core")
		return err
	}
	ob.logger.ExternalChainWatcher.Info().Msgf("ZetaSent event detected and reported: PostSend zeta tx: %s", zetaHash)
	return nil
}

func (ob *EVMChainClient) CheckReceiptForCoinTypeERC20(txHash string) error {
	custody, err := ob.GetERC20CustodyContract()
	if err != nil {
		return err
	}
	hash := ethcommon.HexToHash(txHash)
	receipt, err := ob.EvmClient.TransactionReceipt(context.Background(), hash)
	if err != nil {
		return err
	}
	var msg *types.MsgVoteOnObservedInboundTx
	for _, log := range receipt.Logs {
		zetaDeposited, err := custody.ParseDeposited(*log)
		if err != nil {
			msg := ob.GetInboundVoteForDepositedEvent(zetaDeposited)
			if msg != nil {
				break
			}
		}
	}
	zetaHash, err := ob.zetaClient.PostSend(PostSendEVMGasLimit, msg)
	if err != nil {
		ob.logger.ExternalChainWatcher.Error().Err(err).Msg("error posting to zeta core")
		return err
	}
	ob.logger.ExternalChainWatcher.Info().Msgf("ZetaSent event detected and reported: PostSend zeta tx: %s", zetaHash)
	return nil
}

func (ob *EVMChainClient) CheckReceiptForCoinTypeGas(txHash string) error {
	hash := ethcommon.HexToHash(txHash)
	tx, isPending, err := ob.EvmClient.TransactionByHash(context.Background(), hash)
	if err != nil {
		return err
	}
	if isPending {
		return errors.New("tx is still pending")
	}

	receipt, err := ob.EvmClient.TransactionReceipt(context.Background(), hash)
	if err != nil {
		ob.logger.ExternalChainWatcher.Err(err).Msg("TransactionReceipt error")
		return err
	}
	if receipt.Status != 1 { // 1: successful, 0: failed
		ob.logger.ExternalChainWatcher.Info().Msgf("tx %s failed; don't act", tx.Hash().Hex())
		return errors.New("Tx not successfull yet")
	}
	block, err := ob.EvmClient.BlockByNumber(context.Background(), receipt.BlockNumber)
	if err != nil {
		ob.logger.ExternalChainWatcher.Err(err).Msg("BlockByNumber error")
		return err
	}
	from, err := ob.EvmClient.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
	if err != nil {
		ob.logger.ExternalChainWatcher.Err(err).Msg("TransactionSender error; trying local recovery (assuming LondonSigner dynamic fee tx type) of sender address")
		signer := ethtypes.NewLondonSigner(big.NewInt(ob.chain.ChainId))
		from, err = signer.Sender(tx)
		if err != nil {
			ob.logger.ExternalChainWatcher.Err(err).Msg("local recovery of sender address failed")
			return err
		}
	}
	msg := ob.GetInboundVoteForTokenSentToTSS(tx.Hash(), tx.Value(), receipt, from, tx.Data())
	zetaHash, err := ob.zetaClient.PostSend(PostSendEVMGasLimit, msg)
	if err != nil {
		ob.logger.ExternalChainWatcher.Error().Err(err).Msg("error posting to zeta core")
		return err
	}
	ob.logger.ExternalChainWatcher.Info().Msgf("Gas Deposit detected and reported: PostSend zeta tx: %s", zetaHash)
	return nil
}
