package observer

import (
	"bytes"
	"encoding/hex"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/gatewayevm.sol"
	"github.com/zeta-chain/zetacore/pkg/chains"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/evm"
	"math/big"
	"strings"
)

// ParseOutboundEventV2 parses an event from an outbound with protocol contract v2
func ParseOutboundEventV2(
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
	case evm.OutboundTypeGasWithdraw:
		// simple transfer, no need to parse event
		return transaction.Value(), chains.ReceiveStatus_success, nil
	//case evm.OutboundTypeERC20Withdraw:
	//	return big.NewInt(0), chains.ReceiveStatus_failed, nil
	case evm.OutboundTypeGasWithdrawAndCall:
		return ParseAndCheckGatewayExecuted(cctx, receipt, gatewayAddr, gateway)
		//case evm.OutboundTypeERC20WithdrawAndCall:
		//	return big.NewInt(0), chains.ReceiveStatus_failed, nil
	}
	return big.NewInt(0), chains.ReceiveStatus_failed, fmt.Errorf("unsupported outbound type %d", outboundType)
}

// ParseAndCheckGatewayExecuted parses and checks the gateway execute event
func ParseAndCheckGatewayExecuted(
	cctx *crosschaintypes.CrossChainTx,
	receipt *ethtypes.Receipt,
	gatewayAddr ethcommon.Address,
	gateway *gatewayevm.GatewayEVM,
) (*big.Int, chains.ReceiveStatus, error) {
	params := cctx.GetCurrentOutboundParam()

	for _, vLog := range receipt.Logs {
		executed, err := gateway.GatewayEVMFilterer.ParseExecuted(*vLog)
		if err == nil {
			// basic event check
			if err := evm.ValidateEvmTxLog(vLog, gatewayAddr, receipt.TxHash.Hex(), evm.TopicsGatewayExecuted); err != nil {
				return big.NewInt(0), chains.ReceiveStatus_failed, errors.Wrap(err, "failed to validate gateway executed event")
			}
			// destination
			if !strings.EqualFold(executed.Destination.Hex(), params.Receiver) {
				return big.NewInt(0), chains.ReceiveStatus_failed, fmt.Errorf("receiver address mismatch in event, want %s got %s",
					params.Receiver, executed.Destination.Hex())
			}
			// amount
			if executed.Value.Cmp(params.Amount.BigInt()) != 0 {
				return big.NewInt(0), chains.ReceiveStatus_failed, fmt.Errorf("amount mismatch in event, want %s got %s",
					params.Amount.String(), executed.Value.String())
			}
			// data
			if err := checkCCTXMessage(executed.Data, cctx.RelayedMessage); err != nil {
				return big.NewInt(0), chains.ReceiveStatus_failed, err
			}

			return executed.Value, chains.ReceiveStatus_success, nil
		}
	}

	return big.NewInt(0), chains.ReceiveStatus_failed, errors.New("gateway execute event not found")
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
