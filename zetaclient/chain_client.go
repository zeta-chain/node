package zetaclient

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"math"
	"math/big"
	"os"
	"path/filepath"
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
)

const (
	PosKey = "PosKey"

	SecondsPerDay = 86400
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
	uniswapV3Abi     *abi.ABI
	uniswapV2Abi     *abi.ABI
	Client           *ethclient.Client
	bridge           *MetachainBridge
	Tss              TSSSigner
	LastBlock        uint64
	confCount        uint64 // must wait this many blocks to be considered "confirmed"
	BlockTime        uint64 // block time in seconds
	txWatchList      map[ethcommon.Hash]string
	mu               *sync.Mutex
	db               *leveldb.DB
	sampleLoger      *zerolog.Logger
	metrics          *metrics.Metrics
	nonceTxHashesMap map[int][]string
	nonceTx          map[int]ethtypes.Receipt
	OutTxChan        chan OutTx // send to this channel if you want something back!

	getZetaExchangeRate func() (float64, error)
}

// Return configuration based on supplied target chain
func NewChainObserver(chain common.Chain, bridge *MetachainBridge, tss TSSSigner, dbpath string, metrics *metrics.Metrics) (*ChainObserver, error) {
	ob := ChainObserver{}
	ob.chain = chain
	ob.mu = &sync.Mutex{}
	sampled := log.Sample(&zerolog.BasicSampler{N: 10})
	ob.sampleLoger = &sampled
	ob.bridge = bridge
	ob.txWatchList = make(map[ethcommon.Hash]string)
	ob.Tss = tss
	ob.metrics = metrics
	ob.nonceTxHashesMap = make(map[int][]string)
	ob.OutTxChan = make(chan OutTx, 100)

	// create metric counters
	err := ob.RegisterPromCounter("rpc_getLogs_count", "Number of getLogs")
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

	// Initialize chain specific setup
	switch chain {
	case common.MumbaiChain:
		ob.mpiAddress = config.Chains[common.MumbaiChain.String()].ConnectorContractAddress
		ob.endpoint = config.MUMBAI_ENDPOINT
		ob.ticker = time.NewTicker(time.Duration(MaxInt(config.POLY_BLOCK_TIME, 12)) * time.Second)
		ob.confCount = config.POLYGON_CONFIRMATION_COUNT
		ob.uniswapV3Abi = &uniswapV3ABI
		ob.BlockTime = config.POLY_BLOCK_TIME

	case common.GoerliChain:
		ob.mpiAddress = config.Chains[common.GoerliChain.String()].ConnectorContractAddress
		ob.endpoint = config.GOERLI_ENDPOINT
		ob.ticker = time.NewTicker(time.Duration(MaxInt(config.ETH_BLOCK_TIME, 12)) * time.Second)
		ob.confCount = config.ETH_CONFIRMATION_COUNT
		ob.uniswapV3Abi = &uniswapV3ABI
		ob.BlockTime = config.ETH_BLOCK_TIME

	case common.BSCTestnetChain:
		ob.mpiAddress = config.Chains[common.BSCTestnetChain.String()].ConnectorContractAddress
		ob.endpoint = config.BSCTESTNET_ENDPOINT
		ob.ticker = time.NewTicker(time.Duration(MaxInt(config.BSC_BLOCK_TIME, 12)) * time.Second)
		ob.confCount = config.BSC_CONFIRMATION_COUNT
		ob.uniswapV2Abi = &uniswapV2ABI
		ob.BlockTime = config.BSC_BLOCK_TIME

	case common.RopstenChain:
		ob.mpiAddress = config.Chains[common.RopstenChain.String()].ConnectorContractAddress
		ob.endpoint = config.ROPSTEN_ENDPOINT
		ob.ticker = time.NewTicker(time.Duration(MaxInt(config.ROPSTEN_BLOCK_TIME, 12)) * time.Second)
		ob.confCount = config.ROPSTEN_CONFIRMATION_COUNT
		ob.uniswapV3Abi = &uniswapV3ABI
		ob.BlockTime = config.ROPSTEN_BLOCK_TIME

	}

	// Dial the mpiAddress
	log.Info().Msgf("Chain %s endpoint %s", ob.chain, ob.endpoint)
	client, err := ethclient.Dial(ob.endpoint)
	if err != nil {
		log.Err(err).Msg("eth Client Dial")
		return nil, err
	}
	ob.Client = client

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
				log.Error().Err(err).Msg("error writing ob.LastBlock to db: ")
			} else {
				ob.LastBlock, _ = binary.Uvarint(buf)
			}
		}

		{
			path := fmt.Sprintf("%s/%s.nonceTxHashesMap", dbpath, chain.String())
			jsonFile, err := os.Open(path)
			defer jsonFile.Close()
			if err != nil {
				log.Error().Err(err).Msgf("error opening %s", path)
			} else {
				dec := json.NewDecoder(jsonFile)
				err = dec.Decode(&ob.nonceTxHashesMap)
				if err != nil {
					log.Error().Err(err).Msgf("error opening %s", path)
				}
			}
		}

		{
			path := fmt.Sprintf("%s/%s.nonceTx", dbpath, chain.String())
			jsonFile, err := os.Open(path)
			defer jsonFile.Close()
			if err != nil {
				log.Error().Err(err).Msgf("error opening %s", path)
			} else {
				dec := json.NewDecoder(jsonFile)
				err = dec.Decode(&ob.nonceTx)
				if err != nil {
					log.Error().Err(err).Msgf("error opening %s", path)
				}
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
	for range ob.ticker.C {
		err := ob.observeChain()
		if err != nil {
			log.Err(err).Msg("observeChain error on " + ob.chain.String())
			continue
		}
	}
}

func (ob *ChainObserver) WatchGasPrice() {
	gasTicker := time.NewTicker(60 * time.Second)
	for range gasTicker.C {
		err := ob.PostGasPrice()
		if err != nil {
			log.Err(err).Msg("PostGasPrice error on " + ob.chain.String())
			continue
		}
	}
}

func (ob *ChainObserver) WatchExchangeRate() {
	gasTicker := time.NewTicker(60 * time.Second)
	for range gasTicker.C {
		var price *big.Int
		var err error
		var bn uint64
		if ob.chain == common.GoerliChain || ob.chain == common.MumbaiChain || ob.chain == common.RopstenChain {
			price, bn, err = ob.GetZetaExchangeRateUniswapV3()
		} else if ob.chain == common.BSCTestnetChain {
			price, bn, err = ob.GetZetaExchangeRateUniswapV2()
		}
		if err != nil {
			log.Err(err).Msg("GetZetaExchangeRate error on " + ob.chain.String())
			continue
		}
		price_f, _ := big.NewFloat(0).SetInt(price).Float64()
		log.Info().Msgf("%s: gasAsset/zeta rate %f", ob.chain, price_f/1e18)
		priceInHex := fmt.Sprintf("0x%x", price)

		_, err = ob.bridge.PostZetaConversionRate(ob.chain, priceInHex, bn)
		if err != nil {
			log.Err(err).Msg("PostZetaConversionRate error on " + ob.chain.String())
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
	ob.sampleLoger.Info().Msgf("%s current block %d, querying from %d to %d, %d blocks left to catch up, watching MPI address %s", ob.chain, header.Number.Uint64(), ob.LastBlock+1, toBlock, int(toBlock)-int(confirmedBlockNum), ethcommon.HexToAddress(ob.mpiAddress))

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
	err := ob.db.Close()
	if err != nil {
		log.Error().Err(err).Msg("error closing db")
	}

	userDir, _ := os.UserHomeDir()
	dbpath := filepath.Join(userDir, ".zetaclient/chainobserver")
	{
		path := fmt.Sprintf("%s/%s.nonceTxHashesMap", dbpath, ob.chain.String())
		log.Info().Msgf("writing to %s", path)
		jsonFile, err := os.Open(path)
		defer jsonFile.Close()
		if err != nil {
			log.Error().Err(err).Msgf("error opening %s", path)
		} else {
			enc := json.NewEncoder(jsonFile)
			err = enc.Encode(ob.nonceTxHashesMap)
			if err != nil {
				log.Error().Err(err).Msgf("error opening %s", path)
			}
		}
	}
	{
		path := fmt.Sprintf("%s/%s.nonceTx", dbpath, ob.chain.String())
		log.Info().Msgf("writing to %s", path)
		jsonFile, err := os.Open(path)
		defer jsonFile.Close()
		if err != nil {
			log.Error().Err(err).Msgf("error opening %s", path)
		} else {
			enc := json.NewEncoder(jsonFile)
			err = enc.Encode(ob.nonceTx)
			if err != nil {
				log.Error().Err(err).Msgf("error opening %s", path)
			}
		}
	}
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
				fmt.Printf("Found sendHash %s on chain %s\n", sendHash, ob.chain)
				retval, err := ob.connectorAbi.Unpack("ZetaReceived", vLog.Data)
				if err != nil {
					fmt.Println("error unpacking ZetaReceived")
					continue
				}
				fmt.Printf("Topic 0 (event hash): %s\n", vLog.Topics[0])
				fmt.Printf("Topic 1 (origin chain id): %d\n", vLog.Topics[1])
				fmt.Printf("Topic 2 (dest address): %s\n", vLog.Topics[2])
				fmt.Printf("Topic 3 (sendHash): %s\n", vLog.Topics[3])
				fmt.Printf("txhash: %s, blocknum %d\n", vLog.TxHash, vLog.BlockNumber)

				if vLog.BlockNumber+config.ETH_CONFIRMATION_COUNT < ob.LastBlock {
					fmt.Printf("Confirmed! Sending PostConfirmation to zetacore...\n")
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
						log.Err(err).Msg("error posting confirmation to meta core")
						continue
					}
					fmt.Printf("Zeta tx hash: %s\n", metaHash)
					return true, true, nil
				} else {
					fmt.Printf("Included in block but not yet confirmed! included in block %d, current block %d\n", vLog.BlockNumber, ob.LastBlock)
					return true, false, nil
				}
			case logZetaRevertedSignatureHash.Hex():
				fmt.Printf("Found (revert tx) sendHash %s on chain %s\n", sendHash, ob.chain)
				retval, err := ob.connectorAbi.Unpack("ZetaReverted", vLog.Data)
				if err != nil {
					fmt.Println("error unpacking ZetaReverted")
					continue
				}
				fmt.Printf("Topic 0 (event hash): %s\n", vLog.Topics[0])
				fmt.Printf("Topic 1 (dest chain id): %d\n", vLog.Topics[1])
				fmt.Printf("Topic 2 (dest address): %s\n", vLog.Topics[2])
				fmt.Printf("Topic 3 (sendHash): %s\n", vLog.Topics[3])
				fmt.Printf("txhash: %s, blocknum %d\n", vLog.TxHash, vLog.BlockNumber)

				if vLog.BlockNumber+config.ETH_CONFIRMATION_COUNT < ob.LastBlock {
					fmt.Printf("Confirmed! Sending PostConfirmation to zetacore...\n")
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
					fmt.Printf("Zeta tx hash: %s\n", metaHash)
					return true, true, nil
				} else {
					fmt.Printf("Included in block but not yet confirmed! included in block %d, current block %d\n", vLog.BlockNumber, ob.LastBlock)
					return true, false, nil
				}
			}
		}
	} else if found && receipt.Status == 0 {
		//FIXME: check nonce here by getTransaction RPC
		zetaTxHash, err := ob.bridge.PostReceiveConfirmation(sendHash, receipt.TxHash.Hex(), receipt.BlockNumber.Uint64(), "", common.ReceiveStatus_Failed, ob.chain.String())
		if err != nil {
			log.Error().Err(err).Msgf("PostReceiveConfirmation error in WatchTxHashWithTimeout; zeta tx hash %s", zetaTxHash)
		}
		return true, true, nil
	}

	return false, false, fmt.Errorf("IsSendOutTxProcessed: unknown chain %s", ob.chain)
}

// return the ratio GAS(ETH, BNB, MATIC, etc)/ZETA from Uniswap v3
// return price (gasasset/zeta), blockNum, error
func (ob *ChainObserver) GetZetaExchangeRateUniswapV3() (*big.Int, uint64, error) {
	TIME_WINDOW := 600 // time weighted average price over last 10min (600s) period
	input, err := ob.uniswapV3Abi.Pack("observe", []uint32{0, uint32(TIME_WINDOW)})
	if err != nil {
		return nil, 0, fmt.Errorf("fail to pack observe")
	}

	fromAddr := ethcommon.HexToAddress(config.TSS_TEST_ADDRESS)
	toAddr := ethcommon.HexToAddress(config.Chains[ob.chain.String()].PoolContractAddress)
	res, err := ob.Client.CallContract(context.TODO(), ethereum.CallMsg{
		From: fromAddr,
		To:   &toAddr,
		Data: input,
	}, nil)
	if err != nil {
		log.Err(err).Msgf("%s CallContract error", ob.chain)
		return nil, 0, err
	}
	bn, err := ob.Client.BlockNumber(context.TODO())
	if err != nil {
		log.Err(err).Msgf("%s BlockNumber error", ob.chain)
		return nil, 0, err
	}
	output, err := ob.uniswapV3Abi.Unpack("observe", res)
	if err != nil || len(output) != 2 {
		log.Err(err).Msgf("%s Unpack error or len(output) (%d) != 2", ob.chain, len(output))
		return nil, 0, err
	}
	cumTicks := *abi.ConvertType(output[0], new([2]*big.Int)).(*[2]*big.Int)
	tickDiff := big.NewInt(0).Div(big.NewInt(0).Sub(cumTicks[0], cumTicks[1]), big.NewInt(int64(TIME_WINDOW)))
	price := math.Pow(1.0001, float64(tickDiff.Int64())) * 1e18 // price is fixed point with decimal 18
	v, _ := big.NewFloat(price).Int(nil)
	return v, bn, nil
}

// return the ratio GAS(ETH, BNB, MATIC, etc)/ZETA from Uniswap v2 and its clone
// return price (gasasset/zeta), blockNum, error
func (ob *ChainObserver) GetZetaExchangeRateUniswapV2() (*big.Int, uint64, error) {
	input, err := ob.uniswapV2Abi.Pack("getReserves")
	if err != nil {
		return nil, 0, fmt.Errorf("fail to pack getReserves")
	}

	fromAddr := ethcommon.HexToAddress(config.TSS_TEST_ADDRESS)
	toAddr := ethcommon.HexToAddress(config.Chains[ob.chain.String()].PoolContractAddress)
	res, err := ob.Client.CallContract(context.TODO(), ethereum.CallMsg{
		From: fromAddr,
		To:   &toAddr,
		Data: input,
	}, nil)
	if err != nil {
		log.Err(err).Msgf("%s CallContract error", ob.chain)
		return nil, 0, err
	}
	bn, err := ob.Client.BlockNumber(context.TODO())
	if err != nil {
		log.Err(err).Msgf("%s BlockNumber error", ob.chain)
		return nil, 0, err
	}
	output, err := ob.uniswapV2Abi.Unpack("getReserves", res)
	if err != nil || len(output) != 3 {
		log.Err(err).Msgf("%s Unpack error or len(output) (%d) != 3", ob.chain, len(output))
		return nil, 0, err
	}
	reserve0 := *abi.ConvertType(output[0], new(*big.Int)).(**big.Int)
	reserve1 := *abi.ConvertType(output[1], new(*big.Int)).(**big.Int)
	r0, acc0 := big.NewFloat(0).SetInt(reserve0).Float64()
	r1, acc1 := big.NewFloat(0).SetInt(reserve1).Float64()

	if r0 <= 0 || r1 <= 0 || acc0 != big.Exact || acc1 != big.Exact {
		log.Err(err).Msgf("%s inexact conversion acc0=%s acc1=%s r0=%d r1=%d", ob.chain, acc0, acc1, reserve0, reserve1)
		return nil, 0, err
	}
	v, _ := big.NewFloat(r0 / r1 * 1.0e18).Int(nil)
	return v, bn, nil
}

// this function periodically checks all the txhash in potential outbound txs
func (ob *ChainObserver) observeOutTx() {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		select {
		case outTx := <-ob.OutTxChan:
			ob.nonceTxHashesMap[outTx.Nonce] = append(ob.nonceTxHashesMap[outTx.Nonce], outTx.TxHash)
		default:
			for nonce, txhashes := range ob.nonceTxHashesMap {
				log.Info().Msgf("observeOutTx: %s nonce %d, len %d", ob.chain, nonce, len(txhashes))
				for _, txhash := range txhashes {
					receipt, err := ob.queryTxByHash(txhash, nonce)
					if err == nil { // confirmed
						log.Info().Msgf("observeOutTx: %s nonce %d, txhash %s confirmed", ob.chain, nonce, txhash)
						delete(ob.nonceTxHashesMap, nonce)
						ob.nonceTx[nonce] = *receipt
						break
					}
					time.Sleep(1 * time.Second)
				}
			}
		}
	}
}

// return the status of txHash
// receipt nil, err non-nil: txHash not found
// receipt non-nil, err non-nil: txHash found but not confirmed
// receipt non-nil, err nil: txHash confirmed
func (ob *ChainObserver) queryTxByHash(txHash string, nonce int) (*ethtypes.Receipt, error) {
	receipt, err := ob.Client.TransactionReceipt(context.TODO(), ethcommon.HexToHash(txHash))
	if err != nil {
		log.Warn().Err(err).Msgf("%s %s TransactionReceipt err", ob.chain, txHash)
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
	ob.OutTxChan <- OutTx{
		TxHash:   txHash,
		Nonce:    nonce,
		SendHash: sendHash,
	}
}
