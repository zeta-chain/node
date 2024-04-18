package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// SetRateLimiterFlags set the rate limiter flags in the store
func (k Keeper) SetRateLimiterFlags(ctx sdk.Context, rateLimiterFlags types.RateLimiterFlags) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RateLimiterFlagsKey))
	b := k.cdc.MustMarshal(&rateLimiterFlags)
	store.Set([]byte{0}, b)
}

// GetRateLimiterFlags returns the rate limiter flags
func (k Keeper) GetRateLimiterFlags(ctx sdk.Context) (val types.RateLimiterFlags, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RateLimiterFlagsKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
