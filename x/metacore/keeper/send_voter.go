package keeper

import (
	"github.com/Meta-Protocol/metacore/x/metacore/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetSendVoter set a specific sendVoter in the store from its index
func (k Keeper) SetSendVoter(ctx sdk.Context, sendVoter types.SendVoter) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SendVoterKey))
	b := k.cdc.MustMarshalBinaryBare(&sendVoter)
	store.Set(types.KeyPrefix(sendVoter.Index), b)
}

// GetSendVoter returns a sendVoter from its index
func (k Keeper) GetSendVoter(ctx sdk.Context, index string) (val types.SendVoter, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SendVoterKey))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshalBinaryBare(b, &val)
	return val, true
}

// RemoveSendVoter removes a sendVoter from the store
func (k Keeper) RemoveSendVoter(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SendVoterKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllSendVoter returns all sendVoter
func (k Keeper) GetAllSendVoter(ctx sdk.Context) (list []types.SendVoter) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SendVoterKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.SendVoter
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
