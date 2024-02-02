package keeper_test

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/assert"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestTSSGet(t *testing.T) {
	k, ctx := keepertest.ObserverKeeper(t)
	tss := sample.Tss()
	k.SetTSS(ctx, tss)
	tssQueried, found := k.GetTSS(ctx)
	assert.True(t, found)
	assert.Equal(t, tss, tssQueried)

}
func TestTSSRemove(t *testing.T) {
	k, ctx := keepertest.ObserverKeeper(t)
	tss := sample.Tss()
	k.SetTSS(ctx, tss)
	k.RemoveTSS(ctx)
	_, found := k.GetTSS(ctx)
	assert.False(t, found)
}

func TestTSSQuerySingle(t *testing.T) {
	k, ctx := keepertest.ObserverKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	//msgs := createTSS(keeper, ctx, 1)
	tss := sample.Tss()
	k.SetTSS(ctx, tss)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetTSSRequest
		response *types.QueryGetTSSResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetTSSRequest{},
			response: &types.QueryGetTSSResponse{TSS: tss},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			response, err := k.TSS(wctx, tc.request)
			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
			} else {
				assert.Equal(t, tc.response, response)
			}
		})
	}
}

func TestTSSQueryHistory(t *testing.T) {
	keeper, ctx := keepertest.ObserverKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	for _, tc := range []struct {
		desc          string
		tssCount      int
		foundPrevious bool
		err           error
	}{
		{
			desc:          "1 Tss addresses",
			tssCount:      1,
			foundPrevious: false,
			err:           nil,
		},
		{
			desc:          "10 Tss addresses",
			tssCount:      10,
			foundPrevious: true,
			err:           nil,
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			tssList := sample.TssList(tc.tssCount)
			for _, tss := range tssList {
				keeper.SetTSS(ctx, tss)
				keeper.SetTSSHistory(ctx, tss)
			}
			request := &types.QueryTssHistoryRequest{}
			response, err := keeper.TssHistory(wctx, request)
			if tc.err != nil {
				assert.ErrorIs(t, err, tc.err)
			} else {
				assert.Equal(t, len(tssList), len(response.TssList))
				prevTss, found := keeper.GetPreviousTSS(ctx)
				assert.Equal(t, tc.foundPrevious, found)
				if found {
					assert.Equal(t, tssList[len(tssList)-2], prevTss)
				}
			}
		})
	}
}

func TestKeeper_TssHistory(t *testing.T) {
	t.Run("Get tss history paginated by limit", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		tssList := sample.TssList(10)
		for _, tss := range tssList {
			k.SetTSSHistory(ctx, tss)
		}
		rst, pageRes, err := k.GetAllTSSPaginated(ctx, &query.PageRequest{Limit: 20, CountTotal: true})
		assert.NoError(t, err)
		sort.Slice(tssList, func(i, j int) bool {
			return tssList[i].FinalizedZetaHeight < tssList[j].FinalizedZetaHeight
		})
		sort.Slice(rst, func(i, j int) bool {
			return rst[i].FinalizedZetaHeight < rst[j].FinalizedZetaHeight
		})
		assert.Equal(t, tssList, rst)
		assert.Equal(t, len(tssList), int(pageRes.Total))
	})
	t.Run("Get tss history paginated by offset", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		tssList := sample.TssList(100)
		offset := 20
		for _, tss := range tssList {
			k.SetTSSHistory(ctx, tss)
		}
		rst, pageRes, err := k.GetAllTSSPaginated(ctx, &query.PageRequest{Offset: uint64(offset), CountTotal: true})
		assert.NoError(t, err)
		sort.Slice(tssList, func(i, j int) bool {
			return tssList[i].FinalizedZetaHeight < tssList[j].FinalizedZetaHeight
		})
		sort.Slice(rst, func(i, j int) bool {
			return rst[i].FinalizedZetaHeight < rst[j].FinalizedZetaHeight
		})
		assert.Subset(t, tssList, rst)
		assert.Equal(t, len(tssList)-offset, len(rst))
		assert.Equal(t, len(tssList), int(pageRes.Total))
	})
	t.Run("Get all TSS without pagination", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
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
		assert.Equal(t, tssList, rst)
	})
	t.Run("Get historical TSS", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		tssList := sample.TssList(100)
		for _, tss := range tssList {
			k.SetTSSHistory(ctx, tss)
		}
		r := rand.Intn((len(tssList)-1)-0) + 0
		tss, found := k.GetHistoricalTssByFinalizedHeight(ctx, tssList[r].FinalizedZetaHeight)
		assert.True(t, found)
		assert.Equal(t, tssList[r], tss)
	})
}
