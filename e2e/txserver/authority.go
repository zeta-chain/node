package txserver

import (
	"fmt"

	e2eutils "github.com/zeta-chain/node/e2e/utils"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
)

// AddAuthorization adds a new authorization in the authority module for admin message
func (zts ZetaTxServer) AddAuthorization(msgURL string) error {
	// retrieve account
	accAdmin, err := zts.clientCtx.Keyring.Key(e2eutils.AdminPolicyName)
	if err != nil {
		return err
	}
	addrAdmin, err := accAdmin.GetAddress()
	if err != nil {
		return err
	}

	// add new authorization
	_, err = zts.BroadcastTx(e2eutils.AdminPolicyName, authoritytypes.NewMsgAddAuthorization(
		addrAdmin.String(),
		msgURL,
		authoritytypes.PolicyType_groupAdmin,
	))
	if err != nil {
		return fmt.Errorf("failed to add authorization: %w", err)
	}

	return nil
}
