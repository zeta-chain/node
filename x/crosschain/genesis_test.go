package crosschain_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		ZetaAccounting: sample.ZetaAccounting(t, "sample"),
		OutTxTrackerList: []types.OutTxTracker{
			sample.OutTxTracker(t, "0"),
			sample.OutTxTracker(t, "1"),
			sample.OutTxTracker(t, "2"),
		},
		InTxTrackerList: []types.InTxTracker{
			sample.InTxTracker(t, "0"),
			sample.InTxTracker(t, "1"),
			sample.InTxTracker(t, "2"),
		},
		FinalizedInbounds: []string{
			sample.Hash().String(),
			sample.Hash().String(),
			sample.Hash().String(),
		},
		GasPriceList: []*types.GasPrice{
			sample.GasPrice(t, "0"),
			sample.GasPrice(t, "1"),
			sample.GasPrice(t, "2"),
		},
		CrossChainTxs: []*types.CrossChainTx{
			sample.CrossChainTx(t, "0"),
			sample.CrossChainTx(t, "1"),
			sample.CrossChainTx(t, "2"),
		},
		LastBlockHeightList: []*types.LastBlockHeight{
			sample.LastBlockHeight(t, "0"),
			sample.LastBlockHeight(t, "1"),
			sample.LastBlockHeight(t, "2"),
		},
		InTxHashToCctxList: []types.InTxHashToCctx{
			sample.InTxHashToCctx(t, "0x0"),
			sample.InTxHashToCctx(t, "0x1"),
			sample.InTxHashToCctx(t, "0x2"),
		},
		RateLimiterFlags: sample.RateLimiterFlags(),
	}

	// Init and export
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)
	crosschain.InitGenesis(ctx, *k, genesisState)
	got := crosschain.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	// Compare genesis after init and export
	nullify.Fill(&genesisState)
	nullify.Fill(got)
	require.Equal(t, genesisState, *got)
}
