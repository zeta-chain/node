package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// ChainNonces methods
// The object stores the current nonce for the chain

// SetChainNonces set a specific chainNonces in the store from its index
func (k Keeper) SetChainNonces(ctx sdk.Context, chainNonces types.ChainNonces) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ChainNoncesKey))
	b := k.cdc.MustMarshal(&chainNonces)
	store.Set(types.KeyPrefix(chainNonces.Index), b)
}

// GetChainNonces returns a chainNonces from its index
func (k Keeper) GetChainNonces(ctx sdk.Context, index string) (val types.ChainNonces, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ChainNoncesKey))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveChainNonces removes a chainNonces from the store
func (k Keeper) RemoveChainNonces(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ChainNoncesKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllChainNonces returns all chainNonces
func (k Keeper) GetAllChainNonces(ctx sdk.Context) (list []types.ChainNonces) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ChainNoncesKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.ChainNonces
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
