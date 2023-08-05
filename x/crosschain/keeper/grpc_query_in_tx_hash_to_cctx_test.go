package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestInTxHashToCctxQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.CrosschainKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNInTxHashToCctx(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetInTxHashToCctxRequest
		response *types.QueryGetInTxHashToCctxResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetInTxHashToCctxRequest{
				InTxHash: msgs[0].InTxHash,
			},
			response: &types.QueryGetInTxHashToCctxResponse{InTxHashToCctx: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetInTxHashToCctxRequest{
				InTxHash: msgs[1].InTxHash,
			},
			response: &types.QueryGetInTxHashToCctxResponse{InTxHashToCctx: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetInTxHashToCctxRequest{
				InTxHash: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.InTxHashToCctx(wctx, tc.request)
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

func TestInTxHashToCctxQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.CrosschainKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNInTxHashToCctx(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllInTxHashToCctxRequest {
		return &types.QueryAllInTxHashToCctxRequest{
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
			resp, err := keeper.InTxHashToCctxAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.InTxHashToCctx), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.InTxHashToCctx),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.InTxHashToCctxAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.InTxHashToCctx), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.InTxHashToCctx),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.InTxHashToCctxAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.InTxHashToCctx),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.InTxHashToCctxAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
