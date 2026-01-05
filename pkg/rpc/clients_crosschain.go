package rpc

import (
	"context"

	"cosmossdk.io/errors"
	"google.golang.org/grpc"

	"github.com/zeta-chain/node/x/crosschain/types"
)

// 32MB
var maxSizeOption = grpc.MaxCallRecvMsgSize(32 * 1024 * 1024)

// GetBlockHeight returns the zetachain block height
func (c *Clients) GetBlockHeight(ctx context.Context) (int64, error) {
	resp, err := c.Crosschain.LastZetaHeight(ctx, &types.QueryLastZetaHeightRequest{})
	if err != nil {
		return 0, err
	}

	return resp.Height, nil
}

// GetAbortedZetaAmount returns the amount of zeta that has been aborted
func (c *Clients) GetAbortedZetaAmount(ctx context.Context) (string, error) {
	resp, err := c.Crosschain.ZetaAccounting(ctx, &types.QueryZetaAccountingRequest{})
	if err != nil {
		return "", errors.Wrap(err, "failed to get aborted zeta amount")
	}

	return resp.AbortedZetaAmount, nil
}

// GetRateLimiterFlags returns the rate limiter flags
func (c *Clients) GetRateLimiterFlags(ctx context.Context) (types.RateLimiterFlags, error) {
	resp, err := c.Crosschain.RateLimiterFlags(ctx, &types.QueryRateLimiterFlagsRequest{})
	if err != nil {
		return types.RateLimiterFlags{}, errors.Wrap(err, "failed to get rate limiter flags")
	}

	return resp.RateLimiterFlags, nil
}

// GetRateLimiterInput returns input data for the rate limit checker
func (c *Clients) GetRateLimiterInput(ctx context.Context, window int64) (*types.QueryRateLimiterInputResponse, error) {
	in := &types.QueryRateLimiterInputRequest{Window: window}

	resp, err := c.Crosschain.RateLimiterInput(ctx, in, maxSizeOption)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get rate limiter input")
	}

	return resp, nil
}

// GetAllCctx returns all cross chain transactions
func (c *Clients) GetAllCctx(ctx context.Context) ([]*types.CrossChainTx, error) {
	resp, err := c.Crosschain.CctxAll(ctx, &types.QueryAllCctxRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all cross chain transactions")
	}

	return resp.CrossChainTx, nil
}

func (c *Clients) GetCctxByHash(ctx context.Context, sendHash string) (*types.CrossChainTx, error) {
	in := &types.QueryGetCctxRequest{Index: sendHash}
	resp, err := c.Crosschain.Cctx(ctx, in)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get cctx by hash")
	}

	return resp.CrossChainTx, nil
}

// GetCctxByNonce returns a cross chain transaction by nonce
func (c *Clients) GetCctxByNonce(ctx context.Context, chainID int64, nonce uint64) (*types.CrossChainTx, error) {
	resp, err := c.Crosschain.CctxByNonce(ctx, &types.QueryGetCctxByNonceRequest{
		ChainID: chainID,
		Nonce:   nonce,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get cctx by nonce")
	}

	return resp.CrossChainTx, nil
}

// ListPendingCCTXWithinRateLimit returns a list of pending cctxs that do not exceed the outbound rate limit
//   - The max size of the list is crosschainkeeper.MaxPendingCctxs
//   - The returned `rateLimitExceeded` flag indicates if the rate limit is exceeded or not
func (c *Clients) ListPendingCCTXWithinRateLimit(
	ctx context.Context,
) (*types.QueryListPendingCctxWithinRateLimitResponse, error) {
	in := &types.QueryListPendingCctxWithinRateLimitRequest{}

	resp, err := c.Crosschain.ListPendingCctxWithinRateLimit(ctx, in, maxSizeOption)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pending cctxs within rate limit")
	}

	return resp, nil
}

// ListPendingCCTX returns a list of pending cctxs for a given chain
//   - The max size of the list is crosschainkeeper.MaxPendingCctxs
func (c *Clients) ListPendingCCTX(ctx context.Context, chainID int64) ([]*types.CrossChainTx, uint64, error) {
	in := &types.QueryListPendingCctxRequest{ChainId: chainID}

	resp, err := c.Crosschain.ListPendingCctx(ctx, in, maxSizeOption)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get pending cctxs")
	}

	return resp.CrossChainTx, resp.TotalPending, nil
}

// GetOutboundTracker returns the outbound tracker for a chain and nonce
func (c *Clients) GetOutboundTracker(ctx context.Context, chainID int64, nonce uint64) (*types.OutboundTracker, error) {
	in := &types.QueryGetOutboundTrackerRequest{ChainID: chainID, Nonce: nonce}

	resp, err := c.Crosschain.OutboundTracker(ctx, in)
	if err != nil {
		return nil, err
	}

	return &resp.OutboundTracker, nil
}

// GetInboundTrackersForChain returns the inbound trackers for a chain
func (c *Clients) GetInboundTrackersForChain(ctx context.Context, chainID int64) ([]types.InboundTracker, error) {
	in := &types.QueryAllInboundTrackerByChainRequest{ChainId: chainID}

	resp, err := c.Crosschain.InboundTrackerAllByChain(ctx, in)
	if err != nil {
		return nil, err
	}

	return resp.InboundTracker, nil
}
