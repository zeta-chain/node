package keeper

import (
	"github.com/stretchr/testify/require"

	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	"github.com/zeta-chain/zetacore/x/observer/types"
)

// Keeper Tests
func createTestKeygen(keeper *Keeper, ctx sdk.Context) types.Keygen {
	item := types.Keygen{
		BlockNumber: 10,
	}
	keeper.SetKeygen(ctx, item)
	return item
}

func TestKeygenGet(t *testing.T) {
	keeper, ctx := SetupKeeper(t)
	item := createTestKeygen(keeper, ctx)
	rst, found := keeper.GetKeygen(ctx)
	assert.True(t, found)
	assert.Equal(t, item, rst)
}
func TestKeygenRemove(t *testing.T) {
	keeper, ctx := SetupKeeper(t)
	createTestKeygen(keeper, ctx)
	keeper.RemoveKeygen(ctx)
	_, found := keeper.GetKeygen(ctx)
	assert.False(t, found)
}

// Querier Tests

func TestKeygenQuery(t *testing.T) {
	keeper, ctx := SetupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	item := createTestKeygen(keeper, ctx)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetKeygenRequest
		response *types.QueryGetKeygenResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetKeygenRequest{},
			response: &types.QueryGetKeygenResponse{Keygen: &item},
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Keygen(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.Equal(t, tc.response, response)
			}
		})
	}
}
