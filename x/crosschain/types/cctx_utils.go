package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// There can only be one active outtx.
// OutboundTxParams[0] is the original outtx, if it reverts, then
// OutboundTxParams[1] is the new outtx.
func (m *CrossChainTx) GetCurrentOutTxParam() *OutboundTxParams {
	if len(m.OutboundTxParams) == 0 {
		return &OutboundTxParams{}
	}
	return m.OutboundTxParams[len(m.OutboundTxParams)-1]
}

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
