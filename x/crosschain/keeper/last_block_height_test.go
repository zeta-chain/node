package keeper

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func createNLastBlockHeight(keeper *Keeper, ctx sdk.Context, n int) []types.LastBlockHeight {
	items := make([]types.LastBlockHeight, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Index = fmt.Sprintf("%d", i)
		keeper.SetLastBlockHeight(ctx, items[i])
	}
	return items
}

func TestLastBlockHeightGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNLastBlockHeight(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetLastBlockHeight(ctx, item.Index)
		assert.True(t, found)
		assert.Equal(t, item, rst)
	}
}
func TestLastBlockHeightRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNLastBlockHeight(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveLastBlockHeight(ctx, item.Index)
		_, found := keeper.GetLastBlockHeight(ctx, item.Index)
		assert.False(t, found)
	}
}

func TestLastBlockHeightGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNLastBlockHeight(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllLastBlockHeight(ctx))
}

//Querier Test

func TestLastBlockHeightQuerySingle(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNLastBlockHeight(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetLastBlockHeightRequest
		response *types.QueryGetLastBlockHeightResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetLastBlockHeightRequest{Index: msgs[0].Index},
			response: &types.QueryGetLastBlockHeightResponse{LastBlockHeight: &msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetLastBlockHeightRequest{Index: msgs[1].Index},
			response: &types.QueryGetLastBlockHeightResponse{LastBlockHeight: &msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetLastBlockHeightRequest{Index: "missing"},
			err:     status.Error(codes.InvalidArgument, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.LastBlockHeight(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.Equal(t, tc.response, response)
			}
		})
	}
}

func TestLastBlockHeightQueryPaginated(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNLastBlockHeight(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllLastBlockHeightRequest {
		return &types.QueryAllLastBlockHeightRequest{
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
			resp, err := keeper.LastBlockHeightAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			for j := i; j < len(msgs) && j < i+step; j++ {
				assert.Equal(t, &msgs[j], resp.LastBlockHeight[j-i])
			}
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.LastBlockHeightAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			for j := i; j < len(msgs) && j < i+step; j++ {
				assert.Equal(t, &msgs[j], resp.LastBlockHeight[j-i])
			}
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.LastBlockHeightAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.LastBlockHeightAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
