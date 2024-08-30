package emissions

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/emissions/keeper"
	"github.com/zeta-chain/node/x/emissions/types"
)

// InitGenesis initializes the emissions module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(fmt.Sprintf("invalid emissions module params: %v\n", genState.Params))
	}

	for _, we := range genState.WithdrawableEmissions {
		k.SetWithdrawableEmission(ctx, we)
	}
}

// ExportGenesis returns the emissions module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	var genesis types.GenesisState
	params, found := k.GetParams(ctx)
	if !found {
		params = types.Params{}
	}
	genesis.Params = params
	genesis.WithdrawableEmissions = k.GetAllWithdrawableEmission(ctx)

	return &genesis
}
