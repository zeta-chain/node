package observer

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"strings"

	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/onrik/ethrpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.non-eth.sol"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/pkg/constant"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/evm"
	"github.com/zeta-chain/zetacore/zetaclient/compliance"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	clientcontext "github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
	"github.com/zeta-chain/zetacore/zetaclient/zetacore"
)

// WatchInTx watches evm chain for incoming txs and post votes to zetacore
func (ob *Observer) WatchInTx() {
	ticker, err := clienttypes.NewDynamicTicker(fmt.Sprintf("EVM_WatchInTx_%d", ob.chain.ChainId), ob.GetChainParams().InTxTicker)
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
			if !clientcontext.IsInboundObservationEnabled(ob.coreContext, ob.GetChainParams()) {
				sampledLogger.Info().Msgf("WatchInTx: inbound observation is disabled for chain %d", ob.chain.ChainId)
				continue
			}
			err := ob.ObserveInTX(sampledLogger)
			if err != nil {
				ob.logger.InTx.Err(err).Msg("WatchInTx: observeInTX error")
			}
			ticker.UpdateInterval(ob.GetChainParams().InTxTicker, ob.logger.InTx)
		case <-ob.stop:
			ob.logger.InTx.Info().Msgf("WatchInTx stopped for chain %d", ob.chain.ChainId)
			return
		}
	}
}

// WatchIntxTracker gets a list of Inbound tracker suggestions from zeta-core at each tick and tries to check if the in-tx was confirmed.
// If it was, it tries to broadcast the confirmation vote. If this zeta client has previously broadcast the vote, the tx would be rejected
func (ob *Observer) WatchIntxTracker() {
	ticker, err := clienttypes.NewDynamicTicker(
		fmt.Sprintf("EVM_WatchIntxTracker_%d", ob.chain.ChainId),
		ob.GetChainParams().InTxTicker,
	)
	if err != nil {
		ob.logger.InTx.Err(err).Msg("error creating ticker")
		return
	}
	defer ticker.Stop()

	ob.logger.InTx.Info().Msgf("Inbound tracker watcher started for chain %d", ob.chain.ChainId)
	for {
		select {
		case <-ticker.C():
			if !clientcontext.IsInboundObservationEnabled(ob.coreContext, ob.GetChainParams()) {
				continue
			}
			err := ob.ProcessInboundTrackers()
			if err != nil {
				ob.logger.InTx.Err(err).Msg("ProcessInboundTrackers error")
			}
			ticker.UpdateInterval(ob.GetChainParams().InTxTicker, ob.logger.InTx)
		case <-ob.stop:
			ob.logger.InTx.Info().Msgf("WatchIntxTracker stopped for chain %d", ob.chain.ChainId)
			return
		}
	}
}

// ProcessInboundTrackers processes inbound trackers from zetacore
func (ob *Observer) ProcessInboundTrackers() error {
	trackers, err := ob.zetacoreClient.GetInboundTrackersForChain(ob.chain.ChainId)
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
		ob.logger.InTx.Info().Msgf("checking tracker for intx %s chain %d", tracker.TxHash, ob.chain.ChainId)

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

func (ob *Observer) ObserveInTX(sampledLogger zerolog.Logger) error {
	// get and update latest block height
	blockNumber, err := ob.evmClient.BlockNumber(context.Background())
	if err != nil {
		return err
	}
	if blockNumber < ob.GetLastBlockHeight() {
		return fmt.Errorf("observeInTX: block number should not decrease: current %d last %d", blockNumber, ob.GetLastBlockHeight())
	}
	ob.SetLastBlockHeight(blockNumber)

	// increment prom counter
	metrics.GetBlockByNumberPerChain.WithLabelValues(ob.chain.ChainName.String()).Inc()

	// skip if current height is too low
	if blockNumber < ob.GetChainParams().ConfirmationCount {
		return fmt.Errorf("observeInTX: skipping observer, current block number %d is too low", blockNumber)
	}
	confirmedBlockNum := blockNumber - ob.GetChainParams().ConfirmationCount

	// skip if no new block is confirmed
	lastScanned := ob.GetLastBlockHeightScanned()
	if lastScanned >= confirmedBlockNum {
		sampledLogger.Debug().Msgf("observeInTX: skipping observer, no new block is produced for chain %d", ob.chain.ChainId)
		return nil
	}

	// get last scanned block height (we simply use same height for all 3 events ZetaSent, Deposited, TssRecvd)
	// Note: using different heights for each event incurs more complexity (metrics, db, etc) and not worth it
	startBlock, toBlock := ob.calcBlockRangeToScan(confirmedBlockNum, lastScanned, config.MaxBlocksPerPeriod)

	// task 1:  query evm chain for zeta sent logs (read at most 100 blocks in one go)
	lastScannedZetaSent := ob.ObserveZetaSent(startBlock, toBlock)

	// task 2: query evm chain for deposited logs (read at most 100 blocks in one go)
	lastScannedDeposited := ob.ObserveERC20Deposited(startBlock, toBlock)

	// task 3: query the incoming tx to TSS address (read at most 100 blocks in one go)
	lastScannedTssRecvd := ob.ObserverTSSReceive(startBlock, toBlock)

	// note: using lowest height for all 3 events is not perfect, but it's simple and good enough
	lastScannedLowest := lastScannedZetaSent
	if lastScannedDeposited < lastScannedLowest {
		lastScannedLowest = lastScannedDeposited
	}
	if lastScannedTssRecvd < lastScannedLowest {
		lastScannedLowest = lastScannedTssRecvd
	}

	// update last scanned block height for all 3 events (ZetaSent, Deposited, TssRecvd), ignore db error
	if lastScannedLowest > lastScanned {
		sampledLogger.Info().Msgf("observeInTX: lasstScanned heights for chain %d ZetaSent %d ERC20Deposited %d TssRecvd %d",
			ob.chain.ChainId, lastScannedZetaSent, lastScannedDeposited, lastScannedTssRecvd)
		ob.SetLastBlockHeightScanned(lastScannedLowest)
		if err := ob.db.Save(clienttypes.ToLastBlockSQLType(lastScannedLowest)).Error; err != nil {
			ob.logger.InTx.Error().Err(err).Msgf("observeInTX: error writing lastScannedLowest %d to db", lastScannedLowest)
		}
	}
	return nil
}

// ObserveZetaSent queries the ZetaSent event from the connector contract and posts to zetacore
// returns the last block successfully scanned
func (ob *Observer) ObserveZetaSent(startBlock, toBlock uint64) uint64 {
	// filter ZetaSent logs
	addrConnector, connector, err := ob.GetConnectorContract()
	if err != nil {
		ob.logger.Chain.Warn().Err(err).Msgf("ObserveZetaSent: GetConnectorContract error:")
		return startBlock - 1 // lastScanned
	}
	iter, err := connector.FilterZetaSent(&bind.FilterOpts{
		Start:   startBlock,
		End:     &toBlock,
		Context: context.TODO(),
	}, []ethcommon.Address{}, []*big.Int{})
	if err != nil {
		ob.logger.Chain.Warn().Err(err).Msgf(
			"ObserveZetaSent: FilterZetaSent error from block %d to %d for chain %d", startBlock, toBlock, ob.chain.ChainId)
		return startBlock - 1 // lastScanned
	}

	// collect and sort events by block number, then tx index, then log index (ascending)
	events := make([]*zetaconnector.ZetaConnectorNonEthZetaSent, 0)
	for iter.Next() {
		// sanity check tx event
		err := evm.ValidateEvmTxLog(&iter.Event.Raw, addrConnector, "", evm.TopicsZetaSent)
		if err == nil {
			events = append(events, iter.Event)
			continue
		}
		ob.logger.InTx.Warn().Err(err).Msgf("ObserveZetaSent: invalid ZetaSent event in tx %s on chain %d at height %d",
			iter.Event.Raw.TxHash.Hex(), ob.chain.ChainId, iter.Event.Raw.BlockNumber)
	}
	sort.SliceStable(events, func(i, j int) bool {
		if events[i].Raw.BlockNumber == events[j].Raw.BlockNumber {
			if events[i].Raw.TxIndex == events[j].Raw.TxIndex {
				return events[i].Raw.Index < events[j].Raw.Index
			}
			return events[i].Raw.TxIndex < events[j].Raw.TxIndex
		}
		return events[i].Raw.BlockNumber < events[j].Raw.BlockNumber
	})

	// increment prom counter
	metrics.GetFilterLogsPerChain.WithLabelValues(ob.chain.ChainName.String()).Inc()

	// post to zetacore
	beingScanned := uint64(0)
	guard := make(map[string]bool)
	for _, event := range events {
		// remember which block we are scanning (there could be multiple events in the same block)
		if event.Raw.BlockNumber > beingScanned {
			beingScanned = event.Raw.BlockNumber
		}
		// guard against multiple events in the same tx
		if guard[event.Raw.TxHash.Hex()] {
			ob.logger.InTx.Warn().Msgf("ObserveZetaSent: multiple remote call events detected in tx %s", event.Raw.TxHash)
			continue
		}
		guard[event.Raw.TxHash.Hex()] = true

		msg := ob.BuildInboundVoteMsgForZetaSentEvent(event)
		if msg != nil {
			_, err = ob.PostVoteInbound(msg, coin.CoinType_Zeta, zetacore.PostVoteInboundMessagePassingExecutionGasLimit)
			if err != nil {
				return beingScanned - 1 // we have to re-scan from this block next time
			}
		}
	}
	// successful processed all events in [startBlock, toBlock]
	return toBlock
}

// ObserveERC20Deposited queries the ERC20CustodyDeposited event from the ERC20Custody contract and posts to zetacore
// returns the last block successfully scanned
func (ob *Observer) ObserveERC20Deposited(startBlock, toBlock uint64) uint64 {
	// filter ERC20CustodyDeposited logs
	addrCustody, erc20custodyContract, err := ob.GetERC20CustodyContract()
	if err != nil {
		ob.logger.InTx.Warn().Err(err).Msgf("ObserveERC20Deposited: GetERC20CustodyContract error:")
		return startBlock - 1 // lastScanned
	}

	iter, err := erc20custodyContract.FilterDeposited(&bind.FilterOpts{
		Start:   startBlock,
		End:     &toBlock,
		Context: context.TODO(),
	}, []ethcommon.Address{})
	if err != nil {
		ob.logger.InTx.Warn().Err(err).Msgf(
			"ObserveERC20Deposited: FilterDeposited error from block %d to %d for chain %d", startBlock, toBlock, ob.chain.ChainId)
		return startBlock - 1 // lastScanned
	}

	// collect and sort events by block number, then tx index, then log index (ascending)
	events := make([]*erc20custody.ERC20CustodyDeposited, 0)
	for iter.Next() {
		// sanity check tx event
		err := evm.ValidateEvmTxLog(&iter.Event.Raw, addrCustody, "", evm.TopicsDeposited)
		if err == nil {
			events = append(events, iter.Event)
			continue
		}
		ob.logger.InTx.Warn().Err(err).Msgf("ObserveERC20Deposited: invalid Deposited event in tx %s on chain %d at height %d",
			iter.Event.Raw.TxHash.Hex(), ob.chain.ChainId, iter.Event.Raw.BlockNumber)
	}
	sort.SliceStable(events, func(i, j int) bool {
		if events[i].Raw.BlockNumber == events[j].Raw.BlockNumber {
			if events[i].Raw.TxIndex == events[j].Raw.TxIndex {
				return events[i].Raw.Index < events[j].Raw.Index
			}
			return events[i].Raw.TxIndex < events[j].Raw.TxIndex
		}
		return events[i].Raw.BlockNumber < events[j].Raw.BlockNumber
	})

	// increment prom counter
	metrics.GetFilterLogsPerChain.WithLabelValues(ob.chain.ChainName.String()).Inc()

	// post to zeatcore
	guard := make(map[string]bool)
	beingScanned := uint64(0)
	for _, event := range events {
		// remember which block we are scanning (there could be multiple events in the same block)
		if event.Raw.BlockNumber > beingScanned {
			beingScanned = event.Raw.BlockNumber
		}
		tx, _, err := ob.TransactionByHash(event.Raw.TxHash.Hex())
		if err != nil {
			ob.logger.InTx.Error().Err(err).Msgf(
				"ObserveERC20Deposited: error getting transaction for intx %s chain %d", event.Raw.TxHash, ob.chain.ChainId)
			return beingScanned - 1 // we have to re-scan from this block next time
		}
		sender := ethcommon.HexToAddress(tx.From)

		// guard against multiple events in the same tx
		if guard[event.Raw.TxHash.Hex()] {
			ob.logger.InTx.Warn().Msgf("ObserveERC20Deposited: multiple remote call events detected in tx %s", event.Raw.TxHash)
			continue
		}
		guard[event.Raw.TxHash.Hex()] = true

		msg := ob.BuildInboundVoteMsgForDepositedEvent(event, sender)
		if msg != nil {
			_, err = ob.PostVoteInbound(msg, coin.CoinType_ERC20, zetacore.PostVoteInboundExecutionGasLimit)
			if err != nil {
				return beingScanned - 1 // we have to re-scan from this block next time
			}
		}
	}
	// successful processed all events in [startBlock, toBlock]
	return toBlock
}

// ObserverTSSReceive queries the incoming gas asset to TSS address and posts to zetacore
// returns the last block successfully scanned
func (ob *Observer) ObserverTSSReceive(startBlock, toBlock uint64) uint64 {
	// query incoming gas asset
	for bn := startBlock; bn <= toBlock; bn++ {
		// post new block header (if any) to zetacore and ignore error
		// TODO: consider having a independent ticker(from TSS scaning) for posting block headers
		// https://github.com/zeta-chain/node/issues/1847
		blockHeaderVerification, found := ob.coreContext.GetBlockHeaderEnabledChains(ob.chain.ChainId)
		if found && blockHeaderVerification.Enabled {
			// post block header for supported chains
			err := ob.postBlockHeader(toBlock)
			if err != nil {
				ob.logger.InTx.Error().Err(err).Msg("error posting block header")
			}
		}

		// observe TSS received gas token in block 'bn'
		err := ob.ObserveTSSReceiveInBlock(bn)
		if err != nil {
			ob.logger.InTx.Error().Err(err).Msgf("ObserverTSSReceive: error observing TSS received token in block %d for chain %d", bn, ob.chain.ChainId)
			return bn - 1 // we have to re-scan from this block next time
		}
	}
	// successful processed all gas asset deposits in [startBlock, toBlock]
	return toBlock
}

// CheckAndVoteInboundTokenZeta checks and votes on the given inbound Zeta token
func (ob *Observer) CheckAndVoteInboundTokenZeta(tx *ethrpc.Transaction, receipt *ethtypes.Receipt, vote bool) (string, error) {
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
			err = evm.ValidateEvmTxLog(&event.Raw, addrConnector, tx.Hash, evm.TopicsZetaSent)
			if err == nil {
				msg = ob.BuildInboundVoteMsgForZetaSentEvent(event)
			} else {
				ob.logger.InTx.Error().Err(err).Msgf("CheckEvmTxLog error on intx %s chain %d", tx.Hash, ob.chain.ChainId)
				return "", err
			}
			break // only one event is allowed per tx
		}
	}
	if msg == nil {
		// no event, restricted tx, etc.
		ob.logger.InTx.Info().Msgf("no ZetaSent event found for intx %s chain %d", tx.Hash, ob.chain.ChainId)
		return "", nil
	}
	if vote {
		return ob.PostVoteInbound(msg, coin.CoinType_Zeta, zetacore.PostVoteInboundMessagePassingExecutionGasLimit)
	}

	return msg.Digest(), nil
}

// CheckAndVoteInboundTokenERC20 checks and votes on the given inbound ERC20 token
func (ob *Observer) CheckAndVoteInboundTokenERC20(tx *ethrpc.Transaction, receipt *ethtypes.Receipt, vote bool) (string, error) {
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
			err = evm.ValidateEvmTxLog(&zetaDeposited.Raw, addrCustory, tx.Hash, evm.TopicsDeposited)
			if err == nil {
				msg = ob.BuildInboundVoteMsgForDepositedEvent(zetaDeposited, sender)
			} else {
				ob.logger.InTx.Error().Err(err).Msgf("CheckEvmTxLog error on intx %s chain %d", tx.Hash, ob.chain.ChainId)
				return "", err
			}
			break // only one event is allowed per tx
		}
	}
	if msg == nil {
		// no event, donation, restricted tx, etc.
		ob.logger.InTx.Info().Msgf("no Deposited event found for intx %s chain %d", tx.Hash, ob.chain.ChainId)
		return "", nil
	}
	if vote {
		return ob.PostVoteInbound(msg, coin.CoinType_ERC20, zetacore.PostVoteInboundExecutionGasLimit)
	}

	return msg.Digest(), nil
}

// CheckAndVoteInboundTokenGas checks and votes on the given inbound gas token
func (ob *Observer) CheckAndVoteInboundTokenGas(tx *ethrpc.Transaction, receipt *ethtypes.Receipt, vote bool) (string, error) {
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
		ob.logger.InTx.Info().Msgf("no vote message built for intx %s chain %d", tx.Hash, ob.chain.ChainId)
		return "", nil
	}
	if vote {
		return ob.PostVoteInbound(msg, coin.CoinType_Gas, zetacore.PostVoteInboundExecutionGasLimit)
	}

	return msg.Digest(), nil
}

// PostVoteInbound posts a vote for the given vote message
func (ob *Observer) PostVoteInbound(msg *types.MsgVoteOnObservedInboundTx, coinType coin.CoinType, retryGasLimit uint64) (string, error) {
	txHash := msg.InTxHash
	chainID := ob.chain.ChainId
	zetaHash, ballot, err := ob.zetacoreClient.PostVoteInbound(zetacore.PostVoteInboundGasLimit, retryGasLimit, msg)
	if err != nil {
		ob.logger.InTx.Err(err).Msgf("intx detected: error posting vote for chain %d token %s intx %s", chainID, coinType, txHash)
		return "", err
	} else if zetaHash != "" {
		ob.logger.InTx.Info().Msgf("intx detected: chain %d token %s intx %s vote %s ballot %s", chainID, coinType, txHash, zetaHash, ballot)
	} else {
		ob.logger.InTx.Info().Msgf("intx detected: chain %d token %s intx %s already voted on ballot %s", chainID, coinType, txHash, ballot)
	}

	return ballot, err
}

// HasEnoughConfirmations checks if the given receipt has enough confirmations
func (ob *Observer) HasEnoughConfirmations(receipt *ethtypes.Receipt, lastHeight uint64) bool {
	confHeight := receipt.BlockNumber.Uint64() + ob.GetChainParams().ConfirmationCount
	return lastHeight >= confHeight
}

// BuildInboundVoteMsgForDepositedEvent builds a inbound vote message for a Deposited event
func (ob *Observer) BuildInboundVoteMsgForDepositedEvent(event *erc20custody.ERC20CustodyDeposited, sender ethcommon.Address) *types.MsgVoteOnObservedInboundTx {
	// compliance check
	maybeReceiver := ""
	parsedAddress, _, err := chains.ParseAddressAndData(hex.EncodeToString(event.Message))
	if err == nil && parsedAddress != (ethcommon.Address{}) {
		maybeReceiver = parsedAddress.Hex()
	}
	if config.ContainRestrictedAddress(sender.Hex(), clienttypes.BytesToEthHex(event.Recipient), maybeReceiver) {
		compliance.PrintComplianceLog(ob.logger.InTx, ob.logger.Compliance,
			false, ob.chain.ChainId, event.Raw.TxHash.Hex(), sender.Hex(), clienttypes.BytesToEthHex(event.Recipient), "ERC20")
		return nil
	}

	// donation check
	if bytes.Equal(event.Message, []byte(constant.DonationMessage)) {
		ob.logger.InTx.Info().Msgf("thank you rich folk for your donation! tx %s chain %d", event.Raw.TxHash.Hex(), ob.chain.ChainId)
		return nil
	}
	message := hex.EncodeToString(event.Message)
	ob.logger.InTx.Info().Msgf("ERC20CustodyDeposited inTx detected on chain %d tx %s block %d from %s value %s message %s",
		ob.chain.ChainId, event.Raw.TxHash.Hex(), event.Raw.BlockNumber, sender.Hex(), event.Amount.String(), message)

	return zetacore.GetInBoundVoteMessage(
		sender.Hex(),
		ob.chain.ChainId,
		"",
		clienttypes.BytesToEthHex(event.Recipient),
		ob.zetacoreClient.Chain().ChainId,
		sdkmath.NewUintFromBigInt(event.Amount),
		hex.EncodeToString(event.Message),
		event.Raw.TxHash.Hex(),
		event.Raw.BlockNumber,
		1_500_000,
		coin.CoinType_ERC20,
		event.Asset.String(),
		ob.zetacoreClient.GetKeys().GetOperatorAddress().String(),
		event.Raw.Index,
	)
}

// BuildInboundVoteMsgForZetaSentEvent builds a inbound vote message for a ZetaSent event
func (ob *Observer) BuildInboundVoteMsgForZetaSentEvent(event *zetaconnector.ZetaConnectorNonEthZetaSent) *types.MsgVoteOnObservedInboundTx {
	destChain := chains.GetChainFromChainID(event.DestinationChainId.Int64())
	if destChain == nil {
		ob.logger.InTx.Warn().Msgf("chain id not supported  %d", event.DestinationChainId.Int64())
		return nil
	}
	destAddr := clienttypes.BytesToEthHex(event.DestinationAddress)

	// compliance check
	sender := event.ZetaTxSenderAddress.Hex()
	if config.ContainRestrictedAddress(sender, destAddr, event.SourceTxOriginAddress.Hex()) {
		compliance.PrintComplianceLog(ob.logger.InTx, ob.logger.Compliance,
			false, ob.chain.ChainId, event.Raw.TxHash.Hex(), sender, destAddr, "Zeta")
		return nil
	}

	if !destChain.IsZetaChain() {
		paramsDest, found := ob.coreContext.GetEVMChainParams(destChain.ChainId)
		if !found {
			ob.logger.InTx.Warn().Msgf("chain id not present in EVMChainParams  %d", event.DestinationChainId.Int64())
			return nil
		}

		if strings.EqualFold(destAddr, paramsDest.ZetaTokenContractAddress) {
			ob.logger.InTx.Warn().Msgf("potential attack attempt: %s destination address is ZETA token contract address %s", destChain, destAddr)
			return nil
		}
	}
	message := base64.StdEncoding.EncodeToString(event.Message)
	ob.logger.InTx.Info().Msgf("ZetaSent inTx detected on chain %d tx %s block %d from %s value %s message %s",
		ob.chain.ChainId, event.Raw.TxHash.Hex(), event.Raw.BlockNumber, sender, event.ZetaValueAndGas.String(), message)

	return zetacore.GetInBoundVoteMessage(
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
		ob.zetacoreClient.GetKeys().GetOperatorAddress().String(),
		event.Raw.Index,
	)
}

// BuildInboundVoteMsgForTokenSentToTSS builds a inbound vote message for a token sent to TSS
func (ob *Observer) BuildInboundVoteMsgForTokenSentToTSS(tx *ethrpc.Transaction, sender ethcommon.Address, blockNumber uint64) *types.MsgVoteOnObservedInboundTx {
	message := tx.Input

	// compliance check
	maybeReceiver := ""
	parsedAddress, _, err := chains.ParseAddressAndData(message)
	if err == nil && parsedAddress != (ethcommon.Address{}) {
		maybeReceiver = parsedAddress.Hex()
	}
	if config.ContainRestrictedAddress(sender.Hex(), maybeReceiver) {
		compliance.PrintComplianceLog(ob.logger.InTx, ob.logger.Compliance,
			false, ob.chain.ChainId, tx.Hash, sender.Hex(), sender.Hex(), "Gas")
		return nil
	}

	// donation check
	// #nosec G703 err is already checked
	data, _ := hex.DecodeString(message)
	if bytes.Equal(data, []byte(constant.DonationMessage)) {
		ob.logger.InTx.Info().Msgf("thank you rich folk for your donation! tx %s chain %d", tx.Hash, ob.chain.ChainId)
		return nil
	}
	ob.logger.InTx.Info().Msgf("TSS inTx detected on chain %d tx %s block %d from %s value %s message %s",
		ob.chain.ChainId, tx.Hash, blockNumber, sender.Hex(), tx.Value.String(), message)

	return zetacore.GetInBoundVoteMessage(
		sender.Hex(),
		ob.chain.ChainId,
		sender.Hex(),
		sender.Hex(),
		ob.zetacoreClient.Chain().ChainId,
		sdkmath.NewUintFromBigInt(&tx.Value),
		message,
		tx.Hash,
		blockNumber,
		90_000,
		coin.CoinType_Gas,
		"",
		ob.zetacoreClient.GetKeys().GetOperatorAddress().String(),
		0, // not a smart contract call
	)
}

// ObserveTSSReceiveInBlock queries the incoming gas asset to TSS address in a single block and posts votes
func (ob *Observer) ObserveTSSReceiveInBlock(blockNumber uint64) error {
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

// calcBlockRangeToScan calculates the next range of blocks to scan
func (ob *Observer) calcBlockRangeToScan(latestConfirmed, lastScanned, batchSize uint64) (uint64, uint64) {
	startBlock := lastScanned + 1
	toBlock := lastScanned + batchSize
	if toBlock > latestConfirmed {
		toBlock = latestConfirmed
	}
	return startBlock, toBlock
}
