package keeper

import (
	"fmt"
	"sort"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func createNInTxTracker(keeper *Keeper, ctx sdk.Context, n int, chainID int64) []types.InTxTracker {
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
	keeper, ctx := setupKeeper(t)
	intxTrackers := createNInTxTracker(keeper, ctx, 10, 5)
	t.Run("Get InTx trackers one by one", func(t *testing.T) {
		for _, item := range intxTrackers {
			rst, found := keeper.GetInTxTracker(ctx, item.ChainId, item.TxHash)
			require.True(t, found)
			require.Equal(t, item, rst)
		}
	})
	t.Run("Get all InTx trackers", func(t *testing.T) {
		rst := keeper.GetAllInTxTracker(ctx)
		require.Equal(t, intxTrackers, rst)
	})
	t.Run("Get all InTx trackers for chain", func(t *testing.T) {
		intxTrackersNew := createNInTxTracker(keeper, ctx, 100, 6)
		rst := keeper.GetAllInTxTrackerForChain(ctx, 6)
		sort.SliceStable(rst, func(i, j int) bool {
			return rst[i].TxHash < rst[j].TxHash
		})
		sort.SliceStable(intxTrackersNew, func(i, j int) bool {
			return intxTrackersNew[i].TxHash < intxTrackersNew[j].TxHash
		})
		require.Equal(t, intxTrackersNew, rst)
	})
	t.Run("Get all InTx trackers for chain paginated by limit", func(t *testing.T) {
		intxTrackers = createNInTxTracker(keeper, ctx, 100, 6)
		rst, pageRes, err := keeper.GetAllInTxTrackerForChainPaginated(ctx, 6, &query.PageRequest{Limit: 10, CountTotal: true})
		require.NoError(t, err)
		require.Subset(t, nullify.Fill(intxTrackers), nullify.Fill(rst))
		require.Equal(t, len(intxTrackers), int(pageRes.Total))
	})
	t.Run("Get all InTx trackers for chain paginated by offset", func(t *testing.T) {
		intxTrackers = createNInTxTracker(keeper, ctx, 100, 6)
		rst, pageRes, err := keeper.GetAllInTxTrackerForChainPaginated(ctx, 6, &query.PageRequest{Offset: 10, CountTotal: true})
		require.NoError(t, err)
		require.Subset(t, nullify.Fill(intxTrackers), nullify.Fill(rst))
		require.Equal(t, len(intxTrackers), int(pageRes.Total))
	})
}
