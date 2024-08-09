package evm

import (
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// OutboundTypes enumerate the different types of outbound transactions
// NOTE: only used for v2 protocol contracts and currently excludes ZETA withdraws
type OutboundTypes int

const (
	// OutboundTypeUnknown is an unknown outbound transaction
	OutboundTypeUnknown OutboundTypes = iota

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

	// OutboundTypeGasWithdrawAndRevert is a gas withdraw and revert call
	OutboundTypeGasWithdrawAndRevert

	// OutboundTypeERC20WithdrawAndRevert is an ERC20 withdraw and revert call
	OutboundTypeERC20WithdrawAndRevert
)

// ParseOutboundTypeFromCCTX returns the outbound type from the CCTX
// TODO: address revert
func ParseOutboundTypeFromCCTX(cctx types.CrossChainTx) OutboundTypes {
	switch cctx.InboundParams.CoinType {
	case coin.CoinType_Gas:
		switch cctx.CctxStatus.Status {
		case types.CctxStatus_PendingOutbound:
			if len(cctx.RelayedMessage) == 0 {
				return OutboundTypeGasWithdraw
			} else {
				return OutboundTypeGasWithdrawAndCall
			}
		}
	case coin.CoinType_ERC20:
		switch cctx.CctxStatus.Status {
		case types.CctxStatus_PendingOutbound:
			if len(cctx.RelayedMessage) == 0 {
				return OutboundTypeERC20Withdraw
			} else {
				return OutboundTypeERC20WithdrawAndCall
			}
		}
	}

	return OutboundTypeUnknown
}
