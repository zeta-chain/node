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

func (k Keeper) TSSAll(c context.Context, req *types.QueryAllTSSRequest) (*types.QueryAllTSSResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var tSSs []*types.TSS
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	tSSStore := prefix.NewStore(store, types.KeyPrefix(types.TSSKey))

	pageRes, err := query.Paginate(tSSStore, req.Pagination, func(key []byte, value []byte) error {
		var tSS types.TSS
		if err := k.cdc.UnmarshalBinaryBare(value, &tSS); err != nil {
			return err
		}

		tSSs = append(tSSs, &tSS)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllTSSResponse{TSS: tSSs, Pagination: pageRes}, nil
}

func (k Keeper) TSS(c context.Context, req *types.QueryGetTSSRequest) (*types.QueryGetTSSResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetTSS(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetTSSResponse{TSS: &val}, nil
}
