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
	WalletConflicts []string
	Time            int64
	TimeReceived    int64
	Details         []byte // btcjson.GetTransactionDetailsResult
	Hex             string
}

type PendingUTXOSQLType struct {
	gorm.Model
	Key  string
	UTXO btcjson.ListUnspentResult `gorm:"embedded"`
}

type TransactionResultSQLType struct {
	gorm.Model
	Key string
	Tx  TransactionResultDB `gorm:"embedded"`
}

func ToTransactionResultDB(txResult btcjson.GetTransactionResult) (TransactionResultDB, error) {
	details, err := json.Marshal(txResult.Details)
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
		WalletConflicts: txResult.WalletConflicts,
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
		WalletConflicts: txResult.WalletConflicts,
		Time:            txResult.Time,
		TimeReceived:    txResult.TimeReceived,
		Details:         nil,
		Hex:             txResult.Hex,
	}
	err := json.Unmarshal(txResult.Details, &res.Details)
	return res, err
}
