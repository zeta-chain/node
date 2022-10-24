package zetaobserver

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetaobserver/keeper"
	"github.com/zeta-chain/zetacore/x/zetaobserver/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	genesisObservers := genState.Observers
	types.VerifyObserverMapper(genesisObservers)
	for _, mapper := range genesisObservers {
		k.SetObserverMapper(ctx, mapper)
	}
	k.SetParams(ctx, types.DefaultParams())
	k.SetSupportedChain(ctx, types.SupportedChains{ChainList: []types.ObserverChain{
		types.ObserverChain_Eth,
		types.ObserverChain_Polygon,
		types.ObserverChain_Bsc,
		types.ObserverChain_Goerli,
		types.ObserverChain_Ropsten,
		types.ObserverChain_Baobab,
	}})
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
