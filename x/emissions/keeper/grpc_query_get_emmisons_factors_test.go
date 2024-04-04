package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

func TestKeeper_GetEmissionsFactors(t *testing.T) {
	k, ctx, _, _ := keepertest.EmissionsKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)

	res, err := k.GetEmissionsFactors(wctx, nil)
	require.NoError(t, err)

	reservesFactor, bondFactor, durationFactor := k.GetBlockRewardComponents(ctx)
	expectedRes := &types.QueryGetEmissionsFactorsResponse{
		ReservesFactor: reservesFactor.String(),
		BondFactor:     bondFactor.String(),
		DurationFactor: durationFactor.String(),
	}
	require.Equal(t, expectedRes, res)
}
