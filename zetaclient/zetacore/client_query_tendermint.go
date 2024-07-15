package zetacore

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"

	"github.com/zeta-chain/zetacore/pkg/retry"
)

// GetLatestZetaBlock returns the latest zeta block
func (c *Client) GetLatestZetaBlock(ctx context.Context) (*tmservice.Block, error) {
	res, err := c.client.tendermint.GetLatestBlock(ctx, &tmservice.GetLatestBlockRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get latest zeta block")
	}

	return res.SdkBlock, nil
}

// GetNodeInfo returns the node info
func (c *Client) GetNodeInfo(ctx context.Context) (*tmservice.GetNodeInfoResponse, error) {
	var err error

	res, err := retry.DoTypedWithRetry(func() (*tmservice.GetNodeInfoResponse, error) {
		return c.client.tendermint.GetNodeInfo(ctx, &tmservice.GetNodeInfoRequest{})
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to get node info")
	}

	return res, nil
}
