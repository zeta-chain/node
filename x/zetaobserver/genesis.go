package zetaobserver

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetaobserver/keeper"
	"github.com/zeta-chain/zetacore/x/zetaobserver/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	for _, mapper := range GetGenesisObserverList() {
		k.SetObserverMapper(ctx, mapper)
	}
	k.SetParams(ctx, types.DefaultParams())
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	params := k.GetParams(ctx)
	return &types.GenesisState{
		Voters:    k.GetAllVoters(ctx),
		Observers: k.GetAllObserverMappers(ctx),
		Params:    &params,
	}
}
