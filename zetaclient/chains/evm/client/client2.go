package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

// Block EVM block.
type Block struct {
	Number           Int64         `json:"number"`
	Hash             string        `json:"hash"`
	ParentHash       string        `json:"parentHash"`
	Nonce            string        `json:"nonce"`
	Sha3Uncles       string        `json:"sha3Uncles"`
	LogsBloom        string        `json:"logsBloom"`
	TransactionsRoot string        `json:"transactionsRoot"`
	StateRoot        string        `json:"stateRoot"`
	Miner            string        `json:"miner"`
	Difficulty       BigInt        `json:"difficulty"`
	TotalDifficulty  BigInt        `json:"totalDifficulty"`
	ExtraData        string        `json:"extraData"`
	Size             Int64         `json:"size"`
	GasLimit         Int64         `json:"gasLimit"`
	GasUsed          Int64         `json:"gasUsed"`
	Timestamp        Int64         `json:"timestamp"`
	Uncles           []string      `json:"uncles"`
	Transactions     []Transaction `json:"transactions"`
}

// Transaction EVM transaction.
type Transaction struct {
	Hash             string `json:"hash"`
	Nonce            Int64  `json:"nonce"`
	BlockHash        string `json:"blockHash"`
	BlockNumber      *Int64 `json:"blockNumber"`
	TransactionIndex *Int64 `json:"transactionIndex"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            BigInt `json:"value"`
	Gas              Int64  `json:"gas"`
	GasPrice         BigInt `json:"gasPrice"`
	Input            string `json:"input"`
}

type Int64 int64

func (i *Int64) UnmarshalJSON(data []byte) error {
	result, err := parseInt64(string(bytes.Trim(data, `"`)))
	if err != nil {
		return errors.Wrap(err, "unable to parse int64")
	}

	*i = Int64(result)

	return nil
}

type BigInt big.Int

func (i *BigInt) UnmarshalJSON(data []byte) error {
	result, err := parseBigInt(string(bytes.Trim(data, `"`)))
	if err != nil {
		return errors.Wrap(err, "unable to parse big.Int")
	}

	*i = BigInt(*result)

	return nil
}

// BlockByNumber2 is alternative to geth BlockByNumber that supports NON-ETH chains.
// For example, OP stack has different tx types that result in err="transaction type not supported"
//
// See https://github.com/zeta-chain/node/issues/3386
// See https://github.com/ethereum/go-ethereum/issues/29407
func (c *Client) BlockByNumber2(ctx context.Context, blockNumber *big.Int) (*Block, error) {
	raw, err := c.call(ctx, "eth_getBlockByNumber", hexutil.EncodeBig(blockNumber), true)
	if err != nil {
		return nil, errors.Wrapf(err, "block %d", blockNumber.Uint64())
	}

	return blockFromRaw(raw)
}

// TransactionByHash2 is alternative to geth TransactionByHash that supports NON-ETH chains.
// See BlockByNumber2.
func (c *Client) TransactionByHash2(ctx context.Context, hash string) (*Transaction, error) {
	raw, err := c.call(ctx, "eth_getTransactionByHash", hash)
	if err != nil {
		return nil, errors.Wrapf(err, "transaction %s", hash)
	}

	var tx Transaction
	if err := json.Unmarshal(raw, &tx); err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal transaction")
	}

	return &tx, nil
}

func (c *Client) call(ctx context.Context, method string, args ...any) (json.RawMessage, error) {
	var raw json.RawMessage

	err := c.Client.Client().CallContext(ctx, &raw, method, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "method %s failed", method)
	}

	return raw, nil
}

func blockFromRaw(raw json.RawMessage) (*Block, error) {
	var block Block
	if err := json.Unmarshal(raw, &block); err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal block")
	}

	return &block, nil
}

// parseInt parse hex string value to int
func parseInt64(raw string) (int64, error) {
	raw = strings.TrimPrefix(raw, "0x")

	i, err := strconv.ParseInt(raw, 16, 64)
	if err != nil {
		return 0, err
	}

	return i, nil
}

func parseBigInt(value string) (*big.Int, error) {
	num := big.NewInt(0)

	_, ok := num.SetString(strings.TrimPrefix(value, "0x"), 16)
	if !ok {
		return nil, fmt.Errorf("unable to parse big.Int from %s", value)
	}

	return num, nil
}
