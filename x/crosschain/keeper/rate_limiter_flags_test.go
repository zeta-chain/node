package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
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
) (fungibletypes.ForeignCoins, types.AssetRate) {
	// create foreign coin
	foreignCoin := sample.ForeignCoins(t, zrc20Addr)
	foreignCoin.Asset = asset
	foreignCoin.ForeignChainId = chainID
	foreignCoin.Decimals = decimals
	foreignCoin.CoinType = coinType

	// create corresponding asset rate
	assetRate := sample.CustomAssetRate(
		foreignCoin.ForeignChainId,
		foreignCoin.Asset,
		foreignCoin.Decimals,
		foreignCoin.CoinType,
		rate,
	)

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

func TestKeeper_GetRateLimiterAssetRateList(t *testing.T) {
	k, ctx, _, zk := keepertest.CrosschainKeeper(t)

	// create test flags
	chainID := chains.GoerliLocalnet.ChainId
	zrc20GasAddr := sample.EthAddress().Hex()
	zrc20ERC20Addr1 := sample.EthAddress().Hex()
	zrc20ERC20Addr2 := sample.EthAddress().Hex()
	testflags := types.RateLimiterFlags{
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

	// asset rates not found before setting flags
	flags, assetRates, found := k.GetRateLimiterAssetRateList(ctx)
	require.False(t, found)
	require.Equal(t, types.RateLimiterFlags{}, flags)
	require.Nil(t, assetRates)

	// set flags
	k.SetRateLimiterFlags(ctx, testflags)

	// add gas coin
	gasCoin, gasAssetRate := createForeignCoinAndAssetRate(
		t,
		zrc20GasAddr,
		"",
		chainID,
		18,
		coin.CoinType_Gas,
		sdk.NewDec(1),
	)
	zk.FungibleKeeper.SetForeignCoins(ctx, gasCoin)

	// add 1st erc20 coin
	erc20Coin1, erc20AssetRate1 := createForeignCoinAndAssetRate(
		t,
		zrc20ERC20Addr1,
		sample.EthAddress().Hex(),
		chainID,
		8,
		coin.CoinType_ERC20,
		sdk.NewDec(2),
	)
	zk.FungibleKeeper.SetForeignCoins(ctx, erc20Coin1)

	// add 2nd erc20 coin
	erc20Coin2, erc20AssetRate2 := createForeignCoinAndAssetRate(
		t,
		zrc20ERC20Addr2,
		sample.EthAddress().Hex(),
		chainID,
		6,
		coin.CoinType_ERC20,
		sdk.NewDec(3),
	)
	zk.FungibleKeeper.SetForeignCoins(ctx, erc20Coin2)

	// get rates
	flags, assetRates, found = k.GetRateLimiterAssetRateList(ctx)
	require.True(t, found)
	require.Equal(t, testflags, flags)
	require.EqualValues(t, []types.AssetRate{gasAssetRate, erc20AssetRate1, erc20AssetRate2}, assetRates)
}
