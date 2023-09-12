package observer

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// InitGenesis initializes the observer module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	genesisObservers := genState.Observers
	observerCount := uint64(0)
	for _, mapper := range genesisObservers {
		if mapper != nil {
			k.SetObserverMapper(ctx, mapper)
			observerCount += uint64(len(mapper.ObserverList))
		}
	}

	// If core params are defined set them, otherwise set default
	if len(genState.CoreParamsList.CoreParams) > 0 {
		k.SetCoreParams(ctx, genState.CoreParamsList)
	} else {
		k.SetCoreParams(ctx, types.GetCoreParams())
	}

	// Set all the nodeAccount
	for _, elem := range genState.NodeAccountList {
		if elem != nil {
			k.SetNodeAccount(ctx, *elem)
		}
	}

	params := types.DefaultParams()
	if genState.Params != nil {
		params = *genState.Params
	}
	k.SetParams(ctx, params)

	// Set if defined
	if genState.CrosschainFlags != nil {
		k.SetCrosschainFlags(ctx, *genState.CrosschainFlags)
	} else {
		k.SetCrosschainFlags(ctx, types.CrosschainFlags{IsInboundEnabled: true, IsOutboundEnabled: true})
	}

	// Set if defined
	if genState.Keygen != nil {
		k.SetKeygen(ctx, *genState.Keygen)
	}

	ballotListForHeight := make(map[int64][]string)
	if len(genState.Ballots) > 0 {
		for _, ballot := range genState.Ballots {
			if ballot != nil {
				k.SetBallot(ctx, ballot)
				ballotListForHeight[ballot.BallotCreationHeight] = append(ballotListForHeight[ballot.BallotCreationHeight], ballot.BallotIdentifier)
			}
		}
	}

	for height, ballotList := range ballotListForHeight {
		k.SetBallotList(ctx, &types.BallotListForHeight{
			Height:           height,
			BallotsIndexList: ballotList,
		})
	}

	if genState.LastObserverCount != nil {
		k.SetLastObserverCount(ctx, genState.LastObserverCount)
	} else {
		k.SetLastObserverCount(ctx, &types.LastObserverCount{LastChangeHeight: 0, Count: observerCount})
	}
}

// ExportGenesis returns the observer module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	params := k.GetParams(ctx)

	coreParams, found := k.GetAllCoreParams(ctx)
	if !found {
		coreParams = types.CoreParamsList{}
	}

	// Get all node accounts
	nodeAccountList := k.GetAllNodeAccount(ctx)
	nodeAccounts := make([]*types.NodeAccount, len(nodeAccountList))
	for i, elem := range nodeAccountList {
		elem := elem
		nodeAccounts[i] = &elem
	}

	// Get all permissionFlags
	pf := types.CrosschainFlags{IsInboundEnabled: true, IsOutboundEnabled: true}
	permissionFlags, found := k.GetCrosschainFlags(ctx)
	if found {
		pf = permissionFlags
	}

	kn := &types.Keygen{}
	keygen, found := k.GetKeygen(ctx)
	if found {
		kn = &keygen
	}

	oc := &types.LastObserverCount{}
	observerCount, found := k.GetLastObserverCount(ctx)
	if found {
		oc = &observerCount
	}

	return &types.GenesisState{
		Ballots:           k.GetAllBallots(ctx),
		Observers:         k.GetAllObserverMappers(ctx),
		CoreParamsList:    coreParams,
		Params:            &params,
		NodeAccountList:   nodeAccounts,
		CrosschainFlags:   &pf,
		Keygen:            kn,
		LastObserverCount: oc,
	}
}
