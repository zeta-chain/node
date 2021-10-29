package metaclientd

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Constants
const (
	// Routers
	ETH_ROUTER  = "0xaEf30949004FAAfcAd2c0e47c0491D91Dd76f7AA"
	POLY_ROUTER = "0x853dFA5D44870aE376F7D2464dA47F17Ee4570d1"
	BSC_ROUTER  = "0x853dFA5D44870aE376F7D2464dA47F17Ee4570d1"

	// API Endpoints
	ETH_ENDPOINT  = "https://ropsten.infura.io/v3/90705d89baca4c2f9a8178f86d30c4f8"
	POLY_ENDPOINT = "https://speedy-nodes-nyc.moralis.io/43709093db18f2726a03c254/polygon/mainnet/archive"
	BSC_ENDPOINT  = "https://speedy-nodes-nyc.moralis.io/43709093db18f2726a03c254/bsc/testnet/archive"

	// Ticker timers
	ETH_BLOCK_TIME  = 15
	POLY_BLOCK_TIME = 3
	BSC_BLOCK_TIME  = 3

	// ABIs
	META_LOCK_ABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"receiver\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"chainid\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"LockSend\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Unlock\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOracleAddres\",\"type\":\"address\"}],\"name\":\"changeOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"receiver\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"chainid\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"lockSend\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceTSSAddressUpdater\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"unlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"updateTSSAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"metaAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_TSSAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_TSSAddressUpdater\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_OracleUpdater\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"getLockedAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"OracleUpdater\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TSSAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TSSAddressUpdater\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
	META_ABI      = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"chainid\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"BurnSend\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"burnee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MBurnt\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"mintee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"MMinted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"chainid\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"burn_send\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burnFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"mintee\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"burnee\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"permit_burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"mintee\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"permit_mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newMinter\",\"type\":\"address\"}],\"name\":\"setMinter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_minterAdmin\",\"type\":\"address\"}],\"name\":\"setMinterAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newOracle\",\"type\":\"address\"}],\"name\":\"setOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialSupply\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"BURN_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_TOTAL_SUPPLY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MINT_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minterAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oracleAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
)

// Chain configuration struct
// Filled with above constants depending on chain
type MetaObserver struct {
	chain     string
	router    string
	endpoint  string
	ticker    *time.Ticker
	abi       string
	client    *ethclient.Client
	bridge    *MetachainBridge
	lastBlock *big.Int
}

// Return configuration based on supplied target chain
func (mo *MetaObserver) InitMetaObserver(chain string) {
	// Initialize constants
	switch chain {
	case "Polygon":
		mo.chain = chain
		mo.router = POLY_ROUTER
		mo.endpoint = POLY_ENDPOINT
		mo.ticker = time.NewTicker(time.Duration(POLY_BLOCK_TIME) * time.Second)
		mo.abi = META_ABI
	case "Ethereum":
		mo.chain = chain
		mo.router = ETH_ROUTER
		mo.endpoint = ETH_ENDPOINT
		mo.ticker = time.NewTicker(time.Duration(ETH_BLOCK_TIME) * time.Second)
		mo.abi = META_LOCK_ABI
	case "BSC":
		mo.chain = chain
		mo.router = BSC_ROUTER
		mo.endpoint = BSC_ENDPOINT
		mo.ticker = time.NewTicker(time.Duration(BSC_BLOCK_TIME) * time.Second)
		mo.abi = META_ABI
	}

	// Initialize Bridge
	err := mo.createBridge()
	if err != nil {
		log.Fatal(err)
		return
	}
}

func (mo *MetaObserver) WatchRouter() {
	// Dial the router
	client, err := ethclient.Dial(mo.endpoint)
	if err != nil {
		log.Fatal(err)
	}

	// Set observer client
	mo.client = client

	// Set the latest block to begin scan
	mo.setLastBlock()

	// At each tick, query the router
	for range mo.ticker.C {
		err := mo.queryRouter()
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func (mo *MetaObserver) queryRouter() error {
	// Get most recent block number from client
	header, err := mo.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}

	// Generate query
	query := ethereum.FilterQuery{
		Addresses: []ethcommon.Address{ethcommon.HexToAddress(mo.router)},
		FromBlock: mo.lastBlock,
		ToBlock:   header.Number,
	}

	// Finally query the for the logs
	logs, err := mo.client.FilterLogs(context.Background(), query)
	if err != nil {
		return err
	}

	// Read in ABI
	contractAbi, err := abi.JSON(strings.NewReader(mo.abi))
	if err != nil {
		return err
	}

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
	mo.lastBlock = header.Number

	// Pull out arguments from logs
	for _, vLog := range logs {
		fmt.Println("Transaction Hash: ", vLog.TxHash.Hex())
		fmt.Println("TxBlockNumber: ", vLog.BlockNumber)

		switch vLog.Topics[0].Hex() {
		case logLockSendSignatureHash.Hex():
			returnVal, err := contractAbi.Unpack("LockSend", vLog.Data)
			if err != nil {
				fmt.Println("error unpacking LockSend")
				continue
			}

			// PostSend to meta core
			metaHash, err := mo.bridge.PostSend(
				returnVal[0].(ethcommon.Address).String(),
				mo.chain,
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
			metaHash, err := mo.bridge.PostSend(
				returnVal[0].(ethcommon.Address).String(),
				mo.chain,
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
			var rxAddress string = returnVal[0].(ethcommon.Address).String()
			var mMint string = returnVal[1].(*big.Int).String()
			metaHash, err := mo.bridge.PostReceiveConfirmation(
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

	return nil
}

func (mo *MetaObserver) setLastBlock() {
	// Check metacore for last block checked and set initial last block
	var useLastBlockHeight bool = true
	lastBlockHeights, err := mo.bridge.GetLastBlockHeight()

	// If there's an error retrieving blocks,
	// continue and use latest block - 10
	if err != nil {
		log.Fatal(err)
		useLastBlockHeight = false
	}

	// If metacore reports no previous blocks, use
	// latest block - 10
	if len(lastBlockHeights) <= 0 {
		useLastBlockHeight = false
	}

	if useLastBlockHeight {
		// Last block from meta core
		mo.lastBlock = big.NewInt(int64(lastBlockHeights[0].LastSendHeight))
	} else {
		// Most recent block - 10
		header, err := mo.client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			log.Fatal(err)
		}
		mo.lastBlock = big.NewInt(0).Sub(header.Number, big.NewInt(int64(10)))
	}
}

func (mo *MetaObserver) createBridge() error {
	// TODO: How do we properly set these values?
	signerName := "alice"
	signerPass := "password"

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	chainHomeFoler := filepath.Join(homeDir, ".metacore")

	kb, _, err := GetKeyringKeybase(chainHomeFoler, signerName, signerPass)
	if err != nil {
		return err
	}

	k := NewKeysWithKeybase(kb, signerName, signerPass)

	chainIP := "127.0.0.1"
	bridge, err := NewMetachainBridge(k, chainIP, "alice")
	if err != nil {
		return err
	}

	mo.bridge = bridge

	return nil
}
