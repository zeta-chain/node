package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// SetInboundHashToCctx set a specific inboundHashToCctx in the store from its index
func (k Keeper) SetInboundHashToCctx(ctx sdk.Context, inboundHashToCctx types.InboundHashToCctx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InboundHashToCctxKeyPrefix))
	b := k.cdc.MustMarshal(&inboundHashToCctx)
	store.Set(types.InboundHashToCctxKey(
		inboundHashToCctx.InboundHash,
	), b)
}

// GetInboundHashToCctx returns a inboundHashToCctx from its index
func (k Keeper) GetInboundHashToCctx(
	ctx sdk.Context,
	inboundHash string,

) (val types.InboundHashToCctx, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InboundHashToCctxKeyPrefix))

	b := store.Get(types.InboundHashToCctxKey(
		inboundHash,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveInboundHashToCctx removes a inboundHashToCctx from the store
func (k Keeper) RemoveInboundHashToCctx(
	ctx sdk.Context,
	inboundHash string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InboundHashToCctxKeyPrefix))
	store.Delete(types.InboundHashToCctxKey(
		inboundHash,
	))
}

// GetAllInboundHashToCctx returns all inboundHashToCctx
func (k Keeper) GetAllInboundHashToCctx(ctx sdk.Context) (list []types.InboundHashToCctx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InboundHashToCctxKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.InboundHashToCctx
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
