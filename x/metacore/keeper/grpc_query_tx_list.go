package keeper

import (
	"context"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) TxList(c context.Context, req *types.QueryGetTxRequest) (*types.QueryGetTxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetTxList(ctx)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	len := int64(len(val.Tx))
	to := req.To
	from := req.From

	if to > len || to == -1 {
		to = len
	}
	if from > to {
		from = to
	}
	if len == 0 {
		return &types.QueryGetTxResponse{Tx: []*types.Tx{}}, nil
	}
	return &types.QueryGetTxResponse{Tx: val.Tx[from:to]}, nil
}
