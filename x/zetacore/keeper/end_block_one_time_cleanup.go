package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	"strconv"
)

func (k Keeper) AbortStaleSends(goCtx context.Context) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	height := ctx.BlockHeight()
	sends := k.GetAllSendSinceBlock(ctx, []types.SendStatus{types.SendStatus_PendingOutbound, types.SendStatus_PendingRevert})

	for _, send := range sends {
		// 28800 blocks is ~48hours
		if send.FinalizedMetaHeight+28800 < uint64(height) {
			nonceString := strconv.Itoa(int(send.Nonce))
			var outTxID string
			if send.Status == types.SendStatus_PendingOutbound {
				outTxID = fmt.Sprintf("%s-%s", send.ReceiverChain, nonceString)
			} else if send.Status == types.SendStatus_PendingRevert {
				outTxID = fmt.Sprintf("%s-%s", send.SenderChain, nonceString)
			}
			k.RemoveOutTxTracker(ctx, outTxID)
			oldStatus := send.Status
			send.Status = types.SendStatus_Aborted
			k.SendMigrateStatus(ctx, *send, oldStatus)
			k.SetSend(ctx, *send)
			ctx.Logger().Info(fmt.Sprintf("Aborted send %s, outTxID %s", send.Index, outTxID))
		}
	}
}
