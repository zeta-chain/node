package keeper

import (
	"github.com/Meta-Protocol/metacore/x/metacore/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetTxoutConfirmation set a specific txoutConfirmation in the store from its index
func (k Keeper) SetTxoutConfirmation(ctx sdk.Context, txoutConfirmation types.TxoutConfirmation) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TxoutConfirmationKey))
	b := k.cdc.MustMarshalBinaryBare(&txoutConfirmation)
	store.Set(types.KeyPrefix(txoutConfirmation.Index), b)
}

// GetTxoutConfirmation returns a txoutConfirmation from its index
func (k Keeper) GetTxoutConfirmation(ctx sdk.Context, index string) (val types.TxoutConfirmation, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TxoutConfirmationKey))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshalBinaryBare(b, &val)
	return val, true
}

// RemoveTxoutConfirmation removes a txoutConfirmation from the store
func (k Keeper) RemoveTxoutConfirmation(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TxoutConfirmationKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllTxoutConfirmation returns all txoutConfirmation
func (k Keeper) GetAllTxoutConfirmation(ctx sdk.Context) (list []types.TxoutConfirmation) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TxoutConfirmationKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.TxoutConfirmation
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
