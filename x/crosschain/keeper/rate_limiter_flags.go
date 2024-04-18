package keeper

import (
	"strings"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// hardcoded rate limiter flags
var rateLimitFlags = types.RateLimiterFlags{
	Enabled: true,
	Window:  100,                   // 100 zeta blocks, 10 minutes
	Rate:    math.NewUint(2000000), // 2000 ZETA
	Conversions: []types.Conversion{
		// ETH
		{
			Zrc20: "0x13A0c5930C028511Dc02665E7285134B6d11A5f4",
			Rate:  sdk.NewDec(2500),
		},
		// USDT
		{
			Zrc20: "0xbD1e64A22B9F92D9Ce81aA9B4b0fFacd80215564",
			Rate:  sdk.MustNewDecFromStr("0.8"),
			//sdk.NewDec(0.8),
		},
		// BTC
		{
			Zrc20: "0x8f56682c2b8b2e3d4f6f7f7d6f3c01b3f6f6a7d6",
			Rate:  sdk.NewDec(50000),
		},
	},
}

// SetRatelimiterFlags set the rate limiter flags in the store
func (k Keeper) SetRatelimiterFlags(ctx sdk.Context, crosschainFlags types.RateLimiterFlags) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RateLimiterFlagsKey))
	b := k.cdc.MustMarshal(&crosschainFlags)
	store.Set([]byte{0}, b)
}

// GetRatelimiterFlags read the rate limiter flags from the store
func (k Keeper) GetRatelimiterFlags(_ sdk.Context) (val types.RateLimiterFlags, found bool) {
	return rateLimitFlags, true
	// store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RateLimiterFlagsKey))

	// b := store.Get([]byte{0})
	// if b == nil {
	// 	return val, false
	// }

	// k.cdc.MustUnmarshal(b, &val)
	// return val, true
}

// GetRatelimiterRates returns two maps of foreign coins and their rates
// The 1st map: foreign chain id -> gas coin rate
// The 2nd map: foreign chain id -> erc20 asset -> erc20 coin rate
func (k Keeper) GetRatelimiterRates(ctx sdk.Context) (map[int64]sdk.Dec, map[int64]map[string]sdk.Dec) {
	rateLimitFlags, _ := k.GetRatelimiterFlags(ctx)

	// the result maps
	gasCoinRates := make(map[int64]sdk.Dec)
	erc20CoinRates := make(map[int64]map[string]sdk.Dec)

	// loop through the rate limiter flags to get the rate
	for _, conversion := range rateLimitFlags.Conversions {
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
