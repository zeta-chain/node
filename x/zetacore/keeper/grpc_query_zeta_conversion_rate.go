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

func (k Keeper) ZetaConversionRateAll(c context.Context, req *types.QueryAllZetaConversionRateRequest) (*types.QueryAllZetaConversionRateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var zetaConversionRates []types.ZetaConversionRate
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	zetaConversionRateStore := prefix.NewStore(store, types.KeyPrefix(types.ZetaConversionRateKeyPrefix))

	pageRes, err := query.Paginate(zetaConversionRateStore, req.Pagination, func(key []byte, value []byte) error {
		var zetaConversionRate types.ZetaConversionRate
		if err := k.cdc.Unmarshal(value, &zetaConversionRate); err != nil {
			return err
		}

		zetaConversionRates = append(zetaConversionRates, zetaConversionRate)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllZetaConversionRateResponse{ZetaConversionRate: zetaConversionRates, Pagination: pageRes}, nil
}

func (k Keeper) ZetaConversionRate(c context.Context, req *types.QueryGetZetaConversionRateRequest) (*types.QueryGetZetaConversionRateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetZetaConversionRate(
		ctx,
		req.Index,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetZetaConversionRateResponse{ZetaConversionRate: val}, nil
}
