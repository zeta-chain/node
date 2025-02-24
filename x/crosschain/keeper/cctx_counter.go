package keeper

import (
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/crosschain/types"
)

func (k Keeper) getCounterValueStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.CounterValueKey))
}

func (k Keeper) getCounterIndexStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CounterIndexKey))
}

// GetCctxCounter retrieves the current counter value
//
// Will return 0 if the counter has never been set
func (k Keeper) GetCctxCounter(ctx sdk.Context) uint64 {
	store := k.getCounterValueStore(ctx)
	storedCounter := store.Get([]byte(types.CounterValueKey))

	return sdk.BigEndianToUint64(storedCounter)
}

// SetCctxCounter updates the current counter value
//
// This should only be used in setCctxCounterIndex and in state import
func (k Keeper) SetCctxCounter(ctx sdk.Context, val uint64) {
	store := k.getCounterValueStore(ctx)
	store.Set([]byte(types.CounterValueKey), sdk.Uint64ToBigEndian(val))
}

// getNextCctxCounter retrieves and increments the counter for ordering
func (k Keeper) getNextCctxCounter(ctx sdk.Context) uint64 {
	storedCounter := k.GetCctxCounter(ctx)
	nextCounter := storedCounter + 1
	k.SetCctxCounter(ctx, nextCounter)
	return nextCounter
}

// setCctxCounterIndex sets a new CCTX in the counter index
//
// note that we use the raw bytes in the index rather than the hex encoded bytes
// like in the main store
func (k Keeper) setCctxCounterIndex(ctx sdk.Context, cctx types.CrossChainTx) {
	counterIndexStore := k.getCounterIndexStore(ctx)
	nextCounter := k.getNextCctxCounter(ctx)

	cctxIndex, err := cctx.GetCCTXIndexBytes()
	if err != nil {
		k.Logger(ctx).Error("get cctx index bytes", "err", err)
		return
	}

	// must use big endian so most significant bytes are first for sortability
	nextCounterBytes := sdk.Uint64ToBigEndian(nextCounter)
	counterIndexStore.Set(nextCounterBytes, cctxIndex[:])
}
