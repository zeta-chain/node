package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// CRUD
func (k Keeper) SetNonceToCctx(ctx sdk.Context, nonceToCctx types.NonceToCctx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NonceToCctxKeyPrefix))
	b := k.cdc.MustMarshal(&nonceToCctx)
	store.Set(types.KeyPrefix(fmt.Sprintf("%s-%d-%d", nonceToCctx.Tss, nonceToCctx.ChainId, nonceToCctx.Nonce)), b)
}

func (k Keeper) GetNonceToCctx(ctx sdk.Context, tss string, chainId int64, nonce int64) (val types.NonceToCctx, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NonceToCctxKeyPrefix))

	b := store.Get(types.KeyPrefix(fmt.Sprintf("%s-%d-%d", tss, chainId, nonce)))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) RemoveNonceToCctx(ctx sdk.Context, nonceToCctx types.NonceToCctx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NonceToCctxKeyPrefix))
	store.Delete(types.KeyPrefix(fmt.Sprintf("%s-%d-%d", nonceToCctx.Tss, nonceToCctx.ChainId, nonceToCctx.Nonce)))
}

func (k Keeper) SetPendingNonces(ctx sdk.Context, pendingNonces types.PendingNonces) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PendingNoncesKeyPrefix))
	b := k.cdc.MustMarshal(&pendingNonces)
	store.Set(types.KeyPrefix(fmt.Sprintf("%s-%d", pendingNonces.Tss, pendingNonces.ChainId)), b)
}

func (k Keeper) GetPendingNonces(ctx sdk.Context, tss string, chainId uint64) (val types.PendingNonces, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PendingNoncesKeyPrefix))

	b := store.Get(types.KeyPrefix(fmt.Sprintf("%s-%d", tss, chainId)))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) RemovePendingNonces(ctx sdk.Context, pendingNonces types.PendingNonces) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PendingNoncesKeyPrefix))
	store.Delete(types.KeyPrefix(fmt.Sprintf("%s-%d", pendingNonces.Tss, pendingNonces.ChainId)))
}
