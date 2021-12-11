package keeper

import (
	"context"
	"fmt"
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
	if ctx.BlockHeight() != 175100 {
		return
	}
	fmt.Println("Purge ETH outbound with nonce >=45")
	sends := k.GetAllSend(ctx)
	fmt.Println("Looping through %d sends...", len(sends))
	for _, send := range sends {
		if send.Status == types.SendStatus_Abort && send.Nonce >= 45 && send.SenderChain == common.ETHChain.String() {
			fmt.Printf("removing send (abort) with nonce %d\n", send.Nonce)
			k.RemoveSend(ctx, send.Index)
		}
		if send.Status == types.SendStatus_Finalized && send.Nonce >= 45 && send.ReceiverChain == common.ETHChain.String() {
			fmt.Printf("removing send (finalized) with nonce %d\n", send.Nonce)
			k.RemoveSend(ctx, send.Index)
		}
	}
}
