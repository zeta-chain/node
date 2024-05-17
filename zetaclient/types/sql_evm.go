package types

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gorm.io/gorm"
)

// EVM Chain observer types ----------------------------------->

const LastBlockNumID = 0xBEEF

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
	ChainID *big.Int `gorm:"embedded"`
	Nonce   uint64
	To      *common.Address
	Hash    common.Hash

	// Serialized go-ethereum transaction
	TransactionData []byte
}

// Relational mapping types:

type ReceiptSQLType struct {
	gorm.Model
	Identifier string
	Receipt    ReceiptDB `gorm:"embedded"`
}

type TransactionSQLType struct {
	gorm.Model
	Identifier  string
	Transaction TransactionDB `gorm:"embedded"`
}

type LastBlockSQLType struct {
	gorm.Model
	Num uint64
}

// Type translation functions:

func ToReceiptDBType(receipt *ethtypes.Receipt) (ReceiptDB, error) {
	logs, err := json.Marshal(receipt.Logs)
	if err != nil {
		return ReceiptDB{}, err
	}
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
	}, nil
}

func FromReceiptDBType(receipt ReceiptDB) (*ethtypes.Receipt, error) {
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
	err := json.Unmarshal(receipt.Logs, &res.Logs)
	return res, err
}

func ToReceiptSQLType(receipt *ethtypes.Receipt, index string) (*ReceiptSQLType, error) {
	r, err := ToReceiptDBType(receipt)
	if err != nil {
		return nil, err
	}
	return &ReceiptSQLType{
		Identifier: index,
		Receipt:    r,
	}, nil
}

func ToTransactionDBType(transaction *ethtypes.Transaction) (TransactionDB, error) {
	data, err := transaction.MarshalBinary()
	if err != nil {
		return TransactionDB{}, err
	}
	return TransactionDB{
		Type:            transaction.Type(),
		ChainID:         transaction.ChainId(),
		Nonce:           transaction.Nonce(),
		To:              transaction.To(),
		Hash:            transaction.Hash(),
		TransactionData: data,
	}, nil
}

func FromTransactionDBType(transaction TransactionDB) (*ethtypes.Transaction, error) {
	res := &ethtypes.Transaction{}
	err := res.UnmarshalBinary(transaction.TransactionData)
	return res, err
}

func ToTransactionSQLType(transaction *ethtypes.Transaction, index string) (*TransactionSQLType, error) {
	trans, err := ToTransactionDBType(transaction)
	if err != nil {
		return nil, err
	}
	return &TransactionSQLType{
		Identifier:  index,
		Transaction: trans,
	}, nil
}

func ToLastBlockSQLType(lastBlock uint64) *LastBlockSQLType {
	return &LastBlockSQLType{
		Model: gorm.Model{ID: LastBlockNumID},
		Num:   lastBlock,
	}
}
