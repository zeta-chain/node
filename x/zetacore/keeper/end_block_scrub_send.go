package keeper

import (
	"context"
)

// ScrubGasPriceOfStuckOutTx change (increase) the gas price of a scheduled Send which has been stuck.
// Stuck tx is one that is not confirmed in 100 blocks (roughly 10min)
// Scrub stuck Send with current gas price, if current gas price is much higher (roughtly 20% higher)
// Emit Scrub event.
// TODO: This loops through all sends; it should only loop through "pending" sends
func (k Keeper) ScrubGasPriceOfStuckOutTx(goCtx context.Context) {
	//ctx := sdk.UnwrapSDKContext(goCtx)
	//height := ctx.BlockHeight()
	//if height%100 == 0 { // every 100 blocks, roughly 10min
	//	store := ctx.KVStore(k.storeKey)
	//	sendStore := prefix.NewStore(store, types.KeyPrefix(types.SendKey))
	//	iterator := sdk.KVStorePrefixIterator(sendStore, []byte{})
	//	defer iterator.Close()
	//	for ; iterator.Valid(); iterator.Next() {
	//		var send types.CrossChainTx
	//		k.cdc.MustUnmarshal(iterator.Value(), &send)
	//		// if the status of send is pending, which means Finalized/Revert
	//		if send.CctxStatus.Status == types.SendStatus_PendingOutbound || send.Status == types.SendStatus_PendingRevert {
	//			if height-int64(send.FinalizedMetaHeight) > 100 { // stuck send
	//				var chain string
	//				if send.Status == types.SendStatus_PendingOutbound {
	//					chain = send.ReceiverChain
	//				} else if send.Status == types.SendStatus_PendingRevert {
	//					chain = send.SenderChain
	//				}
	//				gasPrice, isFound := k.GetGasPrice(ctx, chain)
	//				if !isFound {
	//					continue
	//				}
	//				mi := gasPrice.MedianIndex
	//				newGasPrice := big.NewInt(0).SetUint64(gasPrice.Prices[mi])
	//				oldGasPrice, ok := big.NewInt(0).SetString(send.GasPrice, 10)
	//				if !ok {
	//					k.Logger(ctx).Error("failed to parse old gas price")
	//					continue
	//				}
	//				// do nothing if new gas price is even lower than old price
	//				if newGasPrice.Cmp(oldGasPrice) < 0 {
	//					continue
	//				}
	//				targetGasPrice := oldGasPrice.Mul(oldGasPrice, big.NewInt(4))
	//				targetGasPrice = targetGasPrice.Div(targetGasPrice, big.NewInt(3)) // targetGasPrice = oldGasPrice * 1.2
	//				// if current new price is not much higher; make it at least 20% higher
	//				// otherwise replacement tx will be rejected by the node
	//				if newGasPrice.Cmp(targetGasPrice) < 0 {
	//					newGasPrice = targetGasPrice
	//				}
	//				send.GasPrice = newGasPrice.String()
	//				k.SetCrossChainTx(ctx, send)
	//				ctx.EventManager().EmitEvent(
	//					sdk.NewEvent(sdk.EventTypeMessage,
	//						sdk.NewAttribute(sdk.AttributeKeyModule, "zetacore"),
	//						sdk.NewAttribute(types.SubTypeKey, types.SendScrubbed),
	//						sdk.NewAttribute(types.SendHash, send.Index),
	//						sdk.NewAttribute("OldGasPrice", fmt.Sprintf("%d", oldGasPrice)),
	//						sdk.NewAttribute("NewGasPrice", fmt.Sprintf("%d", newGasPrice)),
	//						sdk.NewAttribute("Chain", chain),
	//						sdk.NewAttribute("Nonce", fmt.Sprintf("%d", send.Nonce)),
	//					),
	//				)
	//			}
	//		}
	//	}
	//}
}
