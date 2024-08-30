package rpc

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"

	"github.com/zeta-chain/node/pkg/retry"
)

// GetLatestZetaBlock returns the latest zeta block
func (c *Clients) GetLatestZetaBlock(ctx context.Context) (*tmservice.Block, error) {
	res, err := c.Tendermint.GetLatestBlock(ctx, &tmservice.GetLatestBlockRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get latest zeta block")
	}

	return res.SdkBlock, nil
}

// GetNodeInfo returns the node info
func (c *Clients) GetNodeInfo(ctx context.Context) (*tmservice.GetNodeInfoResponse, error) {
	var err error

	res, err := retry.DoTypedWithRetry(func() (*tmservice.GetNodeInfoResponse, error) {
		return c.Tendermint.GetNodeInfo(ctx, &tmservice.GetNodeInfoRequest{})
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to get node info")
	}

	return res, nil
}
