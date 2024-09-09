package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/node/pkg/proofs"
	"github.com/zeta-chain/node/x/lightclient/types"
)

// BlockHeaderAll queries all block headers
func (k Keeper) BlockHeaderAll(
	c context.Context,
	req *types.QueryAllBlockHeaderRequest,
) (*types.QueryAllBlockHeaderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	store := ctx.KVStore(k.storeKey)
	blockHeaderStore := prefix.NewStore(store, types.KeyPrefix(types.BlockHeaderKey))

	var blockHeaders []proofs.BlockHeader
	pageRes, err := query.Paginate(blockHeaderStore, req.Pagination, func(_ []byte, value []byte) error {
		var blockHeader proofs.BlockHeader
		if err := k.cdc.Unmarshal(value, &blockHeader); err != nil {
			return err
		}

		blockHeaders = append(blockHeaders, blockHeader)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryAllBlockHeaderResponse{BlockHeaders: blockHeaders, Pagination: pageRes}, nil
}

// BlockHeader queries block header by hash
func (k Keeper) BlockHeader(
	c context.Context,
	req *types.QueryGetBlockHeaderRequest,
) (*types.QueryGetBlockHeaderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	header, found := k.GetBlockHeader(sdk.UnwrapSDKContext(c), req.BlockHash)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetBlockHeaderResponse{BlockHeader: &header}, nil
}
