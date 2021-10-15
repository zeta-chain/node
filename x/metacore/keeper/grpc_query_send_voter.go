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

func (k Keeper) SendVoterAll(c context.Context, req *types.QueryAllSendVoterRequest) (*types.QueryAllSendVoterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var sendVoters []*types.SendVoter
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	sendVoterStore := prefix.NewStore(store, types.KeyPrefix(types.SendVoterKey))

	pageRes, err := query.Paginate(sendVoterStore, req.Pagination, func(key []byte, value []byte) error {
		var sendVoter types.SendVoter
		if err := k.cdc.UnmarshalBinaryBare(value, &sendVoter); err != nil {
			return err
		}

		sendVoters = append(sendVoters, &sendVoter)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllSendVoterResponse{SendVoter: sendVoters, Pagination: pageRes}, nil
}

func (k Keeper) SendVoter(c context.Context, req *types.QueryGetSendVoterRequest) (*types.QueryGetSendVoterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetSendVoter(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetSendVoterResponse{SendVoter: &val}, nil
}
