package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

// GetParamSetIfExists get all parameters as types.Params if they exist
// non existent parameters will return zero values
func (k Keeper) GetParamSetIfExists(ctx sdk.Context) (params types.Params) {
	k.paramStore.GetParamSetIfExists(ctx, &params)
	return
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramStore.SetParamSet(ctx, &params)
}
