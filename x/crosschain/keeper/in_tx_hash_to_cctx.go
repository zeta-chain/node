package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// SetInTxHashToCctx set a specific inTxHashToCctx in the store from its index
func (k Keeper) SetInTxHashToCctx(ctx sdk.Context, inTxHashToCctx types.InTxHashToCctx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InTxHashToCctxKeyPrefix))
	b := k.cdc.MustMarshal(&inTxHashToCctx)
	store.Set(types.InTxHashToCctxKey(
		inTxHashToCctx.InTxHash,
	), b)
}

// GetInTxHashToCctx returns a inTxHashToCctx from its index
func (k Keeper) GetInTxHashToCctx(
	ctx sdk.Context,
	inTxHash string,

) (val types.InTxHashToCctx, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InTxHashToCctxKeyPrefix))

	b := store.Get(types.InTxHashToCctxKey(
		inTxHash,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveInTxHashToCctx removes a inTxHashToCctx from the store
func (k Keeper) RemoveInTxHashToCctx(
	ctx sdk.Context,
	inTxHash string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InTxHashToCctxKeyPrefix))
	store.Delete(types.InTxHashToCctxKey(
		inTxHash,
	))
}

// GetAllInTxHashToCctx returns all inTxHashToCctx
func (k Keeper) GetAllInTxHashToCctx(ctx sdk.Context) (list []types.InTxHashToCctx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InTxHashToCctxKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.InTxHashToCctx
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
