package keeper_test

import (
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/nullify"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestForeignCoinsQuerySingle(t *testing.T) {
	keeper, ctx, _, _ := keepertest.FungibleKeeper(t)
	msgs := createNForeignCoins(t, keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetForeignCoinsRequest
		response *types.QueryGetForeignCoinsResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetForeignCoinsRequest{
				Index: msgs[0].Zrc20ContractAddress,
			},
			response: &types.QueryGetForeignCoinsResponse{ForeignCoins: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetForeignCoinsRequest{
				Index: msgs[1].Zrc20ContractAddress,
			},
			response: &types.QueryGetForeignCoinsResponse{ForeignCoins: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetForeignCoinsRequest{
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
			response, err := keeper.ForeignCoins(ctx, tc.request)
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

func TestForeignCoinsFromAsset(t *testing.T) {
	keeper, ctx, _, _ := keepertest.FungibleKeeper(t)
	msgs := createNForeignCoins(t, keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetForeignCoinsFromAssetRequest
		response *types.QueryGetForeignCoinsFromAssetResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetForeignCoinsFromAssetRequest{
				ChainId: msgs[0].ForeignChainId,
				Asset:   msgs[0].Asset,
			},
			response: &types.QueryGetForeignCoinsFromAssetResponse{ForeignCoins: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetForeignCoinsFromAssetRequest{
				ChainId: msgs[1].ForeignChainId,
				Asset:   msgs[1].Asset,
			},
			response: &types.QueryGetForeignCoinsFromAssetResponse{ForeignCoins: msgs[1]},
		},
		{
			desc: "Not found",
			request: &types.QueryGetForeignCoinsFromAssetRequest{
				ChainId: msgs[0].ForeignChainId + 1,
				Asset:   msgs[0].Asset,
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc:    "Invalid request",
			request: nil,
			err:     status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.ForeignCoinsFromAsset(ctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, nullify.Fill(tc.response), nullify.Fill(response))
			}
		})
	}
}

func TestForeignCoinsQueryPaginated(t *testing.T) {
	keeper, ctx, _, _ := keepertest.FungibleKeeper(t)
	msgs := createNForeignCoins(t, keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllForeignCoinsRequest {
		return &types.QueryAllForeignCoinsRequest{
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
			resp, err := keeper.ForeignCoinsAll(ctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ForeignCoins), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.ForeignCoins),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.ForeignCoinsAll(ctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.ForeignCoins), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.ForeignCoins),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.ForeignCoinsAll(ctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.ForeignCoins),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.ForeignCoinsAll(ctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
