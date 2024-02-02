package keeper_test

import (
	"fmt"
	"sort"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func createNInTxTracker(keeper *keeper.Keeper, ctx sdk.Context, n int, chainID int64) []types.InTxTracker {
	items := make([]types.InTxTracker, n)
	for i := range items {
		items[i].TxHash = fmt.Sprintf("TxHash-%d", i)
		items[i].ChainId = chainID
		items[i].CoinType = common.CoinType_Gas
		keeper.SetInTxTracker(ctx, items[i])
	}
	return items
}
func TestKeeper_GetAllInTxTrackerForChain(t *testing.T) {
	t.Run("Get InTx trackers one by one", func(t *testing.T) {
		keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
		intxTrackers := createNInTxTracker(keeper, ctx, 10, 5)
		for _, item := range intxTrackers {
			rst, found := keeper.GetInTxTracker(ctx, item.ChainId, item.TxHash)
			assert.True(t, found)
			assert.Equal(t, item, rst)
		}
	})
	t.Run("Get all InTx trackers", func(t *testing.T) {
		keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
		intxTrackers := createNInTxTracker(keeper, ctx, 10, 5)
		rst := keeper.GetAllInTxTracker(ctx)
		assert.Equal(t, intxTrackers, rst)
	})
	t.Run("Get all InTx trackers for chain", func(t *testing.T) {
		keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
		intxTrackersNew := createNInTxTracker(keeper, ctx, 100, 6)
		rst := keeper.GetAllInTxTrackerForChain(ctx, 6)
		sort.SliceStable(rst, func(i, j int) bool {
			return rst[i].TxHash < rst[j].TxHash
		})
		sort.SliceStable(intxTrackersNew, func(i, j int) bool {
			return intxTrackersNew[i].TxHash < intxTrackersNew[j].TxHash
		})
		assert.Equal(t, intxTrackersNew, rst)
	})
	t.Run("Get all InTx trackers for chain paginated by limit", func(t *testing.T) {
		keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
		intxTrackers := createNInTxTracker(keeper, ctx, 100, 6)
		rst, pageRes, err := keeper.GetAllInTxTrackerForChainPaginated(ctx, 6, &query.PageRequest{Limit: 10, CountTotal: true})
		assert.NoError(t, err)
		assert.Subset(t, nullify.Fill(intxTrackers), nullify.Fill(rst))
		assert.Equal(t, len(intxTrackers), int(pageRes.Total))
	})
	t.Run("Get all InTx trackers for chain paginated by offset", func(t *testing.T) {
		keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
		intxTrackers := createNInTxTracker(keeper, ctx, 100, 6)
		rst, pageRes, err := keeper.GetAllInTxTrackerForChainPaginated(ctx, 6, &query.PageRequest{Offset: 10, CountTotal: true})
		assert.NoError(t, err)
		assert.Subset(t, nullify.Fill(intxTrackers), nullify.Fill(rst))
		assert.Equal(t, len(intxTrackers), int(pageRes.Total))
	})
	t.Run("Get all InTx trackers paginated by limit", func(t *testing.T) {
		keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
		intxTrackers := append(createNInTxTracker(keeper, ctx, 10, 6), createNInTxTracker(keeper, ctx, 10, 7)...)
		rst, pageRes, err := keeper.GetAllInTxTrackerPaginated(ctx, &query.PageRequest{Limit: 20, CountTotal: true})
		assert.NoError(t, err)
		assert.Subset(t, nullify.Fill(intxTrackers), nullify.Fill(rst))
		assert.Equal(t, len(intxTrackers), int(pageRes.Total))
	})
	t.Run("Get all InTx trackers paginated by offset", func(t *testing.T) {
		keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
		intxTrackers := append(createNInTxTracker(keeper, ctx, 100, 6), createNInTxTracker(keeper, ctx, 100, 7)...)
		rst, pageRes, err := keeper.GetAllInTxTrackerPaginated(ctx, &query.PageRequest{Offset: 10, CountTotal: true})
		assert.NoError(t, err)
		assert.Subset(t, nullify.Fill(intxTrackers), nullify.Fill(rst))
		assert.Equal(t, len(intxTrackers), int(pageRes.Total))
	})
	t.Run("Delete InTxTracker", func(t *testing.T) {
		keeper, ctx, _, _ := keepertest.CrosschainKeeper(t)
		intxTrackers := createNInTxTracker(keeper, ctx, 10, 5)
		trackers := keeper.GetAllInTxTracker(ctx)
		for _, item := range trackers {
			keeper.RemoveInTxTrackerIfExists(ctx, item.ChainId, item.TxHash)
		}

		intxTrackers = createNInTxTracker(keeper, ctx, 10, 6)
		for _, item := range intxTrackers {
			keeper.RemoveInTxTrackerIfExists(ctx, item.ChainId, item.TxHash)
		}
		rst := keeper.GetAllInTxTrackerForChain(ctx, 6)
		assert.Equal(t, 0, len(rst))
	})
}
