package crosschain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func EndBlocker(ctx sdk.Context, keeper keeper.Keeper) {
	if ctx.BlockHeight()%types.BlocksPerHour == 0 {
		completedCctx := keeper.GetAllCctxByStatuses(ctx, []types.CctxStatus{
			types.CctxStatus_OutboundMined,
			types.CctxStatus_Aborted})
		pendingCCTX := keeper.GetAllCctxByStatuses(ctx, []types.CctxStatus{
			types.CctxStatus_PendingOutbound,
			types.CctxStatus_PendingInbound,
			types.CctxStatus_PendingRevert,
		})
		pendingOutTxTrackers := keeper.GetAllOutTxTracker(ctx)
		for _, cctx := range completedCctx {
			keeper.RemoveCrossChainTx(ctx, cctx.Index, cctx.CctxStatus.Status)
		}
		for _, cctx := range pendingCCTX {
			if IsCCTXExpired(ctx, cctx) {
				keeper.RemoveCrossChainTx(ctx, cctx.Index, cctx.CctxStatus.Status)
			}
		}
		for _, outTxTracker := range pendingOutTxTrackers {
			if IsOutTxTrackerExpired(ctx, outTxTracker) {
				keeper.RemoveOutTxTracker(ctx, outTxTracker.ChainId, outTxTracker.Nonce)
			}
		}
	}
}

func IsCCTXExpired(ctx sdk.Context, cctx *types.CrossChainTx) bool {
	if int64(cctx.InboundTxParams.InboundTxFinalizedZetaHeight)+types.BlocksPerDay <= ctx.BlockHeight() {
		return true
	}
	return false
}

func IsOutTxTrackerExpired(ctx sdk.Context, outTxTracker types.OutTxTracker) bool {
	if int64(outTxTracker.CreationHeight)+types.BlocksPerDay <= ctx.BlockHeight() {
		return true
	}
	return false
}
