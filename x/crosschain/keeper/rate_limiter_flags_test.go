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
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

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

func TestKeeper_GetRateLimiterRates(t *testing.T) {
	k, ctx, _, zk := keepertest.CrosschainKeeper(t)

	// create test flags
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

	chainID := chains.GoerliLocalnetChain().ChainId

	// add gas coin
	fcGas := sample.ForeignCoins(t, zrc20GasAddr)
	fcGas.CoinType = coin.CoinType_Gas
	fcGas.ForeignChainId = chainID
	zk.FungibleKeeper.SetForeignCoins(ctx, fcGas)

	// add two erc20 coins
	asset1 := sample.EthAddress().Hex()
	fcERC20 := sample.ForeignCoins(t, zrc20ERC20Addr1)
	fcERC20.Asset = asset1
	fcERC20.ForeignChainId = chainID
	zk.FungibleKeeper.SetForeignCoins(ctx, fcERC20)

	asset2 := sample.EthAddress().Hex()
	fcERC20 = sample.ForeignCoins(t, zrc20ERC20Addr2)
	fcERC20.Asset = asset2
	fcERC20.ForeignChainId = chainID
	zk.FungibleKeeper.SetForeignCoins(ctx, fcERC20)

	// set flags
	k.SetRateLimiterFlags(ctx, flags)
	r, found := k.GetRateLimiterFlags(ctx)
	require.True(t, found)
	require.Equal(t, flags, r)

	// get rates
	gasRates, erc20Rates := k.GetRateLimiterRates(ctx)
	require.Equal(t, 1, len(gasRates))
	require.Equal(t, 1, len(erc20Rates))
	require.Equal(t, sdk.NewDec(1), gasRates[chainID])
	require.Equal(t, 2, len(erc20Rates[chainID]))
	require.Equal(t, sdk.NewDec(2), erc20Rates[chainID][strings.ToLower(asset1)])
	require.Equal(t, sdk.NewDec(3), erc20Rates[chainID][strings.ToLower(asset2)])
}
