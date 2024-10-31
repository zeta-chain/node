package authority

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

const (
	// #nosec G101 not a hardcoded credential
	opWeightMsgUpdateClientParams          = "op_weight_msg_update_client_params"
	defaultWeightMsgUpdateClientParams int = 100
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	//observers := make([]string, len(simState.Accounts))
	//for _, account := range simState.Accounts {
	//	observers = append(observers, account.Address.String())
	//}
	//observerGenesis := types.DefaultGenesis()
	//observerGenesis.Observers = types.ObserverSet{ObserverList: observers}
	//simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(observerGenesis)

}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

func (AppModule) ProposalMsgs(_ module.SimulationState) []simtypes.WeightedProposalMsg {
	return nil
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgUpdateClientParams int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgUpdateClientParams, &weightMsgUpdateClientParams, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateClientParams = defaultWeightMsgUpdateClientParams
		},
	)

	return operations
}
