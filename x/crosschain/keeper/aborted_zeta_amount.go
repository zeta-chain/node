package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func (k Keeper) SetAbortedZetaAmount(ctx sdk.Context, abortedZetaAmount types.AbortedZetaAmount) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&abortedZetaAmount)
	store.Set([]byte(types.AbortedZetaAmountKey), b)
}

func (k Keeper) GetAbortedZetaAmount(ctx sdk.Context) (val types.AbortedZetaAmount, found bool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.AbortedZetaAmountKey))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) AddAbortedZetaAmount(ctx sdk.Context, amount sdkmath.Uint) {
	abortedZetaAmount, found := k.GetAbortedZetaAmount(ctx)
	if !found {
		abortedZetaAmount = types.AbortedZetaAmount{
			Amount: amount,
		}
	} else {
		abortedZetaAmount.Amount = abortedZetaAmount.Amount.Add(amount)
	}
	k.SetAbortedZetaAmount(ctx, abortedZetaAmount)
}
