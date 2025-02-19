package v5

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/emissions/types"
)

type EmissionsKeeper interface {
	SetParams(ctx sdk.Context, params types.Params) error
	GetParams(ctx sdk.Context) (types.Params, bool)
}

// MigrateStore migrates the store from v3 to v4
// The v3 params are copied to the v4 params, and the v4 params are set in the store
// v4 params removes unused parameters from v3; these values are discarded.
// v4 introduces a new parameter, BlockRewardAmount, which is set to the default value
func MigrateStore(
	ctx sdk.Context,
	emissionsKeeper EmissionsKeeper,
) error {
	defaultParams := types.DefaultParams()
	params, found := emissionsKeeper.GetParams(ctx)
	if found {
		params.PendingBallotsBufferBlocks = defaultParams.PendingBallotsBufferBlocks
	}
	return emissionsKeeper.SetParams(ctx, params)
}
