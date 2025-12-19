package local

import (
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// replaceObserver performs the flow to replace an existing observer with a new validator while reusing the exiting tss key-shard
func replaceObserver(r *runner.E2ERunner, reuseTSSFrom string) {
	stakeToBecomeValidator(r)
	addGrantsWithHotkey(r, "zetaclient-new-validator", reuseTSSFrom)
	updateObserver(r, reuseTSSFrom)
}

// updateObserver updates the observer set to replace the old operator address with the new one.
func updateObserver(r *runner.E2ERunner, oldZetaclient string) {
	oldInfo, err := utils.FetchHotkeyAddress(oldZetaclient)
	require.NoError(r, err)
	oldObserver := oldInfo.ObserverAddress

	newInfo, err := utils.FetchHotkeyAddress("zetaclient-new-validator")
	require.NoError(r, err)
	newObserver := newInfo.ObserverAddress

	msg := observertypes.NewMsgUpdateObserver(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		oldObserver,
		newObserver,
		observertypes.ObserverUpdateReason_AdminUpdate,
	)

	_, err = r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, msg)
	require.NoError(r, err, "failed to broadcast MsgUpdateObserver")
}
