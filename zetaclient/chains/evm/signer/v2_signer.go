package signer

import (
	"context"
	"fmt"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/evm/common"
)

// SignOutboundFromCCTXV2 signs an outbound transaction from a CCTX with protocol contract v2
func (signer *Signer) SignOutboundFromCCTXV2(
	ctx context.Context,
	cctx *types.CrossChainTx,
	outboundData *OutboundData,
) (*ethtypes.Transaction, error) {
	outboundType := common.ParseOutboundTypeFromCCTX(*cctx)

	// V2 arbitrary calls that route through GatewayEVM.execute are not signed:
	// the destination contract and calldata are caller-controlled, so we
	// cancel via a TSS-to-TSS zero-value self-transfer that consumes the
	// assigned nonce. This avoids the head-of-line blocking that would occur
	// if the CCTX were left perpetually unsigned.
	//
	// Only OutboundTypeCall and OutboundTypeGasWithdrawAndCall reach
	// signGatewayExecute. The other arbitrary-call routes
	// (ERC20Custody.withdrawAndCall) constrain the destination to invoke
	// typed `Callable.onCall` and remain enabled. Plain withdraws emit
	// `isArbitraryCall=true` from GatewayZEVM.withdraw() too, but they don't
	// reach signGatewayExecute either.
	if outboundData.callOptions.IsArbitraryCall && common.IsGatewayExecuteOutbound(outboundType) {
		return signer.SignCancel(ctx, outboundData)
	}
	switch outboundType {
	case common.OutboundTypeGasWithdraw, common.OutboundTypeGasWithdrawRevert:
		return signer.SignGasWithdraw(ctx, outboundData)
	case common.OutboundTypeERC20Withdraw, common.OutboundTypeERC20WithdrawRevert:
		return signer.signERC20CustodyWithdraw(ctx, outboundData)
	case common.OutboundTypeERC20WithdrawAndCall:
		return signer.signERC20CustodyWithdrawAndCall(ctx, outboundData)
	case common.OutboundTypeZetaWithdrawRevert:
		return signer.signZetaConnectorWithdraw(ctx, outboundData)
	// Add when implementing Zeta withdraws
	//common.OutboundTypeZetaWithdraw and common.OutboundTypeZetaWithdrawAndCall
	case common.OutboundTypeGasWithdrawAndCall, common.OutboundTypeCall:
		// both gas withdraw and call and no-asset call uses gateway execute
		// no-asset call simply hash msg.value == 0
		return signer.signGatewayExecute(ctx, outboundData)
	case common.OutboundTypeGasWithdrawRevertAndCallOnRevert:
		return signer.signGatewayExecuteRevert(ctx, cctx.InboundParams.Sender, outboundData)
	case common.OutboundTypeERC20WithdrawRevertAndCallOnRevert:
		return signer.signERC20CustodyWithdrawRevert(ctx, cctx.InboundParams.Sender, outboundData)
	case common.OutboundTypeZetaWithdrawRevertAndCallOnRevert:
		return signer.signZetaConnectorWithdrawRevert(ctx, cctx.InboundParams.Sender, outboundData)
	}
	return nil, fmt.Errorf("unsupported outbound type %d", outboundType)
}
