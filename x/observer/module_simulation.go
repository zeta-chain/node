package observer

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/zeta-chain/zetacore/testutil/sample"
	observersimulation "github.com/zeta-chain/zetacore/x/observer/simulation"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = observersimulation.FindAccount
	_ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

/* #nosec */
const (
	opWeightMsgUpdateClientParams          = "op_weight_msg_update_client_params"
	defaultWeightMsgUpdateClientParams int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	defaultParams := types.DefaultParams()
	observerGenesis := types.GenesisState{
		Params: &defaultParams,
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&observerGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RandomizedParams creates randomized  param changes for the simulator
func (am AppModule) RandomizedParams(_ *rand.Rand) []simtypes.ParamChange {

	return []simtypes.ParamChange{}
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

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
