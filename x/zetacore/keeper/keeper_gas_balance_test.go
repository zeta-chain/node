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

// Keeper Tests
func createNGasBalance(keeper *Keeper, ctx sdk.Context, n int) []types.GasBalance {
	items := make([]types.GasBalance, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].Index = fmt.Sprintf("%d", i)
		keeper.SetGasBalance(ctx, items[i])
	}
	return items
}

func TestGasBalanceGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNGasBalance(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetGasBalance(ctx, item.Index)
		assert.True(t, found)
		assert.Equal(t, item, rst)
	}
}
func TestGasBalanceRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNGasBalance(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveGasBalance(ctx, item.Index)
		_, found := keeper.GetGasBalance(ctx, item.Index)
		assert.False(t, found)
	}
}

func TestGasBalanceGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNGasBalance(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllGasBalance(ctx))
}

// Querier Tests
func TestGasBalanceQuerySingle(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNGasBalance(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetGasBalanceRequest
		response *types.QueryGetGasBalanceResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetGasBalanceRequest{Index: msgs[0].Index},
			response: &types.QueryGetGasBalanceResponse{GasBalance: &msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetGasBalanceRequest{Index: msgs[1].Index},
			response: &types.QueryGetGasBalanceResponse{GasBalance: &msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetGasBalanceRequest{Index: "missing"},
			err:     status.Error(codes.InvalidArgument, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.GasBalance(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.Equal(t, tc.response, response)
			}
		})
	}
}

func TestGasBalanceQueryPaginated(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNGasBalance(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllGasBalanceRequest {
		return &types.QueryAllGasBalanceRequest{
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
			resp, err := keeper.GasBalanceAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			for j := i; j < len(msgs) && j < i+step; j++ {
				assert.Equal(t, &msgs[j], resp.GasBalance[j-i])
			}
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.GasBalanceAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			for j := i; j < len(msgs) && j < i+step; j++ {
				assert.Equal(t, &msgs[j], resp.GasBalance[j-i])
			}
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.GasBalanceAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.GasBalanceAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
