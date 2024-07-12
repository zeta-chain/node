package zetacore

import (
	"context"

	"cosmossdk.io/errors"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/retry"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// GetCrosschainFlags returns the crosschain flags
func (c *Client) GetCrosschainFlags(ctx context.Context) (types.CrosschainFlags, error) {
	resp, err := c.client.observer.CrosschainFlags(ctx, &types.QueryGetCrosschainFlagsRequest{})
	if err != nil {
		return types.CrosschainFlags{}, err
	}

	return resp.CrosschainFlags, nil
}

// GetSupportedChains returns the supported chains
func (c *Client) GetSupportedChains(ctx context.Context) ([]chains.Chain, error) {
	resp, err := c.client.observer.SupportedChains(ctx, &types.QuerySupportedChains{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get supported chains")
	}

	return resp.GetChains(), nil
}

// GetChainParams returns all the chain params
func (c *Client) GetChainParams(ctx context.Context) ([]*types.ChainParams, error) {
	in := &types.QueryGetChainParamsRequest{}

	resp, err := retry.DoTypedWithRetry(func() (*types.QueryGetChainParamsResponse, error) {
		return c.client.observer.GetChainParams(ctx, in)
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to get chain params")
	}

	return resp.ChainParams.ChainParams, nil
}

// GetChainParamsForChainID returns the chain params for a given chain ID
func (c *Client) GetChainParamsForChainID(
	ctx context.Context,
	externalChainID int64,
) (*types.ChainParams, error) {
	in := &types.QueryGetChainParamsForChainRequest{ChainId: externalChainID}

	resp, err := c.client.observer.GetChainParamsForChain(ctx, in)
	if err != nil {
		return &types.ChainParams{}, err
	}

	return resp.ChainParams, nil
}

// GetObserverList returns the list of observers
func (c *Client) GetObserverList(ctx context.Context) ([]string, error) {
	in := &types.QueryObserverSet{}

	resp, err := retry.DoTypedWithRetry(func() (*types.QueryObserverSetResponse, error) {
		return c.client.observer.ObserverSet(ctx, in)
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to get observer list")
	}

	return resp.Observers, nil
}

// GetBallotByID returns a ballot by ID
func (c *Client) GetBallotByID(ctx context.Context, id string) (*types.QueryBallotByIdentifierResponse, error) {
	in := &types.QueryBallotByIdentifierRequest{BallotIdentifier: id}

	return c.client.observer.BallotByIdentifier(ctx, in)
}

// GetNonceByChain returns the nonce by chain
func (c *Client) GetNonceByChain(ctx context.Context, chain chains.Chain) (types.ChainNonces, error) {
	in := &types.QueryGetChainNoncesRequest{Index: chain.ChainName.String()}

	resp, err := c.client.observer.ChainNonces(ctx, in)
	if err != nil {
		return types.ChainNonces{}, errors.Wrap(err, "failed to get nonce by chain")
	}

	return resp.ChainNonces, nil
}

// GetKeyGen returns the keygen
func (c *Client) GetKeyGen(ctx context.Context) (*types.Keygen, error) {
	in := &types.QueryGetKeygenRequest{}

	resp, err := retry.DoTypedWithRetry(func() (*types.QueryGetKeygenResponse, error) {
		return c.client.observer.Keygen(ctx, in)
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to get keygen")
	}

	return resp.GetKeygen(), nil
}

// GetAllNodeAccounts returns all node accounts
func (c *Client) GetAllNodeAccounts(ctx context.Context) ([]*types.NodeAccount, error) {
	resp, err := c.client.observer.NodeAccountAll(ctx, &types.QueryAllNodeAccountRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all node accounts")
	}

	c.logger.Debug().Int("node_account.len", len(resp.NodeAccount)).Msg("GetAllNodeAccounts: OK")

	return resp.NodeAccount, nil
}

// GetBallot returns a ballot by ID
func (c *Client) GetBallot(
	ctx context.Context,
	ballotIdentifier string,
) (*types.QueryBallotByIdentifierResponse, error) {
	in := &types.QueryBallotByIdentifierRequest{BallotIdentifier: ballotIdentifier}

	resp, err := c.client.observer.BallotByIdentifier(ctx, in)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ballot")
	}

	return resp, nil
}

// GetCurrentTSS returns the current TSS
func (c *Client) GetCurrentTSS(ctx context.Context) (types.TSS, error) {
	resp, err := c.client.observer.TSS(ctx, &types.QueryGetTSSRequest{})
	if err != nil {
		return types.TSS{}, errors.Wrap(err, "failed to get current tss")
	}

	return resp.TSS, nil
}

// GetEVMTSSAddress returns the EVM TSS address.
func (c *Client) GetEVMTSSAddress(ctx context.Context) (string, error) {
	resp, err := c.client.observer.GetTssAddress(ctx, &types.QueryGetTssAddressRequest{})
	if err != nil {
		return "", errors.Wrap(err, "failed to get eth tss address")
	}

	return resp.Eth, nil
}

// GetBTCTSSAddress returns the BTC TSS address
func (c *Client) GetBTCTSSAddress(ctx context.Context, chainID int64) (string, error) {
	in := &types.QueryGetTssAddressRequest{BitcoinChainId: chainID}

	resp, err := c.client.observer.GetTssAddress(ctx, in)
	if err != nil {
		return "", errors.Wrap(err, "failed to get btc tss address")
	}
	return resp.Btc, nil
}

// GetTSSHistory returns the TSS history
func (c *Client) GetTSSHistory(ctx context.Context) ([]types.TSS, error) {
	resp, err := c.client.observer.TssHistory(ctx, &types.QueryTssHistoryRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tss history")
	}

	return resp.TssList, nil
}

// GetPendingNonces returns the pending nonces
func (c *Client) GetPendingNonces(ctx context.Context) (*types.QueryAllPendingNoncesResponse, error) {
	resp, err := c.client.observer.PendingNoncesAll(ctx, &types.QueryAllPendingNoncesRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pending nonces")
	}

	return resp, nil
}

// GetPendingNoncesByChain returns the pending nonces for a chain and current tss address
func (c *Client) GetPendingNoncesByChain(ctx context.Context, chainID int64) (types.PendingNonces, error) {
	in := &types.QueryPendingNoncesByChainRequest{ChainId: chainID}

	resp, err := c.client.observer.PendingNoncesByChain(ctx, in)
	if err != nil {
		return types.PendingNonces{}, errors.Wrap(err, "failed to get pending nonces by chain")
	}

	return resp.PendingNonces, nil
}

// HasVoted returns whether an observer has voted
func (c *Client) HasVoted(ctx context.Context, ballotIndex string, voterAddress string) (bool, error) {
	in := &types.QueryHasVotedRequest{
		BallotIdentifier: ballotIndex,
		VoterAddress:     voterAddress,
	}

	resp, err := c.client.observer.HasVoted(ctx, in)
	if err != nil {
		return false, errors.Wrap(err, "failed to check if observer has voted")
	}

	return resp.HasVoted, nil
}
