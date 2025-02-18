package v3

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/authority/types"
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
	//
	var (
		authorizationList            = types.DefaultAuthorizationsList()
		updateZRC20NameAuthorization = types.Authorization{
			MsgUrl:           "/zetachain.zetacore.fungible.MsgUpdateZRC20Name",
			AuthorizedPolicy: types.PolicyType_groupAdmin,
		}
		removeInboundAuthorization = types.Authorization{
			MsgUrl:           "/zetachain.zetacore.crosschain.MsgRemoveInboundTracker",
			AuthorizedPolicy: types.PolicyType_groupEmergency,
		}
		updateOperationalChainParamsAuthorization = types.Authorization{
			MsgUrl:           "/zetachain.zetacore.observer.MsgUpdateOperationalChainParams",
			AuthorizedPolicy: types.PolicyType_groupOperational,
		}
		updateChainParamsAuthorization = types.Authorization{
			MsgUrl:           "/zetachain.zetacore.observer.MsgUpdateChainParams",
			AuthorizedPolicy: types.PolicyType_groupAdmin,
		}
	)

	// Fetch the current authorization list, if found use that instead of default list
	al, found := keeper.GetAuthorizationList(ctx)
	if found {
		authorizationList = al
	}

	// Add the new authorization
	authorizationList.SetAuthorization(updateZRC20NameAuthorization)
	authorizationList.SetAuthorization(removeInboundAuthorization)
	authorizationList.SetAuthorization(updateOperationalChainParamsAuthorization)
	authorizationList.SetAuthorization(updateChainParamsAuthorization)

	// Validate the authorization list
	err := authorizationList.Validate()
	if err != nil {
		return err
	}

	// Set the new authorization list
	keeper.SetAuthorizationList(ctx, authorizationList)
	return nil
}
