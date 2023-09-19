package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SetNodeAccount set a specific nodeAccount in the store from its index
func (k Keeper) SetBlockHeader(ctx sdk.Context, header types.BlockHeader) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlockHeaderKey))
	b := k.cdc.MustMarshal(&header)
	store.Set(header.Hash, b)
}

// GetNodeAccount returns a nodeAccount from its index
func (k Keeper) GetBlockHeader(ctx sdk.Context, hash []byte) (val types.BlockHeader, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlockHeaderKey))

	b := store.Get(hash)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveNodeAccount removes a nodeAccount from the store
func (k Keeper) RemoveBlockHeader(ctx sdk.Context, hash []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlockHeaderKey))
	store.Delete(hash)
}

// GRPC querier for block header
func (k Keeper) GetAllBlockHeaders(c context.Context, req *types.QueryAllBlockHeaderRequest) (*types.QueryAllBlockHeaderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)
	blockHeaderStore := prefix.NewStore(store, types.KeyPrefix(types.BlockHeaderKey))

	var blockHeaders []*types.BlockHeader
	pageRes, err := query.Paginate(blockHeaderStore, req.Pagination, func(key []byte, value []byte) error {
		var blockHeader types.BlockHeader
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

func (k Keeper) GetBlockHeaderByHash(c context.Context, req *types.QueryGetBlockHeaderByHashRequest) (*types.QueryGetBlockHeaderByHashResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if len(req.BlockHash) != 32 {
		return nil, status.Error(codes.InvalidArgument, "invalid hash")
	}

	header, found := k.GetBlockHeader(sdk.UnwrapSDKContext(c), req.BlockHash)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetBlockHeaderByHashResponse{BlockHeader: &header}, nil
}
