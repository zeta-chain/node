package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestNewParams(t *testing.T) {
	params := NewParams()

	// Verifying all parameters to ensure they are set correctly.
	require.Equal(t, "1.25", params.MaxBondFactor, "MaxBondFactor should be set to 1.25")
	require.Equal(t, "0.75", params.MinBondFactor, "MinBondFactor should be set to 0.75")
	require.Equal(t, "6.00", params.AvgBlockTime, "AvgBlockTime should be set to 6.00")
	require.Equal(t, "00.67", params.TargetBondRatio, "TargetBondRatio should be set to 00.67")
	require.Equal(t, "00.50", params.ValidatorEmissionPercentage, "ValidatorEmissionPercentage should be set to 00.50")
	require.Equal(t, "00.25", params.ObserverEmissionPercentage, "ObserverEmissionPercentage should be set to 00.25")
	require.Equal(t, "00.25", params.TssSignerEmissionPercentage, "TssSignerEmissionPercentage should be set to 00.25")
	require.Equal(t, "0.001877876953694702", params.DurationFactorConstant, "DurationFactorConstant should be set to 0.001877876953694702")

	require.Equal(t, sdkmath.NewInt(100000000000000000), params.ObserverSlashAmount, "ObserverSlashAmount should be set to 100000000000000000")
}

func TestDefaultParams(t *testing.T) {
	params := DefaultParams()

	// The default parameters should match those set in NewParams
	require.Equal(t, NewParams(), params)
}

func TestValidateDurationFactorConstant(t *testing.T) {
	require.NoError(t, validateDurationFactorConstant("1"))
	require.Error(t, validateDurationFactorConstant(1))
}

func TestValidateMaxBondFactor(t *testing.T) {
	require.Error(t, validateMaxBondFactor(1))
	require.NoError(t, validateMaxBondFactor("1.00"))
	require.NoError(t, validateMaxBondFactor("1.25"))
	require.Error(t, validateMaxBondFactor("1.30")) // Should fail as it's higher than 1.25
}

func TestValidateMinBondFactor(t *testing.T) {
	require.Error(t, validateMinBondFactor(1))
	require.NoError(t, validateMinBondFactor("0.75"))
	require.Error(t, validateMinBondFactor("0.50")) // Should fail as it's lower than 0.75
}

func TestValidateAvgBlockTime(t *testing.T) {
	require.Error(t, validateAvgBlockTime(6))
	require.Error(t, validateAvgBlockTime("invalid"))
	require.NoError(t, validateAvgBlockTime("6.00"))
	require.Error(t, validateAvgBlockTime("-1")) // Negative time should fail
	require.Error(t, validateAvgBlockTime("0"))  // Zero should also fail
}

func TestValidateTargetBondRatio(t *testing.T) {
	require.Error(t, validateTargetBondRatio(0.5))
	require.NoError(t, validateTargetBondRatio("0.50"))
	require.Error(t, validateTargetBondRatio("-0.01")) // Less than 0 percent should fail
	require.Error(t, validateTargetBondRatio("1.01"))  // More than 100 percent should fail
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

func TestParamsString(t *testing.T) {
	params := DefaultParams()
	out, err := yaml.Marshal(params)
	require.NoError(t, err)
	require.Equal(t, string(out), params.String())
}
