package keeper_test

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

// createForeignCoinAndAssetRate creates foreign coin and corresponding asset rate
func createForeignCoinAndAssetRate(
	t *testing.T,
	zrc20Addr string,
	asset string,
	chainID int64,
	decimals uint32,
	coinType coin.CoinType,
	rate sdk.Dec,
) (fungibletypes.ForeignCoins, *types.AssetRate) {
	// create foreign coin
	foreignCoin := sample.ForeignCoins(t, zrc20Addr)
	foreignCoin.Asset = asset
	foreignCoin.ForeignChainId = chainID
	foreignCoin.Decimals = decimals
	foreignCoin.CoinType = coinType

	// create corresponding asset rate
	assetRate := &types.AssetRate{
		ChainId:  foreignCoin.ForeignChainId,
		Asset:    strings.ToLower(foreignCoin.Asset),
		Decimals: foreignCoin.Decimals,
		CoinType: foreignCoin.CoinType,
		Rate:     rate,
	}

	return foreignCoin, assetRate
}

func TestKeeper_GetRateLimiterFlags(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)

	// not found
	_, found := k.GetRateLimiterFlags(ctx)
	require.False(t, found)

	flags := sample.RateLimiterFlags()

	k.SetRateLimiterFlags(ctx, flags)
	r, found := k.GetRateLimiterFlags(ctx)
	require.True(t, found)
	require.Equal(t, flags, r)
}

func TestKeeper_GetRateLimiterRateList(t *testing.T) {
	k, ctx, _, zk := keepertest.CrosschainKeeper(t)

	// create test flags
	chainID := chains.GoerliLocalnetChain.ChainId
	zrc20GasAddr := sample.EthAddress().Hex()
	zrc20ERC20Addr1 := sample.EthAddress().Hex()
	zrc20ERC20Addr2 := sample.EthAddress().Hex()
	flags := types.RateLimiterFlags{
		Rate: sdk.NewUint(100),
		Conversions: []types.Conversion{
			{
				Zrc20: zrc20GasAddr,
				Rate:  sdk.NewDec(1),
			},
			{
				Zrc20: zrc20ERC20Addr1,
				Rate:  sdk.NewDec(2),
			},
			{
				Zrc20: zrc20ERC20Addr2,
				Rate:  sdk.NewDec(3),
			},
		},
	}

	// set flags
	k.SetRateLimiterFlags(ctx, flags)

	// add gas coin
	gasCoin, gasAssetRate := createForeignCoinAndAssetRate(t, zrc20GasAddr, "", chainID, 18, coin.CoinType_Gas, sdk.NewDec(1))
	zk.FungibleKeeper.SetForeignCoins(ctx, gasCoin)

	// add 1st erc20 coin
	erc20Coin1, erc20AssetRate1 := createForeignCoinAndAssetRate(t, zrc20ERC20Addr1, sample.EthAddress().Hex(), chainID, 8, coin.CoinType_ERC20, sdk.NewDec(2))
	zk.FungibleKeeper.SetForeignCoins(ctx, erc20Coin1)

	// add 2nd erc20 coin
	erc20Coin2, erc20AssetRate2 := createForeignCoinAndAssetRate(t, zrc20ERC20Addr2, sample.EthAddress().Hex(), chainID, 6, coin.CoinType_ERC20, sdk.NewDec(3))
	zk.FungibleKeeper.SetForeignCoins(ctx, erc20Coin2)

	// get rates
	assetRates := k.GetRateLimiterAssetRateList(ctx)
	require.EqualValues(t, []*types.AssetRate{gasAssetRate, erc20AssetRate1, erc20AssetRate2}, assetRates)
}

func TestBuildAssetRateMapFromList(t *testing.T) {
	// define asset rate list
	assetRates := []*types.AssetRate{
		{
			ChainId:  1,
			Asset:    "eth",
			Decimals: 18,
			CoinType: coin.CoinType_Gas,
			Rate:     sdk.NewDec(1),
		},
		{
			ChainId:  1,
			Asset:    "usdt",
			Decimals: 6,
			CoinType: coin.CoinType_ERC20,
			Rate:     sdk.NewDec(2),
		},
		{
			ChainId:  2,
			Asset:    "btc",
			Decimals: 8,
			CoinType: coin.CoinType_Gas,
			Rate:     sdk.NewDec(3),
		},
	}

	// build asset rate map
	gasAssetRateMap, erc20AssetRateMap := keeper.BuildAssetRateMapFromList(assetRates)

	// check length
	require.Equal(t, 2, len(gasAssetRateMap))
	require.Equal(t, 1, len(erc20AssetRateMap))
	require.Equal(t, 1, len(erc20AssetRateMap[1]))

	// check gas asset rate map
	require.EqualValues(t, assetRates[0], gasAssetRateMap[1])
	require.EqualValues(t, assetRates[2], gasAssetRateMap[2])

	// check erc20 asset rate map
	require.EqualValues(t, assetRates[1], erc20AssetRateMap[1]["usdt"])
}
