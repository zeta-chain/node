package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/mirror/types"
)

func TestERC20TokenPairsQuery(t *testing.T) {
	keeper, ctx := keepertest.MirrorKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	item := createTestERC20TokenPairs(keeper, ctx)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetERC20TokenPairsRequest
		response *types.QueryGetERC20TokenPairsResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetERC20TokenPairsRequest{},
			response: &types.QueryGetERC20TokenPairsResponse{ERC20TokenPairs: item},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.ERC20TokenPairs(wctx, tc.request)
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
