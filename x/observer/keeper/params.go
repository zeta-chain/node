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

func (k Keeper) SetClientParamsByChainID(ctx sdk.Context, chainID int64, clientParams types.ClientParams) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClientParamsByChainID))
	b := k.cdc.MustMarshal(&clientParams)
	key := types.KeyPrefix(fmt.Sprintf("%d", chainID))
	fmt.Println("Setting ChainID", chainID, "to", key)
	store.Set(key, b)
}

func (k Keeper) GetClientParamsByChainID(ctx sdk.Context, chainID int64) (val types.ClientParams, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ClientParamsByChainID))
	key := types.KeyPrefix(fmt.Sprintf("%d", chainID))
	fmt.Println("Getting ChainID", chainID, "    ", key)
	b := store.Get(key)
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
