package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// CRUD
func (k Keeper) SetPendingTxQueue(ctx sdk.Context, pendingTxQueue types.PendingTxQueue) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PendingTxQueueKeyPrefix))
	b := k.cdc.MustMarshal(&pendingTxQueue)
	store.Set(types.KeyPrefix(fmt.Sprintf("%d", pendingTxQueue.ChainId)), b)
}

// GetKeygen returns keygen
func (k Keeper) GetPendingTxQueue(ctx sdk.Context, chainId uint64) (val types.PendingTxQueue, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PendingTxQueueKeyPrefix))

	b := store.Get(types.KeyPrefix(fmt.Sprintf("%d", chainId)))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveKeygen removes keygen from the store
func (k Keeper) RemovePendingTxQueue(ctx sdk.Context, chainId uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PendingTxQueueKeyPrefix))
	store.Delete(types.KeyPrefix(fmt.Sprintf("%d", chainId)))
}

// GetAllChainNonces returns all chainNonces
func (k Keeper) GetAllPendingTxQueues(ctx sdk.Context) (list []*types.PendingTxQueue) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PendingTxQueueKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PendingTxQueue
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, &val)
	}

	return
}

// PendingTx
func (k Keeper) SetPendingTx(ctx sdk.Context, pendingTx types.PendingTx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PendingTxKeyPrefix))
	b := k.cdc.MustMarshal(&pendingTx)
	store.Set(types.KeyPrefix(pendingTx.CctxIndex), b)
}

func (k Keeper) GetPendingTx(ctx sdk.Context, cctxIndex string) (val types.PendingTx, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PendingTxKeyPrefix))

	b := store.Get(types.KeyPrefix(cctxIndex))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
