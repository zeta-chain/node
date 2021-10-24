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

func (k Keeper) ReceiveAll(c context.Context, req *types.QueryAllReceiveRequest) (*types.QueryAllReceiveResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var receives []*types.Receive
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	receiveStore := prefix.NewStore(store, types.KeyPrefix(types.ReceiveKey))

	pageRes, err := query.Paginate(receiveStore, req.Pagination, func(key []byte, value []byte) error {
		var receive types.Receive
		if err := k.cdc.UnmarshalBinaryBare(value, &receive); err != nil {
			return err
		}

		receives = append(receives, &receive)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllReceiveResponse{Receive: receives, Pagination: pageRes}, nil
}

func (k Keeper) Receive(c context.Context, req *types.QueryGetReceiveRequest) (*types.QueryGetReceiveResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetReceive(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetReceiveResponse{Receive: &val}, nil
}
