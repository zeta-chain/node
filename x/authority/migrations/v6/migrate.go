package v6

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/authority/types"
)

type authorityKeeper interface {
	SetAuthorizationList(ctx sdk.Context, list types.AuthorizationList)
	GetAuthorizationList(ctx sdk.Context) (val types.AuthorizationList, found bool)
}

// MigrateStore migrates the authority module state from the consensus version 5 to 6
// It ensures that the newly added message MsgRemoveObserver is authorized under the Admin policy.
func MigrateStore(
	ctx sdk.Context,
	keeper authorityKeeper,
) error {
	var (
		authorizationList           = types.DefaultAuthorizationsList()
		removeObserverAuthorization = types.Authorization{
			MsgUrl:           "/zetachain.zetacore.observer.MsgRemoveObserver",
			AuthorizedPolicy: types.PolicyType_groupAdmin,
		}
	)

	al, found := keeper.GetAuthorizationList(ctx)
	if found {
		authorizationList = al
	}

	authorizationList.SetAuthorization(removeObserverAuthorization)
	if err := authorizationList.Validate(); err != nil {
		return err
	}
	keeper.SetAuthorizationList(ctx, authorizationList)
	return nil
}
