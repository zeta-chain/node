package keeper

import (
	"strconv"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/exp/slices"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// SetGasPrice set a specific gasPrice in the store from its index
func (k Keeper) SetGasPrice(ctx sdk.Context, gasPrice types.GasPrice) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GasPriceKey))
	gasPrice.Index = strconv.FormatInt(gasPrice.ChainId, 10)
	b := k.cdc.MustMarshal(&gasPrice)
	store.Set(types.KeyPrefix(gasPrice.Index), b)
}

// GetGasPrice returns a gasPrice from its index or false if it doesn't exist.
func (k Keeper) GetGasPrice(ctx sdk.Context, chainID int64) (types.GasPrice, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GasPriceKey))

	b := store.Get(types.KeyPrefix(strconv.FormatInt(chainID, 10)))
	if b == nil {
		return types.GasPrice{}, false
	}

	var val types.GasPrice
	k.cdc.MustUnmarshal(b, &val)

	return val, true
}

// GetMedianGasPriceInUint returns the median gas price (and median priority fee) from the store
// or false if it doesn't exist.
func (k Keeper) GetMedianGasPriceInUint(ctx sdk.Context, chainID int64) (math.Uint, math.Uint, bool) {
	entity, found := k.GetGasPrice(ctx, chainID)
	if !found {
		return math.ZeroUint(), math.ZeroUint(), false
	}

	var (
		gasPrice    = math.NewUint(entity.Prices[entity.MedianIndex])
		priorityFee = math.NewUint(medianValue(entity.PriorityFees))
	)

	return gasPrice, priorityFee, true
}

// medianValue returns the median value of a slice
// example: [ 1 7 5 2 3 6 4 ] => [ 1 2 3 4 5 6 7 ] => 4
func medianValue(items []uint64) uint64 {
	switch len(items) {
	case 0:
		return 0
	case 1:
		return items[0]
	}

	// We don't want to modify the original slice
	copiedItems := make([]uint64, len(items))
	copy(copiedItems, items)

	slices.Sort(copiedItems)
	mv := copiedItems[len(copiedItems)/2]

	// We don't need the copy anymore
	//nolint:ineffassign // let's help garbage collector :)
	copiedItems = nil

	return mv
}

// RemoveGasPrice removes a gasPrice from the store
func (k Keeper) RemoveGasPrice(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GasPriceKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllGasPrice returns all gasPrice
func (k Keeper) GetAllGasPrice(ctx sdk.Context) (list []types.GasPrice) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GasPriceKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.GasPrice
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
