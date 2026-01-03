package v6

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/authority/types"
)

type authorityKeeper interface {
	SetAuthorizationList(ctx sdk.Context, list types.AuthorizationList)
	GetAuthorizationList(ctx sdk.Context) (val types.AuthorizationList, found bool)
}

// MigrateStore migrates the authority module state from the consensus version 4 to 5
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
		whitelistAuthorization = types.Authorization{
			MsgUrl:           "/zetachain.zetacore.crosschain.MsgWhitelistAsset",
			AuthorizedPolicy: types.PolicyType_groupAdmin,
		}
	)

	al, found := keeper.GetAuthorizationList(ctx)
	if found {
		authorizationList = al
	}

	authorizationList.SetAuthorization(removeObserverAuthorization)
	authorizationList.SetAuthorization(whitelistAuthorization)

	// Validate the authorization list
	err := authorizationList.Validate()
	if err != nil {
		return err
	}
	keeper.SetAuthorizationList(ctx, authorizationList)
	return nil
}
