package keeper

import (
	"github.com/Meta-Protocol/metacore/x/metacore/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetTxinVoter set a specific txinVoter in the store from its index
func (k Keeper) SetTxinVoter(ctx sdk.Context, txinVoter types.TxinVoter) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TxinVoterKey))
	b := k.cdc.MustMarshalBinaryBare(&txinVoter)
	store.Set(types.KeyPrefix(txinVoter.Index), b)
}

// GetTxinVoter returns a txinVoter from its index
func (k Keeper) GetTxinVoter(ctx sdk.Context, index string) (val types.TxinVoter, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TxinVoterKey))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshalBinaryBare(b, &val)
	return val, true
}

// RemoveTxinVoter removes a txinVoter from the store
func (k Keeper) RemoveTxinVoter(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TxinVoterKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllTxinVoter returns all txinVoter
func (k Keeper) GetAllTxinVoter(ctx sdk.Context) (list []types.TxinVoter) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TxinVoterKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.TxinVoter
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
