package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/node/x/crosschain/types"
)

// RateLimiterFlags queries the rate limiter flags
func (k Keeper) RateLimiterFlags(
	c context.Context,
	req *types.QueryRateLimiterFlagsRequest,
) (*types.QueryRateLimiterFlagsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	rateLimiterFlags, found := k.GetRateLimiterFlags(ctx)
	if !found {
		return nil, status.Error(codes.Internal, "not found")
	}

	return &types.QueryRateLimiterFlagsResponse{RateLimiterFlags: rateLimiterFlags}, nil
}
