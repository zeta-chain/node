package keeper

import (
	"context"
	"fmt"
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

func (k Keeper) ZEVMGetTransaction(c context.Context, req *types.QueryZEVMGetTransactionRequest) (*types.QueryZEVMGetTransactionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	//_ := sdk.UnwrapSDKContext(c)
	rpcclient := types.RPCClient
	if rpcclient == nil {
		return nil, status.Error(codes.Internal, "rpc client is not initialized")
	}
	query := fmt.Sprintf("ethereum_tx.ethereumTxHash='%s'", req.Hash)
	res, err := rpcclient.TxSearch(c, query, false, nil, nil, "asc")
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	result := ""
	hash := ""
	if len(res.Txs) > 0 {
		tx := res.Txs[0]
		result = fmt.Sprintf("%x", tx.TxResult.Log)
		hash = tx.Hash.String()
	}
	return &types.QueryZEVMGetTransactionResponse{
		Value: fmt.Sprintf("%s", result),
		Hash:  hash,
	}, nil
}
