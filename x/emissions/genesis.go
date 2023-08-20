package emissions

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/emissions/keeper"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

// InitGenesis initializes the emissions module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, ak types.AccountKeeper, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)

	// This create module accounts if not created
	// TODO: find a proper way  to create module accounts?
	ak.GetModuleAccount(ctx, types.ModuleName)
	ak.GetModuleAccount(ctx, types.UndistributedTssRewardsPool)
	ak.GetModuleAccount(ctx, types.UndistributedObserverRewardsPool)
}

// ExportGenesis returns the emissions module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	var genesis *types.GenesisState
	genesis.Params = k.GetParams(ctx)

	return genesis
}
