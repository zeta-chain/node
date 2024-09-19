package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/node/x/crosschain/types"
)

func (k Keeper) InboundTrackerAllByChain(
	goCtx context.Context,
	request *types.QueryAllInboundTrackerByChainRequest,
) (*types.QueryAllInboundTrackerByChainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var inTxTrackers []types.InboundTracker
	inTxTrackers, pageRes, err := k.GetAllInboundTrackerForChainPaginated(ctx, request.ChainId, request.Pagination)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryAllInboundTrackerByChainResponse{InboundTracker: inTxTrackers, Pagination: pageRes}, nil
}

func (k Keeper) InboundTrackerAll(
	goCtx context.Context,
	req *types.QueryAllInboundTrackersRequest,
) (*types.QueryAllInboundTrackersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var inTxTrackers []types.InboundTracker
	inTxTrackers, pageRes, err := k.GetAllInboundTrackerPaginated(ctx, req.Pagination)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryAllInboundTrackersResponse{InboundTracker: inTxTrackers, Pagination: pageRes}, nil
}

func (k Keeper) InboundTracker(
	goCtx context.Context,
	req *types.QueryInboundTrackerRequest,
) (*types.QueryInboundTrackerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	inTxTracker, found := k.GetInboundTracker(ctx, req.ChainId, req.TxHash)
	if !found {
		return nil, status.Errorf(
			codes.NotFound,
			"Inbound tracker not found for ChainID: %d, TxHash: %s",
			req.ChainId,
			req.TxHash,
		)
	}

	return &types.QueryInboundTrackerResponse{InboundTracker: inTxTracker}, nil
}
