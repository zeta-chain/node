package types

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestParamKeyTable(t *testing.T) {
	kt := ParamKeyTable()

	ps := Params{}
	for _, psp := range ps.ParamSetPairs() {
		require.PanicsWithValue(t, "duplicate parameter key", func() {
			kt.RegisterType(psp)
		})
	}
}

func TestParamSetPairs(t *testing.T) {
	params := DefaultParams()
	pairs := params.ParamSetPairs()

	require.Equal(t, 3, len(pairs), "The number of param set pairs should match the expected count")

	assertParamSetPair(t, pairs, KeyPrefix(ObserverParamsKey), &params.ObserverParams, validateVotingThresholds)
	assertParamSetPair(t, pairs, KeyPrefix(AdminPolicyParamsKey), &params.AdminPolicy, validateAdminPolicy)
	assertParamSetPair(t, pairs, KeyPrefix(BallotMaturityBlocksParamsKey), &params.BallotMaturityBlocks, validateBallotMaturityBlocks)
}

func assertParamSetPair(t *testing.T, pairs paramtypes.ParamSetPairs, key []byte, expectedValue interface{}, valFunc paramtypes.ValueValidatorFn) {
	for _, pair := range pairs {
		if string(pair.Key) == string(key) {
			require.Equal(t, expectedValue, pair.Value, "Value does not match for key %s", string(key))

			actualValFunc := pair.ValidatorFn
			require.Equal(t, reflect.ValueOf(valFunc).Pointer(), reflect.ValueOf(actualValFunc).Pointer(), "Val func doesnt match for key %s", string(key))
			return
		}
	}

	t.Errorf("Key %s not found in ParamSetPairs", string(key))
}

func TestParamsString(t *testing.T) {
	params := DefaultParams()
	out, err := yaml.Marshal(params)
	require.NoError(t, err)
	require.Equal(t, string(out), params.String())
}

func TestValidateVotingThresholds(t *testing.T) {
	require.Error(t, validateVotingThresholds("invalid"))

	params := DefaultParams()
	require.NoError(t, validateVotingThresholds(params.ObserverParams))

	params.ObserverParams[0].BallotThreshold = sdk.MustNewDecFromStr("1.1")
	require.Error(t, validateVotingThresholds(params.ObserverParams))
}

func TestValidateAdminPolicy(t *testing.T) {
	require.Error(t, validateAdminPolicy("invalid"))
	require.NoError(t, validateAdminPolicy([]*Admin_Policy{}))
}

func TestValidateBallotMaturityBlocks(t *testing.T) {
	require.Error(t, validateBallotMaturityBlocks("invalid"))
	require.NoError(t, validateBallotMaturityBlocks(int64(1)))
}
