package zetaclient

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/syndtr/goleveldb/leveldb/util"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/zetacore/contracts/evm"
	metricsPkg "github.com/zeta-chain/zetacore/zetaclient/metrics"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"

	cctxtypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
)

const (
	PosKey                 = "PosKey"
	NonceTxHashesKeyPrefix = "NonceTxHashes-"
	NonceTxKeyPrefix       = "NonceTx-"
)

var errEmptyBlock = fmt.Errorf("server returned empty transaction list but block header indicates transactions")

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
	EvmClient                 *ethclient.Client
	zetaClient                *ZetaCoreBridge
	Tss                       TSSSigner
	lastBlock                 uint64
	confCount                 uint64 // must wait this many blocks to be considered "confirmed"
	BlockTime                 uint64 // block time in seconds
	txWatchList               map[ethcommon.Hash]string
	mu                        *sync.Mutex
	db                        *leveldb.DB
	sampleLogger              *zerolog.Logger
	outTXConfirmedReceipts    map[int]*ethtypes.Receipt
	outTXConfirmedTransaction map[int]*ethtypes.Transaction
	MinNonce                  int64
	MaxNonce                  int64
	OutTxChan                 chan OutTx // send to this channel if you want something back!
	ZetaPriceQuerier          ZetaPriceQuerier
	stop                      chan struct{}
	fileLogger                *zerolog.Logger // for critical info
	logger                    zerolog.Logger
}

var _ ChainClient = (*EVMChainClient)(nil)

// Return configuration based on supplied target chain
func NewEVMChainClient(chain common.Chain, bridge *ZetaCoreBridge, tss TSSSigner, dbpath string, metrics *metricsPkg.Metrics) (*EVMChainClient, error) {
	ob := EVMChainClient{
		ChainMetrics: NewChainMetrics(chain.String(), metrics),
	}
	ob.stop = make(chan struct{})
	ob.chain = chain
	ob.mu = &sync.Mutex{}
	sampled := log.Sample(&zerolog.BasicSampler{N: 10})
	ob.sampleLogger = &sampled
	ob.logger = log.With().Str("chain", chain.String()).Logger()
	ob.zetaClient = bridge
	ob.txWatchList = make(map[ethcommon.Hash]string)
	ob.Tss = tss
	ob.outTXConfirmedReceipts = make(map[int]*ethtypes.Receipt)
	ob.outTXConfirmedTransaction = make(map[int]*ethtypes.Transaction)
	ob.OutTxChan = make(chan OutTx, 100)
	addr := ethcommon.HexToAddress(config.Chains[chain.String()].ConnectorContractAddress)
	if addr == ethcommon.HexToAddress("0x0") {
		return nil, fmt.Errorf("connector contract address %s not configured for chain %s", config.Chains[chain.String()].ConnectorContractAddress, chain.String())
	}
	ob.ConnectorAddress = addr
	ob.endpoint = config.Chains[chain.String()].Endpoint
	logFile, err := os.OpenFile(ob.chain.String()+"_debug.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

	if err != nil {
		// Can we log an error before we have our logger? :)
		log.Error().Err(err).Msgf("there was an error creating a logFile chain %s", ob.chain.String())
	}
	fileLogger := zerolog.New(logFile).With().Logger()
	ob.fileLogger = &fileLogger

	// initialize the Client
	ob.logger.Info().Msgf("Chain %s endpoint %s", ob.chain, ob.endpoint)
	client, err := ethclient.Dial(ob.endpoint)
	if err != nil {
		ob.logger.Error().Err(err).Msg("eth Client Dial")
		return nil, err
	}
	ob.EvmClient = client

	// initialize the connector
	connector, err := evm.NewConnector(addr, ob.EvmClient)
	if err != nil {
		ob.logger.Error().Err(err).Msg("Connector")
		return nil, err
	}
	ob.Connector = connector

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

	uniswapV3ABI, err := abi.JSON(strings.NewReader(config.UNISWAPV3POOL))
	if err != nil {
		return nil, err
	}
	uniswapV2ABI, err := abi.JSON(strings.NewReader(config.PANCAKEPOOL))
	if err != nil {
		return nil, err
	}

	uniswapv3Querier, uniswapv2Querier, dummyQuerior := ob.GetPriceQueriers(chain.String(), uniswapV3ABI, uniswapV2ABI)
	ob.SetChainDetails(chain, uniswapv3Querier, uniswapv2Querier)
	if os.Getenv("DUMMY_PRICE") != "" {
		ob.logger.Info().Msg("Using dummy price of 1:1")
		ob.ZetaPriceQuerier = dummyQuerior
	}

	if dbpath != "" {
		err := ob.BuildBlockIndex(dbpath, chain.String())
		if err != nil {
			return nil, err
		}
		ob.BuildReceiptsMap()

	}
	ob.logger.Info().Msgf("%s: start scanning from block %d", chain, ob.GetLastBlockHeight())

	return &ob, nil
}

func (ob *EVMChainClient) Start() {
	go ob.ExternalChainWatcher() // Observes external Chains for incoming trasnactions
	go ob.WatchGasPrice()        // Observes external Chains for Gas prices and posts to core
	go ob.WatchExchangeRate()    // Observers ZetaPriceQuerier for Zeta prices and posts to core
	go ob.observeOutTx()
}

func (ob *EVMChainClient) Stop() {
	ob.logger.Info().Msgf("ob %s is stopping", ob.chain)
	close(ob.stop) // this notifies all goroutines to stop

	ob.logger.Info().Msg("closing ob.db")
	err := ob.db.Close()
	if err != nil {
		ob.logger.Error().Err(err).Msg("error closing db")
	}

	ob.logger.Info().Msgf("%s observer stopped", ob.chain)
}

// returns: isIncluded, isConfirmed, Error
// If isConfirmed, it also post to ZetaCore
func (ob *EVMChainClient) IsSendOutTxProcessed(sendHash string, nonce int, fromOrToZeta bool) (bool, bool, error) {
	ob.mu.Lock()
	receipt, found1 := ob.outTXConfirmedReceipts[nonce]
	transaction, found2 := ob.outTXConfirmedTransaction[nonce]
	ob.mu.Unlock()
	found := found1 && found2
	sendID := fmt.Sprintf("%s/%d", ob.chain.String(), nonce)
	logger := ob.logger.With().Str("sendID", sendID).Logger()
	if fromOrToZeta {
		if found && receipt.Status == 1 {
			zetaHash, err := ob.zetaClient.PostReceiveConfirmation(
				sendHash,
				receipt.TxHash.Hex(),
				receipt.BlockNumber.Uint64(),
				transaction.Value(),
				common.ReceiveStatus_Success,
				ob.chain.String(),
				nonce,
				common.CoinType_Gas,
			)
			if err != nil {
				logger.Error().Err(err).Msg("error posting confirmation to meta core")
			}
			logger.Info().Msgf("Zeta tx hash: %s\n", zetaHash)
			return true, true, nil
		} else if found && receipt.Status == 0 { // the same as below events flow
			logger.Info().Msgf("Found (failed tx) sendHash %s on chain %s txhash %s", sendHash, ob.chain, receipt.TxHash.Hex())
			zetaTxHash, err := ob.zetaClient.PostReceiveConfirmation(sendHash, receipt.TxHash.Hex(), receipt.BlockNumber.Uint64(), big.NewInt(0), common.ReceiveStatus_Failed, ob.chain.String(), nonce, common.CoinType_Gas)
			if err != nil {
				logger.Error().Err(err).Msgf("PostReceiveConfirmation error in WatchTxHashWithTimeout; zeta tx hash %s", zetaTxHash)
			}
			logger.Info().Msgf("Zeta tx hash: %s", zetaTxHash)
			return true, true, nil
		}
	} else {
		if found && receipt.Status == 1 {
			logs := receipt.Logs
			for _, vLog := range logs {
				receivedLog, err := ob.Connector.ConnectorFilterer.ParseZetaReceived(*vLog)
				if err == nil {
					logger.Info().Msgf("Found (outTx) sendHash %s on chain %s txhash %s", sendHash, ob.chain, vLog.TxHash.Hex())
					if vLog.BlockNumber+ob.confCount < ob.GetLastBlockHeight() {
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
							ob.chain.String(),
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
					logger.Info().Msgf("Included; %d blocks before confirmed! chain %s nonce %d", int(vLog.BlockNumber+ob.confCount)-int(ob.GetLastBlockHeight()), ob.chain, nonce)
					return true, false, nil
				}
				revertedLog, err := ob.Connector.ConnectorFilterer.ParseZetaReverted(*vLog)
				if err == nil {
					logger.Info().Msgf("Found (revertTx) sendHash %s on chain %s txhash %s", sendHash, ob.chain, vLog.TxHash.Hex())
					if vLog.BlockNumber+ob.confCount < ob.GetLastBlockHeight() {
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
							ob.chain.String(),
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
					logger.Info().Msgf("Included; %d blocks before confirmed! chain %s nonce %d", int(vLog.BlockNumber+ob.confCount)-int(ob.GetLastBlockHeight()), ob.chain, nonce)
					return true, false, nil
				}
			}
		} else if found && receipt.Status == 0 {
			//FIXME: check nonce here by getTransaction RPC
			logger.Info().Msgf("Found (failed tx) sendHash %s on chain %s txhash %s", sendHash, ob.chain, receipt.TxHash.Hex())
			zetaTxHash, err := ob.zetaClient.PostReceiveConfirmation(sendHash, receipt.TxHash.Hex(), receipt.BlockNumber.Uint64(), big.NewInt(0), common.ReceiveStatus_Failed, ob.chain.String(), nonce, common.CoinType_Zeta)
			if err != nil {
				logger.Error().Err(err).Msgf("PostReceiveConfirmation error in WatchTxHashWithTimeout; zeta tx hash %s", zetaTxHash)
			}
			logger.Info().Msgf("Zeta tx hash: %s", zetaTxHash)
			return true, true, nil
		}
	}

	return false, false, nil
}

// FIXME: there's a chance that a txhash in OutTxChan may not deliver when Stop() is called
// observeOutTx periodically checks all the txhash in potential outbound txs
func (ob *EVMChainClient) observeOutTx() {
	logger := ob.logger
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			trackers, err := ob.zetaClient.GetAllOutTxTrackerByChain(ob.chain)
			if err != nil {
				return
			}
			outTimeout := time.After(90 * time.Second)
		TRACKERLOOP:
			for _, tracker := range trackers {
				nonceInt, err := strconv.Atoi(tracker.Nonce)
				if err != nil {
					return
				}
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
							ob.outTXConfirmedReceipts[nonceInt] = receipt
							ob.outTXConfirmedTransaction[nonceInt] = transaction
							value, err := receipt.MarshalJSON()
							if err != nil {
								logger.Error().Err(err).Msgf("receipt marshal error %s", receipt.TxHash.Hex())
							}
							ob.mu.Unlock()
							err = ob.db.Put([]byte(NonceTxKeyPrefix+fmt.Sprintf("%d", nonceInt)), value, nil)
							if err != nil {
								logger.Error().Err(err).Msgf("PurgeTxHashWatchList: error putting nonce %d tx hashes %s to db", nonceInt, receipt.TxHash.Hex())
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
	receipt, err1 := ob.EvmClient.TransactionReceipt(ctxt, ethcommon.HexToHash(txHash))
	transaction, _, err2 := ob.EvmClient.TransactionByHash(ctxt, ethcommon.HexToHash(txHash))

	if err1 != nil || err2 != nil {
		if err1 != ethereum.NotFound {
			logger.Warn().Err(err1).Msg("TransactionReceipt/TransactionByHash error")
		}
		return nil, nil, err1
	} else if receipt.BlockNumber.Uint64()+ob.confCount > ob.GetLastBlockHeight() {
		log.Warn().Msgf("included but not confirmed: receipt block %d, current block %d", receipt.BlockNumber, ob.GetLastBlockHeight())
		return nil, nil, fmt.Errorf("included but not confirmed")
	} else { // confirmed outbound tx
		return receipt, transaction, nil
	}
}

func (ob *EVMChainClient) SetLastBlockHeight(block uint64) {
	atomic.StoreUint64(&ob.lastBlock, block)
}

func (ob *EVMChainClient) GetLastBlockHeight() uint64 {
	return atomic.LoadUint64(&ob.lastBlock)
}

func (ob *EVMChainClient) ExternalChainWatcher() {
	// At each tick, query the Connector contract
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
	if confirmedBlockNum <= ob.GetLastBlockHeight() {
		ob.sampleLogger.Info().Msg("Skipping observer , No new block is produced ")
		return nil
	}
	lastBlock := ob.GetLastBlockHeight()
	startBlock := lastBlock + 1
	toBlock := lastBlock + config.MaxBlocksPerPeriod // read at most 10 blocks in one go
	if toBlock >= confirmedBlockNum {
		toBlock = confirmedBlockNum
	}
	ob.sampleLogger.Info().Msgf("%s current block %d, querying from %d to %d, %d blocks left to catch up, watching MPI address %s", ob.chain, header.Number.Uint64(), ob.GetLastBlockHeight()+1, toBlock, int(toBlock)-int(confirmedBlockNum), ob.ConnectorAddress.Hex())

	// Finally query the for the logs
	logs, err := ob.Connector.FilterZetaSent(&bind.FilterOpts{
		Start:   startBlock,
		End:     &toBlock,
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

		destChain := config.FindChainByID(event.DestinationChainId)
		destAddr := clienttypes.BytesToEthHex(event.DestinationAddress)
		if strings.EqualFold(destAddr, config.Chains[destChain].ZETATokenContractAddress) {
			ob.logger.Warn().Msgf("potential attack attempt: %s destination address is ZETA token contract address %s", destChain, destAddr)
		}
		zetaHash, err := ob.zetaClient.PostSend(
			event.ZetaTxSenderAddress.Hex(),
			ob.chain.String(),
			clienttypes.BytesToEthHex(event.DestinationAddress),
			config.FindChainByID(event.DestinationChainId),
			event.ZetaValueAndGas.String(),
			event.ZetaValueAndGas.String(),
			base64.StdEncoding.EncodeToString(event.Message),
			event.Raw.TxHash.Hex(),
			event.Raw.BlockNumber,
			event.DestinationGasLimit.Uint64(),
			common.CoinType_Zeta,
			PostSendNonEVMGasLimit,
		)
		if err != nil {
			ob.logger.Error().Err(err).Msg("error posting to zeta core")
			continue
		}
		ob.logger.Info().Msgf("ZetaSent event detected and reported: PostSend zeta tx: %s", zetaHash)
	}

	// ============= query the incoming tx to TSS address ==============
	tssAddress := ob.Tss.EVMAddress()
	// query incoming gas asset
	for bn := startBlock; bn <= toBlock; bn++ {
		//block, err := ob.EvmClient.BlockByNumber(context.Background(), big.NewInt(int64(bn)))
		block, err := ob.EvmClient.BlockByNumber(context.Background(), big.NewInt(int64(bn)))
		if err != nil {
			//TODO: this is very hacky becaue klatyn uses different empty tx hash as ethereum:
			// see: https://github.com/klaytn/klaytn/blob/febce7b01a616a556423704cf9faa7da4bc4753f/client/klay_client.go#L119
			if ob.chain == common.BaobabChain && strings.Contains(err.Error(), errEmptyBlock.Error()) {
			} else {
				ob.logger.Error().Err(err).Msgf("error getting block: %d", bn)
			}
			continue
		}
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
					ob.logger.Err(err).Msg("TransactionSender")
					continue
				}
				ob.logger.Info().Msgf("TSS inTx detected: %s, blocknum %d", tx.Hash().Hex(), receipt.BlockNumber)
				ob.logger.Info().Msgf("TSS inTx value: %s", tx.Value().String())
				ob.logger.Info().Msgf("TSS inTx from: %s", from.Hex())
				message := ""
				if len(tx.Data()) != 0 {
					message = hex.EncodeToString(tx.Data())
				}
				zetaHash, err := ob.zetaClient.PostSend(
					from.Hex(),
					ob.chain.String(),
					from.Hex(),
					"ZETA",
					tx.Value().String(),
					tx.Value().String(),
					message,
					tx.Hash().Hex(),
					receipt.BlockNumber.Uint64(),
					90_000,
					common.CoinType_Gas,
					PostSendEVMGasLimit,
				)
				if err != nil {
					ob.logger.Error().Err(err).Msg("error posting to zeta core")
					continue
				}
				ob.logger.Info().Msgf("ZetaSent event detected and reported: PostSend zeta tx: %s", zetaHash)
			}
		}
	}
	// ============= end of query the incoming tx to TSS address ==============

	//ob.LastBlock = toBlock
	ob.SetLastBlockHeight(toBlock)
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, toBlock)
	err = ob.db.Put([]byte(PosKey), buf[:n], nil)
	if err != nil {
		ob.logger.Error().Err(err).Msg("error writing toBlock to db")
	}
	return nil
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
	gasTicker := time.NewTicker(60 * time.Second)
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
	//if chainOb.chain == common.ETHChain {
	//	input, err := chainOb.connectorAbi.Pack("getLockedAmount")
	//	if err != nil {
	//		return fmt.Errorf("fail to getLockedAmount")
	//	}
	//	bn, err := chainOb.Client.BlockNumber(context.TODO())
	//	if err != nil {
	//		log.Err(err).Msgf("%s BlockNumber error", chainOb.chain)
	//		return err
	//	}
	//	fromAddr := ethcommon.HexToAddress(config.TSS_TEST_ADDRESS)
	//	toAddr := ethcommon.HexToAddress(config.ETH_MPI_ADDRESS)
	//	res, err := chainOb.Client.CallContract(context.TODO(), ethereum.CallMsg{
	//		From: fromAddr,
	//		To:   &toAddr,
	//		Data: input,
	//	}, big.NewInt(0).SetUint64(bn))
	//	if err != nil {
	//		log.Err(err).Msgf("%s CallContract error", chainOb.chain)
	//		return err
	//	}
	//	output, err := chainOb.connectorAbi.Unpack("getLockedAmount", res)
	//	if err != nil {
	//		log.Err(err).Msgf("%s Unpack error", chainOb.chain)
	//		return err
	//	}
	//	lockedAmount := *connectorAbi.ConvertType(output[0], new(*big.Int)).(**big.Int)
	//	//fmt.Printf("ETH: block %d: lockedAmount %d\n", bn, lockedAmount)
	//	supply = lockedAmount.String()
	//
	//} else if chainOb.chain == common.BSCChain {
	//	input, err := chainOb.connectorAbi.Pack("totalSupply")
	//	if err != nil {
	//		return fmt.Errorf("fail to totalSupply")
	//	}
	//	bn, err := chainOb.Client.BlockNumber(context.TODO())
	//	if err != nil {
	//		log.Err(err).Msgf("%s BlockNumber error", chainOb.chain)
	//		return err
	//	}
	//	fromAddr := ethcommon.HexToAddress(config.TSS_TEST_ADDRESS)
	//	toAddr := ethcommon.HexToAddress(config.BSC_MPI_ADDRESS)
	//	res, err := chainOb.Client.CallContract(context.TODO(), ethereum.CallMsg{
	//		From: fromAddr,
	//		To:   &toAddr,
	//		Data: input,
	//	}, big.NewInt(0).SetUint64(bn))
	//	if err != nil {
	//		log.Err(err).Msgf("%s CallContract error", chainOb.chain)
	//		return err
	//	}
	//	output, err := chainOb.connectorAbi.Unpack("totalSupply", res)
	//	if err != nil {
	//		log.Err(err).Msgf("%s Unpack error", chainOb.chain)
	//		return err
	//	}
	//	totalSupply := *connectorAbi.ConvertType(output[0], new(*big.Int)).(**big.Int)
	//	//fmt.Printf("BSC: block %d: totalSupply %d\n", bn, totalSupply)
	//	supply = totalSupply.String()
	//} else if chainOb.chain == common.POLYGONChain {
	//	input, err := chainOb.connectorAbi.Pack("totalSupply")
	//	if err != nil {
	//		return fmt.Errorf("fail to totalSupply")
	//	}
	//	bn, err := chainOb.Client.BlockNumber(context.TODO())
	//	if err != nil {
	//		log.Err(err).Msgf("%s BlockNumber error", chainOb.chain)
	//		return err
	//	}
	//	fromAddr := ethcommon.HexToAddress(config.TSS_TEST_ADDRESS)
	//	toAddr := ethcommon.HexToAddress(config.POLYGON_MPI_ADDRESS)
	//	res, err := chainOb.Client.CallContract(context.TODO(), ethereum.CallMsg{
	//		From: fromAddr,
	//		To:   &toAddr,
	//		Data: input,
	//	}, big.NewInt(0).SetUint64(bn))
	//	if err != nil {
	//		log.Err(err).Msgf("%s CallContract error", chainOb.chain)
	//		return err
	//	}
	//	output, err := chainOb.connectorAbi.Unpack("totalSupply", res)
	//	if err != nil {
	//		log.Err(err).Msgf("%s Unpack error", chainOb.chain)
	//		return err
	//	}
	//	totalSupply := *connectorAbi.ConvertType(output[0], new(*big.Int)).(**big.Int)
	//	//fmt.Printf("BSC: block %d: totalSupply %d\n", bn, totalSupply)
	//	supply = totalSupply.String()
	//} else {
	//	log.Error().Msgf("chain not supported %s", chainOb.chain)
	//	return fmt.Errorf("unsupported chain %s", chainOb.chain)
	//}

	_, err = ob.zetaClient.PostGasPrice(ob.chain, gasPrice.Uint64(), supply, blockNum)
	if err != nil {
		ob.logger.Err(err).Msg("PostGasPrice:")
		return err
	}

	//bal, err := chainOb.Client.BalanceAt(context.TODO(), chainOb.Tss.EVMAddress(), nil)
	//if err != nil {
	//	log.Err(err).Msg("BalanceAt:")
	//	return err
	//}
	//_, err = chainOb.zetaClient.PostGasBalance(chainOb.chain, bal.String(), blockNum)
	//if err != nil {
	//	log.Err(err).Msg("PostGasBalance:")
	//	return err
	//}
	return nil
}

// query ZetaCore about the last block that it has heard from a specific chain.
// return 0 if not existent.
func (ob *EVMChainClient) getLastHeight() uint64 {
	lastheight, err := ob.zetaClient.GetLastBlockHeightByChain(ob.chain)
	if err != nil {
		ob.logger.Warn().Err(err).Msgf("getLastHeight")
		return 0
	}
	return lastheight.LastSendHeight
}

func (ob *EVMChainClient) WatchExchangeRate() {
	ticker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-ticker.C:
			price, bn, err := ob.ZetaPriceQuerier.GetZetaPrice()
			if err != nil {
				ob.logger.Error().Err(err).Msg("GetZetaExchangeRate error")
				continue
			}
			priceInHex := fmt.Sprintf("0x%x", price)
			_, err = ob.zetaClient.PostZetaConversionRate(ob.chain, priceInHex, bn)
			if err != nil {
				ob.logger.Error().Err(err).Msg("PostZetaConversionRate error")
			}
		case <-ob.stop:
			ob.logger.Info().Msg("WatchExchangeRate stopped")
			return
		}
	}
}

func (ob *EVMChainClient) BuildBlockIndex(dbpath, chain string) error {
	logger := ob.logger
	path := fmt.Sprintf("%s/%s", dbpath, chain) // e.g. ~/.zetaclient/ETH
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return err
	}
	ob.db = db
	envvar := ob.chain.String() + "_SCAN_FROM"
	scanFromBlock := os.Getenv(envvar)
	if scanFromBlock != "" {
		logger.Info().Msgf("envvar %s is set; scan from  block %s", envvar, scanFromBlock)
		if scanFromBlock == "latest" {
			header, err := ob.EvmClient.HeaderByNumber(context.Background(), nil)
			if err != nil {
				return err
			}
			ob.SetLastBlockHeight(header.Number.Uint64())
		} else {
			scanFromBlockInt, err := strconv.ParseInt(scanFromBlock, 10, 64)
			if err != nil {
				return err
			}
			ob.SetLastBlockHeight(uint64(scanFromBlockInt))
		}
	} else { // last observed block
		buf, err := db.Get([]byte(PosKey), nil)
		if err != nil {
			logger.Info().Msg("db PosKey does not exist; read from ZetaCore")
			ob.SetLastBlockHeight(ob.getLastHeight())
			// if ZetaCore does not have last heard block height, then use current
			if ob.GetLastBlockHeight() == 0 {
				header, err := ob.EvmClient.HeaderByNumber(context.Background(), nil)
				if err != nil {
					return err
				}
				ob.SetLastBlockHeight(header.Number.Uint64())
			}
			buf2 := make([]byte, binary.MaxVarintLen64)
			n := binary.PutUvarint(buf2, ob.GetLastBlockHeight())
			err := db.Put([]byte(PosKey), buf2[:n], nil)
			if err != nil {
				logger.Error().Err(err).Msg("error writing ob.LastBlock to db: ")
			}
		} else {
			lastBlock, _ := binary.Uvarint(buf)
			ob.SetLastBlockHeight(lastBlock)
		}
	}
	return nil
}

func (ob *EVMChainClient) BuildReceiptsMap() {
	logger := ob.logger
	iter := ob.db.NewIterator(util.BytesPrefix([]byte(NonceTxKeyPrefix)), nil)
	for iter.Next() {
		key := string(iter.Key())
		nonce, err := strconv.ParseInt(key[len(NonceTxKeyPrefix):], 10, 64)
		if err != nil {
			logger.Error().Err(err).Msgf("error parsing nonce: %s", key)
			continue
		}
		var receipt ethtypes.Receipt
		err = receipt.UnmarshalJSON(iter.Value())
		if err != nil {
			logger.Error().Err(err).Msgf("error unmarshalling receipt: %s", key)
			continue
		}
		ob.outTXConfirmedReceipts[int(nonce)] = &receipt
		//log.Info().Msgf("chain %s reading nonce %d with receipt of tx %s", ob.chain, nonce, receipt.TxHash.Hex())
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		logger.Error().Err(err).Msg("error iterating over db")
	}
}

func (ob *EVMChainClient) GetPriceQueriers(chain string, uniswapV3ABI, uniswapV2ABI abi.ABI) (*UniswapV3ZetaPriceQuerier, *UniswapV2ZetaPriceQuerier, *DummyZetaPriceQuerier) {
	uniswapv3querier := &UniswapV3ZetaPriceQuerier{
		UniswapV3Abi:        &uniswapV3ABI,
		Client:              ob.EvmClient,
		PoolContractAddress: ethcommon.HexToAddress(config.Chains[chain].PoolContractAddress),
		Chain:               ob.chain,
		TokenOrder:          config.Chains[chain].PoolTokenOrder,
	}
	uniswapv2querier := &UniswapV2ZetaPriceQuerier{
		UniswapV2Abi:        &uniswapV2ABI,
		Client:              ob.EvmClient,
		PoolContractAddress: ethcommon.HexToAddress(config.Chains[chain].PoolContractAddress),
		Chain:               ob.chain,
		TokenOrder:          config.Chains[chain].PoolTokenOrder,
	}
	dummyQuerier := &DummyZetaPriceQuerier{
		Chain:  ob.chain,
		Client: ob.EvmClient,
	}
	return uniswapv3querier, uniswapv2querier, dummyQuerier
}

func (ob *EVMChainClient) SetChainDetails(chain common.Chain,
	uniswapv3querier *UniswapV3ZetaPriceQuerier,
	uniswapv2querier *UniswapV2ZetaPriceQuerier) {
	MinObInterval := 24
	switch chain {
	case common.MumbaiChain:
		ob.ticker = time.NewTicker(time.Duration(MaxInt(config.PolygonBlockTime, MinObInterval)) * time.Second)
		ob.confCount = config.PolygonConfirmationCount
		ob.BlockTime = config.PolygonBlockTime

	case common.GoerliChain:
		ob.ticker = time.NewTicker(time.Duration(MaxInt(config.EthBlockTime, MinObInterval)) * time.Second)
		ob.confCount = config.EthConfirmationCount
		ob.BlockTime = config.EthBlockTime

	case common.BSCTestnetChain:
		ob.ticker = time.NewTicker(time.Duration(MaxInt(config.BscBlockTime, MinObInterval)) * time.Second)
		ob.confCount = config.BscConfirmationCount
		ob.BlockTime = config.BscBlockTime

	case common.BaobabChain:
		ob.ticker = time.NewTicker(time.Duration(MaxInt(config.EthBlockTime, MinObInterval)) * time.Second)
		ob.confCount = config.EthConfirmationCount
		ob.BlockTime = config.EthBlockTime

	case common.BTCTestnetChain:
		ob.ticker = time.NewTicker(time.Duration(MaxInt(config.EthBlockTime, MinObInterval)) * time.Second)
		ob.confCount = config.BtcConfirmationCount
		ob.BlockTime = config.EthBlockTime

	case common.Ganache:
		ob.ticker = time.NewTicker(time.Duration(MaxInt(config.RopstenBlockTime, MinObInterval)) * time.Second)
		ob.confCount = 0
		ob.BlockTime = 1
	}
	switch config.Chains[chain.String()].PoolContract {
	case clienttypes.UniswapV2:
		ob.ZetaPriceQuerier = uniswapv2querier
	case clienttypes.UniswapV3:
		ob.ZetaPriceQuerier = uniswapv3querier
	default:
		ob.logger.Error().Msgf("unknown pool contract type: %d", config.Chains[chain.String()].PoolContract)
	}
}

func (ob *EVMChainClient) SetMinAndMaxNonce(trackers []cctxtypes.OutTxTracker) error {
	minNonce, maxNonce := int64(-1), int64(0)
	for _, tracker := range trackers {
		conv, err := strconv.Atoi(tracker.Nonce)
		if err != nil {
			return err
		}
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
