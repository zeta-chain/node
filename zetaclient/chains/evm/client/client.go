package client

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/zetaclient/metrics"
)

type Client struct {
	*ethclient.Client

	signer eth.Signer
}

// NewFromEndpoint new Client constructor based on endpoint URL.
func NewFromEndpoint(ctx context.Context, endpoint string) (*Client, error) {
	httpClient, err := metrics.GetInstrumentedHTTPClient(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get instrumented HTTP client")
	}

	ethRPC, err := rpc.DialOptions(ctx, endpoint, rpc.WithHTTPClient(httpClient))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to dial EVM client (endpoint %q)", endpoint)
	}

	client := ethclient.NewClient(ethRPC)

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get chain ID")
	}

	ethSigner := eth.LatestSignerForChainID(chainID)

	return New(client, ethSigner), nil
}

// New Client constructor.
func New(client *ethclient.Client, signer eth.Signer) *Client {
	return &Client{client, signer}
}

// IsTxConfirmed checks whether txHash settled on-chain && has at least X blocks of confirmations.
func (c *Client) IsTxConfirmed(ctx context.Context, txHash string, confirmations uint64) (bool, error) {
	if confirmations == 0 {
		return false, errors.New("confirmations must be greater than 0")
	}

	hash := common.HexToHash(txHash)

	// query the tx
	_, isPending, err := c.TransactionByHash(ctx, hash)
	switch {
	case err != nil:
		return false, errors.Wrapf(err, "error getting transaction for tx %s", txHash)
	case isPending:
		return false, nil
	}

	// query receipt
	receipt, err := c.TransactionReceipt(ctx, hash)
	switch {
	case err != nil:
		return false, errors.Wrapf(err, "error getting transaction receipt for tx %s", txHash)
	case receipt == nil:
		// should not happen
		return false, errors.Errorf("receipt is nil for tx %s", txHash)
	}

	// query last block height
	blockNumber, err := c.BlockNumber(ctx)
	switch {
	case err != nil:
		return false, errors.Wrap(err, "error getting block number")
	case blockNumber < receipt.BlockNumber.Uint64():
		// should not happen
		return false, nil
	}

	// check confirmations
	blocksConfirmed := 1 + (blockNumber - receipt.BlockNumber.Uint64())

	return blocksConfirmed >= confirmations, nil
}

// HealthCheck returns the latest block time in UTC.
func (c *Client) HealthCheck(ctx context.Context) (time.Time, error) {
	// query latest block header
	header, err := c.HeaderByNumber(ctx, nil)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "unable to get latest block header")
	}

	// convert block time to UTC
	// #nosec G115 always in range
	blockTime := time.Unix(int64(header.Time), 0).UTC()

	return blockTime, nil
}

// Overrides c.Client ChainID method.
func (c *Client) Signer() eth.Signer {
	return c.signer
}
