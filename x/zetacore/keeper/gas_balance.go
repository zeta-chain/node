package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

// SetGasBalance set a specific gasBalance in the store from its index
func (k Keeper) SetGasBalance(ctx sdk.Context, gasBalance types.GasBalance) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GasBalanceKey))
	b := k.cdc.MustMarshalBinaryBare(&gasBalance)
	store.Set(types.KeyPrefix(gasBalance.Index), b)
}

// GetGasBalance returns a gasBalance from its index
func (k Keeper) GetGasBalance(ctx sdk.Context, index string) (val types.GasBalance, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GasBalanceKey))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshalBinaryBare(b, &val)
	return val, true
}

// RemoveGasBalance removes a gasBalance from the store
func (k Keeper) RemoveGasBalance(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GasBalanceKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllGasBalance returns all gasBalance
func (k Keeper) GetAllGasBalance(ctx sdk.Context) (list []types.GasBalance) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GasBalanceKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.GasBalance
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
