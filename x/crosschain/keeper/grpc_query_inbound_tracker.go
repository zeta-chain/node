package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) InboundTrackerAllByChain(goCtx context.Context, request *types.QueryAllInboundTrackerByChainRequest) (*types.QueryAllInboundTrackerByChainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var inTxTrackers []types.InboundTracker
	inTxTrackers, pageRes, err := k.GetAllInboundTrackerForChainPaginated(ctx, request.ChainId, request.Pagination)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryAllInboundTrackerByChainResponse{InboundTracker: inTxTrackers, Pagination: pageRes}, nil
}

func (k Keeper) InboundTrackerAll(goCtx context.Context, req *types.QueryAllInboundTrackersRequest) (*types.QueryAllInboundTrackersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var inTxTrackers []types.InboundTracker
	inTxTrackers, pageRes, err := k.GetAllInboundTrackerPaginated(ctx, req.Pagination)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryAllInboundTrackersResponse{InboundTracker: inTxTrackers, Pagination: pageRes}, nil
}
