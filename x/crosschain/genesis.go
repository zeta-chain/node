package crosschain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// InitGenesis initializes the crosschain module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Params
	k.SetParams(ctx, genState.Params)

	// Set all the outTxTracker
	for _, elem := range genState.OutTxTrackerList {
		k.SetOutTxTracker(ctx, elem)
	}

	// Set all the inTxTracker
	for _, elem := range genState.InTxTrackerList {
		k.SetInTxTracker(ctx, elem)
	}

	// Set all the inTxHashToCctx
	for _, elem := range genState.InTxHashToCctxList {
		k.SetInTxHashToCctx(ctx, elem)
	}

	// Set all the gasPrice
	for _, elem := range genState.GasPriceList {
		if elem != nil {
			k.SetGasPrice(ctx, *elem)
		}
	}

	// Set all the chain nonces
	for _, elem := range genState.ChainNoncesList {
		if elem != nil {
			k.SetChainNonces(ctx, *elem)
		}
	}

	// Set all the last block heights
	for _, elem := range genState.LastBlockHeightList {
		if elem != nil {
			k.SetLastBlockHeight(ctx, *elem)
		}
	}

	// Set all the cross-chain txs
	for _, elem := range genState.CrossChainTxs {
		if elem != nil {
			k.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, *elem)
		}
	}

	if genState.Tss != nil {
		if genState.Tss != nil {
			k.SetTSS(ctx, *genState.Tss)
		}
		for _, chain := range common.DefaultChainsList() {
			k.SetPendingNonces(ctx, types.PendingNonces{
				NonceLow:  0,
				NonceHigh: 0,
				ChainId:   chain.ChainId,
				Tss:       genState.Tss.TssPubkey,
			})
		}
		for _, elem := range genState.TssHistory {
			k.SetTSSHistory(ctx, elem)
		}
	}
}

// ExportGenesis returns the crosschain module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	var genesis types.GenesisState

	genesis.Params = k.GetParams(ctx)
	genesis.OutTxTrackerList = k.GetAllOutTxTracker(ctx)
	genesis.InTxHashToCctxList = k.GetAllInTxHashToCctx(ctx)
	genesis.InTxTrackerList = k.GetAllInTxTracker(ctx)

	// Get tss
	tss, found := k.GetTSS(ctx)
	if found {
		genesis.Tss = &tss
	}

	// Get all gas prices
	gasPriceList := k.GetAllGasPrice(ctx)
	for _, elem := range gasPriceList {
		elem := elem
		genesis.GasPriceList = append(genesis.GasPriceList, &elem)
	}

	// Get all chain nonces
	chainNoncesList := k.GetAllChainNonces(ctx)
	for _, elem := range chainNoncesList {
		elem := elem
		genesis.ChainNoncesList = append(genesis.ChainNoncesList, &elem)
	}

	// Get all last block heights
	lastBlockHeightList := k.GetAllLastBlockHeight(ctx)
	for _, elem := range lastBlockHeightList {
		elem := elem
		genesis.LastBlockHeightList = append(genesis.LastBlockHeightList, &elem)
	}

	// Get all send
	sendList := k.GetAllCrossChainTx(ctx)
	for _, elem := range sendList {
		elem := elem
		genesis.CrossChainTxs = append(genesis.CrossChainTxs, &elem)
	}

	genesis.TssHistory = k.GetAllTSS(ctx)

	return &genesis
}
