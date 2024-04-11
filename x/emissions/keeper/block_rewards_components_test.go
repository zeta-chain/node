package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	emissionskeeper "github.com/zeta-chain/zetacore/x/emissions/keeper"
	emissionstypes "github.com/zeta-chain/zetacore/x/emissions/types"
)

func TestKeeper_CalculateFixedValidatorRewards(t *testing.T) {
	tt := []struct {
		name                 string
		blockTimeInSecs      string
		expectedBlockRewards sdk.Dec
		wantErr              bool
	}{
		{
			name:            "Invalid block time",
			blockTimeInSecs: "",
			wantErr:         true,
		},
		{
			name:                 "Block Time 5.7",
			blockTimeInSecs:      "5.7",
			expectedBlockRewards: sdk.MustNewDecFromStr("9620949074074074074.074070733466756687"),
		},
		{
			name:                 "Block Time 6",
			blockTimeInSecs:      "6",
			expectedBlockRewards: sdk.MustNewDecFromStr("10127314814814814814.814814814814814815"),
		},
		{
			name:                 "Block Time 3",
			blockTimeInSecs:      "3",
			expectedBlockRewards: sdk.MustNewDecFromStr("5063657407407407407.407407407407407407"),
		},
		{
			name:                 "Block Time 2",
			blockTimeInSecs:      "2",
			expectedBlockRewards: sdk.MustNewDecFromStr("3375771604938271604.938271604938271605"),
		},
		{
			name:                 "Block Time 8",
			blockTimeInSecs:      "8",
			expectedBlockRewards: sdk.MustNewDecFromStr("13503086419753086419.753086419753086420"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			blockRewards, err := emissionskeeper.CalculateFixedValidatorRewards(tc.blockTimeInSecs)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expectedBlockRewards, blockRewards)
		})
	}
}

func TestKeeper_GetFixedBlockRewards(t *testing.T) {
	k, _, _, _ := keepertest.EmissionsKeeper(t)
	fixedBlockRewards, err := k.GetFixedBlockRewards()
	require.NoError(t, err)
	require.Equal(t, emissionstypes.BlockReward, fixedBlockRewards)
}

func TestKeeper_GetBlockRewardComponent(t *testing.T) {
	t.Run("should return all 0s if reserves factor is 0", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionKeeperWithMockOptions(t, keepertest.EmissionMockOptions{
			UseBankMock: true,
		})

		bankMock := keepertest.GetEmissionsBankMock(t, k)
		bankMock.On("GetBalance",
			ctx, mock.Anything, config.BaseDenom).
			Return(sdk.NewCoin(config.BaseDenom, math.NewInt(0)), nil).Once()

		reservesFactor, bondFactor, durationFactor := k.GetBlockRewardComponents(ctx)
		require.Equal(t, sdk.ZeroDec(), reservesFactor)
		require.Equal(t, sdk.ZeroDec(), bondFactor)
		require.Equal(t, sdk.ZeroDec(), durationFactor)
	})

	t.Run("should return if reserves factor is not 0", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionKeeperWithMockOptions(t, keepertest.EmissionMockOptions{
			UseBankMock: true,
		})

		bankMock := keepertest.GetEmissionsBankMock(t, k)
		bankMock.On("GetBalance",
			ctx, mock.Anything, config.BaseDenom).
			Return(sdk.NewCoin(config.BaseDenom, math.NewInt(1)), nil).Once()

		reservesFactor, bondFactor, durationFactor := k.GetBlockRewardComponents(ctx)
		require.Equal(t, sdk.OneDec(), reservesFactor)
		// bonded ratio is 0
		require.Equal(t, sdk.ZeroDec(), bondFactor)
		// non 0 value returned
		require.NotEqual(t, sdk.ZeroDec(), durationFactor)
		require.Positive(t, durationFactor.BigInt().Int64())
	})
}

func TestKeeper_GetBondFactor(t *testing.T) {
	t.Run("should return 0 if current bond ratio is 0", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)

		bondFactor := k.GetBondFactor(ctx, k.GetStakingKeeper())
		require.Equal(t, sdk.ZeroDec(), bondFactor)
	})

	t.Run("should return max bond factor if bond factor exceeds max bond factor", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionKeeperWithMockOptions(t, keepertest.EmissionMockOptions{
			UseStakingMock: true,
		})

		params := emissionstypes.DefaultParams()
		params.TargetBondRatio = "0.5"
		params.MaxBondFactor = "1.1"
		params.MinBondFactor = "0.9"
		err := k.SetParams(ctx, params)
		require.NoError(t, err)

		stakingMock := keepertest.GetEmissionsStakingMock(t, k)
		stakingMock.On("BondedRatio", ctx).Return(sdk.MustNewDecFromStr("0.25"))
		bondFactor := k.GetBondFactor(ctx, k.GetStakingKeeper())
		require.Equal(t, sdk.MustNewDecFromStr(params.MaxBondFactor), bondFactor)
	})

	t.Run("should return min bond factor if bond factor below min bond factor", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionKeeperWithMockOptions(t, keepertest.EmissionMockOptions{
			UseStakingMock: true,
		})

		params := emissionstypes.DefaultParams()
		params.TargetBondRatio = "0.5"
		params.MaxBondFactor = "1.1"
		params.MinBondFactor = "0.9"
		err := k.SetParams(ctx, params)
		require.NoError(t, err)

		stakingMock := keepertest.GetEmissionsStakingMock(t, k)
		stakingMock.On("BondedRatio", ctx).Return(sdk.MustNewDecFromStr("0.75"))
		bondFactor := k.GetBondFactor(ctx, k.GetStakingKeeper())
		require.Equal(t, sdk.MustNewDecFromStr(params.MinBondFactor), bondFactor)
	})

	t.Run("should return calculated bond factor if bond factor in range", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionKeeperWithMockOptions(t, keepertest.EmissionMockOptions{
			UseStakingMock: true,
		})

		params := emissionstypes.DefaultParams()
		params.TargetBondRatio = "0.5"
		params.MaxBondFactor = "1.1"
		params.MinBondFactor = "0.9"
		err := k.SetParams(ctx, params)
		require.NoError(t, err)

		stakingMock := keepertest.GetEmissionsStakingMock(t, k)
		stakingMock.On("BondedRatio", ctx).Return(sdk.MustNewDecFromStr("0.5"))
		bondFactor := k.GetBondFactor(ctx, k.GetStakingKeeper())
		require.Equal(t, sdk.OneDec(), bondFactor)
	})
}

func TestKeeper_GetDurationFactor(t *testing.T) {
	t.Run("should return duration factor 0 if duration factor constant is 0", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		params := emissionstypes.DefaultParams()
		params.DurationFactorConstant = "0"
		err := k.SetParams(ctx, params)
		require.NoError(t, err)

		duractionFactor := k.GetDurationFactor(ctx)
		require.Equal(t, sdk.ZeroDec(), duractionFactor)
	})

	t.Run("should return duration factor for default params", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		duractionFactor := k.GetDurationFactor(ctx)
		// hardcoding actual expected value for default params, it will change if logic changes
		require.Equal(t, sdk.MustNewDecFromStr("0.000000004346937374"), duractionFactor)
	})
}
