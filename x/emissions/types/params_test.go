package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestNewParams(t *testing.T) {
	params := NewParams()

	require.Equal(t, "00.50", params.ValidatorEmissionPercentage, "ValidatorEmissionPercentage should be set to 00.50")
	require.Equal(t, "00.25", params.ObserverEmissionPercentage, "ObserverEmissionPercentage should be set to 00.25")
	require.Equal(t, "00.25", params.TssSignerEmissionPercentage, "TssSignerEmissionPercentage should be set to 00.25")
	require.Equal(
		t,
		sdkmath.NewInt(100000000000000000),
		params.ObserverSlashAmount,
		"ObserverSlashAmount should be set to 100000000000000000",
	)
	require.Equal(t, int64(300), params.BallotMaturityBlocks, "BallotMaturityBlocks should be set to 300")
	require.Equal(t, BlockReward, params.BlockRewardAmount, "BlockRewardAmount should be set to 0")
}

func TestDefaultParams(t *testing.T) {
	params := DefaultParams()

	// The default parameters should match those set in NewParams
	require.Equal(t, NewParams(), params)
}

func TestValidateValidatorEmissionPercentage(t *testing.T) {
	require.Error(t, validateValidatorEmissionPercentage(0.5))
	require.NoError(t, validateValidatorEmissionPercentage("0.50"))
	require.Error(t, validateValidatorEmissionPercentage("-0.50")) // Less than 0 percent should fail
	require.Error(t, validateValidatorEmissionPercentage("1.01"))  // More than 100 percent should fail
}

func TestValidateObserverEmissionPercentage(t *testing.T) {
	require.Error(t, validateObserverEmissionPercentage(0.25))
	require.NoError(t, validateObserverEmissionPercentage("0.25"))
	require.Error(t, validateObserverEmissionPercentage("-0.50")) // Less than 0 percent should fail
	require.Error(t, validateObserverEmissionPercentage("1.01"))  // More than 100 percent should fail
}

func TestValidateTssEmissionPercentage(t *testing.T) {
	require.Error(t, validateTssEmissionPercentage(0.25))
	require.NoError(t, validateTssEmissionPercentage("0.25"))
	require.Error(t, validateTssEmissionPercentage("-0.25")) // Less than 0 percent should fail
	require.Error(t, validateTssEmissionPercentage("1.01"))  // More than 100 percent should fail
}

func TestValidateObserverSlashAmount(t *testing.T) {
	require.Error(t, validateObserverSlashAmount(10))
	require.Error(t, validateObserverSlashAmount("10"))
	require.Error(t, validateObserverSlashAmount(sdkmath.NewInt(-10))) // Less than 0
	require.Error(t, validateObserverSlashAmount(nil))
	require.NoError(t, validateObserverSlashAmount(sdkmath.NewInt(10)))
}

func TestValidateBallotMaturityBlocks(t *testing.T) {
	require.Error(t, validateBallotMaturityBlocks("10"))
	require.Error(t, validateBallotMaturityBlocks(-100))
	require.NoError(t, validateBallotMaturityBlocks(int64(100)))
}

func TestValidateBlockRewardAmount(t *testing.T) {
	require.Error(t, validateBlockRewardsAmount("0.50"))
	require.Error(t, validateBlockRewardsAmount("-0.50"))
	require.Error(t, validateBlockRewardsAmount(sdkmath.LegacyMustNewDecFromStr("-0.50")))
	require.Error(t, validateBlockRewardsAmount(nil))
	require.NoError(t, validateBlockRewardsAmount(sdkmath.LegacyMustNewDecFromStr("0.50")))
	require.NoError(t, validateBlockRewardsAmount(sdkmath.LegacyZeroDec()))
	require.NoError(t, validateBlockRewardsAmount(BlockReward))
}

func TestValidate(t *testing.T) {
	t.Run("should validate", func(t *testing.T) {
		params := NewParams()
		require.NoError(t, params.Validate())
	})

	t.Run("should error for invalid validator emissions percentage", func(t *testing.T) {
		params := NewParams()
		params.ValidatorEmissionPercentage = "-1.30"
		require.Error(t, params.Validate())
	})

	t.Run("should error for invalid observer emissions percentage", func(t *testing.T) {
		params := NewParams()
		params.ObserverEmissionPercentage = "-1.30"
		require.Error(t, params.Validate())
	})

	t.Run("should error for invalid tss emissions percentage", func(t *testing.T) {
		params := NewParams()
		params.TssSignerEmissionPercentage = "-1.30"
		require.Error(t, params.Validate())
	})

	t.Run("should error for invalid observer slash amount", func(t *testing.T) {
		params := NewParams()
		params.ObserverSlashAmount = sdkmath.NewInt(-10)
		require.Error(t, params.Validate())
	})

	t.Run("should error for invalid ballot maturity blocks", func(t *testing.T) {
		params := NewParams()
		params.BallotMaturityBlocks = -100
		require.Error(t, params.Validate())
	})

	t.Run("should error for negative block reward amount", func(t *testing.T) {
		params := NewParams()
		params.BlockRewardAmount = sdkmath.LegacyMustNewDecFromStr("-1.30")
		require.ErrorContains(t, params.Validate(), "block reward amount must not be negative")
	})

	t.Run("should error if pending ballots buffer blocks is negative", func(t *testing.T) {
		params := NewParams()
		params.PendingBallotsDeletionBufferBlocks = -100
		require.Error(t, params.Validate())
	})
}
func TestParamsString(t *testing.T) {
	params := DefaultParams()
	out, err := yaml.Marshal(params)
	require.NoError(t, err)
	require.Equal(t, string(out), params.String())
}
