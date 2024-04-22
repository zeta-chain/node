package keeper

import (
	"context"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListPendingCctxWithinRateLimit returns a list of pending cctxs that do not exceed the outbound rate limit
// a limit for the number of cctxs to return can be specified or the default is MaxPendingCctxs
func (k Keeper) ListPendingCctxWithinRateLimit(c context.Context, req *types.QueryListPendingCctxWithinRateLimitRequest) (*types.QueryListPendingCctxWithinRateLimitResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// use default MaxPendingCctxs if not specified or too high
	limit := req.Limit
	if limit == 0 || limit > MaxPendingCctxs {
		limit = MaxPendingCctxs
	}
	ctx := sdk.UnwrapSDKContext(c)

	// define a few variables to be used in the query loops
	limitExceeded := false
	totalPending := uint64(0)
	totalCctxValueInZeta := sdk.NewDec(0)
	cctxs := make([]*types.CrossChainTx, 0)
	chains := k.zetaObserverKeeper.GetSupportedForeignChains(ctx)

	// check rate limit flags to decide if we should apply rate limit
	applyLimit := true
	rateLimitFlags, found := k.GetRateLimiterFlags(ctx)
	if !found || !rateLimitFlags.Enabled {
		applyLimit = false
	}

	// fallback to non-rate-limited query if rate limiter is disabled
	if !applyLimit {
		for _, chain := range chains {
			resp, err := k.ListPendingCctx(ctx, &types.QueryListPendingCctxRequest{ChainId: chain.ChainId, Limit: limit})
			if err == nil {
				cctxs = append(cctxs, resp.CrossChainTx...)
				totalPending += resp.TotalPending
			}
		}
		return &types.QueryListPendingCctxWithinRateLimitResponse{
			CrossChainTx:      cctxs,
			TotalPending:      totalPending,
			RateLimitExceeded: false,
		}, nil
	}

	// get current height and tss
	height := ctx.BlockHeight()
	if height <= 0 {
		return nil, status.Error(codes.OutOfRange, "height out of range")
	}
	tss, found := k.zetaObserverKeeper.GetTSS(ctx)
	if !found {
		return nil, observertypes.ErrTssNotFound
	}

	// calculate the rate limiter sliding window left boundary (inclusive)
	leftWindowBoundary := height - rateLimitFlags.Window
	if leftWindowBoundary < 0 {
		leftWindowBoundary = 0
	}

	// get the conversion rates for all foreign coins
	var gasCoinRates map[int64]sdk.Dec
	var erc20CoinRates map[int64]map[string]sdk.Dec
	var foreignCoinMap map[int64]map[string]fungibletypes.ForeignCoins
	var rateLimitInZeta sdk.Dec
	if applyLimit {
		gasCoinRates, erc20CoinRates = k.GetRateLimiterRates(ctx)
		foreignCoinMap = k.fungibleKeeper.GetAllForeignCoinMap(ctx)
		rateLimitInZeta = sdk.NewDecFromBigInt(rateLimitFlags.Rate.BigInt())
	}

	// the criteria to stop adding cctxs to the rpc response
	maxCCTXsReached := func() bool {
		// #nosec G701 len always positive
		return uint32(len(cctxs)) >= limit
	}

	// if a cctx falls within the rate limiter window
	isCctxInWindow := func(cctx *types.CrossChainTx) bool {
		// #nosec G701 checked positive
		return cctx.InboundTxParams.InboundTxObservedExternalHeight >= uint64(leftWindowBoundary)
	}

	// query pending nonces for each foreign chain
	// Note: The pending nonces could change during the RPC call, so query them beforehand
	pendingNoncesMap := make(map[int64]*observertypes.PendingNonces)
	for _, chain := range chains {
		pendingNonces, found := k.GetObserverKeeper().GetPendingNonces(ctx, tss.TssPubkey, chain.ChainId)
		if !found {
			return nil, status.Error(codes.Internal, "pending nonces not found")
		}
		pendingNoncesMap[chain.ChainId] = &pendingNonces
	}

	// query backwards for potential missed pending cctxs for each foreign chain
LoopBackwards:
	for _, chain := range chains {
		// we should at least query 1000 prior to find any pending cctx that we might have missed
		// this logic is needed because a confirmation of higher nonce will automatically update the p.NonceLow
		// therefore might mask some lower nonce cctx that is still pending.
		pendingNonces := pendingNoncesMap[chain.ChainId]
		startNonce := pendingNonces.NonceLow - 1
		endNonce := pendingNonces.NonceLow - 1000
		if endNonce < 0 {
			endNonce = 0
		}

		// add the pending nonces to the total pending
		// Note: the `totalPending` may not be accurate only if the rate limiter triggers early exit
		// `totalPending` is now used for metrics only and it's okay to trade off accuracy for performance
		// #nosec G701 always in range
		totalPending += uint64(pendingNonces.NonceHigh - pendingNonces.NonceLow)

		// query cctx by nonce backwards to the left boundary of the rate limit sliding window
		for nonce := startNonce; nonce >= 0; nonce-- {
			cctx, err := getCctxByChainIDAndNonce(k, ctx, tss.TssPubkey, chain.ChainId, nonce)
			if err != nil {
				return nil, err
			}
			inWindow := isCctxInWindow(cctx)

			// We should at least go backwards by 1000 nonces to pick up missed pending cctxs
			// We might go even further back if rate limiter is enabled and the endNonce hasn't hit the left window boundary yet
			// There are two criteria to stop scanning backwards:
			// criteria #1: we'll stop at the left window boundary if the `endNonce` hasn't hit it yet
			if nonce < endNonce && !inWindow {
				break
			}
			// criteria #2: we should finish the RPC call if the rate limit is exceeded
			if inWindow && rateLimitExceeded(chain.ChainId, cctx, gasCoinRates, erc20CoinRates, foreignCoinMap, &totalCctxValueInZeta, rateLimitInZeta) {
				limitExceeded = true
				break LoopBackwards
			}

			// only take a `limit` number of pending cctxs as result but still count the total pending cctxs
			if IsPending(cctx) {
				totalPending++
				if !maxCCTXsReached() {
					cctxs = append(cctxs, cctx)
				}
			}
		}
	}

	// query forwards for pending cctxs for each foreign chain
LoopForwards:
	for _, chain := range chains {
		// query the pending cctxs in range [NonceLow, NonceHigh)
		pendingNonces := pendingNoncesMap[chain.ChainId]
		for nonce := pendingNonces.NonceLow; nonce < pendingNonces.NonceHigh; nonce++ {
			cctx, err := getCctxByChainIDAndNonce(k, ctx, tss.TssPubkey, chain.ChainId, nonce)
			if err != nil {
				return nil, err
			}
			inWindow := isCctxInWindow(cctx)

			// only take a `limit` number of pending cctxs as result
			if maxCCTXsReached() {
				break LoopForwards
			}
			// criteria #2: we should finish the RPC call if the rate limit is exceeded
			if inWindow && rateLimitExceeded(chain.ChainId, cctx, gasCoinRates, erc20CoinRates, foreignCoinMap, &totalCctxValueInZeta, rateLimitInZeta) {
				limitExceeded = true
				break LoopForwards
			}
			cctxs = append(cctxs, cctx)
		}
	}

	// sort the cctxs by chain ID and nonce (lower nonce holds higher priority for scheduling)
	sort.Slice(cctxs, func(i, j int) bool {
		if cctxs[i].GetCurrentOutTxParam().ReceiverChainId == cctxs[j].GetCurrentOutTxParam().ReceiverChainId {
			return cctxs[i].GetCurrentOutTxParam().OutboundTxTssNonce < cctxs[j].GetCurrentOutTxParam().OutboundTxTssNonce
		}
		return cctxs[i].GetCurrentOutTxParam().ReceiverChainId < cctxs[j].GetCurrentOutTxParam().ReceiverChainId
	})

	return &types.QueryListPendingCctxWithinRateLimitResponse{
		CrossChainTx:      cctxs,
		TotalPending:      totalPending,
		ValueWithinWindow: totalCctxValueInZeta.TruncateInt().Uint64(),
		RateLimitExceeded: limitExceeded,
	}, nil
}

// ConvertCctxValue converts the value of the cctx in ZETA using given conversion rates
func ConvertCctxValue(
	chainID int64,
	cctx *types.CrossChainTx,
	gasCoinRates map[int64]sdk.Dec,
	erc20CoinRates map[int64]map[string]sdk.Dec,
	foreignCoinMap map[int64]map[string]fungibletypes.ForeignCoins,
) sdk.Dec {
	var rate sdk.Dec
	var decimals uint64
	switch cctx.InboundTxParams.CoinType {
	case coin.CoinType_Zeta:
		// no conversion needed for ZETA
		amountCctx := sdk.NewDecFromBigInt(cctx.GetCurrentOutTxParam().Amount.BigInt())
		return amountCctx.Quo(sdk.NewDec(10).Power(18))
	case coin.CoinType_Gas:
		rate = gasCoinRates[chainID]
	case coin.CoinType_ERC20:
		// get the ERC20 coin rate
		_, found := erc20CoinRates[chainID]
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
	if rate.IsNil() || rate.LTE(sdk.NewDec(0)) {
		return sdk.NewDec(0)
	}

	// get foreign coin decimals
	foreignCoinFromChainMap, found := foreignCoinMap[chainID]
	if !found {
		// skip if no coin found for this chainID
		return sdk.NewDec(0)
	}
	foreignCoin, found := foreignCoinFromChainMap[strings.ToLower(cctx.InboundTxParams.Asset)]
	if !found {
		// skip if no coin found for this Asset
		return sdk.NewDec(0)
	}
	decimals = uint64(foreignCoin.Decimals)

	// given decimals = 6, the `oneZrc20` amount will be 10^6 = 1000000
	oneZrc20 := sdk.NewDec(10).Power(decimals)

	// convert asset amount into ZETA
	// step 1: convert the amount into ZRC20 integer amount
	// step 2: convert the ZRC20 integer amount into decimal amount
	// given amountCctx = 2000000, rate = 0.8, decimals = 6
	// the amountZrc20 = 2000000 * 0.8 = 1600000, the amountZeta = 1600000 / 1000000 = 1.6
	amountCctx := sdk.NewDecFromBigInt(cctx.GetCurrentOutTxParam().Amount.BigInt())
	amountZrc20 := amountCctx.Mul(rate)
	amountZeta := amountZrc20.Quo(oneZrc20)
	return amountZeta
}

// rateLimitExceeded accumulates the cctx value and then checks if the rate limit is exceeded
// returns true if the rate limit is exceeded
func rateLimitExceeded(
	chainID int64,
	cctx *types.CrossChainTx,
	gasCoinRates map[int64]sdk.Dec,
	erc20CoinRates map[int64]map[string]sdk.Dec,
	foreignCoinMap map[int64]map[string]fungibletypes.ForeignCoins,
	currentCctxValue *sdk.Dec,
	rateLimitValue sdk.Dec,
) bool {
	amountZeta := ConvertCctxValue(chainID, cctx, gasCoinRates, erc20CoinRates, foreignCoinMap)
	*currentCctxValue = currentCctxValue.Add(amountZeta)
	return currentCctxValue.GT(rateLimitValue)
}
