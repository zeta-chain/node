package keeper

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/crosschain/types"
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

// GetRateLimiterAssetRateList returns a list of all foreign asset rate
func (k Keeper) GetRateLimiterAssetRateList(
	ctx sdk.Context,
) (flags types.RateLimiterFlags, assetRates []types.AssetRate, found bool) {
	flags, found = k.GetRateLimiterFlags(ctx)
	if !found {
		return flags, nil, false
	}

	// loop through the rate limiter flags to get the rate
	for _, conversion := range flags.Conversions {
		fCoin, found := k.fungibleKeeper.GetForeignCoins(ctx, conversion.Zrc20)
		if !found {
			continue
		}

		// add asset rate to list
		assetRates = append(assetRates, types.AssetRate{
			ChainId:  fCoin.ForeignChainId,
			Asset:    strings.ToLower(fCoin.Asset),
			Decimals: fCoin.Decimals,
			CoinType: fCoin.CoinType,
			Rate:     conversion.Rate,
		})
	}
	return flags, assetRates, true
}
