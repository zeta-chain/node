package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"net/http"
	"os"
	"strconv"
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
		HttpClient: &http.Client{},
	}
	resp := client.EthGetBlockByNumber(uint64(bn), false)
	var jsonObject map[string]interface{}
	if resp.Error != nil {
		fmt.Printf("Error: %s (code %d)\n", resp.Error.Message, resp.Error.Code)
		panic(resp.Error.Message)
	} else {
		//fmt.Printf("Result: %s\n", string(resp.Result))
		err := json.Unmarshal(resp.Result, &jsonObject)
		prettyJSON, err := json.MarshalIndent(jsonObject, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("Result: %s\n", string(prettyJSON))
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
	} else {
		jsonObject = make(map[string]interface{})
		err := json.Unmarshal(tx.Result, &jsonObject)
		prettyJSON, err := json.MarshalIndent(jsonObject, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Printf("Result: %s\n", string(prettyJSON))
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
	fmt.Printf("Receipt: %+v\n", receipt)
}

type EthClient struct {
	Endpoint   string
	HttpClient *http.Client
}

func (c *EthClient) EthGetBlockByNumber(blockNum uint64, verbose bool) *Response {
	client := c.HttpClient
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
	// Decode the response from JSON
	var rpcResp Response
	err = json.NewDecoder(resp.Body).Decode(&rpcResp)
	if err != nil {
		panic(err)
	}

	return &rpcResp
}

func (c *EthClient) EthGetTransactionReceipt(txhash string) *Response {
	client := c.HttpClient
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
	// Decode the response from JSON
	var rpcResp Response
	err = json.NewDecoder(resp.Body).Decode(&rpcResp)
	if err != nil {
		panic(err)
	}

	return &rpcResp
}
