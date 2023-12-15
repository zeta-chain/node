package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ObserversByChain(goCtx context.Context, req *types.QueryObserversByChainRequest) (*types.QueryObserversByChainResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO move parsing to client
	// https://github.com/zeta-chain/node/issues/867

	chain := k.GetSupportedChainFromChainID(ctx, req.ChainId)
	if chain == nil {
		return &types.QueryObserversByChainResponse{}, types.ErrSupportedChains
	}
	mapper, found := k.GetObserverMapper(ctx, chain)
	if !found {
		return &types.QueryObserversByChainResponse{}, types.ErrObserverNotPresent
	}
	return &types.QueryObserversByChainResponse{Observers: mapper.ObserverList}, nil
}

func (k Keeper) AllObserverMappers(goCtx context.Context, req *types.QueryAllObserverMappersRequest) (*types.QueryAllObserverMappersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	mappers := k.GetAllObserverMappers(ctx)
	return &types.QueryAllObserverMappersResponse{ObserverMappers: mappers}, nil
}

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
