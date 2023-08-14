package keeper

import (
	"strconv"
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
func createNGasPrice(keeper *Keeper, ctx sdk.Context, n int) []types.GasPrice {
	items := make([]types.GasPrice, n)
	for i := range items {
		items[i].Creator = "any"
		items[i].ChainId = int64(i)
		items[i].Index = strconv.FormatInt(int64(i), 10)
		keeper.SetGasPrice(ctx, items[i])
	}
	return items
}

func TestGasPriceGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNGasPrice(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetGasPrice(ctx, item.ChainId)
		assert.True(t, found)
		assert.Equal(t, item, rst)
	}
}
func TestGasPriceRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNGasPrice(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveGasPrice(ctx, item.Index)
		_, found := keeper.GetGasPrice(ctx, item.ChainId)
		assert.False(t, found)
	}
}

func TestGasPriceGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNGasPrice(keeper, ctx, 10)
	assert.Equal(t, items, keeper.GetAllGasPrice(ctx))
}

// Querier Tests

func TestGasPriceQuerySingle(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNGasPrice(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetGasPriceRequest
		response *types.QueryGetGasPriceResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetGasPriceRequest{Index: msgs[0].Index},
			response: &types.QueryGetGasPriceResponse{GasPrice: &msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetGasPriceRequest{Index: msgs[1].Index},
			response: &types.QueryGetGasPriceResponse{GasPrice: &msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetGasPriceRequest{Index: "1000000000"},
			err:     status.Error(codes.InvalidArgument, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.GasPrice(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.Equal(t, tc.response, response)
			}
		})
	}
}

func TestGasPriceQueryPaginated(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNGasPrice(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllGasPriceRequest {
		return &types.QueryAllGasPriceRequest{
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
			resp, err := keeper.GasPriceAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			for j := i; j < len(msgs) && j < i+step; j++ {
				assert.Equal(t, &msgs[j], resp.GasPrice[j-i])
			}
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.GasPriceAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			for j := i; j < len(msgs) && j < i+step; j++ {
				assert.Equal(t, &msgs[j], resp.GasPrice[j-i])
			}
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.GasPriceAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.GasPriceAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
