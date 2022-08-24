package keeper_test

import (
	"fmt"
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

func createNOutTxTracker(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.OutTxTracker {
	items := make([]types.OutTxTracker, n)
	for i := range items {
		items[i].Index = fmt.Sprintf("testchain-%d", i)

		keeper.SetOutTxTracker(ctx, items[i])
	}
	return items
}

func TestOutTxTrackerGet(t *testing.T) {
	keeper, ctx := keepertest.ZetacoreKeeper(t)
	items := createNOutTxTracker(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetOutTxTracker(ctx,
			item.Index,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestOutTxTrackerRemove(t *testing.T) {
	keeper, ctx := keepertest.ZetacoreKeeper(t)
	items := createNOutTxTracker(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveOutTxTracker(ctx,
			item.Index,
		)
		_, found := keeper.GetOutTxTracker(ctx,
			item.Index,
		)
		require.False(t, found)
	}
}

func TestOutTxTrackerGetAll(t *testing.T) {
	keeper, ctx := keepertest.ZetacoreKeeper(t)
	items := createNOutTxTracker(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllOutTxTracker(ctx)),
	)
}
