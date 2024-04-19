package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestGetAllAuthzZetaclientTxTypes(t *testing.T) {
	require.Equal(t, []string{"/crosschain.MsgVoteGasPrice",
		"/crosschain.MsgVoteOnObservedInboundTx",
		"/crosschain.MsgVoteOnObservedOutboundTx",
		"/crosschain.MsgAddToOutTxTracker",
		"/observer.MsgVoteTSS",
		"/observer.MsgAddBlameVote",
		"/observer.MsgVoteBlockHeader"},
		crosschaintypes.GetAllAuthzZetaclientTxTypes())
}
