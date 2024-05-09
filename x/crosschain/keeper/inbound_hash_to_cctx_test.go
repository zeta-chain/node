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

func createNInboundHashToCctx(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.InboundHashToCctx {
	items := make([]types.InboundHashToCctx, n)
	for i := range items {
		items[i].InboundHash = strconv.Itoa(i)

		keeper.SetInboundHashToCctx(ctx, items[i])
	}
	return items
}

func TestInTxHashToCctxGet(t *testing.T) {
	keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
	items := createNInboundHashToCctx(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetInboundHashToCctx(ctx,
			item.InboundHash,
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
	items := createNInboundHashToCctx(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveInboundHashToCctx(ctx,
			item.InboundHash,
		)
		_, found := keeper.GetInboundHashToCctx(ctx,
			item.InboundHash,
		)
		require.False(t, found)
	}
}

func TestInTxHashToCctxGetAll(t *testing.T) {
	keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
	items := createNInboundHashToCctx(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllInboundHashToCctx(ctx)),
	)
}
