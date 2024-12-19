package zetacore

import (
	"context"
	"sort"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
)

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

	resp, err := c.Crosschain.OutboundTrackerAllByChain(ctx, in)
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
