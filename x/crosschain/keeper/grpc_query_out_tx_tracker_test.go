package keeper_test

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
