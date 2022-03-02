package zetaclient

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/types"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	PosKey = "PosKey"
)

// DELETE
//var logLockSendSignature = []byte("LockSend(address,string,uint256,uint256,string,bytes)")
//var logLockSendSignatureHash = crypto.Keccak256Hash(logLockSendSignature)
//
//var logUnlockSignature = []byte("Unlock(address,uint256,bytes32)")
//var logUnlockSignatureHash = crypto.Keccak256Hash(logUnlockSignature)
//
//var logBurnSendSignature = []byte("BurnSend(address,string,uint256,uint256,string,bytes)")
//var logBurnSendSignatureHash = crypto.Keccak256Hash(logBurnSendSignature)
//
//var logMMintedSignature = []byte("MMinted(address,uint256,bytes32)")
//var logMMintedSignatureHash = crypto.Keccak256Hash(logMMintedSignature)

// END DELETE

//   event ZetaMessageSendEvent(address indexed sender, uint16 destChainID, bytes destContract, uint zetaAmount, uint gasLimit, bytes message, bytes zetaParams);
var logMPISendSignature = []byte("ZetaMessageSendEvent(address,uint16,bytes,uint256,uint256,bytes,bytes)")
var logMPISendSignatureHash = crypto.Keccak256Hash(logMPISendSignature)

var logMPIReceiveSignature = []byte("ZetaMessageReceiveEvent(bytes,uint16,address,uint256,bytes,bytes32)")
var logMPIReceiveSignatureHash = crypto.Keccak256Hash(logMPIReceiveSignature)

var topics = make([][]ethcommon.Hash, 1)

// Chain configuration struct
// Filled with above constants depending on chain
type ChainObserver struct {
	chain      common.Chain
	mpiAddress string
	endpoint   string
	ticker     *time.Ticker
	abiString  string
	abi        *abi.ABI // token contract ABI on non-ethereum chain; zetalocker on ethereum
	//zetaAbi     *abi.ABI // only useful for ethereum; the token contract
	Client      *ethclient.Client
	bridge      *MetachainBridge
	Tss         TSSSigner
	LastBlock   uint64
	confCount   uint64 // must wait this many blocks to be considered "confirmed"
	txWatchList map[ethcommon.Hash]string
	mu          *sync.Mutex
	db          *leveldb.DB
	sampleLoger *zerolog.Logger
}

// Return configuration based on supplied target chain
func NewChainObserver(chain common.Chain, bridge *MetachainBridge, tss TSSSigner, dbpath string) (*ChainObserver, error) {
	chainOb := ChainObserver{}
	chainOb.mu = &sync.Mutex{}
	sampled := log.Sample(&zerolog.BasicSampler{N: 10})
	chainOb.sampleLoger = &sampled
	chainOb.bridge = bridge
	chainOb.txWatchList = make(map[ethcommon.Hash]string)
	chainOb.Tss = tss

	ethEndPoint := os.Getenv("ETH_ENDPOINT")
	if ethEndPoint != "" {
		config.ETH_ENDPOINT = ethEndPoint
		log.Info().Msgf("ETH_ENDPOINT: %s", ethEndPoint)
	}
	bscEndPoint := os.Getenv("BSC_ENDPOINT")
	if bscEndPoint != "" {
		config.BSC_ENDPOINT = bscEndPoint
	}
	polygonEndPoint := os.Getenv("POLYGON_ENDPOINT")
	if polygonEndPoint != "" {
		config.POLY_ENDPOINT = polygonEndPoint
	}

	// Initialize constants
	switch chain {
	case common.POLYGONChain:
		chainOb.chain = chain
		chainOb.mpiAddress = config.POLYGON_TOKEN_ADDRESS
		chainOb.endpoint = config.POLY_ENDPOINT
		chainOb.ticker = time.NewTicker(time.Duration(config.POLY_BLOCK_TIME) * time.Second)
		chainOb.abiString = config.MPI_ABI_STRING
		chainOb.confCount = config.POLYGON_CONFIRMATION_COUNT

	case common.ETHChain:
		chainOb.chain = chain
		chainOb.mpiAddress = config.Chains["ETH"].MPIContractAddress
		chainOb.endpoint = config.ETH_ENDPOINT
		chainOb.ticker = time.NewTicker(time.Duration(config.ETH_BLOCK_TIME) * time.Second)
		chainOb.abiString = config.MPI_ABI_STRING
		chainOb.confCount = config.ETH_CONFIRMATION_COUNT

	case common.BSCChain:
		chainOb.chain = chain
		chainOb.mpiAddress = config.Chains["BSC"].MPIContractAddress
		chainOb.endpoint = config.BSC_ENDPOINT
		chainOb.ticker = time.NewTicker(time.Duration(config.BSC_BLOCK_TIME) * time.Second)
		chainOb.abiString = config.MPI_ABI_STRING
		chainOb.confCount = config.BSC_CONFIRMATION_COUNT
	}
	contractABI, err := abi.JSON(strings.NewReader(chainOb.abiString))
	if err != nil {
		return nil, err
	}
	chainOb.abi = &contractABI
	//if chain == common.ETHChain {
	//	tokenABI, err := abi.JSON(strings.NewReader(config.ETH_ZETA_ABI))
	//	if err != nil {
	//		return nil, err
	//	}
	//	chainOb.zetaAbi = &tokenABI
	//}

	// Dial the mpiAddress
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
			log.Info().Msg("db PosKey does not exsit; read from ZetaCore")
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
			db.Put([]byte(PosKey), buf2[:n], nil)
		} else {
			chainOb.LastBlock, _ = binary.Uvarint(buf)
		}
	}
	log.Info().Msgf("%s: start scanning from block %d", chain, chainOb.LastBlock)

	if bridge != nil {
		_, err = bridge.GetNonceByChain(chain)
		if err != nil { // if Nonce of Chain is not found in ZetaCore; report it
			nonce, err := client.NonceAt(context.TODO(), tss.Address(), nil)
			if err != nil {
				log.Err(err).Msg("NonceAt")
				return nil, err
			}
			log.Debug().Msgf("signer %s Posting Nonce of chain %s of nonce %d", bridge.GetKeys().signerName, chain, nonce)
			_, err = bridge.PostNonce(chain, nonce)
			if err != nil {
				log.Err(err).Msg("PostNonce")
				return nil, err
			}
		}
	}

	// this is shared structure to query logs by sendHash
	//topics[2] = make([]ethcommon.Hash, 1)

	return &chainOb, nil
}

func (chainOb *ChainObserver) WatchRouter() {
	// At each tick, query the mpiAddress
	for range chainOb.ticker.C {
		err := chainOb.observeChain()
		if err != nil {
			log.Err(err).Msg("observeChain error")
			continue
		}
		chainOb.observeFailedTx()
	}
}

func (chainOb *ChainObserver) WatchGasPrice() {
	for range chainOb.ticker.C {
		err := chainOb.PostGasPrice()
		if err != nil {
			log.Err(err).Msg("PostGasPrice error")
			continue
		}
	}
}

func (chainOb *ChainObserver) AddTxToWatchList(txhash string, sendhash string) {
	hash := ethcommon.HexToHash(txhash)
	chainOb.mu.Lock()
	chainOb.txWatchList[hash] = sendhash
	chainOb.mu.Unlock()
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
	//	input, err := chainOb.abi.Pack("getLockedAmount")
	//	if err != nil {
	//		return fmt.Errorf("fail to getLockedAmount")
	//	}
	//	bn, err := chainOb.Client.BlockNumber(context.TODO())
	//	if err != nil {
	//		log.Err(err).Msgf("%s BlockNumber error", chainOb.chain)
	//		return err
	//	}
	//	fromAddr := ethcommon.HexToAddress(config.TSS_TEST_ADDRESS)
	//	toAddr := ethcommon.HexToAddress(config.ETH_ZETALOCK_ADDRESS)
	//	res, err := chainOb.Client.CallContract(context.TODO(), ethereum.CallMsg{
	//		From: fromAddr,
	//		To:   &toAddr,
	//		Data: input,
	//	}, big.NewInt(0).SetUint64(bn))
	//	if err != nil {
	//		log.Err(err).Msgf("%s CallContract error", chainOb.chain)
	//		return err
	//	}
	//	output, err := chainOb.abi.Unpack("getLockedAmount", res)
	//	if err != nil {
	//		log.Err(err).Msgf("%s Unpack error", chainOb.chain)
	//		return err
	//	}
	//	lockedAmount := *abi.ConvertType(output[0], new(*big.Int)).(**big.Int)
	//	//fmt.Printf("ETH: block %d: lockedAmount %d\n", bn, lockedAmount)
	//	supply = lockedAmount.String()
	//
	//} else if chainOb.chain == common.BSCChain {
	//	input, err := chainOb.abi.Pack("totalSupply")
	//	if err != nil {
	//		return fmt.Errorf("fail to totalSupply")
	//	}
	//	bn, err := chainOb.Client.BlockNumber(context.TODO())
	//	if err != nil {
	//		log.Err(err).Msgf("%s BlockNumber error", chainOb.chain)
	//		return err
	//	}
	//	fromAddr := ethcommon.HexToAddress(config.TSS_TEST_ADDRESS)
	//	toAddr := ethcommon.HexToAddress(config.BSC_TOKEN_ADDRESS)
	//	res, err := chainOb.Client.CallContract(context.TODO(), ethereum.CallMsg{
	//		From: fromAddr,
	//		To:   &toAddr,
	//		Data: input,
	//	}, big.NewInt(0).SetUint64(bn))
	//	if err != nil {
	//		log.Err(err).Msgf("%s CallContract error", chainOb.chain)
	//		return err
	//	}
	//	output, err := chainOb.abi.Unpack("totalSupply", res)
	//	if err != nil {
	//		log.Err(err).Msgf("%s Unpack error", chainOb.chain)
	//		return err
	//	}
	//	totalSupply := *abi.ConvertType(output[0], new(*big.Int)).(**big.Int)
	//	//fmt.Printf("BSC: block %d: totalSupply %d\n", bn, totalSupply)
	//	supply = totalSupply.String()
	//} else if chainOb.chain == common.POLYGONChain {
	//	input, err := chainOb.abi.Pack("totalSupply")
	//	if err != nil {
	//		return fmt.Errorf("fail to totalSupply")
	//	}
	//	bn, err := chainOb.Client.BlockNumber(context.TODO())
	//	if err != nil {
	//		log.Err(err).Msgf("%s BlockNumber error", chainOb.chain)
	//		return err
	//	}
	//	fromAddr := ethcommon.HexToAddress(config.TSS_TEST_ADDRESS)
	//	toAddr := ethcommon.HexToAddress(config.POLYGON_TOKEN_ADDRESS)
	//	res, err := chainOb.Client.CallContract(context.TODO(), ethereum.CallMsg{
	//		From: fromAddr,
	//		To:   &toAddr,
	//		Data: input,
	//	}, big.NewInt(0).SetUint64(bn))
	//	if err != nil {
	//		log.Err(err).Msgf("%s CallContract error", chainOb.chain)
	//		return err
	//	}
	//	output, err := chainOb.abi.Unpack("totalSupply", res)
	//	if err != nil {
	//		log.Err(err).Msgf("%s Unpack error", chainOb.chain)
	//		return err
	//	}
	//	totalSupply := *abi.ConvertType(output[0], new(*big.Int)).(**big.Int)
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

	bal, err := chainOb.Client.BalanceAt(context.TODO(), chainOb.Tss.Address(), nil)
	_, err = chainOb.bridge.PostGasBalance(chainOb.chain, bal.String(), blockNum)
	if err != nil {
		log.Err(err).Msg("PostGasBalance:")
		return err
	}
	return nil
}

func (chainOb *ChainObserver) observeFailedTx() {
	chainOb.mu.Lock()
	//for txhash, sendHash := range chainOb.txWatchList {
	//	receipt, err := chainOb.Client.TransactionReceipt(context.TODO(), txhash)
	//	if err != nil {
	//		continue
	//	}
	//	if receipt.Status == 0 { // failed tx
	//		log.Debug().Msgf("failed tx receipts: txhash %s sendHash %s", txhash.Hex(), sendHash)
	//		_, err = chainOb.bridge.PostReceiveConfirmation(sendHash, txhash.Hex(), receipt.BlockNumber.Uint64(), "", common.ReceiveStatus_Failed, chainOb.chain.String())
	//		if err != nil {
	//			log.Err(err).Msg("failed tx: PostReceiveConfirmation error ")
	//		}
	//	} else {
	//
	//	}
	//	delete(chainOb.txWatchList, txhash)
	//}
	chainOb.mu.Unlock()
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

	topics[0] = []ethcommon.Hash{logMPISendSignatureHash}

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

	// Read in ABI
	contractAbi := chainOb.abi

	// Pull out arguments from logs
	for _, vLog := range logs {
		log.Info().Msgf("TxBlockNumber %d Transaction Hash: %s topic %s\n", vLog.BlockNumber, vLog.TxHash.Hex()[:6], vLog.Topics[0].Hex()[:6])

		switch vLog.Topics[0].Hex() {
		case logMPISendSignatureHash.Hex():
			vals, err := contractAbi.Unpack("ZetaMessageSendEvent", vLog.Data)
			if err != nil {
				log.Err(err).Msg("error unpacking ZetaMessageSendEvent")
				continue
			}
			sender := vLog.Topics[1]
			destChainID := vals[0].(uint16)
			destContract := vals[1].([]byte)
			zetaAmount := vals[2].(*big.Int)
			gasLimit := vals[3].(*big.Int)
			message := vals[4].([]byte)
			zetaParams := vals[5].([]byte)

			_ = gasLimit
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
	chainOb.db.Put([]byte(PosKey), buf[:n], nil)
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

func (chainOb *ChainObserver) IsSendOutTxProcessed(sendHash string) (bool, error) {
	recvTopics := make([][]ethcommon.Hash, 4)
	recvTopics[3] = []ethcommon.Hash{ethcommon.HexToHash(sendHash)}
	recvTopics[0] = []ethcommon.Hash{logMPIReceiveSignatureHash}
	query := ethereum.FilterQuery{
		Addresses: []ethcommon.Address{ethcommon.HexToAddress(config.Chains[chainOb.chain.String()].MPIContractAddress)},
		FromBlock: big.NewInt(0), // LastBlock has been processed;
		ToBlock:   nil,
		Topics:    recvTopics,
	}
	logs, err := chainOb.Client.FilterLogs(context.Background(), query)
	if err != nil {
		return false, fmt.Errorf("IsSendOutTxProcessed: Client FilterLog fail %w", err)
	}
	if len(logs) == 0 {
		return false, nil
	}
	if len(logs) > 1 {
		log.Fatal().Msgf("More than two logs with send hash %s", sendHash)
		log.Fatal().Msgf("First one: %+v\nSecond one:%+v\n", logs[0], logs[1])
		return true, fmt.Errorf("More than two logs with send hash %s", sendHash)
	}
	for _, vLog := range logs {
		switch vLog.Topics[0].Hex() {
		case logMPIReceiveSignatureHash.Hex():
			fmt.Printf("Found sendHash %s on chain %s\n", sendHash, chainOb.chain)
			retval, err := chainOb.abi.Unpack("ZetaMessageReceiveEvent", vLog.Data)
			if err != nil {
				fmt.Println("error unpacking Unlock")
				continue
			}
			fmt.Printf("Topic 0: %s\n", vLog.Topics[0])
			fmt.Printf("Topic 1: %s\n", vLog.Topics[1])
			fmt.Printf("Topic 2: %s\n", vLog.Topics[2])
			//fmt.Printf("data: %d\n", retval[0].(*big.Int))
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
			} else {
				fmt.Printf("Included in block but not yet confirmed! included in block %d, current block %d\n", vLog.BlockNumber, chainOb.LastBlock)
				return false, nil
			}
			return true, nil
		}
	}

	return false, fmt.Errorf("IsSendOutTxProcessed: unknown chain %s", chainOb.chain)
}
