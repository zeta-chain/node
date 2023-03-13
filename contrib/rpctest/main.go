package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeta-chain/zetacore/contracts/evm/zetaeth"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	ZetaEthPriv = "9D00E4D7A8A14384E01CD90B83745BCA847A66AD8797A9904A200C28C2648E64"
)

type Request struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type Response struct {
	Result json.RawMessage `json:"result"`
	Error  *Error          `json:"error"`
	ID     int             `json:"id"`
}
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <blocknum>\n", os.Args[0])
		os.Exit(1)
	}
	fmt.Printf("Start testing the zEVM ETH JSON-RPC for all txs...\n")
	fmt.Printf("Test1: simple gas voter tx\n")

	bn, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	if false {
		// THIS WOULD NOT WORK: see https://github.com/zeta-chain/zeta-node/pull/445
		// USE RAW JSON-RPC INSTEAD
		zevmClient, err := ethclient.Dial("http://localhost:8545")
		if err != nil {
			panic(err)
		}

		block, err := zevmClient.BlockByNumber(context.Background(), big.NewInt(int64(bn)))
		if err != nil {
			panic(err)
		}

		fmt.Printf("Block number: %d, num of txs %d (should be 1)\n", block.Number(), len(block.Transactions()))
	}

	client := &EthClient{
		Endpoint:   "http://localhost:8545",
		HTTPClient: &http.Client{},
	}
	resp := client.EthGetBlockByNumber(uint64(bn), false)
	var jsonObject map[string]interface{}
	if resp.Error != nil {
		fmt.Printf("Error: %s (code %d)\n", resp.Error.Message, resp.Error.Code)
		panic(resp.Error.Message)
	} else {
		//fmt.Printf("Result: %s\n", string(resp.Result))
		err = json.Unmarshal(resp.Result, &jsonObject)
		if err != nil {
			panic(err)
		}
	}

	txs, ok := jsonObject["transactions"].([]interface{})
	if !ok || len(txs) != 1 {
		panic("Wrong number of txs")
	}
	txhash, ok := txs[0].(string)
	if !ok {
		panic("Wrong tx type")
	}
	fmt.Printf("Tx hash: %s\n", txhash)
	tx := client.EthGetTransactionReceipt(txhash)
	if tx.Error != nil {
		fmt.Printf("Error: %s (code %d)\n", tx.Error.Message, tx.Error.Code)
		panic(tx.Error.Message)
	}

	// tx receipt can be queried by ethclient queries.
	zevmClient, err := ethclient.Dial(client.Endpoint)
	if err != nil {
		panic(err)
	}
	receipt, err := zevmClient.TransactionReceipt(context.Background(), ethcommon.HexToHash(txhash))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Receipt status: %+v\n", receipt.Status)

	// HeaderByHash works; BlockByHash does not work;
	// main offending RPC is the transaction type; we have custom type id 56
	// which is not recognized by the go-ethereum client.
	blockHeader, err := zevmClient.HeaderByNumber(context.Background(), big.NewInt(int64(bn)))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Block header TxHash: %+v\n", blockHeader.TxHash)

	chainid, err := zevmClient.ChainID(context.Background())
	if err != nil {
		panic(err)
	}
	zetaEthPrivKey, err := crypto.HexToECDSA(ZetaEthPriv)
	if err != nil {
		panic(err)
	}
	zevmAuth, err := bind.NewKeyedTransactorWithChainID(zetaEthPrivKey, chainid)
	if err != nil {
		panic(err)
	}
	zetaContractAddress, tx2, zetaContract, err := zetaeth.DeployZetaEth(zevmAuth, zevmClient, big.NewInt(2_100_000_000))
	_, _ = zetaContractAddress, zetaContract
	if err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Second)
	receipt, err = zevmClient.TransactionReceipt(context.Background(), tx2.Hash())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deploy EthZeta Contract Receipt: %+v\n", receipt)
	receipt2 := client.EthGetTransactionReceipt(tx2.Hash().Hex())
	if receipt2.Error != nil {
		fmt.Printf("Error: %s (code %d)\n", receipt2.Error.Message, receipt2.Error.Code)
		panic(tx.Error.Message)
	} else {
		jsonObject = make(map[string]interface{})
		err = json.Unmarshal(receipt2.Result, &jsonObject)
		if err != nil {
			panic(err)
		}
		prettyJSON, err := json.MarshalIndent(jsonObject, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("Result: %s\n", string(prettyJSON))
	}
	fmt.Printf("ZetaEth Contract Address: %s\n", zetaContractAddress.Hex())
	if zetaContractAddress != receipt.ContractAddress {
		panic(fmt.Sprintf("Contract address mismatch: wanted %s, got %s", zetaContractAddress, receipt.ContractAddress))
	}

}

type EthClient struct {
	Endpoint   string
	HTTPClient *http.Client
}

func (c *EthClient) EthGetBlockByNumber(blockNum uint64, verbose bool) *Response {
	client := c.HTTPClient
	hexBlockNum := fmt.Sprintf("0x%x", blockNum)
	req := &Request{
		Jsonrpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params: []interface{}{
			hexBlockNum,
			verbose,
		},
		ID: 1,
	}

	// Encode the request to JSON
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(req)
	if err != nil {
		panic(err)
	}
	// Create a new HTTP request
	httpReq, err := http.NewRequest("POST", c.Endpoint, buf)
	if err != nil {
		panic(err)
	}
	// Set the content type header
	httpReq.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	resp, err := client.Do(httpReq)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close() //#nosec
	// Decode the response from JSON
	var rpcResp Response
	err = json.NewDecoder(resp.Body).Decode(&rpcResp)
	if err != nil {
		panic(err)
	}

	return &rpcResp
}

func (c *EthClient) EthGetTransactionReceipt(txhash string) *Response {
	client := c.HTTPClient
	req := &Request{
		Jsonrpc: "2.0",
		Method:  "eth_getTransactionReceipt",
		Params: []interface{}{
			txhash,
		},
		ID: 1,
	}

	// Encode the request to JSON
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(req)
	if err != nil {
		panic(err)
	}
	// Create a new HTTP request
	httpReq, err := http.NewRequest("POST", c.Endpoint, buf)
	if err != nil {
		panic(err)
	}
	// Set the content type header
	httpReq.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	resp, err := client.Do(httpReq)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close() //#nosec
	// Decode the response from JSON
	var rpcResp Response
	err = json.NewDecoder(resp.Body).Decode(&rpcResp)
	if err != nil {
		panic(err)
	}

	return &rpcResp
}
