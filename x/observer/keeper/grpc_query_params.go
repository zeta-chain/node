package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	return &types.QueryParamsResponse{
		Params: k.GetParams(ctx)}, nil
}

func (k Keeper) GetCoreParamsForChain(goCtx context.Context, req *types.QueryGetCoreParamsForChainRequest) (*types.QueryGetCoreParamsForChainResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	clientParams, found := k.GetCoreParamsByChainID(ctx, req.ChainID)
	if !found {
		return nil, status.Error(codes.NotFound, "client params not found")
	}
	return &types.QueryGetCoreParamsForChainResponse{
		CoreParams: &clientParams,
	}, nil
}
