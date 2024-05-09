package evm

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	sdkmath "cosmossdk.io/math"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/onrik/ethrpc"
	"github.com/pkg/errors"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.non-eth.sol"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/pkg/constant"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/compliance"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	corecontext "github.com/zeta-chain/zetacore/zetaclient/core_context"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
	"github.com/zeta-chain/zetacore/zetaclient/zetabridge"
	"golang.org/x/net/context"
)

// WatchInboundTracker gets a list of Inbound tracker suggestions from zeta-core at each tick and tries to check if the in-tx was confirmed.
// If it was, it tries to broadcast the confirmation vote. If this zeta client has previously broadcast the vote, the tx would be rejected
func (ob *ChainClient) WatchInboundTracker() {
	ticker, err := clienttypes.NewDynamicTicker(
		fmt.Sprintf("EVM_WatchInboundTracker_%d", ob.chain.ChainId),
		ob.GetChainParams().InboundTicker,
	)
	if err != nil {
		ob.logger.Inbound.Err(err).Msg("error creating ticker")
		return
	}
	defer ticker.Stop()

	ob.logger.Inbound.Info().Msgf("Intx tracker watcher started for chain %d", ob.chain.ChainId)
	for {
		select {
		case <-ticker.C():
			if !corecontext.IsInboundObservationEnabled(ob.coreContext, ob.GetChainParams()) {
				continue
			}
			err := ob.ObserveInboundTrackers()
			if err != nil {
				ob.logger.Inbound.Err(err).Msg("ObserveTrackerSuggestions error")
			}
			ticker.UpdateInterval(ob.GetChainParams().InboundTicker, ob.logger.Inbound)
		case <-ob.stop:
			ob.logger.Inbound.Info().Msg("ExternalChainWatcher for inboundTrackerSuggestions stopped")
			return
		}
	}
}

// ObserveInboundTrackers observes the inbound trackers for the chain
func (ob *ChainClient) ObserveInboundTrackers() error {
	trackers, err := ob.zetaBridge.GetInboundTrackersForChain(ob.chain.ChainId)
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
		ob.logger.Inbound.Info().Msgf("checking tracker for intx %s chain %d", tracker.TxHash, ob.chain.ChainId)

		// check and vote on inbound tx
		switch tracker.CoinType {
		case coin.CoinType_Zeta:
			_, err = ob.CheckAndVoteInboundTokenZeta(tx, receipt, true)
		case coin.CoinType_ERC20:
			_, err = ob.CheckAndVoteInboundTokenERC20(tx, receipt, true)
		case coin.CoinType_Gas:
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
	var msg *types.MsgVoteInbound
	for _, log := range receipt.Logs {
		event, err := connector.ParseZetaSent(*log)
		if err == nil && event != nil {
			// sanity check tx event
			err = ValidateEvmTxLog(&event.Raw, addrConnector, tx.Hash, TopicsZetaSent)
			if err == nil {
				msg = ob.BuildInboundVoteMsgForZetaSentEvent(event)
			} else {
				ob.logger.Inbound.Error().Err(err).Msgf("CheckEvmTxLog error on intx %s chain %d", tx.Hash, ob.chain.ChainId)
				return "", err
			}
			break // only one event is allowed per tx
		}
	}
	if msg == nil {
		// no event, restricted tx, etc.
		ob.logger.Inbound.Info().Msgf("no ZetaSent event found for intx %s chain %d", tx.Hash, ob.chain.ChainId)
		return "", nil
	}
	if vote {
		return ob.PostVoteInbound(msg, coin.CoinType_Zeta, zetabridge.PostVoteInboundMessagePassingExecutionGasLimit)
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
	var msg *types.MsgVoteInbound
	for _, log := range receipt.Logs {
		zetaDeposited, err := custody.ParseDeposited(*log)
		if err == nil && zetaDeposited != nil {
			// sanity check tx event
			err = ValidateEvmTxLog(&zetaDeposited.Raw, addrCustory, tx.Hash, TopicsDeposited)
			if err == nil {
				msg = ob.BuildInboundVoteMsgForDepositedEvent(zetaDeposited, sender)
			} else {
				ob.logger.Inbound.Error().Err(err).Msgf("CheckEvmTxLog error on intx %s chain %d", tx.Hash, ob.chain.ChainId)
				return "", err
			}
			break // only one event is allowed per tx
		}
	}
	if msg == nil {
		// no event, donation, restricted tx, etc.
		ob.logger.Inbound.Info().Msgf("no Deposited event found for intx %s chain %d", tx.Hash, ob.chain.ChainId)
		return "", nil
	}
	if vote {
		return ob.PostVoteInbound(msg, coin.CoinType_ERC20, zetabridge.PostVoteInboundExecutionGasLimit)
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
		ob.logger.Inbound.Info().Msgf("no vote message built for intx %s chain %d", tx.Hash, ob.chain.ChainId)
		return "", nil
	}
	if vote {
		return ob.PostVoteInbound(msg, coin.CoinType_Gas, zetabridge.PostVoteInboundExecutionGasLimit)
	}

	return msg.Digest(), nil
}

// PostVoteInbound posts a vote for the given vote message
func (ob *ChainClient) PostVoteInbound(msg *types.MsgVoteInbound, coinType coin.CoinType, retryGasLimit uint64) (string, error) {
	txHash := msg.InboundHash
	chainID := ob.chain.ChainId
	zetaHash, ballot, err := ob.zetaBridge.PostVoteInbound(zetabridge.PostVoteInboundGasLimit, retryGasLimit, msg)
	if err != nil {
		ob.logger.Inbound.Err(err).Msgf("intx detected: error posting vote for chain %d token %s intx %s", chainID, coinType, txHash)
		return "", err
	} else if zetaHash != "" {
		ob.logger.Inbound.Info().Msgf("intx detected: chain %d token %s intx %s vote %s ballot %s", chainID, coinType, txHash, zetaHash, ballot)
	} else {
		ob.logger.Inbound.Info().Msgf("intx detected: chain %d token %s intx %s already voted on ballot %s", chainID, coinType, txHash, ballot)
	}

	return ballot, err
}

// HasEnoughConfirmations checks if the given receipt has enough confirmations
func (ob *ChainClient) HasEnoughConfirmations(receipt *ethtypes.Receipt, lastHeight uint64) bool {
	confHeight := receipt.BlockNumber.Uint64() + ob.GetChainParams().ConfirmationCount
	return lastHeight >= confHeight
}

// BuildInboundVoteMsgForDepositedEvent builds a inbound vote message for a Deposited event
func (ob *ChainClient) BuildInboundVoteMsgForDepositedEvent(event *erc20custody.ERC20CustodyDeposited, sender ethcommon.Address) *types.MsgVoteInbound {
	// compliance check
	maybeReceiver := ""
	parsedAddress, _, err := chains.ParseAddressAndData(hex.EncodeToString(event.Message))
	if err == nil && parsedAddress != (ethcommon.Address{}) {
		maybeReceiver = parsedAddress.Hex()
	}
	if config.ContainRestrictedAddress(sender.Hex(), clienttypes.BytesToEthHex(event.Recipient), maybeReceiver) {
		compliance.PrintComplianceLog(ob.logger.Inbound, ob.logger.Compliance,
			false, ob.chain.ChainId, event.Raw.TxHash.Hex(), sender.Hex(), clienttypes.BytesToEthHex(event.Recipient), "ERC20")
		return nil
	}

	// donation check
	if bytes.Equal(event.Message, []byte(constant.DonationMessage)) {
		ob.logger.Inbound.Info().Msgf("thank you rich folk for your donation! tx %s chain %d", event.Raw.TxHash.Hex(), ob.chain.ChainId)
		return nil
	}
	message := hex.EncodeToString(event.Message)
	ob.logger.Inbound.Info().Msgf("ERC20CustodyDeposited inTx detected on chain %d tx %s block %d from %s value %s message %s",
		ob.chain.ChainId, event.Raw.TxHash.Hex(), event.Raw.BlockNumber, sender.Hex(), event.Amount.String(), message)

	return zetabridge.GetInboundVoteMessage(
		sender.Hex(),
		ob.chain.ChainId,
		"",
		clienttypes.BytesToEthHex(event.Recipient),
		ob.zetaBridge.ZetaChain().ChainId,
		sdkmath.NewUintFromBigInt(event.Amount),
		hex.EncodeToString(event.Message),
		event.Raw.TxHash.Hex(),
		event.Raw.BlockNumber,
		1_500_000,
		coin.CoinType_ERC20,
		event.Asset.String(),
		ob.zetaBridge.GetKeys().GetOperatorAddress().String(),
		event.Raw.Index,
	)
}

// BuildInboundVoteMsgForZetaSentEvent builds a inbound vote message for a ZetaSent event
func (ob *ChainClient) BuildInboundVoteMsgForZetaSentEvent(event *zetaconnector.ZetaConnectorNonEthZetaSent) *types.MsgVoteInbound {
	destChain := chains.GetChainFromChainID(event.DestinationChainId.Int64())
	if destChain == nil {
		ob.logger.Inbound.Warn().Msgf("chain id not supported  %d", event.DestinationChainId.Int64())
		return nil
	}
	destAddr := clienttypes.BytesToEthHex(event.DestinationAddress)

	// compliance check
	sender := event.ZetaTxSenderAddress.Hex()
	if config.ContainRestrictedAddress(sender, destAddr, event.SourceTxOriginAddress.Hex()) {
		compliance.PrintComplianceLog(ob.logger.Inbound, ob.logger.Compliance,
			false, ob.chain.ChainId, event.Raw.TxHash.Hex(), sender, destAddr, "Zeta")
		return nil
	}

	if !destChain.IsZetaChain() {
		paramsDest, found := ob.coreContext.GetEVMChainParams(destChain.ChainId)
		if !found {
			ob.logger.Inbound.Warn().Msgf("chain id not present in EVMChainParams  %d", event.DestinationChainId.Int64())
			return nil
		}

		if strings.EqualFold(destAddr, paramsDest.ZetaTokenContractAddress) {
			ob.logger.Inbound.Warn().Msgf("potential attack attempt: %s destination address is ZETA token contract address %s", destChain, destAddr)
			return nil
		}
	}
	message := base64.StdEncoding.EncodeToString(event.Message)
	ob.logger.Inbound.Info().Msgf("ZetaSent inTx detected on chain %d tx %s block %d from %s value %s message %s",
		ob.chain.ChainId, event.Raw.TxHash.Hex(), event.Raw.BlockNumber, sender, event.ZetaValueAndGas.String(), message)

	return zetabridge.GetInboundVoteMessage(
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
		coin.CoinType_Zeta,
		"",
		ob.zetaBridge.GetKeys().GetOperatorAddress().String(),
		event.Raw.Index,
	)
}

// BuildInboundVoteMsgForTokenSentToTSS builds a inbound vote message for a token sent to TSS
func (ob *ChainClient) BuildInboundVoteMsgForTokenSentToTSS(tx *ethrpc.Transaction, sender ethcommon.Address, blockNumber uint64) *types.MsgVoteInbound {
	message := tx.Input

	// compliance check
	maybeReceiver := ""
	parsedAddress, _, err := chains.ParseAddressAndData(message)
	if err == nil && parsedAddress != (ethcommon.Address{}) {
		maybeReceiver = parsedAddress.Hex()
	}
	if config.ContainRestrictedAddress(sender.Hex(), maybeReceiver) {
		compliance.PrintComplianceLog(ob.logger.Inbound, ob.logger.Compliance,
			false, ob.chain.ChainId, tx.Hash, sender.Hex(), sender.Hex(), "Gas")
		return nil
	}

	// donation check
	// #nosec G703 err is already checked
	data, _ := hex.DecodeString(message)
	if bytes.Equal(data, []byte(constant.DonationMessage)) {
		ob.logger.Inbound.Info().Msgf("thank you rich folk for your donation! tx %s chain %d", tx.Hash, ob.chain.ChainId)
		return nil
	}
	ob.logger.Inbound.Info().Msgf("TSS inTx detected on chain %d tx %s block %d from %s value %s message %s",
		ob.chain.ChainId, tx.Hash, blockNumber, sender.Hex(), tx.Value.String(), message)

	return zetabridge.GetInboundVoteMessage(
		sender.Hex(),
		ob.chain.ChainId,
		sender.Hex(),
		sender.Hex(),
		ob.zetaBridge.ZetaChain().ChainId,
		sdkmath.NewUintFromBigInt(&tx.Value),
		message,
		tx.Hash,
		blockNumber,
		90_000,
		coin.CoinType_Gas,
		"",
		ob.zetaBridge.GetKeys().GetOperatorAddress().String(),
		0, // not a smart contract call
	)
}

// ObserveTSSReceiveInBlock queries the incoming gas asset to TSS address in a single block and posts votes
func (ob *ChainClient) ObserveTSSReceiveInBlock(blockNumber uint64) error {
	block, err := ob.GetBlockByNumberCached(blockNumber)
	if err != nil {
		return errors.Wrapf(err, "error getting block %d for chain %d", blockNumber, ob.chain.ChainId)
	}

	for i := range block.Transactions {
		tx := block.Transactions[i]
		if ethcommon.HexToAddress(tx.To) == ob.Tss.EVMAddress() {
			receipt, err := ob.evmClient.TransactionReceipt(context.Background(), ethcommon.HexToHash(tx.Hash))
			if err != nil {
				return errors.Wrapf(err, "error getting receipt for intx %s chain %d", tx.Hash, ob.chain.ChainId)
			}

			_, err = ob.CheckAndVoteInboundTokenGas(&tx, receipt, true)
			if err != nil {
				return errors.Wrapf(err, "error checking and voting inbound gas asset for intx %s chain %d", tx.Hash, ob.chain.ChainId)
			}
		}
	}
	return nil
}
