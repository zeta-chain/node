package emissions_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/emissions"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		WithdrawableEmissions: []types.WithdrawableEmissions{
			sample.WithdrawableEmissions(t),
			sample.WithdrawableEmissions(t),
			sample.WithdrawableEmissions(t),
		},
	}

	// Init and export
	k, ctx := keepertest.EmissionsKeeper(t)
	emissions.InitGenesis(ctx, *k, genesisState)
	got := emissions.ExportGenesis(ctx, *k)
	assert.NotNil(t, got)

	// Compare genesis after init and export
	nullify.Fill(&genesisState)
	nullify.Fill(got)
	assert.Equal(t, genesisState, *got)
}
