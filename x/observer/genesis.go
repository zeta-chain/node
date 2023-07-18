package observer

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	genesisObservers := genState.Observers
	types.VerifyObserverMapper(genesisObservers)
	for _, mapper := range genesisObservers {
		k.SetObserverMapper(ctx, mapper)
	}
	k.SetCoreParams(ctx, types.GetCoreParams())
	// Set all the nodeAccount
	for _, elem := range genState.NodeAccountList {
		k.SetNodeAccount(ctx, *elem)
	}
	k.SetParams(ctx, types.DefaultParams())
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	params := k.GetParams(ctx)
	// Get all nodeAccount
	nodeAccountList := k.GetAllNodeAccount(ctx)
	nodeAccounts := make([]*types.NodeAccount, len(nodeAccountList))
	for i, elem := range nodeAccountList {
		nodeAccounts[i] = &elem // #nosec G601
	}
	return &types.GenesisState{
		Ballots:         k.GetAllBallots(ctx),
		Observers:       k.GetAllObserverMappers(ctx),
		Params:          &params,
		NodeAccountList: nodeAccounts,
	}
}
