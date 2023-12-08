package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// ChainNonces methods
// The object stores the current nonce for the chain

// SetChainNonces set a specific chainNonces in the store from its index
func (k Keeper) SetChainNonces(ctx sdk.Context, chainNonces types.ChainNonces) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ChainNoncesKey))
	b := k.cdc.MustMarshal(&chainNonces)
	store.Set(types.KeyPrefix(chainNonces.Index), b)
}

// GetChainNonces returns a chainNonces from its index
func (k Keeper) GetChainNonces(ctx sdk.Context, index string) (val types.ChainNonces, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ChainNoncesKey))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveChainNonces removes a chainNonces from the store
func (k Keeper) RemoveChainNonces(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ChainNoncesKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllChainNonces returns all chainNonces
func (k Keeper) GetAllChainNonces(ctx sdk.Context) (list []types.ChainNonces) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ChainNoncesKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ChainNonces
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// NonceToCctx methods
// The object stores the mapping from nonce to cross chain tx

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

func (k Keeper) RemoveNonceToCctx(ctx sdk.Context, nonceToCctx types.NonceToCctx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.NonceToCctxKeyPrefix))
	store.Delete(types.KeyPrefix(fmt.Sprintf("%s-%d-%d", nonceToCctx.Tss, nonceToCctx.ChainId, nonceToCctx.Nonce)))
}

// PendingNonces methods
// The object stores the pending nonces for the chain

func (k Keeper) SetPendingNonces(ctx sdk.Context, pendingNonces types.PendingNonces) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PendingNoncesKeyPrefix))
	b := k.cdc.MustMarshal(&pendingNonces)
	store.Set(types.KeyPrefix(fmt.Sprintf("%s-%d", pendingNonces.Tss, pendingNonces.ChainId)), b)
}

func (k Keeper) GetPendingNonces(ctx sdk.Context, tss string, chainID int64) (val types.PendingNonces, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PendingNoncesKeyPrefix))

	b := store.Get(types.KeyPrefix(fmt.Sprintf("%s-%d", tss, chainID)))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetAllPendingNonces(ctx sdk.Context) (list []*types.PendingNonces) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PendingNoncesKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.PendingNonces
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, &val)
	}

	return
}

func (k Keeper) RemovePendingNonces(ctx sdk.Context, pendingNonces types.PendingNonces) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PendingNoncesKeyPrefix))
	store.Delete(types.KeyPrefix(fmt.Sprintf("%s-%d", pendingNonces.Tss, pendingNonces.ChainId)))
}

// Utility functions

func (k Keeper) RemoveFromPendingNonces(ctx sdk.Context, tssPubkey string, chainID int64, nonce int64) {
	p, found := k.GetPendingNonces(ctx, tssPubkey, chainID)
	if found && nonce >= p.NonceLow && nonce <= p.NonceHigh {
		p.NonceLow = nonce + 1
		k.SetPendingNonces(ctx, p)
	}
}

func (k Keeper) SetTssAndUpdateNonce(ctx sdk.Context, tss observerTypes.TSS) {
	k.zetaObserverKeeper.SetTSS(ctx, tss)
	// initialize the nonces and pending nonces of all enabled chains
	supportedChains := k.zetaObserverKeeper.GetParams(ctx).GetSupportedChains()
	for _, chain := range supportedChains {
		chainNonce := types.ChainNonces{
			Index:   chain.ChainName.String(),
			ChainId: chain.ChainId,
			Nonce:   0,
			// #nosec G701 always positive
			FinalizedHeight: uint64(ctx.BlockHeight()),
		}
		k.SetChainNonces(ctx, chainNonce)

		p := types.PendingNonces{
			NonceLow:  0,
			NonceHigh: 0,
			ChainId:   chain.ChainId,
			Tss:       tss.TssPubkey,
		}
		k.SetPendingNonces(ctx, p)
	}
}
