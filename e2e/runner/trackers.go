package runner

import (
	"github.com/stretchr/testify/require"

	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// EnsureNoTrackers ensures that there are no trackers left on zetacore
func (r *E2ERunner) EnsureNoTrackers() {
	// get all trackers
	res, err := r.CctxClient.OutTxTrackerAll(
		r.Ctx,
		&crosschaintypes.QueryAllOutboundTrackerRequest{},
	)
	require.NoError(r, err)
	require.Empty(r, res.OutboundTracker, "there should be no trackers at the end of the test")
}
