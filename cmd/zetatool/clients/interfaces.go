package clients

import (
	"context"
	"math/big"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayevm.sol"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/zetaconnector.non-eth.sol"

	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	btcobserver "github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
)

// EVMClient defines the interface for EVM chain interactions
type EVMClient interface {
	// Basic transaction queries
	TransactionByHash(ctx context.Context, hash string) (*ethtypes.Transaction, bool, error)
	TransactionReceipt(ctx context.Context, hash string) (*ethtypes.Receipt, error)
	ChainID(ctx context.Context) (*big.Int, error)
	BlockNumber(ctx context.Context) (uint64, error)
	TransactionSender(ctx context.Context, tx *ethtypes.Transaction, blockHash ethcommon.Hash, txIndex uint) (ethcommon.Address, error)

	// Contract event parsing
	ParseConnectorZetaSent(log ethtypes.Log, connectorAddr string) (*zetaconnector.ZetaConnectorNonEthZetaSent, error)
	ParseCustodyDeposited(log ethtypes.Log, custodyAddr string) (*erc20custody.ERC20CustodyDeposited, error)
	ParseGatewayDeposited(log ethtypes.Log, gatewayAddr string) (*gatewayevm.GatewayEVMDeposited, error)
	ParseGatewayDepositedAndCalled(log ethtypes.Log, gatewayAddr string) (*gatewayevm.GatewayEVMDepositedAndCalled, error)
	ParseGatewayCalled(log ethtypes.Log, gatewayAddr string) (*gatewayevm.GatewayEVMCalled, error)
}

// BitcoinClient defines the interface for Bitcoin chain interactions
type BitcoinClient interface {
	Ping(ctx context.Context) error
	GetRawTransactionVerbose(ctx context.Context, txHash *chainhash.Hash) (*btcjson.TxRawResult, error)
	GetBlockVerbose(ctx context.Context, blockHash *chainhash.Hash) (*btcjson.GetBlockVerboseTxResult, error)
	GetBtcEventWithWitness(ctx context.Context, tx btcjson.TxRawResult, tssAddress string, blockNumber uint64, feeRateMultiplier float64, logger zerolog.Logger, netParams *chaincfg.Params) (*btcobserver.BTCInboundEvent, error)
}

// SolanaClient defines the interface for Solana chain interactions
type SolanaClient interface {
	GetTransaction(ctx context.Context, signature solana.Signature) (*solrpc.GetTransactionResult, error)
	ProcessTransactionResultWithAddressLookups(ctx context.Context, txResult *solrpc.GetTransactionResult, logger zerolog.Logger, signature solana.Signature) *solana.Transaction
	FilterInboundEvents(txResult *solrpc.GetTransactionResult, gatewayID solana.PublicKey, chainID int64, logger zerolog.Logger, tx *solana.Transaction) ([]*clienttypes.InboundEvent, error)
}

// SuiClient defines the interface for Sui chain interactions
type SuiClient interface {
	GetBalance(ctx context.Context, address string) (uint64, error)
}

// TONClient defines the interface for TON chain interactions
type TONClient interface {
	GetAccountBalance(ctx context.Context, address string) (uint64, error)
}

// ZetacoreClient defines the read-only interface for querying ZetaCore.
// This interface contains only the methods that zetatool needs for CCTX tracking.
type ZetacoreClient interface {
	// CCTX queries
	GetCctxByHash(ctx context.Context, hash string) (*crosschaintypes.CrossChainTx, error)
	InboundHashToCctxData(ctx context.Context, hash string) (*crosschaintypes.QueryInboundHashToCctxDataResponse, error)

	// Tracker queries
	GetOutboundTracker(ctx context.Context, chainID int64, nonce uint64) (*crosschaintypes.OutboundTracker, error)

	// Chain params
	GetChainParamsForChainID(ctx context.Context, chainID int64) (*observertypes.ChainParams, error)

	// TSS queries
	GetTssAddress(ctx context.Context, btcChainID int64) (*observertypes.QueryGetTssAddressResponse, error)
	GetEVMTSSAddress(ctx context.Context) (string, error)
	GetBTCTSSAddress(ctx context.Context, chainID int64) (string, error)

	// Ballot queries
	GetBallotByID(ctx context.Context, id string) (*observertypes.QueryBallotByIdentifierResponse, error)
}
