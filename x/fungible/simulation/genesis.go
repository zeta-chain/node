package simulation

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/zeta-chain/node/x/fungible/types"
)

func RandomizedGenState(simState *module.SimulationState) {
	// We do not have any params that we need to randomize for this module
	// The default state is empty which is sufficient for now , this can be modified later weh we add operations that need existing state data to be processed
	fungibleGenesis := types.DefaultGenesis()
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(fungibleGenesis)
}
