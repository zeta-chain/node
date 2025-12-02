package crosschain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// InitGenesis initializes the crosschain module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetZetaAccounting(ctx, genState.ZetaAccounting)

	// Set all the outbound tracker
	for _, elem := range genState.OutboundTrackerList {
		k.SetOutboundTracker(ctx, elem)
	}

	// Set all the inTxTracker
	for _, elem := range genState.InboundTrackerList {
		k.SetInboundTracker(ctx, elem)
	}

	// Set all the inTxHashToCctx
	for _, elem := range genState.InboundHashToCctxList {
		k.SetInboundHashToCctx(ctx, elem)
	}

	// Set all the gasPrice
	for _, elem := range genState.GasPriceList {
		if elem != nil {
			k.SetGasPrice(ctx, *elem)
		}
	}

	// Set all the last block heights
	for _, elem := range genState.LastBlockHeightList {
		if elem != nil {
			k.SetLastBlockHeight(ctx, *elem)
		}
	}

	// Set the cross-chain transactions only,
	// We don't need to call SaveCCTXUpdate as the other fields are being set already
	for _, elem := range genState.CrossChainTxs {
		if elem != nil {
			cctx := *elem
			k.SetCrossChainTx(ctx, cctx)
		}
	}

	for _, elem := range genState.FinalizedInbounds {
		k.SetFinalizedInbound(ctx, elem)
	}

	k.SetRateLimiterFlags(ctx, genState.RateLimiterFlags)

	k.SetCctxCounter(ctx, genState.Counter)
}

// ExportGenesis returns the crosschain module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	var genesis types.GenesisState

	genesis.OutboundTrackerList = k.GetAllOutboundTracker(ctx)
	genesis.InboundHashToCctxList = k.GetAllInboundHashToCctx(ctx)
	genesis.InboundTrackerList = k.GetAllInboundTracker(ctx)

	// Get all gas prices
	gasPriceList := k.GetAllGasPrice(ctx)
	for _, elem := range gasPriceList {
		genesis.GasPriceList = append(genesis.GasPriceList, &elem)
	}

	// Get all last block heights
	lastBlockHeightList := k.GetAllLastBlockHeight(ctx)
	for _, elem := range lastBlockHeightList {
		genesis.LastBlockHeightList = append(genesis.LastBlockHeightList, &elem)
	}

	// Get all send
	sendList := k.GetAllCrossChainTx(ctx)
	for _, elem := range sendList {
		genesis.CrossChainTxs = append(genesis.CrossChainTxs, &elem)
	}

	amount, found := k.GetZetaAccounting(ctx)
	if found {
		genesis.ZetaAccounting = amount
	}
	genesis.FinalizedInbounds = k.GetAllFinalizedInbound(ctx)

	rateLimiterFlags, found := k.GetRateLimiterFlags(ctx)
	if found {
		genesis.RateLimiterFlags = rateLimiterFlags
	}

	genesis.Counter = k.GetCctxCounter(ctx)

	return &genesis
}
