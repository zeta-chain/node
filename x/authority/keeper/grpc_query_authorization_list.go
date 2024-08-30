package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/node/x/authority/types"
)

// AuthorizationList returns the list of authorizations
func (k Keeper) AuthorizationList(c context.Context,
	req *types.QueryAuthorizationListRequest,
) (*types.QueryAuthorizationListResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	authorizationList, found := k.GetAuthorizationList(ctx)
	if !found {
		return nil, status.Error(codes.Internal, types.ErrAuthorizationListNotFound.Error())
	}

	return &types.QueryAuthorizationListResponse{AuthorizationList: authorizationList}, nil
}

// Authorization returns the authorization for a given message URL
func (k Keeper) Authorization(c context.Context,
	req *types.QueryAuthorizationRequest,
) (*types.QueryAuthorizationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	err := types.ValidateMsgURL(req.MsgUrl)
	if err != nil {
		return nil, err
	}

	authorizationList, found := k.GetAuthorizationList(ctx)
	if !found {
		return nil, status.Error(codes.Internal, types.ErrAuthorizationListNotFound.Error())
	}

	authorization, err := authorizationList.GetAuthorizedPolicy(req.MsgUrl)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAuthorizationResponse{Authorization: types.Authorization{
		MsgUrl:           req.MsgUrl,
		AuthorizedPolicy: authorization,
	}}, nil
}
