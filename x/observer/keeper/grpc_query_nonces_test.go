package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestChainNoncesQuerySingle(t *testing.T) {
	k, ctx := keepertest.ObserverKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	chainNonces := sample.ChainNoncesList(t, 2)
	for _, nonce := range chainNonces {
		k.SetChainNonces(ctx, nonce)
	}
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetChainNoncesRequest
		response *types.QueryGetChainNoncesResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetChainNoncesRequest{Index: chainNonces[0].Index},
			response: &types.QueryGetChainNoncesResponse{ChainNonces: chainNonces[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetChainNoncesRequest{Index: chainNonces[1].Index},
			response: &types.QueryGetChainNoncesResponse{ChainNonces: chainNonces[1]},
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
			response, err := k.ChainNonces(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.Equal(t, tc.response, response)
			}
		})
	}
}

func TestChainNoncesQueryPaginated(t *testing.T) {
	k, ctx := keepertest.ObserverKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	chainNonces := sample.ChainNoncesList(t, 5)
	for _, nonce := range chainNonces {
		k.SetChainNonces(ctx, nonce)
	}

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
		for i := 0; i < len(chainNonces); i += step {
			resp, err := k.ChainNoncesAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			for j := i; j < len(chainNonces) && j < i+step; j++ {
				require.Equal(t, chainNonces[j], resp.ChainNonces[j-i])
			}
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(chainNonces); i += step {
			resp, err := k.ChainNoncesAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			for j := i; j < len(chainNonces) && j < i+step; j++ {
				require.Equal(t, chainNonces[j], resp.ChainNonces[j-i])
			}
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := k.ChainNoncesAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(chainNonces), int(resp.Pagination.Total))
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := k.ChainNoncesAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
