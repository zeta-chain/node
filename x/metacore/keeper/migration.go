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
	if ctx.BlockHeight() != 175340 {
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
	fmt.Println("restoring ETH nonce to 45")
	nonce, found := k.GetChainNonces(ctx, common.ETHChain.String())
	if found {
		fmt.Println("found nonce; restoring it...")
		nonce.Nonce = 45
		k.SetChainNonces(ctx, nonce)
	}

	fmt.Println("manually confirming POLYGON send with nonce==22")
	send, found := k.GetSend(ctx, "0xfd197a093a82ff211284b89fe106ec9251e1ed14ca7a0a37393cca95c517c014")
	if found {
		fmt.Println("found send 0xfd19; changing its status to mined...")
		send.Status = types.SendStatus_Mined
		k.SetSend(ctx, send)
	}
}
