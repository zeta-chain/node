package keeper

import (
	"context"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) InTxAll(c context.Context, req *types.QueryAllInTxRequest) (*types.QueryAllInTxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var inTxs []*types.InTx
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	inTxStore := prefix.NewStore(store, types.KeyPrefix(types.InTxKey))

	pageRes, err := query.Paginate(inTxStore, req.Pagination, func(key []byte, value []byte) error {
		var inTx types.InTx
		if err := k.cdc.UnmarshalBinaryBare(value, &inTx); err != nil {
			return err
		}

		inTxs = append(inTxs, &inTx)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllInTxResponse{InTx: inTxs, Pagination: pageRes}, nil
}

func (k Keeper) InTx(c context.Context, req *types.QueryGetInTxRequest) (*types.QueryGetInTxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetInTx(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetInTxResponse{InTx: &val}, nil
}
