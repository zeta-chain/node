package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/zetacore/x/authority/types"
)

func (k Keeper) AuthorizationList(c context.Context,
	req *types.QueryAuthorizationListRequest,
) (*types.QueryAuthorizationListResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	authorizationList, found := k.GetAuthorizationList(ctx)
	if !found {
		return nil, types.ErrAuthorizationListNotFound
	}
	return &types.QueryAuthorizationListResponse{AuthorizationList: authorizationList}, nil
}

func (k Keeper) Authorization(c context.Context,
	req *types.QueryAuthorizationRequest,
) (*types.QueryAuthorizationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	err := types.ValidateMsgURL(req.MsgUrl)
	if err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(c)
	authorizationList, found := k.GetAuthorizationList(ctx)
	if !found {
		return nil, types.ErrAuthorizationListNotFound
	}
	authorization, err := authorizationList.GetAuthorizedPolicy(req.MsgUrl)
	if err != nil {
		return nil, err
	}
	return &types.QueryAuthorizationResponse{Authorization: types.Authorization{
		MsgUrl:           req.MsgUrl,
		AuthorizedPolicy: authorization,
	}}, nil
}
