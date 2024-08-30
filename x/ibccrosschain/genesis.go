package ibccrosschain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/ibccrosschain/keeper"
	"github.com/zeta-chain/node/x/ibccrosschain/types"
)

// InitGenesis initializes the ibccrosschain module's state from a provided genesis state
func InitGenesis(_ sdk.Context, _ keeper.Keeper, _ types.GenesisState) {}

// ExportGenesis returns the ibccrosschain module's exported genesis.
func ExportGenesis(_ sdk.Context, _ keeper.Keeper) *types.GenesisState {
	var genesis types.GenesisState
	return &genesis
}
