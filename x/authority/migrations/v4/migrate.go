package v4

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/authority/types"
)

type authorityKeeper interface {
	SetAuthorizationList(ctx sdk.Context, list types.AuthorizationList)
	GetAuthorizationList(ctx sdk.Context) (val types.AuthorizationList, found bool)
}

// MigrateStore migrates the authority module state from the consensus version 3 to 4
func MigrateStore(
	ctx sdk.Context,
	keeper authorityKeeper,
) error {
	var (
		authorizationList                    = types.DefaultAuthorizationsList()
		burnFungibleModuleAssetAuthorization = types.Authorization{
			MsgUrl:           "/zetachain.zetacore.fungible.MsgBurnFungibleModuleAsset",
			AuthorizedPolicy: types.PolicyType_groupAdmin,
		}
	)

	// Fetch the current authorization list, if found use that instead of default list
	al, found := keeper.GetAuthorizationList(ctx)
	if found {
		authorizationList = al
	}

	// Add the new authorization
	authorizationList.SetAuthorization(burnFungibleModuleAssetAuthorization)

	// Validate the authorization list
	err := authorizationList.Validate()
	if err != nil {
		return err
	}

	// Set the new authorization list
	keeper.SetAuthorizationList(ctx, authorizationList)
	return nil
}
