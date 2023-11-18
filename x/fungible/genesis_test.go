package fungible_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/nullify"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		ForeignCoinList: []types.ForeignCoin{
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
