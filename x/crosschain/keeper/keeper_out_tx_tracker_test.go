package keeper_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/zeta-chain/zetacore/x/crosschain/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

// Keeper Tests
func createNOutTxTracker(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.OutTxTracker {
	items := make([]types.OutTxTracker, n)
	for i := range items {
		items[i].ChainId = int64(i)
		items[i].Nonce = uint64(i)
		items[i].Index = fmt.Sprintf("%d-%d", items[i].ChainId, items[i].Nonce)

		keeper.SetOutTxTracker(ctx, items[i])
	}
	return items
}

func TestOutTxTrackerGet(t *testing.T) {
	keeper, ctx := keepertest.CrosschainKeeper(t)
	items := createNOutTxTracker(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetOutTxTracker(ctx,
			item.ChainId,
			item.Nonce,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestOutTxTrackerRemove(t *testing.T) {
	k, ctx := keepertest.CrosschainKeeper(t)
	items := createNOutTxTracker(k, ctx, 10)
	for _, item := range items {
		k.RemoveOutTxTracker(ctx,
			item.ChainId,
			item.Nonce,
		)
		_, found := k.GetOutTxTracker(ctx,
			item.ChainId,
			item.Nonce,
		)
		require.False(t, found)
	}
}

func TestOutTxTrackerGetAll(t *testing.T) {
	keeper, ctx := keepertest.CrosschainKeeper(t)
	items := createNOutTxTracker(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllOutTxTracker(ctx)),
	)
}

// Querier Tests

//func TestOutTxTrackerQuerySingle(t *testing.T) {
//	keeper, ctx := keepertest.ZetacoreKeeper(t)
//	wctx := sdk.WrapSDKContext(ctx)
//	msgs := createNOutTxTracker(keeper, ctx, 2)
//	for _, tc := range []struct {
//		desc     string
//		request  *types.QueryGetOutTxTrackerRequest
//		response *types.QueryGetOutTxTrackerResponse
//		err      error
//	}{
//		{
//			desc: "First",
//			request: &types.QueryGetOutTxTrackerRequest{
//				Index: msgs[0].Index,
//			},
//			response: &types.QueryGetOutTxTrackerResponse{OutTxTracker: msgs[0]},
//		},
//		{
//			desc: "Second",
//			request: &types.QueryGetOutTxTrackerRequest{
//				Index: msgs[1].Index,
//			},
//			response: &types.QueryGetOutTxTrackerResponse{OutTxTracker: msgs[1]},
//		},
//		{
//			desc: "KeyNotFound",
//			request: &types.QueryGetOutTxTrackerRequest{
//				Index: strconv.Itoa(100000),
//			},
//			err: status.Error(codes.NotFound, "not found"),
//		},
//		{
//			desc: "InvalidRequest",
//			err:  status.Error(codes.InvalidArgument, "invalid request"),
//		},
//	} {
//		t.Run(tc.desc, func(t *testing.T) {
//			response, err := keeper.OutTxTracker(wctx, tc.request)
//			if tc.err != nil {
//				require.ErrorIs(t, err, tc.err)
//			} else {
//				require.NoError(t, err)
//				require.Equal(t,
//					nullify.Fill(tc.response),
//					nullify.Fill(response),
//				)
//			}
//		})
//	}
//}
//
//func TestOutTxTrackerQueryPaginated(t *testing.T) {
//	keeper, ctx := keepertest.ZetacoreKeeper(t)
//	wctx := sdk.WrapSDKContext(ctx)
//	msgs := createNOutTxTracker(keeper, ctx, 5)
//
//	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllOutTxTrackerRequest {
//		return &types.QueryAllOutTxTrackerRequest{
//			Pagination: &query.PageRequest{
//				Key:        next,
//				Offset:     offset,
//				Limit:      limit,
//				CountTotal: total,
//			},
//		}
//	}
//	t.Run("ByOffset", func(t *testing.T) {
//		step := 2
//		for i := 0; i < len(msgs); i += step {
//			resp, err := keeper.OutTxTrackerAll(wctx, request(nil, uint64(i), uint64(step), false))
//			require.NoError(t, err)
//			require.LessOrEqual(t, len(resp.OutTxTracker), step)
//			require.Subset(t,
//				nullify.Fill(msgs),
//				nullify.Fill(resp.OutTxTracker),
//			)
//		}
//	})
//	t.Run("ByKey", func(t *testing.T) {
//		step := 2
//		var next []byte
//		for i := 0; i < len(msgs); i += step {
//			resp, err := keeper.OutTxTrackerAll(wctx, request(next, 0, uint64(step), false))
//			require.NoError(t, err)
//			require.LessOrEqual(t, len(resp.OutTxTracker), step)
//			require.Subset(t,
//				nullify.Fill(msgs),
//				nullify.Fill(resp.OutTxTracker),
//			)
//			next = resp.Pagination.NextKey
//		}
//	})
//	t.Run("Total", func(t *testing.T) {
//		resp, err := keeper.OutTxTrackerAll(wctx, request(nil, 0, 0, true))
//		require.NoError(t, err)
//		require.Equal(t, len(msgs), int(resp.Pagination.Total))
//		require.ElementsMatch(t,
//			nullify.Fill(msgs),
//			nullify.Fill(resp.OutTxTracker),
//		)
//	})
//	t.Run("InvalidRequest", func(t *testing.T) {
//		_, err := keeper.OutTxTrackerAll(wctx, nil)
//		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
//	})
//}
