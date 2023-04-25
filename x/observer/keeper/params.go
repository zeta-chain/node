package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

func (k Keeper) SetCoreParamsByChainID(ctx sdk.Context,coreParams []types.CoreParams) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&coreParams)
	key := types.KeyPrefix(fmt.Sprintf("%s", types.AllCoreParams))
	store.Set(key, b)
}

func (k Keeper) GetAllCoreParams(ctx sdk.Context) (val []types.CoreParams) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AllCoreParams))
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(""))
	for ; iterator.Valid(); iterator.Next() {
		var item types.CoreParams
		k.cdc.MustUnmarshal(iterator.Value(), &item)
		val = append(val, item)
	}
	return
}

}

func (k Keeper) GetCoreParamsByChainID(ctx sdk.Context, chainID int64) (val types.CoreParams, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AllCoreParams))
	key := types.KeyPrefix(fmt.Sprintf("%d", chainID))
	b := store.Get(key)
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
