package observer

import (
	"fmt"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/evm"
	"math/big"
)

// ParseOutboundEventV2 parses an event from an outbound with protocol contract v2
func ParseOutboundEventV2(
	cctx *types.CrossChainTx,
	receipt *ethtypes.Receipt,
	transaction *ethtypes.Transaction,
) (*big.Int, chains.ReceiveStatus, error) {
	receivedValue := big.NewInt(0)
	status := chains.ReceiveStatus_failed
	if receipt.Status == ethtypes.ReceiptStatusSuccessful {
		receivedValue = transaction.Value()
		status = chains.ReceiveStatus_success
	}

	outboundType := evm.ParseOutboundTypeFromCCTX(*cctx)
	switch outboundType {
	case evm.OutboundTypeGasWithdraw:
		return receivedValue, status, nil
	case evm.OutboundTypeERC20Withdraw:
		return receivedValue, status, nil
	case evm.OutboundTypeGasWithdrawAndCall:
		return receivedValue, status, nil
	case evm.OutboundTypeERC20WithdrawAndCall:
		return receivedValue, status, nil
	}
	return receivedValue, status, fmt.Errorf("unsupported outbound type %d", outboundType)
}
