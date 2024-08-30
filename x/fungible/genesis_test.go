package fungible_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/nullify"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/fungible"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		ForeignCoinsList: []types.ForeignCoins{
			sample.ForeignCoins(t, sample.EthAddress().String()),
			sample.ForeignCoins(t, sample.EthAddress().String()),
			sample.ForeignCoins(t, sample.EthAddress().String()),
		},
		SystemContract: sample.SystemContract(),
	}

	// Init and export
	k, ctx, _, _ := keepertest.FungibleKeeper(t)
	fungible.InitGenesis(ctx, *k, genesisState)
	got := fungible.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	// Compare genesis after init and export
	nullify.Fill(&genesisState)
	nullify.Fill(got)
	require.Equal(t, genesisState, *got)
}
