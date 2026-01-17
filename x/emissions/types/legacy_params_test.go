package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/zeta-chain/node/x/emissions/types"
)

func TestLegacyString(t *testing.T) {
	t.Run("should return correct string representation of legacy params", func(t *testing.T) {
		params := types.LegacyParams{
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
		out, err := yaml.Marshal(params)
		require.NoError(t, err)
		require.Equal(t, string(out), params.String())
	})
}
