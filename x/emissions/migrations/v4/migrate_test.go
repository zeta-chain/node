package v4_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	v4 "github.com/zeta-chain/node/x/emissions/migrations/v4"
	"github.com/zeta-chain/node/x/emissions/types"
)

func TestMigrateStore(t *testing.T) {
	t.Run("should successfully migrate to new params in mainnet", func(t *testing.T) {
		//Arrange
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		cdc := k.GetCodec()
		emissionsStoreKey := k.GetStoreKey()
		mainnetParams := LegacyMainnetParams()
		err := SetLegacyParams(ctx, emissionsStoreKey, cdc, mainnetParams)
		require.NoError(t, err)

		//Act
		err = v4.MigrateStore(ctx, k)
		require.NoError(t, err)

		//Assert
		params, found := k.GetParams(ctx)
		require.True(t, found)
		require.Equal(t, mainnetParams.ValidatorEmissionPercentage, params.ValidatorEmissionPercentage)
		require.Equal(t, mainnetParams.ObserverEmissionPercentage, params.ObserverEmissionPercentage)
		require.Equal(t, mainnetParams.TssSignerEmissionPercentage, params.TssSignerEmissionPercentage)
		require.Equal(t, mainnetParams.ObserverSlashAmount, params.ObserverSlashAmount)
		require.Equal(t, mainnetParams.BallotMaturityBlocks, params.BallotMaturityBlocks)
		require.Equal(t, types.BlockReward, params.BlockRewardAmount)
	})

	t.Run("should successfully migrate to new params in testnet", func(t *testing.T) {
		//Arrange
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		cdc := k.GetCodec()
		emissionsStoreKey := k.GetStoreKey()
		testNetParams := LegacyTestNetParams()
		err := SetLegacyParams(ctx, emissionsStoreKey, cdc, testNetParams)
		require.NoError(t, err)

		//Act
		err = v4.MigrateStore(ctx, k)
		require.NoError(t, err)

		//Assert
		params, found := k.GetParams(ctx)
		require.True(t, found)
		require.Equal(t, testNetParams.ValidatorEmissionPercentage, params.ValidatorEmissionPercentage)
		require.Equal(t, testNetParams.ObserverEmissionPercentage, params.ObserverEmissionPercentage)
		require.Equal(t, testNetParams.TssSignerEmissionPercentage, params.TssSignerEmissionPercentage)
		require.Equal(t, testNetParams.ObserverSlashAmount, params.ObserverSlashAmount)
		require.Equal(t, testNetParams.BallotMaturityBlocks, params.BallotMaturityBlocks)
		require.Equal(t, types.BlockReward, params.BlockRewardAmount)
	})

	t.Run(
		"should successfully migrate using default values if legacy param for ValidatorEmissionPercentage is not available",
		func(t *testing.T) {
			//Arrange
			k, ctx, _, _ := keepertest.EmissionsKeeper(t)
			cdc := k.GetCodec()
			emissionsStoreKey := k.GetStoreKey()

			mainnetParams := LegacyMainnetParams()
			mainnetParams.ValidatorEmissionPercentage = ""
			err := SetLegacyParams(ctx, emissionsStoreKey, cdc, mainnetParams)
			require.NoError(t, err)

			//Act
			err = v4.MigrateStore(ctx, k)
			require.NoError(t, err)

			//Assert
			defaultParams := types.DefaultParams()
			params, found := k.GetParams(ctx)
			require.True(t, found)
			require.Equal(t, defaultParams.ValidatorEmissionPercentage, params.ValidatorEmissionPercentage)
			require.Equal(t, mainnetParams.ObserverEmissionPercentage, params.ObserverEmissionPercentage)
			require.Equal(t, mainnetParams.TssSignerEmissionPercentage, params.TssSignerEmissionPercentage)
			require.Equal(t, mainnetParams.ObserverSlashAmount, params.ObserverSlashAmount)
			require.Equal(t, mainnetParams.BallotMaturityBlocks, params.BallotMaturityBlocks)
			require.Equal(t, types.BlockReward, params.BlockRewardAmount)
		},
	)

	t.Run(
		"should successfully migrate using default values if legacy param for ObserverEmissionPercentage is not available",
		func(t *testing.T) {
			//Arrange
			k, ctx, _, _ := keepertest.EmissionsKeeper(t)
			cdc := k.GetCodec()
			emissionsStoreKey := k.GetStoreKey()

			mainnetParams := LegacyMainnetParams()
			mainnetParams.ObserverEmissionPercentage = ""
			err := SetLegacyParams(ctx, emissionsStoreKey, cdc, mainnetParams)
			require.NoError(t, err)

			//Act
			err = v4.MigrateStore(ctx, k)
			require.NoError(t, err)

			//Assert
			defaultParams := types.DefaultParams()
			params, found := k.GetParams(ctx)
			require.True(t, found)
			require.Equal(t, mainnetParams.ValidatorEmissionPercentage, params.ValidatorEmissionPercentage)
			require.Equal(t, defaultParams.ObserverEmissionPercentage, params.ObserverEmissionPercentage)
			require.Equal(t, mainnetParams.TssSignerEmissionPercentage, params.TssSignerEmissionPercentage)
			require.Equal(t, mainnetParams.ObserverSlashAmount, params.ObserverSlashAmount)
			require.Equal(t, mainnetParams.BallotMaturityBlocks, params.BallotMaturityBlocks)
			require.Equal(t, types.BlockReward, params.BlockRewardAmount)
		},
	)

	t.Run(
		"should successfully migrate using default values if legacy param for TssSignerEmissionPercentage is not available",
		func(t *testing.T) {
			//Arrange
			k, ctx, _, _ := keepertest.EmissionsKeeper(t)
			cdc := k.GetCodec()
			emissionsStoreKey := k.GetStoreKey()

			mainnetParams := LegacyMainnetParams()
			mainnetParams.TssSignerEmissionPercentage = ""
			err := SetLegacyParams(ctx, emissionsStoreKey, cdc, mainnetParams)
			require.NoError(t, err)

			//Act
			err = v4.MigrateStore(ctx, k)
			require.NoError(t, err)

			//Assert
			defaultParams := types.DefaultParams()
			params, found := k.GetParams(ctx)
			require.True(t, found)
			require.Equal(t, mainnetParams.ValidatorEmissionPercentage, params.ValidatorEmissionPercentage)
			require.Equal(t, mainnetParams.ObserverEmissionPercentage, params.ObserverEmissionPercentage)
			require.Equal(t, defaultParams.TssSignerEmissionPercentage, params.TssSignerEmissionPercentage)
			require.Equal(t, mainnetParams.ObserverSlashAmount, params.ObserverSlashAmount)
			require.Equal(t, mainnetParams.BallotMaturityBlocks, params.BallotMaturityBlocks)
			require.Equal(t, types.BlockReward, params.BlockRewardAmount)
		},
	)

	t.Run(
		"should successfully migrate using default values if legacy param for ObserverSlashAmount is not available",
		func(t *testing.T) {
			//Arrange
			k, ctx, _, _ := keepertest.EmissionsKeeper(t)
			cdc := k.GetCodec()
			emissionsStoreKey := k.GetStoreKey()

			mainnetParams := LegacyMainnetParams()
			mainnetParams.ObserverSlashAmount = sdkmath.NewInt(-1)
			err := SetLegacyParams(ctx, emissionsStoreKey, cdc, mainnetParams)
			require.NoError(t, err)

			//Act
			err = v4.MigrateStore(ctx, k)
			require.NoError(t, err)

			//Assert
			defaultParams := types.DefaultParams()
			params, found := k.GetParams(ctx)
			require.True(t, found)
			require.Equal(t, mainnetParams.ValidatorEmissionPercentage, params.ValidatorEmissionPercentage)
			require.Equal(t, mainnetParams.ObserverEmissionPercentage, params.ObserverEmissionPercentage)
			require.Equal(t, mainnetParams.TssSignerEmissionPercentage, params.TssSignerEmissionPercentage)
			require.Equal(t, defaultParams.ObserverSlashAmount.String(), params.ObserverSlashAmount.String())
			require.Equal(t, mainnetParams.BallotMaturityBlocks, params.BallotMaturityBlocks)
			require.Equal(t, types.BlockReward, params.BlockRewardAmount)
		},
	)

	t.Run(
		"should successfully migrate using default values if legacy param for BallotMaturityBlocks is not available",
		func(t *testing.T) {
			//Arrange
			k, ctx, _, _ := keepertest.EmissionsKeeper(t)
			cdc := k.GetCodec()
			emissionsStoreKey := k.GetStoreKey()

			mainnetParams := LegacyMainnetParams()
			mainnetParams.BallotMaturityBlocks = -1
			err := SetLegacyParams(ctx, emissionsStoreKey, cdc, mainnetParams)
			require.NoError(t, err)

			//Act
			err = v4.MigrateStore(ctx, k)
			require.NoError(t, err)

			//Assert
			defaultParams := types.DefaultParams()
			params, found := k.GetParams(ctx)
			require.True(t, found)
			require.Equal(t, mainnetParams.ValidatorEmissionPercentage, params.ValidatorEmissionPercentage)
			require.Equal(t, mainnetParams.ObserverEmissionPercentage, params.ObserverEmissionPercentage)
			require.Equal(t, mainnetParams.TssSignerEmissionPercentage, params.TssSignerEmissionPercentage)
			require.Equal(t, mainnetParams.ObserverSlashAmount, params.ObserverSlashAmount)
			require.Equal(t, defaultParams.BallotMaturityBlocks, params.BallotMaturityBlocks)
			require.Equal(t, types.BlockReward, params.BlockRewardAmount)
		},
	)

	t.Run("fail to migrate if legacy params are not found", func(t *testing.T) {
		//Arrange
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		store := ctx.KVStore(k.GetStoreKey())
		store.Delete(types.KeyPrefix(types.ParamsKey))

		//Act
		err := v4.MigrateStore(ctx, k)

		//Assert
		require.ErrorIs(t, err, types.ErrMigrationFailed)
		require.ErrorContains(t, err, "failed to get legacy params")
	})

	// This scenario is hypothetical as the legacy params have valid values.
	t.Run("fail to migrate if params are not valid", func(t *testing.T) {
		//Arrange
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		cdc := k.GetCodec()
		emissionsStoreKey := k.GetStoreKey()
		mainnetParams := LegacyMainnetParams()
		mainnetParams.TssSignerEmissionPercentage = "2.0"
		err := SetLegacyParams(ctx, emissionsStoreKey, cdc, mainnetParams)
		require.NoError(t, err)

		//Act
		err = v4.MigrateStore(ctx, k)

		//Assert
		require.ErrorIs(t, err, types.ErrMigrationFailed)
		require.ErrorContains(t, err, "tss emission percentage cannot be more than 100 percent")
	})
}

func TestGetParamsLegacy(t *testing.T) {
	t.Run("should successfully get legacy params", func(t *testing.T) {
		//Arrange
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		cdc := k.GetCodec()
		emissionsStoreKey := k.GetStoreKey()
		mainnetParams := LegacyMainnetParams()
		err := SetLegacyParams(ctx, emissionsStoreKey, cdc, mainnetParams)
		require.NoError(t, err)

		//Act
		params, found := v4.GetParamsLegacy(ctx, emissionsStoreKey, cdc)

		//Assert
		require.True(t, found)
		require.Equal(t, mainnetParams, params)
	})

	t.Run("should return false if legacy params are not found", func(t *testing.T) {
		//Arrange
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		store := ctx.KVStore(k.GetStoreKey())
		store.Delete(types.KeyPrefix(types.ParamsKey))

		//Act
		params, found := v4.GetParamsLegacy(ctx, k.GetStoreKey(), k.GetCodec())

		//Assert
		require.False(t, found)
		require.Equal(t, types.LegacyParams{}, params)
	})

	t.Run("should return false if unable to unmarshal legacy params", func(t *testing.T) {
		//Arrange
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		cdc := k.GetCodec()
		emissionsStoreKey := k.GetStoreKey()
		store := ctx.KVStore(emissionsStoreKey)
		store.Set(types.KeyPrefix(types.ParamsKey), []byte{0x00})

		//Act
		params, found := v4.GetParamsLegacy(ctx, emissionsStoreKey, cdc)

		//Assert
		require.False(t, found)
		require.Equal(t, types.LegacyParams{}, params)
	})
}

func SetLegacyParams(
	ctx sdk.Context,
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	params types.LegacyParams,
) error {
	store := ctx.KVStore(storeKey)
	bz, err := cdc.Marshal(&params)
	if err != nil {
		return err
	}

	store.Set(types.KeyPrefix(types.ParamsKey), bz)
	return nil
}

// https://zetachain-api.lavenderfive.com/zeta-chain/emissions/params
func LegacyMainnetParams() types.LegacyParams {
	return types.LegacyParams{
		MaxBondFactor:               "1.25",
		MinBondFactor:               "0.75",
		AvgBlockTime:                "6.00",
		TargetBondRatio:             "0.67",
		ObserverEmissionPercentage:  "0.125",
		ValidatorEmissionPercentage: "0.75",
		TssSignerEmissionPercentage: "0.125",
		DurationFactorConstant:      "0.001877876953694702",
		ObserverSlashAmount:         sdkmath.NewIntFromUint64(100000000000000000),
		BallotMaturityBlocks:        100,
	}
}

// https://zetachain-testnet-api.itrocket.net/zeta-chain/emissions/params
func LegacyTestNetParams() types.LegacyParams {
	return types.LegacyParams{
		MaxBondFactor:               "1.25",
		MinBondFactor:               "0.75",
		AvgBlockTime:                "6.00",
		TargetBondRatio:             "0.67",
		ObserverEmissionPercentage:  "0.05",
		ValidatorEmissionPercentage: "0.90",
		TssSignerEmissionPercentage: "0.05",
		DurationFactorConstant:      "0.001877876953694702",
		ObserverSlashAmount:         sdkmath.NewIntFromUint64(100000000000000000),
		BallotMaturityBlocks:        100,
	}
}
