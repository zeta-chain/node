package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/zeta-chain/zetacore/x/observer/types"
)

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

func (k Keeper) GetAllPendingNoncesPaginated(
	ctx sdk.Context,
	pagination *query.PageRequest,
) (list []types.PendingNonces, pageRes *query.PageResponse, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PendingNoncesKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()
	pageRes, err = query.Paginate(store, pagination, func(_ []byte, value []byte) error {
		var val types.PendingNonces
		if err := k.cdc.Unmarshal(value, &val); err != nil {
			return err
		}
		list = append(list, val)
		return nil
	})

	return
}

func (k Keeper) GetAllPendingNonces(ctx sdk.Context) (list []types.PendingNonces, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PendingNoncesKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.PendingNonces
		err := k.cdc.Unmarshal(iterator.Value(), &val)
		if err != nil {
			return nil, err
		}
		list = append(list, val)
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

func (k Keeper) SetTssAndUpdateNonce(ctx sdk.Context, tss types.TSS) {
	k.SetTSS(ctx, tss)
	// initialize the nonces and pending nonces of all enabled chains
	supportedChains := k.GetSupportedChains(ctx)
	for _, chain := range supportedChains {
		chainNonce := types.ChainNonces{
			Index:   chain.ChainName.String(),
			ChainId: chain.ChainId,
			Nonce:   0,
			// #nosec G115 always positive
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
