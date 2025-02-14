package v3

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/authority/types"
)

const (
	MsgURLRemoveInboundTracker    = "/zetachain.zetacore.crosschain.MsgRemoveInboundTracker"
	MsgURLDisableFastConfirmation = "/zetachain.zetacore.observer.MsgDisableFastConfirmation"
)

type authorityKeeper interface {
	SetAuthorizationList(ctx sdk.Context, list types.AuthorizationList)
	GetAuthorizationList(ctx sdk.Context) (val types.AuthorizationList, found bool)
}

// MigrateStore migrates the authority module state from the consensus version 2 to 3
func MigrateStore(
	ctx sdk.Context,
	keeper authorityKeeper,
) error {
	var (
		authorizationList          = types.DefaultAuthorizationsList()
		removeInboundAuthorization = types.Authorization{
			MsgUrl:           MsgURLRemoveInboundTracker,
			AuthorizedPolicy: types.PolicyType_groupEmergency,
		}
		disableFastConfirmationAuthorization = types.Authorization{
			MsgUrl:           MsgURLDisableFastConfirmation,
			AuthorizedPolicy: types.PolicyType_groupEmergency,
		}
	)

	// Fetch the current authorization list, if found use that instead of default list
	al, found := keeper.GetAuthorizationList(ctx)
	if found {
		authorizationList = al
	}

	// Add the new authorization
	authorizationList.SetAuthorization(removeInboundAuthorization)
	authorizationList.SetAuthorization(disableFastConfirmationAuthorization)

	// Validate the authorization list
	err := authorizationList.Validate()
	if err != nil {
		return err
	}

	// Set the new authorization list
	keeper.SetAuthorizationList(ctx, authorizationList)
	return nil
}
