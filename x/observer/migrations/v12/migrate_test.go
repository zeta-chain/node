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
	t.Run("can migrate chain params and crosschain flags", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// set chain params
		testChainParams := getTestChainParams()
		k.SetChainParamsList(ctx, testChainParams)

		// set crosschain flags with V2 ZETA enabled
		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled:  true,
			IsOutboundEnabled: true,
			IsV2ZetaEnabled:   true,
		})

		// ACT
		err := v12.MigrateStore(ctx, *k)

		// ASSERT
		require.NoError(t, err)

		// verify chain params were migrated
		newChainParams, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Len(t, newChainParams.ChainParams, len(testChainParams.ChainParams))
		for _, param := range newChainParams.ChainParams {
			require.Equal(t, uint64(100), param.StabilityPoolPercentage)
			require.True(t, param.GasPriceMultiplier.IsPositive())
		}

		// verify V2 ZETA flows were disabled
		flags, found := k.GetCrosschainFlags(ctx)
		require.True(t, found)
		require.False(t, flags.IsV2ZetaEnabled)
		require.True(t, flags.IsInboundEnabled)
		require.True(t, flags.IsOutboundEnabled)
	})

	t.Run("returns error when chain params not found", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// don't set any chain params

		// ACT
		err := v12.MigrateStore(ctx, *k)

		// ASSERT
		require.ErrorIs(t, err, types.ErrChainParamsNotFound)
	})

	t.Run("returns error when chain not found in chain params", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// set chain params with invalid chain ID
		testChainParams := getTestChainParams()
		testChainParams.ChainParams[0].ChainId = 999999999
		k.SetChainParamsList(ctx, testChainParams)

		// ACT
		err := v12.MigrateStore(ctx, *k)

		// ASSERT
		require.ErrorIs(t, err, types.ErrSupportedChains)
	})

	t.Run("returns error when chain params validation fails", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// set chain params with invalid values that will fail validation
		testChainParams := getTestChainParams()
		testChainParams.ChainParams[0].InboundTicker = 0
		k.SetChainParamsList(ctx, testChainParams)

		// ACT
		err := v12.MigrateStore(ctx, *k)

		// ASSERT
		require.ErrorIs(t, err, types.ErrInvalidChainParams)
	})
}

func TestUpdateChainParams(t *testing.T) {
	t.Run("can update gas price multiplier and stability pool percentage", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// set chain params
		testChainParams := getTestChainParams()
		k.SetChainParamsList(ctx, testChainParams)

		// ensure the chain params are set correctly
		oldChainParams, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Equal(t, testChainParams, oldChainParams)

		// migrate the store
		err := v12.UpdateChainParams(ctx, *k)
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

	t.Run("returns error if chain params not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		err := v12.UpdateChainParams(ctx, *k)

		require.ErrorIs(t, err, types.ErrChainParamsNotFound)
	})

	t.Run("returns error if chain not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		testChainParams := getTestChainParams()
		testChainParams.ChainParams[0].ChainId = 1000000000
		k.SetChainParamsList(ctx, testChainParams)

		err := v12.UpdateChainParams(ctx, *k)

		require.ErrorIs(t, err, types.ErrSupportedChains)
		require.ErrorContains(t, err, "chain 1000000000 not found")
	})

	t.Run("returns error if chain params validation fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		testChainParams := getTestChainParams()
		testChainParams.ChainParams[0].InboundTicker = 0
		k.SetChainParamsList(ctx, testChainParams)

		err := v12.UpdateChainParams(ctx, *k)

		require.ErrorIs(t, err, types.ErrInvalidChainParams)
	})
}

func TestUpdateCrosschainFlags(t *testing.T) {
	t.Run("disables V2 ZETA flows when flags exist", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// set crosschain flags with V2 ZETA enabled
		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled:  true,
			IsOutboundEnabled: true,
			IsV2ZetaEnabled:   true,
		})

		// ACT
		v12.UpdateCrosschainFlags(ctx, *k)

		// ASSERT
		flags, found := k.GetCrosschainFlags(ctx)
		require.True(t, found)
		require.False(t, flags.IsV2ZetaEnabled)
		// other flags should be preserved
		require.True(t, flags.IsInboundEnabled)
		require.True(t, flags.IsOutboundEnabled)
	})

	t.Run("sets default flags when not found", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		// ensure no flags are set
		_, found := k.GetCrosschainFlags(ctx)
		require.False(t, found)

		// ACT
		v12.UpdateCrosschainFlags(ctx, *k)

		// ASSERT
		flags, found := k.GetCrosschainFlags(ctx)
		require.True(t, found)
		require.False(t, flags.IsV2ZetaEnabled)
		// verify default flags are used
		require.True(t, flags.IsInboundEnabled)
		require.True(t, flags.IsOutboundEnabled)
		require.NotNil(t, flags.GasPriceIncreaseFlags)
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
