package e2etests

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ZEVMRPC provides thin wrappers over eth JSON-RPC using CallContext, plus helpers to parse common responses.
type ZEVMRPC struct {
	client *ethclient.Client
}

func NewZEVMRPC(client *ethclient.Client) *ZEVMRPC {
	return &ZEVMRPC{client: client}
}

// Call wraps Client().CallContext for convenience.
func (z *ZEVMRPC) Call(ctx context.Context, result any, method string, params ...any) error {
	return z.client.Client().CallContext(ctx, result, method, params...)
}

func (z *ZEVMRPC) EthGetTransactionByHash(ctx context.Context, txHash common.Hash) (*TxByHash, error) {
	var raw json.RawMessage
	if err := z.Call(ctx, &raw, "eth_getTransactionByHash", txHash); err != nil {
		return nil, err
	}

	m, err := raw.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var txByHash TxByHash
	err = json.Unmarshal(m, &txByHash)
	if err != nil {
		return nil, err
	}

	return &txByHash, nil
}

func (z *ZEVMRPC) EthGetTransactionReceipt(ctx context.Context, txHash common.Hash) (*TxReceipt, error) {
	var raw json.RawMessage
	if err := z.Call(ctx, &raw, "eth_getTransactionReceipt", txHash); err != nil {
		return nil, err
	}

	m, err := raw.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var txReceipt TxReceipt
	err = json.Unmarshal(m, &txReceipt)
	if err != nil {
		return nil, err
	}

	return &txReceipt, nil
}

func (z *ZEVMRPC) EthGetBlockByNumber(ctx context.Context, blockNumber *big.Int, fullTx bool) (*Block, error) {
	var raw json.RawMessage
	if err := z.Call(ctx, &raw, "eth_getBlockByNumber", hexutil.EncodeBig(blockNumber), fullTx); err != nil {
		return nil, err
	}

	m, err := raw.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var block Block
	err = json.Unmarshal(m, &block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

func (z *ZEVMRPC) EthGetBlockByHash(ctx context.Context, blockHash common.Hash, fullTx bool) (*Block, error) {
	var raw json.RawMessage
	if err := z.Call(ctx, &raw, "eth_getBlockByHash", blockHash, fullTx); err != nil {
		return nil, err
	}

	m, err := raw.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var block Block
	err = json.Unmarshal(m, &block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

func (z *ZEVMRPC) DebugTraceTransaction(ctx context.Context, txHash common.Hash) (*TraceTx, error) {
	var raw json.RawMessage
	if err := z.Call(ctx, &raw, "debug_traceTransaction", txHash, map[string]any{"tracer": "callTracer"}); err != nil {
		return nil, err
	}

	m, err := raw.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var traceTx TraceTx
	err = json.Unmarshal(m, &traceTx)
	if err != nil {
		return nil, err
	}

	return &traceTx, nil
}

func (z *ZEVMRPC) DebugTraceBlockByNumber(ctx context.Context, blockNumber *big.Int) (TraceBlock, error) {
	var raw json.RawMessage
	if err := z.Call(ctx, &raw, "debug_traceBlockByNumber", hexutil.EncodeBig(blockNumber), map[string]any{"tracer": "callTracer"}); err != nil {
		return nil, err
	}

	m, err := raw.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var traceBlock TraceBlock
	err = json.Unmarshal(m, &traceBlock)
	if err != nil {
		return nil, err
	}
	return traceBlock, nil
}

type TraceTx struct {
	From    string    `json:"from"`
	Gas     string    `json:"gas"`
	GasUsed string    `json:"gasUsed"`
	Input   string    `json:"input"`
	Output  string    `json:"output,omitempty"`
	To      string    `json:"to"`
	Type    string    `json:"type"`
	Value   string    `json:"value,omitempty"`
	Calls   []TraceTx `json:"calls,omitempty"` // recursive children
}

type TraceBlock []struct {
	Result *TraceTx `json:"result,omitempty"`
}

type TxByHash struct {
	BlockHash            string `json:"blockHash"`
	BlockNumber          string `json:"blockNumber"`
	From                 string `json:"from"`
	Gas                  string `json:"gas"`
	GasPrice             string `json:"gasPrice"`
	MaxFeePerGas         string `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas"`
	Hash                 string `json:"hash"`
	Input                string `json:"input"`
	Nonce                string `json:"nonce"`
	To                   string `json:"to"`
	TransactionIndex     string `json:"transactionIndex"`
	Value                string `json:"value"`
	Type                 string `json:"type"`
	AccessList           []any  `json:"accessList"`
	ChainID              string `json:"chainId"`
	V                    string `json:"v"`
	R                    string `json:"r"`
	S                    string `json:"s"`
}

type TxReceipt struct {
	BlockHash         string `json:"blockHash"`
	BlockNumber       string `json:"blockNumber"`
	ContractAddress   any    `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	EffectiveGasPrice string `json:"effectiveGasPrice"`
	From              string `json:"from"`
	GasUsed           string `json:"gasUsed"`
	Logs              []struct {
		Address          string   `json:"address"`
		Topics           []string `json:"topics"`
		Data             string   `json:"data"`
		BlockNumber      string   `json:"blockNumber"`
		TransactionHash  string   `json:"transactionHash"`
		TransactionIndex string   `json:"transactionIndex"`
		BlockHash        string   `json:"blockHash"`
		LogIndex         string   `json:"logIndex"`
		Removed          bool     `json:"removed"`
	} `json:"logs"`
	LogsBloom        string `json:"logsBloom"`
	Status           string `json:"status"`
	To               string `json:"to"`
	TransactionHash  string `json:"transactionHash"`
	TransactionIndex string `json:"transactionIndex"`
	Type             string `json:"type"`
}

type Block struct {
	BaseFeePerGas    string   `json:"baseFeePerGas"`
	Difficulty       string   `json:"difficulty"`
	ExtraData        string   `json:"extraData"`
	GasLimit         string   `json:"gasLimit"`
	GasUsed          string   `json:"gasUsed"`
	Hash             string   `json:"hash"`
	LogsBloom        string   `json:"logsBloom"`
	Miner            string   `json:"miner"`
	MixHash          string   `json:"mixHash"`
	Nonce            string   `json:"nonce"`
	Number           string   `json:"number"`
	ParentHash       string   `json:"parentHash"`
	ReceiptsRoot     string   `json:"receiptsRoot"`
	Sha3Uncles       string   `json:"sha3Uncles"`
	Size             string   `json:"size"`
	StateRoot        string   `json:"stateRoot"`
	Timestamp        string   `json:"timestamp"`
	TotalDifficulty  string   `json:"totalDifficulty"`
	Transactions     []string `json:"transactions"`
	TransactionsRoot string   `json:"transactionsRoot"`
	Uncles           []any    `json:"uncles"`
}
