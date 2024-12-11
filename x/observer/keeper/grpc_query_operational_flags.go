package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/node/x/observer/types"
)

func (k Keeper) OperationalFlags(
	c context.Context,
	req *types.QueryOperationalFlagsRequest,
) (*types.QueryOperationalFlagsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	// ignoring found is intentional
	operationalFlags, _ := k.GetOperationalFlags(ctx)
	return &types.QueryOperationalFlagsResponse{
		OperationalFlags: operationalFlags,
	}, nil
}
