package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) NodeAccountAll(c context.Context, req *types.QueryAllNodeAccountRequest) (*types.QueryAllNodeAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var nodeAccounts []*types.NodeAccount
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	nodeAccountStore := prefix.NewStore(store, types.KeyPrefix(types.NodeAccountKey))

	pageRes, err := query.Paginate(nodeAccountStore, req.Pagination, func(key []byte, value []byte) error {
		var nodeAccount types.NodeAccount
		if err := k.cdc.Unmarshal(value, &nodeAccount); err != nil {
			return err
		}

		nodeAccounts = append(nodeAccounts, &nodeAccount)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllNodeAccountResponse{NodeAccount: nodeAccounts, Pagination: pageRes}, nil
}

func (k Keeper) NodeAccount(c context.Context, req *types.QueryGetNodeAccountRequest) (*types.QueryGetNodeAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetNodeAccount(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetNodeAccountResponse{NodeAccount: &val}, nil
}
