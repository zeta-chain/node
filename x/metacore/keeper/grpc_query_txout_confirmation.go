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

func (k Keeper) TxoutConfirmationAll(c context.Context, req *types.QueryAllTxoutConfirmationRequest) (*types.QueryAllTxoutConfirmationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var txoutConfirmations []*types.TxoutConfirmation
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	txoutConfirmationStore := prefix.NewStore(store, types.KeyPrefix(types.TxoutConfirmationKey))

	pageRes, err := query.Paginate(txoutConfirmationStore, req.Pagination, func(key []byte, value []byte) error {
		var txoutConfirmation types.TxoutConfirmation
		if err := k.cdc.UnmarshalBinaryBare(value, &txoutConfirmation); err != nil {
			return err
		}

		txoutConfirmations = append(txoutConfirmations, &txoutConfirmation)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllTxoutConfirmationResponse{TxoutConfirmation: txoutConfirmations, Pagination: pageRes}, nil
}

func (k Keeper) TxoutConfirmation(c context.Context, req *types.QueryGetTxoutConfirmationRequest) (*types.QueryGetTxoutConfirmationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetTxoutConfirmation(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetTxoutConfirmationResponse{TxoutConfirmation: &val}, nil
}
