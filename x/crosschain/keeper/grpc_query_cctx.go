package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

const (
	// MaxPendingCctxs is the maximum number of pending cctxs that can be queried
	MaxPendingCctxs = 500

	// MaxLookbackNonce is the maximum number of nonces to look back to find missed pending cctxs
	MaxLookbackNonce = 1000
)

func (k Keeper) ZetaAccounting(
	c context.Context,
	_ *types.QueryZetaAccountingRequest,
) (*types.QueryZetaAccountingResponse, error) {
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
	sendStore := prefix.NewStore(store, types.KeyPrefix(types.CCTXKey))

	pageRes, err := query.Paginate(sendStore, req.Pagination, func(_ []byte, value []byte) error {
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

func (k Keeper) CctxByNonce(
	c context.Context,
	req *types.QueryGetCctxByNonceRequest,
) (*types.QueryGetCctxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	tss, found := k.GetObserverKeeper().GetTSS(ctx)
	if !found {
		return nil, status.Error(codes.Internal, "tss not found")
	}
	// #nosec G115 always in range
	cctx, err := getCctxByChainIDAndNonce(k, ctx, tss.TssPubkey, req.ChainID, int64(req.Nonce))
	if err != nil {
		return nil, err
	}

	return &types.QueryGetCctxResponse{CrossChainTx: cctx}, nil
}

// ListPendingCctx returns a list of pending cctxs and the total number of pending cctxs
// a limit for the number of cctxs to return can be specified or the default is MaxPendingCctxs
func (k Keeper) ListPendingCctx(
	c context.Context,
	req *types.QueryListPendingCctxRequest,
) (*types.QueryListPendingCctxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// use default MaxPendingCctxs if not specified or too high
	limit := req.Limit
	if limit == 0 || limit > MaxPendingCctxs {
		limit = MaxPendingCctxs
	}
	ctx := sdk.UnwrapSDKContext(c)

	// query the nonces that are pending
	tss, found := k.zetaObserverKeeper.GetTSS(ctx)
	if !found {
		return nil, observertypes.ErrTssNotFound
	}
	pendingNonces, found := k.GetObserverKeeper().GetPendingNonces(ctx, tss.TssPubkey, req.ChainId)
	if !found {
		return nil, status.Error(codes.Internal, "pending nonces not found")
	}

	cctxs := make([]*types.CrossChainTx, 0)
	maxCCTXsReached := func() bool {
		// #nosec G115 len always positive
		return uint32(len(cctxs)) >= limit
	}

	totalPending := uint64(0)

	// now query the previous nonces up to 1000 prior to find any pending cctx that we might have missed
	// need this logic because a confirmation of higher nonce will automatically update the p.NonceLow
	// therefore might mask some lower nonce cctx that is still pending.
	startNonce := pendingNonces.NonceLow - MaxLookbackNonce
	if startNonce < 0 {
		startNonce = 0
	}
	for i := startNonce; i < pendingNonces.NonceLow; i++ {
		cctx, err := getCctxByChainIDAndNonce(k, ctx, tss.TssPubkey, req.ChainId, i)
		if err != nil {
			return nil, err
		}

		// only take a `limit` number of pending cctxs as result but still count the total pending cctxs
		if IsPending(cctx) {
			totalPending++
			if !maxCCTXsReached() {
				cctxs = append(cctxs, cctx)
			}
		}
	}

	// add the pending nonces to the total pending
	// #nosec G115 always in range
	totalPending += uint64(pendingNonces.NonceHigh - pendingNonces.NonceLow)

	// now query the pending nonces that we know are pending
	for i := pendingNonces.NonceLow; i < pendingNonces.NonceHigh && !maxCCTXsReached(); i++ {
		cctx, err := getCctxByChainIDAndNonce(k, ctx, tss.TssPubkey, req.ChainId, i)
		if err != nil {
			return nil, err
		}
		cctxs = append(cctxs, cctx)
	}

	return &types.QueryListPendingCctxResponse{
		CrossChainTx: cctxs,
		TotalPending: totalPending,
	}, nil
}

// getCctxByChainIDAndNonce returns the cctx by chainID and nonce
func getCctxByChainIDAndNonce(
	k Keeper,
	ctx sdk.Context,
	tssPubkey string,
	chainID int64,
	nonce int64,
) (*types.CrossChainTx, error) {
	nonceToCctx, found := k.GetObserverKeeper().GetNonceToCctx(ctx, tssPubkey, chainID, nonce)
	if !found {
		return nil, status.Error(
			codes.Internal,
			fmt.Sprintf("nonceToCctx not found: chainid %d, nonce %d", chainID, nonce),
		)
	}
	cctx, found := k.GetCrossChainTx(ctx, nonceToCctx.CctxIndex)
	if !found {
		return nil, status.Error(codes.Internal, fmt.Sprintf("cctx not found: index %s", nonceToCctx.CctxIndex))
	}
	return &cctx, nil
}
