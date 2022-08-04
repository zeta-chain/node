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

		Erc20List: []types.Erc20{
			{
				Id: 0,
			},
			{
				Id: 1,
			},
		},
		Erc20Count: 2,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.MirrorKeeper(t)
	mirror.InitGenesis(ctx, *k, genesisState)
	got := mirror.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.Erc20List, got.Erc20List)
	require.Equal(t, genesisState.Erc20Count, got.Erc20Count)
	// this line is used by starport scaffolding # genesis/test/assert
}
