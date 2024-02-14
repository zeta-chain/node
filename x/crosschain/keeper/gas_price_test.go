package keeper

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// Keeper Tests
func createNGasPrice(keeper *Keeper, ctx sdk.Context, n int) []types.GasPrice {
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
	keeper, ctx := setupKeeper(t)
	items := createNGasPrice(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetGasPrice(ctx, item.ChainId)
		require.True(t, found)
		require.Equal(t, item, rst)
	}
}
func TestGasPriceRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNGasPrice(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveGasPrice(ctx, item.Index)
		_, found := keeper.GetGasPrice(ctx, item.ChainId)
		require.False(t, found)
	}
}

func TestGasPriceGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNGasPrice(keeper, ctx, 10)
	require.Equal(t, items, keeper.GetAllGasPrice(ctx))
}
