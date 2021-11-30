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

func (k Keeper) GasBalanceAll(c context.Context, req *types.QueryAllGasBalanceRequest) (*types.QueryAllGasBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var gasBalances []*types.GasBalance
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	gasBalanceStore := prefix.NewStore(store, types.KeyPrefix(types.GasBalanceKey))

	pageRes, err := query.Paginate(gasBalanceStore, req.Pagination, func(key []byte, value []byte) error {
		var gasBalance types.GasBalance
		if err := k.cdc.UnmarshalBinaryBare(value, &gasBalance); err != nil {
			return err
		}

		gasBalances = append(gasBalances, &gasBalance)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllGasBalanceResponse{GasBalance: gasBalances, Pagination: pageRes}, nil
}

func (k Keeper) GasBalance(c context.Context, req *types.QueryGetGasBalanceRequest) (*types.QueryGetGasBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetGasBalance(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetGasBalanceResponse{GasBalance: &val}, nil
}
