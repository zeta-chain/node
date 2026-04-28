package common

import (
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// OutboundType enumerate the different types of outbound transactions
// NOTE: only used for v2 protocol contracts and currently excludes ZETA withdraws
type OutboundType int

const (
	// OutboundTypeUnknown is an unknown outbound transaction
	OutboundTypeUnknown OutboundType = iota

	// OutboundTypeGasWithdraw is a gas withdraw transaction
	OutboundTypeGasWithdraw

	// OutboundTypeERC20Withdraw is an ERC20 withdraw transaction
	OutboundTypeERC20Withdraw

	// OutboundTypeGasWithdrawAndCall is a gas withdraw and call transaction
	OutboundTypeGasWithdrawAndCall

	// OutboundTypeERC20WithdrawAndCall is an ERC20 withdraw and call transaction
	OutboundTypeERC20WithdrawAndCall

	// OutboundTypeCall is a no-asset call transaction
	OutboundTypeCall

	// OutboundTypeGasWithdrawRevert is a gas withdraw revert
	OutboundTypeGasWithdrawRevert

	// OutboundTypeGasWithdrawRevertAndCallOnRevert is a gas withdraw revert and call on revert
	OutboundTypeGasWithdrawRevertAndCallOnRevert

	// OutboundTypeERC20WithdrawRevert is an ERC20 withdraw revert
	OutboundTypeERC20WithdrawRevert

	// OutboundTypeERC20WithdrawRevertAndCallOnRevert is an ERC20 withdraw revert and call on revert
	OutboundTypeERC20WithdrawRevertAndCallOnRevert

	OutboundTypeZetaWithdrawRevertAndCallOnRevert

	OutboundTypeZetaWithdrawRevert

	OutboundTypeZetaWithdrawAndCall

	OutboundTypeZetaWithdraw
)

// IsArbitraryCallCancellable reports whether the outbound type, when paired
// with CallOptions.IsArbitraryCall=true, would produce a TSS-signed call
// containing MessageContext.Sender == address(0) and thereby reach
// GatewayEVM's arbitrary-call branch (_executeArbitraryCall) on the
// destination chain.
//
// Three signer entry points translate IsArbitraryCall into a zero sender
// inside MessageContext (see OutboundData.MessageContext) and pass it into
// a contract that branches on sender:
//   - signGatewayExecute → GatewayEVM.execute
//     reached by OutboundTypeCall and OutboundTypeGasWithdrawAndCall
//   - signERC20CustodyWithdrawAndCall → ERC20Custody.withdrawAndCall
//     → GatewayEVM.executeWithERC20
//     reached by OutboundTypeERC20WithdrawAndCall
//   - signZetaConnectorWithdrawAndCall → ZetaConnector.withdrawAndCall
//     → GatewayEVM.executeWithERC20
//     reached by OutboundTypeZetaWithdrawAndCall
//
// All three end at GatewayEVM's arbitrary-call branch, which performs a raw
// destination.call(data) with msg.sender == GatewayEVM and is not selector-
// filtered against ERC20 authority (transferFrom / approve / permit). CCTXs
// with IsArbitraryCall=true on any of these outbound types must be cancelled
// at the signer rather than relayed to the destination chain.
func IsArbitraryCallCancellable(t OutboundType) bool {
	return t == OutboundTypeCall ||
		t == OutboundTypeGasWithdrawAndCall ||
		t == OutboundTypeERC20WithdrawAndCall ||
		t == OutboundTypeZetaWithdrawAndCall
}

// ParseOutboundTypeFromCCTX returns the outbound type from the CCTX
func ParseOutboundTypeFromCCTX(cctx types.CrossChainTx) OutboundType {
	switch cctx.InboundParams.CoinType {
	case coin.CoinType_Gas:
		switch cctx.CctxStatus.Status {
		case types.CctxStatus_PendingOutbound:
			if cctx.InboundParams.IsCrossChainCall {
				return OutboundTypeGasWithdrawAndCall
			}

			return OutboundTypeGasWithdraw
		case types.CctxStatus_PendingRevert:
			if cctx.RevertOptions.CallOnRevert {
				return OutboundTypeGasWithdrawRevertAndCallOnRevert
			}

			return OutboundTypeGasWithdrawRevert
		}
	case coin.CoinType_ERC20:
		switch cctx.CctxStatus.Status {
		case types.CctxStatus_PendingOutbound:
			if cctx.InboundParams.IsCrossChainCall {
				return OutboundTypeERC20WithdrawAndCall
			}

			return OutboundTypeERC20Withdraw
		case types.CctxStatus_PendingRevert:
			if cctx.RevertOptions.CallOnRevert {
				return OutboundTypeERC20WithdrawRevertAndCallOnRevert
			}

			return OutboundTypeERC20WithdrawRevert
		}
	case coin.CoinType_Zeta:
		switch cctx.CctxStatus.Status {
		case types.CctxStatus_PendingOutbound:
			if cctx.InboundParams.IsCrossChainCall {
				return OutboundTypeZetaWithdrawAndCall
			}
			return OutboundTypeZetaWithdraw
		case types.CctxStatus_PendingRevert:
			if cctx.RevertOptions.CallOnRevert {
				return OutboundTypeZetaWithdrawRevertAndCallOnRevert
			}
			return OutboundTypeZetaWithdrawRevert
		}
	case coin.CoinType_NoAssetCall:
		if cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound {
			return OutboundTypeCall
		}
	}

	return OutboundTypeUnknown
}
