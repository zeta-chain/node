package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func getOutTrackerIndex(chainID int64, nonce uint64) string {
	return fmt.Sprintf("%d-%d", chainID, nonce)
}

// SetOutTxTracker set a specific outTxTracker in the store from its index
func (k Keeper) SetOutTxTracker(ctx sdk.Context, outTxTracker types.OutTxTracker) {
	outTxTracker.Index = getOutTrackerIndex(outTxTracker.ChainId, outTxTracker.Nonce)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OutTxTrackerKeyPrefix))
	b := k.cdc.MustMarshal(&outTxTracker)
	store.Set(types.OutTxTrackerKey(
		outTxTracker.Index,
	), b)
}

// GetOutTxTracker returns a outTxTracker from its index
func (k Keeper) GetOutTxTracker(
	ctx sdk.Context,
	chainID int64,
	nonce uint64,

) (val types.OutTxTracker, found bool) {
	index := getOutTrackerIndex(chainID, nonce)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OutTxTrackerKeyPrefix))

	b := store.Get(types.OutTxTrackerKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveOutTxTracker removes a outTxTracker from the store
func (k Keeper) RemoveOutTxTracker(
	ctx sdk.Context,
	chainID int64,
	nonce uint64,

) {
	index := getOutTrackerIndex(chainID, nonce)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OutTxTrackerKeyPrefix))
	store.Delete(types.OutTxTrackerKey(
		index,
	))
}

// GetAllOutTxTracker returns all outTxTracker
func (k Keeper) GetAllOutTxTracker(ctx sdk.Context) (list []types.OutTxTracker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OutTxTrackerKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.OutTxTracker
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
