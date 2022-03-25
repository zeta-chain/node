package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

// SetInTx set a specific inTx in the store from its index
func (k Keeper) SetInTx(ctx sdk.Context, inTx types.InTx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InTxKey))
	b := k.cdc.MustMarshalBinaryBare(&inTx)
	store.Set(types.KeyPrefix(inTx.Index), b)
}

// GetInTx returns a inTx from its index
func (k Keeper) GetInTx(ctx sdk.Context, index string) (val types.InTx, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InTxKey))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshalBinaryBare(b, &val)
	return val, true
}

// RemoveInTx removes a inTx from the store
func (k Keeper) RemoveInTx(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InTxKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllInTx returns all inTx
func (k Keeper) GetAllInTx(ctx sdk.Context) (list []types.InTx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InTxKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.InTx
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
