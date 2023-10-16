package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return
}

func (k Keeper) GetParamsIsExists(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSetIfExists(ctx, &params)
	return
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

func (k Keeper) SetCoreParams(ctx sdk.Context, coreParams common.CoreParamsList) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&coreParams)
	key := types.KeyPrefix(fmt.Sprintf("%s", types.AllCoreParams))
	store.Set(key, b)
}

func (k Keeper) GetAllCoreParams(ctx sdk.Context) (val common.CoreParamsList, found bool) {
	found = false
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.KeyPrefix(fmt.Sprintf("%s", types.AllCoreParams)))
	if b == nil {
		return
	}
	found = true
	k.cdc.MustUnmarshal(b, &val)
	return
}

func (k Keeper) GetCoreParamsByChainID(ctx sdk.Context, chainID int64) (*common.CoreParams, bool) {
	allCoreParams, found := k.GetAllCoreParams(ctx)
	if !found {
		return &common.CoreParams{}, false
	}
	for _, coreParams := range allCoreParams.CoreParams {
		if coreParams.ChainId == chainID {
			return coreParams, true
		}
	}
	return &common.CoreParams{}, false
}
