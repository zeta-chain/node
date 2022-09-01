package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

// SetSend set a specific send in the store from its index
func (k Keeper) SetSend(ctx sdk.Context, send types.Send) {
	p := types.KeyPrefix(fmt.Sprintf("%s-%d", types.SendKey, send.Status))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
	b := k.cdc.MustMarshal(&send)
	store.Set(types.KeyPrefix(send.Index), b)
}

// GetSend returns a send from its index
func (k Keeper) GetSend(ctx sdk.Context, index string, status types.SendStatus) (val types.Send, found bool) {
	p := types.KeyPrefix(fmt.Sprintf("%s-%d", types.SendKey, status))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetSendMultipleStatus(ctx sdk.Context, index string, status []types.SendStatus) (val types.Send, found bool) {
	for _, s := range status {
		p := types.KeyPrefix(fmt.Sprintf("%s-%d", types.SendKey, s))
		store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
		send := store.Get(types.KeyPrefix(index))
		if send != nil {
			k.cdc.MustUnmarshal(send, &val)
			return val, true
		}
	}
	return val, false
}

// RemoveSend removes a send from the store
func (k Keeper) RemoveSend(ctx sdk.Context, index string, status types.SendStatus) {
	p := types.KeyPrefix(fmt.Sprintf("%s-%d", types.SendKey, status))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
	store.Delete(types.KeyPrefix(index))
}

// GetAllSend returns all send
func (k Keeper) GetAllSend(ctx sdk.Context) []types.Send {
	var list []types.Send
	for i := 0; i <= 6; i++ {
		list = append(list, k.GetAllSendWithStatus(ctx, types.SendStatus(i))...)
	}
	return list
}

func (k Keeper) GetAllSendWithStatus(ctx sdk.Context, status types.SendStatus) (list []types.Send) {
	p := types.KeyPrefix(fmt.Sprintf("%s-%d", types.SendKey, status))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.Send
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}
	return
}

func (k Keeper) GetAllPendingOutBoundSend(ctx sdk.Context) (list []*types.Send) {
	p := types.KeyPrefix(fmt.Sprintf("%s-%d", types.SendKey, types.SendStatus_PendingOutbound))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.Send
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, &val)
	}
	return
}

func (k Keeper) GetAllPendingInBoundSend(ctx sdk.Context) (list []*types.Send) {
	p := types.KeyPrefix(fmt.Sprintf("%s-%d", types.SendKey, types.SendStatus_PendingInbound))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.Send
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, &val)
	}
	return
}
