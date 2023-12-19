package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// SetLastBlockHeight set a specific lastBlockHeight in the store from its index
func (k Keeper) SetLastBlockHeight(ctx sdk.Context, lastBlockHeight types.LastBlockHeight) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LastBlockHeightKey))
	b := k.cdc.MustMarshal(&lastBlockHeight)
	store.Set(types.KeyPrefix(lastBlockHeight.Index), b)
}

// GetLastBlockHeight returns a lastBlockHeight from its index
func (k Keeper) GetLastBlockHeight(ctx sdk.Context, index string) (val types.LastBlockHeight, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LastBlockHeightKey))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveLastBlockHeight removes a lastBlockHeight from the store
func (k Keeper) RemoveLastBlockHeight(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LastBlockHeightKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllLastBlockHeight returns all lastBlockHeight
func (k Keeper) GetAllLastBlockHeight(ctx sdk.Context) (list []types.LastBlockHeight) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LastBlockHeightKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.LastBlockHeight
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
