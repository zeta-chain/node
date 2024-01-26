package ethrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// EthError - ethereum error
type EthError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err EthError) Error() string {
	return fmt.Sprintf("Error %d (%s)", err.Code, err.Message)
}

type ethResponse struct {
	ID      int             `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *EthError       `json:"error"`
}

type ethRequest struct {
	ID      int           `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// EthRPC - Ethereum rpc client
type EthRPC struct {
	url    string
	client *http.Client
}

// New create new rpc client with given url
func New(url string, options ...func(rpc *EthRPC)) *EthRPC {
	rpc := &EthRPC{
		url:    url,
		client: http.DefaultClient,
	}
	for _, option := range options {
		option(rpc)
	}

	return rpc
}

// NewEthRPC create new rpc client with given url
func NewEthRPC(url string, options ...func(rpc *EthRPC)) *EthRPC {
	return New(url, options...)
}

// URL returns client url
func (rpc *EthRPC) URL() string {
	return rpc.url
}

// Call returns raw response of method call
func (rpc *EthRPC) Call(method string, params ...interface{}) (json.RawMessage, error) {
	request := ethRequest{
		ID:      1,
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	response, err := rpc.client.Post(rpc.url, "application/json", bytes.NewBuffer(body))
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	resp := new(ethResponse)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, *resp.Error
	}

	return resp.Result, nil

}

// RawCall returns raw response of method call (Deprecated)
func (rpc *EthRPC) RawCall(method string, params ...interface{}) (json.RawMessage, error) {
	return rpc.Call(method, params...)
}

func (rpc *EthRPC) getBlock(method string, withTransactions bool, params ...interface{}) (*Block, error) {
	result, err := rpc.RawCall(method, params...)
	if err != nil {
		return nil, err
	}
	if bytes.Equal(result, []byte("null")) {
		return nil, nil
	}

	var response proxyBlock
	if withTransactions {
		response = new(proxyBlockWithTransactions)
	} else {
		response = new(proxyBlockWithoutTransactions)
	}

	err = json.Unmarshal(result, response)
	if err != nil {
		return nil, err
	}

	block := response.toBlock()
	return &block, nil
}

// EthGetBlockByNumber returns information about a block by block number.
func (rpc *EthRPC) EthGetBlockByNumber(number uint64, withTransactions bool) (*Block, error) {
	return rpc.getBlock("eth_getBlockByNumber", withTransactions, IntToHex(number), withTransactions)
}
