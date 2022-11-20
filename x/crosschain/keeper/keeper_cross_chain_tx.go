package keeper

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) CctxChangePrefixStore(ctx sdk.Context, send types.CrossChainTx, oldStatus types.CctxStatus) {
	// Defensive Programming :Remove first set later
	_, found := k.GetCrossChainTx(ctx, send.Index, oldStatus)
	if found {
		k.RemoveCrossChainTx(ctx, send.Index, oldStatus)
	}
	k.SetCrossChainTx(ctx, send)
}

// SetCrossChainTx set a specific send in the store from its index
func (k Keeper) SetCrossChainTx(ctx sdk.Context, send types.CrossChainTx) {
	p := types.KeyPrefix(fmt.Sprintf("%s-%d", types.SendKey, send.CctxStatus.Status))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
	b := k.cdc.MustMarshal(&send)
	store.Set(types.KeyPrefix(send.Index), b)
}

// GetCrossChainTx returns a send from its index
func (k Keeper) GetCrossChainTx(ctx sdk.Context, index string, status types.CctxStatus) (val types.CrossChainTx, found bool) {
	p := types.KeyPrefix(fmt.Sprintf("%s-%d", types.SendKey, status))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveCrossChainTx removes a send from the store
func (k Keeper) RemoveCrossChainTx(ctx sdk.Context, index string, status types.CctxStatus) {
	p := types.KeyPrefix(fmt.Sprintf("%s-%d", types.SendKey, status))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
	store.Delete(types.KeyPrefix(index))
}

func (k Keeper) GetCctxByIndexAndStatuses(ctx sdk.Context, index string, status []types.CctxStatus) (val types.CrossChainTx, found bool) {
	for _, s := range status {
		p := types.KeyPrefix(fmt.Sprintf("%s-%d", types.SendKey, s))
		store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
		send := store.Get(types.KeyPrefix(index))
		if send != nil {
			k.cdc.MustUnmarshal(send, &val)
			return val, true
		}
	}
	return val, false
}

func (k Keeper) GetCctxByIndex(ctx sdk.Context, index string) (val types.CrossChainTx, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SendKey))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true

}

// GetAllCrossChainTx returns all cctx
func (k Keeper) GetAllCctxByStatuses(ctx sdk.Context, status []types.CctxStatus) (list []*types.CrossChainTx) {
	var sends []*types.CrossChainTx
	for _, s := range status {
		p := types.KeyPrefix(fmt.Sprintf("%s-%d", types.SendKey, s))
		store := prefix.NewStore(ctx.KVStore(k.storeKey), p)
		iterator := sdk.KVStorePrefixIterator(store, []byte{})
		defer iterator.Close()
		for ; iterator.Valid(); iterator.Next() {
			var val types.CrossChainTx
			k.cdc.MustUnmarshal(iterator.Value(), &val)
			sends = append(sends, &val)
		}
	}
	return sends
}

// Queries

func (k Keeper) CctxAll(c context.Context, req *types.QueryAllCctxRequest) (*types.QueryAllCctxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	var sends []*types.CrossChainTx
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	sendStore := prefix.NewStore(store, types.KeyPrefix(types.SendKey))

	pageRes, err := query.Paginate(sendStore, req.Pagination, func(key []byte, value []byte) error {
		var send types.CrossChainTx
		if err := k.cdc.Unmarshal(value, &send); err != nil {
			return err
		}
		sends = append(sends, &send)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllCctxResponse{CrossChainTx: sends, Pagination: pageRes}, nil
}

func (k Keeper) Cctx(c context.Context, req *types.QueryGetCctxRequest) (*types.QueryGetCctxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetCctxByIndexAndStatuses(ctx, req.Index, types.AllStatus())
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetCctxResponse{CrossChainTx: &val}, nil
}

func (k Keeper) CctxAllPending(c context.Context, req *types.QueryAllCctxPendingRequest) (*types.QueryAllCctxPendingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	sends := k.GetAllCctxByStatuses(ctx, []types.CctxStatus{types.CctxStatus_PendingOutbound, types.CctxStatus_PendingRevert})
	return &types.QueryAllCctxPendingResponse{CrossChainTx: sends}, nil
}

func (k Keeper) CreateNewCCTX(ctx sdk.Context, msg *types.MsgVoteOnObservedInboundTx, index string, s types.CctxStatus) types.CrossChainTx {
	inboundParams := &types.InBoundTxParams{
		Sender:                          msg.Sender,
		SenderChain:                     msg.SenderChain,
		InBoundTxObservedHash:           msg.InTxHash,
		InBoundTxObservedExternalHeight: msg.InBlockHeight,
		InBoundTxFinalizedZetaHeight:    0,
		InBoundTXBallotIndex:            index,
	}

	outBoundParams := &types.OutBoundTxParams{
		Receiver:                         msg.Receiver,
		ReceiverChain:                    msg.ReceiverChain,
		Broadcaster:                      0,
		OutBoundTxHash:                   "",
		OutBoundTxTSSNonce:               0,
		OutBoundTxGasLimit:               msg.GasLimit,
		OutBoundTxGasPrice:               "",
		OutBoundTXBallotIndex:            "",
		OutBoundTxFinalizedZetaHeight:    0,
		OutBoundTxObservedExternalHeight: 0,
	}
	status := &types.Status{
		Status:              s,
		StatusMessage:       "",
		LastUpdateTimestamp: ctx.BlockHeader().Time.Unix(),
	}
	newCctx := types.CrossChainTx{
		Creator:          msg.Creator,
		Index:            index,
		ZetaBurnt:        sdk.NewUintFromString(msg.ZetaBurnt),
		ZetaMint:         sdk.ZeroUint(),
		ZetaFees:         sdk.ZeroUint(),
		RelayedMessage:   msg.Message,
		CctxStatus:       status,
		InBoundTxParams:  inboundParams,
		OutBoundTxParams: outBoundParams,
	}
	return newCctx
}

// Cctx Pending Queue
// SetCrossChainTx set a specific send in the store from its index
func (k Keeper) SetCctxPendingQueue(ctx sdk.Context, queue types.CctxPendingQueue) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CctxPendingQueueKeyPrefix))
	b := k.cdc.MustMarshal(&queue)
	store.Set(types.KeyPrefix(queue.Index), b)
}

// GetCrossChainTx returns a send from its index
func (k Keeper) GetCctxPendingQueue(ctx sdk.Context, index string) (val types.CctxPendingQueue, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CctxPendingQueueKeyPrefix))
	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveNodeAccount removes a nodeAccount from the store
func (k Keeper) RemoveCctxPendingQueue(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CctxPendingQueueKeyPrefix))
	store.Delete(types.KeyPrefix(index))
}

// GetAllNodeAccount returns all nodeAccount
func (k Keeper) GetAllCctxPendingQueue(ctx sdk.Context) (list []types.CctxPendingQueue) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CctxPendingQueueKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.CctxPendingQueue
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetAllNodeAccount returns all nodeAccount
func (k Keeper) GetAllCctxPendingQueueByChain(ctx sdk.Context, chain string, limit int64) (list []*types.CrossChainTx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CctxPendingQueueKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte(chain))

	defer iterator.Close()

	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		if i == limit {
			break
		}
		var val types.CctxPendingQueue
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		cctx, found := k.GetCctxByIndex(ctx, val.CctxIndex)
		if found {
			list = append(list, &cctx)
		}
		i++
	}
	return
}

// gRPC service function
func (k Keeper) CctxAllPendingQueue(c context.Context, req *types.QueryAllCctxPendingQueueRequest) (*types.QueryAllCctxPendingQueueResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	cctxs := k.GetAllCctxPendingQueueByChain(ctx, req.Chain, req.Limit)
	return &types.QueryAllCctxPendingQueueResponse{
		CrossChainTx: cctxs,
	}, nil
}
