package zetaclient

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"math"
	"math/big"
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

// Chain configuration struct
// Filled with above constants depending on chain
type ChainObserver struct {
	chain        common.Chain
	mpiAddress   string
	endpoint     string
	ticker       *time.Ticker
	connectorAbi *abi.ABI // token contract ABI on non-ethereum chain; zetalocker on ethereum
	uniswapV3Abi *abi.ABI
	uniswapV2Abi *abi.ABI
	//zetaAbi     *connectorAbi.ABI // only useful for ethereum; the token contract
	Client      *ethclient.Client
	bridge      *MetachainBridge
	Tss         TSSSigner
	LastBlock   uint64
	confCount   uint64 // must wait this many blocks to be considered "confirmed"
	txWatchList map[ethcommon.Hash]string
	mu          *sync.Mutex
	db          *leveldb.DB
	sampleLoger *zerolog.Logger
	metrics     *metrics.Metrics

	getZetaExchangeRate func() (float64, error)
}

// Return configuration based on supplied target chain
func NewChainObserver(chain common.Chain, bridge *MetachainBridge, tss TSSSigner, dbpath string, metrics *metrics.Metrics) (*ChainObserver, error) {
	chainOb := ChainObserver{}
	chainOb.chain = chain
	chainOb.mu = &sync.Mutex{}
	sampled := log.Sample(&zerolog.BasicSampler{N: 10})
	chainOb.sampleLoger = &sampled
	chainOb.bridge = bridge
	chainOb.txWatchList = make(map[ethcommon.Hash]string)
	chainOb.Tss = tss
	chainOb.metrics = metrics

	// create metric counters
	err := chainOb.RegisterPromCounter("rpc_getLogs_count", "Number of getLogs")
	if err != nil {
		return nil, err
	}

	// initialize the pool ABI
	mpiABI, err := abi.JSON(strings.NewReader(config.CONNECTOR_ABI_STRING))
	if err != nil {
		return nil, err
	}
	chainOb.connectorAbi = &mpiABI
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
		chainOb.mpiAddress = config.Chains[common.MumbaiChain.String()].ConnectorContractAddress
		chainOb.endpoint = config.MUMBAI_ENDPOINT
		chainOb.ticker = time.NewTicker(time.Duration(MaxInt(config.POLY_BLOCK_TIME, 12)) * time.Second)
		chainOb.confCount = config.POLYGON_CONFIRMATION_COUNT
		chainOb.uniswapV3Abi = &uniswapV3ABI

	case common.GoerliChain:
		chainOb.mpiAddress = config.Chains[common.GoerliChain.String()].ConnectorContractAddress
		chainOb.endpoint = config.GOERLI_ENDPOINT
		chainOb.ticker = time.NewTicker(time.Duration(MaxInt(config.ETH_BLOCK_TIME, 12)) * time.Second)
		chainOb.confCount = config.ETH_CONFIRMATION_COUNT
		chainOb.uniswapV3Abi = &uniswapV3ABI

	case common.BSCTestnetChain:
		chainOb.mpiAddress = config.Chains[common.BSCTestnetChain.String()].ConnectorContractAddress
		chainOb.endpoint = config.BSCTESTNET_ENDPOINT
		chainOb.ticker = time.NewTicker(time.Duration(MaxInt(config.BSC_BLOCK_TIME, 12)) * time.Second)
		chainOb.confCount = config.BSC_CONFIRMATION_COUNT
		chainOb.uniswapV2Abi = &uniswapV2ABI

	case common.RopstenChain:
		chainOb.mpiAddress = config.Chains[common.RopstenChain.String()].ConnectorContractAddress
		chainOb.endpoint = config.ROPSTEN_ENDPOINT
		chainOb.ticker = time.NewTicker(time.Duration(MaxInt(config.ROPSTEN_BLOCK_TIME, 12)) * time.Second)
		chainOb.confCount = config.ROPSTEN_CONFIRMATION_COUNT
		chainOb.uniswapV3Abi = &uniswapV3ABI
	}

	// Dial the mpiAddress
	log.Info().Msgf("Chain %s endpoint %s", chainOb.chain, chainOb.endpoint)
	client, err := ethclient.Dial(chainOb.endpoint)
	if err != nil {
		log.Err(err).Msg("eth Client Dial")
		return nil, err
	}
	chainOb.Client = client

	if dbpath != "" {
		path := fmt.Sprintf("%s/%s", dbpath, chain.String()) // e.g. ~/.zetaclient/ETH
		db, err := leveldb.OpenFile(path, nil)

		if err != nil {
			return nil, err
		}
		chainOb.db = db
		buf, err := db.Get([]byte(PosKey), nil)
		if err != nil {
			log.Info().Msg("db PosKey does not exist; read from ZetaCore")
			chainOb.LastBlock = chainOb.getLastHeight()
			// if ZetaCore does not have last heard block height, then use current
			if chainOb.LastBlock == 0 {
				header, err := chainOb.Client.HeaderByNumber(context.Background(), nil)
				if err != nil {
					return nil, err
				}
				chainOb.LastBlock = header.Number.Uint64()
			}
			buf2 := make([]byte, binary.MaxVarintLen64)
			n := binary.PutUvarint(buf2, chainOb.LastBlock)
			err := db.Put([]byte(PosKey), buf2[:n], nil)
			log.Error().Err(err).Msg("error writing chainOb.LastBlock to db: ")
		} else {
			chainOb.LastBlock, _ = binary.Uvarint(buf)
		}
	}
	log.Info().Msgf("%s: start scanning from block %d", chain, chainOb.LastBlock)

	// this is shared structure to query logs by sendHash
	log.Info().Msgf("Chain %s logZetaReceivedSignatureHash %s", chainOb.chain, logZetaReceivedSignatureHash.Hex())

	return &chainOb, nil
}

func (chainOb *ChainObserver) GetPromCounter(name string) (prometheus.Counter, error) {
	if cnt, found := metrics.Counters[chainOb.chain.String()+name]; found {
		return cnt, nil
	} else {
		return nil, errors.New("counter not found")
	}
}

func (chainOb *ChainObserver) RegisterPromCounter(name string, help string) error {
	cntName := chainOb.chain.String() + name
	return chainOb.metrics.RegisterCounter(cntName, help)
}

func (chainOb *ChainObserver) Start() {
	go chainOb.WatchRouter()
	go chainOb.WatchGasPrice()
	go chainOb.WatchExchangeRate()
}

func (chainOb *ChainObserver) PostNonceIfNotRecorded() error {
	bridge := chainOb.bridge
	client := chainOb.Client
	tss := chainOb.Tss
	chain := chainOb.chain

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

func (chainOb *ChainObserver) WatchRouter() {
	// At each tick, query the mpiAddress
	for range chainOb.ticker.C {
		err := chainOb.observeChain()
		if err != nil {
			log.Err(err).Msg("observeChain error on " + chainOb.chain.String())
			continue
		}
	}
}

func (chainOb *ChainObserver) WatchGasPrice() {
	gasTicker := time.NewTicker(60 * time.Second)
	for range gasTicker.C {
		err := chainOb.PostGasPrice()
		if err != nil {
			log.Err(err).Msg("PostGasPrice error on " + chainOb.chain.String())
			continue
		}
	}
}

func (chainOb *ChainObserver) WatchExchangeRate() {
	gasTicker := time.NewTicker(60 * time.Second)
	for range gasTicker.C {
		var price *big.Int
		var err error
		var bn uint64
		if chainOb.chain == common.GoerliChain || chainOb.chain == common.MumbaiChain || chainOb.chain == common.RopstenChain {
			price, bn, err = chainOb.GetZetaExchangeRateUniswapV3()
		} else if chainOb.chain == common.BSCTestnetChain {
			price, bn, err = chainOb.GetZetaExchangeRateUniswapV2()
		}
		if err != nil {
			log.Err(err).Msg("GetZetaExchangeRate error on " + chainOb.chain.String())
			continue
		}
		price_f, _ := big.NewFloat(0).SetInt(price).Float64()
		log.Info().Msgf("%s: gasAsset/zeta rate %f", chainOb.chain, price_f/1e18)
		priceInHex := fmt.Sprintf("0x%x", price)

		_, err = chainOb.bridge.PostZetaConversionRate(chainOb.chain, priceInHex, bn)
		if err != nil {
			log.Err(err).Msg("PostZetaConversionRate error on " + chainOb.chain.String())
		}
	}
}

func (chainOb *ChainObserver) PostGasPrice() error {
	// GAS PRICE
	gasPrice, err := chainOb.Client.SuggestGasPrice(context.TODO())
	if err != nil {
		log.Err(err).Msg("PostGasPrice:")
		return err
	}
	blockNum, err := chainOb.Client.BlockNumber(context.TODO())
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

	_, err = chainOb.bridge.PostGasPrice(chainOb.chain, gasPrice.Uint64(), supply, blockNum)
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

func (chainOb *ChainObserver) observeChain() error {
	header, err := chainOb.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}
	// "confirmed" current block number
	confirmedBlockNum := header.Number.Uint64() - chainOb.confCount
	// skip if no new block is produced.
	if confirmedBlockNum <= chainOb.LastBlock {
		return nil
	}
	toBlock := chainOb.LastBlock + config.MAX_BLOCKS_PER_PERIOD // read at most 10 blocks in one go
	if toBlock >= confirmedBlockNum {
		toBlock = confirmedBlockNum
	}

	topics[0] = []ethcommon.Hash{logZetaSentSignatureHash}

	query := ethereum.FilterQuery{
		Addresses: []ethcommon.Address{ethcommon.HexToAddress(chainOb.mpiAddress)},
		FromBlock: big.NewInt(0).SetUint64(chainOb.LastBlock + 1), // LastBlock has been processed;
		ToBlock:   big.NewInt(0).SetUint64(toBlock),
		Topics:    topics,
	}
	//log.Debug().Msgf("signer %s block from %d to %d", chainOb.bridge.GetKeys().signerName, query.FromBlock, query.ToBlock)
	chainOb.sampleLoger.Info().Msgf("%s current block %d, querying from %d to %d, %d blocks left to catch up, watching MPI address %s", chainOb.chain, header.Number.Uint64(), chainOb.LastBlock+1, toBlock, int(toBlock)-int(confirmedBlockNum), ethcommon.HexToAddress(chainOb.mpiAddress))

	// Finally query the for the logs
	logs, err := chainOb.Client.FilterLogs(context.Background(), query)
	if err != nil {
		return err
	}
	cnt, err := chainOb.GetPromCounter("rpc_getLogs_count")
	if err != nil {
		return err
	}
	cnt.Inc()

	// Read in ABI
	contractAbi := chainOb.connectorAbi

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

			metaHash, err := chainOb.bridge.PostSend(
				ethcommon.HexToAddress(sender.Hex()).Hex(),
				chainOb.chain.String(),
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

	chainOb.LastBlock = toBlock
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, toBlock)
	err = chainOb.db.Put([]byte(PosKey), buf[:n], nil)
	if err != nil {
		log.Error().Err(err).Msg("error writing toBlock to db")
	}
	return nil
}

// query ZetaCore about the last block that it has heard from a specific chain.
// return 0 if not existent.
func (chainOb *ChainObserver) getLastHeight() uint64 {
	lastheight, err := chainOb.bridge.GetLastBlockHeightByChain(chainOb.chain)
	if err != nil {
		log.Warn().Err(err).Msgf("getLastHeight")
		return 0
	}
	return lastheight.LastSendHeight
}

// query the base gas price for the block number bn.
func (chainOb *ChainObserver) GetBaseGasPrice() *big.Int {
	gasPrice, err := chainOb.Client.SuggestGasPrice(context.TODO())
	if err != nil {
		log.Err(err).Msg("GetBaseGasPrice")
		return nil
	}
	return gasPrice
}

func (chainOb *ChainObserver) Stop() error {
	return chainOb.db.Close()
}

// returns: isIncluded, isConfirmed, Error
// If isConfirmed, it also post to ZetaCore
func (chainOb *ChainObserver) IsSendOutTxProcessed(sendHash string) (bool, bool, error) {
	recvTopics := make([][]ethcommon.Hash, 4)
	recvTopics[3] = []ethcommon.Hash{ethcommon.HexToHash(sendHash)}
	recvTopics[0] = []ethcommon.Hash{logZetaReceivedSignatureHash, logZetaRevertedSignatureHash}
	//fromBlock := big.NewInt(int64(chainOb.LastBlock - 3*SecondsPerDay/config.Chains[chainOb.chain.String()].BlockTime))
	query := ethereum.FilterQuery{
		Addresses: []ethcommon.Address{ethcommon.HexToAddress(config.Chains[chainOb.chain.String()].ConnectorContractAddress)},
		FromBlock: big.NewInt(0), // LastBlock from 3 days ago
		ToBlock:   nil,
		Topics:    recvTopics,
	}
	//log.Info().Msgf("%s getLogs: from %d to %d", chainOb.chain, fromBlock, chainOb.LastBlock)
	logs, err := chainOb.Client.FilterLogs(context.Background(), query)
	if err != nil {
		return false, false, fmt.Errorf("[%s] IsSendOutTxProcessed(sendHash %s): Client FilterLog fail %w", chainOb.chain, sendHash, err)
	}
	cnt, err := chainOb.GetPromCounter("rpc_getLogs_count")
	if err != nil {
		log.Error().Err(err).Msg("prometheus counter error")
	}
	cnt.Inc()

	if len(logs) == 0 {
		return false, false, nil
	}
	if len(logs) > 1 {
		log.Fatal().Msgf("More than two logs with send hash %s", sendHash)
		log.Fatal().Msgf("First one: %+v\nSecond one:%+v\n", logs[0], logs[1])
	}
	for _, vLog := range logs {
		switch vLog.Topics[0].Hex() {
		case logZetaReceivedSignatureHash.Hex():
			fmt.Printf("Found sendHash %s on chain %s\n", sendHash, chainOb.chain)
			retval, err := chainOb.connectorAbi.Unpack("ZetaReceived", vLog.Data)
			if err != nil {
				fmt.Println("error unpacking ZetaReceived")
				continue
			}
			fmt.Printf("Topic 0 (event hash): %s\n", vLog.Topics[0])
			fmt.Printf("Topic 1 (origin chain id): %d\n", vLog.Topics[1])
			fmt.Printf("Topic 2 (dest address): %s\n", vLog.Topics[2])
			fmt.Printf("Topic 3 (sendHash): %s\n", vLog.Topics[3])
			fmt.Printf("txhash: %s, blocknum %d\n", vLog.TxHash, vLog.BlockNumber)

			if vLog.BlockNumber+config.ETH_CONFIRMATION_COUNT < chainOb.LastBlock {
				fmt.Printf("Confirmed! Sending PostConfirmation to zetacore...\n")
				sendhash := vLog.Topics[3].Hex()
				//var rxAddress string = ethcommon.HexToAddress(vLog.Topics[1].Hex()).Hex()
				var mMint string = retval[1].(*big.Int).String()
				metaHash, err := chainOb.bridge.PostReceiveConfirmation(
					sendhash,
					vLog.TxHash.Hex(),
					vLog.BlockNumber,
					mMint,
					common.ReceiveStatus_Success,
					chainOb.chain.String(),
				)
				if err != nil {
					log.Err(err).Msg("error posting confirmation to meta core")
					continue
				}
				fmt.Printf("Zeta tx hash: %s\n", metaHash)
				return true, true, nil
			} else {
				fmt.Printf("Included in block but not yet confirmed! included in block %d, current block %d\n", vLog.BlockNumber, chainOb.LastBlock)
				return true, false, nil
			}
		case logZetaRevertedSignatureHash.Hex():
			fmt.Printf("Found (revert tx) sendHash %s on chain %s\n", sendHash, chainOb.chain)
			retval, err := chainOb.connectorAbi.Unpack("ZetaReverted", vLog.Data)
			if err != nil {
				fmt.Println("error unpacking ZetaReverted")
				continue
			}
			fmt.Printf("Topic 0 (event hash): %s\n", vLog.Topics[0])
			fmt.Printf("Topic 1 (dest chain id): %d\n", vLog.Topics[1])
			fmt.Printf("Topic 2 (dest address): %s\n", vLog.Topics[2])
			fmt.Printf("Topic 3 (sendHash): %s\n", vLog.Topics[3])
			fmt.Printf("txhash: %s, blocknum %d\n", vLog.TxHash, vLog.BlockNumber)

			if vLog.BlockNumber+config.ETH_CONFIRMATION_COUNT < chainOb.LastBlock {
				fmt.Printf("Confirmed! Sending PostConfirmation to zetacore...\n")
				sendhash := vLog.Topics[3].Hex()
				var mMint string = retval[2].(*big.Int).String()
				metaHash, err := chainOb.bridge.PostReceiveConfirmation(
					sendhash,
					vLog.TxHash.Hex(),
					vLog.BlockNumber,
					mMint,
					common.ReceiveStatus_Success,
					chainOb.chain.String(),
				)
				if err != nil {
					log.Err(err).Msg("error posting confirmation to meta core")
					continue
				}
				fmt.Printf("Zeta tx hash: %s\n", metaHash)
				return true, true, nil
			} else {
				fmt.Printf("Included in block but not yet confirmed! included in block %d, current block %d\n", vLog.BlockNumber, chainOb.LastBlock)
				return true, false, nil
			}
		}
	}

	return false, false, fmt.Errorf("IsSendOutTxProcessed: unknown chain %s", chainOb.chain)
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

// watch outbound tx
// returns whether outbound tx is successful or failure
func (chainOb *ChainObserver) WatchTxHashWithTimeout(txid string, sendHash string) bool {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Minute)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msgf("TIMEOUT: watching outTx %s on chain %s", txid, chainOb.chain)
			return false
		default:
			receipt, err := chainOb.Client.TransactionReceipt(context.TODO(), ethcommon.HexToHash(txid))
			if err != nil {
				log.Error().Err(err).Msgf("error watching outTx %s on chain %s", txid, chainOb.chain)
			} else {
				if receipt.Status == 1 { // 1: success
					log.Info().Msgf("SUCCESS: watching outTx %s on chain %s", txid, chainOb.chain)
					return true
				} else if receipt.Status == 0 { // tx mined but failed; should revert
					log.Info().Msgf("FAILED: watching outTx %s on chain %s", txid, chainOb.chain)
					zetaTxHash, err := chainOb.bridge.PostReceiveConfirmation(sendHash, txid, receipt.BlockNumber.Uint64(), "", common.ReceiveStatus_Failed, chainOb.chain.String())
					if err != nil {
						log.Error().Err(err).Msgf("PostReceiveConfirmation error in WatchTxHashWithTimeout; zeta tx hash %s", zetaTxHash)
					}
					return false
				}
				return false
			}
			time.Sleep(10 * time.Second)
		}
	}
}
