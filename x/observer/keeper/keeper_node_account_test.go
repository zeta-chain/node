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

	"github.com/zeta-chain/zetacore/x/observer/types"
)

// Keeper Tests

func createNNodeAccount(keeper *Keeper, ctx sdk.Context, n int) []types.NodeAccount {
	items := make([]types.NodeAccount, n)
	for i := range items {
		items[i].Operator = fmt.Sprintf("%d", i)
		keeper.SetNodeAccount(ctx, items[i])
	}
	return items
}

func TestNodeAccountGet(t *testing.T) {
	keeper, ctx := SetupKeeper(t)
	items := createNNodeAccount(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetNodeAccount(ctx, item.Operator)
		assert.True(t, found)
		assert.Equal(t, item, rst)
	}
}
func TestNodeAccountRemove(t *testing.T) {
	keeper, ctx := SetupKeeper(t)
	items := createNNodeAccount(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveNodeAccount(ctx, item.Operator)
		_, found := keeper.GetNodeAccount(ctx, item.Operator)
		assert.False(t, found)
	}
}

func TestNodeAccountGetAll(t *testing.T) {
	keeper, ctx := SetupKeeper(t)
	items := createNNodeAccount(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllNodeAccount(ctx))
}

// Querier Tests

func TestNodeAccountQuerySingle(t *testing.T) {
	keeper, ctx := SetupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNNodeAccount(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetNodeAccountRequest
		response *types.QueryGetNodeAccountResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetNodeAccountRequest{Index: msgs[0].Operator},
			response: &types.QueryGetNodeAccountResponse{NodeAccount: &msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetNodeAccountRequest{Index: msgs[1].Operator},
			response: &types.QueryGetNodeAccountResponse{NodeAccount: &msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetNodeAccountRequest{Index: "missing"},
			err:     status.Error(codes.InvalidArgument, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.NodeAccount(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.Equal(t, tc.response, response)
			}
		})
	}
}

func TestNodeAccountQueryPaginated(t *testing.T) {
	keeper, ctx := SetupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNNodeAccount(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllNodeAccountRequest {
		return &types.QueryAllNodeAccountRequest{
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
			resp, err := keeper.NodeAccountAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			for j := i; j < len(msgs) && j < i+step; j++ {
				assert.Equal(t, &msgs[j], resp.NodeAccount[j-i])
			}
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.NodeAccountAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			for j := i; j < len(msgs) && j < i+step; j++ {
				assert.Equal(t, &msgs[j], resp.NodeAccount[j-i])
			}
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.NodeAccountAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.NodeAccountAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
