package keeper

import (
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func createTSS(keeper *Keeper, ctx sdk.Context) types.TSS {
	tss := types.TSS{
		TssPubkey:           "tssPubkey0",
		TssParticipantList:  []string{"tssParticipantList0"},
		OperatorAddressList: []string{"operatorAddressList0"},
		KeyGenZetaHeight:    100,
		FinalizedZetaHeight: 110,
	}
	keeper.SetTSS(ctx, tss)
	return tss
}

func TestTSSGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	tssSaved := createTSS(keeper, ctx)
	tss, found := keeper.GetTSS(ctx)
	assert.True(t, found)
	assert.Equal(t, tssSaved, tss)

}
func TestTSSRemove(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	_ = createTSS(keeper, ctx)
	keeper.RemoveTSS(ctx)
	_, found := keeper.GetTSS(ctx)
	assert.False(t, found)

}

func TestTSSQuerySingle(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createTSS(keeper, ctx)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetTSSRequest
		response *types.QueryGetTSSResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetTSSRequest{},
			response: &types.QueryGetTSSResponse{TSS: &msgs},
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
