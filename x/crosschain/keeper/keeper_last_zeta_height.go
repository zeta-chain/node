package keeper

import (
	"context"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) LastZetaHeight(goCtx context.Context, req *types.QueryLastZetaHeightRequest) (*types.QueryLastZetaHeightResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	height := ctx.BlockHeight()
	if height >= math.MaxInt64 || height < 0 {
		return nil, status.Error(codes.OutOfRange, "height out of range")
	}
	return &types.QueryLastZetaHeightResponse{
		Height: height,
	}, nil
}
