package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ZetaAccounting(c context.Context, _ *types.QueryZetaAccountingRequest) (*types.QueryZetaAccountingResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	amount, found := k.GetZetaAccounting(ctx)
	if !found {
		return nil, status.Error(codes.Internal, "aborted zeta amount not found")
	}
	return &types.QueryZetaAccountingResponse{AbortedZetaAmount: amount.AbortedZetaAmount}, nil
}

func (k Keeper) CctxAll(c context.Context, req *types.QueryAllCctxRequest) (*types.QueryAllCctxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	var sends []*types.CrossChainTx
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	sendStore := prefix.NewStore(store, types.KeyPrefix(types.SendKey))

	pageRes, err := query.Paginate(sendStore, req.Pagination, func(key []byte, value []byte) error {
		var send types.CrossChainTx
		if err := k.cdc.Unmarshal(value, &send); err != nil {
			return err
		}
		sends = append(sends, &send)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllCctxResponse{CrossChainTx: sends, Pagination: pageRes}, nil
}

func (k Keeper) CctxByStatus(c context.Context, req *types.QueryCctxByStatusRequest) (*types.QueryCctxByStatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	p := types.KeyPrefix(fmt.Sprintf("%s", types.SendKey))
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)

	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()
	cctxList := make([]types.CrossChainTx, 0)
	for ; iterator.Valid(); iterator.Next() {
		var val types.CrossChainTx
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		if val.CctxStatus.Status == req.Status {
			cctxList = append(cctxList, val)
		}
	}

	return &types.QueryCctxByStatusResponse{CrossChainTx: cctxList}, nil
}

func (k Keeper) Cctx(c context.Context, req *types.QueryGetCctxRequest) (*types.QueryGetCctxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetCrossChainTx(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetCctxResponse{CrossChainTx: &val}, nil
}

func (k Keeper) CctxByNonce(c context.Context, req *types.QueryGetCctxByNonceRequest) (*types.QueryGetCctxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	tss, found := k.GetTSS(ctx)
	if !found {
		return nil, status.Error(codes.Internal, "tss not found")
	}
	// #nosec G701 always in range
	res, found := k.GetNonceToCctx(ctx, tss.TssPubkey, req.ChainID, int64(req.Nonce))
	if !found {
		return nil, status.Error(codes.Internal, fmt.Sprintf("nonceToCctx not found: nonce %d, chainid %d", req.Nonce, req.ChainID))
	}
	val, found := k.GetCrossChainTx(ctx, res.CctxIndex)
	if !found {
		return nil, status.Error(codes.Internal, fmt.Sprintf("cctx not found: index %s", res.CctxIndex))
	}

	return &types.QueryGetCctxResponse{CrossChainTx: &val}, nil
}

func (k Keeper) CctxAllPending(c context.Context, req *types.QueryAllCctxPendingRequest) (*types.QueryAllCctxPendingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	tss, found := k.GetTSS(ctx)
	if !found {
		return nil, status.Error(codes.Internal, "tss not found")
	}
	p, found := k.GetPendingNonces(ctx, tss.TssPubkey, req.ChainId)
	if !found {
		return nil, status.Error(codes.Internal, "pending nonces not found")
	}
	sends := make([]*types.CrossChainTx, 0)

	// now query the previous nonces up to 100 prior to find any pending cctx that we might have missed
	// need this logic because a confirmation of higher nonce will automatically update the p.NonceLow
	// therefore might mask some lower nonce cctx that is still pending.
	startNonce := p.NonceLow - 1000
	if startNonce < 0 {
		startNonce = 0
	}
	for i := startNonce; i < p.NonceLow; i++ {
		res, found := k.GetNonceToCctx(ctx, tss.TssPubkey, req.ChainId, i)
		if !found {
			return nil, status.Error(codes.Internal, fmt.Sprintf("nonceToCctx not found: nonce %d, chainid %d", i, req.ChainId))
		}
		send, found := k.GetCrossChainTx(ctx, res.CctxIndex)
		if !found {
			return nil, status.Error(codes.Internal, fmt.Sprintf("cctx not found: index %s", res.CctxIndex))
		}
		if send.CctxStatus.Status == types.CctxStatus_PendingOutbound || send.CctxStatus.Status == types.CctxStatus_PendingRevert {
			sends = append(sends, &send)
		}
	}

	// now query the pending nonces that we know are pending
	for i := p.NonceLow; i < p.NonceHigh; i++ {
		ntc, found := k.GetNonceToCctx(ctx, tss.TssPubkey, req.ChainId, i)
		if !found {
			return nil, status.Error(codes.Internal, "nonceToCctx not found")
		}
		cctx, found := k.GetCrossChainTx(ctx, ntc.CctxIndex)
		if !found {
			return nil, status.Error(codes.Internal, "cctxIndex not found")
		}
		sends = append(sends, &cctx)
	}

	return &types.QueryAllCctxPendingResponse{CrossChainTx: sends}, nil
}
