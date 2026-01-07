package clients

import (
	"context"
	"fmt"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	ethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/zeta-chain/node/cmd/zetatool/config"
	"github.com/zeta-chain/node/pkg/chains"
)

// EVMClientAdapter wraps ethclient.Client to implement EVMClient interface
type EVMClientAdapter struct {
	client *ethclient.Client
}

// NewEVMClientAdapter creates a new EVMClientAdapter from an RPC URL
func NewEVMClientAdapter(rpcURL string) (*EVMClientAdapter, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to EVM RPC: %w", err)
	}
	return &EVMClientAdapter{client: client}, nil
}

// NewEVMClientForChain creates an EVM client for a specific chain using config
func NewEVMClientForChain(chain chains.Chain, cfg *config.Config) (*ethclient.Client, error) {
	rpcURL := ResolveEVMRPC(chain, cfg)
	if rpcURL == "" {
		return nil, fmt.Errorf("rpc not found for chain %d network %s", chain.ChainId, chain.Network)
	}
	rpcClient, err := ethrpc.DialHTTP(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to eth rpc: %w", err)
	}
	return ethclient.NewClient(rpcClient), nil
}

// ResolveEVMRPC returns the RPC URL for a given EVM chain
func ResolveEVMRPC(chain chains.Chain, cfg *config.Config) string {
	return map[chains.Network]string{
		chains.Network_eth:        cfg.EthereumRPC,
		chains.Network_base:       cfg.BaseRPC,
		chains.Network_polygon:    cfg.PolygonRPC,
		chains.Network_bsc:        cfg.BscRPC,
		chains.Network_arbitrum:   cfg.ArbitrumRPC,
		chains.Network_optimism:   cfg.OptimismRPC,
		chains.Network_avalanche:  cfg.AvalancheRPC,
		chains.Network_worldchain: cfg.WorldRPC,
	}[chain.Network]
}

// GetEvmTx retrieves a transaction and its receipt by hash
func GetEvmTx(
	ctx context.Context,
	client *ethclient.Client,
	txHash string,
	chainID int64,
) (*ethtypes.Transaction, *ethtypes.Receipt, error) {
	hash := ethcommon.HexToHash(txHash)
	tx, isPending, err := client.TransactionByHash(ctx, hash)
	if err != nil {
		return nil, nil, fmt.Errorf("tx not found on chain: %w, chainID: %d", err, chainID)
	}
	if isPending {
		return nil, nil, fmt.Errorf("tx is still pending on chain: %d", chainID)
	}
	receipt, err := client.TransactionReceipt(ctx, hash)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get receipt: %w, tx hash: %s", err, txHash)
	}
	return tx, receipt, nil
}

// IsTxConfirmed checks if a transaction has enough confirmations
func IsTxConfirmed(ctx context.Context, client *ethclient.Client, txHash string, confirmations uint64) (bool, error) {
	if confirmations == 0 {
		return false, fmt.Errorf("confirmations must be greater than 0")
	}

	hash := ethcommon.HexToHash(txHash)

	_, isPending, err := client.TransactionByHash(ctx, hash)
	if err != nil {
		return false, fmt.Errorf("error getting transaction: %w", err)
	}
	if isPending {
		return false, nil
	}

	receipt, err := client.TransactionReceipt(ctx, hash)
	if err != nil {
		return false, fmt.Errorf("error getting receipt: %w", err)
	}
	if receipt == nil {
		return false, fmt.Errorf("receipt is nil for tx %s", txHash)
	}

	currentBlock, err := client.BlockNumber(ctx)
	if err != nil {
		return false, fmt.Errorf("error getting block number: %w", err)
	}

	txBlock := receipt.BlockNumber.Uint64()
	return currentBlock >= txBlock+confirmations, nil
}

// GetEVMBalance fetches the native token balance for an address on an EVM chain
func GetEVMBalance(ctx context.Context, rpcURL string, address ethcommon.Address) (*big.Int, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	return client.BalanceAt(ctx, address, nil)
}

// FormatEVMBalance converts wei to ETH with 9 decimal places
func FormatEVMBalance(wei *big.Int) string {
	if wei == nil {
		return "0.000000000"
	}

	weiFloat := new(big.Float).SetInt(wei)
	divisor := new(big.Float).SetInt(big.NewInt(params.Ether))
	eth := new(big.Float).Quo(weiFloat, divisor)

	return eth.Text('f', 9)
}

// TransactionByHash retrieves a transaction by its hash (EVMClientAdapter method)
func (e *EVMClientAdapter) TransactionByHash(ctx context.Context, hash string) (*ethtypes.Transaction, bool, error) {
	txHash := ethcommon.HexToHash(hash)
	return e.client.TransactionByHash(ctx, txHash)
}

// TransactionReceipt retrieves a transaction receipt by transaction hash (EVMClientAdapter method)
func (e *EVMClientAdapter) TransactionReceipt(ctx context.Context, hash string) (*ethtypes.Receipt, error) {
	txHash := ethcommon.HexToHash(hash)
	return e.client.TransactionReceipt(ctx, txHash)
}

// ChainID returns the chain ID (EVMClientAdapter method)
func (e *EVMClientAdapter) ChainID(ctx context.Context) (*big.Int, error) {
	return e.client.ChainID(ctx)
}

// GetRawClient returns the underlying ethclient.Client for advanced operations
func (e *EVMClientAdapter) GetRawClient() *ethclient.Client {
	return e.client
}
