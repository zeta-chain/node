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

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/evm/common"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/compliance"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

// ObserverOutbound processes outbound trackers
func (ob *Observer) ObserverOutbound(ctx context.Context) error {
	chainID := ob.Chain().ChainId
	trackers, err := ob.ZetacoreClient().GetAllOutboundTrackerByChain(ctx, ob.Chain().ChainId, interfaces.Ascending)
	if err != nil {
		return errors.Wrap(err, "GetAllOutboundTrackerByChain error")
	}

	// prepare logger fields
	logger := ob.Logger().Outbound.With().
		Str(logs.FieldMethod, "ProcessOutboundTrackers").
		Int64(logs.FieldChain, chainID).
		Logger()

	// process outbound trackers
	for _, tracker := range trackers {
		// go to next tracker if this one already has a confirmed tx
		nonce := tracker.Nonce
		if ob.IsTxConfirmed(nonce) {
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
				logger.Info().Msgf("confirmed outbound %s for chain %d nonce %d", txHash.TxHash, chainID, nonce)
				if txCount > 1 {
					logger.Error().
						Msgf("checkConfirmedTx passed, txCount %d chain %d nonce %d receipt %v tx %v", txCount, chainID, nonce, receipt, tx)
				}
			}
		}

		// should be only one txHash confirmed for each nonce.
		if txCount == 1 {
			ob.SetTxNReceipt(nonce, outboundReceipt, outbound)
		} else if txCount > 1 {
			// should not happen. We can't tell which txHash is true. It might happen (e.g. bug, glitchy/hacked endpoint)
			ob.Logger().Outbound.Error().Msgf("WatchOutbound: confirmed multiple (%d) outbound for chain %d nonce %d", txCount, chainID, nonce)
		} else {
			if tracker.MaxReached() {
				ob.Logger().Outbound.Error().Msgf("WatchOutbound: outbound tracker is full of hashes for chain %d nonce %d", chainID, nonce)
			}
		}
	}

	return nil
}

// PostVoteOutbound posts vote to zetacore for the confirmed outbound
func (ob *Observer) PostVoteOutbound(
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
	)

	const gasLimit = zetacore.PostVoteOutboundGasLimit

	var retryGasLimit uint64
	if msg.Status == chains.ReceiveStatus_failed {
		retryGasLimit = zetacore.PostVoteOutboundRevertGasLimit
	}

	// post vote to zetacore
	logFields := map[string]any{
		"chain":    chainID,
		"nonce":    nonce,
		"outbound": receipt.TxHash.String(),
	}
	zetaTxHash, ballot, err := ob.ZetacoreClient().PostVoteOutbound(ctx, gasLimit, retryGasLimit, msg)
	if err != nil {
		logger.Error().
			Err(err).
			Fields(logFields).
			Msgf("PostVoteOutbound: error posting vote for chain %d", chainID)
		return
	}

	// print vote tx hash and ballot
	if zetaTxHash != "" {
		logFields["vote"] = zetaTxHash
		logFields["ballot"] = ballot
		logger.Info().Fields(logFields).Msgf("PostVoteOutbound: posted vote for chain %d", chainID)
	}
}

// VoteOutboundIfConfirmed checks outbound status and returns (continueKeysign, error)
func (ob *Observer) VoteOutboundIfConfirmed(
	ctx context.Context,
	cctx *crosschaintypes.CrossChainTx,
) (bool, error) {
	// skip if outbound is not confirmed
	nonce := cctx.GetCurrentOutboundParam().TssNonce
	if !ob.IsTxConfirmed(nonce) {
		return true, nil
	}
	receipt, transaction := ob.GetTxNReceipt(nonce)
	sendID := fmt.Sprintf("%d-%d", ob.Chain().ChainId, nonce)
	logger := ob.Logger().Outbound.With().Str("sendID", sendID).Logger()

	// get connector and erce20Custody contracts
	connectorAddr, connector, err := ob.GetConnectorContract()
	if err != nil {
		return true, errors.Wrapf(err, "error getting zeta connector for chain %d", ob.Chain().ChainId)
	}
	custodyAddr, custody, err := ob.GetERC20CustodyContract()
	if err != nil {
		return true, errors.Wrapf(err, "error getting erc20 custody for chain %d", ob.Chain().ChainId)
	}
	gatewayAddr, gateway, err := ob.GetGatewayContract()
	if err != nil {
		return true, errors.Wrap(err, "error getting gateway for chain")
	}
	_, custodyV2, err := ob.GetERC20CustodyV2Contract()
	if err != nil {
		return true, errors.Wrapf(err, "error getting erc20 custody v2 for chain %d", ob.Chain().ChainId)
	}

	// define a few common variables
	var receiveValue *big.Int
	var receiveStatus chains.ReceiveStatus
	cointype := cctx.InboundParams.CoinType

	// compliance check, special handling the cancelled cctx
	if compliance.IsCctxRestricted(cctx) {
		// use cctx's amount to bypass the amount check in zetacore
		receiveValue = cctx.GetCurrentOutboundParam().Amount.BigInt()
		receiveStatus := chains.ReceiveStatus_failed
		if receipt.Status == ethtypes.ReceiptStatusSuccessful {
			receiveStatus = chains.ReceiveStatus_success
		}
		ob.PostVoteOutbound(ctx, cctx.Index, receipt, transaction, receiveValue, receiveStatus, nonce, cointype, logger)
		return false, nil
	}

	// parse the received value from the outbound receipt
	receiveValue, receiveStatus, err = parseOutboundReceivedValue(
		cctx,
		receipt,
		transaction,
		cointype,
		connectorAddr,
		connector,
		custodyAddr,
		custody,
		custodyV2,
		gatewayAddr,
		gateway,
	)
	if err != nil {
		logger.Error().
			Err(err).
			Msgf("VoteOutboundIfConfirmed: error parsing outbound event for chain %d txhash %s", ob.Chain().ChainId, receipt.TxHash)
		return true, err
	}

	// post vote to zetacore
	ob.PostVoteOutbound(ctx, cctx.Index, receipt, transaction, receiveValue, receiveStatus, nonce, cointype, logger)
	return false, nil
}

// parseOutboundReceivedValue parses the received value and status from the outbound receipt
// The receivd value is the amount of Zeta/ERC20/Gas token (released from connector/custody/TSS) sent to the receiver
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
		return parseOutboundEventV2(cctx, receipt, transaction, custodyAddress, custodyV2, gatewayAddress, gateway)
	}

	// parse receive value from the outbound receipt for Zeta and ERC20
	switch cointype {
	case coin.CoinType_Zeta:
		if receipt.Status == ethtypes.ReceiptStatusSuccessful {
			receivedLog, revertedLog, err := ParseAndCheckZetaEvent(cctx, receipt, connectorAddress, connector)
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
			withdrawn, err := ParseAndCheckWithdrawnEvent(cctx, receipt, custodyAddress, custody)
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

// ParseAndCheckZetaEvent parses and checks ZetaReceived/ZetaReverted event from the outbound receipt
// It either returns an ZetaReceived or an ZetaReverted event, or an error if no event found
func ParseAndCheckZetaEvent(
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

// ParseAndCheckWithdrawnEvent parses and checks erc20 Withdrawn event from the outbound receipt
func ParseAndCheckWithdrawnEvent(
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

// FilterTSSOutbound filters the outbounds from TSS address to supplement outbound trackers
func (ob *Observer) FilterTSSOutbound(ctx context.Context, startBlock, toBlock uint64) {
	// filters the outbounds from TSS address block by block
	for bn := startBlock; bn <= toBlock; bn++ {
		ob.FilterTSSOutboundInBlock(ctx, bn)
	}
}

// FilterTSSOutboundInBlock filters the outbounds in a single block to supplement outbound trackers
func (ob *Observer) FilterTSSOutboundInBlock(ctx context.Context, blockNumber uint64) {
	// query block and ignore error (we don't rescan as we are only supplementing outbound trackers)
	block, err := ob.GetBlockByNumberCached(blockNumber)
	if err != nil {
		ob.Logger().
			Outbound.Error().
			Err(err).
			Msgf("error getting block %d for chain %d", blockNumber, ob.Chain().ChainId)
		return
	}

	for i := range block.Transactions {
		tx := block.Transactions[i]
		if ethcommon.HexToAddress(tx.From) == ob.TSS().PubKey().AddressEVM() {
			// #nosec G115 nonce always positive
			nonce := uint64(tx.Nonce)
			if !ob.IsTxConfirmed(nonce) {
				if receipt, txx, ok := ob.checkConfirmedTx(ctx, tx.Hash, nonce); ok {
					ob.SetTxNReceipt(nonce, receipt, txx)
					ob.Logger().
						Outbound.Info().
						Msgf("TSS outbound detected on chain %d nonce %d tx %s", ob.Chain().ChainId, nonce, tx.Hash)
				}
			}
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
	ob.LastBlock()
	// check confirmations
	lastHeight, err := ob.evmClient.BlockNumber(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("BlockNumber error")
		return nil, nil, false
	}
	if !ob.HasEnoughConfirmations(receipt, lastHeight) {
		logger.Debug().
			Msgf("tx included but not confirmed, receipt block %d current block %d", receipt.BlockNumber.Uint64(), lastHeight)
		return nil, nil, false
	}

	// cross-check tx inclusion against the block
	// Note: a guard for false BlockNumber in receipt. The blob-carrying tx won't come here
	err = ob.CheckTxInclusion(transaction, receipt)
	if err != nil {
		logger.Error().Err(err).Msg("CheckTxInclusion error")
		return nil, nil, false
	}

	return receipt, transaction, true
}
