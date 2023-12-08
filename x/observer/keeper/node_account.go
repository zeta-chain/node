package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// SetNodeAccount set a specific nodeAccount in the store from its index
func (k Keeper) SetNodeAccount(ctx sdk.Context, nodeAccount types.NodeAccount) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NodeAccountKey))
	b := k.cdc.MustMarshal(&nodeAccount)
	store.Set(types.KeyPrefix(nodeAccount.Operator), b)
}

// GetNodeAccount returns a nodeAccount from its index
func (k Keeper) GetNodeAccount(ctx sdk.Context, index string) (val types.NodeAccount, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NodeAccountKey))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveNodeAccount removes a nodeAccount from the store
func (k Keeper) RemoveNodeAccount(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NodeAccountKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllNodeAccount returns all nodeAccount
func (k Keeper) GetAllNodeAccount(ctx sdk.Context) (list []types.NodeAccount) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NodeAccountKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.NodeAccount
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
