package metaclient

import (
	"context"
	"fmt"
	"github.com/Meta-Protocol/metacore/metaclient/config"
	"github.com/rs/zerolog/log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Chain configuration struct
// Filled with above constants depending on chain
type ChainObserver struct {
	chain     string
	router    string
	endpoint  string
	ticker    *time.Ticker
	abiString string
	abi 	  *abi.ABI
	client    *ethclient.Client
	bridge    *MetachainBridge
	lastBlock *big.Int
}

// Return configuration based on supplied target chain
func  NewChainObserver(chain string, bridge *MetachainBridge) (*ChainObserver, error){
	chainOb := ChainObserver{}
	chainOb.bridge = bridge

	// Initialize constants
	switch chain {
	case "Polygon":
		chainOb.chain = chain
		chainOb.router = config.POLY_ROUTER
		chainOb.endpoint = config.POLY_ENDPOINT
		chainOb.ticker = time.NewTicker(time.Duration(config.POLY_BLOCK_TIME) * time.Second)
		chainOb.abiString = config.META_ABI
	case "Ethereum":
		chainOb.chain = chain
		chainOb.router = config.ETH_ROUTER
		chainOb.endpoint = config.ETH_ENDPOINT
		chainOb.ticker = time.NewTicker(time.Duration(config.ETH_BLOCK_TIME) * time.Second)
		chainOb.abiString = config.META_LOCK_ABI
	case "BSC":
		chainOb.chain = chain
		chainOb.router = config.BSC_ROUTER
		chainOb.endpoint = config.BSC_ENDPOINT
		chainOb.ticker = time.NewTicker(time.Duration(config.BSC_BLOCK_TIME) * time.Second)
		chainOb.abiString = config.BSC_META_ABI
	}
	abi, err := abi.JSON(strings.NewReader(chainOb.abiString))
	if err != nil {
		return nil, err
	}
	chainOb.abi = &abi

	return &chainOb, nil
}

func (chainOb *ChainObserver) WatchRouter() {
	// Dial the router
	client, err := ethclient.Dial(chainOb.endpoint)
	if err != nil {
		log.Err(err).Msg("eth client Dial")
		return
	}

	// Set observer client
	chainOb.client = client

	// Set the latest block to begin scan
	chainOb.setLastBlock()

	// At each tick, query the router
	for range chainOb.ticker.C {
		err := chainOb.queryRouter()
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func (chainOb *ChainObserver) queryRouter() error {
	// Get most recent block number from client
	header, err := chainOb.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}

	// Generate query
	//fmt.Printf("chain %s start observing from %d\n", chainOb.chain, chainOb.lastBlock)
	query := ethereum.FilterQuery{
		Addresses: []ethcommon.Address{ethcommon.HexToAddress(chainOb.router)},
		FromBlock: chainOb.lastBlock.Add(chainOb.lastBlock, big.NewInt(1)),
		ToBlock:   header.Number,
	}
	log.Debug().Msgf("block from %d to %d", query.FromBlock, query.ToBlock)

	// Finally query the for the logs
	logs, err := chainOb.client.FilterLogs(context.Background(), query)
	if err != nil {
		return err
	}

	// Read in ABI
	contractAbi := chainOb.abi


	// LockSend event signature
	logLockSendSignature := []byte("LockSend(address,string,uint256,string,bytes)")
	logLockSendSignatureHash := crypto.Keccak256Hash(logLockSendSignature)

	// Unlock event signature
	logUnlockSignature := []byte("Unlock(address,uint256)")
	logUnlockSignatureHash := crypto.Keccak256Hash(logUnlockSignature)

	// BurnSend event signature
	logBurnSendSignature := []byte("BurnSend(address,address,uint256,uint256,string)")
	logBurnSendSignatureHash := crypto.Keccak256Hash(logBurnSendSignature)

	// MMinted event signature
	logMMintedSignature := []byte("MMinted(address,uint256)")
	logMMintedSignatureHash := crypto.Keccak256Hash(logMMintedSignature)

	// Update last block
	chainOb.lastBlock = header.Number

	// Pull out arguments from logs
	for _, vLog := range logs {
		log.Info().Msgf("TxBlockNumber %d Transaction Hash: %s\n", vLog.BlockNumber, vLog.TxHash.Hex())

		switch vLog.Topics[0].Hex() {
		case logLockSendSignatureHash.Hex():
			returnVal, err := contractAbi.Unpack("LockSend", vLog.Data)
			if err != nil {
				fmt.Println("error unpacking LockSend")
				continue
			}

			// PostSend to meta core
			metaHash, err := chainOb.bridge.PostSend(
				returnVal[0].(ethcommon.Address).String(),
				chainOb.chain,
				returnVal[1].(string),
				returnVal[3].(string),
				returnVal[2].(*big.Int).String(),
				"0",
				string(returnVal[4].([]uint8)), // TODO: figure out appropriate format for message
				vLog.TxHash.Hex(),
				vLog.BlockNumber,
			)
			if err != nil {
				fmt.Println("error posting to meta core")
				continue
			}

			fmt.Println("PostSend metahash: ", metaHash)
		case logBurnSendSignatureHash.Hex():
			returnVal, err := contractAbi.Unpack("BurnSend", vLog.Data)
			if err != nil {
				fmt.Println("error unpacking LockSend")
				continue
			}

			// PostSend to meta core
			metaHash, err := chainOb.bridge.PostSend(
				returnVal[0].(ethcommon.Address).String(),
				chainOb.chain,
				returnVal[1].(ethcommon.Address).String(),
				returnVal[3].(*big.Int).String(),
				returnVal[2].(*big.Int).String(),
				"0",
				returnVal[4].(string), // TODO: figure out appropriate format for message
				vLog.TxHash.Hex(),
				vLog.BlockNumber,
			)
			if err != nil {
				fmt.Println("error posting to meta core")
				continue
			}

			fmt.Println("PostSend metahash: ", metaHash)
		case logUnlockSignatureHash.Hex():
			returnVal, err := contractAbi.Unpack("Unlock", vLog.Data)
			if err != nil {
				fmt.Println("error unpacking LockSend")
				continue
			}

			// Post confirmation to meta core
			var sendHash, outTxHash string

			// sendHash = empty string for now
			// outTxHash = tx hash returned by signer.MMint
			var rxAddress string = returnVal[0].(ethcommon.Address).String()
			var mMint string = returnVal[1].(*big.Int).String()
			metaHash, err := chainOb.bridge.PostReceiveConfirmation(
				sendHash,
				outTxHash,
				vLog.BlockNumber,
				mMint,
			)
			if err != nil {
				fmt.Println("error posting confirmation to meta score")
				continue
			}

			fmt.Println("Receiver Address: ", rxAddress)
			fmt.Println("Post confirmation meta hash: ", metaHash)
		case logMMintedSignatureHash.Hex():
			// TODO: Handle MMinted
			fmt.Println("Observed MMinted")
		}
	}
	chainOb.lastBlock = header.Number

	return nil
}

func (chainOb *ChainObserver) setLastBlock() {

}
