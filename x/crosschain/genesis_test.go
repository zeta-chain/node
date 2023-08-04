package crosschain_test

import (
	"github.com/zeta-chain/zetacore/x/crosschain"
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		OutTxTrackerList: []types.OutTxTracker{
			{
				Index: "0",
			},
			{
				Index: "1",
			},
		},
		InTxHashToCctxList: []types.InTxHashToCctx{
			{
				InTxHash: "0",
			},
			{
				InTxHash: "1",
			},
		},
		//PermissionFlags: &types.PermissionFlags{
		//	IsInboundEnabled: true,
		//},
	}

	k, ctx := keepertest.ZetacoreKeeper(t)
	crosschain.InitGenesis(ctx, *k, genesisState)
	got := crosschain.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.OutTxTrackerList, got.OutTxTrackerList)
	require.ElementsMatch(t, genesisState.InTxHashToCctxList, got.InTxHashToCctxList)
	//require.Equal(t, genesisState.PermissionFlags, got.PermissionFlags)
}
