package evm

import (
	"fmt"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.non-eth.sol"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
)

// ParseZetaEvent either returns a parsed ZetaReceived or ZetaReverted event
func ParseZetaEvent(
	receipt *ethtypes.Receipt,
	connectorAddr ethcommon.Address,
	connector *zetaconnector.ZetaConnectorNonEth,
	sendHash string,
	txHash string) (*zetaconnector.ZetaConnectorNonEthZetaReceived, *zetaconnector.ZetaConnectorNonEthZetaReverted, error) {
	for _, vLog := range receipt.Logs {
		receivedEvent, err := connector.ZetaConnectorNonEthFilterer.ParseZetaReceived(*vLog)
		if err == nil {
			// sanity check tx event
			err = ValidateEvmTxLog(vLog, connectorAddr, txHash, TopicsZetaReceived)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "error validating ZetaReceived event from outtx %s", txHash)
			}
			if vLog.Topics[3].Hex() != sendHash {
				return nil, nil, fmt.Errorf("cctx index mismatch in ZetaReceived event from outtx %s, want %s got %s", txHash, sendHash, vLog.Topics[3].Hex())
			}
			return receivedEvent, nil, nil
		}
		revertedEvent, err := connector.ZetaConnectorNonEthFilterer.ParseZetaReverted(*vLog)
		if err == nil {
			// sanity check tx event
			err = ValidateEvmTxLog(vLog, connectorAddr, receipt.TxHash.Hex(), TopicsZetaReverted)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "error validating ZetaReverted event from outtx %s", txHash)
			}
			return nil, revertedEvent, nil
		}
	}
	return nil, nil, fmt.Errorf("no ZetaReceived/ZetaReverted event found in outtx %s", txHash)
}

// ParseERC20WithdrawnEvent returns a parsed ERC20CustodyWithdrawn event from the outtx receipt
func ParseERC20WithdrawnEvent(
	receipt *ethtypes.Receipt,
	custodyAddr ethcommon.Address,
	custody *erc20custody.ERC20Custody,
	txHash string) (*erc20custody.ERC20CustodyWithdrawn, error) {
	for _, vLog := range receipt.Logs {
		withdrawnEvent, err := custody.ParseWithdrawn(*vLog)
		if err == nil {
			// sanity check tx event
			err = ValidateEvmTxLog(vLog, custodyAddr, receipt.TxHash.Hex(), TopicsWithdrawn)
			if err != nil {
				return nil, errors.Wrapf(err, "error validating ERC20CustodyWithdrawn event from outtx %s", txHash)
			}
			return withdrawnEvent, nil
		}
	}
	return nil, fmt.Errorf("no ERC20CustodyWithdrawn event found in outtx %s", txHash)
}

// ParseOuttxReceivedValue parses the received value from the outtx receipt
func ParseOuttxReceivedValue(
	receipt *ethtypes.Receipt,
	transaction *ethtypes.Transaction,
	cointype coin.CoinType,
	connectorAddress ethcommon.Address,
	connector *zetaconnector.ZetaConnectorNonEth,
	custodyAddress ethcommon.Address,
	custody *erc20custody.ERC20Custody,
	sendHash string,
	txHash string) (chains.ReceiveStatus, *big.Int, error) {
	// determine the receive status and value
	var receiveStatus chains.ReceiveStatus
	var receiveValue *big.Int
	switch receipt.Status {
	case ethtypes.ReceiptStatusSuccessful:
		receiveStatus = chains.ReceiveStatus_Success
		receiveValue = transaction.Value()
	case ethtypes.ReceiptStatusFailed:
		receiveStatus = chains.ReceiveStatus_Failed
		receiveValue = big.NewInt(0)
	default:
		// https://docs.nethereum.com/en/latest/nethereum-receipt-status/
		return chains.ReceiveStatus_Failed, nil, fmt.Errorf("unknown tx receipt status %d for outtx %s", receipt.Status, receipt.TxHash)
	}

	// parse receive value from the outtx receipt for Zeta and ERC20
	switch cointype {
	case coin.CoinType_Zeta:
		if receipt.Status == ethtypes.ReceiptStatusSuccessful {
			receivedLog, revertedLog, err := ParseZetaEvent(receipt, connectorAddress, connector, sendHash, txHash)
			if err != nil {
				return chains.ReceiveStatus_Failed, nil, err
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
			withdrawn, err := ParseERC20WithdrawnEvent(receipt, custodyAddress, custody, txHash)
			if err != nil {
				return chains.ReceiveStatus_Failed, nil, err
			}
			// use the value in Withdrawn event for vote message
			receiveValue = withdrawn.Amount
		}
	case coin.CoinType_Gas, coin.CoinType_Cmd:
		// nothing to do for CoinType_Gas/CoinType_Cmd, no need to parse event
	default:
		return chains.ReceiveStatus_Failed, nil, fmt.Errorf("unknown coin type %s for outtx %s", cointype, txHash)
	}
	return receiveStatus, receiveValue, nil
}
