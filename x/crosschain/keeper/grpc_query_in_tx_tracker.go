package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) InTxTrackerAllByChain(goCtx context.Context, request *types.QueryAllInTxTrackerByChainRequest) (*types.QueryAllInTxTrackerByChainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var inTxTrackers []types.InTxTracker
	inTxTrackers, pageRes, err := k.GetAllInTxTrackerForChainPaginated(ctx, request.ChainId, request.Pagination)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryAllInTxTrackerByChainResponse{InTxTracker: inTxTrackers, Pagination: pageRes}, nil
}

func (k Keeper) InTxTrackerAll(goCtx context.Context, _ *types.QueryAllInTxTrackersRequest) (*types.QueryAllInTxTrackersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var inTxTrackers []types.InTxTracker
	inTxTrackers = k.GetAllInTxTracker(ctx)
	return &types.QueryAllInTxTrackersResponse{InTxTracker: inTxTrackers}, nil
}
