package keeper_test

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestKeeper_GetTSS(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)
	tss := sample.Tss()
	k.SetTSS(ctx, tss)
	tssQueried, found := k.GetTSS(ctx)
	require.True(t, found)
	require.Equal(t, tss, tssQueried)

}
func TestKeeper_RemoveTSS(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)
	tss := sample.Tss()
	k.SetTSS(ctx, tss)
	k.RemoveTSS(ctx)
	_, found := k.GetTSS(ctx)
	require.False(t, found)
}

func TestKeeper_CheckIfTssPubkeyHasBeenGenerated(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)
	tss := sample.Tss()

	generated, found := k.CheckIfTssPubkeyHasBeenGenerated(ctx, tss.TssPubkey)
	require.False(t, found)
	require.Equal(t, types.TSS{}, generated)

	k.AppendTss(ctx, tss)

	generated, found = k.CheckIfTssPubkeyHasBeenGenerated(ctx, tss.TssPubkey)
	require.True(t, found)
	require.Equal(t, tss, generated)
}

func TestKeeper_GetHistoricalTssByFinalizedHeight(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)
	tssList := sample.TssList(100)
	r := rand.Intn((len(tssList)-1)-0) + 0
	_, found := k.GetHistoricalTssByFinalizedHeight(ctx, tssList[r].FinalizedZetaHeight)
	require.False(t, found)

	for _, tss := range tssList {
		k.SetTSSHistory(ctx, tss)
	}
	tss, found := k.GetHistoricalTssByFinalizedHeight(ctx, tssList[r].FinalizedZetaHeight)
	require.True(t, found)
	require.Equal(t, tssList[r], tss)
}

func TestKeeper_TssHistory(t *testing.T) {
	t.Run("Get tss history paginated by limit", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		tssList := sample.TssList(10)
		for _, tss := range tssList {
			k.SetTSSHistory(ctx, tss)
		}
		rst, pageRes, err := k.GetAllTSSPaginated(ctx, &query.PageRequest{Limit: 20, CountTotal: true})
		require.NoError(t, err)
		sort.Slice(tssList, func(i, j int) bool {
			return tssList[i].FinalizedZetaHeight < tssList[j].FinalizedZetaHeight
		})
		sort.Slice(rst, func(i, j int) bool {
			return rst[i].FinalizedZetaHeight < rst[j].FinalizedZetaHeight
		})
		require.Equal(t, tssList, rst)
		require.Equal(t, len(tssList), int(pageRes.Total))
	})
	t.Run("Get tss history paginated by offset", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		tssList := sample.TssList(100)
		offset := 20
		for _, tss := range tssList {
			k.SetTSSHistory(ctx, tss)
		}
		rst, pageRes, err := k.GetAllTSSPaginated(ctx, &query.PageRequest{Offset: uint64(offset), CountTotal: true})
		require.NoError(t, err)
		sort.Slice(tssList, func(i, j int) bool {
			return tssList[i].FinalizedZetaHeight < tssList[j].FinalizedZetaHeight
		})
		sort.Slice(rst, func(i, j int) bool {
			return rst[i].FinalizedZetaHeight < rst[j].FinalizedZetaHeight
		})
		require.Subset(t, tssList, rst)
		require.Equal(t, len(tssList)-offset, len(rst))
		require.Equal(t, len(tssList), int(pageRes.Total))
	})
	t.Run("Get all TSS without pagination", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		tssList := sample.TssList(100)
		for _, tss := range tssList {
			k.SetTSSHistory(ctx, tss)
		}
		rst := k.GetAllTSS(ctx)
		sort.Slice(tssList, func(i, j int) bool {
			return tssList[i].FinalizedZetaHeight < tssList[j].FinalizedZetaHeight
		})
		sort.Slice(rst, func(i, j int) bool {
			return rst[i].FinalizedZetaHeight < rst[j].FinalizedZetaHeight
		})
		require.Equal(t, tssList, rst)
	})
}
