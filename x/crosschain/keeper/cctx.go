package keeper

import (
	"slices"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/x/crosschain/types"
	observerTypes "github.com/zeta-chain/node/x/observer/types"
)

// SaveCCTXUpdate does the following things in one function:

// 1. Set the Nonce to Cctx mapping
// A new mapping between a nonce and a cctx index should be created only when we add a new outbound to an existing cctx.
// When adding a new outbound , the only two conditions are
// - The cctx is in CctxStatus_PendingOutbound , which means the first outbound has been added, and we need to set the nonce for that
// - The cctx is in CctxStatus_PendingRevert , which means the second outbound has been added, and we need to set the nonce for that

// 2. Set the cctx in the store

// 3. Update the mapping inboundHash -> cctxIndex
// A new value is added to the mapping when a single inbound hash is connected to multiple cctx indexes
// If the inbound hash to cctx mapping does not exist, a new mapping is created and the cctx index is added to the list of cctx indexes

// 4. update the zeta accounting
// Zeta-accounting is updated aborted cctxs of cointtype zeta.When a cctx is aborted it means that `GetAbortedAmount`
//of zeta is locked and cannot be used.

func (k Keeper) SaveCCTXUpdate(
	ctx sdk.Context,
	cctx types.CrossChainTx,
	tssPubkey string,
) {
	k.setNonceToCCTX(ctx, cctx, tssPubkey)
	k.SetCrossChainTx(ctx, cctx)
	k.updateInboundHashToCCTX(ctx, cctx)
	k.updateZetaAccounting(ctx, cctx)
}

// updateInboundHashToCCTX updates the mapping between an inbound hash and a cctx index.
// A new index is added to the list of cctx indexes if it is not already present
func (k Keeper) updateInboundHashToCCTX(
	ctx sdk.Context,
	cctx types.CrossChainTx,
) {
	in, _ := k.GetInboundHashToCctx(ctx, cctx.InboundParams.ObservedHash)
	in.InboundHash = cctx.InboundParams.ObservedHash
	if !slices.Contains(in.CctxIndex, cctx.Index) {
		in.CctxIndex = append(in.CctxIndex, cctx.Index)
	}
	k.SetInboundHashToCctx(ctx, in)
}

// updateZetaAccounting updates the zeta accounting with the amount of zeta that was locked in an aborted cctx
func (k Keeper) updateZetaAccounting(
	ctx sdk.Context,
	cctx types.CrossChainTx,
) {
	if cctx.CctxStatus.Status == types.CctxStatus_Aborted &&
		cctx.InboundParams.CoinType == coin.CoinType_Zeta &&
		!cctx.CctxStatus.IsAbortRefunded {
		k.AddZetaAbortedAmount(ctx, GetAbortedAmount(cctx))
	}
}

// setNonceToCCTX updates the mapping between a nonce and a cctx index if the cctx is in a PendingOutbound or PendingRevert state
func (k Keeper) setNonceToCCTX(
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
}

// SetCrossChainTx set a specific cctx in the store from its index
func (k Keeper) SetCrossChainTx(ctx sdk.Context, cctx types.CrossChainTx) {
	// only set the updated timestamp if the block height is >0 to allow
	// for a genesis import
	if cctx.CctxStatus != nil && ctx.BlockHeight() > 0 {
		cctx.CctxStatus.LastUpdateTimestamp = ctx.BlockHeader().Time.Unix()
	}
	p := types.KeyPrefix(types.CCTXKey)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
	b := k.cdc.MustMarshal(&cctx)
	cctxIndex := types.KeyPrefix(cctx.Index)

	isUpdate := store.Has(cctxIndex)
	store.Set(cctxIndex, b)

	if !isUpdate {
		k.setCctxCounterIndex(ctx, cctx)
	}
}

// GetCrossChainTx returns a cctx from its index
func (k Keeper) GetCrossChainTx(ctx sdk.Context, index string) (val types.CrossChainTx, found bool) {
	p := types.KeyPrefix(types.CCTXKey)
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
	p := types.KeyPrefix(types.CCTXKey)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)

	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

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
	p := types.KeyPrefix(types.CCTXKey)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
	store.Delete(types.KeyPrefix(index))
}
