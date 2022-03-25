package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

// SetTSSVoter set a specific tSSVoter in the store from its index
func (k Keeper) SetTSSVoter(ctx sdk.Context, tSSVoter types.TSSVoter) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TSSVoterKey))
	b := k.cdc.MustMarshalBinaryBare(&tSSVoter)
	store.Set(types.KeyPrefix(tSSVoter.Index), b)
}

// GetTSSVoter returns a tSSVoter from its index
func (k Keeper) GetTSSVoter(ctx sdk.Context, index string) (val types.TSSVoter, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TSSVoterKey))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshalBinaryBare(b, &val)
	return val, true
}

// RemoveTSSVoter removes a tSSVoter from the store
func (k Keeper) RemoveTSSVoter(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TSSVoterKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllTSSVoter returns all tSSVoter
func (k Keeper) GetAllTSSVoter(ctx sdk.Context) (list []types.TSSVoter) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TSSVoterKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.TSSVoter
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
