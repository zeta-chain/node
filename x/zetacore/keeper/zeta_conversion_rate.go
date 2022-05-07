package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

// SetZetaConversionRate set a specific zetaConversionRate in the store from its index
func (k Keeper) SetZetaConversionRate(ctx sdk.Context, zetaConversionRate types.ZetaConversionRate) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ZetaConversionRateKeyPrefix))
	b := k.cdc.MustMarshal(&zetaConversionRate)
	store.Set(types.ZetaConversionRateKey(
		zetaConversionRate.Index,
	), b)
}

// GetZetaConversionRate returns a zetaConversionRate from its index
func (k Keeper) GetZetaConversionRate(
	ctx sdk.Context,
	index string,

) (val types.ZetaConversionRate, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ZetaConversionRateKeyPrefix))

	b := store.Get(types.ZetaConversionRateKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveZetaConversionRate removes a zetaConversionRate from the store
func (k Keeper) RemoveZetaConversionRate(
	ctx sdk.Context,
	index string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ZetaConversionRateKeyPrefix))
	store.Delete(types.ZetaConversionRateKey(
		index,
	))
}

// GetAllZetaConversionRate returns all zetaConversionRate
func (k Keeper) GetAllZetaConversionRate(ctx sdk.Context) (list []types.ZetaConversionRate) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ZetaConversionRateKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ZetaConversionRate
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
