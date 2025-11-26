package keeper_test

import (
	"math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestLastBlockHeightQuerySingle(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNLastBlockHeight(k, ctx, 2)
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
		t.Run(tc.desc, func(t *testing.T) {
			response, err := k.LastBlockHeight(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.Equal(t, tc.response, response)
			}
		})
	}
}

func TestLastBlockHeightLimits(t *testing.T) {
	t.Run("should err if last send height is max int", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)
		k.SetLastBlockHeight(ctx, types.LastBlockHeight{
			Index:              "index",
			LastOutboundHeight: math.MaxInt64,
		})

		res, err := k.LastBlockHeight(wctx, &types.QueryGetLastBlockHeightRequest{
			Index: "index",
		})
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should err if last receive height is max int", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)
		k.SetLastBlockHeight(ctx, types.LastBlockHeight{
			Index:              "index",
			LastOutboundHeight: 10,
			LastInboundHeight:  math.MaxInt64,
		})

		res, err := k.LastBlockHeight(wctx, &types.QueryGetLastBlockHeightRequest{
			Index: "index",
		})
		require.Nil(t, res)
		require.Error(t, err)
	})
}

func TestLastBlockHeightQueryPaginated(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNLastBlockHeight(k, ctx, 5)

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
			resp, err := k.LastBlockHeightAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			for j := i; j < len(msgs) && j < i+step; j++ {
				require.Equal(t, &msgs[j], resp.LastBlockHeight[j-i])
			}
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := k.LastBlockHeightAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			for j := i; j < len(msgs) && j < i+step; j++ {
				require.Equal(t, &msgs[j], resp.LastBlockHeight[j-i])
			}
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := k.LastBlockHeightAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := k.LastBlockHeightAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
