package keeper

import (
	"context"
	"sort"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SendAll(c context.Context, req *types.QueryAllSendRequest) (*types.QueryAllSendResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var sends []*types.Send
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	sendStore := prefix.NewStore(store, types.KeyPrefix(types.SendKey))

	pageRes, err := query.Paginate(sendStore, req.Pagination, func(key []byte, value []byte) error {
		var send types.Send
		if err := k.cdc.Unmarshal(value, &send); err != nil {
			return err
		}

		sends = append(sends, &send)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllSendResponse{Send: sends, Pagination: pageRes}, nil
}

func (k Keeper) Send(c context.Context, req *types.QueryGetSendRequest) (*types.QueryGetSendResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetSendAllStatus(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetSendResponse{Send: &val}, nil
}

func (k Keeper) SendAllPending(c context.Context, req *types.QueryAllSendPendingRequest) (*types.QueryAllSendPendingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	sends := k.GetAllSend(ctx, []types.SendStatus{types.SendStatus_PendingOutbound, types.SendStatus_PendingRevert})

	return &types.QueryAllSendPendingResponse{Send: sends}, nil
}

//Deprecated:SendAllLegacy
func (k Keeper) SendAllLegacy(c context.Context, req *types.QueryAllSendLegacyRequest) (*types.QueryAllSendLegacyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var sends []*types.Send
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	sendStore := prefix.NewStore(store, types.KeyPrefix(types.SendKey))

	pageRes, err := query.Paginate(sendStore, req.Pagination, func(key []byte, value []byte) error {
		var send types.Send
		if err := k.cdc.Unmarshal(value, &send); err != nil {
			return err
		}
		sends = append(sends, &send)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllSendLegacyResponse{Send: sends, Pagination: pageRes}, nil
}

//Deprecated:SendLegacy
func (k Keeper) SendLegacy(c context.Context, req *types.QueryGetSendLegacyRequest) (*types.QueryGetSendLegacyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetSendLegacy(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetSendLegacyResponse{Send: &val}, nil
}

//Deprecated:SendAllPendingLegacy
func (k Keeper) SendAllPendingLegacy(c context.Context, req *types.QueryAllSendPendingLegacyRequest) (*types.QueryAllSendPendingLegacyResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	var sends []*types.Send
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	sendStore := prefix.NewStore(store, types.KeyPrefix(types.SendKey))
	iterator := sdk.KVStorePrefixIterator(sendStore, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Send
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		// if the status of send is pending, which means Finalized/Revert
		if val.Status == types.SendStatus_PendingOutbound || val.Status == types.SendStatus_PendingRevert {
			sends = append(sends, &val)
		}
	}
	sort.SliceStable(sends,
		func(i, j int) bool {
			if sends[i].FinalizedMetaHeight == sends[j].FinalizedMetaHeight {
				return sends[i].Nonce < sends[j].Nonce
			}
			return sends[i].FinalizedMetaHeight < sends[j].FinalizedMetaHeight
		})

	return &types.QueryAllSendPendingLegacyResponse{Send: sends}, nil
}
