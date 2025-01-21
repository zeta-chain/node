package client

import (
	"context"
	"math/big"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/metrics"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

type Client struct {
	interfaces.EVMRPCClient
	ethtypes.Signer
}

func NewFromEndpoint(ctx context.Context, endpoint string) (*Client, error) {
	if endpoint == testutils.MockEVMRPCEndpoint {
		chainID := big.NewInt(chains.Ethereum.ChainId)
		ethSigner := ethtypes.NewLondonSigner(chainID)
		client := &mocks.EVMRPCClient{}

		return New(client, ethSigner), nil
	}

	httpClient, err := metrics.GetInstrumentedHTTPClient(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get instrumented HTTP client")
	}

	rpc, err := ethrpc.DialOptions(ctx, endpoint, ethrpc.WithHTTPClient(httpClient))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to dial EVM client (endpoint %q)", endpoint)
	}

	client := ethclient.NewClient(rpc)

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get chain ID")
	}

	ethSigner := ethtypes.LatestSignerForChainID(chainID)

	return New(client, ethSigner), nil
}

func New(client interfaces.EVMRPCClient, signer ethtypes.Signer) *Client {
	return &Client{client, signer}
}

func (c *Client) IsTxConfirmed(ctx context.Context, txHash string, confirmations uint64) (bool, error) {
	hash := ethcommon.HexToHash(txHash)

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
	lastHeight, err := c.BlockNumber(ctx)
	switch {
	case err != nil:
		return false, errors.Wrap(err, "error getting block number")
	case lastHeight < receipt.BlockNumber.Uint64():
		// check confirmations
		return false, nil
	}

	blocksConfirmed := 1 + (lastHeight - receipt.BlockNumber.Uint64())

	return blocksConfirmed >= confirmations, nil
}

// HealthCheck asserts RPC health. Returns the latest block time in UTC.
func (c *Client) HealthCheck(ctx context.Context) (time.Time, error) {
	// query latest block number
	bn, err := c.BlockNumber(ctx)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "RPC failed on BlockNumber, RPC down?")
	}

	// query suggested gas price
	if _, err = c.EVMRPCClient.SuggestGasPrice(ctx); err != nil {
		return time.Time{}, errors.Wrap(err, "RPC failed on SuggestGasPrice, RPC down?")
	}

	// query latest block header
	header, err := c.EVMRPCClient.HeaderByNumber(ctx, new(big.Int).SetUint64(bn))
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "RPC failed on HeaderByNumber(%d), RPC down?", bn)
	}

	// convert block time to UTC
	// #nosec G115 always in range
	blockTime := time.Unix(int64(header.Time), 0).UTC()

	return blockTime, nil
}
