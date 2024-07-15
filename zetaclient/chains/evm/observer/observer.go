// Package observer implements the EVM chain observer
package observer

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/onrik/ethrpc"
	"github.com/pkg/errors"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zeta.non-eth.sol"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.non-eth.sol"

	"github.com/zeta-chain/zetacore/pkg/bg"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/base"
	"github.com/zeta-chain/zetacore/zetaclient/chains/evm"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
)

var _ interfaces.ChainObserver = &Observer{}

// Observer is the observer for evm chains
type Observer struct {
	// base.Observer implements the base chain observer
	base.Observer

	// evmClient is the EVM client for the observed chain
	evmClient interfaces.EVMRPCClient

	// evmJSONRPC is the EVM JSON RPC client for the observed chain
	evmJSONRPC interfaces.EVMJSONRPCClient

	// outboundPendingTransactions is the map to index pending transactions by hash
	outboundPendingTransactions map[string]*ethtypes.Transaction

	// outboundConfirmedReceipts is the map to index confirmed receipts by hash
	outboundConfirmedReceipts map[string]*ethtypes.Receipt

	// outboundConfirmedTransactions is the map to index confirmed transactions by hash
	outboundConfirmedTransactions map[string]*ethtypes.Transaction
}

// NewObserver returns a new EVM chain observer
func NewObserver(
	ctx context.Context,
	evmCfg config.EVMConfig,
	evmClient interfaces.EVMRPCClient,
	chainParams observertypes.ChainParams,
	zetacoreClient interfaces.ZetacoreClient,
	tss interfaces.TSSSigner,
	dbpath string,
	logger base.Logger,
	ts *metrics.TelemetryServer,
) (*Observer, error) {
	// create base observer
	baseObserver, err := base.NewObserver(
		evmCfg.Chain,
		chainParams,
		zetacoreClient,
		tss,
		base.DefaultBlockCacheSize,
		base.DefaultHeaderCacheSize,
		ts,
		logger,
	)
	if err != nil {
		return nil, err
	}

	// create evm observer
	ob := &Observer{
		Observer:                      *baseObserver,
		evmClient:                     evmClient,
		evmJSONRPC:                    ethrpc.NewEthRPC(evmCfg.Endpoint),
		outboundPendingTransactions:   make(map[string]*ethtypes.Transaction),
		outboundConfirmedReceipts:     make(map[string]*ethtypes.Receipt),
		outboundConfirmedTransactions: make(map[string]*ethtypes.Transaction),
	}

	// open database and load data
	err = ob.LoadDB(ctx, dbpath)
	if err != nil {
		return nil, err
	}

	return ob, nil
}

// WithEvmClient attaches a new evm client to the observer
func (ob *Observer) WithEvmClient(client interfaces.EVMRPCClient) {
	ob.evmClient = client
}

// WithEvmJSONRPC attaches a new evm json rpc client to the observer
func (ob *Observer) WithEvmJSONRPC(client interfaces.EVMJSONRPCClient) {
	ob.evmJSONRPC = client
}

// SetChainParams sets the chain params for the observer
// Note: chain params is accessed concurrently
func (ob *Observer) SetChainParams(params observertypes.ChainParams) {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	ob.WithChainParams(params)
}

// GetChainParams returns the chain params for the observer
// Note: chain params is accessed concurrently
func (ob *Observer) GetChainParams() observertypes.ChainParams {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	return ob.ChainParams()
}

// GetConnectorContract returns the non-Eth connector address and binder
func (ob *Observer) GetConnectorContract() (ethcommon.Address, *zetaconnector.ZetaConnectorNonEth, error) {
	addr := ethcommon.HexToAddress(ob.GetChainParams().ConnectorContractAddress)
	contract, err := zetaconnector.NewZetaConnectorNonEth(addr, ob.evmClient)
	return addr, contract, err
}

// GetConnectorContractEth returns the Eth connector address and binder
func (ob *Observer) GetConnectorContractEth() (ethcommon.Address, *zetaconnectoreth.ZetaConnectorEth, error) {
	addr := ethcommon.HexToAddress(ob.GetChainParams().ConnectorContractAddress)
	contract, err := FetchConnectorContractEth(addr, ob.evmClient)
	return addr, contract, err
}

// GetERC20CustodyContract returns ERC20Custody contract address and binder
func (ob *Observer) GetERC20CustodyContract() (ethcommon.Address, *erc20custody.ERC20Custody, error) {
	addr := ethcommon.HexToAddress(ob.GetChainParams().Erc20CustodyContractAddress)
	contract, err := erc20custody.NewERC20Custody(addr, ob.evmClient)
	return addr, contract, err
}

// FetchConnectorContractEth returns the Eth connector address and binder
// TODO(revamp): move this to a contract package
func FetchConnectorContractEth(
	addr ethcommon.Address,
	client interfaces.EVMRPCClient,
) (*zetaconnectoreth.ZetaConnectorEth, error) {
	return zetaconnectoreth.NewZetaConnectorEth(addr, client)
}

// FetchZetaTokenContract returns the non-Eth ZETA token binder
// TODO(revamp): move this to a contract package
func FetchZetaTokenContract(
	addr ethcommon.Address,
	client interfaces.EVMRPCClient,
) (*zeta.ZetaNonEth, error) {
	return zeta.NewZetaNonEth(addr, client)
}

// Start all observation routines for the evm chain
func (ob *Observer) Start(ctx context.Context) {
	ob.Logger().Chain.Info().Msgf("observer is starting for chain %d", ob.Chain().ChainId)

	bg.Work(ctx, ob.WatchInbound, bg.WithName("WatchInbound"), bg.WithLogger(ob.Logger().Inbound))
	bg.Work(ctx, ob.WatchOutbound, bg.WithName("WatchOutbound"), bg.WithLogger(ob.Logger().Outbound))
	bg.Work(ctx, ob.WatchGasPrice, bg.WithName("WatchGasPrice"), bg.WithLogger(ob.Logger().GasPrice))
	bg.Work(ctx, ob.WatchInboundTracker, bg.WithName("WatchInboundTracker"), bg.WithLogger(ob.Logger().Inbound))
	bg.Work(ctx, ob.WatchRPCStatus, bg.WithName("WatchRPCStatus"), bg.WithLogger(ob.Logger().Chain))
}

// WatchRPCStatus watches the RPC status of the evm chain
// TODO(revamp): move ticker to ticker file
// TODO(revamp): move inner logic to a separate function
func (ob *Observer) WatchRPCStatus(ctx context.Context) error {
	ob.Logger().Chain.Info().Msgf("Starting RPC status check for chain %d", ob.Chain().ChainId)
	ticker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-ticker.C:
			if !ob.GetChainParams().IsSupported {
				continue
			}
			bn, err := ob.evmClient.BlockNumber(ctx)
			if err != nil {
				ob.Logger().Chain.Error().Err(err).Msg("RPC Status Check error: RPC down?")
				continue
			}
			gasPrice, err := ob.evmClient.SuggestGasPrice(ctx)
			if err != nil {
				ob.Logger().Chain.Error().Err(err).Msg("RPC Status Check error: RPC down?")
				continue
			}
			header, err := ob.evmClient.HeaderByNumber(ctx, new(big.Int).SetUint64(bn))
			if err != nil {
				ob.Logger().Chain.Error().Err(err).Msg("RPC Status Check error: RPC down?")
				continue
			}
			// #nosec G115 always in range
			blockTime := time.Unix(int64(header.Time), 0).UTC()
			elapsedSeconds := time.Since(blockTime).Seconds()
			if elapsedSeconds > 100 {
				ob.Logger().Chain.Warn().
					Msgf("RPC Status Check warning: RPC stale or chain stuck (check explorer)? Latest block %d timestamp is %.0fs ago", bn, elapsedSeconds)
				continue
			}
			ob.Logger().Chain.Info().
				Msgf("[OK] RPC status: latest block num %d, timestamp %s ( %.0fs ago), suggested gas price %d", header.Number, blockTime.String(), elapsedSeconds, gasPrice.Uint64())
		case <-ob.StopChannel():
			return nil
		}
	}
}

// SetPendingTx sets the pending transaction in memory
func (ob *Observer) SetPendingTx(nonce uint64, transaction *ethtypes.Transaction) {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	ob.outboundPendingTransactions[ob.GetTxID(nonce)] = transaction
}

// GetPendingTx gets the pending transaction from memory
func (ob *Observer) GetPendingTx(nonce uint64) *ethtypes.Transaction {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	return ob.outboundPendingTransactions[ob.GetTxID(nonce)]
}

// SetTxNReceipt sets the receipt and transaction in memory
func (ob *Observer) SetTxNReceipt(nonce uint64, receipt *ethtypes.Receipt, transaction *ethtypes.Transaction) {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	delete(ob.outboundPendingTransactions, ob.GetTxID(nonce)) // remove pending transaction, if any
	ob.outboundConfirmedReceipts[ob.GetTxID(nonce)] = receipt
	ob.outboundConfirmedTransactions[ob.GetTxID(nonce)] = transaction
}

// GetTxNReceipt gets the receipt and transaction from memory
func (ob *Observer) GetTxNReceipt(nonce uint64) (*ethtypes.Receipt, *ethtypes.Transaction) {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	receipt := ob.outboundConfirmedReceipts[ob.GetTxID(nonce)]
	transaction := ob.outboundConfirmedTransactions[ob.GetTxID(nonce)]
	return receipt, transaction
}

// IsTxConfirmed returns true if there is a confirmed tx for 'nonce'
func (ob *Observer) IsTxConfirmed(nonce uint64) bool {
	ob.Mu().Lock()
	defer ob.Mu().Unlock()
	return ob.outboundConfirmedReceipts[ob.GetTxID(nonce)] != nil &&
		ob.outboundConfirmedTransactions[ob.GetTxID(nonce)] != nil
}

// CheckTxInclusion returns nil only if tx is included at the position indicated by the receipt ([block, index])
func (ob *Observer) CheckTxInclusion(tx *ethtypes.Transaction, receipt *ethtypes.Receipt) error {
	block, err := ob.GetBlockByNumberCached(receipt.BlockNumber.Uint64())
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
		ob.RemoveCachedBlock(receipt.BlockNumber.Uint64()) // clean stale block from cache
		return fmt.Errorf("transaction at index %d has different hash %s, txHash %s nonce %d block %d",
			receipt.TransactionIndex, txAtIndex.Hash, tx.Hash(), tx.Nonce(), receipt.BlockNumber.Uint64())
	}

	return nil
}

// WatchGasPrice watches evm chain for gas prices and post to zetacore
// TODO(revamp): move ticker to ticker file
// TODO(revamp): move inner logic to a separate function
func (ob *Observer) WatchGasPrice(ctx context.Context) error {
	// report gas price right away as the ticker takes time to kick in
	err := ob.PostGasPrice(ctx)
	if err != nil {
		ob.Logger().GasPrice.Error().Err(err).Msgf("PostGasPrice error for chain %d", ob.Chain().ChainId)
	}

	// start gas price ticker
	ticker, err := clienttypes.NewDynamicTicker(
		fmt.Sprintf("EVM_WatchGasPrice_%d", ob.Chain().ChainId),
		ob.GetChainParams().GasPriceTicker,
	)
	if err != nil {
		ob.Logger().GasPrice.Error().Err(err).Msg("NewDynamicTicker error")
		return err
	}
	ob.Logger().GasPrice.Info().Msgf("WatchGasPrice started for chain %d with interval %d",
		ob.Chain().ChainId, ob.GetChainParams().GasPriceTicker)

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			if !ob.GetChainParams().IsSupported {
				continue
			}
			err = ob.PostGasPrice(ctx)
			if err != nil {
				ob.Logger().GasPrice.Error().Err(err).Msgf("PostGasPrice error for chain %d", ob.Chain().ChainId)
			}
			ticker.UpdateInterval(ob.GetChainParams().GasPriceTicker, ob.Logger().GasPrice)
		case <-ob.StopChannel():
			ob.Logger().GasPrice.Info().Msg("WatchGasPrice stopped")
			return nil
		}
	}
}

// PostGasPrice posts gas price to zetacore
// TODO(revamp): move to gas price file
func (ob *Observer) PostGasPrice(ctx context.Context) error {
	// GAS PRICE
	gasPrice, err := ob.evmClient.SuggestGasPrice(ctx)
	if err != nil {
		ob.Logger().GasPrice.Err(err).Msg("Err SuggestGasPrice:")
		return err
	}
	blockNum, err := ob.evmClient.BlockNumber(ctx)
	if err != nil {
		ob.Logger().GasPrice.Err(err).Msg("Err Fetching Most recent Block : ")
		return err
	}

	// SUPPLY
	supply := "100" // lockedAmount on ETH, totalSupply on other chains

	zetaHash, err := ob.ZetacoreClient().PostVoteGasPrice(ctx, ob.Chain(), gasPrice.Uint64(), supply, blockNum)
	if err != nil {
		ob.Logger().GasPrice.Err(err).Msg("PostGasPrice to zetacore failed")
		return err
	}
	_ = zetaHash

	return nil
}

// TransactionByHash query transaction by hash via JSON-RPC
// TODO(revamp): update this method as a pure RPC method that takes two parameters (jsonRPC, and txHash) and move to upper package to file rpc.go
func (ob *Observer) TransactionByHash(txHash string) (*ethrpc.Transaction, bool, error) {
	tx, err := ob.evmJSONRPC.EthGetTransactionByHash(txHash)
	if err != nil {
		return nil, false, err
	}
	err = evm.ValidateEvmTransaction(tx)
	if err != nil {
		return nil, false, err
	}
	return tx, tx.BlockNumber == nil, nil
}

// GetBlockHeaderCached get block header by number from cache
func (ob *Observer) GetBlockHeaderCached(ctx context.Context, blockNumber uint64) (*ethtypes.Header, error) {
	if result, ok := ob.HeaderCache().Get(blockNumber); ok {
		if header, ok := result.(*ethtypes.Header); ok {
			return header, nil
		}
		return nil, errors.New("cached value is not of type *ethtypes.Header")
	}
	header, err := ob.evmClient.HeaderByNumber(ctx, new(big.Int).SetUint64(blockNumber))
	if err != nil {
		return nil, err
	}
	ob.HeaderCache().Add(blockNumber, header)
	return header, nil
}

// GetBlockByNumberCached get block by number from cache
// returns block, ethrpc.Block, isFallback, isSkip, error
func (ob *Observer) GetBlockByNumberCached(blockNumber uint64) (*ethrpc.Block, error) {
	if result, ok := ob.BlockCache().Get(blockNumber); ok {
		if block, ok := result.(*ethrpc.Block); ok {
			return block, nil
		}
		return nil, errors.New("cached value is not of type *ethrpc.Block")
	}
	if blockNumber > math.MaxInt32 {
		return nil, fmt.Errorf("block number %d is too large", blockNumber)
	}
	// #nosec G115 always in range, checked above
	block, err := ob.BlockByNumber(int(blockNumber))
	if err != nil {
		return nil, err
	}
	ob.BlockCache().Add(blockNumber, block)
	return block, nil
}

// RemoveCachedBlock remove block from cache
func (ob *Observer) RemoveCachedBlock(blockNumber uint64) {
	ob.BlockCache().Remove(blockNumber)
}

// BlockByNumber query block by number via JSON-RPC
func (ob *Observer) BlockByNumber(blockNumber int) (*ethrpc.Block, error) {
	block, err := ob.evmJSONRPC.EthGetBlockByNumber(blockNumber, true)
	if err != nil {
		return nil, err
	}
	for i := range block.Transactions {
		err := evm.ValidateEvmTransaction(&block.Transactions[i])
		if err != nil {
			return nil, err
		}
	}
	return block, nil
}

// LoadDB open sql database and load data into EVM observer
// TODO(revamp): move to a db file
func (ob *Observer) LoadDB(ctx context.Context, dbPath string) error {
	if dbPath == "" {
		return errors.New("empty db path")
	}

	// open database
	err := ob.OpenDB(dbPath, "")
	if err != nil {
		return errors.Wrapf(err, "error OpenDB for chain %d", ob.Chain().ChainId)
	}

	// run auto migration
	// transaction and receipt tables are used nowhere but we still run migration in case they are needed in future
	err = ob.DB().AutoMigrate(
		&clienttypes.ReceiptSQLType{},
		&clienttypes.TransactionSQLType{},
	)
	if err != nil {
		return errors.Wrapf(err, "error AutoMigrate for chain %d", ob.Chain().ChainId)
	}

	// load last block scanned
	err = ob.LoadLastBlockScanned(ctx)

	return err
}

// LoadLastBlockScanned loads the last scanned block from the database
// TODO(revamp): move to a db file
func (ob *Observer) LoadLastBlockScanned(ctx context.Context) error {
	err := ob.Observer.LoadLastBlockScanned(ob.Logger().Chain)
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
	ob.Logger().Chain.Info().Msgf("chain %d starts scanning from block %d", ob.Chain().ChainId, ob.LastBlockScanned())

	return nil
}

// postBlockHeader posts the block header to zetacore
// TODO(revamp): move to a block header file
func (ob *Observer) postBlockHeader(ctx context.Context, tip uint64) error {
	bn := tip

	chainState, err := ob.ZetacoreClient().GetBlockHeaderChainState(ctx, ob.Chain().ChainId)
	if err == nil && chainState != nil && chainState.EarliestHeight > 0 {
		// #nosec G115 always positive
		bn = uint64(chainState.LatestHeight) + 1 // the next header to post
	}

	if bn > tip {
		return fmt.Errorf("postBlockHeader: must post block confirmed block header: %d > %d", bn, tip)
	}

	header, err := ob.GetBlockHeaderCached(ctx, bn)
	if err != nil {
		ob.Logger().Inbound.Error().Err(err).Msgf("postBlockHeader: error getting block: %d", bn)
		return err
	}
	headerRLP, err := rlp.EncodeToBytes(header)
	if err != nil {
		ob.Logger().Inbound.Error().Err(err).Msgf("postBlockHeader: error encoding block header: %d", bn)
		return err
	}

	_, err = ob.ZetacoreClient().PostVoteBlockHeader(
		ctx,
		ob.Chain().ChainId,
		header.Hash().Bytes(),
		header.Number.Int64(),
		proofs.NewEthereumHeader(headerRLP),
	)
	if err != nil {
		ob.Logger().Inbound.Error().Err(err).Msgf("postBlockHeader: error posting block header: %d", bn)
		return err
	}
	return nil
}
