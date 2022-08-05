package mirror

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/zeta-chain/zetacore/x/mirror/keeper"
	"github.com/zeta-chain/zetacore/x/mirror/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, accountKeeper authkeeper.AccountKeeper, genState types.GenesisState) {
	// Set if defined
	if genState.ERC20TokenPairs != nil {
		k.SetERC20TokenPairs(ctx, *genState.ERC20TokenPairs)
	}
	// ensure erc20 module account is set on genesis
	if acc := accountKeeper.GetModuleAccount(ctx, types.ModuleName); acc == nil {
		// NOTE: shouldn't occur
		panic("the erc20 module account has not been set")
	}
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// Get all eRC20TokenPairs
	eRC20TokenPairs, found := k.GetERC20TokenPairs(ctx)
	if found {
		genesis.ERC20TokenPairs = &eRC20TokenPairs
	}
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
