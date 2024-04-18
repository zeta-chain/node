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

// GetRateLimiterRates returns two maps of foreign coins and their rates
// The 1st map: foreign chain id -> gas coin rate
// The 2nd map: foreign chain id -> erc20 asset -> erc20 coin rate
func (k Keeper) GetRateLimiterRates(ctx sdk.Context) (map[int64]sdk.Dec, map[int64]map[string]sdk.Dec) {
	rateLimiterFlags, _ := k.GetRateLimiterFlags(ctx)

	// the result maps
	gasCoinRates := make(map[int64]sdk.Dec)
	erc20CoinRates := make(map[int64]map[string]sdk.Dec)

	// loop through the rate limiter flags to get the rate
	for _, conversion := range rateLimiterFlags.Conversions {
		fCoin, found := k.fungibleKeeper.GetForeignCoins(ctx, conversion.Zrc20)
		if !found {
			continue
		}

		chainID := fCoin.ForeignChainId
		switch fCoin.CoinType {
		case coin.CoinType_Gas:
			gasCoinRates[chainID] = conversion.Rate
		case coin.CoinType_ERC20:
			if _, found := erc20CoinRates[chainID]; !found {
				erc20CoinRates[chainID] = make(map[string]sdk.Dec)
			}
			erc20CoinRates[chainID][strings.ToLower(fCoin.Asset)] = conversion.Rate
		}
	}
	return gasCoinRates, erc20CoinRates
}
