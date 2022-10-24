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

	// setup system contract & some auxiliary contracts
	// setup uniswap v2 factory
	uniswapV2Factory, err := k.DeployUniswapV2Factory(ctx)
	if err != nil {
		panic(err)
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("UniswapV2Factory", uniswapV2Factory.String()),
		),
	)
	wzeta, err := k.DeployWZETA(ctx)
	if err != nil {
		panic(err)
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("DeployWZetaContract", wzeta.String()),
		),
	)

	router, err := k.DeployUniswapV2Router02(ctx, uniswapV2Factory, wzeta)
	if err != nil {
		panic(err)
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("DeployUniswapV2Router02", router.String()),
		),
	)

	SystemContractAddress, err := k.DeploySystemContract(ctx, wzeta, uniswapV2Factory, router)
	if err != nil {
		panic(err)
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute("SystemContractAddress", SystemContractAddress.String()),
		),
	)

	// set the system contract
	system, _ := k.GetSystemContract(ctx)
	system.SystemContract = SystemContractAddress.String()
	k.SetSystemContract(ctx, system)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.ForeignCoinsList = k.GetAllForeignCoins(ctx)
	// Get all zetaDepositAndCallContract
	system, found := k.GetSystemContract(ctx)
	if found {
		genesis.SystemContract = &system
	}
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
