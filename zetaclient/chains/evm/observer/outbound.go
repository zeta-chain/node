package observer

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.non-eth.sol"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/evm"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/compliance"
	clientcontext "github.com/zeta-chain/zetacore/zetaclient/context"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
)

// GetTxID returns a unique id for outbound tx
func (ob *Observer) GetTxID(nonce uint64) string {
	tssAddr := ob.Tss.EVMAddress().String()
	return fmt.Sprintf("%d-%s-%d", ob.chain.ChainId, tssAddr, nonce)
}

// WatchOutTx watches evm chain for outgoing txs status
func (ob *Observer) WatchOutTx() {
	ticker, err := clienttypes.NewDynamicTicker(fmt.Sprintf("EVM_WatchOutTx_%d", ob.chain.ChainId), ob.GetChainParams().OutTxTicker)
	if err != nil {
		ob.logger.OutTx.Error().Err(err).Msg("error creating ticker")
		return
	}

	ob.logger.OutTx.Info().Msgf("WatchOutTx started for chain %d", ob.chain.ChainId)
	sampledLogger := ob.logger.OutTx.Sample(&zerolog.BasicSampler{N: 10})
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			if !clientcontext.IsOutboundObservationEnabled(ob.coreContext, ob.GetChainParams()) {
				sampledLogger.Info().Msgf("WatchOutTx: outbound observation is disabled for chain %d", ob.chain.ChainId)
				continue
			}
			trackers, err := ob.zetacoreClient.GetAllOutTxTrackerByChain(ob.chain.ChainId, interfaces.Ascending)
			if err != nil {
				continue
			}
			for _, tracker := range trackers {
				nonceInt := tracker.Nonce
				if ob.IsTxConfirmed(nonceInt) { // Go to next tracker if this one already has a confirmed tx
					continue
				}
				txCount := 0
				var outtxReceipt *ethtypes.Receipt
				var outtx *ethtypes.Transaction
				for _, txHash := range tracker.HashList {
					if receipt, tx, ok := ob.checkConfirmedTx(txHash.TxHash, nonceInt); ok {
						txCount++
						outtxReceipt = receipt
						outtx = tx
						ob.logger.OutTx.Info().Msgf("WatchOutTx: confirmed outTx %s for chain %d nonce %d", txHash.TxHash, ob.chain.ChainId, nonceInt)
						if txCount > 1 {
							ob.logger.OutTx.Error().Msgf(
								"WatchOutTx: checkConfirmedTx passed, txCount %d chain %d nonce %d receipt %v transaction %v", txCount, ob.chain.ChainId, nonceInt, outtxReceipt, outtx)
						}
					}
				}
				if txCount == 1 { // should be only one txHash confirmed for each nonce.
					ob.SetTxNReceipt(nonceInt, outtxReceipt, outtx)
				} else if txCount > 1 { // should not happen. We can't tell which txHash is true. It might happen (e.g. glitchy/hacked endpoint)
					ob.logger.OutTx.Error().Msgf("WatchOutTx: confirmed multiple (%d) outTx for chain %d nonce %d", txCount, ob.chain.ChainId, nonceInt)
				}
			}
			ticker.UpdateInterval(ob.GetChainParams().OutTxTicker, ob.logger.OutTx)
		case <-ob.stop:
			ob.logger.OutTx.Info().Msg("WatchOutTx: stopped")
			return
		}
	}
}

// PostVoteOutbound posts vote to zetacore for the confirmed outtx
func (ob *Observer) PostVoteOutbound(
	cctxIndex string,
	receipt *ethtypes.Receipt,
	transaction *ethtypes.Transaction,
	receiveValue *big.Int,
	receiveStatus chains.ReceiveStatus,
	nonce uint64,
	cointype coin.CoinType,
	logger zerolog.Logger,
) {
	chainID := ob.chain.ChainId
	zetaTxHash, ballot, err := ob.zetacoreClient.PostVoteOutbound(
		cctxIndex,
		receipt.TxHash.Hex(),
		receipt.BlockNumber.Uint64(),
		receipt.GasUsed,
		transaction.GasPrice(),
		transaction.Gas(),
		receiveValue,
		receiveStatus,
		ob.chain,
		nonce,
		cointype,
	)
	if err != nil {
		logger.Error().Err(err).Msgf("PostVoteOutbound: error posting vote for chain %d nonce %d outtx %s ", chainID, nonce, receipt.TxHash)
	} else if zetaTxHash != "" {
		logger.Info().Msgf("PostVoteOutbound: posted vote for chain %d nonce %d outtx %s vote %s ballot %s", chainID, nonce, receipt.TxHash, zetaTxHash, ballot)
	}
}

// IsOutboundProcessed checks outtx status and returns (isIncluded, isConfirmed, error)
// It also posts vote to zetacore if the tx is confirmed
func (ob *Observer) IsOutboundProcessed(cctx *crosschaintypes.CrossChainTx, logger zerolog.Logger) (bool, bool, error) {
	// skip if outtx is not confirmed
	nonce := cctx.GetCurrentOutTxParam().OutboundTxTssNonce
	if !ob.IsTxConfirmed(nonce) {
		return false, false, nil
	}
	receipt, transaction := ob.GetTxNReceipt(nonce)
	sendID := fmt.Sprintf("%d-%d", ob.chain.ChainId, nonce)
	logger = logger.With().Str("sendID", sendID).Logger()

	// get connector and erce20Custody contracts
	connectorAddr, connector, err := ob.GetConnectorContract()
	if err != nil {
		return false, false, errors.Wrapf(err, "error getting zeta connector for chain %d", ob.chain.ChainId)
	}
	custodyAddr, custody, err := ob.GetERC20CustodyContract()
	if err != nil {
		return false, false, errors.Wrapf(err, "error getting erc20 custody for chain %d", ob.chain.ChainId)
	}

	// define a few common variables
	var receiveValue *big.Int
	var receiveStatus chains.ReceiveStatus
	cointype := cctx.InboundTxParams.CoinType

	// compliance check, special handling the cancelled cctx
	if compliance.IsCctxRestricted(cctx) {
		// use cctx's amount to bypass the amount check in zetacore
		receiveValue = cctx.GetCurrentOutTxParam().Amount.BigInt()
		receiveStatus := chains.ReceiveStatus_failed
		if receipt.Status == ethtypes.ReceiptStatusSuccessful {
			receiveStatus = chains.ReceiveStatus_success
		}
		ob.PostVoteOutbound(cctx.Index, receipt, transaction, receiveValue, receiveStatus, nonce, cointype, logger)
		return true, true, nil
	}

	// parse the received value from the outtx receipt
	receiveValue, receiveStatus, err = ParseOuttxReceivedValue(cctx, receipt, transaction, cointype, connectorAddr, connector, custodyAddr, custody)
	if err != nil {
		logger.Error().Err(err).Msgf("IsOutboundProcessed: error parsing outtx event for chain %d txhash %s", ob.chain.ChainId, receipt.TxHash)
		return false, false, err
	}

	// post vote to zetacore
	ob.PostVoteOutbound(cctx.Index, receipt, transaction, receiveValue, receiveStatus, nonce, cointype, logger)
	return true, true, nil
}

// ParseAndCheckZetaEvent parses and checks ZetaReceived/ZetaReverted event from the outtx receipt
// It either returns an ZetaReceived or an ZetaReverted event, or an error if no event found
func ParseAndCheckZetaEvent(
	cctx *crosschaintypes.CrossChainTx,
	receipt *ethtypes.Receipt,
	connectorAddr ethcommon.Address,
	connector *zetaconnector.ZetaConnectorNonEth,
) (*zetaconnector.ZetaConnectorNonEthZetaReceived, *zetaconnector.ZetaConnectorNonEthZetaReverted, error) {
	params := cctx.GetCurrentOutTxParam()
	for _, vLog := range receipt.Logs {
		// try parsing ZetaReceived event
		received, err := connector.ZetaConnectorNonEthFilterer.ParseZetaReceived(*vLog)
		if err == nil {
			err = evm.ValidateEvmTxLog(vLog, connectorAddr, receipt.TxHash.Hex(), evm.TopicsZetaReceived)
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
			err = evm.ValidateEvmTxLog(vLog, connectorAddr, receipt.TxHash.Hex(), evm.TopicsZetaReverted)
			if err != nil {
				return nil, nil, errors.Wrap(err, "error validating ZetaReverted event")
			}
			if !strings.EqualFold(ethcommon.BytesToAddress(reverted.DestinationAddress[:]).Hex(), cctx.InboundTxParams.Sender) {
				return nil, nil, fmt.Errorf("receiver address mismatch in ZetaReverted event, want %s got %s",
					cctx.InboundTxParams.Sender, ethcommon.BytesToAddress(reverted.DestinationAddress[:]).Hex())
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

// ParseAndCheckWithdrawnEvent parses and checks erc20 Withdrawn event from the outtx receipt
func ParseAndCheckWithdrawnEvent(
	cctx *crosschaintypes.CrossChainTx,
	receipt *ethtypes.Receipt,
	custodyAddr ethcommon.Address,
	custody *erc20custody.ERC20Custody,
) (*erc20custody.ERC20CustodyWithdrawn, error) {
	params := cctx.GetCurrentOutTxParam()
	for _, vLog := range receipt.Logs {
		withdrawn, err := custody.ParseWithdrawn(*vLog)
		if err == nil {
			err = evm.ValidateEvmTxLog(vLog, custodyAddr, receipt.TxHash.Hex(), evm.TopicsWithdrawn)
			if err != nil {
				return nil, errors.Wrap(err, "error validating Withdrawn event")
			}
			if !strings.EqualFold(withdrawn.Recipient.Hex(), params.Receiver) {
				return nil, fmt.Errorf("receiver address mismatch in Withdrawn event, want %s got %s",
					params.Receiver, withdrawn.Recipient.Hex())
			}
			if !strings.EqualFold(withdrawn.Asset.Hex(), cctx.InboundTxParams.Asset) {
				return nil, fmt.Errorf("asset mismatch in Withdrawn event, want %s got %s",
					cctx.InboundTxParams.Asset, withdrawn.Asset.Hex())
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

// ParseOuttxReceivedValue parses the received value and status from the outtx receipt
// The receivd value is the amount of Zeta/ERC20/Gas token (released from connector/custody/TSS) sent to the receiver
func ParseOuttxReceivedValue(
	cctx *crosschaintypes.CrossChainTx,
	receipt *ethtypes.Receipt,
	transaction *ethtypes.Transaction,
	cointype coin.CoinType,
	connectorAddress ethcommon.Address,
	connector *zetaconnector.ZetaConnectorNonEth,
	custodyAddress ethcommon.Address,
	custody *erc20custody.ERC20Custody,
) (*big.Int, chains.ReceiveStatus, error) {
	// determine the receive status and value
	// https://docs.nethereum.com/en/latest/nethereum-receipt-status/
	receiveValue := big.NewInt(0)
	receiveStatus := chains.ReceiveStatus_failed
	if receipt.Status == ethtypes.ReceiptStatusSuccessful {
		receiveValue = transaction.Value()
		receiveStatus = chains.ReceiveStatus_success
	}

	// parse receive value from the outtx receipt for Zeta and ERC20
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

// checkConfirmedTx checks if a txHash is confirmed
// returns (receipt, transaction, true) if confirmed or (nil, nil, false) otherwise
func (ob *Observer) checkConfirmedTx(txHash string, nonce uint64) (*ethtypes.Receipt, *ethtypes.Transaction, bool) {
	ctxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// query transaction
	transaction, isPending, err := ob.evmClient.TransactionByHash(ctxt, ethcommon.HexToHash(txHash))
	if err != nil {
		log.Error().Err(err).Msgf("confirmTxByHash: error getting transaction for outtx %s chain %d", txHash, ob.chain.ChainId)
		return nil, nil, false
	}
	if transaction == nil { // should not happen
		log.Error().Msgf("confirmTxByHash: transaction is nil for txHash %s nonce %d", txHash, nonce)
		return nil, nil, false
	}

	// check tx sender and nonce
	signer := ethtypes.NewLondonSigner(big.NewInt(ob.chain.ChainId))
	from, err := signer.Sender(transaction)
	if err != nil {
		log.Error().Err(err).Msgf("confirmTxByHash: local recovery of sender address failed for outtx %s chain %d", transaction.Hash().Hex(), ob.chain.ChainId)
		return nil, nil, false
	}
	if from != ob.Tss.EVMAddress() { // must be TSS address
		log.Error().Msgf("confirmTxByHash: sender %s for outtx %s chain %d is not TSS address %s",
			from.Hex(), transaction.Hash().Hex(), ob.chain.ChainId, ob.Tss.EVMAddress().Hex())
		return nil, nil, false
	}
	if transaction.Nonce() != nonce { // must match cctx nonce
		log.Error().Msgf("confirmTxByHash: outtx %s nonce mismatch: wanted %d, got tx nonce %d", txHash, nonce, transaction.Nonce())
		return nil, nil, false
	}

	// save pending transaction
	if isPending {
		ob.SetPendingTx(nonce, transaction)
		return nil, nil, false
	}

	// query receipt
	receipt, err := ob.evmClient.TransactionReceipt(ctxt, ethcommon.HexToHash(txHash))
	if err != nil {
		if err != ethereum.NotFound {
			log.Warn().Err(err).Msgf("confirmTxByHash: TransactionReceipt error, txHash %s nonce %d", txHash, nonce)
		}
		return nil, nil, false
	}
	if receipt == nil { // should not happen
		log.Error().Msgf("confirmTxByHash: receipt is nil for txHash %s nonce %d", txHash, nonce)
		return nil, nil, false
	}

	// check confirmations
	if !ob.HasEnoughConfirmations(receipt, ob.GetLastBlockHeight()) {
		log.Debug().Msgf("confirmTxByHash: txHash %s nonce %d included but not confirmed: receipt block %d, current block %d",
			txHash, nonce, receipt.BlockNumber, ob.GetLastBlockHeight())
		return nil, nil, false
	}

	// cross-check tx inclusion against the block
	// Note: a guard for false BlockNumber in receipt. The blob-carrying tx won't come here
	err = ob.CheckTxInclusion(transaction, receipt)
	if err != nil {
		log.Error().Err(err).Msgf("confirmTxByHash: checkTxInclusion error for txHash %s nonce %d", txHash, nonce)
		return nil, nil, false
	}

	return receipt, transaction, true
}
