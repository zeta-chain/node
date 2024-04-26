package lightclient

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/lightclient/keeper"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

// InitGenesis initializes the lightclient module's state from a provided genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// set block headers
	for _, elem := range genState.BlockHeaders {
		k.SetBlockHeader(ctx, elem)
	}

	// set chain states
	for _, elem := range genState.ChainStates {
		k.SetChainState(ctx, elem)
	}

	// set verification flags for all chains
	for _, elem := range genState.VerificationFlags {
		k.SetVerificationFlags(ctx, elem)
	}
}

// ExportGenesis returns the lightclient module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		BlockHeaders:      k.GetAllBlockHeaders(ctx),
		ChainStates:       k.GetAllChainStates(ctx),
		VerificationFlags: k.GetAllVerificationFlags(ctx),
	}
}
