package keeper

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ScrubGasPriceOfStuckOutTx change (increase) the gas price of a scheduled Send which has been stuck.
// Stuck tx is one that is not confirmed in 100 blocks (roughly 10min)
// Scrub stuck Send with current gas price, if current gas price is much higher (roughtly 20% higher)
// Emit Scrub event.
func (k Keeper) ScrubGasPriceOfStuckOutTx(goCtx context.Context) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	height := ctx.BlockHeight()
	if height%100 == 0 { // every 100 blocks, roughly 10min
		store := ctx.KVStore(k.storeKey)
		sendStore := prefix.NewStore(store, types.KeyPrefix(types.SendKey))
		iterator := sdk.KVStorePrefixIterator(sendStore, []byte{})
		defer iterator.Close()
		for ; iterator.Valid(); iterator.Next() {
			var send types.Send
			k.cdc.MustUnmarshal(iterator.Value(), &send)
			// if the status of send is pending, which means Finalized/Revert
			if send.Status == types.SendStatus_PendingOutbound || send.Status == types.SendStatus_PendingRevert {
				if height-int64(send.FinalizedMetaHeight) > 100 { // stuck send
					var chain string
					if send.Status == types.SendStatus_PendingOutbound {
						chain = send.ReceiverChain
					} else if send.Status == types.SendStatus_PendingRevert {
						chain = send.SenderChain
					}
					gasPrice, isFound := k.GetGasPrice(ctx, chain)
					if !isFound {
						continue
					}
					mi := gasPrice.MedianIndex
					newGasPrice := gasPrice.Prices[mi]
					oldGasPrice, err := strconv.ParseInt(send.GasPrice, 10, 64)
					if err != nil {
						continue
					}
					if float64(newGasPrice) < float64(oldGasPrice)*1.2 {
						newGasPrice = uint64(float64(oldGasPrice) * 1.2)
					}
					send.GasPrice = fmt.Sprintf("%d", newGasPrice)
					k.SetSend(ctx, send)
					ctx.EventManager().EmitEvent(
						sdk.NewEvent(sdk.EventTypeMessage,
							sdk.NewAttribute(sdk.AttributeKeyModule, "zetacore"),
							sdk.NewAttribute(types.SubTypeKey, types.SendScrubbed),
							sdk.NewAttribute(types.SendHash, send.Index),
							sdk.NewAttribute("OldGasPrice", fmt.Sprintf("%d", oldGasPrice)),
							sdk.NewAttribute("NewGasPrice", fmt.Sprintf("%d", newGasPrice)),
							sdk.NewAttribute("Chain", chain),
							sdk.NewAttribute("Nonce", fmt.Sprintf("%d", send.Nonce)),
						),
					)
				}
			}
		}
	}
}
