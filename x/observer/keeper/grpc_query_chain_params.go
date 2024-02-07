package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetChainParamsForChain(
	goCtx context.Context,
	req *types.QueryGetChainParamsForChainRequest,
) (*types.QueryGetChainParamsForChainResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	chainParams, found := k.GetChainParamsByChainID(ctx, req.ChainId)
	if !found {
		return nil, status.Error(codes.NotFound, "chain params not found")
	}
	return &types.QueryGetChainParamsForChainResponse{
		ChainParams: chainParams,
	}, nil
}

func (k Keeper) GetChainParams(
	goCtx context.Context,
	req *types.QueryGetChainParamsRequest,
) (*types.QueryGetChainParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	chainParams, found := k.GetChainParamsList(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "chain params not found")
	}
	return &types.QueryGetChainParamsResponse{
		ChainParams: &chainParams,
	}, nil
}
