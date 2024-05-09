package evm

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.non-eth.sol"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/compliance"
)

// PostVoteOutbound posts vote to zetacore for the confirmed outtx
func (ob *ChainClient) PostVoteOutbound(
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
	zetaTxHash, ballot, err := ob.zetaBridge.PostVoteOutbound(
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
func (ob *ChainClient) IsOutboundProcessed(cctx *crosschaintypes.CrossChainTx, logger zerolog.Logger) (bool, bool, error) {
	// skip if outtx is not confirmed
	nonce := cctx.GetCurrentOutboundParam().TssNonce
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
	cointype := cctx.InboundParams.CoinType

	// compliance check, special handling the cancelled cctx
	if compliance.IsCctxRestricted(cctx) {
		// use cctx's amount to bypass the amount check in zetacore
		receiveValue = cctx.GetCurrentOutboundParam().Amount.BigInt()
		receiveStatus := chains.ReceiveStatus_failed
		if receipt.Status == ethtypes.ReceiptStatusSuccessful {
			receiveStatus = chains.ReceiveStatus_success
		}
		ob.PostVoteOutbound(cctx.Index, receipt, transaction, receiveValue, receiveStatus, nonce, cointype, logger)
		return true, true, nil
	}

	// parse the received value from the outtx receipt
	receiveValue, receiveStatus, err = ParseOutboundReceivedValue(cctx, receipt, transaction, cointype, connectorAddr, connector, custodyAddr, custody)
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
	params := cctx.GetCurrentOutboundParam()
	for _, vLog := range receipt.Logs {
		// try parsing ZetaReceived event
		received, err := connector.ZetaConnectorNonEthFilterer.ParseZetaReceived(*vLog)
		if err == nil {
			err = ValidateEvmTxLog(vLog, connectorAddr, receipt.TxHash.Hex(), TopicsZetaReceived)
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
			err = ValidateEvmTxLog(vLog, connectorAddr, receipt.TxHash.Hex(), TopicsZetaReverted)
			if err != nil {
				return nil, nil, errors.Wrap(err, "error validating ZetaReverted event")
			}
			if !strings.EqualFold(ethcommon.BytesToAddress(reverted.DestinationAddress[:]).Hex(), cctx.InboundParams.Sender) {
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

// ParseAndCheckWithdrawnEvent parses and checks erc20 Withdrawn event from the outtx receipt
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
			err = ValidateEvmTxLog(vLog, custodyAddr, receipt.TxHash.Hex(), TopicsWithdrawn)
			if err != nil {
				return nil, errors.Wrap(err, "error validating Withdrawn event")
			}
			if !strings.EqualFold(withdrawn.Recipient.Hex(), params.Receiver) {
				return nil, fmt.Errorf("receiver address mismatch in Withdrawn event, want %s got %s",
					params.Receiver, withdrawn.Recipient.Hex())
			}
			if !strings.EqualFold(withdrawn.Asset.Hex(), cctx.InboundParams.Asset) {
				return nil, fmt.Errorf("asset mismatch in Withdrawn event, want %s got %s",
					cctx.InboundParams.Asset, withdrawn.Asset.Hex())
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

// ParseOutboundReceivedValue parses the received value and status from the outtx receipt
// The receivd value is the amount of Zeta/ERC20/Gas token (released from connector/custody/TSS) sent to the receiver
func ParseOutboundReceivedValue(
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
