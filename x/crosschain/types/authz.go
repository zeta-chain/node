package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// GetAllAuthzZetaclientTxTypes returns all the authz types for required for zetaclient
func GetAllAuthzZetaclientTxTypes() []string {
	return []string{
		sdk.MsgTypeURL(&MsgGasPriceVoter{}),
		sdk.MsgTypeURL(&MsgVoteOnObservedInboundTx{}),
		sdk.MsgTypeURL(&MsgVoteOnObservedOutboundTx{}),
		sdk.MsgTypeURL(&MsgCreateTSSVoter{}),
		sdk.MsgTypeURL(&MsgAddToOutTxTracker{}),
		sdk.MsgTypeURL(&observertypes.MsgAddBlameVote{}),
		sdk.MsgTypeURL(&observertypes.MsgAddBlockHeader{}),
	}
}
