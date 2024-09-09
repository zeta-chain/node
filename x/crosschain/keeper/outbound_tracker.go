package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/crosschain/types"
)

func getOutboundTrackerIndex(chainID int64, nonce uint64) string {
	return fmt.Sprintf("%d-%d", chainID, nonce)
}

// SetOutboundTracker set a specific outbound tracker in the store from its index
func (k Keeper) SetOutboundTracker(ctx sdk.Context, outboundTracker types.OutboundTracker) {
	outboundTracker.Index = getOutboundTrackerIndex(outboundTracker.ChainId, outboundTracker.Nonce)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OutboundTrackerKeyPrefix))
	b := k.cdc.MustMarshal(&outboundTracker)
	store.Set(types.OutboundTrackerKey(
		outboundTracker.Index,
	), b)
}

// GetOutboundTracker returns an outbound tracker from its index
func (k Keeper) GetOutboundTracker(
	ctx sdk.Context,
	chainID int64,
	nonce uint64,
) (val types.OutboundTracker, found bool) {
	index := getOutboundTrackerIndex(chainID, nonce)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OutboundTrackerKeyPrefix))

	b := store.Get(types.OutboundTrackerKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveOutboundTrackerFromStore removes a outbound tracker from the store
func (k Keeper) RemoveOutboundTrackerFromStore(
	ctx sdk.Context,
	chainID int64,
	nonce uint64,
) {
	index := getOutboundTrackerIndex(chainID, nonce)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OutboundTrackerKeyPrefix))
	store.Delete(types.OutboundTrackerKey(
		index,
	))
}

// GetAllOutboundTracker returns all outbound trackers
func (k Keeper) GetAllOutboundTracker(ctx sdk.Context) (list []types.OutboundTracker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OutboundTrackerKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.OutboundTracker
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
