package keeper

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/emissions/types"
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

func (k Keeper) AddRewards(ctx sdk.Context, address string, amount sdkmath.Int) {
	we, found := k.GetWithdrawableEmission(ctx, address)
	if !found {
		we = types.WithdrawableEmissions{
			Address: address,
			Amount:  amount,
		}
	} else {
		we.Amount = we.Amount.Add(amount)
	}
	k.SetWithdrawableEmission(ctx, we)
}

// SlashRewards slashes the rewards of a given address, if the address has no rewards left, it will set the rewards to 0.
/* This function is a basic implementation of slashing, it will be improved in the future, based on further discussions.
Improvements will include:
- Add a jailing mechanism
- Slash observer below 0 , or remove from observer list if their rewards are below 0
*/

func (k Keeper) SlashRewards(ctx sdk.Context, address string, amount sdkmath.Int) {
	we, found := k.GetWithdrawableEmission(ctx, address)
	if !found {
		we = types.WithdrawableEmissions{
			Address: address,
			Amount:  sdk.ZeroInt(),
		}
	} else {
		slashedRewards := we.Amount.Sub(amount)
		if slashedRewards.IsNegative() {
			we.Amount = sdkmath.ZeroInt()
		}
		we.Amount = we.Amount.Sub(amount)
	}
	k.SetWithdrawableEmission(ctx, we)
}
