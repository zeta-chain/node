package emissions_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/nullify"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/emissions"
	"github.com/zeta-chain/node/x/emissions/types"
)

func TestGenesis(t *testing.T) {
	t.Run("should init and export for valid state", func(t *testing.T) {
		params := types.DefaultParams()

		genesisState := types.GenesisState{
			Params: params,
			WithdrawableEmissions: []types.WithdrawableEmissions{
				sample.WithdrawableEmissions(t),
				sample.WithdrawableEmissions(t),
				sample.WithdrawableEmissions(t),
			},
		}

		// Init and export
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		emissions.InitGenesis(ctx, *k, genesisState)
		got := emissions.ExportGenesis(ctx, *k)
		require.NotNil(t, got)

		// Compare genesis after init and export
		nullify.Fill(&genesisState)
		nullify.Fill(got)
		require.Equal(t, genesisState, *got)
	})

	t.Run("should error for invalid params", func(t *testing.T) {
		params := types.DefaultParams()
		params.ObserverSlashAmount = sdk.NewInt(-1)

		genesisState := types.GenesisState{
			Params: params,
			WithdrawableEmissions: []types.WithdrawableEmissions{
				sample.WithdrawableEmissions(t),
				sample.WithdrawableEmissions(t),
				sample.WithdrawableEmissions(t),
			},
		}

		// Init and export
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		require.Panics(t, func() {
			emissions.InitGenesis(ctx, *k, genesisState)
		})
	})
}
