package observer

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	genesisObservers := genState.Observers
	types.VerifyObserverMapper(genesisObservers)
	for _, mapper := range genesisObservers {
		k.SetObserverMapper(ctx, mapper)
	}
	allCoreParams := types.GetCoreParams()
	for _, chain := range common.DefaultChainsList() {
		if _, ok := allCoreParams[chain.ChainId]; !ok {
			panic("Core params not found for chain id : " + chain.ChainName.String())
		}
		k.SetCoreParamsByChainID(ctx, chain.ChainId, allCoreParams[chain.ChainId])
	}
	k.SetParams(ctx, types.DefaultParams())
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	params := k.GetParams(ctx)
	return &types.GenesisState{
		Ballots:   k.GetAllBallots(ctx),
		Observers: k.GetAllObserverMappers(ctx),
		Params:    &params,
	}
}
