package types

import (
	"gorm.io/gorm"
)

const (
	// LastBlockNumID is the identifier to access the last block number in the database
	LastBlockNumID = 0xBEEF

	// LastTxHashID is the identifier to access the last transaction hash in the database
	LastTxHashID = 0xBEF0
)

// LastBlockSQLType is a model for storing the last block number
type LastBlockSQLType struct {
	gorm.Model
	Num uint64
}

// LastTransactionSQLType is a model for storing the last transaction hash
type LastTransactionSQLType struct {
	gorm.Model
	Hash string
}

// AuxStringSQLType is a model for storing auxiliary string data
type AuxStringSQLType struct {
	gorm.Model
	Key   string `gorm:"column:key_name;uniqueIndex;not null"`
	Value string
}

// ToLastBlockSQLType converts a last block number to a LastBlockSQLType
func ToLastBlockSQLType(lastBlock uint64) *LastBlockSQLType {
	return &LastBlockSQLType{
		Model: gorm.Model{ID: LastBlockNumID},
		Num:   lastBlock,
	}
}

// ToLastTxHashSQLType converts a last transaction hash to a LastTransactionSQLType
func ToLastTxHashSQLType(lastTx string) *LastTransactionSQLType {
	return &LastTransactionSQLType{
		Model: gorm.Model{ID: LastTxHashID},
		Hash:  lastTx,
	}
}

// ToAuxStringSQLType converts given key and value to a AuxStringSQLType
func ToAuxStringSQLType(key, value string) *AuxStringSQLType {
	return &AuxStringSQLType{
		Key:   key,
		Value: value,
	}
}
