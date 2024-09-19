package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/authority/types"
)

// ChainInfo queries chain info
func (k Keeper) ChainInfo(
	c context.Context,
	req *types.QueryGetChainInfoRequest,
) (*types.QueryGetChainInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	// fetch chain info
	// if the object has not been initialized, return an object containing an empty list
	chainInfo, found := k.GetChainInfo(ctx)
	if !found {
		chainInfo.Chains = []chains.Chain{}
	}

	return &types.QueryGetChainInfoResponse{ChainInfo: chainInfo}, nil
}
