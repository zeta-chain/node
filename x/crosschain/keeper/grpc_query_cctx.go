package keeper

import (
	"context"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
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
	res, found := k.GetObserverKeeper().GetNonceToCctx(ctx, tss.TssPubkey, req.ChainID, int64(req.Nonce))
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
// a limit for the number of cctxs to return can be specified or the default is MaxPendingCctxs
func (k Keeper) CctxListPending(c context.Context, req *types.QueryListCctxPendingRequest) (*types.QueryListCctxPendingResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// check limit and use default MaxPendingCctxs if not specified
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
	pendingNonces, found := k.GetObserverKeeper().GetPendingNonces(ctx, tss.TssPubkey, req.ChainId)
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
		nonceToCctx, found := k.GetObserverKeeper().GetNonceToCctx(ctx, tss.TssPubkey, req.ChainId, i)
		if !found {
			return nil, status.Error(codes.Internal, fmt.Sprintf("nonceToCctx not found: nonce %d, chainid %d", i, req.ChainId))
		}
		cctx, found := k.GetCrossChainTx(ctx, nonceToCctx.CctxIndex)
		if !found {
			return nil, status.Error(codes.Internal, fmt.Sprintf("cctx not found: index %s", nonceToCctx.CctxIndex))
		}

		// only take a `limit` number of pending cctxs as result but still count the total pending cctxs
		if IsPending(cctx) {
			totalPending++
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
		nonceToCctx, found := k.GetObserverKeeper().GetNonceToCctx(ctx, tss.TssPubkey, req.ChainId, i)
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

// CctxListPendingWithinRateLimit returns a list of pending cctxs that do not exceed the outbound rate limit
// a limit for the number of cctxs to return can be specified or the default is MaxPendingCctxs
func (k Keeper) CctxListPendingWithinRateLimit(c context.Context, req *types.QueryListCctxPendingWithRateLimitRequest) (*types.QueryListCctxPendingWithRateLimitResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// check limit and use default MaxPendingCctxs if not specified
	if req.Limit > MaxPendingCctxs {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("limit exceeds max limit of %d", MaxPendingCctxs))
	}
	limit := req.Limit
	if limit == 0 {
		limit = MaxPendingCctxs
	}

	// get current height and tss
	ctx := sdk.UnwrapSDKContext(c)
	height := ctx.BlockHeight()
	if height <= 0 {
		return nil, status.Error(codes.OutOfRange, "height out of range")
	}
	tss, found := k.zetaObserverKeeper.GetTSS(ctx)
	if !found {
		return nil, status.Error(codes.Internal, "tss not found")
	}

	// check rate limit flags to decide if we should apply rate limit
	applyLimit := true
	rateLimitFlags, found := k.GetRatelimiterFlags(ctx)
	if !found || !rateLimitFlags.Enabled {
		applyLimit = false
	}

	// calculate the rate limiter sliding window left boundary (inclusive)
	leftWindowBoundary := height - rateLimitFlags.Window
	if leftWindowBoundary < 0 {
		leftWindowBoundary = 0
	}

	// get the conversion rates for all foreign coins
	var gasCoinRates map[int64]sdk.Dec
	var erc20CoinRates map[int64]map[string]sdk.Dec
	var erc20Coins map[int64]map[string]fungibletypes.ForeignCoins
	var rateLimitInZeta sdk.Dec
	if applyLimit {
		gasCoinRates, erc20CoinRates = k.GetRatelimiterRates(ctx)
		erc20Coins = k.fungibleKeeper.GetAllForeignERC20CoinMap(ctx)
		rateLimitInZeta = sdk.NewDecFromBigInt(rateLimitFlags.Rate.BigInt())
	}

	// define a few variables to be used in the below loops
	limitExceeded := false
	totalPending := uint64(0)
	totalCctxValueInZeta := sdk.NewDec(0)
	cctxs := make([]*types.CrossChainTx, 0)
	pendingNoncesMap := make(map[int64]*observertypes.PendingNonces)

	// the criteria to stop adding cctxs to the rpc response
	maxCCTXsReached := func() bool {
		// #nosec G701 len always positive
		return uint32(len(cctxs)) >= limit
	}

	// query pending nonces for each supported chain
	// Note: The pending nonces could change during the RPC call, so query them beforehand
	chains := k.zetaObserverKeeper.GetSupportedChains(ctx)
	for _, chain := range chains {
		if chain.IsZetaChain() {
			continue
		}
		pendingNonces, found := k.GetObserverKeeper().GetPendingNonces(ctx, tss.TssPubkey, chain.ChainId)
		if !found {
			return nil, status.Error(codes.Internal, "pending nonces not found")
		}
		pendingNoncesMap[chain.ChainId] = &pendingNonces
	}

	// query backwards for potential missed pending cctxs for each supported chain
LoopBackwards:
	for _, chain := range chains {
		if chain.IsZetaChain() {
			continue
		}

		// we should at least query 1000 prior to find any pending cctx that we might have missed
		// this logic is needed because a confirmation of higher nonce will automatically update the p.NonceLow
		// therefore might mask some lower nonce cctx that is still pending.
		pendingNonces := pendingNoncesMap[chain.ChainId]
		startNonce := pendingNonces.NonceLow - 1
		endNonce := pendingNonces.NonceLow - 1000
		if endNonce < 0 {
			endNonce = 0
		}

		// query cctx by nonce backwards to the left boundary of the rate limit sliding window
		for nonce := startNonce; nonce >= 0; nonce-- {
			nonceToCctx, found := k.GetObserverKeeper().GetNonceToCctx(ctx, tss.TssPubkey, chain.ChainId, nonce)
			if !found {
				return nil, status.Error(codes.Internal, fmt.Sprintf("nonceToCctx not found: chainid %d, nonce %d", chain.ChainId, nonce))
			}
			cctx, found := k.GetCrossChainTx(ctx, nonceToCctx.CctxIndex)
			if !found {
				return nil, status.Error(codes.Internal, fmt.Sprintf("cctx not found: index %s", nonceToCctx.CctxIndex))
			}

			// We should at least go backwards by 1000 nonces to pick up missed pending cctxs
			// We might go even further back if rate limiter is enabled and the endNonce hasn't hit the left window boundary yet
			// There are three criteria to stop scanning backwards:
			// criteria #1: if rate limiter is disabled, we should stop exactly on the `endNonce`
			if !applyLimit && nonce < endNonce {
				break
			}
			if applyLimit {
				// criteria #2: if rate limiter is enabled, we'll stop at the left window boundary if the `endNonce` hasn't hit it yet
				// #nosec G701 always positive
				if nonce < endNonce && cctx.InboundTxParams.InboundTxObservedExternalHeight < uint64(leftWindowBoundary) {
					break
				}
				// criteria #3: if rate limiter is enabled, we should finish the RPC call if the rate limit is exceeded
				if rateLimitExceeded(chain.ChainId, &cctx, gasCoinRates, erc20CoinRates, erc20Coins, &totalCctxValueInZeta, rateLimitInZeta) {
					limitExceeded = true
					break LoopBackwards
				}
			}

			// only take a `limit` number of pending cctxs as result but still count the total pending cctxs
			if IsPending(cctx) {
				totalPending++
				if !maxCCTXsReached() {
					cctxs = append(cctxs, &cctx)
				}
			}
		}

		// add the pending nonces to the total pending
		// Note: the `totalPending` may not be accurate only if the rate limiter triggers early exit
		// `totalPending` is now used for metrics only and it's okay to trade off accuracy for performance
		// #nosec G701 always in range
		totalPending += uint64(pendingNonces.NonceHigh - pendingNonces.NonceLow)
	}

	// query forwards for pending cctxs for each supported chain
LoopForwards:
	for _, chain := range chains {
		if chain.IsZetaChain() {
			continue
		}

		// query the pending cctxs in range [NonceLow, NonceHigh)
		pendingNonces := pendingNoncesMap[chain.ChainId]
		for i := pendingNonces.NonceLow; i < pendingNonces.NonceHigh; i++ {
			nonceToCctx, found := k.GetObserverKeeper().GetNonceToCctx(ctx, tss.TssPubkey, chain.ChainId, i)
			if !found {
				return nil, status.Error(codes.Internal, "nonceToCctx not found")
			}
			cctx, found := k.GetCrossChainTx(ctx, nonceToCctx.CctxIndex)
			if !found {
				return nil, status.Error(codes.Internal, "cctxIndex not found")
			}

			// only take a `limit` number of pending cctxs as result
			if maxCCTXsReached() {
				break LoopForwards
			}
			// criteria #3: if rate limiter is enabled, we should finish the RPC call if the rate limit is exceeded
			if applyLimit && rateLimitExceeded(chain.ChainId, &cctx, gasCoinRates, erc20CoinRates, erc20Coins, &totalCctxValueInZeta, rateLimitInZeta) {
				limitExceeded = true
				break LoopForwards
			}
			cctxs = append(cctxs, &cctx)
		}
	}

	return &types.QueryListCctxPendingWithRateLimitResponse{
		CrossChainTx:      cctxs,
		TotalPending:      totalPending,
		RateLimitExceeded: limitExceeded,
	}, nil
}

// convertCctxValue converts the value of the cctx in ZETA using given conversion rates
func convertCctxValue(
	chainID int64,
	cctx *types.CrossChainTx,
	gasCoinRates map[int64]sdk.Dec,
	erc20CoinRates map[int64]map[string]sdk.Dec,
	erc20Coins map[int64]map[string]fungibletypes.ForeignCoins,
) sdk.Dec {
	var rate sdk.Dec
	var decimals uint64
	switch cctx.InboundTxParams.CoinType {
	case coin.CoinType_Zeta:
		// no conversion needed for ZETA
		rate = sdk.NewDec(1)
	case coin.CoinType_Gas:
		rate = gasCoinRates[chainID]
	case coin.CoinType_ERC20:
		// get the ERC20 coin decimals
		_, found := erc20Coins[chainID]
		if !found {
			// skip if no coin found for this chainID
			return sdk.NewDec(0)
		}
		fCoin, found := erc20Coins[chainID][strings.ToLower(cctx.InboundTxParams.Asset)]
		if !found {
			// skip if no coin found for this Asset
			return sdk.NewDec(0)
		}
		// #nosec G701 always in range
		decimals = uint64(fCoin.Decimals)

		// get the ERC20 coin rate
		_, found = erc20CoinRates[chainID]
		if !found {
			// skip if no rate found for this chainID
			return sdk.NewDec(0)
		}
		rate = erc20CoinRates[chainID][strings.ToLower(cctx.InboundTxParams.Asset)]
	default:
		// skip CoinType_Cmd
		return sdk.NewDec(0)
	}

	// should not happen, return 0 to skip if it happens
	if rate.LTE(sdk.NewDec(0)) {
		return sdk.NewDec(0)
	}

	// the reciprocal of `rate` is the amount of zrc20 needed to buy 1 ZETA
	// for example, given rate = 0.8, the reciprocal is 1.25, which means 1.25 ZRC20 can buy 1 ZETA
	// given decimals = 6, the `oneZeta` amount will be 1.25 * 10^6 = 1250000
	oneZrc20 := sdk.NewDec(1).Power(decimals)
	oneZeta := oneZrc20.Quo(rate)

	// convert asset amount into ZETA
	amountCctx := sdk.NewDecFromBigInt(cctx.InboundTxParams.Amount.BigInt())
	amountZeta := amountCctx.Quo(oneZeta)
	return amountZeta
}

// rateLimitExceeded accumulates the cctx value and then checks if the rate limit is exceeded
// returns true if the rate limit is exceeded
func rateLimitExceeded(
	chainID int64,
	cctx *types.CrossChainTx,
	gasCoinRates map[int64]sdk.Dec,
	erc20CoinRates map[int64]map[string]sdk.Dec,
	erc20Coins map[int64]map[string]fungibletypes.ForeignCoins,
	currentCctxValue *sdk.Dec,
	rateLimitValue sdk.Dec,
) bool {
	amountZeta := convertCctxValue(chainID, cctx, gasCoinRates, erc20CoinRates, erc20Coins)
	*currentCctxValue = currentCctxValue.Add(amountZeta)
	return currentCctxValue.GT(rateLimitValue)
}
