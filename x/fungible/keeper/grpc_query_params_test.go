package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	testkeeper "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx, _, _ := testkeeper.FungibleKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	response, err := keeper.Params(wctx, &types.QueryParamsRequest{})
	assert.NoError(t, err)
	assert.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
