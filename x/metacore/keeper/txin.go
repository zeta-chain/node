package keeper

import (
	"github.com/Meta-Protocol/metacore/x/metacore/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetTxin set a specific txin in the store from its index
func (k Keeper) SetTxin(ctx sdk.Context, txin types.Txin) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TxinKey))
	b := k.cdc.MustMarshalBinaryBare(&txin)
	store.Set(types.KeyPrefix(txin.Index), b)
}

// GetTxin returns a txin from its index
func (k Keeper) GetTxin(ctx sdk.Context, index string) (val types.Txin, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TxinKey))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshalBinaryBare(b, &val)
	return val, true
}

// RemoveTxin removes a txin from the store
func (k Keeper) RemoveTxin(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TxinKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllTxin returns all txin
func (k Keeper) GetAllTxin(ctx sdk.Context) (list []types.Txin) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TxinKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Txin
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
