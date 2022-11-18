package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

// SetForeignCoins set a specific foreignCoins in the store from its index
func (k Keeper) SetForeignCoins(ctx sdk.Context, foreignCoins types.ForeignCoins) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ForeignCoinsKeyPrefix))
	b := k.cdc.MustMarshal(&foreignCoins)
	store.Set(types.ForeignCoinsKey(
		foreignCoins.Index,
	), b)
}

// GetForeignCoins returns a foreignCoins from its index
func (k Keeper) GetForeignCoins(
	ctx sdk.Context,
	index string,

) (val types.ForeignCoins, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ForeignCoinsKeyPrefix))

	b := store.Get(types.ForeignCoinsKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveForeignCoins removes a foreignCoins from the store
func (k Keeper) RemoveForeignCoins(
	ctx sdk.Context,
	index string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ForeignCoinsKeyPrefix))
	store.Delete(types.ForeignCoinsKey(
		index,
	))
}

// GetAllForeignCoins returns all foreignCoins
func (k Keeper) GetAllForeignCoins(ctx sdk.Context) (list []types.ForeignCoins) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ForeignCoinsKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ForeignCoins
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
