package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestZetaDepositAndCallContractQuery(t *testing.T) {
	keeper, ctx := keepertest.FungibleKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	item := createTestZetaDepositAndCallContract(keeper, ctx)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetZetaDepositAndCallContractRequest
		response *types.QueryGetZetaDepositAndCallContractResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetZetaDepositAndCallContractRequest{},
			response: &types.QueryGetZetaDepositAndCallContractResponse{ZetaDepositAndCallContract: item},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.ZetaDepositAndCallContract(wctx, tc.request)
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

