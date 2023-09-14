package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	eth "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Prove(c context.Context, req *types.QueryProveRequest) (*types.QueryProveResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	blockHash := eth.HexToHash(req.BlockHash)
	res, found := k.GetBlockHeader(ctx, blockHash.Bytes())
	if !found {
		return nil, status.Error(codes.NotFound, "block header not found")
	}
	var header ethtypes.Header
	err := rlp.DecodeBytes(res.Header, &header)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to decode header: %s", err))
	}
	proven := false
	if found {
		val, err := req.Proof.Verify(header.TxHash, int(req.TxIndex))
		if err == nil {
			var txx ethtypes.Transaction
			err = txx.UnmarshalBinary(val)
			if err != nil {
				return nil, status.Error(codes.Internal, fmt.Sprintf("failed to unmarshal transaction: %s", err))
			}
			proven = true
		}
	}
	return &types.QueryProveResponse{
		Valid: proven,
	}, nil
}
