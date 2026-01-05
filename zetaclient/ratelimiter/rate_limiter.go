// Package ratelimiter provides functionalities for rate limiting the cross-chain transactions
package ratelimiter

import (
	sdkmath "cosmossdk.io/math"

	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// Input is the input data for the rate limiter
type Input struct {
	// zeta chain height
	Height int64

	// the missed cctxs in range [?, NonceLow) across all chains
	CctxsMissed []*crosschaintypes.CrossChainTx

	// the pending cctxs in range [NonceLow, NonceHigh) across all chains
	CctxsPending []*crosschaintypes.CrossChainTx

	// the total value of the past cctxs within window across all chains
	PastCctxsValue sdkmath.Int

	// the total value of the pending cctxs across all chains
	PendingCctxsValue sdkmath.Int

	// the lowest height of the pending (not missed) cctxs across all chains
	LowestPendingCctxHeight int64
}

// Output is the output data for the rate limiter
type Output struct {
	// the cctxs to be scheduled after rate limit check
	CctxsMap map[int64][]*crosschaintypes.CrossChainTx

	// the current sliding window within which the withdrawals are considered by the rate limiter
	CurrentWithdrawWindow int64

	// the current withdraw rate (azeta/block) within the current sliding window
	CurrentWithdrawRate sdkmath.Int

	// whether the current withdraw rate exceeds the given rate limit or not
	RateLimitExceeded bool
}

// NewInput creates a rate limiter input from gRPC response
func NewInput(resp crosschaintypes.QueryRateLimiterInputResponse) (*Input, bool) {
	// parse the past cctxs value from string
	pastCctxsValue, ok := sdkmath.NewIntFromString(resp.PastCctxsValue)
	if !ok {
		return nil, false
	}

	// parse the pending cctxs value from string
	pendingCctxsValue, ok := sdkmath.NewIntFromString(resp.PendingCctxsValue)
	if !ok {
		return nil, false
	}

	return &Input{
		Height:                  resp.Height,
		CctxsMissed:             resp.CctxsMissed,
		CctxsPending:            resp.CctxsPending,
		PastCctxsValue:          pastCctxsValue,
		PendingCctxsValue:       pendingCctxsValue,
		LowestPendingCctxHeight: resp.LowestPendingCctxHeight,
	}, true
}

// IsRateLimiterUsable checks if the rate limiter is usable or not
func IsRateLimiterUsable(rateLimiterFlags crosschaintypes.RateLimiterFlags) bool {
	if !rateLimiterFlags.Enabled {
		return false
	}
	if rateLimiterFlags.Window <= 0 {
		return false
	}
	if rateLimiterFlags.Rate.IsNil() {
		return false
	}
	if rateLimiterFlags.Rate.IsZero() {
		return false
	}
	return true
}

// ApplyRateLimiter applies the rate limiter to the input and produces output
func ApplyRateLimiter(input *Input, window int64, rate sdkmath.Uint) *Output {
	// block limit and the window limit in azeta
	blockLimitInAzeta := sdkmath.NewIntFromBigInt(rate.BigInt())
	windowLimitInAzeta := blockLimitInAzeta.Mul(sdkmath.NewInt(window))

	// invariant: for period of time >= `window`, the zetaclient-side average withdraw rate should be <= `blockLimitInZeta`
	// otherwise, zetaclient should wait for the average rate to drop below `blockLimitInZeta`
	withdrawWindow := window
	withdrawLimitInAzeta := windowLimitInAzeta
	if input.LowestPendingCctxHeight != 0 {
		// If [input.LowestPendingCctxHeight, input.Height] is wider than the given `window`, we should:
		// 1. use the wider window to calculate the average withdraw rate
		// 2. adjust the limit proportionally to fit the wider window
		pendingCctxWindow := input.Height - input.LowestPendingCctxHeight + 1
		if pendingCctxWindow > window {
			withdrawWindow = pendingCctxWindow
			withdrawLimitInAzeta = blockLimitInAzeta.Mul(sdkmath.NewInt(pendingCctxWindow))
		}
	}

	// limit exceeded or not
	totalWithdrawInAzeta := input.PastCctxsValue.Add(input.PendingCctxsValue)
	limitExceeded := totalWithdrawInAzeta.GT(withdrawLimitInAzeta)

	// define the result cctx map to be scheduled
	cctxMap := make(map[int64][]*crosschaintypes.CrossChainTx)

	// addCctxsToMap adds the given cctxs to the cctx map
	addCctxsToMap := func(cctxs []*crosschaintypes.CrossChainTx) {
		for _, cctx := range cctxs {
			chainID := cctx.GetCurrentOutboundParam().ReceiverChainId
			if _, found := cctxMap[chainID]; !found {
				cctxMap[chainID] = make([]*crosschaintypes.CrossChainTx, 0)
			}
			cctxMap[chainID] = append(cctxMap[chainID], cctx)
		}
	}

	// schedule missed cctxs regardless of the `limitExceeded` flag
	addCctxsToMap(input.CctxsMissed)

	// schedule pending cctxs only if `limitExceeded == false`
	if !limitExceeded {
		addCctxsToMap(input.CctxsPending)
	}

	return &Output{
		CctxsMap:              cctxMap,
		CurrentWithdrawWindow: withdrawWindow,
		CurrentWithdrawRate:   totalWithdrawInAzeta.Quo(sdkmath.NewInt(withdrawWindow)),
		RateLimitExceeded:     limitExceeded,
	}
}
