package simulation

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func RandomizedGenState(simState *module.SimulationState) {
	crosschainGenesis := types.DefaultGenesis()
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(crosschainGenesis)
}
