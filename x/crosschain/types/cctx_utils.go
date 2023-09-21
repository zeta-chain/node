package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// GetCurrentOutTxParam returns the current outbound tx params.
// There can only be one active outtx.
// OutboundTxParams[0] is the original outtx, if it reverts, then
// OutboundTxParams[1] is the new outtx.
func (m CrossChainTx) GetCurrentOutTxParam() *OutboundTxParams {
	if len(m.OutboundTxParams) == 0 {
		return &OutboundTxParams{}
	}
	return m.OutboundTxParams[len(m.OutboundTxParams)-1]
}

// IsCurrentOutTxRevert returns true if the current outbound tx is the revert tx.
func (m CrossChainTx) IsCurrentOutTxRevert() bool {
	return len(m.OutboundTxParams) == 2
}

// OriginalDestinationChainID returns the original destination of the outbound tx, reverted or not
// If there is no outbound tx, return -1
func (m CrossChainTx) OriginalDestinationChainID() int64 {
	if len(m.OutboundTxParams) == 0 {
		return -1
	}
	return m.OutboundTxParams[0].ReceiverChainId
}

// GetAllAuthzZetaclientTxTypes returns all the authz types for zetaclient
func GetAllAuthzZetaclientTxTypes() []string {
	return []string{
		sdk.MsgTypeURL(&MsgNonceVoter{}),
		sdk.MsgTypeURL(&MsgGasPriceVoter{}),
		sdk.MsgTypeURL(&MsgVoteOnObservedInboundTx{}),
		sdk.MsgTypeURL(&MsgVoteOnObservedOutboundTx{}),
		sdk.MsgTypeURL(&MsgSetNodeKeys{}),
		sdk.MsgTypeURL(&MsgCreateTSSVoter{}),
		sdk.MsgTypeURL(&MsgAddToOutTxTracker{}),
		sdk.MsgTypeURL(&MsgSetNodeKeys{}),
		sdk.MsgTypeURL(&types.MsgAddBlameVote{}),
		sdk.MsgTypeURL(&types.MsgAddBlockHeader{}),
	}
}
