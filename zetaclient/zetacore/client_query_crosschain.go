package zetacore

import (
	"context"
	"sort"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
)

// 32MB
var maxSizeOption = grpc.MaxCallRecvMsgSize(32 * 1024 * 1024)

// GetLastBlockHeight returns the zetachain block height
func (c *Client) GetLastBlockHeight(ctx context.Context) (uint64, error) {
	resp, err := c.client.crosschain.LastBlockHeight(ctx, &types.QueryGetLastBlockHeightRequest{})
	if err != nil {
		return 0, errors.Wrap(err, "failed to get block height")
	}

	return resp.GetLastBlockHeight().LastInboundHeight, nil
}

// GetBlockHeight returns the zetachain block height
func (c *Client) GetBlockHeight(ctx context.Context) (int64, error) {
	resp, err := c.client.crosschain.LastZetaHeight(ctx, &types.QueryLastZetaHeightRequest{})
	if err != nil {
		return 0, err
	}

	return resp.Height, nil
}

// GetAbortedZetaAmount returns the amount of zeta that has been aborted
func (c *Client) GetAbortedZetaAmount(ctx context.Context) (string, error) {
	resp, err := c.client.crosschain.ZetaAccounting(ctx, &types.QueryZetaAccountingRequest{})
	if err != nil {
		return "", errors.Wrap(err, "failed to get aborted zeta amount")
	}

	return resp.AbortedZetaAmount, nil
}

// GetRateLimiterFlags returns the rate limiter flags
func (c *Client) GetRateLimiterFlags(ctx context.Context) (types.RateLimiterFlags, error) {
	resp, err := c.client.crosschain.RateLimiterFlags(ctx, &types.QueryRateLimiterFlagsRequest{})
	if err != nil {
		return types.RateLimiterFlags{}, errors.Wrap(err, "failed to get rate limiter flags")
	}

	return resp.RateLimiterFlags, nil
}

// GetRateLimiterInput returns input data for the rate limit checker
func (c *Client) GetRateLimiterInput(ctx context.Context, window int64) (*types.QueryRateLimiterInputResponse, error) {
	in := &types.QueryRateLimiterInputRequest{Window: window}

	resp, err := c.client.crosschain.RateLimiterInput(ctx, in, maxSizeOption)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get rate limiter input")
	}

	return resp, nil
}

// GetAllCctx returns all cross chain transactions
func (c *Client) GetAllCctx(ctx context.Context) ([]*types.CrossChainTx, error) {
	resp, err := c.client.crosschain.CctxAll(ctx, &types.QueryAllCctxRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all cross chain transactions")
	}

	return resp.CrossChainTx, nil
}

func (c *Client) GetCctxByHash(ctx context.Context, sendHash string) (*types.CrossChainTx, error) {
	in := &types.QueryGetCctxRequest{Index: sendHash}
	resp, err := c.client.crosschain.Cctx(ctx, in)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get cctx by hash")
	}

	return resp.CrossChainTx, nil
}

// GetCctxByNonce returns a cross chain transaction by nonce
func (c *Client) GetCctxByNonce(ctx context.Context, chainID int64, nonce uint64) (*types.CrossChainTx, error) {
	resp, err := c.client.crosschain.CctxByNonce(ctx, &types.QueryGetCctxByNonceRequest{
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
func (c *Client) ListPendingCCTXWithinRateLimit(
	ctx context.Context,
) (*types.QueryListPendingCctxWithinRateLimitResponse, error) {
	in := &types.QueryListPendingCctxWithinRateLimitRequest{}

	resp, err := c.client.crosschain.ListPendingCctxWithinRateLimit(ctx, in, maxSizeOption)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pending cctxs within rate limit")
	}

	return resp, nil
}

// ListPendingCCTX returns a list of pending cctxs for a given chainID
//   - The max size of the list is crosschainkeeper.MaxPendingCctxs
func (c *Client) ListPendingCCTX(ctx context.Context, chainID int64) ([]*types.CrossChainTx, uint64, error) {
	in := &types.QueryListPendingCctxRequest{ChainId: chainID}

	resp, err := c.client.crosschain.ListPendingCctx(ctx, in, maxSizeOption)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get pending cctxs")
	}

	return resp.CrossChainTx, resp.TotalPending, nil
}

// GetOutboundTracker returns the outbound tracker for a chain and nonce
func (c *Client) GetOutboundTracker(
	ctx context.Context,
	chain chains.Chain,
	nonce uint64,
) (*types.OutboundTracker, error) {
	in := &types.QueryGetOutboundTrackerRequest{ChainID: chain.ChainId, Nonce: nonce}

	resp, err := c.client.crosschain.OutboundTracker(ctx, in)
	if err != nil {
		return nil, err
	}

	return &resp.OutboundTracker, nil
}

// GetInboundTrackersForChain returns the inbound trackers for a chain
func (c *Client) GetInboundTrackersForChain(ctx context.Context, chainID int64) ([]types.InboundTracker, error) {
	in := &types.QueryAllInboundTrackerByChainRequest{ChainId: chainID}

	resp, err := c.client.crosschain.InboundTrackerAllByChain(ctx, in)
	if err != nil {
		return nil, err
	}

	return resp.InboundTracker, nil
}

// GetAllOutboundTrackerByChain returns all outbound trackers for a chain
func (c *Client) GetAllOutboundTrackerByChain(
	ctx context.Context,
	chainID int64,
	order interfaces.Order,
) ([]types.OutboundTracker, error) {
	in := &types.QueryAllOutboundTrackerByChainRequest{
		Chain: chainID,
		Pagination: &query.PageRequest{
			Key:        nil,
			Offset:     0,
			Limit:      2000,
			CountTotal: false,
			Reverse:    false,
		},
	}

	resp, err := c.client.crosschain.OutboundTrackerAllByChain(ctx, in)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all outbound trackers")
	}

	if order == interfaces.Ascending {
		sort.SliceStable(resp.OutboundTracker, func(i, j int) bool {
			return resp.OutboundTracker[i].Nonce < resp.OutboundTracker[j].Nonce
		})
	} else if order == interfaces.Descending {
		sort.SliceStable(resp.OutboundTracker, func(i, j int) bool {
			return resp.OutboundTracker[i].Nonce > resp.OutboundTracker[j].Nonce
		})
	}

	return resp.OutboundTracker, nil
}
