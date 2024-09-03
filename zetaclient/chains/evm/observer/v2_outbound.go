package observer

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayevm.sol"

	"github.com/zeta-chain/node/pkg/chains"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/evm"
)

// parseOutboundEventV2 parses an event from an outbound with protocol contract v2
func parseOutboundEventV2(
	cctx *crosschaintypes.CrossChainTx,
	receipt *ethtypes.Receipt,
	transaction *ethtypes.Transaction,
	custodyAddr ethcommon.Address,
	custody *erc20custody.ERC20Custody,
	gatewayAddr ethcommon.Address,
	gateway *gatewayevm.GatewayEVM,
) (*big.Int, chains.ReceiveStatus, error) {
	// return failed status if receipt status is failed
	if receipt.Status == ethtypes.ReceiptStatusFailed {
		return big.NewInt(0), chains.ReceiveStatus_failed, nil
	}

	outboundType := evm.ParseOutboundTypeFromCCTX(*cctx)
	switch outboundType {
	case evm.OutboundTypeGasWithdraw, evm.OutboundTypeGasWithdrawRevert:
		// simple transfer, no need to parse event
		return transaction.Value(), chains.ReceiveStatus_success, nil
	case evm.OutboundTypeERC20Withdraw, evm.OutboundTypeERC20WithdrawRevert:
		return parseAndCheckERC20CustodyWithdraw(cctx, receipt, custodyAddr, custody)
	case evm.OutboundTypeERC20WithdrawAndCall:
		return parseAndCheckERC20CustodyWithdrawAndCall(cctx, receipt, custodyAddr, custody)
	case evm.OutboundTypeGasWithdrawAndCall, evm.OutboundTypeCall:
		// both gas withdraw and call and no-asset call uses gateway execute
		// no-asset call simply hash msg.value == 0
		return parseAndCheckGatewayExecuted(cctx, receipt, gatewayAddr, gateway)
	case evm.OutboundTypeGasWithdrawRevertAndCallOnRevert, evm.OutboundTypeERC20WithdrawRevertAndCallOnRevert:
		return parseAndCheckGatewayReverted(cctx, receipt, gatewayAddr, gateway)
	}
	return big.NewInt(0), chains.ReceiveStatus_failed, fmt.Errorf("unsupported outbound type %d", outboundType)
}

// parseAndCheckGatewayExecuted parses and checks the gateway execute event
func parseAndCheckGatewayExecuted(
	cctx *crosschaintypes.CrossChainTx,
	receipt *ethtypes.Receipt,
	gatewayAddr ethcommon.Address,
	gateway *gatewayevm.GatewayEVM,
) (*big.Int, chains.ReceiveStatus, error) {
	params := cctx.GetCurrentOutboundParam()

	for _, vLog := range receipt.Logs {
		executed, err := gateway.GatewayEVMFilterer.ParseExecuted(*vLog)
		if err != nil {
			continue
		}
		// basic event check
		if err := evm.ValidateEvmTxLog(vLog, gatewayAddr, receipt.TxHash.Hex(), evm.TopicsGatewayExecuted); err != nil {
			return big.NewInt(
					0,
				), chains.ReceiveStatus_failed, errors.Wrap(
					err,
					"failed to validate gateway executed event",
				)
		}
		// destination
		if !strings.EqualFold(executed.Destination.Hex(), params.Receiver) {
			return big.NewInt(
					0,
				), chains.ReceiveStatus_failed, fmt.Errorf(
					"receiver address mismatch in event, want %s got %s",
					params.Receiver,
					executed.Destination.Hex(),
				)
		}
		// amount
		if executed.Value.Cmp(params.Amount.BigInt()) != 0 {
			return big.NewInt(
					0,
				), chains.ReceiveStatus_failed, fmt.Errorf(
					"amount mismatch in event, want %s got %s",
					params.Amount.String(),
					executed.Value.String(),
				)
		}
		// data
		if err := checkCCTXMessage(executed.Data, cctx.RelayedMessage); err != nil {
			return big.NewInt(0), chains.ReceiveStatus_failed, err
		}

		return executed.Value, chains.ReceiveStatus_success, nil
	}

	return big.NewInt(0), chains.ReceiveStatus_failed, errors.New("gateway execute event not found")
}

// parseAndCheckGatewayReverted parses and checks the gateway reverted event
func parseAndCheckGatewayReverted(
	cctx *crosschaintypes.CrossChainTx,
	receipt *ethtypes.Receipt,
	gatewayAddr ethcommon.Address,
	gateway *gatewayevm.GatewayEVM,
) (*big.Int, chains.ReceiveStatus, error) {
	params := cctx.GetCurrentOutboundParam()

	for _, vLog := range receipt.Logs {
		reverted, err := gateway.GatewayEVMFilterer.ParseReverted(*vLog)
		if err != nil {
			continue
		}
		// basic event check
		if err := evm.ValidateEvmTxLog(vLog, gatewayAddr, receipt.TxHash.Hex(), evm.TopicsGatewayReverted); err != nil {
			return big.NewInt(
					0,
				), chains.ReceiveStatus_failed, errors.Wrap(
					err,
					"failed to validate gateway reverte event",
				)
		}
		// destination
		if !strings.EqualFold(reverted.To.Hex(), params.Receiver) {
			return big.NewInt(
					0,
				), chains.ReceiveStatus_failed, fmt.Errorf(
					"receiver address mismatch in event, want %s got %s",
					params.Receiver,
					reverted.To.Hex(),
				)
		}
		// token
		if !strings.EqualFold(reverted.Token.Hex(), cctx.InboundParams.Asset) {
			return big.NewInt(
					0,
				), chains.ReceiveStatus_failed, fmt.Errorf(
					"asset address mismatch in event, want %s got %s",
					cctx.InboundParams.Asset,
					reverted.Token.Hex(),
				)
		}
		// amount
		if reverted.Amount.Cmp(params.Amount.BigInt()) != 0 {
			return big.NewInt(
					0,
				), chains.ReceiveStatus_failed, fmt.Errorf(
					"amount mismatch in event, want %s got %s",
					params.Amount.String(),
					reverted.Amount.String(),
				)
		}

		return reverted.Amount, chains.ReceiveStatus_success, nil
	}

	return big.NewInt(0), chains.ReceiveStatus_failed, errors.New("erc20 custody withdraw event not found")
}

// parseAndCheckERC20CustodyWithdraw parses and checks the ERC20 custody withdraw event
func parseAndCheckERC20CustodyWithdraw(
	cctx *crosschaintypes.CrossChainTx,
	receipt *ethtypes.Receipt,
	custodyAddr ethcommon.Address,
	custody *erc20custody.ERC20Custody,
) (*big.Int, chains.ReceiveStatus, error) {
	params := cctx.GetCurrentOutboundParam()

	for _, vLog := range receipt.Logs {
		withdrawn, err := custody.ERC20CustodyFilterer.ParseWithdrawn(*vLog)
		if err != nil {
			continue
		}
		// basic event check
		if err := evm.ValidateEvmTxLog(vLog, custodyAddr, receipt.TxHash.Hex(), evm.TopicsERC20CustodyWithdraw); err != nil {
			return big.NewInt(
					0,
				), chains.ReceiveStatus_failed, errors.Wrap(
					err,
					"failed to validate erc20 custody withdrawn event",
				)
		}
		// destination
		if !strings.EqualFold(withdrawn.To.Hex(), params.Receiver) {
			return big.NewInt(
					0,
				), chains.ReceiveStatus_failed, fmt.Errorf(
					"receiver address mismatch in event, want %s got %s",
					params.Receiver,
					withdrawn.To.Hex(),
				)
		}
		// token
		if !strings.EqualFold(withdrawn.Token.Hex(), cctx.InboundParams.Asset) {
			return big.NewInt(
					0,
				), chains.ReceiveStatus_failed, fmt.Errorf(
					"asset address mismatch in event, want %s got %s",
					cctx.InboundParams.Asset,
					withdrawn.Token.Hex(),
				)
		}
		// amount
		if withdrawn.Amount.Cmp(params.Amount.BigInt()) != 0 {
			return big.NewInt(
					0,
				), chains.ReceiveStatus_failed, fmt.Errorf(
					"amount mismatch in event, want %s got %s",
					params.Amount.String(),
					withdrawn.Amount.String(),
				)
		}

		return withdrawn.Amount, chains.ReceiveStatus_success, nil
	}

	return big.NewInt(0), chains.ReceiveStatus_failed, errors.New("erc20 custody withdraw event not found")
}

// parseAndCheckERC20CustodyWithdrawAndCall parses and checks the ERC20 custody withdraw and call event
func parseAndCheckERC20CustodyWithdrawAndCall(
	cctx *crosschaintypes.CrossChainTx,
	receipt *ethtypes.Receipt,
	custodyAddr ethcommon.Address,
	custody *erc20custody.ERC20Custody,
) (*big.Int, chains.ReceiveStatus, error) {
	params := cctx.GetCurrentOutboundParam()

	for _, vLog := range receipt.Logs {
		withdrawn, err := custody.ERC20CustodyFilterer.ParseWithdrawnAndCalled(*vLog)
		if err != nil {
			continue
		}
		// basic event check
		if err := evm.ValidateEvmTxLog(vLog, custodyAddr, receipt.TxHash.Hex(), evm.TopicsERC20CustodyWithdrawAndCall); err != nil {
			return big.NewInt(
					0,
				), chains.ReceiveStatus_failed, errors.Wrap(
					err,
					"failed to validate erc20 custody withdraw and call event",
				)
		}
		// destination
		if !strings.EqualFold(withdrawn.To.Hex(), params.Receiver) {
			return big.NewInt(
					0,
				), chains.ReceiveStatus_failed, fmt.Errorf(
					"receiver address mismatch in event, want %s got %s",
					params.Receiver,
					withdrawn.To.Hex(),
				)
		}
		// token
		if !strings.EqualFold(withdrawn.Token.Hex(), cctx.InboundParams.Asset) {
			return big.NewInt(
					0,
				), chains.ReceiveStatus_failed, fmt.Errorf(
					"asset address mismatch in event, want %s got %s",
					cctx.InboundParams.Asset,
					withdrawn.Token.Hex(),
				)
		}
		// amount
		if withdrawn.Amount.Cmp(params.Amount.BigInt()) != 0 {
			return big.NewInt(
					0,
				), chains.ReceiveStatus_failed, fmt.Errorf(
					"amount mismatch in event, want %s got %s",
					params.Amount.String(),
					withdrawn.Amount.String(),
				)
		}
		// data
		if err := checkCCTXMessage(withdrawn.Data, cctx.RelayedMessage); err != nil {
			return big.NewInt(0), chains.ReceiveStatus_failed, err
		}

		return withdrawn.Amount, chains.ReceiveStatus_success, nil
	}

	return big.NewInt(0), chains.ReceiveStatus_failed, errors.New("erc20 custody withdraw and call event not found")
}

// checkCCTXMessage checks the message of cctx with the emitted data of the event
func checkCCTXMessage(emittedData []byte, message string) error {
	messageBytes, err := hex.DecodeString(message)
	if err != nil {
		return errors.Wrap(err, "failed to decode message")
	}
	if !bytes.Equal(emittedData, messageBytes) {
		return fmt.Errorf("message mismatch, want %s got %s", message, hex.EncodeToString(emittedData))
	}
	return nil
}
