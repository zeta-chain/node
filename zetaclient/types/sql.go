package types

import (
	"github.com/btcsuite/btcd/btcjson"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"gorm.io/gorm"
)

// EVM Chain client types ----------------------------------->

const LastBlockNumId = 0xBEEF

type ReceiptSQLType struct {
	gorm.Model
	Nonce   int
	Receipt ethtypes.Receipt `gorm:"embedded"`
}

type TransactionSQLType struct {
	gorm.Model
	Nonce       int
	Transaction ethtypes.Transaction `gorm:"embedded"`
}

type LastBlockSQLType struct {
	gorm.Model
	Num int64
}

func ToReceiptSQLType(receipt *ethtypes.Receipt, nonce int) *ReceiptSQLType {
	return &ReceiptSQLType{
		Nonce:   nonce,
		Receipt: *receipt,
	}
}

func ToTransactionSQLType(transaction *ethtypes.Transaction, nonce int) *TransactionSQLType {
	return &TransactionSQLType{
		Nonce:       nonce,
		Transaction: *transaction,
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
