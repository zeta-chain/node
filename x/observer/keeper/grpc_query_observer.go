package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ShowObserverCount(goCtx context.Context, req *types.QueryShowObserverCountRequest) (*types.QueryShowObserverCountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	lb, found := k.GetLastObserverCount(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "last observer count not found")
	}

	return &types.QueryShowObserverCountResponse{
		LastObserverCount: &lb,
	}, nil
}

func (k Keeper) ObserverSet(goCtx context.Context, req *types.QueryObserverSet) (*types.QueryObserverSetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	observerSet, found := k.GetObserverSet(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "observer set not found")
	}
	return &types.QueryObserverSetResponse{
		Observers: observerSet.ObserverList,
	}, nil

}
