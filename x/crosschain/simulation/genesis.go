package simulation

import (
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/zeta-chain/node/x/crosschain/types"
)

func RandomizedGenState(simState *module.SimulationState) {
	// Randomization is primarily done for params present in the application state
	// We do not need to randomize the genesis state for the crosschain module for now.
	crosschainGenesis := types.DefaultGenesis()
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(crosschainGenesis)
}
