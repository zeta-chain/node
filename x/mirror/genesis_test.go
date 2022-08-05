package mirror_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/mirror"
	"github.com/zeta-chain/zetacore/x/mirror/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		ERC20TokenPairs: &types.ERC20TokenPairs{
			TokenPairs: "6",
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.MirrorKeeper(t)
	mirror.InitGenesis(ctx, *k, genesisState)
	got := mirror.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.ERC20TokenPairs, got.ERC20TokenPairs)
	// this line is used by starport scaffolding # genesis/test/assert
}
