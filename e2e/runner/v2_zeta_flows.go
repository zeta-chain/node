package runner

import (
	"github.com/zeta-chain/node/e2e/utils"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// EnableV2ZETAFlows sends a message to enable V2 ZETA gateway flows
func (r *E2ERunner) EnableV2ZETAFlows() error {
	r.Logger.Print("enabling V2 ZETA gateway flows")

	msg := observertypes.NewMsgUpdateV2ZetaFlows(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		true,
	)
	_, err := r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, msg)
	if err != nil {
		return err
	}

	r.Logger.Print("V2 ZETA gateway flows enabled")
	return nil
}

// IsV2ZETAEnabled checks if V2 ZETA gateway flows are enabled
// Returns false if crosschain flags are not set on the network
func (r *E2ERunner) IsV2ZETAEnabled() bool {
	response, err := r.ObserverClient.CrosschainFlags(r.Ctx, &observertypes.QueryGetCrosschainFlagsRequest{})
	if err != nil {
		return false
	}
	return response.CrosschainFlags.IsV2ZetaEnabled
}
