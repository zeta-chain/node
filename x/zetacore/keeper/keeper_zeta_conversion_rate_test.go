package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/zetacore/keeper"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

// Keeper Tests
func createNZetaConversionRate(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.ZetaConversionRate {
	items := make([]types.ZetaConversionRate, n)
	for i := range items {
		items[i].Index = strconv.Itoa(i)

		keeper.SetZetaConversionRate(ctx, items[i])
	}
	return items
}

func TestZetaConversionRateGet(t *testing.T) {
	keeper, ctx := keepertest.ZetacoreKeeper(t)
	items := createNZetaConversionRate(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetZetaConversionRate(ctx,
			item.Index,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestZetaConversionRateRemove(t *testing.T) {
	keeper, ctx := keepertest.ZetacoreKeeper(t)
	items := createNZetaConversionRate(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveZetaConversionRate(ctx,
			item.Index,
		)
		_, found := keeper.GetZetaConversionRate(ctx,
			item.Index,
		)
		require.False(t, found)
	}
}

func TestZetaConversionRateGetAll(t *testing.T) {
	keeper, ctx := keepertest.ZetacoreKeeper(t)
	items := createNZetaConversionRate(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllZetaConversionRate(ctx)),
	)
}

// Querier Tests

func TestZetaConversionRateQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.ZetacoreKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNZetaConversionRate(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetZetaConversionRateRequest
		response *types.QueryGetZetaConversionRateResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetZetaConversionRateRequest{
				Index: msgs[0].Index,
			},
			response: &types.QueryGetZetaConversionRateResponse{ZetaConversionRate: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetZetaConversionRateRequest{
				Index: msgs[1].Index,
			},
			response: &types.QueryGetZetaConversionRateResponse{ZetaConversionRate: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetZetaConversionRateRequest{
				Index: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.ZetaConversionRate(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestZetaConversionRateQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.ZetacoreKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNZetaConversionRate(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllZetaConversionRateRequest {
		return &types.QueryAllZetaConversionRateRequest{
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
			resp, err := keeper.ZetaConversionRateAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ZetaConversionRate), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.ZetaConversionRate),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.ZetaConversionRateAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ZetaConversionRate), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.ZetaConversionRate),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.ZetaConversionRateAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.ZetaConversionRate),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.ZetaConversionRateAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
