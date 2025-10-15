package zetacore

import (
	"context"
	"sort"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/metrics"
)

func (c *Client) ListPendingCCTX(ctx context.Context, chain chains.Chain) ([]*types.CrossChainTx, uint64, error) {
	list, total, err := c.Clients.ListPendingCCTX(ctx, chain.ChainId)

	if err == nil {
		value := float64(total)

		metrics.PendingTxsPerChain.WithLabelValues(chain.Name).Set(value)
	}

	return list, total, err
}

// GetOutboundTrackers returns all outbound trackers for a chain in ascending order.
func (c *Client) GetOutboundTrackers(ctx context.Context,
	chainID int64,
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

	resp, err := c.Crosschain.OutboundTrackerAllByChain(ctx, in)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all outbound trackers")
	}

	sort.SliceStable(resp.OutboundTracker, func(i, j int) bool {
		return resp.OutboundTracker[i].Nonce < resp.OutboundTracker[j].Nonce
	})

	return resp.OutboundTracker, nil
}
