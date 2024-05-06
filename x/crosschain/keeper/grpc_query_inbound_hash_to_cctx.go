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

func (k Keeper) InboundHashToCctxAll(c context.Context, req *types.QueryAllInboundHashToCctxRequest) (*types.QueryAllInboundHashToCctxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var inTxHashToCctxs []types.InboundHashToCctx
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	inTxHashToCctxStore := prefix.NewStore(store, types.KeyPrefix(types.InboundHashToCctxKeyPrefix))

	pageRes, err := query.Paginate(inTxHashToCctxStore, req.Pagination, func(_ []byte, value []byte) error {
		var inboundHashToCctx types.InboundHashToCctx
		if err := k.cdc.Unmarshal(value, &inboundHashToCctx); err != nil {
			return err
		}

		inTxHashToCctxs = append(inTxHashToCctxs, inboundHashToCctx)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllInboundHashToCctxResponse{InboundHashToCctx: inTxHashToCctxs, Pagination: pageRes}, nil
}

func (k Keeper) InboundHashToCctx(c context.Context, req *types.QueryGetInboundHashToCctxRequest) (*types.QueryGetInboundHashToCctxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetInboundHashToCctx(
		ctx,
		req.InboundHash,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetInboundHashToCctxResponse{InboundHashToCctx: val}, nil
}

// InboundHashToCctxData queries the data of all cctxs indexed by a in tx hash
func (k Keeper) InboundHashToCctxData(
	c context.Context,
	req *types.QueryInboundHashToCctxDataRequest,
) (*types.QueryInboundHashToCctxDataResponse, error) {
	inTxHashToCctxRes, err := k.InboundHashToCctx(c, &types.QueryGetInboundHashToCctxRequest{InboundHash: req.InboundHash})
	if err != nil {
		return nil, err
	}

	cctxs := make([]types.CrossChainTx, len(inTxHashToCctxRes.InboundHashToCctx.CctxIndex))
	ctx := sdk.UnwrapSDKContext(c)
	for i, cctxIndex := range inTxHashToCctxRes.InboundHashToCctx.CctxIndex {
		cctx, found := k.GetCrossChainTx(ctx, cctxIndex)
		if !found {
			// This is an internal error because the cctx should always exist from the index
			return nil, status.Errorf(codes.Internal, "cctx indexed %s doesn't exist", cctxIndex)
		}

		cctxs[i] = cctx
	}

	return &types.QueryInboundHashToCctxDataResponse{CrossChainTxs: cctxs}, nil
}
