package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"

	tmrpcclient "github.com/cometbft/cometbft/rpc/client"
	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"

	evmmempool "github.com/cosmos/evm/mempool"
	"github.com/cosmos/evm/server/config"
	servertypes "github.com/cosmos/evm/server/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/zeta-chain/node/rpc/types"

	"cosmossdk.io/log"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BackendI implements the Cosmos and EVM backend.
type BackendI interface { //nolint: revive
	EVMBackend

	GetConfig() config.Config
}

// EVMBackend implements the functionality shared within ethereum namespaces
// as defined by EIP-1474: https://github.com/ethereum/EIPs/blob/master/EIPS/eip-1474.md
// Implemented by Backend.
type EVMBackend interface {
	// Node specific queries
	Accounts() ([]common.Address, error)
	Syncing() (interface{}, error)
	SetEtherbase(etherbase common.Address) bool
	SetGasPrice(gasPrice hexutil.Big) bool
	ImportRawKey(privkey, password string) (common.Address, error)
	ListAccounts() ([]common.Address, error)
	NewMnemonic(uid string, language keyring.Language, hdPath, bip39Passphrase string, algo keyring.SignatureAlgo) (*keyring.Record, error)
	UnprotectedAllowed() bool
	RPCGasCap() uint64            // global gas cap for eth_call over rpc: DoS protection
	RPCEVMTimeout() time.Duration // global timeout for eth_call over rpc: DoS protection
	RPCTxFeeCap() float64         // RPCTxFeeCap is the global transaction fee(price * gaslimit) cap for send-transaction variants. The unit is ether.
	RPCMinGasPrice() *big.Int

	// Sign Tx
	Sign(address common.Address, data hexutil.Bytes) (hexutil.Bytes, error)
	SendTransaction(args evmtypes.TransactionArgs) (common.Hash, error)
	SignTypedData(address common.Address, typedData apitypes.TypedData) (hexutil.Bytes, error)

	// Blocks Info
	BlockNumber() (hexutil.Uint64, error)
	GetHeaderByNumber(blockNum types.BlockNumber) (map[string]interface{}, error)
	GetHeaderByHash(hash common.Hash) (map[string]interface{}, error)
	GetBlockByNumber(blockNum types.BlockNumber, fullTx bool) (map[string]interface{}, error)
	GetBlockByHash(hash common.Hash, fullTx bool) (map[string]interface{}, error)
	GetBlockTransactionCountByHash(hash common.Hash) *hexutil.Uint
	GetBlockTransactionCountByNumber(blockNum types.BlockNumber) *hexutil.Uint
	CometBlockByNumber(blockNum types.BlockNumber) (*tmrpctypes.ResultBlock, error)
	CometBlockByHash(blockHash common.Hash) (*tmrpctypes.ResultBlock, error)
	BlockNumberFromComet(blockNrOrHash types.BlockNumberOrHash) (types.BlockNumber, error)
	BlockNumberFromCometByHash(blockHash common.Hash) (*big.Int, error)
	EthMsgsFromCometBlock(block *tmrpctypes.ResultBlock, blockRes *tmrpctypes.ResultBlockResults) ([]*evmtypes.MsgEthereumTx, []*types.TxResultAdditionalFields)
	BlockBloomFromCometBlock(blockRes *tmrpctypes.ResultBlockResults) (ethtypes.Bloom, error)
	HeaderByNumber(blockNum types.BlockNumber) (*ethtypes.Header, error)
	HeaderByHash(blockHash common.Hash) (*ethtypes.Header, error)
	RPCBlockFromCometBlock(resBlock *tmrpctypes.ResultBlock, blockRes *tmrpctypes.ResultBlockResults, fullTx bool) (map[string]interface{}, error)
	EthBlockByNumber(blockNum types.BlockNumber) (*ethtypes.Block, error)
	EthBlockFromCometBlock(resBlock *tmrpctypes.ResultBlock, blockRes *tmrpctypes.ResultBlockResults) (*ethtypes.Block, error)
	GetBlockReceipts(blockNrOrHash types.BlockNumberOrHash) ([]map[string]interface{}, error)

	// Account Info
	GetCode(address common.Address, blockNrOrHash types.BlockNumberOrHash) (hexutil.Bytes, error)
	GetBalance(address common.Address, blockNrOrHash types.BlockNumberOrHash) (*hexutil.Big, error)
	GetStorageAt(address common.Address, key string, blockNrOrHash types.BlockNumberOrHash) (hexutil.Bytes, error)
	GetProof(address common.Address, storageKeys []string, blockNrOrHash types.BlockNumberOrHash) (*types.AccountResult, error)
	GetTransactionCount(address common.Address, blockNum types.BlockNumber) (*hexutil.Uint64, error)

	// Chain Info
	ChainID() (*hexutil.Big, error)
	ChainConfig() *params.ChainConfig
	GlobalMinGasPrice() (*big.Int, error)
	BaseFee(blockRes *tmrpctypes.ResultBlockResults) (*big.Int, error)
	CurrentHeader() (*ethtypes.Header, error)
	PendingTransactions() ([]*sdk.Tx, error)
	GetCoinbase() (sdk.AccAddress, error)
	FeeHistory(blockCount math.HexOrDecimal64, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*types.FeeHistoryResult, error)
	SuggestGasTipCap(baseFee *big.Int) (*big.Int, error)

	// Tx Info
	GetTransactionByHash(txHash common.Hash) (*types.RPCTransaction, error)
	GetTxByEthHash(txHash common.Hash) (*servertypes.TxResult, *types.TxResultAdditionalFields, error)
	GetTxByTxIndex(height int64, txIndex uint) (*servertypes.TxResult, *types.TxResultAdditionalFields, error)
	GetTransactionByBlockAndIndex(block *tmrpctypes.ResultBlock, idx hexutil.Uint) (*types.RPCTransaction, error)
	GetTransactionReceipt(hash common.Hash) (map[string]interface{}, error)
	GetTransactionLogs(hash common.Hash) ([]*ethtypes.Log, error)
	GetTransactionByBlockHashAndIndex(hash common.Hash, idx hexutil.Uint) (*types.RPCTransaction, error)
	GetTransactionByBlockNumberAndIndex(blockNum types.BlockNumber, idx hexutil.Uint) (*types.RPCTransaction, error)
	CreateAccessList(args evmtypes.TransactionArgs, blockNrOrHash types.BlockNumberOrHash, overrides *json.RawMessage) (*types.AccessListResult, error)

	// Send Transaction
	Resend(args evmtypes.TransactionArgs, gasPrice *hexutil.Big, gasLimit *hexutil.Uint64) (common.Hash, error)
	SendRawTransaction(data hexutil.Bytes) (common.Hash, error)
	SetTxDefaults(args evmtypes.TransactionArgs) (evmtypes.TransactionArgs, error)
	EstimateGas(args evmtypes.TransactionArgs, blockNrOptional *types.BlockNumber) (hexutil.Uint64, error)
	DoCall(args evmtypes.TransactionArgs, blockNr types.BlockNumber, overrides *json.RawMessage) (*evmtypes.MsgEthereumTxResponse, error)
	GasPrice() (*hexutil.Big, error)

	// Filter API
	GetLogs(hash common.Hash) ([][]*ethtypes.Log, error)
	GetLogsByHeight(height *int64) ([][]*ethtypes.Log, error)
	BloomStatus() (uint64, uint64)

	// TxPool API
	Content() (map[string]map[string]map[string]*types.RPCTransaction, error)
	ContentFrom(address common.Address) (map[string]map[string]*types.RPCTransaction, error)
	Inspect() (map[string]map[string]map[string]string, error)
	Status() (map[string]hexutil.Uint, error)

	// Tracing
	TraceTransaction(hash common.Hash, config *types.TraceConfig) (interface{}, error)
	TraceBlock(height types.BlockNumber, config *types.TraceConfig, block *tmrpctypes.ResultBlock) ([]*evmtypes.TxTraceResult, error)
	TraceCall(args evmtypes.TransactionArgs, blockNrOrHash types.BlockNumberOrHash, config *types.TraceConfig) (interface{}, error)
}

var _ BackendI = (*Backend)(nil)

// ProcessBlocker is a function type that processes a block and its associated data
// for fee history calculation. It takes a Tendermint block, its corresponding
// Ethereum block representation, reward percentiles for fee estimation,
// block results, and a target fee history entry to populate.
//
// Parameters:
//   - tendermintBlock: The raw Tendermint block data
//   - ethBlock: The Ethereum-formatted block representation
//   - rewardPercentiles: Percentiles used for fee reward calculation
//   - tendermintBlockResult: Block execution results from Tendermint
//   - targetOneFeeHistory: The fee history entry to be populated
//
// Returns an error if block processing fails.
type ProcessBlocker func(
	tendermintBlock *tmrpctypes.ResultBlock,
	ethBlock *map[string]interface{},
	rewardPercentiles []float64,
	tendermintBlockResult *tmrpctypes.ResultBlockResults,
	targetOneFeeHistory *types.OneFeeHistory,
) error

// Backend implements the BackendI interface
type Backend struct {
	Ctx                 context.Context
	ClientCtx           client.Context
	RPCClient           tmrpcclient.SignClient
	QueryClient         *types.QueryClient // gRPC query client
	Logger              log.Logger
	EvmChainID          *big.Int
	Cfg                 config.Config
	AllowUnprotectedTxs bool
	Indexer             servertypes.EVMTxIndexer
	ProcessBlocker      ProcessBlocker
	Mempool             *evmmempool.ExperimentalEVMMempool
}

func (b *Backend) GetConfig() config.Config {
	return b.Cfg
}

// NewBackend creates a new Backend instance for cosmos and ethereum namespaces
func NewBackend(
	ctx *server.Context,
	logger log.Logger,
	clientCtx client.Context,
	allowUnprotectedTxs bool,
	indexer servertypes.EVMTxIndexer,
	mempool *evmmempool.ExperimentalEVMMempool,
) *Backend {
	appConf, err := config.GetConfig(ctx.Viper)
	if err != nil {
		panic(err)
	}

	rpcClient, ok := clientCtx.Client.(tmrpcclient.SignClient)
	if !ok {
		panic(fmt.Sprintf("invalid rpc client, expected: tmrpcclient.SignClient, got: %T", clientCtx.Client))
	}

	ethCfg := evmtypes.GetEthChainConfig()
	logger.Info("chain id from eth cfg", "chainId", ethCfg.ChainID.String())

	b := &Backend{
		Ctx:                 context.Background(),
		ClientCtx:           clientCtx,
		RPCClient:           rpcClient,
		QueryClient:         types.NewQueryClient(clientCtx),
		Logger:              logger.With("module", "backend"),
		EvmChainID:          ethCfg.ChainID,
		Cfg:                 appConf,
		AllowUnprotectedTxs: allowUnprotectedTxs,
		Indexer:             indexer,
		Mempool:             mempool,
	}
	b.ProcessBlocker = b.ProcessBlock
	return b
}
