package keeper

import (
	"context"
	"fmt"

	"github.com/btcsuite/btcutil"
	"github.com/zeta-chain/zetacore/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Prove simply checks two things:
// 1. the block header is available
// 2. the proof is good
func (k Keeper) Prove(c context.Context, req *types.QueryProveRequest) (*types.QueryProveResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	blockHash, err := common.StringToHash(req.ChainId, req.BlockHash)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res, found := k.GetBlockHeader(ctx, blockHash)
	if !found {
		return nil, status.Error(codes.NotFound, "block header not found")
	}

	proven := false

	txBytes, err := req.Proof.Verify(res.Header, int(req.TxIndex))
	if err != nil && !common.IsErrorInvalidProof(err) {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err == nil {
		if common.IsEVMChain(req.ChainId) {
			var txx ethtypes.Transaction
			err = txx.UnmarshalBinary(txBytes)
			if err != nil {
				return nil, status.Error(codes.Internal, fmt.Sprintf("failed to unmarshal evm transaction: %s", err))
			}
			if txx.Hash().Hex() != req.TxHash {
				return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("tx hash mismatch: %s != %s", txx.Hash().Hex(), req.TxHash))
			}
			proven = true
		} else if common.IsBitcoinChain(req.ChainId) {
			tx, err := btcutil.NewTxFromBytes(txBytes)
			if err != nil {
				return nil, status.Error(codes.Internal, fmt.Sprintf("failed to unmarshal btc transaction: %s", err))
			}
			if tx.MsgTx().TxHash().String() != req.TxHash {
				return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("tx hash mismatch: %s != %s", tx.MsgTx().TxHash().String(), req.TxHash))
			}
			proven = true
		} else {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid chain id (%d)", req.ChainId))
		}
	}

	return &types.QueryProveResponse{
		Valid: proven,
	}, nil
}
