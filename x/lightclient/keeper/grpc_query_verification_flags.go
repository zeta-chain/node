package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// VerificationFlags implements the Query/VerificationFlags gRPC method
func (k Keeper) VerificationFlags(c context.Context, req *types.QueryVerificationFlagsRequest) (*types.QueryVerificationFlagsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetVerificationFlags(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryVerificationFlagsResponse{VerificationFlags: val}, nil
}
