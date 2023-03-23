package zetaclient

import (
	"context"
	"cosmossdk.io/math"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	math2 "math"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/zetacore/contracts/evm"
	"github.com/zeta-chain/zetacore/contracts/evm/erc20custody"
	metricsPkg "github.com/zeta-chain/zetacore/zetaclient/metrics"

	"github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"

	cctxtypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
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

// Chain configuration struct
// Filled with above constants depending on chain
type EVMChainClient struct {
	*ChainMetrics

	chain                     common.Chain
	endpoint                  string
	ticker                    *time.Ticker
	Connector                 *evm.Connector
	ConnectorAddress          ethcommon.Address
	ERC20Custody              *erc20custody.ERC20Custody
	ERC20CustodyAddress       ethcommon.Address
	EvmClient                 *ethclient.Client
	KlaytnClient              *KlaytnClient
	zetaClient                *ZetaCoreBridge
	Tss                       TSSSigner
	lastBlock                 int64
	confCount                 uint64 // must wait this many blocks to be considered "confirmed"
	BlockTime                 uint64 // block time in seconds
	txWatchList               map[ethcommon.Hash]string
	mu                        *sync.Mutex
	db                        *gorm.DB
	sampleLogger              *zerolog.Logger
	outTXConfirmedReceipts    map[int]*ethtypes.Receipt
	outTXConfirmedTransaction map[int]*ethtypes.Transaction
	MinNonce                  int64
	MaxNonce                  int64
	OutTxChan                 chan OutTx // send to this channel if you want something back!
	stop                      chan struct{}
	fileLogger                *zerolog.Logger // for critical info
	logger                    zerolog.Logger
}

var _ ChainClient = (*EVMChainClient)(nil)

// Return configuration based on supplied target chain
func NewEVMChainClient(chain common.Chain, bridge *ZetaCoreBridge, tss TSSSigner, dbpath string, metrics *metricsPkg.Metrics) (*EVMChainClient, error) {
	ob := EVMChainClient{
		ChainMetrics: NewChainMetrics(chain.ChainName.String(), metrics),
	}
	ob.stop = make(chan struct{})
	ob.chain = chain
	ob.mu = &sync.Mutex{}
	sampled := log.Sample(&zerolog.BasicSampler{N: 10})
	ob.sampleLogger = &sampled
	ob.logger = log.With().Str("chain", chain.ChainName.String()).Logger()
	ob.zetaClient = bridge
	ob.txWatchList = make(map[ethcommon.Hash]string)
	ob.Tss = tss
	ob.outTXConfirmedReceipts = make(map[int]*ethtypes.Receipt)
	ob.outTXConfirmedTransaction = make(map[int]*ethtypes.Transaction)
	ob.OutTxChan = make(chan OutTx, 100)
	addr := ethcommon.HexToAddress(config.ChainConfigs[chain.ChainName.String()].ConnectorContractAddress)
	erc20CustodyAddress := ethcommon.HexToAddress(config.ChainConfigs[chain.ChainName.String()].ERC20CustodyContractAddress)
	if addr == ethcommon.HexToAddress("0x0") {
		return nil, fmt.Errorf("connector contract address %s not configured for chain %s", config.ChainConfigs[chain.String()].ConnectorContractAddress, chain.String())
	}
	ob.ConnectorAddress = addr
	ob.ERC20CustodyAddress = erc20CustodyAddress
	ob.endpoint = config.ChainConfigs[chain.ChainName.String()].Endpoint
	logFile, err := os.OpenFile(ob.chain.ChainName.String()+"_debug.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		// Can we log an error before we have our Logger? :)
		log.Error().Err(err).Msgf("there was an error creating a logFile chain %s", ob.chain.String())
	}
	fileLogger := zerolog.New(logFile).With().Logger()
	ob.fileLogger = &fileLogger

	// initialize the Client
	ob.logger.Info().Msgf("Chain %s endpoint %s", ob.chain.String(), ob.endpoint)
	client, err := ethclient.Dial(ob.endpoint)
	if err != nil {
		ob.logger.Error().Err(err).Msg("eth Client Dial")
		return nil, err
	}
	ob.EvmClient = client

	if chain.IsKlaytnChain() {
		kclient, err := Dial(ob.endpoint)
		if err != nil {
			ob.logger.Error().Err(err).Msg("klaytn Client Dial")
			return nil, err
		}
		ob.KlaytnClient = kclient
	}

	// initialize the connector
	connector, err := evm.NewConnector(addr, ob.EvmClient)
	if err != nil {
		ob.logger.Error().Err(err).Msg("Connector")
		return nil, err
	}
	ob.Connector = connector

	// initialize erc20 custody
	erc20CustodyContract, err := erc20custody.NewERC20Custody(erc20CustodyAddress, ob.EvmClient)
	if err != nil {
		ob.logger.Error().Err(err).Msg("ERC20Custody")
		return nil, err
	}
	ob.ERC20Custody = erc20CustodyContract

	// create metric counters
	err = ob.RegisterPromCounter("rpc_getLogs_count", "Number of getLogs")
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

	ob.SetChainDetails(chain)

	err = ob.LoadDB(dbpath, chain)
	if err != nil {
		return nil, err
	}

	ob.logger.Info().Msgf("%s: start scanning from block %d", chain.String(), ob.GetLastBlockHeight())

	return &ob, nil
}

func (ob *EVMChainClient) Start() {
	go ob.ExternalChainWatcher() // Observes external Chains for incoming trasnactions
	go ob.WatchGasPrice()        // Observes external Chains for Gas prices and posts to core
	go ob.observeOutTx()         // Populates receipts and confirmed outbound transactions
}

func (ob *EVMChainClient) Stop() {
	ob.logger.Info().Msgf("ob %s is stopping", ob.chain.String())
	close(ob.stop) // this notifies all goroutines to stop

	ob.logger.Info().Msg("closing ob.db")
	dbInst, err := ob.db.DB()
	if err != nil {
		ob.logger.Error().Err(err).Msg("error getting database instance")
	}
	err = dbInst.Close()
	if err != nil {
		ob.logger.Error().Err(err).Msg("error closing database")
	}

	ob.logger.Info().Msgf("%s observer stopped", ob.chain.String())
}

// returns: isIncluded, isConfirmed, Error
// If isConfirmed, it also post to ZetaCore
func (ob *EVMChainClient) IsSendOutTxProcessed(sendHash string, nonce int, cointype common.CoinType) (bool, bool, error) {
	ob.mu.Lock()
	receipt, found1 := ob.outTXConfirmedReceipts[nonce]
	transaction, found2 := ob.outTXConfirmedTransaction[nonce]
	ob.mu.Unlock()
	found := found1 && found2
	if !found {
		return false, false, nil
	}
	sendID := fmt.Sprintf("%s-%d", ob.chain.String(), nonce)
	logger := ob.logger.With().Str("sendID", sendID).Logger()
	if cointype == common.CoinType_Gas { // the outbound is a regular Ether/BNB/Matic transfer; no need to check events
		if receipt.Status == 1 {
			zetaHash, err := ob.zetaClient.PostReceiveConfirmation(
				sendHash,
				receipt.TxHash.Hex(),
				receipt.BlockNumber.Uint64(),
				transaction.Value(),
				common.ReceiveStatus_Success,
				ob.chain,
				nonce,
				common.CoinType_Gas,
			)
			if err != nil {
				logger.Error().Err(err).Msg("error posting confirmation to meta core")
			}
			logger.Info().Msgf("Zeta tx hash: %s\n", zetaHash)
			return true, true, nil
		} else if receipt.Status == 0 { // the same as below events flow
			logger.Info().Msgf("Found (failed tx) sendHash %s on chain %s txhash %s", sendHash, ob.chain.String(), receipt.TxHash.Hex())
			zetaTxHash, err := ob.zetaClient.PostReceiveConfirmation(sendHash, receipt.TxHash.Hex(), receipt.BlockNumber.Uint64(), big.NewInt(0), common.ReceiveStatus_Failed, ob.chain, nonce, common.CoinType_Gas)
			if err != nil {
				logger.Error().Err(err).Msgf("PostReceiveConfirmation error in WatchTxHashWithTimeout; zeta tx hash %s", zetaTxHash)
			}
			logger.Info().Msgf("Zeta tx hash: %s", zetaTxHash)
			return true, true, nil
		}
	} else if cointype == common.CoinType_Zeta { // the outbound is a Zeta transfer; need to check events ZetaReceived
		if receipt.Status == 1 {
			logs := receipt.Logs
			for _, vLog := range logs {
				confHeight := vLog.BlockNumber + ob.confCount
				if confHeight < 0 || confHeight >= math2.MaxInt64 {
					return false, false, fmt.Errorf("confHeight is out of range")
				}
				// TODO rewrite this to return early if not confirmed
				receivedLog, err := ob.Connector.ConnectorFilterer.ParseZetaReceived(*vLog)
				if err == nil {
					logger.Info().Msgf("Found (outTx) sendHash %s on chain %s txhash %s", sendHash, ob.chain.String(), vLog.TxHash.Hex())
					if int64(confHeight) < ob.GetLastBlockHeight() {
						logger.Info().Msg("Confirmed! Sending PostConfirmation to zetacore...")
						if len(vLog.Topics) != 4 {
							logger.Error().Msgf("wrong number of topics in log %d", len(vLog.Topics))
							return false, false, fmt.Errorf("wrong number of topics in log %d", len(vLog.Topics))
						}
						sendhash := vLog.Topics[3].Hex()
						//var rxAddress string = ethcommon.HexToAddress(vLog.Topics[1].Hex()).Hex()
						mMint := receivedLog.ZetaValue
						zetaHash, err := ob.zetaClient.PostReceiveConfirmation(
							sendhash,
							vLog.TxHash.Hex(),
							vLog.BlockNumber,
							mMint,
							common.ReceiveStatus_Success,
							ob.chain,
							nonce,
							common.CoinType_Zeta,
						)
						if err != nil {
							logger.Error().Err(err).Msg("error posting confirmation to meta core")
							continue
						}
						logger.Info().Msgf("Zeta tx hash: %s\n", zetaHash)
						return true, true, nil
					}
					logger.Info().Msgf("Included; %d blocks before confirmed! chain %s nonce %d", int(vLog.BlockNumber+ob.confCount)-int(ob.GetLastBlockHeight()), ob.chain.String(), nonce)
					return true, false, nil
				}
				revertedLog, err := ob.Connector.ConnectorFilterer.ParseZetaReverted(*vLog)
				if err == nil {
					logger.Info().Msgf("Found (revertTx) sendHash %s on chain %s txhash %s", sendHash, ob.chain.String(), vLog.TxHash.Hex())
					if int64(confHeight) < ob.GetLastBlockHeight() {
						logger.Info().Msg("Confirmed! Sending PostConfirmation to zetacore...")
						if len(vLog.Topics) != 3 {
							logger.Error().Msgf("wrong number of topics in log %d", len(vLog.Topics))
							return false, false, fmt.Errorf("wrong number of topics in log %d", len(vLog.Topics))
						}
						sendhash := vLog.Topics[2].Hex()
						mMint := revertedLog.RemainingZetaValue
						metaHash, err := ob.zetaClient.PostReceiveConfirmation(
							sendhash,
							vLog.TxHash.Hex(),
							vLog.BlockNumber,
							mMint,
							common.ReceiveStatus_Success,
							ob.chain,
							nonce,
							common.CoinType_Zeta,
						)
						if err != nil {
							logger.Err(err).Msg("error posting confirmation to meta core")
							continue
						}
						logger.Info().Msgf("Zeta tx hash: %s", metaHash)
						return true, true, nil
					}
					logger.Info().Msgf("Included; %d blocks before confirmed! chain %s nonce %d", int(vLog.BlockNumber+ob.confCount)-int(ob.GetLastBlockHeight()), ob.chain.String(), nonce)
					return true, false, nil
				}
			}
		} else if receipt.Status == 0 {
			//FIXME: check nonce here by getTransaction RPC
			logger.Info().Msgf("Found (failed tx) sendHash %s on chain %s txhash %s", sendHash, ob.chain.String(), receipt.TxHash.Hex())
			zetaTxHash, err := ob.zetaClient.PostReceiveConfirmation(sendHash, receipt.TxHash.Hex(), receipt.BlockNumber.Uint64(), big.NewInt(0), common.ReceiveStatus_Failed, ob.chain, nonce, common.CoinType_Zeta)
			if err != nil {
				logger.Error().Err(err).Msgf("PostReceiveConfirmation error in WatchTxHashWithTimeout; zeta tx hash %s", zetaTxHash)
			}
			logger.Info().Msgf("Zeta tx hash: %s", zetaTxHash)
			return true, true, nil
		}
	} else if cointype == common.CoinType_ERC20 {
		if receipt.Status == 1 {
			logs := receipt.Logs
			ERC20Custody, err := erc20custody.NewERC20Custody(ob.ERC20CustodyAddress, ob.EvmClient)
			if err != nil {
				logger.Warn().Msgf("NewERC20Custody err: %s", err)
			}
			for _, vLog := range logs {
				event, err := ERC20Custody.ParseWithdrawn(*vLog)
				confHeight := vLog.BlockNumber + ob.confCount
				if confHeight < 0 || confHeight >= math2.MaxInt64 {
					return false, false, fmt.Errorf("confHeight is out of range")
				}
				if err == nil {
					logger.Info().Msgf("Found (ERC20Custody.Withdrawn Event) sendHash %s on chain %s txhash %s", sendHash, ob.chain.String(), vLog.TxHash.Hex())
					if int64(confHeight) < ob.GetLastBlockHeight() {

						logger.Info().Msg("Confirmed! Sending PostConfirmation to zetacore...")
						zetaHash, err := ob.zetaClient.PostReceiveConfirmation(
							sendHash,
							vLog.TxHash.Hex(),
							vLog.BlockNumber,
							event.Amount,
							common.ReceiveStatus_Success,
							ob.chain,
							nonce,
							common.CoinType_ERC20,
						)
						if err != nil {
							logger.Error().Err(err).Msg("error posting confirmation to meta core")
							continue
						}
						logger.Info().Msgf("Zeta tx hash: %s\n", zetaHash)
						return true, true, nil
					}
					logger.Info().Msgf("Included; %d blocks before confirmed! chain %s nonce %d", int(vLog.BlockNumber+ob.confCount)-int(ob.GetLastBlockHeight()), ob.chain.String(), nonce)
					return true, false, nil
				}
			}
		}
	}

	return false, false, nil
}

// FIXME: there's a chance that a txhash in OutTxChan may not deliver when Stop() is called
// observeOutTx periodically checks all the txhash in potential outbound txs
func (ob *EVMChainClient) observeOutTx() {
	logger := ob.logger
	ticker := time.NewTicker(3 * time.Second) // FIXME: config this
	for {
		select {
		case <-ticker.C:
			trackers, err := ob.zetaClient.GetAllOutTxTrackerByChain(ob.chain)
			if err != nil {
				continue
			}
			outTimeout := time.After(90 * time.Second)
		TRACKERLOOP:
			for _, tracker := range trackers {
				nonceInt := tracker.Nonce
			TXHASHLOOP:
				for _, txHash := range tracker.HashList {
					inTimeout := time.After(3000 * time.Millisecond)
					select {
					case <-outTimeout:
						logger.Warn().Msgf("observeOutTx timeout on nonce %d", nonceInt)
						break TRACKERLOOP
					default:
						receipt, transaction, err := ob.queryTxByHash(txHash.TxHash, int64(nonceInt))
						if err == nil && receipt != nil { // confirmed
							ob.mu.Lock()
							ob.outTXConfirmedReceipts[int(nonceInt)] = receipt
							ob.outTXConfirmedTransaction[int(nonceInt)] = transaction
							ob.mu.Unlock()

							if dbc := ob.db.Create(clienttypes.ToReceiptSQLType(receipt, int(nonceInt))); dbc.Error != nil {
								logger.Error().Err(err).Msgf("PurgeTxHashWatchList: error putting nonce %d tx hashes %s to db", nonceInt, receipt.TxHash.Hex())
							}
							if dbc := ob.db.Create(clienttypes.ToTransactionSQLType(transaction, int(nonceInt))); dbc.Error != nil {
								logger.Error().Err(err).Msgf("PurgeTxHashWatchList: error putting nonce %d tx hashes %s to db", nonceInt, transaction.Hash())
							}

							break TXHASHLOOP
						}
						<-inTimeout
					}
				}
			}
		case <-ob.stop:
			logger.Info().Msg("observeOutTx: stopped")
			return
		}
	}
}

// return the status of txHash
// receipt nil, err non-nil: txHash not found
// receipt nil, err nil: txHash receipt recorded, but may not be confirmed
// receipt non-nil, err nil: txHash confirmed
func (ob *EVMChainClient) queryTxByHash(txHash string, nonce int64) (*ethtypes.Receipt, *ethtypes.Transaction, error) {
	logger := ob.logger.With().Str("txHash", txHash).Int64("nonce", nonce).Logger()
	if ob.outTXConfirmedReceipts[int(nonce)] != nil && ob.outTXConfirmedTransaction[int(nonce)] != nil {
		return nil, nil, fmt.Errorf("queryTxByHash: txHash %s receipts already recorded", txHash)
	}
	ctxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	receipt, err := ob.EvmClient.TransactionReceipt(ctxt, ethcommon.HexToHash(txHash))
	if err != nil {
		if err != ethereum.NotFound {
			logger.Warn().Err(err).Msg("TransactionReceipt/TransactionByHash error")
		}
		return nil, nil, err
	}
	transaction, _, err := ob.EvmClient.TransactionByHash(ctxt, ethcommon.HexToHash(txHash))
	if err != nil {
		return nil, nil, err
	}
	confHeight := receipt.BlockNumber.Uint64() + ob.confCount
	if confHeight < 0 || confHeight >= math2.MaxInt64 {
		return nil, nil, fmt.Errorf("confHeight is out of range")
	}
	if int64(confHeight) > ob.GetLastBlockHeight() {
		log.Warn().Msgf("included but not confirmed: receipt block %d, current block %d", receipt.BlockNumber, ob.GetLastBlockHeight())
		return nil, nil, fmt.Errorf("included but not confirmed")
	}
	return receipt, transaction, nil
}

func (ob *EVMChainClient) SetLastBlockHeight(block int64) {
	if block < 0 {
		panic("lastBlock is negative")
	}
	if block >= math2.MaxInt64 {
		panic("lastBlock is too large")
	}
	atomic.StoreInt64(&ob.lastBlock, block)
}

func (ob *EVMChainClient) GetLastBlockHeight() int64 {
	height := atomic.LoadInt64(&ob.lastBlock)
	if height < 0 {
		panic("lastBlock is negative")
	}
	if height >= math2.MaxInt64 {
		panic("lastBlock is too large")
	}
	return height
}

func (ob *EVMChainClient) ExternalChainWatcher() {
	// At each tick, query the Connector contract
	ob.logger.Info().Msg("ExternalChainWatcher started")
	for {
		select {
		case <-ob.ticker.C:
			err := ob.observeInTX()
			if err != nil {
				ob.logger.Err(err).Msg("observeInTX error")
				continue
			}
		case <-ob.stop:
			ob.logger.Info().Msg("ExternalChainWatcher stopped")
			return
		}
	}
}

func (ob *EVMChainClient) observeInTX() error {
	permssions, err := ob.zetaClient.GetInboundPermissions()
	if err != nil {
		return err
	}
	if !permssions.IsInboundEnabled {
		return errors.New("inbound TXS / Send has been disabled by the protocol")
	}
	header, err := ob.EvmClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}
	counter, err := ob.GetPromCounter("rpc_getBlockByNumber_count")
	if err != nil {
		ob.logger.Error().Err(err).Msg("GetPromCounter:")
	}
	counter.Inc()

	// "confirmed" current block number
	confirmedBlockNum := header.Number.Uint64() - ob.confCount
	// skip if no new block is produced.
	if confirmedBlockNum < 0 || confirmedBlockNum > math2.MaxUint64 {
		ob.sampleLogger.Error().Msg("Skipping observer , confirmedBlockNum is negative or too large ")
		return nil
	}
	if confirmedBlockNum <= uint64(ob.GetLastBlockHeight()) {
		ob.sampleLogger.Info().Msg("Skipping observer , No new block is produced ")
		return nil
	}
	lastBlock := ob.GetLastBlockHeight()
	startBlock := lastBlock + 1
	toBlock := lastBlock + config.MaxBlocksPerPeriod // read at most 10 blocks in one go
	if uint64(toBlock) >= confirmedBlockNum {
		toBlock = int64(confirmedBlockNum)
	}
	//ob.logger.Info().Msgf("%s current block %d, querying from %d to %d, %d blocks left to catch up, watching MPI address %s", ob.chain.String(), header.Number.Uint64(), ob.GetLastBlockHeight()+1, toBlock, int(toBlock)-int(confirmedBlockNum), ob.ConnectorAddress.Hex())

	// Query evm chain for zeta sent logs
	if startBlock < 0 || startBlock >= math2.MaxInt64 {
		return fmt.Errorf("startBlock is negative or too large")
	}
	if toBlock < 0 || toBlock >= math2.MaxInt64 {
		return fmt.Errorf("toBlock is negative or too large")
	}
	tb := uint64(toBlock)
	logs, err := ob.Connector.FilterZetaSent(&bind.FilterOpts{
		Start:   uint64(startBlock),
		End:     &tb,
		Context: context.TODO(),
	}, []ethcommon.Address{}, []*big.Int{})

	if err != nil {
		return err
	}
	cnt, err := ob.GetPromCounter("rpc_getLogs_count")
	if err != nil {
		return err
	}
	cnt.Inc()

	// Pull out arguments from logs
	for logs.Next() {
		event := logs.Event
		ob.logger.Info().Msgf("TxBlockNumber %d Transaction Hash: %s Message : %s", event.Raw.BlockNumber, event.Raw.TxHash, event.Message)
		destChain := common.GetChainFromChainID(event.DestinationChainId.Int64())
		destAddr := clienttypes.BytesToEthHex(event.DestinationAddress)

		if strings.EqualFold(destAddr, config.ChainConfigs[destChain.ChainName.String()].ZETATokenContractAddress) {
			ob.logger.Warn().Msgf("potential attack attempt: %s destination address is ZETA token contract address %s", destChain, destAddr)
		}
		zetaHash, err := ob.zetaClient.PostSend(
			event.ZetaTxSenderAddress.Hex(),
			ob.chain.ChainId,
			event.SourceTxOriginAddress.Hex(),
			clienttypes.BytesToEthHex(event.DestinationAddress),
			destChain.ChainId,
			math.NewUintFromBigInt(event.ZetaValueAndGas),
			base64.StdEncoding.EncodeToString(event.Message),
			event.Raw.TxHash.Hex(),
			event.Raw.BlockNumber,
			event.DestinationGasLimit.Uint64(),
			common.CoinType_Zeta,
			PostSendNonEVMGasLimit,
			"",
		)
		if err != nil {
			ob.logger.Error().Err(err).Msg("error posting to zeta core")
			continue
		}
		ob.logger.Info().Msgf("ZetaSent event detected and reported: PostSend zeta tx: %s", zetaHash)
	}

	// Query evm chain for deposited logs
	if startBlock < 0 || startBlock >= math2.MaxInt64 {
		ob.logger.Error().Msgf("startBlock is out of range: %d", startBlock)
	}
	if toBlock < 0 || toBlock >= math2.MaxInt64 {
		ob.logger.Error().Msgf("toBlock is out of range: %d", toBlock)
	}
	toB := uint64(toBlock)
	depositedLogs, err := ob.ERC20Custody.FilterDeposited(&bind.FilterOpts{
		Start:   uint64(startBlock),
		End:     &toB,
		Context: context.TODO(),
	})

	if err != nil {
		return err
	}
	cnt, err = ob.GetPromCounter("rpc_getLogs_count")
	if err != nil {
		return err
	}
	cnt.Inc()

	// Pull out arguments from logs
	for depositedLogs.Next() {
		event := depositedLogs.Event
		ob.logger.Info().Msgf("TxBlockNumber %d Transaction Hash: %s Message : %s", event.Raw.BlockNumber, event.Raw.TxHash, event.Message)

		zetaHash, err := ob.zetaClient.PostSend(
			"",
			ob.chain.ChainId,
			"",
			clienttypes.BytesToEthHex(event.Recipient),
			config.ChainConfigs[common.ZetaChain().ChainName.String()].Chain.ChainId,
			math.NewUintFromBigInt(event.Amount),
			hex.EncodeToString(event.Message),
			event.Raw.TxHash.Hex(),
			event.Raw.BlockNumber,
			1_500_000,
			common.CoinType_ERC20,
			PostSendEVMGasLimit,
			event.Asset.String(),
		)
		if err != nil {
			ob.logger.Error().Err(err).Msg("error posting to zeta core")
			continue
		}
		ob.logger.Info().Msgf("ZRC20Cusotdy Deposited event detected and reported: PostSend zeta tx: %s", zetaHash)
	}

	// ============= query the incoming tx to TSS address ==============
	tssAddress := ob.Tss.EVMAddress()
	// query incoming gas asset
	if !ob.chain.IsKlaytnChain() {
		for bn := startBlock; bn <= toBlock; bn++ {
			//block, err := ob.EvmClient.BlockByNumber(context.Background(), big.NewInt(int64(bn)))
			block, err := ob.EvmClient.BlockByNumber(context.Background(), big.NewInt(bn))
			if err != nil {
				ob.logger.Error().Err(err).Msgf("error getting block: %d", bn)
				continue
			}
			ob.logger.Debug().Msgf("block %d: num txs: %d", bn, len(block.Transactions()))
			for _, tx := range block.Transactions() {
				if tx.To() == nil {
					continue
				}
				if *tx.To() == tssAddress {
					receipt, err := ob.EvmClient.TransactionReceipt(context.Background(), tx.Hash())
					if err != nil {
						ob.logger.Err(err).Msg("TransactionReceipt error")
						continue
					}
					if receipt.Status != 1 { // 1: successful, 0: failed
						ob.logger.Info().Msgf("tx %s failed; don't act", tx.Hash().Hex())
						continue
					}

					from, err := ob.EvmClient.TransactionSender(context.Background(), tx, block.Hash(), receipt.TransactionIndex)
					if err != nil {
						ob.logger.Err(err).Msg("TransactionSender error; trying local recovery (assuming LondonSigner dynamic fee tx type) of sender address")
						chainConf, found := config.ChainConfigs[ob.chain.String()]
						if !found || chainConf == nil {
							ob.logger.Error().Msgf("chain %s not found in config", ob.chain.String())
							continue
						}
						signer := ethtypes.NewLondonSigner(big.NewInt(chainConf.Chain.ChainId))
						from, err = signer.Sender(tx)
						if err != nil {
							ob.logger.Err(err).Msg("local recovery of sender address failed")
							continue
						}
					}
					zetaHash, err := ob.ReportTokenSentToTSS(tx.Hash(), tx.Value(), receipt, from, tx.Data())
					if err != nil {
						ob.logger.Error().Err(err).Msg("error posting to zeta core")
						continue
					}
					ob.logger.Info().Msgf("Gas Deposit detected and reported: PostSend zeta tx: %s", zetaHash)
				}
			}
		}
	} else { // for Klaytn
		for bn := startBlock; bn <= toBlock; bn++ {
			//block, err := ob.EvmClient.BlockByNumber(context.Background(), big.NewInt(int64(bn)))
			block, err := ob.KlaytnClient.BlockByNumber(context.Background(), big.NewInt(bn))
			if err != nil {
				ob.logger.Error().Err(err).Msgf("error getting block: %d", bn)
				continue
			}
			for _, tx := range block.Transactions {
				if tx.To == nil {
					continue
				}
				if *tx.To == tssAddress {
					receipt, err := ob.EvmClient.TransactionReceipt(context.Background(), tx.Hash)
					if err != nil {
						ob.logger.Err(err).Msg("TransactionReceipt error")
						continue
					}
					if receipt.Status != 1 { // 1: successful, 0: failed
						ob.logger.Info().Msgf("tx %s failed; don't act", tx.Hash.Hex())
						continue
					}

					from := *tx.From
					value := tx.Value.ToInt()

					zetaHash, err := ob.ReportTokenSentToTSS(tx.Hash, value, receipt, from, tx.Input)
					if err != nil {
						ob.logger.Error().Err(err).Msg("error posting to zeta core")
						continue
					}
					ob.logger.Info().Msgf("ZetaSent event detected and reported: PostSend zeta tx: %s", zetaHash)
				}
			}
		}
	}
	// ============= end of query the incoming tx to TSS address ==============
	ob.SetLastBlockHeight(toBlock)
	if dbc := ob.db.Create(clienttypes.ToLastBlockSQLType(ob.GetLastBlockHeight())); dbc.Error != nil {
		ob.logger.Error().Err(err).Msg("error writing toBlock to db")
	}
	return nil
}

func (ob *EVMChainClient) ReportTokenSentToTSS(txhash ethcommon.Hash, value *big.Int, receipt *ethtypes.Receipt, from ethcommon.Address, data []byte) (string, error) {
	ob.logger.Info().Msgf("TSS inTx detected: %s, blocknum %d", txhash.Hex(), receipt.BlockNumber)
	ob.logger.Info().Msgf("TSS inTx value: %s", value.String())
	ob.logger.Info().Msgf("TSS inTx from: %s", from.Hex())
	message := ""
	if len(data) != 0 {
		message = hex.EncodeToString(data)
	}
	zetaHash, err := ob.zetaClient.PostSend(
		from.Hex(),
		ob.chain.ChainId,
		from.Hex(),
		from.Hex(),
		common.ZetaChain().ChainId,
		math.NewUintFromBigInt(value),
		message,
		txhash.Hex(),
		receipt.BlockNumber.Uint64(),
		90_000,
		common.CoinType_Gas,
		PostSendEVMGasLimit,
		"",
	)
	return zetaHash, err
}

// query the base gas price for the block number bn.
func (ob *EVMChainClient) GetBaseGasPrice() *big.Int {
	gasPrice, err := ob.EvmClient.SuggestGasPrice(context.TODO())
	if err != nil {
		ob.logger.Err(err).Msg("GetBaseGasPrice")
		return nil
	}
	return gasPrice
}

func (ob *EVMChainClient) PostNonceIfNotRecorded() error {
	logger := ob.logger
	zetaClient := ob.zetaClient
	evmClient := ob.EvmClient
	tss := ob.Tss
	chain := ob.chain

	_, err := zetaClient.GetNonceByChain(chain)
	if err != nil { // if Nonce of Chain is not found in ZetaCore; report it
		nonce, err := evmClient.NonceAt(context.TODO(), tss.EVMAddress(), nil)
		if err != nil {
			logger.Fatal().Err(err).Msg("NonceAt")
			return err
		}
		pendingNonce, err := evmClient.PendingNonceAt(context.TODO(), tss.EVMAddress())
		if err != nil {
			logger.Fatal().Err(err).Msg("PendingNonceAt")
			return err
		}
		if pendingNonce != nonce {
			logger.Fatal().Msgf("fatal: pending nonce %d != nonce %d", pendingNonce, nonce)
			return fmt.Errorf("pending nonce %d != nonce %d", pendingNonce, nonce)
		}
		if err != nil {
			logger.Fatal().Err(err).Msg("NonceAt")
			return err
		}
		logger.Debug().Msgf("signer %s Posting Nonce of  of nonce %d", zetaClient.GetKeys().signerName, nonce)
		_, err = zetaClient.PostNonce(chain, nonce)
		if err != nil {
			logger.Fatal().Err(err).Msg("PostNonce")
			return err
		}
	}
	return nil
}

func (ob *EVMChainClient) WatchGasPrice() {
	err := ob.PostGasPrice()
	if err != nil {
		ob.logger.Error().Err(err).Msg("PostGasPrice error on " + ob.chain.String())
	}
	gasTicker := time.NewTicker(5 * time.Second) // FIXME: configure this in chainconfig
	for {
		select {
		case <-gasTicker.C:
			err := ob.PostGasPrice()
			if err != nil {
				ob.logger.Error().Err(err).Msg("PostGasPrice error on " + ob.chain.String())
				continue
			}
		case <-ob.stop:
			ob.logger.Info().Msg("WatchGasPrice stopped")
			return
		}
	}
}

func (ob *EVMChainClient) PostGasPrice() error {
	// GAS PRICE
	gasPrice, err := ob.EvmClient.SuggestGasPrice(context.TODO())
	if err != nil {
		ob.logger.Err(err).Msg("PostGasPrice:")
		return err
	}
	blockNum, err := ob.EvmClient.BlockNumber(context.TODO())
	if err != nil {
		ob.logger.Err(err).Msg("PostGasPrice:")
		return err
	}

	// SUPPLY
	var supply string // lockedAmount on ETH, totalSupply on other chains
	supply = "100"

	_, err = ob.zetaClient.PostGasPrice(ob.chain, gasPrice.Uint64(), supply, blockNum)
	if err != nil {
		ob.logger.Err(err).Msg("PostGasPrice:")
		return err
	}

	return nil
}

// query ZetaCore about the last block that it has heard from a specific chain.
// return 0 if not existent.
func (ob *EVMChainClient) getLastHeight() int64 {
	lastheight, err := ob.zetaClient.GetLastBlockHeightByChain(ob.chain)
	if err != nil {
		ob.logger.Warn().Err(err).Msgf("getLastHeight")
		return 0
	}
	return int64(lastheight.LastSendHeight)
}

func (ob *EVMChainClient) BuildBlockIndex() error {
	logger := ob.logger
	envvar := ob.chain.String() + "_SCAN_FROM"
	scanFromBlock := os.Getenv(envvar)
	if scanFromBlock != "" {
		logger.Info().Msgf("envvar %s is set; scan from  block %s", envvar, scanFromBlock)
		if scanFromBlock == clienttypes.EnvVarLatest {
			header, err := ob.EvmClient.HeaderByNumber(context.Background(), nil)
			if err != nil {
				return err
			}
			ob.SetLastBlockHeight(header.Number.Int64())
		} else {
			scanFromBlockInt, err := strconv.ParseInt(scanFromBlock, 10, 64)
			if err != nil {
				return err
			}
			ob.SetLastBlockHeight(scanFromBlockInt)
		}
	} else { // last observed block
		var lastBlockNum clienttypes.LastBlockSQLType
		if dbf := ob.db.First(&lastBlockNum, clienttypes.LastBlockNumId); dbf.Error != nil {
			logger.Info().Msg("db PosKey does not exist; read from ZetaCore")
			ob.SetLastBlockHeight(ob.getLastHeight())
			// if ZetaCore does not have last heard block height, then use current
			if ob.GetLastBlockHeight() == 0 {
				header, err := ob.EvmClient.HeaderByNumber(context.Background(), nil)
				if err != nil {
					return err
				}
				ob.SetLastBlockHeight(header.Number.Int64())
			}
			if dbc := ob.db.Create(clienttypes.ToLastBlockSQLType(ob.GetLastBlockHeight())); dbc.Error != nil {
				logger.Error().Err(dbc.Error).Msg("error writing ob.LastBlock to db: ")
			}
		} else {
			ob.SetLastBlockHeight(lastBlockNum.Num)
		}
	}
	return nil
}

func (ob *EVMChainClient) BuildReceiptsMap() {
	logger := ob.logger
	var receipts []clienttypes.ReceiptSQLType
	dbf := ob.db.Find(&receipts)

	if dbf.Error == nil {
		for _, receipt := range receipts {
			ob.outTXConfirmedReceipts[receipt.Nonce] = &receipt.Receipt
		}
	} else {
		logger.Error().Err(dbf.Error).Msg("error iterating over db")
	}
}

func (ob *EVMChainClient) BuildTransactionsMap() {
	logger := ob.logger
	var transactions []clienttypes.TransactionSQLType
	dbf := ob.db.Find(&transactions)

	if dbf.Error == nil {
		for _, transaction := range transactions {
			ob.outTXConfirmedTransaction[transaction.Nonce] = &transaction.Transaction
		}
	} else {
		logger.Error().Err(dbf.Error).Msg("error iterating over db")
	}
}

// LoadDB open sql database and load data into EVMChainClient
func (ob *EVMChainClient) LoadDB(dbPath string, chain common.Chain) error {
	if dbPath != "" {
		path := fmt.Sprintf("%s/%s", dbPath, chain.String()) //Use "file::memory:?cache=shared" for temp db
		db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		err = db.AutoMigrate(&clienttypes.ReceiptSQLType{},
			&clienttypes.TransactionSQLType{},
			&clienttypes.LastBlockSQLType{})
		if err != nil {
			return err
		}

		ob.db = db
		err = ob.BuildBlockIndex()
		if err != nil {
			return err
		}
		ob.BuildReceiptsMap()
		ob.BuildTransactionsMap()
	}
	return nil
}

func (ob *EVMChainClient) SetChainDetails(chain common.Chain) {
	MinObInterval := 24
	chainconfig := config.ChainConfigs[chain.ChainName.String()]
	ob.confCount = chainconfig.ConfCount
	ob.BlockTime = chainconfig.BlockTime
	ob.ticker = time.NewTicker(time.Duration(MaxInt(int(chainconfig.BlockTime), MinObInterval)) * time.Second)
}

func (ob *EVMChainClient) SetMinAndMaxNonce(trackers []cctxtypes.OutTxTracker) error {
	minNonce, maxNonce := int64(-1), int64(0)
	for _, tracker := range trackers {
		conv := tracker.Nonce
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
