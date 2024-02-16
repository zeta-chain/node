package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func (k Keeper) SetZetaAccounting(ctx sdk.Context, abortedZetaAmount types.ZetaAccounting) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&abortedZetaAmount)
	store.Set([]byte(types.ZetaAccountingKey), b)
}

func (k Keeper) GetZetaAccounting(ctx sdk.Context) (val types.ZetaAccounting, found bool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.ZetaAccountingKey))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) AddZetaAbortedAmount(ctx sdk.Context, amount sdkmath.Uint) {
	zetaAccounting, found := k.GetZetaAccounting(ctx)
	if !found {
		zetaAccounting = types.ZetaAccounting{
			AbortedZetaAmount: amount,
		}
	} else {
		zetaAccounting.AbortedZetaAmount = zetaAccounting.AbortedZetaAmount.Add(amount)
	}
	k.SetZetaAccounting(ctx, zetaAccounting)
}

func (k Keeper) RemoveZetaAbortedAmount(ctx sdk.Context, amount sdkmath.Uint) error {
	zetaAccounting, found := k.GetZetaAccounting(ctx)
	if !found {
		return types.ErrUnableToFindZetaAccounting
	}
	if zetaAccounting.AbortedZetaAmount.LT(amount) {
		return types.ErrInsufficientZetaAmount
	}
	zetaAccounting.AbortedZetaAmount = zetaAccounting.AbortedZetaAmount.Sub(amount)
	k.SetZetaAccounting(ctx, zetaAccounting)
	return nil
}
