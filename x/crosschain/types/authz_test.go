package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestGetAllAuthzZetaclientTxTypes(t *testing.T) {
	require.Equal(t, []string{"/zetachain.zetacore.crosschain.MsgGasPriceVoter",
		"/zetachain.zetacore.crosschain.MsgVoteOnObservedInboundTx",
		"/zetachain.zetacore.crosschain.MsgVoteOnObservedOutboundTx",
		"/zetachain.zetacore.crosschain.MsgCreateTSSVoter",
		"/zetachain.zetacore.crosschain.MsgAddToOutTxTracker",
		"/zetachain.zetacore.observer.MsgAddBlameVote",
		"/zetachain.zetacore.observer.MsgAddBlockHeader"},
		crosschaintypes.GetAllAuthzZetaclientTxTypes())
}
