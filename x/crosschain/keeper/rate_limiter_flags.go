package keeper

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/pkg/coin"
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

// GetRateLimiterAssetRateList returns a list of all foreign asset rate
func (k Keeper) GetRateLimiterAssetRateList(ctx sdk.Context) []*types.AssetRate {
	rateLimiterFlags, _ := k.GetRateLimiterFlags(ctx)

	// the result list
	assetRateList := make([]*types.AssetRate, 0)

	// loop through the rate limiter flags to get the rate
	for _, conversion := range rateLimiterFlags.Conversions {
		fCoin, found := k.fungibleKeeper.GetForeignCoins(ctx, conversion.Zrc20)
		if !found {
			continue
		}

		// add asset rate to list
		assetRateList = append(assetRateList, &types.AssetRate{
			ChainId:  fCoin.ForeignChainId,
			Asset:    strings.ToLower(fCoin.Asset),
			Decimals: fCoin.Decimals,
			CoinType: fCoin.CoinType,
			Rate:     conversion.Rate,
		})
	}
	return assetRateList
}

// BuildAssetRateMapFromList builds maps (foreign chain id -> asset -> rate) from a list of gas and erc20 asset rates
// The 1st map: foreign chain id -> gas coin asset rate
// The 2nd map: foreign chain id -> erc20 asset -> erc20 coin asset rate
func BuildAssetRateMapFromList(assetRates []*types.AssetRate) (map[int64]*types.AssetRate, map[int64]map[string]*types.AssetRate) {
	// the result maps
	gasAssetRateMap := make(map[int64]*types.AssetRate)
	erc20AssetRateMap := make(map[int64]map[string]*types.AssetRate)

	// loop through the asset rates to build the maps
	for _, assetRate := range assetRates {
		chainID := assetRate.ChainId
		if assetRate.CoinType == coin.CoinType_Gas {
			gasAssetRateMap[chainID] = assetRate
		} else {
			if _, found := erc20AssetRateMap[chainID]; !found {
				erc20AssetRateMap[chainID] = make(map[string]*types.AssetRate)
			}
			erc20AssetRateMap[chainID][strings.ToLower(assetRate.Asset)] = assetRate
		}
	}
	return gasAssetRateMap, erc20AssetRateMap
}
