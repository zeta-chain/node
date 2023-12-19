package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetAllBlockHeaders queries all for block header
func (k Keeper) GetAllBlockHeaders(c context.Context, req *types.QueryAllBlockHeaderRequest) (*types.QueryAllBlockHeaderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)
	blockHeaderStore := prefix.NewStore(store, types.KeyPrefix(types.BlockHeaderKey))

	var blockHeaders []*common.BlockHeader
	pageRes, err := query.Paginate(blockHeaderStore, req.Pagination, func(key []byte, value []byte) error {
		var blockHeader common.BlockHeader
		if err := k.cdc.Unmarshal(value, &blockHeader); err != nil {
			return err
		}

		blockHeaders = append(blockHeaders, &blockHeader)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryAllBlockHeaderResponse{BlockHeaders: blockHeaders, Pagination: pageRes}, nil
}

// GetBlockHeaderByHash queries block header by hash
func (k Keeper) GetBlockHeaderByHash(c context.Context, req *types.QueryGetBlockHeaderByHashRequest) (*types.QueryGetBlockHeaderByHashResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	header, found := k.GetBlockHeader(sdk.UnwrapSDKContext(c), req.BlockHash)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetBlockHeaderByHashResponse{BlockHeader: &header}, nil
}

func (k Keeper) GetBlockHeaderStateByChain(c context.Context, req *types.QueryGetBlockHeaderStateRequest) (*types.QueryGetBlockHeaderStateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	state, found := k.GetBlockHeaderState(sdk.UnwrapSDKContext(c), req.ChainId)
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("not found: chain id %d", req.ChainId))
	}

	return &types.QueryGetBlockHeaderStateResponse{BlockHeaderState: &state}, nil
}
