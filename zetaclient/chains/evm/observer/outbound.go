package observer

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"cosmossdk.io/math"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/protocol-contracts/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayevm.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/zetaconnector.non-eth.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/zetaconnectornative.sol"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/evm/common"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/compliance"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

// ProcessOutboundTrackers processes outbound trackers
func (ob *Observer) ProcessOutboundTrackers(ctx context.Context) error {
	trackers, err := ob.ZetacoreClient().GetAllOutboundTrackerByChain(ctx, ob.Chain().ChainId, interfaces.Ascending)
	if err != nil {
		return errors.Wrap(err, "GetAllOutboundTrackerByChain error")
	}

	// keep last block up-to-date
	if err := ob.updateLastBlock(ctx); err != nil {
		return err
	}

	// prepare logger fields
	logger := ob.Logger().Outbound.With().
		Str(logs.FieldMethod, "ProcessOutboundTrackers").
		Logger()

	// process outbound trackers
	for _, tracker := range trackers {
		// go to next tracker if this one already has a confirmed tx
		nonce := tracker.Nonce
		if ob.isTxConfirmed(nonce) {
			continue
		}

		// check each txHash and save tx and receipt if it's legit and confirmed
		txCount := 0
		var outboundReceipt *ethtypes.Receipt
		var outbound *ethtypes.Transaction
		for _, txHash := range tracker.HashList {
			if receipt, tx, ok := ob.checkConfirmedTx(ctx, txHash.TxHash, nonce); ok {
				txCount++
				outboundReceipt = receipt
				outbound = tx

				logger.Info().
					Uint64(logs.FieldNonce, nonce).
					Str(logs.FieldTx, txHash.TxHash).
					Msg("Confirmed outbound")
			}
		}

		switch {
		case txCount == 1:
			ob.setTxNReceipt(nonce, outboundReceipt, outbound)
		case txCount > 1:
			// Unexpected state: multiple transactions exist for a single nonce.
			// This could indicate duplicate transaction broadcasting or unreliable RPC data
			logger.Error().Uint64(logs.FieldNonce, nonce).Msgf("Confirmed multiple (%d) outbound", txCount)
		case tracker.MaxReached():
			logger.Error().Uint64(logs.FieldNonce, nonce).Msg("Outbound tracker is full of hashes")
		}
	}

	return nil
}

// postVoteOutbound posts vote to zetacore for the confirmed outbound
func (ob *Observer) postVoteOutbound(
	ctx context.Context,
	cctxIndex string,
	receipt *ethtypes.Receipt,
	transaction *ethtypes.Transaction,
	receiveValue *big.Int,
	receiveStatus chains.ReceiveStatus,
	nonce uint64,
	coinType coin.CoinType,
	logger zerolog.Logger,
) {
	chainID := ob.Chain().ChainId

	signerAddress := ob.ZetacoreClient().GetKeys().GetOperatorAddress()

	msg := crosschaintypes.NewMsgVoteOutbound(
		signerAddress.String(),
		cctxIndex,
		receipt.TxHash.Hex(),
		receipt.BlockNumber.Uint64(),
		receipt.GasUsed,
		math.NewIntFromBigInt(transaction.GasPrice()),
		transaction.Gas(),
		math.NewUintFromBigInt(receiveValue),
		receiveStatus,
		chainID,
		nonce,
		coinType,
		crosschaintypes.ConfirmationMode_SAFE,
	)

	const gasLimit = zetacore.PostVoteOutboundGasLimit

	retryGasLimit := zetacore.PostVoteOutboundRetryGasLimit
	if msg.Status == chains.ReceiveStatus_failed {
		retryGasLimit = zetacore.PostVoteOutboundRevertGasLimit
	}

	// post vote to zetacore
	logFields := map[string]any{
		logs.FieldNonce: nonce,
		logs.FieldTx:    receipt.TxHash.String(),
	}

	zetaTxHash, ballot, err := ob.ZetacoreClient().PostVoteOutbound(ctx, gasLimit, retryGasLimit, msg)
	if err != nil {
		logger.Error().
			Err(err).
			Fields(logFields).
			Msg("Unable to post outbound vote")
		return
	}

	// print vote tx hash and ballot
	if zetaTxHash != "" {
		logFields[logs.FieldZetaTx] = zetaTxHash
		logFields[logs.FieldBallot] = ballot
		logger.Info().Fields(logFields).Msg("Outbound vote posted")
	}
}

// VoteOutboundIfConfirmed checks outbound status and returns (continueKeysign, error)
func (ob *Observer) VoteOutboundIfConfirmed(
	ctx context.Context,
	cctx *crosschaintypes.CrossChainTx,
) (bool, error) {
	// skip if outbound is not confirmed
	nonce := cctx.GetCurrentOutboundParam().TssNonce
	if !ob.isTxConfirmed(nonce) {
		return true, nil
	}
	receipt, transaction := ob.getTxNReceipt(nonce)
	sendID := fmt.Sprintf("%d-%d", ob.Chain().ChainId, nonce)
	logger := ob.Logger().Outbound.With().Str("sendID", sendID).Logger()
	// get connector and erc20Custody contracts
	// Only one of these connector contracts will be used at one time.
	// V1 cctx's of cointype ZETA would not be processed once the connector is upgraded to V2
	connectorLegacyAddr, connectorLegacy, err := ob.getConnectorLegacyContract()
	if err != nil {
		return true, errors.Wrap(err, "error getting legacy zeta connector")
	}

	connectorAddr, connector, err := ob.getConnectorContract()
	if err != nil {
		return true, errors.Wrap(err, "error getting zeta connector")
	}

	custodyAddr, custody, err := ob.getERC20CustodyContract()
	if err != nil {
		return true, errors.Wrap(err, "error getting erc20 custody")
	}
	gatewayAddr, gateway, err := ob.getGatewayContract()
	if err != nil {
		return true, errors.Wrap(err, "error getting gateway for chain")
	}
	_, custodyV2, err := ob.getERC20CustodyV2Contract()
	if err != nil {
		return true, errors.Wrap(err, "error getting erc20 custody v2 for chain")
	}

	// define a few common variables
	var (
		receiveValue  *big.Int
		receiveStatus chains.ReceiveStatus
		cointype      = cctx.InboundParams.CoinType
	)

	// cancelled transaction means the outbound is failed
	// - set amount to CCTX's amount to bypass amount check in zetacore
	// - set status to failed to revert the CCTX in zetacore
	if compliance.IsCCTXRestricted(cctx) {
		receiveValue = cctx.GetCurrentOutboundParam().Amount.BigInt()
		receiveStatus = chains.ReceiveStatus_failed
		ob.postVoteOutbound(ctx, cctx.Index, receipt, transaction, receiveValue, receiveStatus, nonce, cointype, logger)
		return false, nil
	}

	// parse the received value from the outbound receipt
	receiveValue, receiveStatus, err = parseOutboundReceivedValue(
		cctx,
		receipt,
		transaction,
		cointype,
		connectorLegacyAddr,
		connectorLegacy,
		custodyAddr,
		custody,
		custodyV2,
		gatewayAddr,
		gateway,
		connectorAddr,
		connector,
	)
	if err != nil {
		logger.Error().
			Err(err).
			Msgf("VoteOutboundIfConfirmed: error parsing outbound event for chain %d txhash %s", ob.Chain().ChainId, receipt.TxHash)
		return true, err
	}

	// post vote to zetacore
	ob.postVoteOutbound(ctx, cctx.Index, receipt, transaction, receiveValue, receiveStatus, nonce, cointype, logger)
	return false, nil
}

// parseOutboundReceivedValue parses the received value and status from the outbound receipt
// The received value is the amount of Zeta/ERC20/Gas token (released from connector/custody/TSS) sent to the receiver
// TODO: simplify this function and reduce the number of argument
// https://github.com/zeta-chain/node/issues/2627
// https://github.com/zeta-chain/node/pull/2666#discussion_r1718379784
func parseOutboundReceivedValue(
	cctx *crosschaintypes.CrossChainTx,
	receipt *ethtypes.Receipt,
	transaction *ethtypes.Transaction,
	cointype coin.CoinType,
	connectorAddress ethcommon.Address,
	connector *zetaconnector.ZetaConnectorNonEth,
	custodyAddress ethcommon.Address,
	custody *erc20custody.ERC20Custody,
	custodyV2 *erc20custody.ERC20Custody,
	gatewayAddress ethcommon.Address,
	gateway *gatewayevm.GatewayEVM,
	connectorNativeAddress ethcommon.Address,
	connectorNative *zetaconnectornative.ZetaConnectorNative,
) (*big.Int, chains.ReceiveStatus, error) {
	// determine the receive status and value
	// https://docs.nethereum.com/en/latest/nethereum-receipt-status/
	receiveValue := big.NewInt(0)
	receiveStatus := chains.ReceiveStatus_failed
	if receipt.Status == ethtypes.ReceiptStatusSuccessful {
		receiveValue = transaction.Value()
		receiveStatus = chains.ReceiveStatus_success
	}

	// parse outbound event for protocol contract v2
	if cctx.ProtocolContractVersion == crosschaintypes.ProtocolContractVersion_V2 {
		return parseOutboundEventV2(
			cctx,
			receipt,
			transaction,
			custodyAddress,
			custodyV2,
			gatewayAddress,
			gateway,
			connectorNativeAddress,
			connectorNative,
		)
	}

	// parse receive value from the outbound receipt for Zeta and ERC20
	switch cointype {
	case coin.CoinType_Zeta:
		if receipt.Status == ethtypes.ReceiptStatusSuccessful {
			receivedLog, revertedLog, err := parseAndCheckZetaEvent(cctx, receipt, connectorAddress, connector)
			if err != nil {
				return nil, chains.ReceiveStatus_failed, err
			}
			// use the value in ZetaReceived/ZetaReverted event for vote message
			if receivedLog != nil {
				receiveValue = receivedLog.ZetaValue
			} else if revertedLog != nil {
				receiveValue = revertedLog.RemainingZetaValue
			}
		}
	case coin.CoinType_ERC20:
		if receipt.Status == ethtypes.ReceiptStatusSuccessful {
			withdrawn, err := parseAndCheckWithdrawnEvent(cctx, receipt, custodyAddress, custody)
			if err != nil {
				return nil, chains.ReceiveStatus_failed, err
			}
			// use the value in Withdrawn event for vote message
			receiveValue = withdrawn.Amount
		}
	case coin.CoinType_Gas, coin.CoinType_Cmd:
		// nothing to do for CoinType_Gas/CoinType_Cmd, no need to parse event
	default:
		return nil, chains.ReceiveStatus_failed, fmt.Errorf("unknown coin type %s", cointype)
	}
	return receiveValue, receiveStatus, nil
}

// parseAndCheckZetaEvent parses and checks ZetaReceived/ZetaReverted event from the outbound receipt
// It either returns an ZetaReceived or an ZetaReverted event, or an error if no event found
func parseAndCheckZetaEvent(
	cctx *crosschaintypes.CrossChainTx,
	receipt *ethtypes.Receipt,
	connectorAddr ethcommon.Address,
	connector *zetaconnector.ZetaConnectorNonEth,
) (*zetaconnector.ZetaConnectorNonEthZetaReceived, *zetaconnector.ZetaConnectorNonEthZetaReverted, error) {
	params := cctx.GetCurrentOutboundParam()
	for _, vLog := range receipt.Logs {
		// try parsing ZetaReceived event
		received, err := connector.ZetaConnectorNonEthFilterer.ParseZetaReceived(*vLog)
		if err == nil {
			err = common.ValidateEvmTxLog(vLog, connectorAddr, receipt.TxHash.Hex(), common.TopicsZetaReceived)
			if err != nil {
				return nil, nil, errors.Wrap(err, "error validating ZetaReceived event")
			}
			if !strings.EqualFold(received.DestinationAddress.Hex(), params.Receiver) {
				return nil, nil, fmt.Errorf("receiver address mismatch in ZetaReceived event, want %s got %s",
					params.Receiver, received.DestinationAddress.Hex())
			}
			if received.ZetaValue.Cmp(params.Amount.BigInt()) != 0 {
				return nil, nil, fmt.Errorf("amount mismatch in ZetaReceived event, want %s got %s",
					params.Amount.String(), received.ZetaValue.String())
			}
			if ethcommon.BytesToHash(received.InternalSendHash[:]).Hex() != cctx.Index {
				return nil, nil, fmt.Errorf("cctx index mismatch in ZetaReceived event, want %s got %s",
					cctx.Index, hex.EncodeToString(received.InternalSendHash[:]))
			}
			return received, nil, nil
		}
		// try parsing ZetaReverted event
		reverted, err := connector.ZetaConnectorNonEthFilterer.ParseZetaReverted(*vLog)
		if err == nil {
			err = common.ValidateEvmTxLog(vLog, connectorAddr, receipt.TxHash.Hex(), common.TopicsZetaReverted)
			if err != nil {
				return nil, nil, errors.Wrap(err, "error validating ZetaReverted event")
			}
			if !strings.EqualFold(
				ethcommon.BytesToAddress(reverted.DestinationAddress[:]).Hex(),
				cctx.InboundParams.Sender,
			) {
				return nil, nil, fmt.Errorf("receiver address mismatch in ZetaReverted event, want %s got %s",
					cctx.InboundParams.Sender, ethcommon.BytesToAddress(reverted.DestinationAddress[:]).Hex())
			}
			if reverted.RemainingZetaValue.Cmp(params.Amount.BigInt()) != 0 {
				return nil, nil, fmt.Errorf("amount mismatch in ZetaReverted event, want %s got %s",
					params.Amount.String(), reverted.RemainingZetaValue.String())
			}
			if ethcommon.BytesToHash(reverted.InternalSendHash[:]).Hex() != cctx.Index {
				return nil, nil, fmt.Errorf("cctx index mismatch in ZetaReverted event, want %s got %s",
					cctx.Index, hex.EncodeToString(reverted.InternalSendHash[:]))
			}
			return nil, reverted, nil
		}
	}
	return nil, nil, errors.New("no ZetaReceived/ZetaReverted event found")
}

// parseAndCheckWithdrawnEvent parses and checks erc20 Withdrawn event from the outbound receipt
func parseAndCheckWithdrawnEvent(
	cctx *crosschaintypes.CrossChainTx,
	receipt *ethtypes.Receipt,
	custodyAddr ethcommon.Address,
	custody *erc20custody.ERC20Custody,
) (*erc20custody.ERC20CustodyWithdrawn, error) {
	params := cctx.GetCurrentOutboundParam()
	for _, vLog := range receipt.Logs {
		withdrawn, err := custody.ParseWithdrawn(*vLog)
		if err == nil {
			err = common.ValidateEvmTxLog(vLog, custodyAddr, receipt.TxHash.Hex(), common.TopicsWithdrawn)
			if err != nil {
				return nil, errors.Wrap(err, "error validating Withdrawn event")
			}
			if !strings.EqualFold(withdrawn.To.Hex(), params.Receiver) {
				return nil, fmt.Errorf("receiver address mismatch in Withdrawn event, want %s got %s",
					params.Receiver, withdrawn.To.Hex())
			}
			if !strings.EqualFold(withdrawn.Token.Hex(), cctx.InboundParams.Asset) {
				return nil, fmt.Errorf("asset mismatch in Withdrawn event, want %s got %s",
					cctx.InboundParams.Asset, withdrawn.Token.Hex())
			}
			if withdrawn.Amount.Cmp(params.Amount.BigInt()) != 0 {
				return nil, fmt.Errorf("amount mismatch in Withdrawn event, want %s got %s",
					params.Amount.String(), withdrawn.Amount.String())
			}
			return withdrawn, nil
		}
	}
	return nil, errors.New("no ERC20 Withdrawn event found")
}

// filterTSSOutbound filters the outbounds from TSS address to supplement outbound trackers
func (ob *Observer) filterTSSOutbound(ctx context.Context, startBlock, toBlock uint64) {
	// filters the outbounds from TSS address block by block
	for bn := startBlock; bn <= toBlock; bn++ {
		ob.filterTSSOutboundInBlock(ctx, bn)
	}
}

// filterTSSOutboundInBlock filters the outbounds in a single block to supplement outbound trackers
func (ob *Observer) filterTSSOutboundInBlock(ctx context.Context, blockNumber uint64) {
	// query block and ignore error (we don't rescan as we are only supplementing outbound trackers)
	block, err := ob.GetBlockByNumberCached(ctx, blockNumber)
	if err != nil {
		ob.Logger().
			Outbound.Error().
			Err(err).
			Uint64(logs.FieldBlock, blockNumber).
			Msg("Error getting block")
		return
	}

	for i := range block.Transactions {
		tx := block.Transactions[i]

		// noop
		if ethcommon.HexToAddress(tx.From) != ob.TSS().PubKey().AddressEVM() {
			continue
		}

		// #nosec G115 nonce always positive
		nonce := uint64(tx.Nonce)

		// noop
		if ob.isTxConfirmed(nonce) {
			continue
		}

		if receipt, txx, ok := ob.checkConfirmedTx(ctx, tx.Hash, nonce); ok {
			ob.setTxNReceipt(nonce, receipt, txx)
			ob.Logger().
				Outbound.Info().
				Uint64(logs.FieldNonce, nonce).
				Str(logs.FieldTx, tx.Hash).
				Msg("TSS outbound detected")
		}
	}
}

// checkConfirmedTx checks if a txHash is confirmed
// returns (receipt, transaction, true) if confirmed or (nil, nil, false) otherwise
func (ob *Observer) checkConfirmedTx(
	ctx context.Context,
	txHash string,
	nonce uint64,
) (*ethtypes.Receipt, *ethtypes.Transaction, bool) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// prepare logger
	logger := ob.Logger().Outbound.With().
		Str(logs.FieldMethod, "checkConfirmedTx").
		Int64(logs.FieldChain, ob.Chain().ChainId).
		Uint64(logs.FieldNonce, nonce).
		Str(logs.FieldTx, txHash).
		Logger()

	// query transaction
	transaction, isPending, err := ob.evmClient.TransactionByHash(ctx, ethcommon.HexToHash(txHash))
	if err != nil {
		logger.Error().Err(err).Msg("TransactionByHash error")
		return nil, nil, false
	}
	if transaction == nil { // should not happen
		logger.Error().Msg("transaction is nil")
		return nil, nil, false
	}
	if isPending {
		// should not happen when we are here. The outbound tracker reporter won't report a pending tx.
		logger.Error().Msg("transaction is pending")
		return nil, nil, false
	}

	// check tx sender and nonce
	signer := ethtypes.NewLondonSigner(big.NewInt(ob.Chain().ChainId))
	from, err := signer.Sender(transaction)
	switch {
	case err != nil:
		logger.Error().Err(err).Msg("local recovery of sender address failed")
		return nil, nil, false
	case from != ob.TSS().PubKey().AddressEVM():
		// might be false positive during TSS upgrade for unconfirmed txs
		// Make sure all deposits/withdrawals are paused during TSS upgrade
		logger.Error().Str("tx.sender", from.String()).Msgf("tx sender is not TSS addresses")
		return nil, nil, false
	case transaction.Nonce() != nonce:
		logger.Error().
			Uint64("tx.nonce", transaction.Nonce()).
			Uint64("tracker.nonce", nonce).
			Msg("tx nonce is not matching tracker nonce")
		return nil, nil, false
	}

	// query receipt
	receipt, err := ob.evmClient.TransactionReceipt(ctx, ethcommon.HexToHash(txHash))
	if err != nil {
		logger.Error().Err(err).Msg("TransactionReceipt error")
		return nil, nil, false
	}
	if receipt == nil { // should not happen
		logger.Error().Msg("receipt is nil")
		return nil, nil, false
	}

	// check confirmations
	txBlock := receipt.BlockNumber.Uint64()
	if !ob.IsBlockConfirmedForOutboundSafe(txBlock) {
		logger.Debug().Uint64("tx_block", txBlock).Uint64("last_block", ob.LastBlock()).Msg("tx not confirmed yet")
		return nil, nil, false
	}

	// cross-check tx inclusion against the block
	// Note: a guard for false BlockNumber in receipt. The blob-carrying tx won't come here
	err = ob.checkTxInclusion(ctx, transaction, receipt)
	if err != nil {
		logger.Error().Err(err).Msg("CheckTxInclusion error")
		return nil, nil, false
	}

	return receipt, transaction, true
}
