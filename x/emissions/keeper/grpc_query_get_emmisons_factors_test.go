package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/x/emissions/types"
)

func TestKeeper_GetEmissionsFactors(t *testing.T) {
	t.Run("should return emissions factor", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.GetEmissionsFactors(wctx, nil)
		require.NoError(t, err)

		reservesFactor, bondFactor, durationFactor := k.GetBlockRewardComponents(ctx, types.DefaultParams())
		expectedRes := &types.QueryGetEmissionsFactorsResponse{
			ReservesFactor: reservesFactor.String(),
			BondFactor:     bondFactor.String(),
			DurationFactor: durationFactor.String(),
		}
		require.Equal(t, expectedRes, res)
	})

	t.Run("should fail if params not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionKeeperWithMockOptions(
			t,
			keepertest.EmissionMockOptions{SkipSettingParams: true},
		)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.GetEmissionsFactors(wctx, nil)
		require.Nil(t, res)
		require.Error(t, err)
	})
}
