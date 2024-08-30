package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/node/x/lightclient/types"
)

// ChainStateAll queries all chain statess
func (k Keeper) ChainStateAll(
	c context.Context,
	req *types.QueryAllChainStateRequest,
) (*types.QueryAllChainStateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)
	chainStateStore := prefix.NewStore(store, types.KeyPrefix(types.ChainStateKey))

	var chainStates []types.ChainState
	pageRes, err := query.Paginate(chainStateStore, req.Pagination, func(_ []byte, value []byte) error {
		var chainState types.ChainState
		if err := k.cdc.Unmarshal(value, &chainState); err != nil {
			return err
		}

		chainStates = append(chainStates, chainState)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryAllChainStateResponse{ChainState: chainStates, Pagination: pageRes}, nil
}

// ChainState queries chain state by chain
func (k Keeper) ChainState(
	c context.Context,
	req *types.QueryGetChainStateRequest,
) (*types.QueryGetChainStateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	state, found := k.GetChainState(sdk.UnwrapSDKContext(c), req.ChainId)
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("not found: chain id %d", req.ChainId))
	}

	return &types.QueryGetChainStateResponse{ChainState: &state}, nil
}
