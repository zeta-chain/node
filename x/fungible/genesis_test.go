package fungible_test

import (
	"github.com/zeta-chain/zetacore/testutil/sample"
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/fungible"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		ForeignCoinsList: []types.ForeignCoins{
			sample.ForeignCoins(t),
			sample.ForeignCoins(t),
			sample.ForeignCoins(t),
		},
		SystemContract: sample.SystemContract(),
	}

	// Init and export
	k, ctx := keepertest.FungibleKeeper(t)
	fungible.InitGenesis(ctx, *k, genesisState)
	got := fungible.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	// Compare genesis after init and export
	nullify.Fill(&genesisState)
	nullify.Fill(got)
	require.Equal(t, genesisState, *got)
}
