package v7_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	v7 "github.com/zeta-chain/node/x/emissions/migrations/v7"
	"github.com/zeta-chain/node/x/emissions/types"
)

func TestMigrateStore(t *testing.T) {
	t.Run("successfully migrate store", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)

		// Set up pre-migration (v6) values
		currentParams := types.DefaultParams()
		currentParams.BallotMaturityBlocks = 150
		currentParams.PendingBallotsDeletionBufferBlocks = 216000
		currentParams.BlockRewardAmount = sdkmath.LegacyMustNewDecFromStr("6751543209876543209.876543209876543210")
		k.SetParams(ctx, currentParams)

		// Act
		err := v7.MigrateStore(ctx, k)

		// Assert
		require.NoError(t, err)
		updatedParams, found := k.GetParams(ctx)
		require.True(t, found)

		// Verify migrated values changed from v6 to v7
		require.NotEqual(t, currentParams.BallotMaturityBlocks, updatedParams.BallotMaturityBlocks)
		require.Equal(t, int64(300), updatedParams.BallotMaturityBlocks)

		require.NotEqual(t, currentParams.BlockRewardAmount, updatedParams.BlockRewardAmount)
		require.Equal(
			t,
			sdkmath.LegacyMustNewDecFromStr("3375771604938271604.938271604938271605"),
			updatedParams.BlockRewardAmount,
		)

		require.NotEqual(
			t,
			currentParams.PendingBallotsDeletionBufferBlocks,
			updatedParams.PendingBallotsDeletionBufferBlocks,
		)
		require.Equal(
			t,
			int64(432000),
			updatedParams.PendingBallotsDeletionBufferBlocks,
		)

		// Verify unchanged params remain the same
		require.Equal(t, currentParams.ValidatorEmissionPercentage, updatedParams.ValidatorEmissionPercentage)
		require.Equal(t, currentParams.ObserverEmissionPercentage, updatedParams.ObserverEmissionPercentage)
		require.Equal(t, currentParams.TssSignerEmissionPercentage, updatedParams.TssSignerEmissionPercentage)
		require.Equal(t, currentParams.ObserverSlashAmount, updatedParams.ObserverSlashAmount)
	})
}
