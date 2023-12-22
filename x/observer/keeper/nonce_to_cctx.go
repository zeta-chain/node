package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// NonceToCctx methods
// The object stores the mapping from nonce to cross chain tx

func (k Keeper) RemoveNonceToCctx(ctx sdk.Context, nonceToCctx types.NonceToCctx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NonceToCctxKeyPrefix))
	store.Delete(types.KeyPrefix(fmt.Sprintf("%s-%d-%d", nonceToCctx.Tss, nonceToCctx.ChainId, nonceToCctx.Nonce)))
}

func (k Keeper) SetNonceToCctx(ctx sdk.Context, nonceToCctx types.NonceToCctx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NonceToCctxKeyPrefix))
	b := k.cdc.MustMarshal(&nonceToCctx)
	store.Set(types.KeyPrefix(fmt.Sprintf("%s-%d-%d", nonceToCctx.Tss, nonceToCctx.ChainId, nonceToCctx.Nonce)), b)
}

func (k Keeper) GetNonceToCctx(ctx sdk.Context, tss string, chainID int64, nonce int64) (val types.NonceToCctx, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NonceToCctxKeyPrefix))

	b := store.Get(types.KeyPrefix(fmt.Sprintf("%s-%d-%d", tss, chainID, nonce)))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetAllNonceToCctx(ctx sdk.Context) (list []types.NonceToCctx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NonceToCctxKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.NonceToCctx
		err := k.cdc.Unmarshal(iterator.Value(), &val)
		if err == nil {
			list = append(list, val)
		}

	}

	return

}
