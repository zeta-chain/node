package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestTSSQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.ObserverKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createTSS(keeper, ctx, 1)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetTSSRequest
		response *types.QueryGetTSSResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetTSSRequest{},
			response: &types.QueryGetTSSResponse{TSS: &msgs[len(msgs)-1]},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.TSS(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.Equal(t, tc.response, response)
			}
		})
	}
}

func TestTSSQueryHistory(t *testing.T) {
	keeper, ctx := keepertest.ObserverKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	for _, tc := range []struct {
		desc          string
		tssCount      int
		foundPrevious bool
		err           error
	}{
		{
			desc:          "1 Tss addresses",
			tssCount:      1,
			foundPrevious: false,
			err:           nil,
		},
		{
			desc:          "10 Tss addresses",
			tssCount:      10,
			foundPrevious: true,
			err:           nil,
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			tssList := createTSS(keeper, ctx, tc.tssCount)
			request := &types.QueryTssHistoryRequest{}
			response, err := keeper.TssHistory(wctx, request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.Equal(t, len(tssList), len(response.TssList))
				prevTss, found := keeper.GetPreviousTSS(ctx)
				assert.Equal(t, tc.foundPrevious, found)
				if found {
					assert.Equal(t, tssList[len(tssList)-2], prevTss)
				}
			}
		})
	}
}
