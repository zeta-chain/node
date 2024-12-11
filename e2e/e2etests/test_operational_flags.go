package e2etests

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// TestOperationalFlags tests the functionality of operations flags.
func TestOperationalFlags(r *runner.E2ERunner, _ []string) {
	operationalFlagsRes, err := r.Clients.Zetacore.Observer.OperationalFlags(
		r.Ctx,
		&observertypes.QueryOperationalFlagsRequest{},
	)
	require.NoError(r, err)

	// always set to low height so it's ignored by zetaclient
	nextRestartHeight := operationalFlagsRes.OperationalFlags.RestartHeight + 1

	updateMsg := observertypes.NewMsgUpdateOperationalFlags(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		observertypes.OperationalFlags{
			RestartHeight: nextRestartHeight,
		},
	)

	_, err = r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, updateMsg)
	require.NoError(r, err)

	operationalFlagsRes, err = r.Clients.Zetacore.Observer.OperationalFlags(
		r.Ctx,
		&observertypes.QueryOperationalFlagsRequest{},
	)
	require.NoError(r, err)
	require.Equal(r, nextRestartHeight, operationalFlagsRes.OperationalFlags.RestartHeight)
}
