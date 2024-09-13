package types

import (
	"bytes"
	"encoding/json"
	"math/big"
	"strconv"
	"strings"
	"unsafe"

	"github.com/onrik/ethrpc"
)

// Transaction - transaction object
type Transaction struct {
	Hash             string
	Nonce            int
	BlockHash        string
	BlockNumber      *int
	TransactionIndex *int
	From             string
	To               string
	Value            big.Int
	Gas              int
	GasPrice         big.Int
	Input            string
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Transaction) UnmarshalJSON(data []byte) error {
	proxy := new(ProxyTransaction)
	if err := json.Unmarshal(data, proxy); err != nil {
		return err
	}

	*t = *(*Transaction)(unsafe.Pointer(proxy))

	return nil
}

// Block - block object
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
	Difficulty       big.Int
	TotalDifficulty  big.Int
	ExtraData        string
	Size             int
	GasLimit         int
	GasUsed          int
	Timestamp        int
	Uncles           []string
	Transactions     []Transaction
}

type ProxyTransaction struct {
	Hash             string  `json:"hash"`
	Nonce            hexInt  `json:"nonce"`
	BlockHash        string  `json:"blockHash"`
	BlockNumber      *hexInt `json:"blockNumber"`
	TransactionIndex *hexInt `json:"transactionIndex"`
	From             string  `json:"from"`
	To               string  `json:"to"`
	Value            hexBig  `json:"value"`
	Gas              hexInt  `json:"gas"`
	GasPrice         hexBig  `json:"gasPrice"`
	Input            string  `json:"input"`
}

type hexInt int

// ParseInt parse string value to int
func ParseInt(value string) (int, error) {
	i, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil {
		return 0, err
	}

	return int(i), nil
}

func (i *hexInt) UnmarshalJSON(data []byte) error {
	result, err := ParseInt(string(bytes.Trim(data, `"`)))
	*i = hexInt(result)

	return err
}

type hexBig big.Int

func (i *hexBig) UnmarshalJSON(data []byte) error {
	result, err := ethrpc.ParseBigInt(string(bytes.Trim(data, `"`)))
	*i = hexBig(result)

	return err
}

type ProxyBlockWithTransactions struct {
	Number           hexInt             `json:"number"`
	Hash             string             `json:"hash"`
	ParentHash       string             `json:"parentHash"`
	Nonce            string             `json:"nonce"`
	Sha3Uncles       string             `json:"sha3Uncles"`
	LogsBloom        string             `json:"logsBloom"`
	TransactionsRoot string             `json:"transactionsRoot"`
	StateRoot        string             `json:"stateRoot"`
	Miner            string             `json:"miner"`
	Difficulty       hexBig             `json:"difficulty"`
	TotalDifficulty  hexBig             `json:"totalDifficulty"`
	ExtraData        string             `json:"extraData"`
	Size             hexInt             `json:"size"`
	GasLimit         hexInt             `json:"gasLimit"`
	GasUsed          hexInt             `json:"gasUsed"`
	Timestamp        hexInt             `json:"timestamp"`
	Uncles           []string           `json:"uncles"`
	Transactions     []ProxyTransaction `json:"transactions"`
}

func (proxy *ProxyBlockWithTransactions) ToBlock() *ethrpc.Block {
	block := *(*Block)(unsafe.Pointer(proxy))

	ethrpcBlock := &ethrpc.Block{
		Number:           block.Number,
		Hash:             block.Hash,
		ParentHash:       block.ParentHash,
		Nonce:            block.Nonce,
		Sha3Uncles:       block.Sha3Uncles,
		LogsBloom:        block.LogsBloom,
		TransactionsRoot: block.TransactionsRoot,
		StateRoot:        block.StateRoot,
		Miner:            block.Miner,
		Difficulty:       block.Difficulty,
		TotalDifficulty:  block.TotalDifficulty,
		ExtraData:        block.ExtraData,
		Size:             block.Size,
		GasLimit:         block.GasLimit,
		GasUsed:          block.GasUsed,
		Timestamp:        block.Timestamp,
		Uncles:           block.Uncles,
		Transactions:     make([]ethrpc.Transaction, len(block.Transactions)),
	}

	// copy transactions
	for i, tx := range block.Transactions {
		ethrpcBlock.Transactions[i] = ethrpc.Transaction{
			Hash:             tx.Hash,
			Nonce:            tx.Nonce,
			BlockHash:        tx.BlockHash,
			BlockNumber:      tx.BlockNumber,
			TransactionIndex: tx.TransactionIndex,
			From:             tx.From,
			To:               tx.To,
			Value:            tx.Value,
			Gas:              tx.Gas,
			GasPrice:         tx.GasPrice,
			Input:            tx.Input,
		}
	}

	return ethrpcBlock
}
