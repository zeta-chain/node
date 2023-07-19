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
	// Set if defined
	if genState.PermissionFlags != nil {
		k.SetPermissionFlags(ctx, *genState.PermissionFlags)
	} else {
		k.SetPermissionFlags(ctx, types.PermissionFlags{IsInboundEnabled: true})
	}
	// Set if defined
	if genState.Keygen != nil {
		k.SetKeygen(ctx, *genState.Keygen)
	}

}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	params := k.GetParams(ctx)
	// Get all nodeAccount
	nodeAccountList := k.GetAllNodeAccount(ctx)
	nodeAccounts := make([]*types.NodeAccount, len(nodeAccountList))
	for i, elem := range nodeAccountList {
		nodeAccounts[i] = &elem // #nosec G601 // false positive
	}
	// Get all permissionFlags
	pf := types.PermissionFlags{IsInboundEnabled: true}
	permissionFlags, found := k.GetPermissionFlags(ctx)
	if found {
		pf = permissionFlags
	}
	kn := &types.Keygen{}
	keygen, found := k.GetKeygen(ctx)
	if found {
		kn = &keygen
	}

	return &types.GenesisState{
		Ballots:         k.GetAllBallots(ctx),
		Observers:       k.GetAllObserverMappers(ctx),
		Params:          &params,
		NodeAccountList: nodeAccounts,
		PermissionFlags: &pf,
		Keygen:          kn,
	}
}
