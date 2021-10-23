package metaclientd

import (
	"context"
	"fmt"
	"log"
	"math/big"
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
	POLY_ROUTER = ""
	BSC_ROUTER  = ""

	// API Endpoints
	ETH_ENDPOINT  = "https://ropsten.infura.io/v3/90705d89baca4c2f9a8178f86d30c4f8"
	POLY_ENDPOINT = ""
	BSC_ENDPOINT  = ""

	// Ticker timers
	ETH_BLOCK_TIME  = 15
	POLY_BLOCK_TIME = 3
	BSC_BLOCK_TIME  = 3

	// ABIs
	META_LOCK_ABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"receiver\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"chainid\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"LockSend\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Unlock\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOracleAddres\",\"type\":\"address\"}],\"name\":\"changeOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"receiver\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"chainid\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"message\",\"type\":\"bytes\"}],\"name\":\"lockSend\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceTSSAddressUpdater\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"unlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"updateTSSAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"metaAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_TSSAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_TSSAddressUpdater\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_OracleUpdater\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"getLockedAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"OracleUpdater\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TSSAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TSSAddressUpdater\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
)

// Chain configuration struct
// Filled with above constants depending on chain
type MetaObserver struct {
	router    string
	endpoint  string
	ticker    *time.Ticker
	abi       string
	client    *ethclient.Client
	lastBlock *big.Int
}

// Return configuration based on supplied target chain
func (mo *MetaObserver) initObserver(chain string) {
	switch chain {
	case "Polygon":
		mo.router = POLY_ROUTER
		mo.endpoint = POLY_ENDPOINT
		mo.ticker = time.NewTicker(time.Duration(POLY_BLOCK_TIME) * time.Second)
		// TODO: ABI
	case "Ethereum":
		mo.router = ETH_ROUTER
		mo.endpoint = ETH_ENDPOINT
		mo.ticker = time.NewTicker(time.Duration(ETH_BLOCK_TIME) * time.Second)
		mo.abi = META_LOCK_ABI
	case "BSC":
		mo.router = BSC_ROUTER
		mo.endpoint = BSC_ENDPOINT
		mo.ticker = time.NewTicker(time.Duration(BSC_BLOCK_TIME) * time.Second)
		// TODO: ABI
	}
}

func (mo *MetaObserver) WatchRouter(chain string) {
	// Initialize variables
	mo.initObserver(chain)

	// Dial the router
	client, err := ethclient.Dial(mo.endpoint)
	if err != nil {
		log.Fatal(err)
	}

	// Set observer client
	mo.client = client

	// TODO: Assuming last block observed is zero
	// We need to integrate with meta blockchain to check last
	// observed block

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	mo.lastBlock = big.NewInt(0).Sub(header.Number, big.NewInt(int64(10)))

	for {
		select {
		case t := <-mo.ticker.C:
			fmt.Println("Ticker at ", t)

			// At each tick, query the router
			err := mo.queryRouter()
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}

}

func (mo *MetaObserver) queryRouter() error {
	//router_address := ethcommon.HexToAddress(mo.router)
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
	contractAbi, err := abi.JSON(strings.NewReader(string(mo.abi)))
	if err != nil {
		return err
	}

	// Look for LockSend event
	logLockSendSignature := []byte("LockSend(address,string,uint256,string,bytes)")
	logLockSendSignatureHash := crypto.Keccak256Hash(logLockSendSignature)

	// Update last block
	mo.lastBlock = header.Number

	// Pull out arguments from logs
	for _, vLog := range logs {
		switch vLog.Topics[0].Hex() {
		case logLockSendSignatureHash.Hex():
			returnVal, err := contractAbi.Unpack("LockSend", vLog.Data)
			if err != nil {
				fmt.Println("error unpacking LockSend")
				continue
			}

			// read and validate the transaction
			mo.readAndValidate(returnVal)
		}
	}

	return nil
}

func (mo *MetaObserver) readAndValidate(values []interface{}) {
	fmt.Println("Send Address: ", values[0])
	fmt.Println("Rx Address: ", values[1])
	fmt.Println("Amount: ", values[2])
	fmt.Println("ChainID: ", values[3])
	fmt.Println("Message: ", values[4])
}
