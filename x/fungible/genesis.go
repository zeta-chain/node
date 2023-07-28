package fungible

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/fungible/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState, authKeeper types.AccountKeeper) {
	// Set all the foreignCoins
	for _, elem := range genState.ForeignCoinsList {
		k.SetForeignCoins(ctx, elem)
	}
	// Set if defined
	if genState.SystemContract != nil {
		k.SetSystemContract(ctx, *genState.SystemContract)
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
	// ensure fungible module account is set on genesis
	if acc := authKeeper.GetModuleAccount(ctx, types.ModuleName); acc == nil {
		// NOTE: shouldn't occur
		panic("the fungible module account has not been set")
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// TODO move foreign coins to observer
	// https://github.com/zeta-chain/node/issues/863
	//genesis.ForeignCoinsList = k(ctx)

	// Get all zetaDepositAndCallContract
	system, found := k.GetSystemContract(ctx)
	if found {
		genesis.SystemContract = &system
	}
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
