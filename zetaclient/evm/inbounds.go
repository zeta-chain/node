package evm

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	sdkmath "cosmossdk.io/math"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/onrik/ethrpc"
	"github.com/pkg/errors"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.non-eth.sol"
	clientcommon "github.com/zeta-chain/zetacore/zetaclient/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
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
			err := ob.ObserveIntxTrackers()
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

// ObserveIntxTrackers observes the inbound trackers for the chain
func (ob *ChainClient) ObserveIntxTrackers() error {
	trackers, err := ob.zetaClient.GetInboundTrackersForChain(ob.chain.ChainId)
	if err != nil {
		return err
	}
	for _, tracker := range trackers {
		// query tx and receipt
		tx, _, err := ob.TransactionByHash(tracker.TxHash)
		if err != nil {
			return errors.Wrapf(err, "error getting transaction for intx %s chain %d", tracker.TxHash, ob.chain.ChainId)
		}
		receipt, err := ob.evmClient.TransactionReceipt(context.Background(), ethcommon.HexToHash(tracker.TxHash))
		if err != nil {
			return errors.Wrapf(err, "error getting receipt for intx %s chain %d", tracker.TxHash, ob.chain.ChainId)
		}
		ob.logger.ExternalChainWatcher.Info().Msgf("checking tracker for intx %s chain %d", tracker.TxHash, ob.chain.ChainId)

		// check and vote on inbound tx
		switch tracker.CoinType {
		case common.CoinType_Zeta:
			_, err = ob.CheckAndVoteInboundTokenZeta(tx, receipt, true)
		case common.CoinType_ERC20:
			_, err = ob.CheckAndVoteInboundTokenERC20(tx, receipt, true)
		case common.CoinType_Gas:
			_, err = ob.CheckAndVoteInboundTokenGas(tx, receipt, true)
		default:
			return fmt.Errorf("unknown coin type %s for intx %s chain %d", tracker.CoinType, tx.Hash, ob.chain.ChainId)
		}
		if err != nil {
			return errors.Wrapf(err, "error checking and voting for intx %s chain %d", tx.Hash, ob.chain.ChainId)
		}
	}
	return nil
}

// CheckAndVoteInboundTokenZeta checks and votes on the given inbound Zeta token
func (ob *ChainClient) CheckAndVoteInboundTokenZeta(tx *ethrpc.Transaction, receipt *ethtypes.Receipt, vote bool) (string, error) {
	// check confirmations
	if confirmed := ob.HasEnoughConfirmations(receipt, ob.GetLastBlockHeight()); !confirmed {
		return "", fmt.Errorf("intx %s has not been confirmed yet: receipt block %d", tx.Hash, receipt.BlockNumber.Uint64())
	}
	// get zeta connector contract
	addrConnector, connector, err := ob.GetConnectorContract()
	if err != nil {
		return "", err
	}

	// build inbound vote message and post vote
	var msg *types.MsgVoteOnObservedInboundTx
	for _, log := range receipt.Logs {
		event, err := connector.ParseZetaSent(*log)
		if err == nil && event != nil {
			// sanity check tx event
			err = ValidateEvmTxLog(&event.Raw, addrConnector, tx.Hash, TopicsZetaSent)
			if err == nil {
				msg = ob.BuildInboundVoteMsgForZetaSentEvent(event)
			} else {
				ob.logger.ExternalChainWatcher.Error().Err(err).Msgf("CheckEvmTxLog error on intx %s chain %d", tx.Hash, ob.chain.ChainId)
				return "", err
			}
			break // only one event is allowed per tx
		}
	}
	if msg == nil {
		// no event, restricted tx, etc.
		ob.logger.ExternalChainWatcher.Info().Msgf("no ZetaSent event found for intx %s chain %d", tx.Hash, ob.chain.ChainId)
		return "", nil
	}
	if vote {
		return ob.PostVoteInbound(msg, common.CoinType_Zeta, zetabridge.PostVoteInboundMessagePassingExecutionGasLimit)
	}
	return msg.Digest(), nil
}

// CheckAndVoteInboundTokenERC20 checks and votes on the given inbound ERC20 token
func (ob *ChainClient) CheckAndVoteInboundTokenERC20(tx *ethrpc.Transaction, receipt *ethtypes.Receipt, vote bool) (string, error) {
	// check confirmations
	if confirmed := ob.HasEnoughConfirmations(receipt, ob.GetLastBlockHeight()); !confirmed {
		return "", fmt.Errorf("intx %s has not been confirmed yet: receipt block %d", tx.Hash, receipt.BlockNumber.Uint64())
	}

	// get erc20 custody contract
	addrCustory, custody, err := ob.GetERC20CustodyContract()
	if err != nil {
		return "", err
	}
	sender := ethcommon.HexToAddress(tx.From)

	// build inbound vote message and post vote
	var msg *types.MsgVoteOnObservedInboundTx
	for _, log := range receipt.Logs {
		zetaDeposited, err := custody.ParseDeposited(*log)
		if err == nil && zetaDeposited != nil {
			// sanity check tx event
			err = ValidateEvmTxLog(&zetaDeposited.Raw, addrCustory, tx.Hash, TopicsDeposited)
			if err == nil {
				msg = ob.BuildInboundVoteMsgForDepositedEvent(zetaDeposited, sender)
			} else {
				ob.logger.ExternalChainWatcher.Error().Err(err).Msgf("CheckEvmTxLog error on intx %s chain %d", tx.Hash, ob.chain.ChainId)
				return "", err
			}
			break // only one event is allowed per tx
		}
	}
	if msg == nil {
		// no event, donation, restricted tx, etc.
		ob.logger.ExternalChainWatcher.Info().Msgf("no Deposited event found for intx %s chain %d", tx.Hash, ob.chain.ChainId)
		return "", nil
	}
	if vote {
		return ob.PostVoteInbound(msg, common.CoinType_ERC20, zetabridge.PostVoteInboundExecutionGasLimit)
	}
	return msg.Digest(), nil
}

// CheckAndVoteInboundTokenGas checks and votes on the given inbound gas token
func (ob *ChainClient) CheckAndVoteInboundTokenGas(tx *ethrpc.Transaction, receipt *ethtypes.Receipt, vote bool) (string, error) {
	// check confirmations
	if confirmed := ob.HasEnoughConfirmations(receipt, ob.GetLastBlockHeight()); !confirmed {
		return "", fmt.Errorf("intx %s has not been confirmed yet: receipt block %d", tx.Hash, receipt.BlockNumber.Uint64())
	}
	// checks receiver and tx status
	if ethcommon.HexToAddress(tx.To) != ob.Tss.EVMAddress() {
		return "", fmt.Errorf("tx.To %s is not TSS address", tx.To)
	}
	if receipt.Status != ethtypes.ReceiptStatusSuccessful {
		return "", errors.New("not a successful tx")
	}
	sender := ethcommon.HexToAddress(tx.From)

	// build inbound vote message and post vote
	msg := ob.BuildInboundVoteMsgForTokenSentToTSS(tx, sender, receipt.BlockNumber.Uint64())
	if msg == nil {
		// donation, restricted tx, etc.
		ob.logger.ExternalChainWatcher.Info().Msgf("no vote message built for intx %s chain %d", tx.Hash, ob.chain.ChainId)
		return "", nil
	}
	if vote {
		return ob.PostVoteInbound(msg, common.CoinType_Gas, zetabridge.PostVoteInboundExecutionGasLimit)
	}
	return msg.Digest(), nil
}

// PostVoteInbound posts a vote for the given vote message
func (ob *ChainClient) PostVoteInbound(msg *types.MsgVoteOnObservedInboundTx, coinType common.CoinType, retryGasLimit uint64) (string, error) {
	txHash := msg.InTxHash
	chainID := ob.chain.ChainId
	zetaHash, ballot, err := ob.zetaClient.PostVoteInbound(zetabridge.PostVoteInboundGasLimit, retryGasLimit, msg)
	if err != nil {
		ob.logger.ExternalChainWatcher.Err(err).Msgf("intx detected: error posting vote for chain %d token %s intx %s", chainID, coinType, txHash)
		return "", err
	} else if zetaHash != "" {
		ob.logger.ExternalChainWatcher.Info().Msgf("intx detected: chain %d token %s intx %s vote %s ballot %s", chainID, coinType, txHash, zetaHash, ballot)
	} else {
		ob.logger.ExternalChainWatcher.Info().Msgf("intx detected: chain %d token %s intx %s already voted on ballot %s", chainID, coinType, txHash, ballot)
	}
	return ballot, err
}

// HasEnoughConfirmations checks if the given receipt has enough confirmations
func (ob *ChainClient) HasEnoughConfirmations(receipt *ethtypes.Receipt, lastHeight uint64) bool {
	confHeight := receipt.BlockNumber.Uint64() + ob.GetChainParams().ConfirmationCount
	return lastHeight >= confHeight
}

// BuildInboundVoteMsgForDepositedEvent builds a inbound vote message for a Deposited event
func (ob *ChainClient) BuildInboundVoteMsgForDepositedEvent(event *erc20custody.ERC20CustodyDeposited, sender ethcommon.Address) *types.MsgVoteOnObservedInboundTx {
	// compliance check
	maybeReceiver := ""
	parsedAddress, _, err := common.ParseAddressAndData(hex.EncodeToString(event.Message))
	if err == nil && parsedAddress != (ethcommon.Address{}) {
		maybeReceiver = parsedAddress.Hex()
	}
	if config.ContainRestrictedAddress(sender.Hex(), clienttypes.BytesToEthHex(event.Recipient), maybeReceiver) {
		clientcommon.PrintComplianceLog(ob.logger.ExternalChainWatcher, ob.logger.Compliance,
			false, ob.chain.ChainId, event.Raw.TxHash.Hex(), sender.Hex(), clienttypes.BytesToEthHex(event.Recipient), "ERC20")
		return nil
	}

	// donation check
	if bytes.Equal(event.Message, []byte(common.DonationMessage)) {
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

// BuildInboundVoteMsgForZetaSentEvent builds a inbound vote message for a ZetaSent event
func (ob *ChainClient) BuildInboundVoteMsgForZetaSentEvent(event *zetaconnector.ZetaConnectorNonEthZetaSent) *types.MsgVoteOnObservedInboundTx {
	destChain := common.GetChainFromChainID(event.DestinationChainId.Int64())
	if destChain == nil {
		ob.logger.ExternalChainWatcher.Warn().Msgf("chain id not supported  %d", event.DestinationChainId.Int64())
		return nil
	}
	destAddr := clienttypes.BytesToEthHex(event.DestinationAddress)

	// compliance check
	sender := event.ZetaTxSenderAddress.Hex()
	if config.ContainRestrictedAddress(sender, destAddr, event.SourceTxOriginAddress.Hex()) {
		clientcommon.PrintComplianceLog(ob.logger.ExternalChainWatcher, ob.logger.Compliance,
			false, ob.chain.ChainId, event.Raw.TxHash.Hex(), sender, destAddr, "Zeta")
		return nil
	}

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
		ob.chain.ChainId, event.Raw.TxHash.Hex(), event.Raw.BlockNumber, sender, event.ZetaValueAndGas.String(), message)

	return zetabridge.GetInBoundVoteMessage(
		sender,
		ob.chain.ChainId,
		event.SourceTxOriginAddress.Hex(),
		destAddr,
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

// BuildInboundVoteMsgForTokenSentToTSS builds a inbound vote message for a token sent to TSS
func (ob *ChainClient) BuildInboundVoteMsgForTokenSentToTSS(tx *ethrpc.Transaction, sender ethcommon.Address, blockNumber uint64) *types.MsgVoteOnObservedInboundTx {
	message := tx.Input

	// compliance check
	maybeReceiver := ""
	parsedAddress, _, err := common.ParseAddressAndData(message)
	if err == nil && parsedAddress != (ethcommon.Address{}) {
		maybeReceiver = parsedAddress.Hex()
	}
	if config.ContainRestrictedAddress(sender.Hex(), maybeReceiver) {
		clientcommon.PrintComplianceLog(ob.logger.ExternalChainWatcher, ob.logger.Compliance,
			false, ob.chain.ChainId, tx.Hash, sender.Hex(), sender.Hex(), "Gas")
		return nil
	}

	// donation check
	// #nosec G703 err is already checked
	data, _ := hex.DecodeString(message)
	if bytes.Equal(data, []byte(common.DonationMessage)) {
		ob.logger.ExternalChainWatcher.Info().Msgf("thank you rich folk for your donation! tx %s chain %d", tx.Hash, ob.chain.ChainId)
		return nil
	}
	ob.logger.ExternalChainWatcher.Info().Msgf("TSS inTx detected on chain %d tx %s block %d from %s value %s message %s",
		ob.chain.ChainId, tx.Hash, blockNumber, sender.Hex(), tx.Value.String(), message)

	return zetabridge.GetInBoundVoteMessage(
		sender.Hex(),
		ob.chain.ChainId,
		sender.Hex(),
		sender.Hex(),
		ob.zetaClient.ZetaChain().ChainId,
		sdkmath.NewUintFromBigInt(&tx.Value),
		message,
		tx.Hash,
		blockNumber,
		90_000,
		common.CoinType_Gas,
		"",
		ob.zetaClient.GetKeys().GetOperatorAddress().String(),
		0, // not a smart contract call
	)
}
