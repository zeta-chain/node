package zetaclient

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/big"
	"os"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zeta.non-eth.sol"
	zetaconnectoreth "github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.eth.sol"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	lru "github.com/hashicorp/golang-lru"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/erc20custody.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/evm/zetaconnector.non-eth.sol"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	metricsPkg "github.com/zeta-chain/zetacore/zetaclient/metrics"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TxHashEnvelope struct {
	TxHash string
	Done   chan struct{}
}

type OutTx struct {
	SendHash string
	TxHash   string
	Nonce    int64
}
type EVMLog struct {
	ChainLogger          zerolog.Logger // Parent logger
	ExternalChainWatcher zerolog.Logger // Observes external Chains for incoming trasnactions
	WatchGasPrice        zerolog.Logger // Observes external Chains for Gas prices and posts to core
	ObserveOutTx         zerolog.Logger // Observes external Chains for Outgoing transactions

}

const (
	DonationMessage    = "I am rich!"
	TopicsZetaSent     = 3 // [signature, zetaTxSenderAddress, destinationChainId] https://github.com/zeta-chain/protocol-contracts/blob/d65814debf17648a6c67d757ba03646415842790/contracts/evm/ZetaConnector.base.sol#L34
	TopicsZetaReceived = 4 // [signature, sourceChainId, destinationAddress]       https://github.com/zeta-chain/protocol-contracts/blob/d65814debf17648a6c67d757ba03646415842790/contracts/evm/ZetaConnector.base.sol#L45
	TopicsZetaReverted = 3 // [signature, destinationChainId, internalSendHash]    https://github.com/zeta-chain/protocol-contracts/blob/d65814debf17648a6c67d757ba03646415842790/contracts/evm/ZetaConnector.base.sol#L54
	TopicsWithdrawn    = 3 // [signature, recipient, asset] https://github.com/zeta-chain/protocol-contracts/blob/d65814debf17648a6c67d757ba03646415842790/contracts/evm/ERC20Custody.sol#L43
	TopicsDeposited    = 2 // [signature, asset]            https://github.com/zeta-chain/protocol-contracts/blob/d65814debf17648a6c67d757ba03646415842790/contracts/evm/ERC20Custody.sol#L42
)

// EVMChainClient represents the chain configuration for an EVM chain
// Filled with above constants depending on chain
type EVMChainClient struct {
	*ChainMetrics
	chain                      common.Chain
	evmClient                  EVMRPCClient
	KlaytnClient               KlaytnRPCClient
	zetaClient                 ZetaCoreBridger
	Tss                        TSSSigner
	lastBlockScanned           uint64
	lastBlock                  uint64
	BlockTimeExternalChain     uint64 // block time in seconds
	txWatchList                map[ethcommon.Hash]string
	Mu                         *sync.Mutex
	db                         *gorm.DB
	outTxPendingTransactions   map[string]*ethtypes.Transaction
	outTXConfirmedReceipts     map[string]*ethtypes.Receipt
	outTXConfirmedTransactions map[string]*ethtypes.Transaction
	MinNonce                   int64
	MaxNonce                   int64
	OutTxChan                  chan OutTx // send to this channel if you want something back!
	stop                       chan struct{}
	fileLogger                 *zerolog.Logger // for critical info
	logger                     EVMLog
	cfg                        *config.Config
	params                     observertypes.ChainParams
	ts                         *TelemetryServer

	BlockCache *lru.Cache
}

var _ ChainClient = (*EVMChainClient)(nil)

// NewEVMChainClient returns a new configuration based on supplied target chain
func NewEVMChainClient(
	bridge ZetaCoreBridger,
	tss TSSSigner,
	dbpath string,
	metrics *metricsPkg.Metrics,
	logger zerolog.Logger,
	cfg *config.Config,
	evmCfg config.EVMConfig,
	ts *TelemetryServer,
) (*EVMChainClient, error) {
	ob := EVMChainClient{
		ChainMetrics: NewChainMetrics(evmCfg.Chain.ChainName.String(), metrics),
		ts:           ts,
	}
	chainLogger := logger.With().Str("chain", evmCfg.Chain.ChainName.String()).Logger()
	ob.logger = EVMLog{
		ChainLogger:          chainLogger,
		ExternalChainWatcher: chainLogger.With().Str("module", "ExternalChainWatcher").Logger(),
		WatchGasPrice:        chainLogger.With().Str("module", "WatchGasPrice").Logger(),
		ObserveOutTx:         chainLogger.With().Str("module", "ObserveOutTx").Logger(),
	}
	ob.cfg = cfg
	ob.params = evmCfg.ChainParams
	ob.stop = make(chan struct{})
	ob.chain = evmCfg.Chain
	ob.Mu = &sync.Mutex{}
	ob.zetaClient = bridge
	ob.txWatchList = make(map[ethcommon.Hash]string)
	ob.Tss = tss
	ob.outTxPendingTransactions = make(map[string]*ethtypes.Transaction)
	ob.outTXConfirmedReceipts = make(map[string]*ethtypes.Receipt)
	ob.outTXConfirmedTransactions = make(map[string]*ethtypes.Transaction)
	ob.OutTxChan = make(chan OutTx, 100)

	logFile, err := os.OpenFile(ob.chain.ChainName.String()+"_debug.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Error().Err(err).Msgf("there was an error creating a logFile chain %s", ob.chain.ChainName.String())
	}
	fileLogger := zerolog.New(logFile).With().Logger()
	ob.fileLogger = &fileLogger

	ob.logger.ChainLogger.Info().Msgf("Chain %s endpoint %s", ob.chain.ChainName.String(), evmCfg.Endpoint)
	client, err := ethclient.Dial(evmCfg.Endpoint)
	if err != nil {
		ob.logger.ChainLogger.Error().Err(err).Msg("eth Client Dial")
		return nil, err
	}
	ob.evmClient = client

	ob.BlockCache, err = lru.New(1000)
	if err != nil {
		ob.logger.ChainLogger.Error().Err(err).Msg("failed to create block cache")
		return nil, err
	}

	if ob.chain.IsKlaytnChain() {
		client, err := Dial(evmCfg.Endpoint)
		if err != nil {
			ob.logger.ChainLogger.Err(err).Msg("klaytn Client Dial")
			return nil, err
		}
		ob.KlaytnClient = client
	}

	// create metric counters
	err = ob.RegisterPromCounter("rpc_getFilterLogs_count", "Number of getLogs")
	if err != nil {
		return nil, err
	}
	err = ob.RegisterPromCounter("rpc_getBlockByNumber_count", "Number of getBlockByNumber")
	if err != nil {
		return nil, err
	}
	err = ob.RegisterPromGauge(metricsPkg.PendingTxs, "Number of pending transactions")
	if err != nil {
		return nil, err
	}

	err = ob.LoadDB(dbpath, ob.chain)
	if err != nil {
		return nil, err
	}

	ob.logger.ChainLogger.Info().Msgf("%s: start scanning from block %d", ob.chain.String(), ob.GetLastBlockHeightScanned())

	return &ob, nil
}
func (ob *EVMChainClient) WithChain(chain common.Chain) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.chain = chain
}
func (ob *EVMChainClient) WithLogger(logger zerolog.Logger) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.logger = EVMLog{
		ChainLogger:          logger,
		ExternalChainWatcher: logger.With().Str("module", "ExternalChainWatcher").Logger(),
		WatchGasPrice:        logger.With().Str("module", "WatchGasPrice").Logger(),
		ObserveOutTx:         logger.With().Str("module", "ObserveOutTx").Logger(),
	}
}

func (ob *EVMChainClient) WithEvmClient(client *ethclient.Client) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.evmClient = client
}

func (ob *EVMChainClient) WithZetaClient(bridge *ZetaCoreBridge) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.zetaClient = bridge
}

func (ob *EVMChainClient) WithParams(params observertypes.ChainParams) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.params = params
}

func (ob *EVMChainClient) SetConfig(cfg *config.Config) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.cfg = cfg
}

func (ob *EVMChainClient) SetChainParams(params observertypes.ChainParams) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.params = params
}

func (ob *EVMChainClient) GetChainParams() observertypes.ChainParams {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	return ob.params
}

func (ob *EVMChainClient) GetConnectorContract() (ethcommon.Address, *zetaconnector.ZetaConnectorNonEth, error) {
	addr := ethcommon.HexToAddress(ob.GetChainParams().ConnectorContractAddress)
	contract, err := FetchConnectorContract(addr, ob.evmClient)
	return addr, contract, err
}

func (ob *EVMChainClient) GetConnectorContractEth() (ethcommon.Address, *zetaconnectoreth.ZetaConnectorEth, error) {
	addr := ethcommon.HexToAddress(ob.GetChainParams().ConnectorContractAddress)
	contract, err := FetchConnectorContractEth(addr, ob.evmClient)
	return addr, contract, err
}

func (ob *EVMChainClient) GetZetaTokenNonEthContract() (ethcommon.Address, *zeta.ZetaNonEth, error) {
	addr := ethcommon.HexToAddress(ob.GetChainParams().ZetaTokenContractAddress)
	contract, err := FetchZetaZetaNonEthTokenContract(addr, ob.evmClient)
	return addr, contract, err
}

func (ob *EVMChainClient) GetERC20CustodyContract() (ethcommon.Address, *erc20custody.ERC20Custody, error) {
	addr := ethcommon.HexToAddress(ob.GetChainParams().Erc20CustodyContractAddress)
	contract, err := FetchERC20CustodyContract(addr, ob.evmClient)
	return addr, contract, err
}

func FetchConnectorContract(addr ethcommon.Address, client EVMRPCClient) (*zetaconnector.ZetaConnectorNonEth, error) {
	return zetaconnector.NewZetaConnectorNonEth(addr, client)
}

func FetchConnectorContractEth(addr ethcommon.Address, client EVMRPCClient) (*zetaconnectoreth.ZetaConnectorEth, error) {
	return zetaconnectoreth.NewZetaConnectorEth(addr, client)
}

func FetchZetaZetaNonEthTokenContract(addr ethcommon.Address, client EVMRPCClient) (*zeta.ZetaNonEth, error) {
	return zeta.NewZetaNonEth(addr, client)
}

func FetchERC20CustodyContract(addr ethcommon.Address, client EVMRPCClient) (*erc20custody.ERC20Custody, error) {
	return erc20custody.NewERC20Custody(addr, client)
}

func (ob *EVMChainClient) Start() {
	go ob.ExternalChainWatcherForNewInboundTrackerSuggestions()
	go ob.ExternalChainWatcher() // Observes external Chains for incoming trasnactions
	go ob.WatchGasPrice()        // Observes external Chains for Gas prices and posts to core
	go ob.observeOutTx()         // Populates receipts and confirmed outbound transactions
}

func (ob *EVMChainClient) Stop() {
	ob.logger.ChainLogger.Info().Msgf("ob %s is stopping", ob.chain.String())
	close(ob.stop) // this notifies all goroutines to stop

	ob.logger.ChainLogger.Info().Msg("closing ob.db")
	dbInst, err := ob.db.DB()
	if err != nil {
		ob.logger.ChainLogger.Info().Msg("error getting database instance")
	}
	err = dbInst.Close()
	if err != nil {
		ob.logger.ChainLogger.Error().Err(err).Msg("error closing database")
	}

	ob.logger.ChainLogger.Info().Msgf("%s observer stopped", ob.chain.String())
}

// returns: isIncluded, isConfirmed, Error
// If isConfirmed, it also post to ZetaCore
func (ob *EVMChainClient) IsSendOutTxProcessed(sendHash string, nonce uint64, cointype common.CoinType, logger zerolog.Logger) (bool, bool, error) {
	params := ob.GetChainParams()
	receipt, transaction := ob.GetTxNReceipt(nonce)
	if receipt == nil || transaction == nil { // not confirmed yet
		return false, false, nil
	}

	sendID := fmt.Sprintf("%s-%d", ob.chain.String(), nonce)
	logger = logger.With().Str("sendID", sendID).Logger()
	if cointype == common.CoinType_Cmd {
		recvStatus := common.ReceiveStatus_Failed
		if receipt.Status == 1 {
			recvStatus = common.ReceiveStatus_Success
		}
		zetaTxHash, ballot, err := ob.zetaClient.PostVoteOutbound(
			sendHash,
			receipt.TxHash.Hex(),
			receipt.BlockNumber.Uint64(),
			receipt.GasUsed,
			transaction.GasPrice(),
			transaction.Gas(),
			transaction.Value(),
			recvStatus,
			ob.chain,
			nonce,
			common.CoinType_Cmd,
		)
		if err != nil {
			logger.Error().Err(err).Msgf("error posting confirmation to meta core for cctx %s nonce %d", sendHash, nonce)
		} else if zetaTxHash != "" {
			logger.Info().Msgf("Zeta tx hash: %s cctx %s nonce %d ballot %s", zetaTxHash, sendHash, nonce, ballot)
		}
		return true, true, nil

	} else if cointype == common.CoinType_Gas { // the outbound is a regular Ether/BNB/Matic transfer; no need to check events
		if receipt.Status == 1 {
			zetaTxHash, ballot, err := ob.zetaClient.PostVoteOutbound(
				sendHash,
				receipt.TxHash.Hex(),
				receipt.BlockNumber.Uint64(),
				receipt.GasUsed,
				transaction.GasPrice(),
				transaction.Gas(),
				transaction.Value(),
				common.ReceiveStatus_Success,
				ob.chain,
				nonce,
				common.CoinType_Gas,
			)
			if err != nil {
				logger.Error().Err(err).Msgf("error posting confirmation to meta core for cctx %s nonce %d", sendHash, nonce)
			} else if zetaTxHash != "" {
				logger.Info().Msgf("Zeta tx hash: %s cctx %s nonce %d ballot %s", zetaTxHash, sendHash, nonce, ballot)
			}
			return true, true, nil
		} else if receipt.Status == 0 { // the same as below events flow
			logger.Info().Msgf("Found (failed tx) sendHash %s on chain %s txhash %s", sendHash, ob.chain.String(), receipt.TxHash.Hex())
			zetaTxHash, ballot, err := ob.zetaClient.PostVoteOutbound(
				sendHash,
				receipt.TxHash.Hex(),
				receipt.BlockNumber.Uint64(),
				receipt.GasUsed,
				transaction.GasPrice(),
				transaction.Gas(),
				big.NewInt(0),
				common.ReceiveStatus_Failed,
				ob.chain,
				nonce,
				common.CoinType_Gas,
			)
			if err != nil {
				logger.Error().Err(err).Msgf("PostVoteOutbound error in WatchTxHashWithTimeout; zeta tx hash %s cctx %s nonce %d", zetaTxHash, sendHash, nonce)
			} else if zetaTxHash != "" {
				logger.Info().Msgf("Zeta tx hash: %s cctx %s nonce %d ballot %s", zetaTxHash, sendHash, nonce, ballot)
			}
			return true, true, nil
		}
	} else if cointype == common.CoinType_Zeta { // the outbound is a Zeta transfer; need to check events ZetaReceived
		if receipt.Status == 1 {
			logs := receipt.Logs
			for _, vLog := range logs {
				confHeight := vLog.BlockNumber + params.ConfirmationCount
				// TODO rewrite this to return early if not confirmed
				connectorAddr, connector, err := ob.GetConnectorContract()
				if err != nil {
					return false, false, fmt.Errorf("error getting connector contract: %w", err)
				}
				receivedLog, err := connector.ZetaConnectorNonEthFilterer.ParseZetaReceived(*vLog)
				if err == nil {
					logger.Info().Msgf("Found (outTx) sendHash %s on chain %s txhash %s", sendHash, ob.chain.String(), vLog.TxHash.Hex())
					if confHeight <= ob.GetLastBlockHeight() {
						logger.Info().Msg("Confirmed! Sending PostConfirmation to zetacore...")
						// sanity check tx event
						err = ob.CheckEvmTxLog(vLog, connectorAddr, transaction.Hash().Hex(), TopicsZetaReceived)
						if err != nil {
							logger.Error().Err(err).Msgf("CheckEvmTxLog error on ZetaReceived event, chain %d nonce %d txhash %s", ob.chain.ChainId, nonce, transaction.Hash().Hex())
							return false, false, err
						}
						sendhash := vLog.Topics[3].Hex()
						//var rxAddress string = ethcommon.HexToAddress(vLog.Topics[1].Hex()).Hex()
						mMint := receivedLog.ZetaValue
						zetaTxHash, ballot, err := ob.zetaClient.PostVoteOutbound(
							sendhash,
							vLog.TxHash.Hex(),
							vLog.BlockNumber,
							receipt.GasUsed,
							transaction.GasPrice(),
							transaction.Gas(),
							mMint,
							common.ReceiveStatus_Success,
							ob.chain,
							nonce,
							common.CoinType_Zeta,
						)
						if err != nil {
							logger.Error().Err(err).Msgf("error posting confirmation to meta core for cctx %s nonce %d", sendHash, nonce)
							continue
						} else if zetaTxHash != "" {
							logger.Info().Msgf("Zeta tx hash: %s cctx %s nonce %d ballot %s", zetaTxHash, sendHash, nonce, ballot)
						}
						return true, true, nil
					}
					logger.Info().Msgf("Included; %d blocks before confirmed! chain %s nonce %d", confHeight-ob.GetLastBlockHeight(), ob.chain.String(), nonce)
					return true, false, nil
				}
				revertedLog, err := connector.ZetaConnectorNonEthFilterer.ParseZetaReverted(*vLog)
				if err == nil {
					logger.Info().Msgf("Found (revertTx) sendHash %s on chain %s txhash %s", sendHash, ob.chain.String(), vLog.TxHash.Hex())
					if confHeight <= ob.GetLastBlockHeight() {
						logger.Info().Msg("Confirmed! Sending PostConfirmation to zetacore...")
						// sanity check tx event
						err = ob.CheckEvmTxLog(vLog, connectorAddr, transaction.Hash().Hex(), TopicsZetaReverted)
						if err != nil {
							logger.Error().Err(err).Msgf("CheckEvmTxLog error on ZetaReverted event, chain %d nonce %d txhash %s", ob.chain.ChainId, nonce, transaction.Hash().Hex())
							return false, false, err
						}
						sendhash := vLog.Topics[2].Hex()
						mMint := revertedLog.RemainingZetaValue
						zetaTxHash, ballot, err := ob.zetaClient.PostVoteOutbound(
							sendhash,
							vLog.TxHash.Hex(),
							vLog.BlockNumber,
							receipt.GasUsed,
							transaction.GasPrice(),
							transaction.Gas(),
							mMint,
							common.ReceiveStatus_Success,
							ob.chain,
							nonce,
							common.CoinType_Zeta,
						)
						if err != nil {
							logger.Err(err).Msgf("error posting confirmation to meta core for cctx %s nonce %d", sendHash, nonce)
							continue
						} else if zetaTxHash != "" {
							logger.Info().Msgf("Zeta tx hash: %s cctx %s nonce %d ballot %s", zetaTxHash, sendHash, nonce, ballot)
						}
						return true, true, nil
					}
					logger.Info().Msgf("Included; %d blocks before confirmed! chain %s nonce %d", confHeight-ob.GetLastBlockHeight(), ob.chain.String(), nonce)
					return true, false, nil
				}
			}
		} else if receipt.Status == 0 {
			//FIXME: check nonce here by getTransaction RPC
			logger.Info().Msgf("Found (failed tx) sendHash %s on chain %s txhash %s", sendHash, ob.chain.String(), receipt.TxHash.Hex())
			zetaTxHash, ballot, err := ob.zetaClient.PostVoteOutbound(
				sendHash,
				receipt.TxHash.Hex(),
				receipt.BlockNumber.Uint64(),
				receipt.GasUsed,
				transaction.GasPrice(),
				transaction.Gas(),
				big.NewInt(0),
				common.ReceiveStatus_Failed,
				ob.chain,
				nonce,
				common.CoinType_Zeta,
			)
			if err != nil {
				logger.Error().Err(err).Msgf("error posting confirmation to meta core for cctx %s nonce %d", sendHash, nonce)
			} else if zetaTxHash != "" {
				logger.Info().Msgf("Zeta tx hash: %s cctx %s nonce %d ballot %s", zetaTxHash, sendHash, nonce, ballot)
			}
			return true, true, nil
		}
	} else if cointype == common.CoinType_ERC20 {
		if receipt.Status == 1 {
			logs := receipt.Logs
			addrCustody, ERC20Custody, err := ob.GetERC20CustodyContract()
			if err != nil {
				logger.Warn().Msgf("NewERC20Custody err: %s", err)
			}
			for _, vLog := range logs {
				event, err := ERC20Custody.ParseWithdrawn(*vLog)
				confHeight := vLog.BlockNumber + params.ConfirmationCount
				if err == nil {
					logger.Info().Msgf("Found (ERC20Custody.Withdrawn Event) sendHash %s on chain %s txhash %s", sendHash, ob.chain.String(), vLog.TxHash.Hex())
					// sanity check tx event
					err = ob.CheckEvmTxLog(vLog, addrCustody, transaction.Hash().Hex(), TopicsWithdrawn)
					if err != nil {
						logger.Error().Err(err).Msgf("CheckEvmTxLog error on Withdrawn event, chain %d nonce %d txhash %s", ob.chain.ChainId, nonce, transaction.Hash().Hex())
						return false, false, err
					}
					if confHeight <= ob.GetLastBlockHeight() {
						logger.Info().Msg("Confirmed! Sending PostConfirmation to zetacore...")
						zetaTxHash, ballot, err := ob.zetaClient.PostVoteOutbound(
							sendHash,
							vLog.TxHash.Hex(),
							vLog.BlockNumber,
							receipt.GasUsed,
							transaction.GasPrice(),
							transaction.Gas(),
							event.Amount,
							common.ReceiveStatus_Success,
							ob.chain,
							nonce,
							common.CoinType_ERC20,
						)
						if err != nil {
							logger.Error().Err(err).Msgf("error posting confirmation to meta core for cctx %s nonce %d", sendHash, nonce)
							continue
						} else if zetaTxHash != "" {
							logger.Info().Msgf("Zeta tx hash: %s cctx %s nonce %d ballot %s", zetaTxHash, sendHash, nonce, ballot)
						}
						return true, true, nil
					}
					logger.Info().Msgf("Included; %d blocks before confirmed! chain %s nonce %d", confHeight-ob.GetLastBlockHeight(), ob.chain.String(), nonce)
					return true, false, nil
				}
			}
		} else {
			logger.Info().Msgf("Found (failed tx) sendHash %s on chain %s txhash %s", sendHash, ob.chain.String(), receipt.TxHash.Hex())
			zetaTxHash, ballot, err := ob.zetaClient.PostVoteOutbound(
				sendHash,
				receipt.TxHash.Hex(),
				receipt.BlockNumber.Uint64(),
				receipt.GasUsed,
				transaction.GasPrice(),
				transaction.Gas(),
				big.NewInt(0),
				common.ReceiveStatus_Failed,
				ob.chain,
				nonce,
				common.CoinType_ERC20,
			)
			if err != nil {
				logger.Error().Err(err).Msgf("PostVoteOutbound error in WatchTxHashWithTimeout; zeta tx hash %s", zetaTxHash)
			} else if zetaTxHash != "" {
				logger.Info().Msgf("Zeta tx hash: %s cctx %s nonce %d ballot %s", zetaTxHash, sendHash, nonce, ballot)
			}
			return true, true, nil
		}
	}

	return false, false, nil
}

// FIXME: there's a chance that a txhash in OutTxChan may not deliver when Stop() is called
// observeOutTx periodically checks all the txhash in potential outbound txs
func (ob *EVMChainClient) observeOutTx() {
	// read env variables if set
	timeoutNonce, err := strconv.Atoi(os.Getenv("OS_TIMEOUT_NONCE"))
	if err != nil || timeoutNonce <= 0 {
		timeoutNonce = 100 * 3 // process up to 100 hashes
	}
	ob.logger.ObserveOutTx.Info().Msgf("observeOutTx: using timeoutNonce %d seconds", timeoutNonce)

	ticker, err := NewDynamicTicker(fmt.Sprintf("EVM_observeOutTx_%d", ob.chain.ChainId), ob.GetChainParams().OutTxTicker)
	if err != nil {
		ob.logger.ObserveOutTx.Error().Err(err).Msg("failed to create ticker")
		return
	}

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			trackers, err := ob.zetaClient.GetAllOutTxTrackerByChain(ob.chain.ChainId, Ascending)
			if err != nil {
				continue
			}
			//FIXME: remove this timeout here to ensure that all trackers are queried
			outTimeout := time.After(time.Duration(timeoutNonce) * time.Second)
		TRACKERLOOP:
			for _, tracker := range trackers {
				nonceInt := tracker.Nonce
				if ob.isTxConfirmed(nonceInt) { // Go to next tracker if this one already has a confirmed tx
					continue
				}
				txCount := 0
				var receipt *ethtypes.Receipt
				var transaction *ethtypes.Transaction
				for _, txHash := range tracker.HashList {
					select {
					case <-outTimeout:
						ob.logger.ObserveOutTx.Warn().Msgf("observeOutTx: timeout on chain %d nonce %d", ob.chain.ChainId, nonceInt)
						break TRACKERLOOP
					default:
						if recpt, tx, ok := ob.checkConfirmedTx(txHash.TxHash, nonceInt); ok {
							txCount++
							receipt = recpt
							transaction = tx
							ob.logger.ObserveOutTx.Info().Msgf("observeOutTx: confirmed outTx %s for chain %d nonce %d", txHash.TxHash, ob.chain.ChainId, nonceInt)
							if txCount > 1 {
								ob.logger.ObserveOutTx.Error().Msgf(
									"observeOutTx: checkConfirmedTx passed, txCount %d chain %d nonce %d receipt %v transaction %v", txCount, ob.chain.ChainId, nonceInt, receipt, transaction)
							}
						}
					}
				}
				if txCount == 1 { // should be only one txHash confirmed for each nonce.
					ob.SetTxNReceipt(nonceInt, receipt, transaction)
				} else if txCount > 1 { // should not happen. We can't tell which txHash is true. It might happen (e.g. glitchy/hacked endpoint)
					ob.logger.ObserveOutTx.Error().Msgf("observeOutTx: confirmed multiple (%d) outTx for chain %d nonce %d", txCount, ob.chain.ChainId, nonceInt)
				}
			}
			ticker.UpdateInterval(ob.GetChainParams().OutTxTicker, ob.logger.ObserveOutTx)
		case <-ob.stop:
			ob.logger.ObserveOutTx.Info().Msg("observeOutTx: stopped")
			return
		}
	}
}

// SetPendingTx sets the pending transaction in memory
func (ob *EVMChainClient) SetPendingTx(nonce uint64, transaction *ethtypes.Transaction) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	ob.outTxPendingTransactions[ob.GetTxID(nonce)] = transaction
}

// GetPendingTx gets the pending transaction from memory
func (ob *EVMChainClient) GetPendingTx(nonce uint64) *ethtypes.Transaction {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	return ob.outTxPendingTransactions[ob.GetTxID(nonce)]
}

// SetTxNReceipt sets the receipt and transaction in memory
func (ob *EVMChainClient) SetTxNReceipt(nonce uint64, receipt *ethtypes.Receipt, transaction *ethtypes.Transaction) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	delete(ob.outTxPendingTransactions, ob.GetTxID(nonce)) // remove pending transaction, if any
	ob.outTXConfirmedReceipts[ob.GetTxID(nonce)] = receipt
	ob.outTXConfirmedTransactions[ob.GetTxID(nonce)] = transaction
}

// GetTxNReceipt gets the receipt and transaction from memory
func (ob *EVMChainClient) GetTxNReceipt(nonce uint64) (*ethtypes.Receipt, *ethtypes.Transaction) {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	receipt := ob.outTXConfirmedReceipts[ob.GetTxID(nonce)]
	transaction := ob.outTXConfirmedTransactions[ob.GetTxID(nonce)]
	return receipt, transaction
}

// isTxConfirmed returns true if there is a confirmed tx for 'nonce'
func (ob *EVMChainClient) isTxConfirmed(nonce uint64) bool {
	ob.Mu.Lock()
	defer ob.Mu.Unlock()
	return ob.outTXConfirmedReceipts[ob.GetTxID(nonce)] != nil && ob.outTXConfirmedTransactions[ob.GetTxID(nonce)] != nil
}

// checkConfirmedTx checks if a txHash is confirmed
// returns (receipt, transaction, true) if confirmed or (nil, nil, false) otherwise
func (ob *EVMChainClient) checkConfirmedTx(txHash string, nonce uint64) (*ethtypes.Receipt, *ethtypes.Transaction, bool) {
	ctxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// query transaction
	transaction, isPending, err := ob.evmClient.TransactionByHash(ctxt, ethcommon.HexToHash(txHash))
	if err != nil {
		log.Error().Err(err).Msgf("confirmTxByHash: TransactionByHash error, txHash %s nonce %d", txHash, nonce)
		return nil, nil, false
	}
	if transaction == nil { // should not happen
		log.Error().Msgf("confirmTxByHash: transaction is nil for txHash %s nonce %d", txHash, nonce)
		return nil, nil, false
	}

	// check tx sender and nonce
	signer := ethtypes.NewLondonSigner(big.NewInt(ob.chain.ChainId))
	from, err := signer.Sender(transaction)
	if err != nil {
		log.Error().Err(err).Msgf("confirmTxByHash: local recovery of sender address failed for txHash %s chain %d", transaction.Hash().Hex(), ob.chain.ChainId)
		return nil, nil, false
	}
	if from != ob.Tss.EVMAddress() { // must be TSS address
		log.Error().Msgf("confirmTxByHash: sender %s for txHash %s chain %d is not TSS address %s",
			from.Hex(), transaction.Hash().Hex(), ob.chain.ChainId, ob.Tss.EVMAddress().Hex())
		return nil, nil, false
	}
	if transaction.Nonce() != nonce { // must match cctx nonce
		log.Error().Msgf("confirmTxByHash: txHash %s nonce mismatch: wanted %d, got tx nonce %d", txHash, nonce, transaction.Nonce())
		return nil, nil, false
	}

	// save pending transaction
	if isPending {
		ob.SetPendingTx(nonce, transaction)
		return nil, nil, false
	}

	// query receipt
	receipt, err := ob.evmClient.TransactionReceipt(ctxt, ethcommon.HexToHash(txHash))
	if err != nil {
		if err != ethereum.NotFound {
			log.Warn().Err(err).Msgf("confirmTxByHash: TransactionReceipt error, txHash %s nonce %d", txHash, nonce)
		}
		return nil, nil, false
	}
	if receipt == nil { // should not happen
		log.Error().Msgf("confirmTxByHash: receipt is nil for txHash %s nonce %d", txHash, nonce)
		return nil, nil, false
	}

	// check confirmations
	confHeight := receipt.BlockNumber.Uint64() + ob.GetChainParams().ConfirmationCount
	if confHeight >= math.MaxInt64 {
		log.Error().Msgf("confirmTxByHash: confHeight is too large for txHash %s nonce %d", txHash, nonce)
		return nil, nil, false
	}
	if confHeight > ob.GetLastBlockHeight() {
		log.Info().Msgf("confirmTxByHash: txHash %s nonce %d included but not confirmed: receipt block %d, current block %d",
			txHash, nonce, receipt.BlockNumber, ob.GetLastBlockHeight())
		return nil, nil, false
	}

	return receipt, transaction, true
}

// SetLastBlockHeightScanned set last block height scanned (not necessarily caught up with external block; could be slow/paused)
func (ob *EVMChainClient) SetLastBlockHeightScanned(height uint64) {
	atomic.StoreUint64(&ob.lastBlockScanned, height)
	ob.ts.SetLastScannedBlockNumber(ob.chain.ChainId, height)
}

// GetLastBlockHeightScanned get last block height scanned (not necessarily caught up with external block; could be slow/paused)
func (ob *EVMChainClient) GetLastBlockHeightScanned() uint64 {
	height := atomic.LoadUint64(&ob.lastBlockScanned)
	return height
}

// SetLastBlockHeight set external last block height
func (ob *EVMChainClient) SetLastBlockHeight(height uint64) {
	if height >= math.MaxInt64 {
		panic("lastBlock is too large")
	}
	atomic.StoreUint64(&ob.lastBlock, height)
}

// GetLastBlockHeight get external last block height
func (ob *EVMChainClient) GetLastBlockHeight() uint64 {
	height := atomic.LoadUint64(&ob.lastBlock)
	if height >= math.MaxInt64 {
		panic("lastBlock is too large")
	}
	return height
}

func (ob *EVMChainClient) ExternalChainWatcher() {
	ticker, err := NewDynamicTicker(fmt.Sprintf("EVM_ExternalChainWatcher_%d", ob.chain.ChainId), ob.GetChainParams().InTxTicker)
	if err != nil {
		ob.logger.ExternalChainWatcher.Error().Err(err).Msg("NewDynamicTicker error")
		return
	}

	defer ticker.Stop()
	ob.logger.ExternalChainWatcher.Info().Msg("ExternalChainWatcher started")
	sampledLogger := ob.logger.ExternalChainWatcher.Sample(&zerolog.BasicSampler{N: 10})
	for {
		select {
		case <-ticker.C():
			err := ob.observeInTX(sampledLogger)
			if err != nil {
				ob.logger.ExternalChainWatcher.Err(err).Msg("observeInTX error")
			}
			ticker.UpdateInterval(ob.GetChainParams().InTxTicker, ob.logger.ExternalChainWatcher)
		case <-ob.stop:
			ob.logger.ExternalChainWatcher.Info().Msg("ExternalChainWatcher stopped")
			return
		}
	}
}

// calcBlockRangeToScan calculates the next range of blocks to scan
func (ob *EVMChainClient) calcBlockRangeToScan(latestConfirmed, lastScanned, batchSize uint64) (uint64, uint64) {
	startBlock := lastScanned + 1
	toBlock := lastScanned + batchSize
	if toBlock > latestConfirmed {
		toBlock = latestConfirmed
	}
	return startBlock, toBlock
}

func (ob *EVMChainClient) postBlockHeader(tip uint64) error {
	bn := tip

	res, err := ob.zetaClient.GetBlockHeaderStateByChain(ob.chain.ChainId)
	if err == nil && res.BlockHeaderState != nil && res.BlockHeaderState.EarliestHeight > 0 {
		// #nosec G701 always positive
		bn = uint64(res.BlockHeaderState.LatestHeight) + 1 // the next header to post
	}

	if bn > tip {
		return fmt.Errorf("postBlockHeader: must post block confirmed block header: %d > %d", bn, tip)
	}

	block, err := ob.GetBlockByNumberCached(bn)
	if err != nil {
		ob.logger.ExternalChainWatcher.Error().Err(err).Msgf("postBlockHeader: error getting block: %d", bn)
		return err
	}
	headerRLP, err := rlp.EncodeToBytes(block.Header())
	if err != nil {
		ob.logger.ExternalChainWatcher.Error().Err(err).Msgf("postBlockHeader: error encoding block header: %d", bn)
		return err
	}

	_, err = ob.zetaClient.PostAddBlockHeader(
		ob.chain.ChainId,
		block.Hash().Bytes(),
		block.Number().Int64(),
		common.NewEthereumHeader(headerRLP),
	)
	if err != nil {
		ob.logger.ExternalChainWatcher.Error().Err(err).Msgf("postBlockHeader: error posting block header: %d", bn)
		return err
	}
	return nil
}

func (ob *EVMChainClient) observeInTX(sampledLogger zerolog.Logger) error {
	// make sure inbound TXS / Send is enabled by the protocol
	flags, err := ob.zetaClient.GetCrosschainFlags()
	if err != nil {
		return err
	}
	if !flags.IsInboundEnabled {
		return errors.New("inbound TXS / Send has been disabled by the protocol")
	}

	// get and update latest block height
	header, err := ob.evmClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}
	ob.SetLastBlockHeight(header.Number.Uint64())

	// increment prom counter
	counter, err := ob.GetPromCounter("rpc_getBlockByNumber_count")
	if err != nil {
		ob.logger.ExternalChainWatcher.Error().Err(err).Msg("GetPromCounter:")
	}
	counter.Inc()

	// skip if current height is too low
	if header.Number.Uint64() < ob.GetChainParams().ConfirmationCount {
		return fmt.Errorf("observeInTX: skipping observer, current block number %d is too low", header.Number.Uint64())
	}
	confirmedBlockNum := header.Number.Uint64() - ob.GetChainParams().ConfirmationCount

	// skip if no new block is confirmed
	lastScanned := ob.GetLastBlockHeightScanned()
	if lastScanned >= confirmedBlockNum {
		sampledLogger.Debug().Msgf("observeInTX: skipping observer, no new block is produced for chain %d", ob.chain.ChainId)
		return nil
	}

	// get last scanned block height (we simply use same height for all 3 events ZetaSent, Deposited, TssRecvd)
	// Note: using different heights for each event incurs more complexity (metrics, db, etc) and not worth it
	startBlock, toBlock := ob.calcBlockRangeToScan(confirmedBlockNum, lastScanned, config.MaxBlocksPerPeriod)

	// task 1:  query evm chain for zeta sent logs (read at most 100 blocks in one go)
	lastScannedZetaSent := ob.observeZetaSent(startBlock, toBlock)

	// task 2: query evm chain for deposited logs (read at most 100 blocks in one go)
	lastScannedDeposited := ob.observeERC20Deposited(startBlock, toBlock)

	// task 3: query the incoming tx to TSS address (read at most 100 blocks in one go)
	lastScannedTssRecvd := ob.observeTssRecvd(startBlock, toBlock, flags)

	// note: using lowest height for all 3 events is not perfect, but it's simple and good enough
	lastScannedLowest := lastScannedZetaSent
	if lastScannedDeposited < lastScannedLowest {
		lastScannedLowest = lastScannedDeposited
	}
	if lastScannedTssRecvd < lastScannedLowest {
		lastScannedLowest = lastScannedTssRecvd
	}

	// update last scanned block height for all 3 events (ZetaSent, Deposited, TssRecvd), ignore db error
	if lastScannedLowest > lastScanned {
		sampledLogger.Info().Msgf("observeInTX: lasstScanned heights for chain %d ZetaSent %d ERC20Deposited %d TssRecvd %d",
			ob.chain.ChainId, lastScannedZetaSent, lastScannedDeposited, lastScannedTssRecvd)
		ob.SetLastBlockHeightScanned(lastScannedLowest)
		if err := ob.db.Save(clienttypes.ToLastBlockSQLType(lastScannedLowest)).Error; err != nil {
			ob.logger.ExternalChainWatcher.Error().Err(err).Msgf("observeInTX: error writing lastScannedLowest %d to db", lastScannedLowest)
		}
	}
	return nil
}

// observeZetaSent queries the ZetaSent event from the connector contract and posts to zetacore
// returns the last block successfully scanned
func (ob *EVMChainClient) observeZetaSent(startBlock, toBlock uint64) uint64 {
	// filter ZetaSent logs
	addrConnector, connector, err := ob.GetConnectorContract()
	if err != nil {
		ob.logger.ChainLogger.Warn().Err(err).Msgf("observeZetaSent: GetConnectorContract error:")
		return startBlock - 1 // lastScanned
	}
	iter, err := connector.FilterZetaSent(&bind.FilterOpts{
		Start:   startBlock,
		End:     &toBlock,
		Context: context.TODO(),
	}, []ethcommon.Address{}, []*big.Int{})
	if err != nil {
		ob.logger.ChainLogger.Warn().Err(err).Msgf(
			"observeZetaSent: FilterZetaSent error from block %d to %d for chain %d", startBlock, toBlock, ob.chain.ChainId)
		return startBlock - 1 // lastScanned
	}

	// collect and sort events by block number, then tx index, then log index (ascending)
	events := make([]*zetaconnector.ZetaConnectorNonEthZetaSent, 0)
	for iter.Next() {
		// sanity check tx event
		err := ob.CheckEvmTxLog(&iter.Event.Raw, addrConnector, "", TopicsZetaSent)
		if err == nil {
			events = append(events, iter.Event)
			continue
		}
		ob.logger.ExternalChainWatcher.Warn().Err(err).Msgf("observeZetaSent: invalid ZetaSent event in tx %s on chain %d at height %d",
			iter.Event.Raw.TxHash.Hex(), ob.chain.ChainId, iter.Event.Raw.BlockNumber)
	}
	sort.SliceStable(events, func(i, j int) bool {
		if events[i].Raw.BlockNumber == events[j].Raw.BlockNumber {
			if events[i].Raw.TxIndex == events[j].Raw.TxIndex {
				return events[i].Raw.Index < events[j].Raw.Index
			}
			return events[i].Raw.TxIndex < events[j].Raw.TxIndex
		}
		return events[i].Raw.BlockNumber < events[j].Raw.BlockNumber
	})

	// increment prom counter
	cnt, err := ob.GetPromCounter("rpc_getFilterLogs_count")
	if err != nil {
		ob.logger.ExternalChainWatcher.Error().Err(err).Msg("GetPromCounter:")
	} else {
		cnt.Inc()
	}

	// post to zetacore
	beingScanned := uint64(0)
	for _, event := range events {
		// remember which block we are scanning (there could be multiple events in the same block)
		if event.Raw.BlockNumber > beingScanned {
			beingScanned = event.Raw.BlockNumber
		}
		msg, err := ob.GetInboundVoteMsgForZetaSentEvent(event)
		if err != nil {
			ob.logger.ExternalChainWatcher.Error().Err(err).Msgf(
				"observeZetaSent: error getting inbound vote msg for tx %s chain %d", event.Raw.TxHash.Hex(), ob.chain.ChainId)
			continue
		}
		zetaHash, ballot, err := ob.zetaClient.PostVoteInbound(PostVoteInboundGasLimit, PostVoteInboundMessagePassingExecutionGasLimit, &msg)
		if err != nil {
			ob.logger.ExternalChainWatcher.Error().Err(err).Msgf(
				"observeZetaSent: error posting event to zeta core for tx %s at height %d for chain %d",
				event.Raw.TxHash.Hex(), event.Raw.BlockNumber, ob.chain.ChainId)
			return beingScanned - 1 // we have to re-scan from this block next time
		} else if zetaHash != "" {
			ob.logger.ExternalChainWatcher.Info().Msgf(
				"observeZetaSent: event detected in tx %s at height %d for chain %d, PostVoteInbound zeta tx: %s ballot %s",
				event.Raw.TxHash.Hex(), event.Raw.BlockNumber, ob.chain.ChainId, zetaHash, ballot)
		}
	}
	// successful processed all events in [startBlock, toBlock]
	return toBlock
}

// observeERC20Deposited queries the ERC20CustodyDeposited event from the ERC20Custody contract and posts to zetacore
// returns the last block successfully scanned
func (ob *EVMChainClient) observeERC20Deposited(startBlock, toBlock uint64) uint64 {
	// filter ERC20CustodyDeposited logs
	addrCustody, erc20custodyContract, err := ob.GetERC20CustodyContract()
	if err != nil {
		ob.logger.ExternalChainWatcher.Warn().Err(err).Msgf("observeERC20Deposited: GetERC20CustodyContract error:")
		return startBlock - 1 // lastScanned
	}
	iter, err := erc20custodyContract.FilterDeposited(&bind.FilterOpts{
		Start:   startBlock,
		End:     &toBlock,
		Context: context.TODO(),
	}, []ethcommon.Address{})
	if err != nil {
		ob.logger.ExternalChainWatcher.Warn().Err(err).Msgf(
			"observeERC20Deposited: FilterDeposited error from block %d to %d for chain %d", startBlock, toBlock, ob.chain.ChainId)
		return startBlock - 1 // lastScanned
	}

	// collect and sort events by block number, then tx index, then log index (ascending)
	events := make([]*erc20custody.ERC20CustodyDeposited, 0)
	for iter.Next() {
		// sanity check tx event
		err := ob.CheckEvmTxLog(&iter.Event.Raw, addrCustody, "", TopicsDeposited)
		if err == nil {
			events = append(events, iter.Event)
			continue
		}
		ob.logger.ExternalChainWatcher.Warn().Err(err).Msgf("observeERC20Deposited: invalid Deposited event in tx %s on chain %d at height %d",
			iter.Event.Raw.TxHash.Hex(), ob.chain.ChainId, iter.Event.Raw.BlockNumber)
	}
	sort.SliceStable(events, func(i, j int) bool {
		if events[i].Raw.BlockNumber == events[j].Raw.BlockNumber {
			if events[i].Raw.TxIndex == events[j].Raw.TxIndex {
				return events[i].Raw.Index < events[j].Raw.Index
			}
			return events[i].Raw.TxIndex < events[j].Raw.TxIndex
		}
		return events[i].Raw.BlockNumber < events[j].Raw.BlockNumber
	})

	// increment prom counter
	cnt, err := ob.GetPromCounter("rpc_getFilterLogs_count")
	if err != nil {
		ob.logger.ExternalChainWatcher.Error().Err(err).Msg("GetPromCounter:")
	} else {
		cnt.Inc()
	}

	// post to zetacore
	beingScanned := uint64(0)
	for _, event := range events {
		// remember which block we are scanning (there could be multiple events in the same block)
		if event.Raw.BlockNumber > beingScanned {
			beingScanned = event.Raw.BlockNumber
		}
		msg, err := ob.GetInboundVoteMsgForDepositedEvent(event)
		if err != nil {
			ob.logger.ExternalChainWatcher.Error().Err(err).Msgf(
				"observeERC20Deposited: error getting inbound vote msg for tx %s chain %d", event.Raw.TxHash.Hex(), ob.chain.ChainId)
			continue
		}
		zetaHash, ballot, err := ob.zetaClient.PostVoteInbound(PostVoteInboundGasLimit, PostVoteInboundExecutionGasLimit, &msg)
		if err != nil {
			ob.logger.ExternalChainWatcher.Error().Err(err).Msgf(
				"observeERC20Deposited: error posting event to zeta core for tx %s at height %d for chain %d",
				event.Raw.TxHash.Hex(), event.Raw.BlockNumber, ob.chain.ChainId)
			return beingScanned - 1 // we have to re-scan from this block next time
		} else if zetaHash != "" {
			ob.logger.ExternalChainWatcher.Info().Msgf(
				"observeERC20Deposited: event detected in tx %s at height %d for chain %d, PostVoteInbound zeta tx: %s ballot %s",
				event.Raw.TxHash.Hex(), event.Raw.BlockNumber, ob.chain.ChainId, zetaHash, ballot)
		}
	}
	// successful processed all events in [startBlock, toBlock]
	return toBlock
}

// observeTssRecvd queries the incoming gas asset to TSS address and posts to zetacore
// returns the last block successfully scanned
func (ob *EVMChainClient) observeTssRecvd(startBlock, toBlock uint64, flags observertypes.CrosschainFlags) uint64 {
	// check TSS address (after keygen, ob.Tss.pubkey will be updated)
	tssAddress := ob.Tss.EVMAddress()
	if tssAddress == (ethcommon.Address{}) {
		ob.logger.ExternalChainWatcher.Warn().Msgf("observeTssRecvd: TSS address not set")
		return startBlock - 1 // lastScanned
	}

	// query incoming gas asset
	for bn := startBlock; bn <= toBlock; bn++ {
		// post new block header (if any) to zetacore and ignore error
		// TODO: consider having a independent ticker(from TSS scaning) for posting block headers
		if flags.BlockHeaderVerificationFlags != nil &&
			flags.BlockHeaderVerificationFlags.IsEthTypeChainEnabled &&
			common.IsHeaderSupportedEvmChain(ob.chain.ChainId) { // post block header for supported chains
			err := ob.postBlockHeader(toBlock)
			if err != nil {
				ob.logger.ExternalChainWatcher.Error().Err(err).Msg("error posting block header")
			}
		}

		// TODO: we can track the total number of 'getBlockByNumber' RPC calls made
		block, err := ob.GetBlockByNumberCached(bn)
		if err != nil {
			ob.logger.ExternalChainWatcher.Error().Err(err).Msgf("observeTssRecvd: error getting block %d for chain %d", bn, ob.chain.ChainId)
			return startBlock - 1 // we have to re-scan from this block next time
		}
		for _, tx := range block.Transactions() {
			if tx.To() == nil {
				continue
			}
			if bytes.Equal(tx.Data(), []byte(DonationMessage)) {
				ob.logger.ExternalChainWatcher.Info().Msgf(
					"observeTssRecvd: thank you rich folk for your donation!: %s chain %d", tx.Hash().Hex(), ob.chain.ChainId)
				continue
			}

			if *tx.To() == tssAddress {
				receipt, err := ob.evmClient.TransactionReceipt(context.Background(), tx.Hash())
				if err != nil {
					ob.logger.ExternalChainWatcher.Err(err).Msgf(
						"observeTssRecvd: TransactionReceipt error for tx %s chain %d", tx.Hash().Hex(), ob.chain.ChainId)
					return startBlock - 1 // we have to re-scan this block next time
				}
				if receipt.Status != 1 { // 1: successful, 0: failed
					ob.logger.ExternalChainWatcher.Info().Msgf("observeTssRecvd: tx %s chain %d failed; don't act", tx.Hash().Hex(), ob.chain.ChainId)
					continue
				}

				from, err := ob.evmClient.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
				if err != nil {
					ob.logger.ExternalChainWatcher.Err(err).Msgf("observeTssRecvd: TransactionSender error for tx %s", tx.Hash().Hex())
					// trying local recovery (assuming LondonSigner dynamic fee tx type) of sender address
					signer := ethtypes.NewLondonSigner(big.NewInt(ob.chain.ChainId))
					from, err = signer.Sender(tx)
					if err != nil {
						ob.logger.ExternalChainWatcher.Err(err).Msgf(
							"observeTssRecvd: local recovery of sender address failed for tx %s chain %d", tx.Hash().Hex(), ob.chain.ChainId)
						continue
					}
				}
				msg := ob.GetInboundVoteMsgForTokenSentToTSS(tx.Hash(), tx.Value(), receipt, from, tx.Data())
				if msg == nil {
					continue
				}
				zetaHash, ballot, err := ob.zetaClient.PostVoteInbound(PostVoteInboundGasLimit, PostVoteInboundExecutionGasLimit, msg)
				if err != nil {
					ob.logger.ExternalChainWatcher.Error().Err(err).Msgf(
						"observeTssRecvd: error posting to zeta core for tx %s at height %d for chain %d", tx.Hash().Hex(), bn, ob.chain.ChainId)
					return startBlock - 1 // we have to re-scan this block next time
				} else if zetaHash != "" {
					ob.logger.ExternalChainWatcher.Info().Msgf(
						"observeTssRecvd: gas asset deposit detected in tx %s at height %d for chain %d, PostVoteInbound zeta tx: %s ballot %s",
						tx.Hash().Hex(), bn, ob.chain.ChainId, zetaHash, ballot)
				}
			}
		}
	}
	// successful processed all gas asset deposits in [startBlock, toBlock]
	return toBlock
}

func (ob *EVMChainClient) WatchGasPrice() {

	err := ob.PostGasPrice()
	if err != nil {
		height, err := ob.zetaClient.GetBlockHeight()
		if err != nil {
			ob.logger.WatchGasPrice.Error().Err(err).Msg("GetBlockHeight error")
		} else {
			ob.logger.WatchGasPrice.Error().Err(err).Msgf("PostGasPrice error at zeta block : %d  ", height)
		}
	}

	ticker, err := NewDynamicTicker(fmt.Sprintf("EVM_WatchGasPrice_%d", ob.chain.ChainId), ob.GetChainParams().GasPriceTicker)
	if err != nil {
		ob.logger.WatchGasPrice.Error().Err(err).Msg("NewDynamicTicker error")
		return
	}

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C():
			err = ob.PostGasPrice()
			if err != nil {
				height, err := ob.zetaClient.GetBlockHeight()
				if err != nil {
					ob.logger.WatchGasPrice.Error().Err(err).Msg("GetBlockHeight error")
				} else {
					ob.logger.WatchGasPrice.Error().Err(err).Msgf("PostGasPrice error at zeta block : %d  ", height)
				}
			}
			ticker.UpdateInterval(ob.GetChainParams().GasPriceTicker, ob.logger.WatchGasPrice)
		case <-ob.stop:
			ob.logger.WatchGasPrice.Info().Msg("WatchGasPrice stopped")
			return
		}
	}
}

func (ob *EVMChainClient) PostGasPrice() error {
	// GAS PRICE
	gasPrice, err := ob.evmClient.SuggestGasPrice(context.TODO())
	if err != nil {
		ob.logger.WatchGasPrice.Err(err).Msg("Err SuggestGasPrice:")
		return err
	}
	blockNum, err := ob.evmClient.BlockNumber(context.TODO())
	if err != nil {
		ob.logger.WatchGasPrice.Err(err).Msg("Err Fetching Most recent Block : ")
		return err
	}

	// SUPPLY
	supply := "100" // lockedAmount on ETH, totalSupply on other chains

	zetaHash, err := ob.zetaClient.PostGasPrice(ob.chain, gasPrice.Uint64(), supply, blockNum)
	if err != nil {
		ob.logger.WatchGasPrice.Err(err).Msg("PostGasPrice to zetacore failed")
		return err
	}
	_ = zetaHash
	//ob.logger.WatchGasPrice.Debug().Msgf("PostGasPrice zeta tx: %s", zetaHash)

	return nil
}

func (ob *EVMChainClient) BuildLastBlock() error {
	logger := ob.logger.ChainLogger.With().Str("module", "BuildBlockIndex").Logger()
	envvar := ob.chain.ChainName.String() + "_SCAN_FROM"
	scanFromBlock := os.Getenv(envvar)
	if scanFromBlock != "" {
		logger.Info().Msgf("BuildLastBlock: envvar %s is set; scan from  block %s", envvar, scanFromBlock)
		if scanFromBlock == clienttypes.EnvVarLatest {
			header, err := ob.evmClient.HeaderByNumber(context.Background(), nil)
			if err != nil {
				return err
			}
			ob.SetLastBlockHeightScanned(header.Number.Uint64())
		} else {
			scanFromBlockInt, err := strconv.ParseUint(scanFromBlock, 10, 64)
			if err != nil {
				return err
			}
			ob.SetLastBlockHeightScanned(scanFromBlockInt)
		}
	} else { // last observed block
		var lastBlockNum clienttypes.LastBlockSQLType
		if err := ob.db.First(&lastBlockNum, clienttypes.LastBlockNumID).Error; err != nil {
			logger.Info().Msgf("BuildLastBlock: db PosKey does not exist; read from external chain %s", ob.chain.String())
			header, err := ob.evmClient.HeaderByNumber(context.Background(), nil)
			if err != nil {
				return err
			}
			ob.SetLastBlockHeightScanned(header.Number.Uint64())
			if dbc := ob.db.Save(clienttypes.ToLastBlockSQLType(ob.GetLastBlockHeightScanned())); dbc.Error != nil {
				logger.Error().Err(dbc.Error).Msgf("BuildLastBlock: error writing lastBlockScanned %d to db", ob.GetLastBlockHeightScanned())
			}
		} else {
			ob.SetLastBlockHeightScanned(lastBlockNum.Num)
		}
	}
	return nil
}

func (ob *EVMChainClient) BuildReceiptsMap() error {
	logger := ob.logger
	var receipts []clienttypes.ReceiptSQLType
	if err := ob.db.Find(&receipts).Error; err != nil {
		logger.ChainLogger.Error().Err(err).Msg("error iterating over db")
		return err
	}
	for _, receipt := range receipts {
		r, err := clienttypes.FromReceiptDBType(receipt.Receipt)
		if err != nil {
			return err
		}
		ob.outTXConfirmedReceipts[receipt.Identifier] = r
	}

	return nil
}

func (ob *EVMChainClient) BuildTransactionsMap() error {
	logger := ob.logger
	var transactions []clienttypes.TransactionSQLType
	if err := ob.db.Find(&transactions).Error; err != nil {
		logger.ChainLogger.Error().Err(err).Msg("error iterating over db")
		return err
	}
	for _, transaction := range transactions {
		trans, err := clienttypes.FromTransactionDBType(transaction.Transaction)
		if err != nil {
			return err
		}
		ob.outTXConfirmedTransactions[transaction.Identifier] = trans
	}
	return nil
}

// LoadDB open sql database and load data into EVMChainClient
func (ob *EVMChainClient) LoadDB(dbPath string, chain common.Chain) error {
	if dbPath != "" {
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			err := os.MkdirAll(dbPath, os.ModePerm)
			if err != nil {
				return err
			}
		}
		path := fmt.Sprintf("%s/%s", dbPath, chain.ChainName.String()) //Use "file::memory:?cache=shared" for temp db
		db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		err = db.AutoMigrate(&clienttypes.ReceiptSQLType{},
			&clienttypes.TransactionSQLType{},
			&clienttypes.LastBlockSQLType{})
		if err != nil {
			ob.logger.ChainLogger.Error().Err(err).Msg("error migrating db")
			return err
		}

		ob.db = db
		err = ob.BuildLastBlock()
		if err != nil {
			return err
		}
	}
	return nil
}

func (ob *EVMChainClient) SetMinAndMaxNonce(trackers []types.OutTxTracker) error {
	minNonce, maxNonce := int64(-1), int64(0)
	for _, tracker := range trackers {
		conv := tracker.Nonce
		// #nosec G701 always in range
		intNonce := int64(conv)
		if minNonce == -1 {
			minNonce = intNonce
		}
		if intNonce < minNonce {
			minNonce = intNonce
		}
		if intNonce > maxNonce {
			maxNonce = intNonce
		}
	}
	if minNonce != -1 {
		atomic.StoreInt64(&ob.MinNonce, minNonce)
	}
	if maxNonce > 0 {
		atomic.StoreInt64(&ob.MaxNonce, maxNonce)
	}
	return nil
}

func (ob *EVMChainClient) GetTxID(nonce uint64) string {
	tssAddr := ob.Tss.EVMAddress().String()
	return fmt.Sprintf("%d-%s-%d", ob.chain.ChainId, tssAddr, nonce)
}

func (ob *EVMChainClient) GetBlockByNumberCached(blockNumber uint64) (*ethtypes.Block, error) {
	if block, ok := ob.BlockCache.Get(blockNumber); ok {
		return block.(*ethtypes.Block), nil
	}
	block, err := ob.evmClient.BlockByNumber(context.Background(), new(big.Int).SetUint64(blockNumber))
	if err != nil {
		return nil, err
	}
	ob.BlockCache.Add(blockNumber, block)
	ob.BlockCache.Add(block.Hash(), block)
	return block, nil
}
