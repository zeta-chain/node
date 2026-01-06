// Package observer implements the EVM chain observer
package observer

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayevm.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/zetaconnector.non-eth.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/zetaconnectornative.sol"

	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/evm/client"
	"github.com/zeta-chain/node/zetaclient/chains/evm/common"
	"github.com/zeta-chain/node/zetaclient/metrics"
)

// EVMClient is the interface for the EVM RPC client.
//
// This interface contains functions from go-ethereum's bind.ContractBackend.
// We are intentionally not embedding ContractBackend in EVMClient because we want to make it
// explicit and clear which methods are part of EVMClient.
//
//go:generate mockery --name EVMClient --filename evm_client.go --case underscore --output ../../../testutils/mocks
type EVMClient interface {
	PendingCodeAt(context.Context, ethcommon.Address) ([]byte, error)
	PendingNonceAt(context.Context, ethcommon.Address) (uint64, error)
	SubscribeFilterLogs(context.Context, ethereum.FilterQuery, chan<- eth.Log) (ethereum.Subscription, error)
	EstimateGas(context.Context, ethereum.CallMsg) (uint64, error)
	CodeAt(_ context.Context, contract ethcommon.Address, blockNumber *big.Int) ([]byte, error)
	CallContract(_ context.Context, _ ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	// This is a mutating function that does not get called when zetaclient is in dry-mode.
	SendTransaction(context.Context, *eth.Transaction) error

	HealthCheck(ctx context.Context) (time.Time, error)

	FilterLogs(context.Context, ethereum.FilterQuery) ([]eth.Log, error)

	SuggestGasTipCap(ctx context.Context) (*big.Int, error)

	SuggestGasPrice(context.Context) (*big.Int, error)

	BlockNumber(context.Context) (uint64, error)

	BlockByNumberCustom(context.Context, *big.Int) (*client.Block, error)

	HeaderByNumber(context.Context, *big.Int) (*eth.Header, error)

	TransactionByHash(context.Context, ethcommon.Hash) (_ *eth.Transaction, isPending bool, _ error)

	TransactionByHashCustom(context.Context, string) (*client.Transaction, error)

	TransactionReceipt(context.Context, ethcommon.Hash) (*eth.Receipt, error)
}

// Observer is the observer for evm chains
type Observer struct {
	// base.Observer implements the base chain observer
	*base.Observer

	priorityFeeConfig

	// evmClient is the EVM client for the observed chain
	evmClient EVMClient

	// outboundConfirmedReceipts is the map to index confirmed receipts by hash
	outboundConfirmedReceipts map[string]*eth.Receipt

	// outboundConfirmedTransactions is the map to index confirmed transactions by hash
	outboundConfirmedTransactions map[string]*eth.Transaction
}

// priorityFeeConfig is the configuration for priority fee
type priorityFeeConfig struct {
	// checked indicates whether the observer checked
	// this EVM chain for EIP-1559 (further checks are cached)
	checked bool

	// supported indicates whether this EVM chain supports EIP-1559
	supported bool
}

// New Observer constructor
func New(baseObserver *base.Observer, evmClient EVMClient) (*Observer, error) {
	// create evm observer
	ob := &Observer{
		Observer:                      baseObserver,
		evmClient:                     evmClient,
		outboundConfirmedReceipts:     make(map[string]*eth.Receipt),
		outboundConfirmedTransactions: make(map[string]*eth.Transaction),
		priorityFeeConfig:             priorityFeeConfig{},
	}

	// load last block scanned
	if err := ob.loadLastBlockScanned(context.Background()); err != nil {
		return nil, errors.Wrap(err, "unable to load last block scanned")
	}

	return ob, nil
}

// getConnectorLegacyContract returns the non-Eth connector address and binder
func (ob *Observer) getConnectorLegacyContract() (ethcommon.Address, *zetaconnector.ZetaConnectorNonEth, error) {
	addr := ethcommon.HexToAddress(ob.ChainParams().ConnectorContractAddress)
	contract, err := zetaconnector.NewZetaConnectorNonEth(addr, ob.evmClient)
	return addr, contract, err
}

func (ob *Observer) getConnectorContract() (ethcommon.Address, *zetaconnectornative.ZetaConnectorNative, error) {
	addr := ethcommon.HexToAddress(ob.ChainParams().ConnectorContractAddress)
	contract, err := zetaconnectornative.NewZetaConnectorNative(addr, ob.evmClient)
	return addr, contract, err
}

// getERC20CustodyContract returns ERC20Custody contract address and binder
func (ob *Observer) getERC20CustodyContract() (ethcommon.Address, *erc20custody.ERC20Custody, error) {
	addr := ethcommon.HexToAddress(ob.ChainParams().Erc20CustodyContractAddress)
	contract, err := erc20custody.NewERC20Custody(addr, ob.evmClient)
	return addr, contract, err
}

// getERC20CustodyV2Contract returns ERC20Custody contract address and binder
// NOTE: we use the same address as gateway v1
// this simplify the migration process v1 will be completely removed in the future
// currently the ABI for withdraw is identical, therefore both contract instances can be used
func (ob *Observer) getERC20CustodyV2Contract() (ethcommon.Address, *erc20custody.ERC20Custody, error) {
	addr := ethcommon.HexToAddress(ob.ChainParams().Erc20CustodyContractAddress)
	contract, err := erc20custody.NewERC20Custody(addr, ob.evmClient)
	return addr, contract, err
}

// getGatewayContract returns the gateway contract address and binder
func (ob *Observer) getGatewayContract() (ethcommon.Address, *gatewayevm.GatewayEVM, error) {
	addr := ethcommon.HexToAddress(ob.ChainParams().GatewayAddress)
	contract, err := gatewayevm.NewGatewayEVM(addr, ob.evmClient)
	return addr, contract, err
}

// setTxNReceipt sets the receipt and transaction in memory
func (ob *Observer) setTxNReceipt(nonce uint64, receipt *eth.Receipt, transaction *eth.Transaction) {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	ob.outboundConfirmedReceipts[ob.OutboundID(nonce)] = receipt
	ob.outboundConfirmedTransactions[ob.OutboundID(nonce)] = transaction
}

// getTxNReceipt gets the receipt and transaction from memory
func (ob *Observer) getTxNReceipt(nonce uint64) (*eth.Receipt, *eth.Transaction) {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	receipt := ob.outboundConfirmedReceipts[ob.OutboundID(nonce)]
	transaction := ob.outboundConfirmedTransactions[ob.OutboundID(nonce)]
	return receipt, transaction
}

// isTxConfirmed returns true if there is a confirmed tx for 'nonce'
func (ob *Observer) isTxConfirmed(nonce uint64) bool {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	id := ob.OutboundID(nonce)

	return ob.outboundConfirmedReceipts[id] != nil && ob.outboundConfirmedTransactions[id] != nil
}

// checkTxInclusion returns nil only if tx is included at the position indicated by the receipt ([block, index])
func (ob *Observer) checkTxInclusion(ctx context.Context, tx *eth.Transaction, receipt *eth.Receipt) error {
	block, err := ob.GetBlockByNumberCached(ctx, receipt.BlockNumber.Uint64())
	if err != nil {
		return errors.Wrapf(err, "GetBlockByNumberCached error for block %d txHash %s nonce %d",
			receipt.BlockNumber.Uint64(), tx.Hash(), tx.Nonce())
	}

	// #nosec G115 non negative value
	if receipt.TransactionIndex >= uint(len(block.Transactions)) {
		return fmt.Errorf("transaction index %d out of range [0, %d), txHash %s nonce %d block %d",
			receipt.TransactionIndex, len(block.Transactions), tx.Hash(), tx.Nonce(), receipt.BlockNumber.Uint64())
	}

	txAtIndex := block.Transactions[receipt.TransactionIndex]
	if !strings.EqualFold(txAtIndex.Hash, tx.Hash().Hex()) {
		ob.removeCachedBlock(receipt.BlockNumber.Uint64()) // clean stale block from cache
		return fmt.Errorf("transaction at index %d has different hash %s, txHash %s nonce %d block %d",
			receipt.TransactionIndex, txAtIndex.Hash, tx.Hash(), tx.Nonce(), receipt.BlockNumber.Uint64())
	}

	return nil
}

// transactionByHash query transaction by hash via JSON-RPC
func (ob *Observer) transactionByHash(ctx context.Context, txHash string) (*client.Transaction, bool, error) {
	tx, err := ob.evmClient.TransactionByHashCustom(ctx, txHash)
	if err != nil {
		return nil, false, err
	}
	err = common.ValidateEvmTransaction(tx)
	if err != nil {
		return nil, false, err
	}
	return tx, tx.BlockNumber == nil, nil
}

// GetBlockByNumberCached get block by number from cache
// returns block, ethrpc.Block, isFallback, isSkip, error
func (ob *Observer) GetBlockByNumberCached(ctx context.Context, blockNumber uint64) (*client.Block, error) {
	if result, ok := ob.BlockCache().Get(blockNumber); ok {
		if block, ok := result.(*client.Block); ok {
			return block, nil
		}
		return nil, errors.New("cached value is not of type *client.Block")
	}
	if blockNumber > math.MaxInt32 {
		return nil, fmt.Errorf("block number %d is too large", blockNumber)
	}
	// #nosec G115 always in range, checked above
	block, err := ob.blockByNumber(ctx, int(blockNumber))
	if err != nil {
		return nil, err
	}
	ob.BlockCache().Add(blockNumber, block)
	return block, nil
}

// removeCachedBlock remove block from cache
func (ob *Observer) removeCachedBlock(blockNumber uint64) {
	ob.BlockCache().Remove(blockNumber)
}

// blockByNumber query block by number via JSON-RPC
func (ob *Observer) blockByNumber(ctx context.Context, blockNumber int) (*client.Block, error) {
	block, err := ob.evmClient.BlockByNumberCustom(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		return nil, err
	}
	if block == nil {
		return nil, fmt.Errorf("block not found: %d", blockNumber)
	}
	for i := range block.Transactions {
		err := common.ValidateEvmTransaction(&block.Transactions[i])
		if err != nil {
			return nil, err
		}
	}
	return block, nil
}

// loadLastBlockScanned loads the last scanned block from the database
// TODO(revamp): move to a db file
func (ob *Observer) loadLastBlockScanned(ctx context.Context) error {
	err := ob.Observer.LoadLastBlockScanned()
	if err != nil {
		return errors.Wrapf(err, "error LoadLastBlockScanned for chain %d", ob.Chain().ChainId)
	}

	// observer will scan from the last block when 'lastBlockScanned == 0', this happens when:
	// 1. environment variable is set explicitly to "latest"
	// 2. environment variable is empty and last scanned block is not found in DB
	if ob.LastBlockScanned() == 0 {
		blockNumber, err := ob.evmClient.BlockNumber(ctx)
		if err != nil {
			return errors.Wrapf(err, "error BlockNumber for chain %d", ob.Chain().ChainId)
		}
		ob.WithLastBlockScanned(blockNumber)
	}
	ob.Logger().Chain.Info().
		Uint64("last_block_scanned", ob.LastBlockScanned()).
		Send()

	return nil
}

func (ob *Observer) CheckRPCStatus(ctx context.Context) error {
	blockTime, err := ob.evmClient.HealthCheck(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to check rpc health")
	}

	metrics.ReportBlockLatency(ob.Chain().Name, blockTime)

	return nil
}

// updateLastBlock is a helper function to update the last block number.
// Note: keep last block up-to-date helps to avoid inaccurate confirmation.
func (ob *Observer) updateLastBlock(ctx context.Context) error {
	blockNumber, err := ob.evmClient.BlockNumber(ctx)
	switch {
	case err != nil:
		return errors.Wrap(err, "error getting block number")
	case blockNumber < ob.LastBlock():
		return fmt.Errorf("block number should not decrease: current %d last %d", blockNumber, ob.LastBlock())
	default:
		ob.WithLastBlock(blockNumber)
	}

	// increment prom counter
	metrics.GetBlockNumberPerChain.WithLabelValues(ob.Chain().Name).Inc()

	return nil
}
