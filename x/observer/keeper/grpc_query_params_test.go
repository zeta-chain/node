package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := SetupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	response, err := keeper.Params(wctx, &types.QueryParamsRequest{})
	assert.NoError(t, err)
	assert.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
