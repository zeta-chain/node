package clients

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	zetaevmclient "github.com/zeta-chain/node/zetaclient/chains/evm/client"
)

// EVMClientAdapter wraps ethclient.Client to implement EVMClient interface
type EVMClientAdapter struct {
	client     *ethclient.Client
	zetaClient *zetaevmclient.Client
}

// NewEVMClientAdapter creates a new EVMClientAdapter
func NewEVMClientAdapter(rpcURL string) (*EVMClientAdapter, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to EVM RPC: %w", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	zetaClient := zetaevmclient.New(client, ethtypes.NewLondonSigner(chainID))

	return &EVMClientAdapter{
		client:     client,
		zetaClient: zetaClient,
	}, nil
}

// TransactionByHash retrieves a transaction by its hash
func (e *EVMClientAdapter) TransactionByHash(ctx context.Context, hash string) (*ethtypes.Transaction, bool, error) {
	txHash := common.HexToHash(hash)
	return e.client.TransactionByHash(ctx, txHash)
}

// TransactionReceipt retrieves a transaction receipt by transaction hash
func (e *EVMClientAdapter) TransactionReceipt(ctx context.Context, hash string) (*ethtypes.Receipt, error) {
	txHash := common.HexToHash(hash)
	return e.client.TransactionReceipt(ctx, txHash)
}

// IsTxConfirmed checks if a transaction has enough confirmations
func (e *EVMClientAdapter) IsTxConfirmed(ctx context.Context, txHash string, confirmations uint64) (bool, error) {
	return e.zetaClient.IsTxConfirmed(ctx, txHash, confirmations)
}

// ChainID returns the chain ID
func (e *EVMClientAdapter) ChainID(ctx context.Context) (*big.Int, error) {
	return e.client.ChainID(ctx)
}

// GetRawClient returns the underlying ethclient.Client for advanced operations
func (e *EVMClientAdapter) GetRawClient() *ethclient.Client {
	return e.client
}
