package zetaclient

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/syndtr/goleveldb/leveldb/util"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/types"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	zetatypes "github.com/zeta-chain/zetacore/x/zetacore/types"
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
	mpiAddress       string
	endpoint         string
	ticker           *time.Ticker
	connectorAbi     *abi.ABI // token contract ABI on non-ethereum chain; zetalocker on ethereum
	Client           *ethclient.Client
	bridge           *ZetaCoreBridge
	Tss              TSSSigner
	LastBlock        uint64
	confCount        uint64 // must wait this many blocks to be considered "confirmed"
	BlockTime        uint64 // block time in seconds
	txWatchList      map[ethcommon.Hash]string
	mu               *sync.Mutex
	db               *leveldb.DB
	sampleLogger     *zerolog.Logger
	metrics          *metrics.Metrics
	nonceTxHashesMap map[int][]string
	nonceTx          map[int]*ethtypes.Receipt
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
	ob.bridge = bridge
	ob.txWatchList = make(map[ethcommon.Hash]string)
	ob.Tss = tss
	ob.metrics = metrics
	ob.nonceTxHashesMap = make(map[int][]string)
	ob.nonceTx = make(map[int]*ethtypes.Receipt)
	ob.OutTxChan = make(chan OutTx, 100)
	ob.mpiAddress = config.Chains[chain.String()].ConnectorContractAddress
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
		log.Err(err).Msg("eth Client Dial")
		return nil, err
	}
	ob.Client = client

	// create metric counters
	err = ob.RegisterPromCounter("rpc_getLogs_count", "Number of getLogs")
	if err != nil {
		return nil, err
	}
	err = ob.RegisterPromCounter("rpc_getBlockByNumber_count", "Number of getBlockByNumber")
	if err != nil {
		return nil, err
	}

	// initialize the pool ABI
	mpiABI, err := abi.JSON(strings.NewReader(config.CONNECTOR_ABI_STRING))
	if err != nil {
		return nil, err
	}
	ob.connectorAbi = &mpiABI
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
		Client:              ob.Client,
		PoolContractAddress: ethcommon.HexToAddress(config.Chains[chain.String()].PoolContractAddress),
		Chain:               ob.chain,
	}
	uniswapv2querier := &UniswapV2ZetaPriceQuerier{
		UniswapV2Abi:        &uniswapV2ABI,
		Client:              ob.Client,
		PoolContractAddress: ethcommon.HexToAddress(config.Chains[chain.String()].PoolContractAddress),
		Chain:               ob.chain,
	}
	dummyQuerier := &DummyZetaPriceQuerier{
		Chain:  ob.chain,
		Client: ob.Client,
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
			header, err := ob.Client.HeaderByNumber(context.Background(), nil)
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
					header, err := ob.Client.HeaderByNumber(context.Background(), nil)
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
				ob.nonceTxHashesMap[int(nonce)] = txHashes
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
				ob.nonceTx[int(nonce)] = &receipt
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

func (ob *ChainObserver) GetPromCounter(name string) (prometheus.Counter, error) {
	if cnt, found := metrics.Counters[ob.chain.String()+"_"+name]; found {
		return cnt, nil
	} else {
		return nil, errors.New("counter not found")
	}
}

func (ob *ChainObserver) RegisterPromCounter(name string, help string) error {
	cntName := ob.chain.String() + "_" + name
	return ob.metrics.RegisterCounter(cntName, help)
}

func (ob *ChainObserver) Start() {
	go ob.WatchRouter()
	go ob.WatchGasPrice()
	go ob.WatchExchangeRate()
	go ob.observeOutTx()
}

func (ob *ChainObserver) PostNonceIfNotRecorded() error {
	bridge := ob.bridge
	client := ob.Client
	tss := ob.Tss
	chain := ob.chain

	_, err := bridge.GetNonceByChain(chain)
	if err != nil { // if Nonce of Chain is not found in ZetaCore; report it
		nonce, err := client.NonceAt(context.TODO(), tss.Address(), nil)
		if err != nil {
			log.Err(err).Msg("NonceAt")
			return err
		}
		log.Debug().Msgf("signer %s Posting Nonce of chain %s of nonce %d", bridge.GetKeys().signerName, chain, nonce)
		_, err = bridge.PostNonce(chain, nonce)
		if err != nil {
			log.Err(err).Msg("PostNonce")
			return err
		}
	}

	return nil
}

func (ob *ChainObserver) WatchRouter() {
	// At each tick, query the mpiAddress
	for {
		select {
		case <-ob.ticker.C:
			err := ob.observeChain()
			if err != nil {
				log.Err(err).Msg("observeChain error on " + ob.chain.String())
				continue
			}
		case <-ob.stop:
			log.Info().Msg("WatchRouter stopped")
			return
		}
	}
}

func (ob *ChainObserver) WatchGasPrice() {
	gasTicker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-gasTicker.C:
			err := ob.PostGasPrice()
			if err != nil {
				log.Err(err).Msg("PostGasPrice error on " + ob.chain.String())
				continue
			}
		case <-ob.stop:
			log.Info().Msg("WatchGasPrice stopped")
			return
		}
	}
}

func (ob *ChainObserver) WatchExchangeRate() {
	ticker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-ticker.C:
			price, bn, err := ob.ZetaPriceQuerier.GetZetaPrice()
			if err != nil {
				log.Err(err).Msg("GetZetaExchangeRate error on " + ob.chain.String())
				continue
			}
			priceInHex := fmt.Sprintf("0x%x", price)

			_, err = ob.bridge.PostZetaConversionRate(ob.chain, priceInHex, bn)
			if err != nil {
				log.Err(err).Msg("PostZetaConversionRate error on " + ob.chain.String())
			}
		case <-ob.stop:
			log.Info().Msg("WatchExchangeRate stopped")
			return
		}
	}
}

func (ob *ChainObserver) PostGasPrice() error {
	// GAS PRICE
	gasPrice, err := ob.Client.SuggestGasPrice(context.TODO())
	if err != nil {
		log.Err(err).Msg("PostGasPrice:")
		return err
	}
	blockNum, err := ob.Client.BlockNumber(context.TODO())
	if err != nil {
		log.Err(err).Msg("PostGasPrice:")
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

	_, err = ob.bridge.PostGasPrice(ob.chain, gasPrice.Uint64(), supply, blockNum)
	if err != nil {
		log.Err(err).Msg("PostGasPrice:")
		return err
	}

	//bal, err := chainOb.Client.BalanceAt(context.TODO(), chainOb.Tss.Address(), nil)
	//if err != nil {
	//	log.Err(err).Msg("BalanceAt:")
	//	return err
	//}
	//_, err = chainOb.bridge.PostGasBalance(chainOb.chain, bal.String(), blockNum)
	//if err != nil {
	//	log.Err(err).Msg("PostGasBalance:")
	//	return err
	//}
	return nil
}

func (ob *ChainObserver) observeChain() error {
	header, err := ob.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}
	counter, err := ob.GetPromCounter("rpc_getBlockByNumber_count")
	if err != nil {
		log.Error().Err(err).Msg("GetPromCounter:")
	}
	counter.Inc()

	// "confirmed" current block number
	confirmedBlockNum := header.Number.Uint64() - ob.confCount
	// skip if no new block is produced.
	if confirmedBlockNum <= ob.LastBlock {
		return nil
	}
	toBlock := ob.LastBlock + config.MAX_BLOCKS_PER_PERIOD // read at most 10 blocks in one go
	if toBlock >= confirmedBlockNum {
		toBlock = confirmedBlockNum
	}

	topics[0] = []ethcommon.Hash{logZetaSentSignatureHash}

	query := ethereum.FilterQuery{
		Addresses: []ethcommon.Address{ethcommon.HexToAddress(ob.mpiAddress)},
		FromBlock: big.NewInt(0).SetUint64(ob.LastBlock + 1), // LastBlock has been processed;
		ToBlock:   big.NewInt(0).SetUint64(toBlock),
		Topics:    topics,
	}
	//log.Debug().Msgf("signer %s block from %d to %d", chainOb.bridge.GetKeys().signerName, query.FromBlock, query.ToBlock)
	ob.sampleLogger.Info().Msgf("%s current block %d, querying from %d to %d, %d blocks left to catch up, watching MPI address %s", ob.chain, header.Number.Uint64(), ob.LastBlock+1, toBlock, int(toBlock)-int(confirmedBlockNum), ethcommon.HexToAddress(ob.mpiAddress))

	// Finally query the for the logs
	logs, err := ob.Client.FilterLogs(context.Background(), query)
	if err != nil {
		return err
	}
	cnt, err := ob.GetPromCounter("rpc_getLogs_count")
	if err != nil {
		return err
	}
	cnt.Inc()

	// Read in ABI
	contractAbi := ob.connectorAbi

	// Pull out arguments from logs
	for _, vLog := range logs {
		log.Info().Msgf("TxBlockNumber %d Transaction Hash: %s topic %s\n", vLog.BlockNumber, vLog.TxHash.Hex()[:6], vLog.Topics[0].Hex()[:6])
		switch vLog.Topics[0].Hex() {
		case logZetaSentSignatureHash.Hex():
			vals, err := contractAbi.Unpack("ZetaSent", vLog.Data)
			if err != nil {
				log.Err(err).Msg("error unpacking ZetaMessageSendEvent")
				continue
			}
			sender := vLog.Topics[1]
			destChainID := vals[0].(*big.Int)
			destContract := vals[1].([]byte)
			zetaAmount := vals[2].(*big.Int)
			gasLimit := vals[3].(*big.Int)
			message := vals[4].([]byte)
			zetaParams := vals[5].([]byte)

			_ = zetaParams

			metaHash, err := ob.bridge.PostSend(
				ethcommon.HexToAddress(sender.Hex()).Hex(),
				ob.chain.String(),
				types.BytesToEthHex(destContract),
				config.FindChainByID(destChainID),
				zetaAmount.String(),
				zetaAmount.String(),
				base64.StdEncoding.EncodeToString(message),
				vLog.TxHash.Hex(),
				vLog.BlockNumber,
				gasLimit.Uint64(),
			)
			if err != nil {
				log.Err(err).Msg("error posting to meta core")
				continue
			}
			log.Debug().Msgf("LockSend detected: PostSend metahash: %s", metaHash)
		}
	}

	ob.LastBlock = toBlock
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, toBlock)
	err = ob.db.Put([]byte(PosKey), buf[:n], nil)
	if err != nil {
		log.Error().Err(err).Msg("error writing toBlock to db")
	}
	return nil
}

// query ZetaCore about the last block that it has heard from a specific chain.
// return 0 if not existent.
func (ob *ChainObserver) getLastHeight() uint64 {
	lastheight, err := ob.bridge.GetLastBlockHeightByChain(ob.chain)
	if err != nil {
		log.Warn().Err(err).Msgf("getLastHeight")
		return 0
	}
	return lastheight.LastSendHeight
}

// query the base gas price for the block number bn.
func (ob *ChainObserver) GetBaseGasPrice() *big.Int {
	gasPrice, err := ob.Client.SuggestGasPrice(context.TODO())
	if err != nil {
		log.Err(err).Msg("GetBaseGasPrice")
		return nil
	}
	return gasPrice
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
	receipt, found := ob.nonceTx[nonce]
	if found && receipt.Status == 1 {
		logs := receipt.Logs
		for _, vLog := range logs {
			switch vLog.Topics[0].Hex() {
			case logZetaReceivedSignatureHash.Hex():
				retval, err := ob.connectorAbi.Unpack("ZetaReceived", vLog.Data)
				if err != nil {
					log.Error().Err(err).Msg("error unpacking ZetaReceived")
					continue
				}

				if vLog.BlockNumber+ob.confCount < ob.LastBlock {
					log.Info().Msgf("Found (outTx) sendHash %s on chain %s txhash %s", sendHash, ob.chain, vLog.TxHash.Hex())
					log.Info().Msg("Confirmed! Sending PostConfirmation to zetacore...")
					sendhash := vLog.Topics[3].Hex()
					//var rxAddress string = ethcommon.HexToAddress(vLog.Topics[1].Hex()).Hex()
					var mMint string = retval[1].(*big.Int).String()
					metaHash, err := ob.bridge.PostReceiveConfirmation(
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
					log.Info().Msgf("Zeta tx hash: %s\n", metaHash)
					return true, true, nil
				} else {
					log.Info().Msgf("Included; %d blocks before confirmed! chain %s nonce %d", int(vLog.BlockNumber+ob.confCount)-int(ob.LastBlock), ob.chain, nonce)
					return true, false, nil
				}
			case logZetaRevertedSignatureHash.Hex():
				retval, err := ob.connectorAbi.Unpack("ZetaReverted", vLog.Data)
				if err != nil {
					log.Error().Err(err).Msg("error unpacking ZetaReverted")
					continue
				}

				if vLog.BlockNumber+ob.confCount < ob.LastBlock {
					log.Info().Msgf("Found (revertTx) sendHash %s on chain %s txhash %s", sendHash, ob.chain, vLog.TxHash.Hex())
					log.Info().Msg("Confirmed! Sending PostConfirmation to zetacore...")
					sendhash := vLog.Topics[3].Hex()
					var mMint string = retval[2].(*big.Int).String()
					metaHash, err := ob.bridge.PostReceiveConfirmation(
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
		zetaTxHash, err := ob.bridge.PostReceiveConfirmation(sendHash, receipt.TxHash.Hex(), receipt.BlockNumber.Uint64(), "", common.ReceiveStatus_Failed, ob.chain.String())
		if err != nil {
			log.Error().Err(err).Msgf("PostReceiveConfirmation error in WatchTxHashWithTimeout; zeta tx hash %s", zetaTxHash)
		}
		log.Info().Msgf("Zeta tx hash: %s", zetaTxHash)
		return true, true, nil
	}

	return false, false, fmt.Errorf("IsSendOutTxProcessed: error on chain %s", ob.chain)
}

// this function periodically checks all the txhash in potential outbound txs
// FIXME: there's a chance that a txhash in OutTxChan may not deliver when Stop() is called
func (ob *ChainObserver) observeOutTx() {

	ticker := time.NewTicker(12 * time.Second)
	for {
		select {
		case <-ticker.C:
			minNonce, maxNonce, err := ob.PurgeTxHashWatchList()
			if len(ob.nonceTxHashesMap) > 0 {
				log.Info().Msgf("chain %s outstanding nonce: %d; nonce range [%d,%d]", ob.chain, len(ob.nonceTxHashesMap), minNonce, maxNonce)
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
					txHashes, found := ob.nonceTxHashesMap[nonce]
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
								delete(ob.nonceTxHashesMap, nonce)
								if err = ob.db.Delete([]byte(NonceTxHashesKeyPrefix+fmt.Sprintf("%d", nonce)), nil); err != nil {
									log.Error().Err(err).Msgf("PurgeTxHashWatchList: error deleting nonce %d tx hashes from db", nonce)
								}
								ob.nonceTx[nonce] = receipt
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

// remove txhash from watch list which have no corresponding sendPending in zetacore.
// returns the min/max nonce after purge
func (ob *ChainObserver) PurgeTxHashWatchList() (int, int, error) {
	purgedTxHashCount := 0
	sends, err := ob.bridge.GetAllPendingSend()
	if err != nil {
		return 0, 0, err
	}
	pendingNonces := make(map[int]bool)
	for _, send := range sends {
		if send.Status == zetatypes.SendStatus_PendingRevert && send.SenderChain == ob.chain.String() {
			pendingNonces[int(send.Nonce)] = true
		} else if send.Status == zetatypes.SendStatus_PendingOutbound && send.ReceiverChain == ob.chain.String() {
			pendingNonces[int(send.Nonce)] = true
		}
	}
	tNow := time.Now()
	ob.mu.Lock()
	for nonce, _ := range ob.nonceTxHashesMap {
		if _, found := pendingNonces[nonce]; !found {
			txHashes := ob.nonceTxHashesMap[nonce]
			delete(ob.nonceTxHashesMap, nonce)
			if err = ob.db.Delete([]byte(NonceTxHashesKeyPrefix+fmt.Sprintf("%d", nonce)), nil); err != nil {
				log.Error().Err(err).Msgf("PurgeTxHashWatchList: error deleting nonce %d tx hashes from db", nonce)
			}
			purgedTxHashCount++
			log.Info().Msgf("PurgeTxHashWatchList: chain %s nonce %d removed", ob.chain, nonce)
			ob.fileLogger.Info().Msgf("PurgeTxHashWatchList: chain %s nonce %d removed txhashes %v", ob.chain, nonce, txHashes)
		}
	}
	ob.mu.Unlock()
	if purgedTxHashCount > 0 {
		log.Info().Msgf("PurgeTxHashWatchList: chain %s purged %d txhashes in %v", ob.chain, purgedTxHashCount, time.Since(tNow))
	}
	minNonce, maxNonce := -1, 0
	if len(pendingNonces) > 0 {
		for nonce, _ := range pendingNonces {
			if minNonce == -1 {
				minNonce = nonce
			}
			if nonce < minNonce {
				minNonce = nonce
			}
			if nonce > maxNonce {
				maxNonce = nonce
			}
		}
	}
	return minNonce, maxNonce, nil
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
	receipt, err := ob.Client.TransactionReceipt(ctxt, ethcommon.HexToHash(txHash))
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

func (ob *ChainObserver) AddTxHashToWatchList(txHash string, nonce int, sendHash string) {
	outTx := OutTx{
		TxHash:   txHash,
		Nonce:    nonce,
		SendHash: sendHash,
	}

	if outTx.TxHash != "" { // TODO: this seems unnecessary
		ob.mu.Lock()
		ob.nonceTxHashesMap[outTx.Nonce] = append(ob.nonceTxHashesMap[outTx.Nonce], outTx.TxHash)
		ob.mu.Unlock()
		key := []byte(NonceTxHashesKeyPrefix + fmt.Sprintf("%d", outTx.Nonce))
		value := []byte(strings.Join(ob.nonceTxHashesMap[outTx.Nonce], ","))
		if err := ob.db.Put(key, value, nil); err != nil {
			log.Error().Err(err).Msgf("AddTxHashToWatchList: error adding nonce %d tx hashes to db", outTx.Nonce)
		}

		log.Info().Msgf("add %s nonce %d TxHash watch list length: %d", ob.chain, outTx.Nonce, len(ob.nonceTxHashesMap[outTx.Nonce]))
		ob.fileLogger.Info().Msgf("add %s nonce %d TxHash watch list length: %d", ob.chain, outTx.Nonce, len(ob.nonceTxHashesMap[outTx.Nonce]))
	}
}
