package keeper

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/emissions/types"
)

func (k Keeper) SetWithdrawableEmission(ctx sdk.Context, we types.WithdrawableEmissions) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WithdrawableEmissionsKey))
	b := k.cdc.MustMarshal(&we)
	store.Set([]byte(we.Address), b)
}

func (k Keeper) GetWithdrawableEmission(ctx sdk.Context, address string) (val types.WithdrawableEmissions, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WithdrawableEmissionsKey))
	b := store.Get(types.KeyPrefix(address))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetAllWithdrawableEmission(ctx sdk.Context) (list []types.WithdrawableEmissions) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WithdrawableEmissionsKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.WithdrawableEmissions
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}
	return
}

// AddObserverEmission adds the given amount to the withdrawable emission of a given address.
// If the address does not have a withdrawable emission, it will create a new withdrawable emission with the given amount.
func (k Keeper) AddObserverEmission(ctx sdk.Context, address string, amount sdkmath.Int) {
	we, found := k.GetWithdrawableEmission(ctx, address)
	if !found {
		we = types.WithdrawableEmissions{Address: address, Amount: sdkmath.ZeroInt()}
	}
	we.Amount = we.Amount.Add(amount)
	k.SetWithdrawableEmission(ctx, we)
}

// RemoveWithdrawableEmission removes the given amount from the withdrawable emission of a given address.
// If the amount is greater than the available withdrawable emissionsf or that address it will remove the entire amount from the withdrawable emissions.
// If the amount is negative or zero, it will return an error.
func (k Keeper) RemoveWithdrawableEmission(ctx sdk.Context, address string, amount sdkmath.Int) error {
	we, found := k.GetWithdrawableEmission(ctx, address)
	if !found {
		return types.ErrEmissionsNotFound
	}
	if amount.IsNegative() || amount.IsZero() {
		return types.ErrInvalidAmount.Wrap("amount to be removed is negative or zero")
	}
	if amount.GT(we.Amount) {
		return types.ErrInvalidAmount.Wrap("amount to be removed is greater than the available withdrawable emission")
	}
	we.Amount = we.Amount.Sub(amount)
	k.SetWithdrawableEmission(ctx, we)
	return nil
}

// SlashObserverEmission slashes the rewards of a given address, if the address has no rewards left, it will set the rewards to 0.
// If the address does not have a withdrawable emission, it will create a new withdrawable emission with zero amount.
/* This function is a basic implementation of slashing; it will be improved in the future .
Improvements will include:
- Add a jailing mechanism
- Slash observer below 0, or remove from an observer list if their rewards are below 0
*/
// https://github.com/zeta-chain/node/issues/945
func (k Keeper) SlashObserverEmission(ctx sdk.Context, address string, slashAmount sdkmath.Int) {
	we, found := k.GetWithdrawableEmission(ctx, address)
	if !found {
		we = types.WithdrawableEmissions{Address: address, Amount: sdkmath.ZeroInt()}
	} else {
		we.Amount = we.Amount.Sub(slashAmount)
		if we.Amount.IsNegative() {
			we.Amount = sdkmath.ZeroInt()
		}
	}
	k.SetWithdrawableEmission(ctx, we)
}
