package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

func createNTSS(keeper *Keeper, ctx sdk.Context, n int) []types.TSS {
	items := make([]types.TSS, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Index = fmt.Sprintf("%d", i)
		keeper.SetTSS(ctx, items[i])
	}
	return items
}

func TestTSSGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTSS(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetTSS(ctx, item.Index)
		assert.True(t, found)
		assert.Equal(t, item, rst)
	}
}
func TestTSSRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTSS(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveTSS(ctx, item.Index)
		_, found := keeper.GetTSS(ctx, item.Index)
		assert.False(t, found)
	}
}

func TestTSSGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNTSS(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllTSS(ctx))
}

func TestTSSQuerySingle(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNTSS(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetTSSRequest
		response *types.QueryGetTSSResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetTSSRequest{Index: msgs[0].Index},
			response: &types.QueryGetTSSResponse{TSS: &msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetTSSRequest{Index: msgs[1].Index},
			response: &types.QueryGetTSSResponse{TSS: &msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetTSSRequest{Index: "missing"},
			err:     status.Error(codes.InvalidArgument, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.TSS(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.Equal(t, tc.response, response)
			}
		})
	}
}

// Querier Test

func TestTSSQueryPaginated(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNTSS(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllTSSRequest {
		return &types.QueryAllTSSRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.TSSAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			for j := i; j < len(msgs) && j < i+step; j++ {
				assert.Equal(t, &msgs[j], resp.TSS[j-i])
			}
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.TSSAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			for j := i; j < len(msgs) && j < i+step; j++ {
				assert.Equal(t, &msgs[j], resp.TSS[j-i])
			}
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.TSSAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.TSSAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
