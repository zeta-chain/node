package keeper

import (
	"context"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	MAX_TX_QUERY = 100
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
	from := req.From
	to := req.To
	if req.Last != 0 {
		to = len
		from = to - req.Last
	}
	from, to = validateFromTo(from, to, len)

	return &types.QueryGetTxResponse{Tx: val.Tx[from:to], Length: len}, nil
}

// make sure that [from, to) <= [0, len)
func validateFromTo(from, to, len int64) (int64, int64) {
	//if from < to-MAX_TX_QUERY {
	//	from = to - MAX_TX_QUERY
	//}
	if from < 0 {
		from = 0
	}

	if to > len || to < 0 {
		to = len
	}
	if from > to {
		from = to
	}
	return from, to
}
