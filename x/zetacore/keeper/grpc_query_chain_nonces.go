package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ChainNoncesAll(c context.Context, req *types.QueryAllChainNoncesRequest) (*types.QueryAllChainNoncesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var chainNoncess []*types.ChainNonces
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	chainNoncesStore := prefix.NewStore(store, types.KeyPrefix(types.ChainNoncesKey))

	pageRes, err := query.Paginate(chainNoncesStore, req.Pagination, func(key []byte, value []byte) error {
		var chainNonces types.ChainNonces
		if err := k.cdc.UnmarshalBinaryBare(value, &chainNonces); err != nil {
			return err
		}

		chainNoncess = append(chainNoncess, &chainNonces)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllChainNoncesResponse{ChainNonces: chainNoncess, Pagination: pageRes}, nil
}

func (k Keeper) ChainNonces(c context.Context, req *types.QueryGetChainNoncesRequest) (*types.QueryGetChainNoncesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetChainNonces(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetChainNoncesResponse{ChainNonces: &val}, nil
}
