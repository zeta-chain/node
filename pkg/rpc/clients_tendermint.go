package rpc

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"

	"github.com/zeta-chain/node/pkg/retry"
)

// GetLatestZetaBlock returns the latest zeta block
func (c *Clients) GetLatestZetaBlock(ctx context.Context) (*cmtservice.Block, error) {
	res, err := c.Tendermint.GetLatestBlock(ctx, &cmtservice.GetLatestBlockRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get latest zeta block")
	}

	return res.SdkBlock, nil
}

// GetNodeInfo returns the node info
func (c *Clients) GetNodeInfo(ctx context.Context) (*cmtservice.GetNodeInfoResponse, error) {
	var err error

	res, err := retry.DoTypedWithRetry(func() (*cmtservice.GetNodeInfoResponse, error) {
		return c.Tendermint.GetNodeInfo(ctx, &cmtservice.GetNodeInfoRequest{})
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to get node info")
	}

	return res, nil
}

func (c *Clients) GetSyncing(ctx context.Context) (bool, error) {
	res, err := c.Tendermint.GetSyncing(ctx, &cmtservice.GetSyncingRequest{})
	if err != nil {
		return false, errors.Wrap(err, "failed to get syncing status")
	}
	return res.Syncing, nil

}
