package lightclient

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/lightclient/keeper"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

// InitGenesis initializes the lightclient module's state from a provided genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {}

// ExportGenesis returns the lightclient module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	var genesis types.GenesisState
	return &genesis
}
