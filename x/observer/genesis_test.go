package observer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		Params:    &params,
		Tss:       &tss,
		BlameList: sample.BlameRecordsList(t, 10),
		Ballots: []*types.Ballot{
			sample.Ballot(t, "0"),
			sample.Ballot(t, "1"),
			sample.Ballot(t, "2"),
		},
		Observers: sample.ObserverSet(3),
		NodeAccountList: []*types.NodeAccount{
			sample.NodeAccount(),
			sample.NodeAccount(),
			sample.NodeAccount(),
		},
		CrosschainFlags:   types.DefaultCrosschainFlags(),
		Keygen:            sample.Keygen(t),
		ChainParamsList:   sample.ChainParamsList(),
		LastObserverCount: sample.LastObserverCount(10),
		TssFundMigrators:  []types.TssFundMigratorInfo{sample.TssFundsMigrator(1), sample.TssFundsMigrator(2)},
		ChainNonces: []types.ChainNonces{
			sample.ChainNonces(t, "0"),
			sample.ChainNonces(t, "1"),
			sample.ChainNonces(t, "2"),
		},
		PendingNonces: sample.PendingNoncesList(t, "sample", 20),
		NonceToCctx:   sample.NonceToCctxList(t, "sample", 20),
	}

	// Init and export
	k, ctx := keepertest.ObserverKeeper(t)
	observer.InitGenesis(ctx, *k, genesisState)
	got := observer.ExportGenesis(ctx, *k)
	assert.NotNil(t, got)

	// Compare genesis after init and export
	nullify.Fill(&genesisState)
	nullify.Fill(got)
	assert.Equal(t, genesisState, *got)
}
