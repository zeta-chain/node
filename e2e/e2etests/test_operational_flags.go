package e2etests

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// TestOperationalFlags tests the functionality of operations flags.
func TestOperationalFlags(r *runner.E2ERunner, _ []string) {
	_, err := r.Clients.Zetacore.Observer.OperationalFlags(
		r.Ctx,
		&observertypes.QueryOperationalFlagsRequest{},
	)
	require.NoError(r, err)

	currentHeight, err := r.Clients.Zetacore.GetBlockHeight(r.Ctx)
	require.NoError(r, err)

	// schedule a restart for 5 blocks in the future
	restartHeight := currentHeight + 5
	updateMsg := observertypes.NewMsgUpdateOperationalFlags(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		observertypes.OperationalFlags{
			RestartHeight: restartHeight,
		},
	)

	_, err = r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, updateMsg)
	require.NoError(r, err)

	operationalFlagsRes, err := r.Clients.Zetacore.Observer.OperationalFlags(
		r.Ctx,
		&observertypes.QueryOperationalFlagsRequest{},
	)
	require.NoError(r, err)
	require.Equal(r, restartHeight, operationalFlagsRes.OperationalFlags.RestartHeight)

	// TODO: wait for restart height + 2 then test that start timestamp metric has increased
}
