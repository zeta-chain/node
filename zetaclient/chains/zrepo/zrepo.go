package zrepo

import (
	"context"
	"errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
)

// Each chain must instantiate its own ZetaRepo.
type ZetaRepo struct {
	client interfaces.ZetacoreClient

	connectedChain chains.Chain
}

func New(client interfaces.ZetacoreClient, connectedChain chains.Chain) *ZetaRepo {
	return &ZetaRepo{
		client:         client,
		connectedChain: connectedChain,
	}
}

func (repo *ZetaRepo) GetCCTX(ctx context.Context, nonce uint64) (*types.CrossChainTx, error) {
	return repo.client.GetCctxByNonce(ctx, repo.connectedChain.ChainId, nonce)
}

func (repo *ZetaRepo) GetOutboundTrackers(ctx context.Context) ([]types.OutboundTracker, error) {
	trackers, err := repo.client.GetOutboundTrackers(ctx, repo.connectedChain.ChainId)
	if err != nil {
		return nil, errors.Join(ErrGetOutboundTrackers, err)
	}
	return trackers, nil
}

func (repo *ZetaRepo) VoteOutbound(ctx context.Context,
	gasLimit uint64,
	retryGasLimit uint64,
	msg *types.MsgVoteOutbound,
) (string, string, error) {
	return repo.client.PostVoteOutbound(ctx, gasLimit, retryGasLimit, msg)
}
