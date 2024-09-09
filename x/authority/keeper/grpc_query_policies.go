package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/node/x/authority/types"
)

// Policies queries policies
func (k Keeper) Policies(
	c context.Context,
	req *types.QueryGetPoliciesRequest,
) (*types.QueryGetPoliciesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	// fetch policies
	policies, found := k.GetPolicies(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "policies not found")
	}

	return &types.QueryGetPoliciesResponse{Policies: policies}, nil
}
