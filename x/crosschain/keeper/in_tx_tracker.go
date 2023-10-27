package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func getInTrackerKey(chainID int64, txHash string) string {
	return fmt.Sprintf("%d-%s", chainID, txHash)
}

// SetInTxTracker set a specific InTxTracker in the store from its index
func (k Keeper) SetInTxTracker(ctx sdk.Context, InTxTracker types.InTxTracker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InTxTrackerKeyPrefix))
	b := k.cdc.MustMarshal(&InTxTracker)
	key := types.KeyPrefix(getInTrackerKey(InTxTracker.ChainId, InTxTracker.TxHash))
	store.Set(key, b)
}

// GetInTxTracker returns a InTxTracker from its index
func (k Keeper) GetInTxTracker(ctx sdk.Context, chainID int64, txHash string) (val types.InTxTracker, found bool) {
	key := getInTrackerKey(chainID, txHash)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InTxTrackerKeyPrefix))
	b := store.Get(types.KeyPrefix(key))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) RemoveInTxTrackerIfExists(ctx sdk.Context, chainID int64, txHash string) {
	key := getInTrackerKey(chainID, txHash)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InTxTrackerKeyPrefix))
	if store.Has(types.KeyPrefix(key)) {
		store.Delete(types.KeyPrefix(key))
	}
}
func (k Keeper) GetAllInTxTracker(ctx sdk.Context) (list []types.InTxTracker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InTxTrackerKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.InTxTracker
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}
	return list
}

func (k Keeper) GetAllInTxTrackerForChain(ctx sdk.Context, chainID int64) (list []types.InTxTracker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InTxTrackerKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte(fmt.Sprintf("%d-", chainID)))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.InTxTracker
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}
	return list
}

func (k Keeper) GetAllInTxTrackerForChainPaginated(ctx sdk.Context, chainID int64, pagination *query.PageRequest) (inTxTrackers []types.InTxTracker, pageRes *query.PageResponse, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(fmt.Sprintf("%s", types.InTxTrackerKeyPrefix)))
	chainStore := prefix.NewStore(store, types.KeyPrefix(fmt.Sprintf("%d-", chainID)))
	pageRes, err = query.Paginate(chainStore, pagination, func(key []byte, value []byte) error {
		var inTxTracker types.InTxTracker
		if err := k.cdc.Unmarshal(value, &inTxTracker); err != nil {
			return err
		}
		inTxTrackers = append(inTxTrackers, inTxTracker)
		return nil
	})
	return
}
