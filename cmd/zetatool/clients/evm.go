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
	"github.com/zeta-chain/protocol-contracts-evm/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayevm.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/zetaconnector.non-eth.sol"

	"github.com/zeta-chain/node/cmd/zetatool/config"
	"github.com/zeta-chain/node/pkg/chains"
)

// EVMClientAdapter wraps an ethclient.Client to implement the EVMClient interface
type EVMClientAdapter struct {
	client *ethclient.Client
}

// NewEVMClientAdapter creates a new EVMClientAdapter for a specific chain using config
func NewEVMClientAdapter(chain chains.Chain, cfg *config.Config) (*EVMClientAdapter, error) {
	rpcURL := ResolveEVMRPC(chain, cfg)
	if rpcURL == "" {
		return nil, fmt.Errorf("rpc not found for chain %d network %s", chain.ChainId, chain.Network)
	}
	rpcClient, err := ethrpc.DialHTTP(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to eth rpc: %w", err)
	}
	return &EVMClientAdapter{client: ethclient.NewClient(rpcClient)}, nil
}

// TransactionByHash retrieves a transaction by its hash
func (a *EVMClientAdapter) TransactionByHash(ctx context.Context, hash string) (*ethtypes.Transaction, bool, error) {
	return a.client.TransactionByHash(ctx, ethcommon.HexToHash(hash))
}

// TransactionReceipt retrieves the receipt for a transaction
func (a *EVMClientAdapter) TransactionReceipt(ctx context.Context, hash string) (*ethtypes.Receipt, error) {
	return a.client.TransactionReceipt(ctx, ethcommon.HexToHash(hash))
}

// ChainID returns the chain ID
func (a *EVMClientAdapter) ChainID(ctx context.Context) (*big.Int, error) {
	return a.client.ChainID(ctx)
}

// BlockNumber returns the current block number
func (a *EVMClientAdapter) BlockNumber(ctx context.Context) (uint64, error) {
	return a.client.BlockNumber(ctx)
}

// TransactionSender returns the sender of a transaction
func (a *EVMClientAdapter) TransactionSender(
	ctx context.Context,
	tx *ethtypes.Transaction,
	blockHash ethcommon.Hash,
	txIndex uint,
) (ethcommon.Address, error) {
	return a.client.TransactionSender(ctx, tx, blockHash, txIndex)
}

// ParseConnectorZetaSent parses a ZetaSent event from the connector contract
func (a *EVMClientAdapter) ParseConnectorZetaSent(
	log ethtypes.Log,
	connectorAddr string,
) (*zetaconnector.ZetaConnectorNonEthZetaSent, error) {
	connector, err := zetaconnector.NewZetaConnectorNonEth(ethcommon.HexToAddress(connectorAddr), a.client)
	if err != nil {
		return nil, fmt.Errorf("failed to create connector contract: %w", err)
	}
	return connector.ParseZetaSent(log)
}

// ParseCustodyDeposited parses a Deposited event from the custody contract
func (a *EVMClientAdapter) ParseCustodyDeposited(
	log ethtypes.Log,
	custodyAddr string,
) (*erc20custody.ERC20CustodyDeposited, error) {
	custody, err := erc20custody.NewERC20Custody(ethcommon.HexToAddress(custodyAddr), a.client)
	if err != nil {
		return nil, fmt.Errorf("failed to create custody contract: %w", err)
	}
	return custody.ParseDeposited(log)
}

// ParseGatewayDeposited parses a Deposited event from the gateway contract
func (a *EVMClientAdapter) ParseGatewayDeposited(
	log ethtypes.Log,
	gatewayAddr string,
) (*gatewayevm.GatewayEVMDeposited, error) {
	gateway, err := gatewayevm.NewGatewayEVM(ethcommon.HexToAddress(gatewayAddr), a.client)
	if err != nil {
		return nil, fmt.Errorf("failed to create gateway contract: %w", err)
	}
	return gateway.ParseDeposited(log)
}

// ParseGatewayDepositedAndCalled parses a DepositedAndCalled event from the gateway contract
func (a *EVMClientAdapter) ParseGatewayDepositedAndCalled(
	log ethtypes.Log,
	gatewayAddr string,
) (*gatewayevm.GatewayEVMDepositedAndCalled, error) {
	gateway, err := gatewayevm.NewGatewayEVM(ethcommon.HexToAddress(gatewayAddr), a.client)
	if err != nil {
		return nil, fmt.Errorf("failed to create gateway contract: %w", err)
	}
	return gateway.ParseDepositedAndCalled(log)
}

// ParseGatewayCalled parses a Called event from the gateway contract
func (a *EVMClientAdapter) ParseGatewayCalled(
	log ethtypes.Log,
	gatewayAddr string,
) (*gatewayevm.GatewayEVMCalled, error) {
	gateway, err := gatewayevm.NewGatewayEVM(ethcommon.HexToAddress(gatewayAddr), a.client)
	if err != nil {
		return nil, fmt.Errorf("failed to create gateway contract: %w", err)
	}
	return gateway.ParseCalled(log)
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

// GetEvmTx retrieves a transaction and its receipt by hash using the EVMClient interface
func GetEvmTx(
	ctx context.Context,
	client EVMClient,
	txHash string,
	chainID int64,
) (*ethtypes.Transaction, *ethtypes.Receipt, error) {
	tx, isPending, err := client.TransactionByHash(ctx, txHash)
	if err != nil {
		return nil, nil, fmt.Errorf("tx not found on chain: %w, chainID: %d", err, chainID)
	}
	if isPending {
		return nil, nil, fmt.Errorf("tx is still pending on chain: %d", chainID)
	}
	receipt, err := client.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get receipt: %w, tx hash: %s", err, txHash)
	}
	return tx, receipt, nil
}

// IsTxConfirmed checks if a transaction has enough confirmations using the EVMClient interface
func IsTxConfirmed(ctx context.Context, client EVMClient, txHash string, confirmations uint64) (bool, error) {
	if confirmations == 0 {
		return false, fmt.Errorf("confirmations must be greater than 0")
	}

	_, isPending, err := client.TransactionByHash(ctx, txHash)
	if err != nil {
		return false, fmt.Errorf("error getting transaction: %w", err)
	}
	if isPending {
		return false, nil
	}

	receipt, err := client.TransactionReceipt(ctx, txHash)
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
	return 1+(currentBlock-txBlock) >= confirmations, nil
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
