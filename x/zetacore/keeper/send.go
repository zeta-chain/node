package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

// SetSend set a specific send in the store from its index
func (k Keeper) SetSend(ctx sdk.Context, send types.Send) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SendKey))
	b := k.cdc.MustMarshalBinaryBare(&send)
	store.Set(types.KeyPrefix(send.Index), b)
}

// GetSend returns a send from its index
func (k Keeper) GetSend(ctx sdk.Context, index string) (val types.Send, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SendKey))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshalBinaryBare(b, &val)
	return val, true
}

// RemoveSend removes a send from the store
func (k Keeper) RemoveSend(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SendKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllSend returns all send
func (k Keeper) GetAllSend(ctx sdk.Context) (list []types.Send) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SendKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Send
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
