package crosschain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func CleanupState(ctx sdk.Context, keeper keeper.Keeper) {
	completedCctx := keeper.GetAllCctxByStatuses(ctx, []types.CctxStatus{
		types.CctxStatus_OutboundMined,
		types.CctxStatus_Aborted})
	for _, cctx := range completedCctx {
		keeper.RemoveCrossChainTx(ctx, cctx.Index, cctx.CctxStatus.Status)
	}
}
