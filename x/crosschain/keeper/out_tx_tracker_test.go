package keeper_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// Keeper Tests
func createNOutTxTracker(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.OutTxTracker {
	items := make([]types.OutTxTracker, n)
	for i := range items {
		items[i].ChainId = int64(i)
		items[i].Nonce = uint64(i)
		items[i].Index = fmt.Sprintf("%d-%d", items[i].ChainId, items[i].Nonce)

		keeper.SetOutTxTracker(ctx, items[i])
	}
	return items
}

func TestOutTxTrackerGet(t *testing.T) {
	keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
	items := createNOutTxTracker(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetOutTxTracker(ctx,
			item.ChainId,
			item.Nonce,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestOutTxTrackerRemove(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)
	items := createNOutTxTracker(k, ctx, 10)
	for _, item := range items {
		k.RemoveOutTxTracker(ctx,
			item.ChainId,
			item.Nonce,
		)
		_, found := k.GetOutTxTracker(ctx,
			item.ChainId,
			item.Nonce,
		)
		require.False(t, found)
	}
}

func TestOutTxTrackerGetAll(t *testing.T) {
	keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
	items := createNOutTxTracker(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllOutTxTracker(ctx)),
	)
}
