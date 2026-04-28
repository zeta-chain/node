package signer

import (
	"context"
	"fmt"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/evm/common"
	"github.com/zeta-chain/node/zetaclient/logs"
)

// SignOutboundFromCCTXV2 signs an outbound transaction from a CCTX with protocol contract v2
func (signer *Signer) SignOutboundFromCCTXV2(
	ctx context.Context,
	cctx *types.CrossChainTx,
	outboundData *OutboundData,
) (*ethtypes.Transaction, error) {
	outboundType := common.ParseOutboundTypeFromCCTX(*cctx)

	// V2 arbitrary calls whose CCTX shape lands at GatewayEVM's arbitrary-call
	// branch on the destination chain are not signed. The destination contract
	// and calldata are caller-controlled, so we cancel via a TSS-to-TSS
	// zero-value self-transfer that consumes the assigned nonce. This avoids
	// the head-of-line blocking that would occur if the CCTX were left
	// perpetually unsigned.
	//
	// IsArbitraryCallCancellable enumerates the outbound types whose signer
	// would otherwise produce MessageContext.Sender=0x0 and pass it into
	// GatewayEVM.execute (signGatewayExecute) or GatewayEVM.executeWithERC20
	// (via ERC20Custody.withdrawAndCall / ZetaConnector.withdrawAndCall) —
	// both of which branch on the zero sender to invoke _executeArbitraryCall.
	//
	// Plain withdraws emit `isArbitraryCall=true` from GatewayZEVM.withdraw()
	// per protocol-contracts semantics ("no authenticated sender"), but they
	// don't reach any of those entry points and are not cancelled.
	if outboundData.callOptions.IsArbitraryCall && common.IsArbitraryCallCancellable(outboundType) {
		signer.Logger().Std.Warn().
			Str(logs.FieldCctxIndex, cctx.Index).
			Uint64(logs.FieldNonce, outboundData.nonce).
			Stringer("destination", outboundData.to).
			Int("outbound_type", int(outboundType)).
			Msg("cancelling V2 arbitrary-call CCTX via TSS self-transfer")
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
