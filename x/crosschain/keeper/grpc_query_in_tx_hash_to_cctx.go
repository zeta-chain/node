package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) InTxHashToCctxAll(c context.Context, req *types.QueryAllInTxHashToCctxRequest) (*types.QueryAllInTxHashToCctxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var inTxHashToCctxs []types.InTxHashToCctx
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	inTxHashToCctxStore := prefix.NewStore(store, types.KeyPrefix(types.InTxHashToCctxKeyPrefix))

	pageRes, err := query.Paginate(inTxHashToCctxStore, req.Pagination, func(key []byte, value []byte) error {
		var inTxHashToCctx types.InTxHashToCctx
		if err := k.cdc.Unmarshal(value, &inTxHashToCctx); err != nil {
			return err
		}

		inTxHashToCctxs = append(inTxHashToCctxs, inTxHashToCctx)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllInTxHashToCctxResponse{InTxHashToCctx: inTxHashToCctxs, Pagination: pageRes}, nil
}

func (k Keeper) InTxHashToCctx(c context.Context, req *types.QueryGetInTxHashToCctxRequest) (*types.QueryGetInTxHashToCctxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetInTxHashToCctx(
		ctx,
		req.InTxHash,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetInTxHashToCctxResponse{InTxHashToCctx: val}, nil
}

// InTxHashToCctxData queries the data of all cctxs indexed by a in tx hash
func (k Keeper) InTxHashToCctxData(
	c context.Context,
	req *types.QueryInTxHashToCctxDataRequest,
) (*types.QueryInTxHashToCctxDataResponse, error) {
	inTxHashToCctxRes, err := k.InTxHashToCctx(c, &types.QueryGetInTxHashToCctxRequest{InTxHash: req.InTxHash})
	if err != nil {
		return nil, err
	}

	var cctxs []types.CrossChainTx
	ctx := sdk.UnwrapSDKContext(c)
	for _, cctxIndex := range inTxHashToCctxRes.InTxHashToCctx.CctxIndex {
		cctx, found := k.GetCrossChainTx(ctx, cctxIndex)
		if !found {
			// This is an internal error because the cctx should always exist from the index
			return nil, status.Errorf(codes.Internal, "cctx indexed %s doesn't exist", cctxIndex)
		}

		cctxs = append(cctxs, cctx)
	}

	return &types.QueryInTxHashToCctxDataResponse{CrossChainTxs: cctxs}, nil
}
