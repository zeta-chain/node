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
	// Set all the outTxTracker
	for _, elem := range genState.OutTxTrackerList {
		k.SetOutTxTracker(ctx, elem)
	}
	// Set all the inTxHashToCctx
	for _, elem := range genState.InTxHashToCctxList {
		k.SetInTxHashToCctx(ctx, elem)
	}

	// Set all the gasPrice
	for _, elem := range genState.GasPriceList {
		k.SetGasPrice(ctx, *elem)
	}

	// Set all the chainNonces
	for _, elem := range genState.ChainNoncesList {
		k.SetChainNonces(ctx, *elem)
	}

	// Set all the lastBlockHeight
	for _, elem := range genState.LastBlockHeightList {
		k.SetLastBlockHeight(ctx, *elem)
	}

	// Set all the send
	for _, elem := range genState.CrossChainTxs {
		k.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, *elem)
	}

	if genState.Tss != nil {
		k.SetTSS(ctx, *genState.Tss)
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
	genesis := types.DefaultGenesis()

	genesis.OutTxTrackerList = k.GetAllOutTxTracker(ctx)
	genesis.InTxHashToCctxList = k.GetAllInTxHashToCctx(ctx)

	// Get all keygen

	// Get all tSSVoter
	// TODO : ADD for single TSS

	// Get all gasPrice
	gasPriceList := k.GetAllGasPrice(ctx)
	for _, elem := range gasPriceList {
		elem := elem
		genesis.GasPriceList = append(genesis.GasPriceList, &elem)
	}

	// Get all chainNonces
	chainNoncesList := k.GetAllChainNonces(ctx)
	for _, elem := range chainNoncesList {
		elem := elem
		genesis.ChainNoncesList = append(genesis.ChainNoncesList, &elem)
	}

	// Get all lastBlockHeight
	lastBlockHeightList := k.GetAllLastBlockHeight(ctx)
	for _, elem := range lastBlockHeightList {
		elem := elem
		genesis.LastBlockHeightList = append(genesis.LastBlockHeightList, &elem)
	}

	// Get all send
	sendList := k.GetAllCrossChainTx(ctx)
	for _, elem := range sendList {
		e := elem
		genesis.CrossChainTxs = append(genesis.CrossChainTxs, &e)
	}

	genesis.TssHistory = k.GetAllTSS(ctx)

	return genesis
}

// TODO : Verify genesis import and export
