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

func (k Keeper) TxinAll(c context.Context, req *types.QueryAllTxinRequest) (*types.QueryAllTxinResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var txins []*types.Txin
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	txinStore := prefix.NewStore(store, types.KeyPrefix(types.TxinKey))

	pageRes, err := query.Paginate(txinStore, req.Pagination, func(key []byte, value []byte) error {
		var txin types.Txin
		if err := k.cdc.UnmarshalBinaryBare(value, &txin); err != nil {
			return err
		}

		txins = append(txins, &txin)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllTxinResponse{Txin: txins, Pagination: pageRes}, nil
}

func (k Keeper) Txin(c context.Context, req *types.QueryGetTxinRequest) (*types.QueryGetTxinResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetTxin(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetTxinResponse{Txin: &val}, nil
}
