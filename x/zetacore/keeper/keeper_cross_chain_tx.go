package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SetCrossChainTx set a specific send in the store from its index
func (k Keeper) SetCrossChainTx(ctx sdk.Context, send types.CrossChainTx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SendKey))
	b := k.cdc.MustMarshal(&send)
	store.Set(types.KeyPrefix(send.Index), b)
}

// GetCrossChainTx returns a send from its index
func (k Keeper) GetCrossChainTx(ctx sdk.Context, index string) (val types.CrossChainTx, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SendKey))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveCrossChainTx removes a send from the store
func (k Keeper) RemoveCrossChainTx(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SendKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllCrossChainTx returns all send
func (k Keeper) GetAllCrossChainTx(ctx sdk.Context) (list []types.CrossChainTx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SendKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.CrossChainTx
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// Queries

func (k Keeper) SendAll(c context.Context, req *types.QueryAllSendRequest) (*types.QueryAllSendResponse, error) {
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

	return &types.QueryAllSendResponse{CrossChainTx: sends, Pagination: pageRes}, nil
}

func (k Keeper) Send(c context.Context, req *types.QueryGetSendRequest) (*types.QueryGetSendResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetCrossChainTx(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetSendResponse{CrossChainTx: &val}, nil
}

func (k Keeper) SendAllPending(c context.Context, req *types.QueryAllSendPendingRequest) (*types.QueryAllSendPendingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var sends []*types.CrossChainTx
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	sendStore := prefix.NewStore(store, types.KeyPrefix(types.SendKey))
	iterator := sdk.KVStorePrefixIterator(sendStore, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.CrossChainTx
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		// if the status of send is pending, which means Finalized/Revert
		if val.CctxStatus.Status == types.CctxStatus_PendingOutbound || val.CctxStatus.Status == types.CctxStatus_PendingRevert {
			sends = append(sends, &val)
		}
	}
	return &types.QueryAllSendPendingResponse{CrossChainTx: sends}, nil
}
