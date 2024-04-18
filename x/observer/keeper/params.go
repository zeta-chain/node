package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// GetParams get all parameters
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefix(types.ParamsKey))
	if bz == nil {
		return types.Params{}, false
	}
	k.cdc.MustUnmarshal(bz, &params)
	return params, true
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.KeyPrefix(types.ParamsKey), bz)
	return nil
}
