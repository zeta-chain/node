package observer_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestGenesis(t *testing.T) {
	params := types.DefaultParams()
	tss := sample.Tss()
	genesisState := types.GenesisState{
		Params: &params,
		Tss:    &tss,
		Ballots: []*types.Ballot{
			sample.Ballot(t, "0"),
			sample.Ballot(t, "1"),
			sample.Ballot(t, "2"),
		},
		Observers: []*types.ObserverMapper{
			sample.ObserverMapper(t, "0"),
			sample.ObserverMapper(t, "1"),
			sample.ObserverMapper(t, "2"),
		},
		NodeAccountList: []*types.NodeAccount{
			sample.NodeAccount(),
			sample.NodeAccount(),
			sample.NodeAccount(),
		},
		CrosschainFlags:   types.DefaultCrosschainFlags(),
		Keygen:            sample.Keygen(t),
		LastObserverCount: sample.LastObserverCount(1000),
		CoreParamsList:    sample.CoreParamsList(),
	}

	// Init and export
	k, ctx := keepertest.ObserverKeeper(t)
	observer.InitGenesis(ctx, *k, genesisState)
	got := observer.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	// Compare genesis after init and export
	nullify.Fill(&genesisState)
	nullify.Fill(got)
	require.Equal(t, genesisState, *got)
}
