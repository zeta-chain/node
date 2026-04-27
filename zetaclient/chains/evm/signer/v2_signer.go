package signer

import (
	"fmt"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/evm/common"
)

// SignOutboundFromCCTXV2 signs an outbound transaction from a CCTX with protocol contract v2
func (signer *Signer) SignOutboundFromCCTXV2(
	cctx *types.CrossChainTx,
	outboundData *OutboundData,
) (*ethtypes.Transaction, error) {
	// V2 arbitrary calls that route through GatewayEVM.execute are not signed:
	// the destination contract and calldata are caller-controlled, so we
	// cancel via a TSS-to-TSS zero-value self-transfer that consumes the
	// assigned nonce. This avoids the head-of-line blocking that would occur
	// if the CCTX were left perpetually unsigned.
	//
	// Only OutboundTypeCall and OutboundTypeGasWithdrawAndCall reach
	// signGatewayExecute. The other arbitrary-call routes
	// (ERC20Custody.withdrawAndCall, ZetaConnector.withdrawAndCall) constrain
	// the destination to invoke typed `Callable.onCall` and remain enabled.
	// Plain withdraws emit `isArbitraryCall=true` from GatewayZEVM.withdraw()
	// too, but they don't reach signGatewayExecute either.
	outboundType := common.ParseOutboundTypeFromCCTX(*cctx)
	if outboundData.callOptions.IsArbitraryCall && common.IsGatewayExecuteOutbound(outboundType) {
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
