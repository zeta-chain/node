package authority

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/authority/keeper"
	"github.com/zeta-chain/node/x/authority/types"
)

// InitGenesis initializes the authority module's state from a provided genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetPolicies(ctx, genState.Policies)
	k.SetChainInfo(ctx, genState.ChainInfo)
	k.SetAuthorizationList(ctx, genState.AuthorizationList)
}

// ExportGenesis returns the authority module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	var genesis types.GenesisState

	policies, found := k.GetPolicies(ctx)
	if found {
		genesis.Policies = policies
	}
	authorizationList, found := k.GetAuthorizationList(ctx)
	if found {
		genesis.AuthorizationList = authorizationList
	}

	chainInfo, found := k.GetChainInfo(ctx)
	if found {
		genesis.ChainInfo = chainInfo
	}

	return &genesis
}
