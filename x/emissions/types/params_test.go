package types

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
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

	require.Equal(t, sdk.Int{}, params.ObserverSlashAmount, "ObserverSlashAmount should be initialized but is currently disabled")
}

func TestParamKeyTable(t *testing.T) {
	kt := ParamKeyTable()

	ps := Params{}
	for _, psp := range ps.ParamSetPairs() {
		require.PanicsWithValue(t, "duplicate parameter key", func() {
			kt.RegisterType(psp)
		})
	}
}

func TestDefaultParams(t *testing.T) {
	params := DefaultParams()

	// The default parameters should match those set in NewParams
	require.Equal(t, NewParams(), params)
}

func TestParamSetPairs(t *testing.T) {
	params := DefaultParams()
	pairs := params.ParamSetPairs()

	require.Equal(t, 8, len(pairs), "The number of param set pairs should match the expected count")

	assertParamSetPair(t, pairs, KeyPrefix(ParamMaxBondFactor), "1.25", validateMaxBondFactor)
	assertParamSetPair(t, pairs, KeyPrefix(ParamMinBondFactor), "0.75", validateMinBondFactor)
	assertParamSetPair(t, pairs, KeyPrefix(ParamAvgBlockTime), "6.00", validateAvgBlockTime)
	assertParamSetPair(t, pairs, KeyPrefix(ParamTargetBondRatio), "00.67", validateTargetBondRatio)
	assertParamSetPair(t, pairs, KeyPrefix(ParamValidatorEmissionPercentage), "00.50", validateValidatorEmissionPercentage)
	assertParamSetPair(t, pairs, KeyPrefix(ParamObserverEmissionPercentage), "00.25", validateObserverEmissionPercentage)
	assertParamSetPair(t, pairs, KeyPrefix(ParamTssSignerEmissionPercentage), "00.25", validateTssEmissonPercentage)
	assertParamSetPair(t, pairs, KeyPrefix(ParamDurationFactorConstant), "0.001877876953694702", validateDurationFactorConstant)
}

func assertParamSetPair(t *testing.T, pairs paramtypes.ParamSetPairs, key []byte, expectedValue string, valFunc paramtypes.ValueValidatorFn) {
	for _, pair := range pairs {
		if string(pair.Key) == string(key) {
			actualValue, ok := pair.Value.(*string)
			require.True(t, ok, "Expected value to be of type *string for key %s", string(key))
			require.Equal(t, expectedValue, *actualValue, "Value does not match for key %s", string(key))

			actualValFunc := pair.ValidatorFn
			require.Equal(t, reflect.ValueOf(valFunc).Pointer(), reflect.ValueOf(actualValFunc).Pointer(), "Val func doesnt match for key %s", string(key))
			return
		}
	}

	t.Errorf("Key %s not found in ParamSetPairs", string(key))
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
	require.Error(t, validateTargetBondRatio("1.01"))  // More than 100 percent should fail
	require.Error(t, validateTargetBondRatio("-0.01")) // Less than 0 percent should fail
}

func TestValidateValidatorEmissionPercentage(t *testing.T) {
	require.Error(t, validateValidatorEmissionPercentage(0.5))
	require.Error(t, validateValidatorEmissionPercentage("-0.50")) // Less than 0 percent should fail
	require.NoError(t, validateValidatorEmissionPercentage("0.50"))
	require.Error(t, validateValidatorEmissionPercentage("1.01")) // More than 100 percent should fail
}

func TestValidateObserverEmissionPercentage(t *testing.T) {
	require.Error(t, validateObserverEmissionPercentage(0.25))
	require.Error(t, validateObserverEmissionPercentage("-0.50")) // Less than 0 percent should fail
	require.NoError(t, validateObserverEmissionPercentage("0.25"))
	require.Error(t, validateObserverEmissionPercentage("1.01")) // More than 100 percent should fail
}

func TestValidateTssEmissionPercentage(t *testing.T) {
	require.Error(t, validateTssEmissonPercentage(0.25))
	require.NoError(t, validateTssEmissonPercentage("0.25"))
	require.Error(t, validateTssEmissonPercentage("1.01")) // More than 100 percent should fail
}

func TestParamsString(t *testing.T) {
	params := DefaultParams()
	out, err := yaml.Marshal(params)
	require.NoError(t, err)
	require.Equal(t, string(out), params.String())
}
