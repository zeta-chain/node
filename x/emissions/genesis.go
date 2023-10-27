package emissions

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/emissions/keeper"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

// InitGenesis initializes the emissions module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)

	for _, we := range genState.WithdrawableEmissions {
		k.SetWithdrawableEmission(ctx, we)
	}
}

// ExportGenesis returns the emissions module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	var genesis types.GenesisState
	genesis.Params = k.GetParams(ctx)
	genesis.WithdrawableEmissions = k.GetAllWithdrawableEmission(ctx)

	return &genesis
}
