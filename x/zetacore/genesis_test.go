package zetacore_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/zetacore"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		OutTxTrackerList: []types.OutTxTracker{
			{
				Index: "0",
			},
			{
				Index: "1",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.ZetacoreKeeper(t)
	zetacore.InitGenesis(ctx, *k, genesisState)
	got := zetacore.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.OutTxTrackerList, got.OutTxTrackerList)
	// this line is used by starport scaffolding # genesis/test/assert
}
