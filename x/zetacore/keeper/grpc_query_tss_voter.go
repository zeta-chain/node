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

func (k Keeper) TSSVoterAll(c context.Context, req *types.QueryAllTSSVoterRequest) (*types.QueryAllTSSVoterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var tSSVoters []*types.TSSVoter
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	tSSVoterStore := prefix.NewStore(store, types.KeyPrefix(types.TSSVoterKey))

	pageRes, err := query.Paginate(tSSVoterStore, req.Pagination, func(key []byte, value []byte) error {
		var tSSVoter types.TSSVoter
		if err := k.cdc.Unmarshal(value, &tSSVoter); err != nil {
			return err
		}

		tSSVoters = append(tSSVoters, &tSSVoter)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllTSSVoterResponse{TSSVoter: tSSVoters, Pagination: pageRes}, nil
}

func (k Keeper) TSSVoter(c context.Context, req *types.QueryGetTSSVoterRequest) (*types.QueryGetTSSVoterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetTSSVoter(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetTSSVoterResponse{TSSVoter: &val}, nil
}
