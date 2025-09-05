package v12_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	v12 "github.com/zeta-chain/node/x/observer/migrations/v12"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestMigrateStore(t *testing.T) {
	t.Run("can migrate stability pool percentage", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// Initialize and set test chain params
		testChainParams := sample.ChainParamsList()
		k.SetChainParamsList(ctx, testChainParams)

		// Verify initial setup
		oldChainParams, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Equal(t, testChainParams, oldChainParams)

		// Act
		err := v12.MigrateStore(ctx, *k)
		require.NoError(t, err)

		// Assert
		// Verify chain params were preserved
		newChainParams, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Len(t, newChainParams.ChainParams, len(oldChainParams.ChainParams))

		// Verify the StabilityPoolPercentage field is updated and nothing else changed
		for i, newParam := range newChainParams.ChainParams {
			oldParam := oldChainParams.ChainParams[i]

			// Check StabilityPoolPercentage is set to 60
			require.Equal(t, uint64(60), newParam.StabilityPoolPercentage)

			// Ensure nothing else has changed
			oldParam.StabilityPoolPercentage = 60
			require.Equal(t, newParam, oldParam)
		}
	})

	t.Run("migrate nothing if chain params not found", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// Verify no chain params exist initially
		allChainParams, found := k.GetChainParamsList(ctx)
		require.False(t, found)
		require.Empty(t, allChainParams.ChainParams)

		// Act
		err := v12.MigrateStore(ctx, *k)

		// Assert
		require.ErrorIs(t, err, types.ErrChainParamsNotFound)

		// Verify nothing has changed
		allChainParams, found = k.GetChainParamsList(ctx)
		require.False(t, found)
		require.Empty(t, allChainParams.ChainParams)
	})

	t.Run("migrate nothing if chain params list validation fails", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// Create invalid chain params
		testChainParams := sample.ChainParamsList()
		testChainParams.ChainParams[0].InboundTicker = 0 // Make first chain params invalid

		// Set up invalid chain params
		k.SetChainParamsList(ctx, testChainParams)

		// Verify initial setup
		oldChainParams, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Equal(t, testChainParams, oldChainParams)

		// Act
		err := v12.MigrateStore(ctx, *k)

		// Assert
		require.ErrorIs(t, err, types.ErrInvalidChainParams)

		// Verify nothing has changed
		newChainParams, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Equal(t, oldChainParams, newChainParams)
	})
}
