package keeper

import (
	"context"
	"github.com/Meta-Protocol/metacore/common"
	"github.com/Meta-Protocol/metacore/x/metacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// for mysterious reason, an outbound with nonce 44 on Goerli
// is skipped, preventing all subsquent oubtound.
// this function simply deletes all Send with scheduled nonce >=45
// on Goerli
func (k Keeper) CleanupEthNonce44Mess(goCtx context.Context) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if ctx.BlockHeader().Height != 175035 {
		return
	}
	sends := k.GetAllSend(ctx)
	for _, send := range sends {
		if send.Status == types.SendStatus_Abort && send.Nonce >= 45 && send.SenderChain == common.ETHChain.String() {
			k.RemoveSend(ctx, send.Index)
		}
		if send.Status == types.SendStatus_Finalized && send.Nonce >= 45 && send.ReceiverChain == common.ETHChain.String() {
			k.RemoveSend(ctx, send.Index)
		}
	}
}
