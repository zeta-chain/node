package keeper

import (
	"context"
	"fmt"

	"github.com/zeta-chain/zetacore/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	eth "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
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

	proven := false

	val, err := req.Proof.Verify(res.Header, int(req.TxIndex))
	if err != nil && !common.IsErrorInvalidProof(err) {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err == nil {
		var txx ethtypes.Transaction
		err = txx.UnmarshalBinary(val)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to unmarshal transaction: %s", err))
		}
		proven = true
	}

	return &types.QueryProveResponse{
		Valid: proven,
	}, nil
}
