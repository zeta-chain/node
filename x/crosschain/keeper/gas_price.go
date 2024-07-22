package keeper

import (
	"strconv"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	slicemath "github.com/zeta-chain/zetacore/pkg/math"
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
		priorityFee = math.NewUint(slicemath.SliceMedianValue(entity.PriorityFees, false))
	)

	return gasPrice, priorityFee, true
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
