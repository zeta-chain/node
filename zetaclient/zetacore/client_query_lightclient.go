package zetacore

import (
	"context"

	"cosmossdk.io/errors"

	"github.com/zeta-chain/zetacore/pkg/proofs"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

// GetBlockHeaderEnabledChains returns the enabled chains for block headers
func (c *Client) GetBlockHeaderEnabledChains(ctx context.Context) ([]types.HeaderSupportedChain, error) {
	resp, err := c.client.light.HeaderEnabledChains(ctx, &types.QueryHeaderEnabledChainsRequest{})
	if err != nil {
		return []types.HeaderSupportedChain{}, err
	}

	return resp.HeaderEnabledChains, nil
}

// GetBlockHeaderChainState returns the block header chain state
func (c *Client) GetBlockHeaderChainState(ctx context.Context, chainID int64) (*types.ChainState, error) {
	in := &types.QueryGetChainStateRequest{ChainId: chainID}

	resp, err := c.client.light.ChainState(ctx, in)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get chain state")
	}

	return resp.ChainState, nil
}

// Prove returns whether a proof is valid
func (c *Client) Prove(
	ctx context.Context,
	blockHash string,
	txHash string,
	txIndex int64,
	proof *proofs.Proof,
	chainID int64,
) (bool, error) {
	in := &types.QueryProveRequest{
		BlockHash: blockHash,
		TxIndex:   txIndex,
		Proof:     proof,
		ChainId:   chainID,
		TxHash:    txHash,
	}

	resp, err := c.client.light.Prove(ctx, in)
	if err != nil {
		return false, errors.Wrap(err, "failed to prove")
	}

	return resp.Valid, nil
}
