package observer

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/observer/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

// InitGenesis initializes the observer module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	if genState.Observers.Len() > 0 {
		k.SetObserverSet(ctx, genState.Observers)
	} else {
		k.SetObserverSet(ctx, types.ObserverSet{})
	}

	if genState.LastObserverCount != nil {
		k.SetLastObserverCount(ctx, genState.LastObserverCount)
	} else {
		k.SetLastObserverCount(ctx, &types.LastObserverCount{LastChangeHeight: 0, Count: genState.Observers.LenUint()})
	}

	// if chain params are defined, set them
	if len(genState.ChainParamsList.ChainParams) > 0 {
		k.SetChainParamsList(ctx, genState.ChainParamsList)
	} else {
		// if no chain params are defined, set localnet chains for test purposes
		btcChainParams := types.GetDefaultBtcRegtestChainParams()
		btcChainParams.IsSupported = true
		goerliChainParams := types.GetDefaultGoerliLocalnetChainParams()
		goerliChainParams.IsSupported = true
		zetaPrivnetChainParams := types.GetDefaultZetaPrivnetChainParams()
		zetaPrivnetChainParams.IsSupported = true
		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				btcChainParams,
				goerliChainParams,
				zetaPrivnetChainParams,
			},
		})
	}

	// Set all the nodeAccount
	for _, elem := range genState.NodeAccountList {
		if elem != nil {
			k.SetNodeAccount(ctx, *elem)
		}
	}

	// Set if defined
	crosschainFlags := types.DefaultCrosschainFlags()
	if genState.CrosschainFlags != nil {
		crosschainFlags.IsOutboundEnabled = genState.CrosschainFlags.IsOutboundEnabled
		crosschainFlags.IsInboundEnabled = genState.CrosschainFlags.IsInboundEnabled
		if genState.CrosschainFlags.GasPriceIncreaseFlags != nil {
			crosschainFlags.GasPriceIncreaseFlags = genState.CrosschainFlags.GasPriceIncreaseFlags
		}
		k.SetCrosschainFlags(ctx, *crosschainFlags)
	} else {
		k.SetCrosschainFlags(ctx, *types.DefaultCrosschainFlags())
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
				ballotListForHeight[ballot.BallotCreationHeight] = append(
					ballotListForHeight[ballot.BallotCreationHeight],
					ballot.BallotIdentifier,
				)
			}
		}
	}

	for height, ballotList := range ballotListForHeight {
		k.SetBallotList(ctx, &types.BallotListForHeight{
			Height:           height,
			BallotsIndexList: ballotList,
		})
	}

	tss := types.TSS{}
	if genState.Tss != nil {
		tss = *genState.Tss
		k.SetTSS(ctx, tss)
	}

	// Set all the pending nonces
	if genState.PendingNonces != nil {
		for _, pendingNonce := range genState.PendingNonces {
			k.SetPendingNonces(ctx, pendingNonce)
		}
	} else {
		for _, chain := range chains.DefaultChainsList() {
			if genState.Tss != nil {
				k.SetPendingNonces(ctx, types.PendingNonces{
					NonceLow:  0,
					NonceHigh: 0,
					ChainId:   chain.ChainId,
					Tss:       tss.TssPubkey,
				})
			}
		}
	}

	for _, elem := range genState.TssHistory {
		k.SetTSSHistory(ctx, elem)
	}

	for _, elem := range genState.TssFundMigrators {
		k.SetFundMigrator(ctx, elem)
	}

	for _, elem := range genState.BlameList {
		k.SetBlame(ctx, elem)
	}

	for _, chainNonce := range genState.ChainNonces {
		k.SetChainNonces(ctx, chainNonce)
	}
	for _, elem := range genState.NonceToCctx {
		k.SetNonceToCctx(ctx, elem)
	}
	k.SetOperationalFlags(ctx, genState.OperationalFlags)
}

// ExportGenesis returns the observer module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	chainParams, found := k.GetChainParamsList(ctx)
	if !found {
		chainParams = types.ChainParamsList{}
	}

	// Get all node accounts
	nodeAccountList := k.GetAllNodeAccount(ctx)
	nodeAccounts := make([]*types.NodeAccount, len(nodeAccountList))
	for i, elem := range nodeAccountList {
		nodeAccounts[i] = &elem
	}

	// Get all crosschain flags
	cf := types.DefaultCrosschainFlags()
	crosschainFlags, found := k.GetCrosschainFlags(ctx)
	if found {
		cf = &crosschainFlags
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

	// Get tss
	tss := &types.TSS{}
	t, found := k.GetTSS(ctx)
	if found {
		tss = &t
	}

	var pendingNonces []types.PendingNonces
	p, err := k.GetAllPendingNonces(ctx)
	if err == nil {
		pendingNonces = p
	}

	os := types.ObserverSet{}
	observers, found := k.GetObserverSet(ctx)
	if found {
		os = observers
	}

	of, _ := k.GetOperationalFlags(ctx)

	return &types.GenesisState{
		Ballots:           k.GetAllBallots(ctx),
		ChainParamsList:   chainParams,
		Observers:         os,
		NodeAccountList:   nodeAccounts,
		CrosschainFlags:   cf,
		Keygen:            kn,
		LastObserverCount: oc,
		Tss:               tss,
		PendingNonces:     pendingNonces,
		TssHistory:        k.GetAllTSS(ctx),
		TssFundMigrators:  k.GetAllTssFundMigrators(ctx),
		BlameList:         k.GetAllBlame(ctx),
		ChainNonces:       k.GetAllChainNonces(ctx),
		NonceToCctx:       k.GetAllNonceToCctx(ctx),
		OperationalFlags:  of,
	}
}
