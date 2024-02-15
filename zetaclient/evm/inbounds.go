package evm

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	sdkmath "cosmossdk.io/math"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.non-eth.sol"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/zetabridge"
	"golang.org/x/net/context"
)

// ExternalChainWatcherForNewInboundTrackerSuggestions At each tick, gets a list of Inbound tracker suggestions from zeta-core and tries to check if the in-tx was confirmed.
// If it was, it tries to broadcast the confirmation vote. If this zeta client has previously broadcast the vote, the tx would be rejected
func (ob *ChainClient) ExternalChainWatcherForNewInboundTrackerSuggestions() {
	ticker, err := clienttypes.NewDynamicTicker(
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

func (ob *ChainClient) ObserveTrackerSuggestions() error {
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

func (ob *ChainClient) CheckReceiptForCoinTypeZeta(txHash string, vote bool) (string, error) {
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

func (ob *ChainClient) CheckReceiptForCoinTypeERC20(txHash string, vote bool) (string, error) {
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

func (ob *ChainClient) CheckReceiptForCoinTypeGas(txHash string, vote bool) (string, error) {
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

// CheckEvmTxLog checks the basics of an EVM tx log
func (ob *ChainClient) CheckEvmTxLog(vLog *ethtypes.Log, wantAddress ethcommon.Address, wantHash string, wantTopics int) error {
	if vLog.Removed {
		return fmt.Errorf("log is removed, chain reorg?")
	}
	if vLog.Address != wantAddress {
		return fmt.Errorf("log emitter address mismatch: want %s got %s", wantAddress.Hex(), vLog.Address.Hex())
	}
	if vLog.TxHash.Hex() == "" {
		return fmt.Errorf("log tx hash is empty: %d %s", vLog.BlockNumber, vLog.TxHash.Hex())
	}
	if wantHash != "" && vLog.TxHash.Hex() != wantHash {
		return fmt.Errorf("log tx hash mismatch: want %s got %s", wantHash, vLog.TxHash.Hex())
	}
	if len(vLog.Topics) != wantTopics {
		return fmt.Errorf("number of topics mismatch: want %d got %d", wantTopics, len(vLog.Topics))
	}
	return nil
}

// HasEnoughConfirmations checks if the given receipt has enough confirmations
func (ob *ChainClient) HasEnoughConfirmations(receipt *ethtypes.Receipt, lastHeight uint64) bool {
	confHeight := receipt.BlockNumber.Uint64() + ob.GetChainParams().ConfirmationCount
	return lastHeight >= confHeight
}

// GetTransactionSender returns the sender of the given transaction
func (ob *ChainClient) GetTransactionSender(tx *ethtypes.Transaction, blockHash ethcommon.Hash, txIndex uint) (ethcommon.Address, error) {
	sender, err := ob.evmClient.TransactionSender(context.Background(), tx, blockHash, txIndex)
	if err != nil {
		// trying local recovery (assuming LondonSigner dynamic fee tx type) of sender address
		signer := ethtypes.NewLondonSigner(big.NewInt(ob.chain.ChainId))
		sender, err = signer.Sender(tx)
		if err != nil {
			ob.logger.ExternalChainWatcher.Err(err).Msgf("can't recover the sender from tx hash %s chain %d", tx.Hash().Hex(), ob.chain.ChainId)
			return ethcommon.Address{}, err
		}
	}
	return sender, nil
}

func (ob *ChainClient) GetInboundVoteMsgForDepositedEvent(event *erc20custody.ERC20CustodyDeposited, sender ethcommon.Address) *types.MsgVoteOnObservedInboundTx {
	if bytes.Equal(event.Message, []byte(DonationMessage)) {
		ob.logger.ExternalChainWatcher.Info().Msgf("thank you rich folk for your donation! tx %s chain %d", event.Raw.TxHash.Hex(), ob.chain.ChainId)
		return nil
	}
	message := hex.EncodeToString(event.Message)
	ob.logger.ExternalChainWatcher.Info().Msgf("ERC20CustodyDeposited inTx detected on chain %d tx %s block %d from %s value %s message %s",
		ob.chain.ChainId, event.Raw.TxHash.Hex(), event.Raw.BlockNumber, sender.Hex(), event.Amount.String(), message)

	return zetabridge.GetInBoundVoteMessage(
		sender.Hex(),
		ob.chain.ChainId,
		"",
		clienttypes.BytesToEthHex(event.Recipient),
		ob.zetaClient.ZetaChain().ChainId,
		sdkmath.NewUintFromBigInt(event.Amount),
		hex.EncodeToString(event.Message),
		event.Raw.TxHash.Hex(),
		event.Raw.BlockNumber,
		1_500_000,
		common.CoinType_ERC20,
		event.Asset.String(),
		ob.zetaClient.GetKeys().GetOperatorAddress().String(),
		event.Raw.Index,
	)
}

func (ob *ChainClient) GetInboundVoteMsgForZetaSentEvent(event *zetaconnector.ZetaConnectorNonEthZetaSent) *types.MsgVoteOnObservedInboundTx {
	destChain := common.GetChainFromChainID(event.DestinationChainId.Int64())
	if destChain == nil {
		ob.logger.ExternalChainWatcher.Warn().Msgf("chain id not supported  %d", event.DestinationChainId.Int64())
		return nil
	}
	destAddr := clienttypes.BytesToEthHex(event.DestinationAddress)
	if !destChain.IsZetaChain() {
		paramsDest, found := ob.coreContext.GetEVMChainParams(destChain.ChainId)
		if !found {
			ob.logger.ExternalChainWatcher.Warn().Msgf("chain id not present in EVMChainParams  %d", event.DestinationChainId.Int64())
			return nil
		}

		if strings.EqualFold(destAddr, paramsDest.ZetaTokenContractAddress) {
			ob.logger.ExternalChainWatcher.Warn().Msgf("potential attack attempt: %s destination address is ZETA token contract address %s", destChain, destAddr)
			return nil
		}
	}
	message := base64.StdEncoding.EncodeToString(event.Message)
	ob.logger.ExternalChainWatcher.Info().Msgf("ZetaSent inTx detected on chain %d tx %s block %d from %s value %s message %s",
		ob.chain.ChainId, event.Raw.TxHash.Hex(), event.Raw.BlockNumber, event.ZetaTxSenderAddress.Hex(), event.ZetaValueAndGas.String(), message)

	return zetabridge.GetInBoundVoteMessage(
		event.ZetaTxSenderAddress.Hex(),
		ob.chain.ChainId,
		event.SourceTxOriginAddress.Hex(),
		clienttypes.BytesToEthHex(event.DestinationAddress),
		destChain.ChainId,
		sdkmath.NewUintFromBigInt(event.ZetaValueAndGas),
		message,
		event.Raw.TxHash.Hex(),
		event.Raw.BlockNumber,
		event.DestinationGasLimit.Uint64(),
		common.CoinType_Zeta,
		"",
		ob.zetaClient.GetKeys().GetOperatorAddress().String(),
		event.Raw.Index,
	)
}

func (ob *ChainClient) GetInboundVoteMsgForTokenSentToTSS(tx *ethtypes.Transaction, sender ethcommon.Address, blockNumber uint64) *types.MsgVoteOnObservedInboundTx {
	if bytes.Equal(tx.Data(), []byte(DonationMessage)) {
		ob.logger.ExternalChainWatcher.Info().Msgf("thank you rich folk for your donation! tx %s chain %d", tx.Hash().Hex(), ob.chain.ChainId)
		return nil
	}
	message := ""
	if len(tx.Data()) != 0 {
		message = hex.EncodeToString(tx.Data())
	}
	ob.logger.ExternalChainWatcher.Info().Msgf("TSS inTx detected on chain %d tx %s block %d from %s value %s message %s",
		ob.chain.ChainId, tx.Hash().Hex(), blockNumber, sender.Hex(), tx.Value().String(), hex.EncodeToString(tx.Data()))

	return zetabridge.GetInBoundVoteMessage(
		sender.Hex(),
		ob.chain.ChainId,
		sender.Hex(),
		sender.Hex(),
		ob.zetaClient.ZetaChain().ChainId,
		sdkmath.NewUintFromBigInt(tx.Value()),
		message,
		tx.Hash().Hex(),
		blockNumber,
		90_000,
		common.CoinType_Gas,
		"",
		ob.zetaClient.GetKeys().GetOperatorAddress().String(),
		0, // not a smart contract call
	)
}
