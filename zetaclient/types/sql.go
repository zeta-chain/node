package types

import (
	"encoding/json"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gorm.io/gorm"
	"math/big"
)

// EVM Chain client types ----------------------------------->

const LastBlockNumId = 0xBEEF

// ReceiptDB : A modified receipt struct that the relational mapping can translate
type ReceiptDB struct {
	// Consensus fields: These fields are defined by the Yellow Paper
	Type              uint8
	PostState         []byte
	Status            uint64
	CumulativeGasUsed uint64
	Bloom             []byte
	Logs              []byte

	// Implementation fields: These fields are added by geth when processing a transaction.
	// They are stored in the chain database.
	TxHash          common.Hash
	ContractAddress common.Address
	GasUsed         uint64

	// Inclusion information: These fields provide information about the inclusion of the
	// transaction corresponding to this receipt.
	BlockHash        common.Hash
	BlockNumber      *big.Int `gorm:"embedded"`
	TransactionIndex uint
}

// TransactionDB : A modified Transaction struct that the relational mapping can translate.
// Inner transaction data is defined as an interface from eth types, so it will be serialized and stored as bytes.
type TransactionDB struct {
	// Data that can be used for queries
	Type    byte
	ChainId *big.Int `gorm:"embedded"`
	Nonce   uint64
	To      *common.Address
	Hash    common.Hash

	// Serialized go-ethereum transaction
	TransactionData []byte
}

// Relational mapping types:

type ReceiptSQLType struct {
	gorm.Model
	Nonce   int
	Receipt ReceiptDB `gorm:"embedded"`
}

type TransactionSQLType struct {
	gorm.Model
	Nonce       int
	Transaction TransactionDB `gorm:"embedded"`
}

type LastBlockSQLType struct {
	gorm.Model
	Num int64
}

// Type translation functions:

func ToReceiptDBType(receipt *ethtypes.Receipt) ReceiptDB {
	logs, _ := json.Marshal(receipt.Logs)
	return ReceiptDB{
		Type:              receipt.Type,
		PostState:         receipt.PostState,
		Status:            receipt.Status,
		CumulativeGasUsed: receipt.CumulativeGasUsed,
		Bloom:             receipt.Bloom.Bytes(),
		Logs:              logs,
		TxHash:            receipt.TxHash,
		ContractAddress:   receipt.ContractAddress,
		GasUsed:           receipt.GasUsed,
		BlockHash:         receipt.BlockHash,
		BlockNumber:       receipt.BlockNumber,
		TransactionIndex:  receipt.TransactionIndex,
	}
}

func FromReceiptDBType(receipt ReceiptDB) *ethtypes.Receipt {
	res := &ethtypes.Receipt{
		Type:              receipt.Type,
		PostState:         receipt.PostState,
		Status:            receipt.Status,
		CumulativeGasUsed: receipt.CumulativeGasUsed,
		Bloom:             ethtypes.BytesToBloom(receipt.Bloom),
		Logs:              nil,
		TxHash:            receipt.TxHash,
		ContractAddress:   receipt.ContractAddress,
		GasUsed:           receipt.GasUsed,
		BlockHash:         receipt.BlockHash,
		BlockNumber:       receipt.BlockNumber,
		TransactionIndex:  receipt.TransactionIndex,
	}
	_ = json.Unmarshal(receipt.Logs, &res.Logs)
	return res
}

func ToReceiptSQLType(receipt *ethtypes.Receipt, nonce int) *ReceiptSQLType {
	return &ReceiptSQLType{
		Nonce:   nonce,
		Receipt: ToReceiptDBType(receipt),
	}
}

func ToTransactionDBType(transaction *ethtypes.Transaction) TransactionDB {
	data, _ := transaction.MarshalBinary()
	return TransactionDB{
		Type:            transaction.Type(),
		ChainId:         transaction.ChainId(),
		Nonce:           transaction.Nonce(),
		To:              transaction.To(),
		Hash:            transaction.Hash(),
		TransactionData: data,
	}
}

func FromTransactionDBType(transaction TransactionDB) *ethtypes.Transaction {
	res := &ethtypes.Transaction{}
	_ = res.UnmarshalBinary(transaction.TransactionData)
	return res
}

func ToTransactionSQLType(transaction *ethtypes.Transaction, nonce int) *TransactionSQLType {
	return &TransactionSQLType{
		Nonce:       nonce,
		Transaction: ToTransactionDBType(transaction),
	}
}

func ToLastBlockSQLType(lastBlock int64) *LastBlockSQLType {
	return &LastBlockSQLType{
		Model: gorm.Model{ID: LastBlockNumId},
		Num:   lastBlock,
	}
}

// BTC Chain client types ----------------------------------->

type PendingUTXOSQLType struct {
	gorm.Model
	Key  string
	UTXO btcjson.ListUnspentResult `gorm:"embedded"`
}
