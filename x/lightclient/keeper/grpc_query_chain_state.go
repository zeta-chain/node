package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetChainStateByChain(c context.Context, req *types.QueryGetChainStateRequest) (*types.QueryGetChainStateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	state, found := k.GetChainState(sdk.UnwrapSDKContext(c), req.ChainId)
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("not found: chain id %d", req.ChainId))
	}

	return &types.QueryGetChainStateResponse{ChainState: &state}, nil
}
