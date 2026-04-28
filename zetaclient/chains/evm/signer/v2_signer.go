package signer

import (
	"fmt"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/evm/common"
	"github.com/zeta-chain/node/zetaclient/logs"
)

// SignOutboundFromCCTXV2 signs an outbound transaction from a CCTX with protocol contract v2
func (signer *Signer) SignOutboundFromCCTXV2(
	cctx *types.CrossChainTx,
	outboundData *OutboundData,
) (*ethtypes.Transaction, error) {
	// V2 arbitrary calls reach GatewayEVM's `_executeArbitraryCall` branch
	// (sender == address(0)), which performs a raw destination.call(data) with
	// msg.sender == GatewayEVM. The destination and calldata are caller-
	// controlled, so we cancel these CCTXs at the signer via a TSS-to-TSS
	// zero-value self-transfer that consumes the assigned nonce. This avoids
	// the head-of-line blocking that would occur if the CCTX were left
	// perpetually unsigned.
	//
	// Three signer entry points produce a sender == 0 MessageContext and
	// land at the same arbitrary-call branch:
	//   - signGatewayExecute → GatewayEVM.execute
	//     (OutboundTypeCall, OutboundTypeGasWithdrawAndCall)
	//   - signERC20CustodyWithdrawAndCall → GatewayEVM.executeWithERC20
	//     (OutboundTypeERC20WithdrawAndCall)
	//   - signZetaConnectorWithdrawAndCall → GatewayEVM.executeWithERC20
	//     (OutboundTypeZetaWithdrawAndCall)
	//
	// Plain withdraws also emit `isArbitraryCall=true` from
	// GatewayZEVM.withdraw(), but they don't reach any of these entry points.
	outboundType := common.ParseOutboundTypeFromCCTX(*cctx)
	if outboundData.callOptions.IsArbitraryCall && common.IsArbitraryCallCancellable(outboundType) {
		signer.Logger().Std.Warn().
			Str(logs.FieldCctxIndex, cctx.Index).
			Uint64(logs.FieldNonce, outboundData.nonce).
			Stringer("destination", outboundData.to).
			Int("outbound_type", int(outboundType)).
			Msg("cancelling V2 arbitrary-call CCTX via TSS self-transfer")
		return signer.SignCancel(outboundData)
	}
	switch outboundType {
	case common.OutboundTypeGasWithdraw, common.OutboundTypeGasWithdrawRevert:
		return signer.SignGasWithdraw(outboundData)
	case common.OutboundTypeERC20Withdraw, common.OutboundTypeERC20WithdrawRevert:
		return signer.signERC20CustodyWithdraw(outboundData)
	case common.OutboundTypeERC20WithdrawAndCall:
		return signer.signERC20CustodyWithdrawAndCall(outboundData)
	case common.OutboundTypeZetaWithdrawRevert, common.OutboundTypeZetaWithdraw:
		return signer.signZetaConnectorWithdraw(outboundData)
	case common.OutboundTypeZetaWithdrawAndCall:
		return signer.signZetaConnectorWithdrawAndCall(outboundData)
	case common.OutboundTypeGasWithdrawAndCall, common.OutboundTypeCall:
		// both gas withdraw and call and no-asset call uses gateway execute
		// no-asset call simply hash msg.value == 0
		return signer.signGatewayExecute(outboundData)
	case common.OutboundTypeGasWithdrawRevertAndCallOnRevert:
		return signer.signGatewayExecuteRevert(cctx.InboundParams.Sender, outboundData)
	case common.OutboundTypeERC20WithdrawRevertAndCallOnRevert:
		return signer.signERC20CustodyWithdrawRevert(cctx.InboundParams.Sender, outboundData)
	case common.OutboundTypeZetaWithdrawRevertAndCallOnRevert:
		return signer.signZetaConnectorWithdrawRevert(cctx.InboundParams.Sender, outboundData)
	}
	return nil, fmt.Errorf("unsupported outbound type %d", outboundType)
}
