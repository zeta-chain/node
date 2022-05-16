package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/zetacore/keeper"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNZetaConversionRate(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.ZetaConversionRate {
	items := make([]types.ZetaConversionRate, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)

		keeper.SetZetaConversionRate(ctx, items[i])
	}
	return items
}

func TestZetaConversionRateGet(t *testing.T) {
	keeper, ctx := keepertest.ZetacoreKeeper(t)
	items := createNZetaConversionRate(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetZetaConversionRate(ctx,
			item.Index,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestZetaConversionRateRemove(t *testing.T) {
	keeper, ctx := keepertest.ZetacoreKeeper(t)
	items := createNZetaConversionRate(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveZetaConversionRate(ctx,
			item.Index,
		)
		_, found := keeper.GetZetaConversionRate(ctx,
			item.Index,
		)
		require.False(t, found)
	}
}

func TestZetaConversionRateGetAll(t *testing.T) {
	keeper, ctx := keepertest.ZetacoreKeeper(t)
	items := createNZetaConversionRate(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllZetaConversionRate(ctx)),
	)
}
