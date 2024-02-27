package keeper

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

func (k Keeper) SetWithdrawEmissions(ctx sdk.Context, we types.WithdrawEmission) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WithdrawEmissionsKey))
	b := k.cdc.MustMarshal(&we)
	store.Set([]byte(we.Address), b)
}

func (k Keeper) GetWithdrawEmissions(ctx sdk.Context, address string) (val types.WithdrawEmission, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WithdrawEmissionsKey))
	b := store.Get(types.KeyPrefix(address))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) DeleteWithdrawEmissions(ctx sdk.Context, address string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WithdrawEmissionsKey))
	store.Delete([]byte(address))
}

func (k Keeper) GetAllWithdrawEmissions(ctx sdk.Context) (list []types.WithdrawEmission) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WithdrawEmissionsKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.WithdrawEmission
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}
	return
}

func (k Keeper) CreateWithdrawEmissions(ctx sdk.Context, address string, amount sdkmath.Int) error {
	emissions, found := k.GetWithdrawableEmission(ctx, address)
	if !found {
		return types.ErrEmissionsNotFound
	}
	if amount.IsNegative() || amount.IsZero() {
		return types.ErrNotEnoughEmissionsAvailable
	}
	if amount.GT(emissions.Amount) {
		amount = emissions.Amount
	}
	k.SetWithdrawEmissions(ctx, types.WithdrawEmission{
		Address: address,
		Amount:  amount,
	})
	return nil
}
