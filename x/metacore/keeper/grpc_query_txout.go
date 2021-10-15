package keeper

import (
	"context"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) TxoutAll(c context.Context, req *types.QueryAllTxoutRequest) (*types.QueryAllTxoutResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var txouts []*types.Txout
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	txoutStore := prefix.NewStore(store, types.KeyPrefix(types.TxoutKey))

	pageRes, err := query.Paginate(txoutStore, req.Pagination, func(key []byte, value []byte) error {
		var txout types.Txout
		if err := k.cdc.UnmarshalBinaryBare(value, &txout); err != nil {
			return err
		}

		txouts = append(txouts, &txout)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllTxoutResponse{Txout: txouts, Pagination: pageRes}, nil
}

func (k Keeper) Txout(c context.Context, req *types.QueryGetTxoutRequest) (*types.QueryGetTxoutResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var txout types.Txout
	ctx := sdk.UnwrapSDKContext(c)

	if !k.HasTxout(ctx, req.Id) {
		return nil, sdkerrors.ErrKeyNotFound
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.TxoutKey))
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetTxoutIDBytes(req.Id)), &txout)

	return &types.QueryGetTxoutResponse{Txout: &txout}, nil
}
