package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestNodeAccountQuerySingle(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)

	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNNodeAccount(k, ctx, 2)
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
		t.Run(tc.desc, func(t *testing.T) {
			response, err := k.NodeAccount(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.Equal(t, tc.response, response)
			}
		})
	}
}

func TestNodeAccountQueryPaginated(t *testing.T) {
	k, ctx, _, _ := keepertest.ObserverKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNNodeAccount(k, ctx, 5)

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
			resp, err := k.NodeAccountAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			for j := i; j < len(msgs) && j < i+step; j++ {
				require.Equal(t, &msgs[j], resp.NodeAccount[j-i])
			}
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := k.NodeAccountAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			for j := i; j < len(msgs) && j < i+step; j++ {
				require.Equal(t, &msgs[j], resp.NodeAccount[j-i])
			}
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := k.NodeAccountAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := k.NodeAccountAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
