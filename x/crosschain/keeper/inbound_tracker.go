package keeper

import (
	"fmt"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/zeta-chain/node/x/crosschain/types"
)

func getInboundTrackerKey(chainID int64, txHash string) string {
	return fmt.Sprintf("%d-%s", chainID, txHash)
}

// SetInboundTracker set a specific InboundTracker in the store from its index
func (k Keeper) SetInboundTracker(ctx sdk.Context, InboundTracker types.InboundTracker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InboundTrackerKeyPrefix))
	b := k.cdc.MustMarshal(&InboundTracker)
	key := types.KeyPrefix(getInboundTrackerKey(InboundTracker.ChainId, InboundTracker.TxHash))
	store.Set(key, b)
}

// GetInboundTracker returns a InboundTracker from its index
func (k Keeper) GetInboundTracker(
	ctx sdk.Context,
	chainID int64,
	txHash string,
) (val types.InboundTracker, found bool) {
	key := getInboundTrackerKey(chainID, txHash)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InboundTrackerKeyPrefix))
	b := store.Get(types.KeyPrefix(key))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) RemoveInboundTrackerIfExists(ctx sdk.Context, chainID int64, txHash string) {
	key := getInboundTrackerKey(chainID, txHash)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InboundTrackerKeyPrefix))
	if store.Has(types.KeyPrefix(key)) {
		store.Delete(types.KeyPrefix(key))
	}
}

func (k Keeper) GetAllInboundTrackerPaginated(
	ctx sdk.Context,
	pagination *query.PageRequest,
) (inTxTrackers []types.InboundTracker, pageRes *query.PageResponse, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InboundTrackerKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	pageRes, err = query.Paginate(store, pagination, func(_ []byte, value []byte) error {
		var inTxTracker types.InboundTracker
		if err := k.cdc.Unmarshal(value, &inTxTracker); err != nil {
			return err
		}
		inTxTrackers = append(inTxTrackers, inTxTracker)
		return nil
	})
	return
}

func (k Keeper) GetAllInboundTracker(ctx sdk.Context) (inTxTrackers []types.InboundTracker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InboundTrackerKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.InboundTracker
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		inTxTrackers = append(inTxTrackers, val)
	}
	return
}

func (k Keeper) GetAllInboundTrackerForChain(ctx sdk.Context, chainID int64) (list []types.InboundTracker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InboundTrackerKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, fmt.Appendf(nil, "%d-", chainID))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.InboundTracker
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}
	return list
}

func (k Keeper) GetAllInboundTrackerForChainPaginated(
	ctx sdk.Context,
	chainID int64,
	pagination *query.PageRequest,
) (inTxTrackers []types.InboundTracker, pageRes *query.PageResponse, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InboundTrackerKeyPrefix))
	chainStore := prefix.NewStore(store, types.KeyPrefix(fmt.Sprintf("%d-", chainID)))
	pageRes, err = query.Paginate(chainStore, pagination, func(_ []byte, value []byte) error {
		var inTxTracker types.InboundTracker
		if err := k.cdc.Unmarshal(value, &inTxTracker); err != nil {
			return err
		}
		inTxTrackers = append(inTxTrackers, inTxTracker)
		return nil
	})
	return
}
