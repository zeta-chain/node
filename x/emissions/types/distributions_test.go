package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/x/emissions/types"
)

func TestKeeper_GetRewardsDistributions(t *testing.T) {
	t.Run("Return fractions of block reward", func(t *testing.T) {
		params := types.NewParams()
		params.ValidatorEmissionPercentage = "0.5"
		params.ObserverEmissionPercentage = "0.25"
		params.TssSignerEmissionPercentage = "0.25"
		val, obs, tss := types.GetRewardsDistributions(params)

		require.EqualValues(t, "1687885802469135802", val.String()) // 0.5 * block reward
		require.EqualValues(t, "843942901234567901", obs.String())  // 0.25 * block reward
		require.EqualValues(t, "843942901234567901", tss.String())  // 0.25 * block reward
	})

	t.Run("Return zero in case of invalid string", func(t *testing.T) {
		val, obs, tss := types.GetRewardsDistributions(types.Params{
			ValidatorEmissionPercentage: "invalid",
			ObserverEmissionPercentage:  "invalid",
			TssSignerEmissionPercentage: "invalid",
		})

		require.True(t, val.IsZero())
		require.True(t, obs.IsZero())
		require.True(t, tss.IsZero())
	})
}
