package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

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
