package metaclient

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/Meta-Protocol/metacore/common"
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
	chain       common.Chain
	router      string
	endpoint    string
	ticker      *time.Ticker
	abiString   string
	abi         *abi.ABI // token contract ABI on non-ethereum chain; zetalocker on ethereum
	zetaAbi     *abi.ABI // only useful for ethereum; the token contract
	client      *ethclient.Client
	bridge      *MetachainBridge
	tss         TSSSigner
	lastBlock   uint64
	confCount   uint64 // must wait this many blocks to be considered "confirmed"
	txWatchList map[ethcommon.Hash]string
}

// Return configuration based on supplied target chain
func NewChainObserver(chain common.Chain, bridge *MetachainBridge, tss TSSSigner) (*ChainObserver, error) {
	chainOb := ChainObserver{}
	chainOb.bridge = bridge
	chainOb.txWatchList = make(map[ethcommon.Hash]string)
	chainOb.tss = tss
	// Initialize constants
	switch chain {
	case common.POLYGONChain:
		chainOb.chain = chain
		chainOb.router = config.POLYGON_TOKEN_ADDRESS
		chainOb.endpoint = config.POLY_ENDPOINT
		chainOb.ticker = time.NewTicker(time.Duration(config.POLY_BLOCK_TIME) * time.Second)
		chainOb.abiString = config.NONETH_ZETA_ABI
		chainOb.confCount = config.POLYGON_CONFIRMATION_COUNT

	case common.ETHChain:
		chainOb.chain = chain
		chainOb.router = config.ETH_ZETALOCK_ADDRESS
		chainOb.endpoint = config.ETH_ENDPOINT
		chainOb.ticker = time.NewTicker(time.Duration(config.ETH_BLOCK_TIME) * time.Second)
		chainOb.abiString = config.ETH_ZETALOCK_ABI
		chainOb.confCount = config.ETH_CONFIRMATION_COUNT

	case common.BSCChain:
		chainOb.chain = chain
		chainOb.router = config.BSC_TOKEN_ADDRESS
		chainOb.endpoint = config.BSC_ENDPOINT
		chainOb.ticker = time.NewTicker(time.Duration(config.BSC_BLOCK_TIME) * time.Second)
		chainOb.abiString = config.NONETH_ZETA_ABI
		chainOb.confCount = config.BSC_CONFIRMATION_COUNT
	}
	contractABI, err := abi.JSON(strings.NewReader(chainOb.abiString))
	if err != nil {
		return nil, err
	}
	chainOb.abi = &contractABI
	if chain == common.ETHChain {
		tokenABI, err := abi.JSON(strings.NewReader(config.ETH_ZETA_ABI))
		if err != nil {
			return nil, err
		}
		chainOb.zetaAbi = &tokenABI
	}

	// Dial the router
	client, err := ethclient.Dial(chainOb.endpoint)
	if err != nil {
		log.Err(err).Msg("eth client Dial")
		return nil, err
	}
	chainOb.client = client
	chainOb.lastBlock = chainOb.setLastBlock()
	// if ZetaCore does not have last heard block height, then use current
	if chainOb.lastBlock == 0 {
		header, err := chainOb.client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			return nil, err
		}
		chainOb.lastBlock = header.Number.Uint64()
	}

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

	return &chainOb, nil
}

func (chainOb *ChainObserver) WatchRouter() {
	// At each tick, query the router
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
	chainOb.txWatchList[hash] = sendhash
}

func (chainOb *ChainObserver) PostGasPrice() error {
	// GAS PRICE
	gasPrice, err := chainOb.client.SuggestGasPrice(context.TODO())
	if err != nil {
		log.Err(err).Msg("PostGasPrice:")
		return err
	}
	blockNum, err := chainOb.client.BlockNumber(context.TODO())
	if err != nil {
		log.Err(err).Msg("PostGasPrice:")
		return err
	}

	// SUPPLY
	var supply string // lockedAmount on ETH, totalSupply on other chains
	if chainOb.chain == common.ETHChain {
		input, err := chainOb.abi.Pack("getLockedAmount")
		if err != nil {
			return fmt.Errorf("fail to getLockedAmount")
		}
		bn, err := chainOb.client.BlockNumber(context.TODO())
		if err != nil {
			log.Err(err).Msgf("%s BlockNumber error", chainOb.chain)
			return err
		}
		fromAddr := ethcommon.HexToAddress(config.TSS_TEST_ADDRESS)
		toAddr := ethcommon.HexToAddress(config.ETH_ZETALOCK_ADDRESS)
		res, err := chainOb.client.CallContract(context.TODO(), ethereum.CallMsg{
			From: fromAddr,
			To:   &toAddr,
			Data: input,
		}, big.NewInt(0).SetUint64(bn))
		if err != nil {
			log.Err(err).Msgf("%s CallContract error", chainOb.chain)
			return err
		}
		output, err := chainOb.abi.Unpack("getLockedAmount", res)
		if err != nil {
			log.Err(err).Msgf("%s Unpack error", chainOb.chain)
			return err
		}
		lockedAmount := *abi.ConvertType(output[0], new(*big.Int)).(**big.Int)
		//fmt.Printf("ETH: block %d: lockedAmount %d\n", bn, lockedAmount)
		supply = lockedAmount.String()

	} else if chainOb.chain == common.BSCChain {
		input, err := chainOb.abi.Pack("totalSupply")
		if err != nil {
			return fmt.Errorf("fail to totalSupply")
		}
		bn, err := chainOb.client.BlockNumber(context.TODO())
		if err != nil {
			log.Err(err).Msgf("%s BlockNumber error", chainOb.chain)
			return err
		}
		fromAddr := ethcommon.HexToAddress(config.TSS_TEST_ADDRESS)
		toAddr := ethcommon.HexToAddress(config.BSC_TOKEN_ADDRESS)
		res, err := chainOb.client.CallContract(context.TODO(), ethereum.CallMsg{
			From: fromAddr,
			To:   &toAddr,
			Data: input,
		}, big.NewInt(0).SetUint64(bn))
		if err != nil {
			log.Err(err).Msgf("%s CallContract error", chainOb.chain)
			return err
		}
		output, err := chainOb.abi.Unpack("totalSupply", res)
		if err != nil {
			log.Err(err).Msgf("%s Unpack error", chainOb.chain)
			return err
		}
		totalSupply := *abi.ConvertType(output[0], new(*big.Int)).(**big.Int)
		//fmt.Printf("BSC: block %d: totalSupply %d\n", bn, totalSupply)
		supply = totalSupply.String()
	} else if chainOb.chain == common.POLYGONChain {
		input, err := chainOb.abi.Pack("totalSupply")
		if err != nil {
			return fmt.Errorf("fail to totalSupply")
		}
		bn, err := chainOb.client.BlockNumber(context.TODO())
		if err != nil {
			log.Err(err).Msgf("%s BlockNumber error", chainOb.chain)
			return err
		}
		fromAddr := ethcommon.HexToAddress(config.TSS_TEST_ADDRESS)
		toAddr := ethcommon.HexToAddress(config.POLYGON_TOKEN_ADDRESS)
		res, err := chainOb.client.CallContract(context.TODO(), ethereum.CallMsg{
			From: fromAddr,
			To:   &toAddr,
			Data: input,
		}, big.NewInt(0).SetUint64(bn))
		if err != nil {
			log.Err(err).Msgf("%s CallContract error", chainOb.chain)
			return err
		}
		output, err := chainOb.abi.Unpack("totalSupply", res)
		if err != nil {
			log.Err(err).Msgf("%s Unpack error", chainOb.chain)
			return err
		}
		totalSupply := *abi.ConvertType(output[0], new(*big.Int)).(**big.Int)
		//fmt.Printf("BSC: block %d: totalSupply %d\n", bn, totalSupply)
		supply = totalSupply.String()
	} else {
		log.Error().Msgf("chain not supported %s", chainOb.chain)
		return fmt.Errorf("unsupported chain %s", chainOb.chain)
	}

	_, err = chainOb.bridge.PostGasPrice(chainOb.chain, gasPrice.Uint64(), supply, blockNum)
	if err != nil {
		log.Err(err).Msg("PostGasPrice:")
		return err
	}

	bal, err := chainOb.client.BalanceAt(context.TODO(), chainOb.tss.Address(), nil)
	_, err = chainOb.bridge.PostGasBalance(chainOb.chain, bal.String(), blockNum)
	if err != nil {
		log.Err(err).Msg("PostGasBalance:")
		return err
	}
	return nil
}

func (chainOb *ChainObserver) observeFailedTx() {
	for txhash, sendHash := range chainOb.txWatchList {
		receipt, err := chainOb.client.TransactionReceipt(context.TODO(), txhash)
		if err != nil {
			continue
		}
		if receipt.Status == 0 { // failed tx
			log.Debug().Msgf("failed tx receipts: txhash %s sendHash %s", txhash.Hex(), sendHash)
			_, err = chainOb.bridge.PostReceiveConfirmation(sendHash, txhash.Hex(), receipt.BlockNumber.Uint64(), "", common.ReceiveStatus_Failed, chainOb.chain.String())
			if err != nil {
				log.Err(err).Msg("failed tx: PostReceiveConfirmation error ")
			}
		} else {
			log.Debug().Msgf("success tx receipts: txhash %s sendHash %s", txhash.Hex(), sendHash)
		}
		delete(chainOb.txWatchList, txhash)
	}
}

func (chainOb *ChainObserver) observeChain() error {
	header, err := chainOb.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}
	// "confirmed" current block number
	confirmedBlockNum := header.Number.Uint64() - chainOb.confCount
	// skip if no new block is produced.
	if confirmedBlockNum <= chainOb.lastBlock {
		return nil
	}
	toBlock := chainOb.lastBlock + 10 // read at most 10 blocks in one go
	if toBlock >= confirmedBlockNum {
		toBlock = confirmedBlockNum
	}
	query := ethereum.FilterQuery{
		Addresses: []ethcommon.Address{ethcommon.HexToAddress(chainOb.router)},
		FromBlock: big.NewInt(0).SetUint64(chainOb.lastBlock + 1), // lastBlock has been processed;
		ToBlock:   big.NewInt(0).SetUint64(toBlock),
	}
	//log.Debug().Msgf("signer %s block from %d to %d", chainOb.bridge.GetKeys().signerName, query.FromBlock, query.ToBlock)

	// Finally query the for the logs
	logs, err := chainOb.client.FilterLogs(context.Background(), query)
	if err != nil {
		return err
	}

	// Read in ABI
	contractAbi := chainOb.abi

	// LockSend event signature
	logLockSendSignature := []byte("LockSend(address,string,uint256,uint256,string,bytes)")
	logLockSendSignatureHash := crypto.Keccak256Hash(logLockSendSignature)

	// Unlock event signature
	logUnlockSignature := []byte("Unlock(address,uint256,bytes32)")
	logUnlockSignatureHash := crypto.Keccak256Hash(logUnlockSignature)

	// BurnSend event signature
	logBurnSendSignature := []byte("BurnSend(address,string,uint256,uint256,string,bytes)")
	logBurnSendSignatureHash := crypto.Keccak256Hash(logBurnSendSignature)

	// MMinted event signature
	logMMintedSignature := []byte("MMinted(address,uint256,bytes32)")
	logMMintedSignatureHash := crypto.Keccak256Hash(logMMintedSignature)

	// Pull out arguments from logs
	for _, vLog := range logs {
		log.Debug().Msgf("TxBlockNumber %d Transaction Hash: %s topic %s\n", vLog.BlockNumber, vLog.TxHash.Hex()[:6], vLog.Topics[0].Hex()[:6])

		switch vLog.Topics[0].Hex() {
		case logLockSendSignatureHash.Hex():
			returnVal, err := contractAbi.Unpack("LockSend", vLog.Data)
			if err != nil {
				log.Err(err).Msg("error unpacking LockSend")
				continue
			}

			// PostSend to meta core
			metaHash, err := chainOb.bridge.PostSend(
				returnVal[0].(ethcommon.Address).String(),
				chainOb.chain.String(),
				returnVal[1].(string),
				returnVal[4].(string),
				returnVal[2].(*big.Int).String(),
				returnVal[3].(*big.Int).String(),
				string(returnVal[5].([]byte)), // TODO: figure out appropriate format for message
				vLog.TxHash.Hex(),
				vLog.BlockNumber,
			)
			if err != nil {
				log.Err(err).Msg("error posting to meta core")
				continue
			}
			log.Debug().Msgf("LockSend detected: PostSend metahash: %s", metaHash)
		case logBurnSendSignatureHash.Hex():
			returnVal, err := contractAbi.Unpack("BurnSend", vLog.Data)
			if err != nil {
				log.Err(err).Msg("error unpacking LockSend")
				continue
			}

			// PostSend to meta core
			metaHash, err := chainOb.bridge.PostSend(
				returnVal[0].(ethcommon.Address).String(),
				chainOb.chain.String(),
				returnVal[1].(string),
				returnVal[3].(string),
				returnVal[2].(*big.Int).String(),
				returnVal[3].(*big.Int).String(),
				string(returnVal[4].([]byte)), // TODO: figure out appropriate format for message
				vLog.TxHash.Hex(),
				vLog.BlockNumber,
			)
			if err != nil {
				log.Err(err).Msg("error posting to meta core")
				continue
			}

			log.Debug().Msgf("BurnSend detected: PostSend metahash: %s", metaHash)
		case logUnlockSignatureHash.Hex():
			returnVal, err := contractAbi.Unpack("Unlock", vLog.Data)
			if err != nil {
				log.Err(err).Msg("error unpacking Unlock")
				continue
			}

			sendhash := returnVal[2].([32]byte)
			sendHash := "0x" + hex.EncodeToString(sendhash[:])
			var rxAddress string = returnVal[0].(ethcommon.Address).String()
			var mMint string = returnVal[1].(*big.Int).String()
			metaHash, err := chainOb.bridge.PostReceiveConfirmation(
				sendHash,
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
			log.Debug().Msgf("Unlock detected; recv %s Post confirmation meta hash %s", rxAddress, metaHash[:6])
			log.Debug().Msgf("Unlocked(sendhash=%s, outTxHash=%s, blockHeight=%d, amount=%s", sendHash[:6], vLog.TxHash.Hex()[:6], vLog.BlockNumber, mMint)

		case logMMintedSignatureHash.Hex():
			returnVal, err := contractAbi.Unpack("MMinted", vLog.Data)
			if err != nil {
				log.Err(err).Msg("error unpacking Unlock")
				continue
			}

			// outTxHash = tx hash returned by signer.MMint
			rxAddress := returnVal[0].(ethcommon.Address).String()
			mMint := returnVal[1].(*big.Int).String()
			sendhash := returnVal[2].([32]byte)
			sendHash := "0x" + hex.EncodeToString(sendhash[:])
			metaHash, err := chainOb.bridge.PostReceiveConfirmation(
				sendHash,
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
			log.Debug().Msgf("MMinted event detected; recv %s Post confirmation meta hash %s", rxAddress, metaHash[:6])
			log.Debug().Msgf("MMinted(sendhash=%s, outTxHash=%s, blockHeight=%d, mMint=%s", sendHash[:6], vLog.TxHash.Hex()[:6], vLog.BlockNumber, mMint)
		}
	}

	chainOb.lastBlock = toBlock

	return nil
}

// query ZetaCore about the last block that it has heard from a specific chain.
// return 0 if not existent.
func (chainOb *ChainObserver) setLastBlock() uint64 {
	lastheight, err := chainOb.bridge.GetLastBlockHeightByChain(chainOb.chain)
	if err != nil {
		log.Warn().Err(err).Msgf("setLastBlock")
		return 0
	}
	return lastheight.LastSendHeight
}

// query the base gas price for the block number bn.
func (chainOb *ChainObserver) GetBaseGasPrice() *big.Int {
	gasPrice, err := chainOb.client.SuggestGasPrice(context.TODO())
	if err != nil {
		log.Err(err).Msg("GetBaseGasPrice")
		return nil
	}
	return gasPrice
}
