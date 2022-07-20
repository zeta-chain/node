package zetaclient

import (
	"context"
	"encoding/binary"
	"fmt"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/zeta-chain/zetacore/contracts/evm"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
)

const (
	PosKey                 = "PosKey"
	NonceTxHashesKeyPrefix = "NonceTxHashes-"
	NonceTxKeyPrefix       = "NonceTx-"
)

//    event ZetaSent(
//        address indexed originSenderAddress,
//        uint256 destinationChainId,
//        bytes destinationAddress,
//        uint256 zetaAmount,
//        uint256 gasLimit,
//        bytes message,
//        bytes zetaParams
//    );
var logZetaSentSignature = []byte("ZetaSent(address,uint256,bytes,uint256,uint256,bytes,bytes)")
var logZetaSentSignatureHash = crypto.Keccak256Hash(logZetaSentSignature)

//    event ZetaReceived(
//        bytes originSenderAddress,
//        uint256 indexed originChainId,
//        address indexed destinationAddress,
//        uint256 zetaAmount,
//        bytes message,
//        bytes32 indexed internalSendHash
//    );
var logZetaReceivedSignature = []byte("ZetaReceived(bytes,uint256,address,uint256,bytes,bytes32)")
var logZetaReceivedSignatureHash = crypto.Keccak256Hash(logZetaReceivedSignature)

//event ZetaReverted(
//address originSenderAddress,
//uint256 originChainId,
//uint256 indexed destinationChainId,
//bytes indexed destinationAddress,
//uint256 zetaAmount,
//bytes message,
//bytes32 indexed internalSendHash
//);
var logZetaRevertedSignature = []byte("ZetaReverted(address,uint256,uint256,bytes,uint256,bytes,bytes32)")
var logZetaRevertedSignatureHash = crypto.Keccak256Hash(logZetaRevertedSignature)

var topics = make([][]ethcommon.Hash, 1)

type TxHashEnvelope struct {
	TxHash string
	Done   chan struct{}
}

type OutTx struct {
	SendHash string
	TxHash   string
	Nonce    int
}

// Chain configuration struct
// Filled with above constants depending on chain
type ChainObserver struct {
	chain            common.Chain
	endpoint         string
	ticker           *time.Ticker
	Connector        *evm.Connector
	ConnectorAddress ethcommon.Address
	EvmClient        *ethclient.Client
	zetaClient       *ZetaCoreBridge
	Tss              TSSSigner
	LastBlock        uint64
	confCount        uint64 // must wait this many blocks to be considered "confirmed"
	BlockTime        uint64 // block time in seconds
	txWatchList      map[ethcommon.Hash]string
	mu               *sync.Mutex
	db               *leveldb.DB
	sampleLogger     *zerolog.Logger
	metrics          *metrics.Metrics
	outTXPending     map[int][]string
	outTXConfirmed   map[int]*ethtypes.Receipt
	MinNonce         int
	MaxNonce         int
	OutTxChan        chan OutTx // send to this channel if you want something back!
	ZetaPriceQuerier ZetaPriceQuerier
	stop             chan struct{}
	wg               sync.WaitGroup

	fileLogger *zerolog.Logger // for critical info
}

// Return configuration based on supplied target chain
func NewChainObserver(chain common.Chain, bridge *ZetaCoreBridge, tss TSSSigner, dbpath string, metrics *metrics.Metrics) (*ChainObserver, error) {
	ob := ChainObserver{}
	ob.stop = make(chan struct{})
	ob.chain = chain
	ob.mu = &sync.Mutex{}
	sampled := log.Sample(&zerolog.BasicSampler{N: 10})
	ob.sampleLogger = &sampled
	ob.zetaClient = bridge
	ob.txWatchList = make(map[ethcommon.Hash]string)
	ob.Tss = tss
	ob.metrics = metrics
	ob.outTXPending = make(map[int][]string)
	ob.outTXConfirmed = make(map[int]*ethtypes.Receipt)
	ob.OutTxChan = make(chan OutTx, 100)
	addr := ethcommon.HexToAddress(config.Chains[chain.String()].ConnectorContractAddress)
	if addr == ethcommon.HexToAddress("0x0") {
		return nil, fmt.Errorf("Connector contract address %s not configured for chain %s", config.Chains[chain.String()].ConnectorContractAddress, chain.String())
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
	log.Info().Msgf("Chain %s endpoint %s", ob.chain, ob.endpoint)
	client, err := ethclient.Dial(ob.endpoint)
	if err != nil {
		log.Error().Err(err).Msg("eth Client Dial")
		return nil, err
	}
	ob.EvmClient = client

	// initialize the connector
	connector, err := evm.NewConnector(addr, ob.EvmClient)
	if err != nil {
		log.Error().Err(err).Msg("Connector")
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

	uniswapV3ABI, err := abi.JSON(strings.NewReader(config.UNISWAPV3POOL))
	if err != nil {
		return nil, err
	}
	uniswapV2ABI, err := abi.JSON(strings.NewReader(config.PANCAKEPOOL))
	if err != nil {
		return nil, err
	}

	// initialize zeta price queriers
	uniswapv3querier := &UniswapV3ZetaPriceQuerier{
		UniswapV3Abi:        &uniswapV3ABI,
		Client:              ob.EvmClient,
		PoolContractAddress: ethcommon.HexToAddress(config.Chains[chain.String()].PoolContractAddress),
		Chain:               ob.chain,
	}
	uniswapv2querier := &UniswapV2ZetaPriceQuerier{
		UniswapV2Abi:        &uniswapV2ABI,
		Client:              ob.EvmClient,
		PoolContractAddress: ethcommon.HexToAddress(config.Chains[chain.String()].PoolContractAddress),
		Chain:               ob.chain,
	}
	dummyQuerier := &DummyZetaPriceQuerier{
		Chain:  ob.chain,
		Client: ob.EvmClient,
	}

	// Initialize chain specific setup
	MIN_OB_INTERVAL := 24 // minimum 24s between observations
	switch chain {
	case common.MumbaiChain:
		ob.ticker = time.NewTicker(time.Duration(MaxInt(config.POLY_BLOCK_TIME, MIN_OB_INTERVAL)) * time.Second)
		ob.confCount = config.POLYGON_CONFIRMATION_COUNT
		ob.ZetaPriceQuerier = uniswapv3querier
		ob.BlockTime = config.POLY_BLOCK_TIME

	case common.GoerliChain:
		ob.ticker = time.NewTicker(time.Duration(MaxInt(config.ETH_BLOCK_TIME, MIN_OB_INTERVAL)) * time.Second)
		ob.confCount = config.ETH_CONFIRMATION_COUNT
		ob.ZetaPriceQuerier = uniswapv3querier
		ob.BlockTime = config.ETH_BLOCK_TIME

	case common.BSCTestnetChain:
		ob.ticker = time.NewTicker(time.Duration(MaxInt(config.BSC_BLOCK_TIME, MIN_OB_INTERVAL)) * time.Second)
		ob.confCount = config.BSC_CONFIRMATION_COUNT
		ob.ZetaPriceQuerier = uniswapv2querier
		ob.BlockTime = config.BSC_BLOCK_TIME

	case common.RopstenChain:
		ob.ticker = time.NewTicker(time.Duration(MaxInt(config.ROPSTEN_BLOCK_TIME, MIN_OB_INTERVAL)) * time.Second)
		ob.confCount = config.ROPSTEN_CONFIRMATION_COUNT
		ob.ZetaPriceQuerier = uniswapv3querier
		ob.BlockTime = config.ROPSTEN_BLOCK_TIME
	}

	if os.Getenv("DUMMY_PRICE") != "" {
		log.Info().Msg("Using dummy price of 1:1")
		ob.ZetaPriceQuerier = dummyQuerier
	}

	if dbpath != "" {
		path := fmt.Sprintf("%s/%s", dbpath, chain.String()) // e.g. ~/.zetaclient/ETH
		db, err := leveldb.OpenFile(path, nil)
		if err != nil {
			return nil, err
		}
		ob.db = db

		envvar := ob.chain.String() + "_SCAN_CURRENT"
		if os.Getenv(envvar) != "" {
			log.Info().Msgf("envvar %s is set; scan from current block", envvar)
			header, err := ob.EvmClient.HeaderByNumber(context.Background(), nil)
			if err != nil {
				return nil, err
			}
			ob.LastBlock = header.Number.Uint64()
		} else { // last observed block
			buf, err := db.Get([]byte(PosKey), nil)
			if err != nil {
				log.Info().Msg("db PosKey does not exist; read from ZetaCore")
				ob.LastBlock = ob.getLastHeight()
				// if ZetaCore does not have last heard block height, then use current
				if ob.LastBlock == 0 {
					header, err := ob.EvmClient.HeaderByNumber(context.Background(), nil)
					if err != nil {
						return nil, err
					}
					ob.LastBlock = header.Number.Uint64()
				}
				buf2 := make([]byte, binary.MaxVarintLen64)
				n := binary.PutUvarint(buf2, ob.LastBlock)
				err := db.Put([]byte(PosKey), buf2[:n], nil)
				if err != nil {
					log.Error().Err(err).Msg("error writing ob.LastBlock to db: ")
				}
			} else {
				ob.LastBlock, _ = binary.Uvarint(buf)
			}
		}

		{
			iter := ob.db.NewIterator(util.BytesPrefix([]byte(NonceTxHashesKeyPrefix)), nil)
			for iter.Next() {
				key := string(iter.Key())
				nonce, err := strconv.ParseInt(key[len(NonceTxHashesKeyPrefix):], 10, 64)
				if err != nil {
					log.Error().Err(err).Msgf("error parsing nonce: %s", key)
					continue
				}
				txHashes := strings.Split(string(iter.Value()), ",")
				ob.outTXPending[int(nonce)] = txHashes
				log.Info().Msgf("reading nonce %d with %d tx hashes", nonce, len(txHashes))
			}
			iter.Release()
			if err = iter.Error(); err != nil {
				log.Error().Err(err).Msg("error iterating over db")
			}
		}

		{
			iter := ob.db.NewIterator(util.BytesPrefix([]byte(NonceTxKeyPrefix)), nil)
			for iter.Next() {
				key := string(iter.Key())
				nonce, err := strconv.ParseInt(key[len(NonceTxKeyPrefix):], 10, 64)
				if err != nil {
					log.Error().Err(err).Msgf("error parsing nonce: %s", key)
					continue
				}
				var receipt ethtypes.Receipt
				err = receipt.UnmarshalJSON(iter.Value())
				if err != nil {
					log.Error().Err(err).Msgf("error unmarshalling receipt: %s", key)
					continue
				}
				ob.outTXConfirmed[int(nonce)] = &receipt
				log.Info().Msgf("chain %s reading nonce %d with receipt of tx %s", ob.chain, nonce, receipt.TxHash.Hex())
			}
			iter.Release()
			if err = iter.Error(); err != nil {
				log.Error().Err(err).Msg("error iterating over db")
			}
		}

	}
	log.Info().Msgf("%s: start scanning from block %d", chain, ob.LastBlock)

	// this is shared structure to query logs by sendHash
	log.Info().Msgf("Chain %s logZetaReceivedSignatureHash %s", ob.chain, logZetaReceivedSignatureHash.Hex())

	return &ob, nil
}

func (ob *ChainObserver) Start() {
	go ob.ExternalChainWatcher() // Observes external Chains for incoming trasnactions
	go ob.WatchGasPrice()        // Observes external Chains for Gas prices and posts to core
	go ob.WatchExchangeRate()    // Observers ZetaPriceQuerier for Zeta prices and posts to core
	go ob.observeOutTx()
}

func (ob *ChainObserver) Stop() {
	log.Info().Msgf("ob %s is stopping", ob.chain)
	close(ob.stop) // this notifies all goroutines to stop

	log.Info().Msg("closing ob.db")
	err := ob.db.Close()
	if err != nil {
		log.Error().Err(err).Msg("error closing db")
	}

	log.Info().Msgf("%s observer stopped", ob.chain)
}

// returns: isIncluded, isConfirmed, Error
// If isConfirmed, it also post to ZetaCore
func (ob *ChainObserver) IsSendOutTxProcessed(sendHash string, nonce int) (bool, bool, error) {
	receipt, found := ob.outTXConfirmed[nonce]
	if found && receipt.Status == 1 {
		logs := receipt.Logs
		for _, vLog := range logs {
			receivedLog, err := ob.Connector.ConnectorFilterer.ParseZetaReceived(*vLog)
			if err == nil {
				log.Info().Msgf("Found (outTx) sendHash %s on chain %s txhash %s", sendHash, ob.chain, vLog.TxHash.Hex())
				if vLog.BlockNumber+ob.confCount < ob.LastBlock {
					log.Info().Msg("Confirmed! Sending PostConfirmation to zetacore...")
					sendhash := vLog.Topics[3].Hex()
					//var rxAddress string = ethcommon.HexToAddress(vLog.Topics[1].Hex()).Hex()
					mMint := receivedLog.ZetaAmount.String()
					zetaHash, err := ob.zetaClient.PostReceiveConfirmation(
						sendhash,
						vLog.TxHash.Hex(),
						vLog.BlockNumber,
						mMint,
						common.ReceiveStatus_Success,
						ob.chain.String(),
					)
					if err != nil {
						log.Error().Err(err).Msg("error posting confirmation to meta core")
						continue
					}
					log.Info().Msgf("Zeta tx hash: %s\n", zetaHash)
					return true, true, nil
				} else {
					log.Info().Msgf("Included; %d blocks before confirmed! chain %s nonce %d", int(vLog.BlockNumber+ob.confCount)-int(ob.LastBlock), ob.chain, nonce)
					return true, false, nil
				}
			}
			revertedLog, err := ob.Connector.ConnectorFilterer.ParseZetaReverted(*vLog)
			if err == nil {
				log.Info().Msgf("Found (revertTx) sendHash %s on chain %s txhash %s", sendHash, ob.chain, vLog.TxHash.Hex())
				if vLog.BlockNumber+ob.confCount < ob.LastBlock {
					log.Info().Msg("Confirmed! Sending PostConfirmation to zetacore...")
					sendhash := vLog.Topics[3].Hex()
					mMint := revertedLog.ZetaAmount.String()
					metaHash, err := ob.zetaClient.PostReceiveConfirmation(
						sendhash,
						vLog.TxHash.Hex(),
						vLog.BlockNumber,
						mMint,
						common.ReceiveStatus_Success,
						ob.chain.String(),
					)
					if err != nil {
						log.Err(err).Msg("error posting confirmation to meta core")
						continue
					}
					log.Info().Msgf("Zeta tx hash: %s", metaHash)
					return true, true, nil
				} else {
					log.Info().Msgf("Included; %d blocks before confirmed! chain %s nonce %d", int(vLog.BlockNumber+ob.confCount)-int(ob.LastBlock), ob.chain, nonce)
					return true, false, nil
				}
			}
		}
	} else if found && receipt.Status == 0 {
		//FIXME: check nonce here by getTransaction RPC
		log.Info().Msgf("Found (failed tx) sendHash %s on chain %s txhash %s", sendHash, ob.chain, receipt.TxHash.Hex())
		zetaTxHash, err := ob.zetaClient.PostReceiveConfirmation(sendHash, receipt.TxHash.Hex(), receipt.BlockNumber.Uint64(), "", common.ReceiveStatus_Failed, ob.chain.String())
		if err != nil {
			log.Error().Err(err).Msgf("PostReceiveConfirmation error in WatchTxHashWithTimeout; zeta tx hash %s", zetaTxHash)
		}
		log.Info().Msgf("Zeta tx hash: %s", zetaTxHash)
		return true, true, nil
	}

	return false, false, fmt.Errorf("IsSendOutTxProcessed: error on chain %s", ob.chain)
}

// FIXME: there's a chance that a txhash in OutTxChan may not deliver when Stop() is called
// observeOutTx periodically checks all the txhash in potential outbound txs
// If it is able to confirm one of txHashes in a outTXPending , it cleans it and saves information to local DB
func (ob *ChainObserver) observeOutTx() {

	ticker := time.NewTicker(12 * time.Second)
	for {
		select {
		case <-ticker.C:
			minNonce, maxNonce, err := ob.PurgeTxHashWatchList()
			if len(ob.outTXPending) > 0 {
				log.Info().Msgf("chain %s outstanding nonce: %d; nonce range [%d,%d]", ob.chain, len(ob.outTXPending), minNonce, maxNonce)
			}
			outTimeout := time.After(12 * time.Second)
			if err == nil {
				ob.MinNonce = minNonce
				ob.MaxNonce = maxNonce
				//log.Warn().Msgf("chain %s MinNonce: %d", ob.chain, ob.MinNonce)
			QUERYLOOP:
				//for nonce, txHashes := range ob.nonceTxHashesMap {
				for nonce := minNonce; nonce <= maxNonce; nonce++ { // ensure lower nonce is queried first
					ob.mu.Lock()
					txHashes, found := ob.outTXPending[nonce]
					txHashesCopy := txHashes
					ob.mu.Unlock()
					if !found {
						continue
					}
				TXHASHLOOP:
					for _, txHash := range txHashesCopy {
						inTimeout := time.After(1000 * time.Millisecond)
						select {
						case <-outTimeout:
							log.Warn().Msgf("QUERYLOOP timouet chain %s nonce %d", ob.chain, nonce)
							break QUERYLOOP
						default:
							receipt, err := ob.queryTxByHash(txHash, nonce)
							if err == nil && receipt != nil { // confirmed
								log.Info().Msgf("observeOutTx: %s nonce %d, txHash %s confirmed", ob.chain, nonce, txHash)
								ob.mu.Lock()
								delete(ob.outTXPending, nonce)
								if err = ob.db.Delete([]byte(NonceTxHashesKeyPrefix+fmt.Sprintf("%d", nonce)), nil); err != nil {
									log.Error().Err(err).Msgf("PurgeTxHashWatchList: error deleting nonce %d tx hashes from db", nonce)
								}
								ob.outTXConfirmed[nonce] = receipt
								value, err := receipt.MarshalJSON()
								if err != nil {
									log.Error().Err(err).Msgf("receipt marshal error %s", receipt.TxHash.Hex())
								}

								ob.mu.Unlock()
								err = ob.db.Put([]byte(NonceTxKeyPrefix+fmt.Sprintf("%d", nonce)), value, nil)
								if err != nil {
									log.Error().Err(err).Msgf("PurgeTxHashWatchList: error putting nonce %d tx hashes %s to db", nonce, receipt.TxHash.Hex())
								}

								break TXHASHLOOP
							}
							<-inTimeout
						}
					}
				}
			} else {
				log.Warn().Err(err).Msg("PurgeTxHashWatchList error")
			}
		case <-ob.stop:
			log.Info().Msg("observeOutTx: stopped")
			return
		}
	}
}

// return the status of txHash
// receipt nil, err non-nil: txHash not found
// receipt non-nil, err non-nil: txHash found but not confirmed
// receipt non-nil, err nil: txHash confirmed
func (ob *ChainObserver) queryTxByHash(txHash string, nonce int) (*ethtypes.Receipt, error) {
	//timeStart := time.Now()
	//defer func() { log.Info().Msgf("queryTxByHash elapsed: %s", time.Since(timeStart)) }()
	ctxt, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	receipt, err := ob.EvmClient.TransactionReceipt(ctxt, ethcommon.HexToHash(txHash))
	if err != nil {
		if err != ethereum.NotFound {
			log.Warn().Err(err).Msgf("%s %s TransactionReceipt err", ob.chain, txHash)
		}
		return nil, err
	} else if receipt.BlockNumber.Uint64()+ob.confCount > ob.LastBlock {
		log.Info().Msgf("%s TransactionReceipt %s mined in block %d but not confirmed; current block num %d", ob.chain, txHash, receipt.BlockNumber.Uint64(), ob.LastBlock)
		return receipt, err
	} else { // confirmed outbound tx
		if receipt.Status == 0 { // failed (reverted tx)
			log.Info().Msgf("%s TransactionReceipt %s nonce %d mined and confirmed, but it's reverted!", ob.chain, txHash, nonce)
		} else if receipt.Status == 1 { // success
			log.Info().Msgf("%s TransactionReceipt %s nonce %d mined and confirmed, and it's successful", ob.chain, txHash, nonce)
		}
		return receipt, nil
	}
}
