package types

import (
	"encoding/json"

	"github.com/btcsuite/btcd/btcjson"
	"gorm.io/gorm"
)

type TransactionResultDB struct {
	Amount          float64
	Fee             float64
	Confirmations   int64
	BlockHash       string
	BlockIndex      int64
	BlockTime       int64
	TxID            string
	WalletConflicts []byte
	Time            int64
	TimeReceived    int64
	Details         []byte // btcjson.GetTransactionDetailsResult
	Hex             string
}

type TransactionResultSQLType struct {
	gorm.Model
	Key string
	Tx  TransactionResultDB `gorm:"embedded"`
}

type OutboundHashSQLType struct {
	gorm.Model
	Key  string
	Hash string
}

func ToTransactionResultDB(txResult btcjson.GetTransactionResult) (TransactionResultDB, error) {
	details, err := json.Marshal(txResult.Details)
	if err != nil {
		return TransactionResultDB{}, err
	}
	conflicts, err := json.Marshal(txResult.WalletConflicts)
	if err != nil {
		return TransactionResultDB{}, err
	}
	return TransactionResultDB{
		Amount:          txResult.Amount,
		Fee:             txResult.Fee,
		Confirmations:   txResult.Confirmations,
		BlockHash:       txResult.BlockHash,
		BlockIndex:      txResult.BlockIndex,
		BlockTime:       txResult.BlockTime,
		TxID:            txResult.TxID,
		WalletConflicts: conflicts,
		Time:            txResult.Time,
		TimeReceived:    txResult.TimeReceived,
		Details:         details,
		Hex:             txResult.Hex,
	}, nil
}

func FromTransactionResultDB(txResult TransactionResultDB) (btcjson.GetTransactionResult, error) {
	res := btcjson.GetTransactionResult{
		Amount:          txResult.Amount,
		Fee:             txResult.Fee,
		Confirmations:   txResult.Confirmations,
		BlockHash:       txResult.BlockHash,
		BlockIndex:      txResult.BlockIndex,
		BlockTime:       txResult.BlockTime,
		TxID:            txResult.TxID,
		WalletConflicts: nil,
		Time:            txResult.Time,
		TimeReceived:    txResult.TimeReceived,
		Details:         nil,
		Hex:             txResult.Hex,
	}
	err := json.Unmarshal(txResult.WalletConflicts, &res.WalletConflicts)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(txResult.Details, &res.Details)
	return res, err
}

func ToTransactionResultSQLType(txResult btcjson.GetTransactionResult, key string) (TransactionResultSQLType, error) {
	txDB, err := ToTransactionResultDB(txResult)
	if err != nil {
		return TransactionResultSQLType{}, err
	}
	return TransactionResultSQLType{
		Key: key,
		Tx:  txDB,
	}, nil
}

func FromTransactionResultSQLType(txSQL TransactionResultSQLType) (btcjson.GetTransactionResult, error) {
	return FromTransactionResultDB(txSQL.Tx)
}

func ToOutboundHashSQLType(hash string, key string) OutboundHashSQLType {
	return OutboundHashSQLType{
		Key:  key,
		Hash: hash,
	}
}
