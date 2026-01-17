package v6_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	v6 "github.com/zeta-chain/node/x/emissions/migrations/v6"
	"github.com/zeta-chain/node/x/emissions/types"
)

func TestMigrateStore(t *testing.T) {
	t.Run("successfully migrate store", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)

		currentParams := types.DefaultParams()
		k.SetParams(ctx, currentParams)

		// Act
		err := v6.MigrateStore(ctx, k)

		// Assert
		// only BallotMaturityBlocks and BlockRewardAmount are updated
		require.NoError(t, err)
		updatedParams, found := k.GetParams(ctx)

		require.NotEqual(t, currentParams.BallotMaturityBlocks, updatedParams.BallotMaturityBlocks)
		require.Equal(t, int64(150), updatedParams.BallotMaturityBlocks)

		require.NotEqual(t, currentParams.BlockRewardAmount, updatedParams.BlockRewardAmount)
		require.Equal(
			t,
			sdkmath.LegacyMustNewDecFromStr("6751543209876543209.876543209876543210"),
			updatedParams.BlockRewardAmount,
		)

		require.NotEqual(
			t,
			currentParams.PendingBallotsDeletionBufferBlocks,
			updatedParams.PendingBallotsDeletionBufferBlocks,
		)
		require.Equal(
			t,
			int64(216000),
			updatedParams.PendingBallotsDeletionBufferBlocks,
		)

		require.True(t, found)
		require.Equal(t, currentParams.ValidatorEmissionPercentage, updatedParams.ValidatorEmissionPercentage)
		require.Equal(t, currentParams.ObserverEmissionPercentage, updatedParams.ObserverEmissionPercentage)
		require.Equal(t, currentParams.TssSignerEmissionPercentage, updatedParams.TssSignerEmissionPercentage)
		require.Equal(t, currentParams.ObserverSlashAmount, updatedParams.ObserverSlashAmount)

	})
}
