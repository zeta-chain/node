package simulation

import (
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/zeta-chain/node/x/observer/types"
)

func RandomizedGenState(simState *module.SimulationState) {
	// We do not have any params that we need to randomize for this module
	observerGenesis := types.DefaultGenesis()
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(observerGenesis)
}
