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
	Number           int
	Hash             string
	ParentHash       string
	Nonce            string
	Sha3Uncles       string
	LogsBloom        string
	TransactionsRoot string
	StateRoot        string
	Miner            string
	Difficulty       *big.Int
	TotalDifficulty  *big.Int
	ExtraData        string
	Size             int
	GasLimit         int
	GasUsed          int
	Timestamp        int
	Uncles           []string
	Transactions     []Transaction
}

// Transaction EVM transaction.
type Transaction struct {
	Hash             string
	Nonce            int
	BlockHash        string
	BlockNumber      *int
	TransactionIndex *int
	From             string
	To               string
	Value            *big.Int
	Gas              int
	GasPrice         *big.Int
	Input            string
}

type hexInt int
type hexBig big.Int

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

	return parseBlock(raw)
}

// TransactionByHash2 is alternative to geth TransactionByHash that supports NON-ETH chains.
// See BlockByNumber2.
func (c *Client) TransactionByHash2(ctx context.Context, hash string) (*Transaction, error) {
	raw, err := c.call(ctx, "eth_getTransactionByHash", hash)
	if err != nil {
		return nil, errors.Wrapf(err, "transaction %s", hash)
	}

	return parseTransaction(raw)
}

func (c *Client) call(ctx context.Context, method string, args ...any) (json.RawMessage, error) {
	var raw json.RawMessage

	err := c.Client.Client().CallContext(ctx, &raw, method, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "method %s failed", method)
	}

	return raw, nil
}

// parses *Block from raw ETH RPC response.
// Note this is NOT the same as json.Unmarshal([]byte, &Block{})
func parseBlock(raw json.RawMessage) (*Block, error) {
	var proxy struct {
		Number           hexInt            `json:"number"`
		Hash             string            `json:"hash"`
		ParentHash       string            `json:"parentHash"`
		Nonce            string            `json:"nonce"`
		Sha3Uncles       string            `json:"sha3Uncles"`
		LogsBloom        string            `json:"logsBloom"`
		TransactionsRoot string            `json:"transactionsRoot"`
		StateRoot        string            `json:"stateRoot"`
		Miner            string            `json:"miner"`
		Difficulty       *hexBig           `json:"difficulty"`
		TotalDifficulty  *hexBig           `json:"totalDifficulty"`
		ExtraData        string            `json:"extraData"`
		Size             hexInt            `json:"size"`
		GasLimit         hexInt            `json:"gasLimit"`
		GasUsed          hexInt            `json:"gasUsed"`
		Timestamp        hexInt            `json:"timestamp"`
		Uncles           []string          `json:"uncles"`
		Transactions     []json.RawMessage `json:"transactions"`
	}

	if err := json.Unmarshal(raw, &proxy); err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal proxy block")
	}

	txs := make([]Transaction, 0, len(proxy.Transactions))
	for i, txRaw := range proxy.Transactions {
		tx, err := parseTransaction(txRaw)
		if err != nil {
			return nil, errors.Wrapf(err, "tx %d", i)
		}

		txs = append(txs, *tx)
	}

	return &Block{
		Number:           int(proxy.Number),
		Hash:             proxy.Hash,
		ParentHash:       proxy.ParentHash,
		Nonce:            proxy.Nonce,
		Sha3Uncles:       proxy.Sha3Uncles,
		LogsBloom:        proxy.LogsBloom,
		TransactionsRoot: proxy.TransactionsRoot,
		StateRoot:        proxy.StateRoot,
		Miner:            proxy.Miner,
		Difficulty:       (*big.Int)(proxy.Difficulty),
		TotalDifficulty:  (*big.Int)(proxy.TotalDifficulty),
		ExtraData:        proxy.ExtraData,
		Size:             int(proxy.Size),
		GasLimit:         int(proxy.GasLimit),
		GasUsed:          int(proxy.GasUsed),
		Timestamp:        int(proxy.Timestamp),
		Uncles:           proxy.Uncles,
		Transactions:     txs,
	}, nil
}

// parses Transaction from raw ETH RPC response.
// Note this is NOT the same as json.Unmarshal([]byte, &Transaction{})
func parseTransaction(data []byte) (*Transaction, error) {
	var proxy struct {
		Hash             string  `json:"hash"`
		Nonce            hexInt  `json:"nonce"`
		BlockHash        string  `json:"blockHash"`
		BlockNumber      *hexInt `json:"blockNumber"`
		TransactionIndex *hexInt `json:"transactionIndex"`
		From             string  `json:"from"`
		To               string  `json:"to"`
		Value            *hexBig `json:"value"`
		Gas              hexInt  `json:"gas"`
		GasPrice         *hexBig `json:"gasPrice"`
		Input            string  `json:"input"`
	}

	if err := json.Unmarshal(data, &proxy); err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal proxy tx")
	}

	return &Transaction{
		Hash:             proxy.Hash,
		Nonce:            int(proxy.Nonce),
		BlockHash:        proxy.BlockHash,
		BlockNumber:      (*int)(proxy.BlockNumber),
		TransactionIndex: (*int)(proxy.TransactionIndex),
		From:             proxy.From,
		To:               proxy.To,
		Value:            (*big.Int)(proxy.Value),
		Gas:              int(proxy.Gas),
		GasPrice:         (*big.Int)(proxy.GasPrice),
		Input:            proxy.Input,
	}, nil
}

func (i *hexInt) UnmarshalJSON(data []byte) error {
	result, err := parseInt(string(bytes.Trim(data, `"`)))
	if err != nil {
		return errors.Wrapf(err, "failed to parse int from %s", data)
	}

	*i = hexInt(result)

	return nil
}

func (i *hexBig) UnmarshalJSON(data []byte) error {
	bi, err := parseBigInt(string(bytes.Trim(data, `"`)))
	if err != nil {
		return errors.Wrapf(err, "failed to parse big.Int from %s", data)
	}

	(*i) = (hexBig)(*bi)

	return nil
}

func parseBigInt(value string) (*big.Int, error) {
	i := big.NewInt(0)
	if _, ok := i.SetString(value, 0); !ok {
		return nil, fmt.Errorf("failed to parse big.Int from %s", value)
	}

	return i, nil
}

func parseInt(value string) (int, error) {
	value = strings.TrimPrefix(value, "0x")

	i, err := strconv.ParseInt(value, 16, 64)
	if err != nil {
		return 0, errors.Wrap(err, value)
	}

	return int(i), nil
}
