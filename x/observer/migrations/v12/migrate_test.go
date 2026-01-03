package v12_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	v12 "github.com/zeta-chain/node/x/observer/migrations/v12"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestMigrateStore(t *testing.T) {
	t.Run("can migrate gas price multiplier and stability pool percentage", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// set chain params
		testChainParams := getTestChainParams()
		k.SetChainParamsList(ctx, testChainParams)

		// ensure the chain params are set correctly
		oldChainParams, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Equal(t, testChainParams, oldChainParams)

		// migrate the store
		err := v12.MigrateStore(ctx, *k)
		require.NoError(t, err)

		// ensure we still have same number of chain params after migration
		newChainParams, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Len(t, newChainParams.ChainParams, len(oldChainParams.ChainParams))

		// compare the old and new chain params
		for i, newParam := range newChainParams.ChainParams {
			oldParam := oldChainParams.ChainParams[i]

			// get the chain
			chain, found := chains.GetChainFromChainID(newParam.ChainId, []chains.Chain{})
			require.True(t, found)

			// get the gas price multiplier for the chain
			gasPriceMultiplier := v12.GetGasPriceMultiplierForChain(chain)

			// ensure the gas price multiplier is set correctly
			require.True(t, gasPriceMultiplier.IsPositive())
			require.True(t, gasPriceMultiplier.Equal(newParam.GasPriceMultiplier))

			// ensure the stability pool percentage is set correctly
			require.Equal(t, uint64(100), newParam.StabilityPoolPercentage)

			// ensure nothing else has changed except the migrated fields
			oldParam.GasPriceMultiplier = gasPriceMultiplier
			oldParam.StabilityPoolPercentage = 100
			require.True(t, types.ChainParamsEqual(*oldParam, *newParam))
		}
	})

	t.Run("migrate nothing if chain params not found", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// ensure no chain params are set
		allChainParams, found := k.GetChainParamsList(ctx)
		require.False(t, found)
		require.Empty(t, allChainParams.ChainParams)

		// migrate the store
		err := v12.MigrateStore(ctx, *k)

		// Assert
		require.ErrorIs(t, err, types.ErrChainParamsNotFound)

		// ensure nothing has changed
		allChainParams, found = k.GetChainParamsList(ctx)
		require.False(t, found)
		require.Empty(t, allChainParams.ChainParams)
	})

	t.Run("migrate nothing if chain not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// set chain params
		testChainParams := getTestChainParams()

		// change the first chain ID to unknown
		testChainParams.ChainParams[0].ChainId = 1000000000

		// set chain params
		k.SetChainParamsList(ctx, testChainParams)

		// ensure the chain params are set correctly
		oldChainParams, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Equal(t, testChainParams, oldChainParams)

		// migrate the store
		err := v12.MigrateStore(ctx, *k)
		require.ErrorIs(t, err, types.ErrSupportedChains)
		require.ErrorContains(t, err, "chain 1000000000 not found")

		// ensure nothing has changed
		newChainParams, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Equal(t, oldChainParams, newChainParams)
	})

	t.Run("migrate nothing if chain params list validation fails", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// get test chain params
		testChainParams := getTestChainParams()

		// make the first chain params invalid
		testChainParams.ChainParams[0].InboundTicker = 0

		// set chain params
		k.SetChainParamsList(ctx, testChainParams)

		// ensure the chain params are set correctly
		oldChainParams, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Equal(t, testChainParams, oldChainParams)

		// migrate the store
		err := v12.MigrateStore(ctx, *k)

		// Assert
		require.ErrorIs(t, err, types.ErrInvalidChainParams)

		// ensure nothing has changed
		newChainParams, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Equal(t, oldChainParams, newChainParams)
	})

}

// getTestChainParams returns a list of chain params for testing
func getTestChainParams() types.ChainParamsList {
	return types.ChainParamsList{
		ChainParams: []*types.ChainParams{
			makeChainParamsZeroGasMultiplier(1),
			makeChainParamsZeroGasMultiplier(56),
			makeChainParamsZeroGasMultiplier(8332),
			makeChainParamsZeroGasMultiplier(7000),
			makeChainParamsZeroGasMultiplier(137),
			makeChainParamsZeroGasMultiplier(8453),
			makeChainParamsZeroGasMultiplier(900),
			makeChainParamsZeroGasMultiplier(42161),
			makeChainParamsZeroGasMultiplier(43114),
			makeChainParamsZeroGasMultiplier(105),
			makeChainParamsZeroGasMultiplier(2015140),
		},
	}
}

// makeChainParamsZeroGasMultiplier creates a sample chain params with zero gas price multiplier
func makeChainParamsZeroGasMultiplier(chainID int64) *types.ChainParams {
	chainParams := sample.ChainParams(chainID)
	chainParams.GasPriceMultiplier = sdkmath.LegacyZeroDec()
	return chainParams
}
