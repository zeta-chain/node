package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/zeta-chain/zetacore/x/observer/types"
)

func createTSS(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.TSS {
	tssList := make([]types.TSS, n)
	for i := 0; i < n; i++ {
		tss := types.TSS{
			TssPubkey:           "tssPubkey",
			TssParticipantList:  []string{"tssParticipantList"},
			OperatorAddressList: []string{"operatorAddressList"},
			KeyGenZetaHeight:    int64(100 + i),
			FinalizedZetaHeight: int64(110 + i),
		}
		keeper.SetTSS(ctx, tss)
		keeper.SetTSSHistory(ctx, tss)
		tssList[i] = tss
	}
	return tssList
}

func TestTSSGet(t *testing.T) {
	keeper, ctx := keepertest.ObserverKeeper(t)
	tssSaved := createTSS(keeper, ctx, 1)
	tss, found := keeper.GetTSS(ctx)
	assert.True(t, found)
	assert.Equal(t, tssSaved[len(tssSaved)-1], tss)

}
func TestTSSRemove(t *testing.T) {
	keeper, ctx := keepertest.ObserverKeeper(t)
	_ = createTSS(keeper, ctx, 1)
	keeper.RemoveTSS(ctx)
	_, found := keeper.GetTSS(ctx)
	assert.False(t, found)
}

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
