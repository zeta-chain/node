package keeper

import (
	"context"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
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
		if err := k.cdc.UnmarshalBinaryBare(value, &send); err != nil {
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

	val, found := k.GetSend(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetSendResponse{Send: &val}, nil
}

func (k Keeper) SendAllPending(c context.Context, req *types.QueryAllSendRequest) (*types.QueryAllSendResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var sends []*types.Send
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	sendStore := prefix.NewStore(store, types.KeyPrefix(types.SendKey))

	pageRes, err := query.Paginate(sendStore, req.Pagination, func(key []byte, value []byte) error {
		var send types.Send
		if err := k.cdc.UnmarshalBinaryBare(value, &send); err != nil {
			return err
		}

		if send.Status == types.SendStatus_Finalized || send.Status == types.SendStatus_Abort {
			sends = append(sends, &send)
		}
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllSendResponse{Send: sends, Pagination: pageRes}, nil
}
