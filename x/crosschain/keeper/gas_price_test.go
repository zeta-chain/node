package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// Keeper Tests
func createNGasPrice(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.GasPrice {
	items := make([]types.GasPrice, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].ChainId = int64(i)
		items[i].Index = strconv.FormatInt(int64(i), 10)
		keeper.SetGasPrice(ctx, items[i])
	}
	return items
}

func TestGasPriceGet(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)
	items := createNGasPrice(k, ctx, 10)
	for _, item := range items {
		rst, found := k.GetGasPrice(ctx, item.ChainId)
		require.True(t, found)
		require.Equal(t, item, rst)
	}
}

func TestGasPriceRemove(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)
	items := createNGasPrice(k, ctx, 10)
	for _, item := range items {
		k.RemoveGasPrice(ctx, item.Index)
		_, found := k.GetGasPrice(ctx, item.ChainId)
		require.False(t, found)
	}
}

func TestGasPriceGetAll(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)
	items := createNGasPrice(k, ctx, 10)
	require.Equal(t, items, k.GetAllGasPrice(ctx))
}
