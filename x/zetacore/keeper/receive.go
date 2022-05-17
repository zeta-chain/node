package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

// SetReceive set a specific receive in the store from its index
func (k Keeper) SetReceive(ctx sdk.Context, receive types.Receive) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ReceiveKey))
	b := k.cdc.MustMarshal(&receive)
	store.Set(types.KeyPrefix(receive.Index), b)
}

// GetReceive returns a receive from its index
func (k Keeper) GetReceive(ctx sdk.Context, index string) (val types.Receive, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ReceiveKey))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveReceive removes a receive from the store
func (k Keeper) RemoveReceive(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ReceiveKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllReceive returns all receive
func (k Keeper) GetAllReceive(ctx sdk.Context) (list []types.Receive) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ReceiveKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Receive
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
