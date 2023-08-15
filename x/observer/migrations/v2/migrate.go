package v2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// MigrateStore migrates the x/observer module state from the consensus version 1 to 2
/* This migration adds a
- new permission flag to the observer module called IsOutboundEnabled
- a new policy Policy_Type_add_observer
*/
func MigrateStore(
	ctx sdk.Context,
	observerKeeper keeper.Keeper,
) error {

	observerKeeper.SetPermissionFlags(ctx, types.PermissionFlags{
		IsInboundEnabled:  true,
		IsOutboundEnabled: true,
	})
	params := observerKeeper.GetParams(ctx)
	params.AdminPolicy = types.DefaultAdminPolicy()
	observerKeeper.SetParams(ctx, params)

	return nil
}
