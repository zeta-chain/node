package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ForeignCoinsAll(c context.Context, req *types.QueryAllForeignCoinsRequest) (*types.QueryAllForeignCoinsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var foreignCoinss []types.ForeignCoins
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	foreignCoinsStore := prefix.NewStore(store, types.KeyPrefix(types.ForeignCoinsKeyPrefix))

	pageRes, err := query.Paginate(foreignCoinsStore, req.Pagination, func(key []byte, value []byte) error {
		var foreignCoins types.ForeignCoins
		if err := k.cdc.Unmarshal(value, &foreignCoins); err != nil {
			return err
		}

		foreignCoinss = append(foreignCoinss, foreignCoins)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllForeignCoinsResponse{ForeignCoins: foreignCoinss, Pagination: pageRes}, nil
}

func (k Keeper) ForeignCoins(c context.Context, req *types.QueryGetForeignCoinsRequest) (*types.QueryGetForeignCoinsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetForeignCoins(
		ctx,
		req.Index,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetForeignCoinsResponse{ForeignCoins: val}, nil
}
