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

// Keeper Tests
func createNChainNonces(keeper *Keeper, ctx sdk.Context, n int) []types.ChainNonces {
	items := make([]types.ChainNonces, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Index = fmt.Sprintf("%d", i)
		keeper.SetChainNonces(ctx, items[i])
	}
	return items
}

func TestChainNoncesGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNChainNonces(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetChainNonces(ctx, item.Index)
		assert.True(t, found)
		assert.Equal(t, item, rst)
	}
}
func TestChainNoncesRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNChainNonces(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveChainNonces(ctx, item.Index)
		_, found := keeper.GetChainNonces(ctx, item.Index)
		assert.False(t, found)
	}
}

func TestChainNoncesGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNChainNonces(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllChainNonces(ctx))
}

//Querier Tests

func TestChainNoncesQuerySingle(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNChainNonces(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetChainNoncesRequest
		response *types.QueryGetChainNoncesResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetChainNoncesRequest{Index: msgs[0].Index},
			response: &types.QueryGetChainNoncesResponse{ChainNonces: &msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetChainNoncesRequest{Index: msgs[1].Index},
			response: &types.QueryGetChainNoncesResponse{ChainNonces: &msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetChainNoncesRequest{Index: "missing"},
			err:     status.Error(codes.InvalidArgument, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.ChainNonces(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.Equal(t, tc.response, response)
			}
		})
	}
}

func TestChainNoncesQueryPaginated(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNChainNonces(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllChainNoncesRequest {
		return &types.QueryAllChainNoncesRequest{
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
			resp, err := keeper.ChainNoncesAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			for j := i; j < len(msgs) && j < i+step; j++ {
				assert.Equal(t, &msgs[j], resp.ChainNonces[j-i])
			}
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.ChainNoncesAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			for j := i; j < len(msgs) && j < i+step; j++ {
				assert.Equal(t, &msgs[j], resp.ChainNonces[j-i])
			}
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.ChainNoncesAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.ChainNoncesAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
