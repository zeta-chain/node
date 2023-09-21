package keeper_test

import (
	"strconv"
	"testing"

	"github.com/zeta-chain/zetacore/x/crosschain/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func createNInTxHashToCctx(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.InTxHashToCctx {
	items := make([]types.InTxHashToCctx, n)
	for i := range items {
		items[i].InTxHash = strconv.Itoa(i)

		keeper.SetInTxHashToCctx(ctx, items[i])
	}
	return items
}

func TestInTxHashToCctxGet(t *testing.T) {
	keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
	items := createNInTxHashToCctx(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetInTxHashToCctx(ctx,
			item.InTxHash,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestInTxHashToCctxRemove(t *testing.T) {
	keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
	items := createNInTxHashToCctx(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveInTxHashToCctx(ctx,
			item.InTxHash,
		)
		_, found := keeper.GetInTxHashToCctx(ctx,
			item.InTxHash,
		)
		require.False(t, found)
	}
}

func TestInTxHashToCctxGetAll(t *testing.T) {
	keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
	items := createNInTxHashToCctx(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllInTxHashToCctx(ctx)),
	)
}
