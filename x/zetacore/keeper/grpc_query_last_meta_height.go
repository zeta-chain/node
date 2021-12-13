package keeper

import (
	"context"

	"github.com/zeta-chain/zetacore/x/zetacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) LastMetaHeight(goCtx context.Context, req *types.QueryLastMetaHeightRequest) (*types.QueryLastMetaHeightResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryLastMetaHeightResponse{
		uint64(ctx.BlockHeight()),
	}, nil
}
