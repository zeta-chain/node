package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/x/crosschain/types"
	observerTypes "github.com/zeta-chain/node/x/observer/types"
)

// SetCctxAndNonceToCctxAndInboundHashToCctx does the following things in one function:
// 1. set the cctx in the store
// 2. set the mapping inboundHash -> cctxIndex , one inboundHash can be connected to multiple cctxindex
// 3. set the mapping nonce => cctx
// 4. update the zeta accounting
func (k Keeper) SetCctxAndNonceToCctxAndInboundHashToCctx(
	ctx sdk.Context,
	cctx types.CrossChainTx,
	tssPubkey string,
) {
	// set mapping nonce => cctxIndex
	if cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound ||
		cctx.CctxStatus.Status == types.CctxStatus_PendingRevert {
		k.GetObserverKeeper().SetNonceToCctx(ctx, observerTypes.NonceToCctx{
			ChainId: cctx.GetCurrentOutboundParam().ReceiverChainId,
			// #nosec G115 always in range
			Nonce:     int64(cctx.GetCurrentOutboundParam().TssNonce),
			CctxIndex: cctx.Index,
			Tss:       tssPubkey,
		})
	}

	k.SetCrossChainTx(ctx, cctx)
	// set mapping inboundHash -> cctxIndex
	in, _ := k.GetInboundHashToCctx(ctx, cctx.InboundParams.ObservedHash)
	in.InboundHash = cctx.InboundParams.ObservedHash
	found := false
	for _, cctxIndex := range in.CctxIndex {
		if cctxIndex == cctx.Index {
			found = true
			break
		}
	}
	if !found {
		in.CctxIndex = append(in.CctxIndex, cctx.Index)
	}
	k.SetInboundHashToCctx(ctx, in)

	if cctx.CctxStatus.Status == types.CctxStatus_Aborted && cctx.InboundParams.CoinType == coin.CoinType_Zeta {
		k.AddZetaAbortedAmount(ctx, GetAbortedAmount(cctx))
	}
}

// SetCrossChainTx set a specific cctx in the store from its index
func (k Keeper) SetCrossChainTx(ctx sdk.Context, cctx types.CrossChainTx) {
	// only set the update timestamp if the block height is >0 to allow
	// for a genesis import
	if cctx.CctxStatus != nil && ctx.BlockHeight() > 0 {
		cctx.CctxStatus.LastUpdateTimestamp = ctx.BlockHeader().Time.Unix()
	}
	p := types.KeyPrefix(fmt.Sprintf("%s", types.CCTXKey))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
	b := k.cdc.MustMarshal(&cctx)
	store.Set(types.KeyPrefix(cctx.Index), b)
}

// GetCrossChainTx returns a cctx from its index
func (k Keeper) GetCrossChainTx(ctx sdk.Context, index string) (val types.CrossChainTx, found bool) {
	p := types.KeyPrefix(fmt.Sprintf("%s", types.CCTXKey))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetAllCrossChainTx returns all cctxs
func (k Keeper) GetAllCrossChainTx(ctx sdk.Context) (list []types.CrossChainTx) {
	p := types.KeyPrefix(fmt.Sprintf("%s", types.CCTXKey))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)

	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.CrossChainTx
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

// RemoveCrossChainTx removes a cctx from the store
func (k Keeper) RemoveCrossChainTx(ctx sdk.Context, index string) {
	p := types.KeyPrefix(fmt.Sprintf("%s", types.CCTXKey))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
	store.Delete(types.KeyPrefix(index))
}
