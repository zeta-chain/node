package rpc

import (
	"context"

	"cosmossdk.io/errors"

	"github.com/zeta-chain/node/pkg/proofs"
	"github.com/zeta-chain/node/x/lightclient/types"
)

// GetBlockHeaderEnabledChains returns the enabled chains for block headers
func (c *Clients) GetBlockHeaderEnabledChains(ctx context.Context) ([]types.HeaderSupportedChain, error) {
	resp, err := c.Lightclient.HeaderEnabledChains(ctx, &types.QueryHeaderEnabledChainsRequest{})
	if err != nil {
		return []types.HeaderSupportedChain{}, err
	}

	return resp.HeaderEnabledChains, nil
}

// GetBlockHeaderChainState returns the block header chain state
func (c *Clients) GetBlockHeaderChainState(ctx context.Context, chainID int64) (*types.ChainState, error) {
	in := &types.QueryGetChainStateRequest{ChainId: chainID}

	resp, err := c.Lightclient.ChainState(ctx, in)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get chain state")
	}

	return resp.ChainState, nil
}

// Prove returns whether a proof is valid
func (c *Clients) Prove(
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

	resp, err := c.Lightclient.Prove(ctx, in)
	if err != nil {
		return false, errors.Wrap(err, "failed to prove")
	}

	return resp.Valid, nil
}
