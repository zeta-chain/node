package rpc

import (
	"context"
	"fmt"

	"cosmossdk.io/errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/retry"
	"github.com/zeta-chain/node/x/observer/types"
)

// GetCrosschainFlags returns the crosschain flags
func (c *Clients) GetCrosschainFlags(ctx context.Context) (types.CrosschainFlags, error) {
	resp, err := c.Observer.CrosschainFlags(ctx, &types.QueryGetCrosschainFlagsRequest{})
	if err != nil {
		return types.CrosschainFlags{}, err
	}

	return resp.CrosschainFlags, nil
}

// GetSupportedChains returns the supported chains
func (c *Clients) GetSupportedChains(ctx context.Context) ([]chains.Chain, error) {
	resp, err := c.Observer.SupportedChains(ctx, &types.QuerySupportedChains{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get supported chains")
	}

	return resp.GetChains(), nil
}

// GetChainParams returns all the chain params
func (c *Clients) GetChainParams(ctx context.Context) ([]*types.ChainParams, error) {
	in := &types.QueryGetChainParamsRequest{}

	resp, err := retry.DoTypedWithRetry(func() (*types.QueryGetChainParamsResponse, error) {
		return c.Observer.GetChainParams(ctx, in)
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to get chain params")
	}

	return resp.ChainParams.ChainParams, nil
}

// GetChainParamsForChainID returns the chain params for a given chain ID
func (c *Clients) GetChainParamsForChainID(
	ctx context.Context,
	externalChainID int64,
) (*types.ChainParams, error) {
	in := &types.QueryGetChainParamsForChainRequest{ChainId: externalChainID}

	resp, err := c.Observer.GetChainParamsForChain(ctx, in)
	if err != nil {
		return &types.ChainParams{}, err
	}

	return resp.ChainParams, nil
}

// GetObserverList returns the list of observers
func (c *Clients) GetObserverList(ctx context.Context) ([]string, error) {
	in := &types.QueryObserverSet{}

	resp, err := retry.DoTypedWithRetry(func() (*types.QueryObserverSetResponse, error) {
		return c.Observer.ObserverSet(ctx, in)
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to get observer list")
	}

	return resp.Observers, nil
}

// GetBallotByID returns a ballot by ID
func (c *Clients) GetBallotByID(ctx context.Context, id string) (*types.QueryBallotByIdentifierResponse, error) {
	in := &types.QueryBallotByIdentifierRequest{BallotIdentifier: id}

	return c.Observer.BallotByIdentifier(ctx, in)
}

// GetNonceByChain returns the nonce by chain
func (c *Clients) GetNonceByChain(ctx context.Context, chain chains.Chain) (types.ChainNonces, error) {
	in := &types.QueryGetChainNoncesRequest{ChainId: chain.ChainId}

	resp, err := c.Observer.ChainNonces(ctx, in)
	if err != nil {
		return types.ChainNonces{}, errors.Wrap(err, "failed to get nonce by chain")
	}

	return resp.ChainNonces, nil
}

// GetKeyGen returns the keygen
func (c *Clients) GetKeyGen(ctx context.Context) (types.Keygen, error) {
	in := &types.QueryGetKeygenRequest{}

	resp, err := retry.DoTypedWithRetry(func() (*types.QueryGetKeygenResponse, error) {
		return c.Observer.Keygen(ctx, in)
	})

	switch {
	case err != nil:
		return types.Keygen{}, errors.Wrap(err, "failed to get keygen")
	case resp.Keygen == nil:
		return types.Keygen{}, fmt.Errorf("keygen is nil")
	}

	return *resp.Keygen, nil
}

// GetAllNodeAccounts returns all node accounts
func (c *Clients) GetAllNodeAccounts(ctx context.Context) ([]*types.NodeAccount, error) {
	resp, err := c.Observer.NodeAccountAll(ctx, &types.QueryAllNodeAccountRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all node accounts")
	}

	return resp.NodeAccount, nil
}

// GetBallot returns a ballot by ID
func (c *Clients) GetBallot(
	ctx context.Context,
	ballotIdentifier string,
) (*types.QueryBallotByIdentifierResponse, error) {
	in := &types.QueryBallotByIdentifierRequest{BallotIdentifier: ballotIdentifier}

	resp, err := c.Observer.BallotByIdentifier(ctx, in)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get ballot")
	}

	return resp, nil
}

// GetEVMTSSAddress returns the current EVM TSS address.
func (c *Clients) GetEVMTSSAddress(ctx context.Context) (string, error) {
	resp, err := c.Observer.GetTssAddress(ctx, &types.QueryGetTssAddressRequest{})
	if err != nil {
		return "", errors.Wrap(err, "failed to get eth tss address")
	}

	return resp.Eth, nil
}

// GetBTCTSSAddress returns the current BTC TSS address
func (c *Clients) GetBTCTSSAddress(ctx context.Context, chainID int64) (string, error) {
	in := &types.QueryGetTssAddressRequest{BitcoinChainId: chainID}

	resp, err := c.Observer.GetTssAddress(ctx, in)
	if err != nil {
		return "", errors.Wrap(err, "failed to get btc tss address")
	}
	return resp.Btc, nil
}

// GetTSS returns the current TSS
func (c *Clients) GetTSS(ctx context.Context) (types.TSS, error) {
	resp, err := c.Observer.TSS(ctx, &types.QueryGetTSSRequest{})
	if err != nil {
		return types.TSS{}, errors.Wrap(err, "failed to get tss")
	}
	return resp.TSS, nil
}

// GetTSSHistory returns the historical list of TSS
func (c *Clients) GetTSSHistory(ctx context.Context) ([]types.TSS, error) {
	resp, err := c.Observer.TssHistory(ctx, &types.QueryTssHistoryRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get tss history")
	}

	return resp.TssList, nil
}

// GetPendingNonces returns the pending nonces
func (c *Clients) GetPendingNonces(ctx context.Context) (*types.QueryAllPendingNoncesResponse, error) {
	resp, err := c.Observer.PendingNoncesAll(ctx, &types.QueryAllPendingNoncesRequest{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pending nonces")
	}

	return resp, nil
}

// GetPendingNoncesByChain returns the pending nonces for a chain and current tss address
func (c *Clients) GetPendingNoncesByChain(ctx context.Context, chainID int64) (types.PendingNonces, error) {
	in := &types.QueryPendingNoncesByChainRequest{ChainId: chainID}

	resp, err := c.Observer.PendingNoncesByChain(ctx, in)
	if err != nil {
		return types.PendingNonces{}, errors.Wrap(err, "failed to get pending nonces by chain")
	}

	return resp.PendingNonces, nil
}

// HasVoted returns whether an observer has voted
func (c *Clients) HasVoted(ctx context.Context, ballotIndex string, voterAddress string) (bool, error) {
	in := &types.QueryHasVotedRequest{
		BallotIdentifier: ballotIndex,
		VoterAddress:     voterAddress,
	}

	resp, err := c.Observer.HasVoted(ctx, in)
	if err != nil {
		return false, errors.Wrap(err, "failed to check if observer has voted")
	}

	return resp.HasVoted, nil
}
