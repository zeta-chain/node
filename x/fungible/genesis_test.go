package fungible_test

import (
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
			{
				Index: "0",
			},
			{
				Index: "1",
			},
		},
		ZetaDepositAndCallContract: &types.ZetaDepositAndCallContract{
			Address: "29",
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.FungibleKeeper(t)
	fungible.InitGenesis(ctx, *k, genesisState)
	got := fungible.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.ForeignCoinsList, got.ForeignCoinsList)
	require.Equal(t, genesisState.ZetaDepositAndCallContract, got.ZetaDepositAndCallContract)
	// this line is used by starport scaffolding # genesis/test/assert
}
