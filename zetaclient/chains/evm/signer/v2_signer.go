package signer

import (
	"context"
	"fmt"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/evm"
)

// SignOutboundFromCCTXV2 signs an outbound transaction from a CCTX with protocol contract v2
func (signer *Signer) SignOutboundFromCCTXV2(
	ctx context.Context,
	cctx *types.CrossChainTx,
	outboundData *OutboundData,
) (*ethtypes.Transaction, error) {
	outboundType := evm.ParseOutboundTypeFromCCTX(*cctx)
	switch outboundType {
	case evm.OutboundTypeGasWithdraw, evm.OutboundTypeGasWithdrawRevert:
		return signer.SignGasWithdraw(ctx, outboundData)
	case evm.OutboundTypeERC20Withdraw, evm.OutboundTypeERC20WithdrawRevert:
		return signer.SignERC20CustodyWithdraw(ctx, outboundData)
	case evm.OutboundTypeERC20WithdrawAndCall:
		return signer.SignERC20CustodyWithdrawAndCall(ctx, outboundData)
	case evm.OutboundTypeGasWithdrawAndCall, evm.OutboundTypeCall:
		// both gas withdraw and call and no-asset call uses gateway execute
		// no-asset call simply hash msg.value == 0
		return signer.SignGatewayExecute(ctx, outboundData)
	case evm.OutboundTypeGasWithdrawRevertAndCallOnRevert:
		return signer.SignGatewayExecuteRevert(ctx, outboundData)
	case evm.OutboundTypeERC20WithdrawRevertAndCallOnRevert:
		return signer.SignERC20CustodyWithdrawRevert(ctx, outboundData)
	}
	return nil, fmt.Errorf("unsupported outbound type %d", outboundType)
}
