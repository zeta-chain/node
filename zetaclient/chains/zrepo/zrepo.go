// Package zrepo provides an abstraction layer for interactions with the zetacore client.
package zrepo

import (
	"context"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/crosschain/types"
)

// ZetaRepo implements the Repository pattern by wrapping a zetacore client.
// Each chain module must instantiate its own ZetaRepo.
type ZetaRepo struct {
	client ZetacoreClient

	connectedChain chains.Chain
}

// New constructs a new ZetaRepo object.
func New(client ZetacoreClient, connectedChain chains.Chain) *ZetaRepo {
	return &ZetaRepo{client, connectedChain}
}

func (repo *ZetaRepo) GetCCTX(ctx context.Context, nonce uint64) (*types.CrossChainTx, error) {
	cctx, err := repo.client.GetCctxByNonce(ctx, repo.connectedChain.ChainId, nonce)
	if err != nil {
		return nil, newClientError(ErrClientGetCCTXByNonce, err)
	}
	return cctx, nil
}

func (repo *ZetaRepo) GetOutboundTrackers(ctx context.Context) ([]types.OutboundTracker, error) {
	trackers, err := repo.client.GetOutboundTrackers(ctx, repo.connectedChain.ChainId)
	if err != nil {
		return nil, newClientError(ErrClientGetOutboundTrackers, err)
	}
	return trackers, nil
}

// VoteOutbound votes on an outbound.
// It returns the hash of the transaction and the index of the ballot.
func (repo *ZetaRepo) VoteOutbound(ctx context.Context,
	gasLimit uint64,
	retryGasLimit uint64,
	msg *types.MsgVoteOutbound,
) (string, string, error) {
	hash, ballot, err := repo.client.PostVoteOutbound(ctx, gasLimit, retryGasLimit, msg)
	if err != nil {
		return "", "", newClientError(ErrClientPostVoteOutbound, err)
	}
	return hash, ballot, nil
}
