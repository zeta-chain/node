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

const (
	// MaxPendingCctxs is the maximum number of pending cctxs that can be queried
	MaxPendingCctxs = 500
)

func (k Keeper) ZetaAccounting(c context.Context, _ *types.QueryZetaAccountingRequest) (*types.QueryZetaAccountingResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	amount, found := k.GetZetaAccounting(ctx)
	if !found {
		return nil, status.Error(codes.Internal, "aborted zeta amount not found")
	}
	return &types.QueryZetaAccountingResponse{
		AbortedZetaAmount: amount.AbortedZetaAmount.String(),
	}, nil
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
	tss, found := k.zetaObserverKeeper.GetTSS(ctx)
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

// CctxListPending returns a list of pending cctxs and the total number of pending cctxs
// a limit for the number of cctxs to return can be specified
// if no limit is specified, the default is MaxPendingCctxs
func (k Keeper) CctxListPending(c context.Context, req *types.QueryListCctxPendingRequest) (*types.QueryListCctxPendingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// check limit
	// if no limit specified, default to MaxPendingCctxs
	if req.Limit > MaxPendingCctxs {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("limit exceeds max limit of %d", MaxPendingCctxs))
	}
	limit := req.Limit
	if limit == 0 {
		limit = MaxPendingCctxs
	}

	ctx := sdk.UnwrapSDKContext(c)

	// query the nonces that are pending
	tss, found := k.zetaObserverKeeper.GetTSS(ctx)
	if !found {
		return nil, status.Error(codes.Internal, "tss not found")
	}
	pendingNonces, found := k.GetPendingNonces(ctx, tss.TssPubkey, req.ChainId)
	if !found {
		return nil, status.Error(codes.Internal, "pending nonces not found")
	}

	cctxs := make([]*types.CrossChainTx, 0)
	maxCCTXsReached := func() bool {
		// #nosec G701 len always positive
		return uint32(len(cctxs)) >= limit
	}

	totalPending := uint64(0)

	// now query the previous nonces up to 1000 prior to find any pending cctx that we might have missed
	// need this logic because a confirmation of higher nonce will automatically update the p.NonceLow
	// therefore might mask some lower nonce cctx that is still pending.
	startNonce := pendingNonces.NonceLow - 1000
	if startNonce < 0 {
		startNonce = 0
	}
	for i := startNonce; i < pendingNonces.NonceLow; i++ {
		nonceToCctx, found := k.GetNonceToCctx(ctx, tss.TssPubkey, req.ChainId, i)
		if !found {
			return nil, status.Error(codes.Internal, fmt.Sprintf("nonceToCctx not found: nonce %d, chainid %d", i, req.ChainId))
		}
		cctx, found := k.GetCrossChainTx(ctx, nonceToCctx.CctxIndex)
		if !found {
			return nil, status.Error(codes.Internal, fmt.Sprintf("cctx not found: index %s", nonceToCctx.CctxIndex))
		}
		if cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound || cctx.CctxStatus.Status == types.CctxStatus_PendingRevert {
			totalPending++

			// we check here if max cctxs is reached because we want to return the total pending cctxs
			// even if we have reached the limit
			if !maxCCTXsReached() {
				cctxs = append(cctxs, &cctx)
			}
		}
	}

	// add the pending nonces to the total pending
	// #nosec G701 always in range
	totalPending += uint64(pendingNonces.NonceHigh - pendingNonces.NonceLow)

	// now query the pending nonces that we know are pending
	for i := pendingNonces.NonceLow; i < pendingNonces.NonceHigh && !maxCCTXsReached(); i++ {
		nonceToCctx, found := k.GetNonceToCctx(ctx, tss.TssPubkey, req.ChainId, i)
		if !found {
			return nil, status.Error(codes.Internal, "nonceToCctx not found")
		}
		cctx, found := k.GetCrossChainTx(ctx, nonceToCctx.CctxIndex)
		if !found {
			return nil, status.Error(codes.Internal, "cctxIndex not found")
		}
		cctxs = append(cctxs, &cctx)
	}

	return &types.QueryListCctxPendingResponse{
		CrossChainTx: cctxs,
		TotalPending: totalPending,
	}, nil
}
