package v2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/authority/types"
)

type authorityKeeper interface {
	SetAuthorizationList(ctx sdk.Context, list types.AuthorizationList)
}

// MigrateStore migrates the authority module state from the consensus version 1 to 2
func MigrateStore(
	ctx sdk.Context,
	keeper authorityKeeper,
) error {
	// It is okay to not validate here, as the authorization list is fixed and will not change
	keeper.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())
	return nil
}
