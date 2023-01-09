package keeper

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"math/big"
)

// ScrubGasPriceOfStuckOutTx change (increase) the gas price of a scheduled Send which has been stuck.
// Stuck tx is one that is not confirmed in 100 blocks (roughly 10min)
// Scrub stuck Send with current gas price, if current gas price is much higher (roughtly 20% higher)
// Emit Scrub event.
// TODO: This loops through all sends; it should only loop through "pending" sends
func (k Keeper) ScrubGasPriceOfStuckOutTx(goCtx context.Context) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	height := ctx.BlockHeight()
	if height%100 == 0 { // every 100 blocks, roughly 10min
		scrubingStatus := []types.CctxStatus{types.CctxStatus_PendingOutbound, types.CctxStatus_PendingRevert}
		for _, status := range scrubingStatus {
			store := ctx.KVStore(k.storeKey)
			p := types.KeyPrefix(fmt.Sprintf("%s-%d", types.SendKey, status))
			k.ScrubUtility(ctx, store, p)
		}
	}
}

func (k Keeper) ScrubUtility(ctx sdk.Context, store sdk.KVStore, p []byte) {
	sendStore := prefix.NewStore(store, p)
	iterator := sdk.KVStorePrefixIterator(sendStore, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var cctx types.CrossChainTx
		k.cdc.MustUnmarshal(iterator.Value(), &cctx)
		// if the status of send is pending, which means Finalized/Revert
		if cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound || cctx.CctxStatus.Status == types.CctxStatus_PendingRevert {
			if ctx.BlockHeight()-int64(cctx.InBoundTxParams.InBoundTxFinalizedZetaHeight) > 100 { // stuck send
				var chainID int64
				if cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound {
					chainID = cctx.OutBoundTxParams.ReceiverChainId
				} else if cctx.CctxStatus.Status == types.CctxStatus_PendingRevert {
					chainID = cctx.InBoundTxParams.SenderChainID
				}
				gasPrice, isFound := k.GetGasPrice(ctx, chainID)
				if !isFound {
					continue
				}
				mi := gasPrice.MedianIndex
				newGasPrice := big.NewInt(0).SetUint64(gasPrice.Prices[mi])
				oldGasPrice, ok := big.NewInt(0).SetString(cctx.OutBoundTxParams.OutBoundTxGasPrice, 10)
				if !ok {
					k.Logger(ctx).Error("failed to parse old gas price")
					continue
				}
				// do nothing if new gas price is even lower than old price
				if newGasPrice.Cmp(oldGasPrice) < 0 {
					continue
				}
				targetGasPrice := oldGasPrice.Mul(oldGasPrice, big.NewInt(4))
				targetGasPrice = targetGasPrice.Div(targetGasPrice, big.NewInt(3)) // targetGasPrice = oldGasPrice * 1.2
				// if current new price is not much higher; make it at least 20% higher
				// otherwise replacement tx will be rejected by the node
				if newGasPrice.Cmp(targetGasPrice) < 0 {
					newGasPrice = targetGasPrice
				}
				cctx.OutBoundTxParams.OutBoundTxGasPrice = newGasPrice.String()
				// No need to migrate as this function does not change the status of Send
				k.SetCrossChainTx(ctx, cctx)
				EmitCCTXScrubbed(ctx, cctx, chainID, oldGasPrice.String(), newGasPrice.String())
			}
		}
	}
}
