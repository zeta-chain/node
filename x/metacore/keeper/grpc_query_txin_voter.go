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

func (k Keeper) TxinVoterAll(c context.Context, req *types.QueryAllTxinVoterRequest) (*types.QueryAllTxinVoterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var txinVoters []*types.TxinVoter
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	txinVoterStore := prefix.NewStore(store, types.KeyPrefix(types.TxinVoterKey))

	pageRes, err := query.Paginate(txinVoterStore, req.Pagination, func(key []byte, value []byte) error {
		var txinVoter types.TxinVoter
		if err := k.cdc.UnmarshalBinaryBare(value, &txinVoter); err != nil {
			return err
		}

		txinVoters = append(txinVoters, &txinVoter)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllTxinVoterResponse{TxinVoter: txinVoters, Pagination: pageRes}, nil
}

func (k Keeper) TxinVoter(c context.Context, req *types.QueryGetTxinVoterRequest) (*types.QueryGetTxinVoterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetTxinVoter(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetTxinVoterResponse{TxinVoter: &val}, nil
}
