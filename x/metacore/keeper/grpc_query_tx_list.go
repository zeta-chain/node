package keeper

import (
	"context"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) TxList(c context.Context, req *types.QueryGetTxListRequest) (*types.QueryGetTxListResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetTxList(ctx)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetTxListResponse{TxList: &val}, nil
}
