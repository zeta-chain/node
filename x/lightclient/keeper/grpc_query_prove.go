package keeper

import (
	"context"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/proofs"
	"github.com/zeta-chain/node/x/lightclient/types"
)

// Prove checks two things:
// 1. the block header is available
// 2. the proof is valid
func (k Keeper) Prove(c context.Context, req *types.QueryProveRequest) (*types.QueryProveResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	// additionalChains is a list of additional chains to search from
	// it is used in the protocol to dynamically support new chains without doing an upgrade
	additionalChains := k.GetAuthorityKeeper().GetAdditionalChainList(ctx)

	blockHash, err := chains.StringToHash(req.ChainId, req.BlockHash, additionalChains)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res, found := k.GetBlockHeader(ctx, blockHash)
	if !found {
		return nil, status.Error(codes.NotFound, "block header not found")
	}

	proven := false

	txBytes, err := req.Proof.Verify(res.Header, int(req.TxIndex))
	if err != nil && !proofs.IsErrorInvalidProof(err) {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if err == nil {
		if chains.IsEVMChain(req.ChainId, additionalChains) {
			var txx ethtypes.Transaction
			err = txx.UnmarshalBinary(txBytes)
			if err != nil {
				return nil, status.Error(codes.Internal, fmt.Sprintf("failed to unmarshal evm transaction: %s", err))
			}
			if txx.Hash().Hex() != req.TxHash {
				return nil, status.Error(
					codes.InvalidArgument,
					fmt.Sprintf("tx hash mismatch: %s != %s", txx.Hash().Hex(), req.TxHash),
				)
			}
			proven = true
		} else if chains.IsBitcoinChain(req.ChainId, additionalChains) {
			tx, err := btcutil.NewTxFromBytes(txBytes)
			if err != nil {
				return nil, status.Error(codes.Internal, fmt.Sprintf("failed to unmarshal btc transaction: %s", err))
			}
			if tx.MsgTx().TxHash().String() != req.TxHash {
				return nil, status.Error(
					codes.InvalidArgument,
					fmt.Sprintf("tx hash mismatch: %s != %s", tx.MsgTx().TxHash().String(), req.TxHash),
				)
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
