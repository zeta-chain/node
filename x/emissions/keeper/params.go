package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/emissions/types"
)

// GetParams get all parameters
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefix(types.ParamsKey))
	if bz == nil {
		return types.Params{}, false
	}
	err := k.cdc.Unmarshal(bz, &params)
	if err != nil {
		return types.Params{}, false
	}

	return params, true
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}

	store.Set(types.KeyPrefix(types.ParamsKey), bz)
	return nil
}
