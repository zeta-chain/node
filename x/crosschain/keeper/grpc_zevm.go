package keeper

import (
	"context"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ZEVMGetTransactionReceipt(c context.Context, req *types.QueryZEVMGetTransactionReceiptRequest) (*types.QueryZEVMGetTransactionReceiptResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.QueryZEVMGetTransactionReceiptResponse{}, nil
}
