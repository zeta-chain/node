package v3_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/x/emissions/exported"
	"github.com/zeta-chain/node/x/emissions/types"

	v3 "github.com/zeta-chain/node/x/emissions/migrations/v3"
)

type mockSubspace struct {
	ps types.Params
}

func newMockSubspace(ps types.Params) mockSubspace {
	return mockSubspace{ps: ps}
}

func (ms mockSubspace) GetParamSet(ctx sdk.Context, ps exported.ParamSet) {
	*ps.(*types.Params) = ms.ps
}

func TestMigrate(t *testing.T) {
	t.Run("should migrate for valid params", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)

		legacyParams := types.Params{
			ValidatorEmissionPercentage: "00.50",
			ObserverEmissionPercentage:  "00.35",
			TssSignerEmissionPercentage: "00.15",
			ObserverSlashAmount:         sdk.ZeroInt(),
		}
		legacySubspace := newMockSubspace(legacyParams)

		require.NoError(t, v3.MigrateStore(ctx, k, legacySubspace))

		params, found := k.GetParams(ctx)
		require.True(t, found)
		legacyParams.ObserverSlashAmount = sdkmath.NewInt(100000000000000000)
		legacyParams.BallotMaturityBlocks = 100
		legacyParams.BlockRewardAmount = types.BlockReward
		require.Equal(t, legacyParams, params)
	})

	t.Run("should migrate if legacy params missing", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)

		legacyParams := types.Params{}
		legacySubspace := newMockSubspace(legacyParams)

		require.NoError(t, v3.MigrateStore(ctx, k, legacySubspace))

		params, found := k.GetParams(ctx)
		require.True(t, found)
		legacyParams = types.DefaultParams()
		legacyParams.ObserverSlashAmount = sdkmath.NewInt(100000000000000000)
		legacyParams.BallotMaturityBlocks = 100
		require.Equal(t, legacyParams, params)
	})

	t.Run("should fail to migrate for invalid params", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)

		legacyParams := types.Params{
			ValidatorEmissionPercentage: "-00.50",
			ObserverEmissionPercentage:  "00.35",
			TssSignerEmissionPercentage: "00.15",
			ObserverSlashAmount:         sdk.ZeroInt(),
		}
		legacySubspace := newMockSubspace(legacyParams)

		err := v3.MigrateStore(ctx, k, legacySubspace)
		require.ErrorContains(t, err, "validator emission percentage cannot be less than 0 percent")
	})
}
