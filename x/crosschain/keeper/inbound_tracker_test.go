package keeper_test

import (
	"fmt"
	"sort"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/coin"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func createNInboundTracker(keeper *keeper.Keeper, ctx sdk.Context, n int, chainID int64) []types.InboundTracker {
	items := make([]types.InboundTracker, n)
	for i := range items {
		items[i].TxHash = fmt.Sprintf("TxHash-%d", i)
		items[i].ChainId = chainID
		items[i].CoinType = coin.CoinType_Gas
		keeper.SetInboundTracker(ctx, items[i])
	}
	return items
}
func TestKeeper_GetAllInboundTrackerForChain(t *testing.T) {
	t.Run("Get Inbound trackers one by one", func(t *testing.T) {
		keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
		intxTrackers := createNInboundTracker(keeper, ctx, 10, 5)
		for _, item := range intxTrackers {
			rst, found := keeper.GetInboundTracker(ctx, item.ChainId, item.TxHash)
			require.True(t, found)
			require.Equal(t, item, rst)
		}
	})
	t.Run("Get all Inbound trackers", func(t *testing.T) {
		keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
		intxTrackers := createNInboundTracker(keeper, ctx, 10, 5)
		rst := keeper.GetAllInboundTracker(ctx)
		require.Equal(t, intxTrackers, rst)
	})
	t.Run("Get all InTx trackers for chain", func(t *testing.T) {
		keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
		intxTrackersNew := createNInboundTracker(keeper, ctx, 100, 6)
		rst := keeper.GetAllInboundTrackerForChain(ctx, 6)
		sort.SliceStable(rst, func(i, j int) bool {
			return rst[i].TxHash < rst[j].TxHash
		})
		sort.SliceStable(intxTrackersNew, func(i, j int) bool {
			return intxTrackersNew[i].TxHash < intxTrackersNew[j].TxHash
		})
		require.Equal(t, intxTrackersNew, rst)
	})
	t.Run("Get all InTx trackers for chain paginated by limit", func(t *testing.T) {
		keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
		intxTrackers := createNInboundTracker(keeper, ctx, 100, 6)
		rst, pageRes, err := keeper.GetAllInboundTrackerForChainPaginated(ctx, 6, &query.PageRequest{Limit: 10, CountTotal: true})
		require.NoError(t, err)
		require.Subset(t, nullify.Fill(intxTrackers), nullify.Fill(rst))
		require.Equal(t, len(intxTrackers), int(pageRes.Total))
	})
	t.Run("Get all InTx trackers for chain paginated by offset", func(t *testing.T) {
		keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
		intxTrackers := createNInboundTracker(keeper, ctx, 100, 6)
		rst, pageRes, err := keeper.GetAllInboundTrackerForChainPaginated(ctx, 6, &query.PageRequest{Offset: 10, CountTotal: true})
		require.NoError(t, err)
		require.Subset(t, nullify.Fill(intxTrackers), nullify.Fill(rst))
		require.Equal(t, len(intxTrackers), int(pageRes.Total))
	})
	t.Run("Get all InTx trackers paginated by limit", func(t *testing.T) {
		keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
		intxTrackers := append(createNInboundTracker(keeper, ctx, 10, 6), createNInboundTracker(keeper, ctx, 10, 7)...)
		rst, pageRes, err := keeper.GetAllInboundTrackerPaginated(ctx, &query.PageRequest{Limit: 20, CountTotal: true})
		require.NoError(t, err)
		require.Subset(t, nullify.Fill(intxTrackers), nullify.Fill(rst))
		require.Equal(t, len(intxTrackers), int(pageRes.Total))
	})
	t.Run("Get all InTx trackers paginated by offset", func(t *testing.T) {
		keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
		intxTrackers := append(createNInboundTracker(keeper, ctx, 100, 6), createNInboundTracker(keeper, ctx, 100, 7)...)
		rst, pageRes, err := keeper.GetAllInboundTrackerPaginated(ctx, &query.PageRequest{Offset: 10, CountTotal: true})
		require.NoError(t, err)
		require.Subset(t, nullify.Fill(intxTrackers), nullify.Fill(rst))
		require.Equal(t, len(intxTrackers), int(pageRes.Total))
	})
	t.Run("Delete InboundTracker", func(t *testing.T) {
		keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
		intxTrackers := createNInboundTracker(keeper, ctx, 10, 5)
		trackers := keeper.GetAllInboundTracker(ctx)
		for _, item := range trackers {
			keeper.RemoveInboundTrackerIfExists(ctx, item.ChainId, item.TxHash)
		}

		intxTrackers = createNInboundTracker(keeper, ctx, 10, 6)
		for _, item := range intxTrackers {
			keeper.RemoveInboundTrackerIfExists(ctx, item.ChainId, item.TxHash)
		}
		rst := keeper.GetAllInboundTrackerForChain(ctx, 6)
		require.Equal(t, 0, len(rst))
	})
}
