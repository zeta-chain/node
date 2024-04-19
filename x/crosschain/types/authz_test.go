package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestGetAllAuthzZetaclientTxTypes(t *testing.T) {
	require.Equal(t, []string{"/zetachain.zetacore.crosschain.MsgVoteGasPrice",
		"/zetachain.zetacore.crosschain.MsgVoteOnObservedInboundTx",
		"/zetachain.zetacore.crosschain.MsgVoteOnObservedOutboundTx",
		"/zetachain.zetacore.crosschain.MsgAddToOutTxTracker",
		"/zetachain.zetacore.observer.MsgVoteTSS",
		"/zetachain.zetacore.observer.MsgAddBlameVote",
		"/zetachain.zetacore.observer.MsgVoteBlockHeader"},
		crosschaintypes.GetAllAuthzZetaclientTxTypes())
}
