package types

import (
	"gorm.io/gorm"
)

const CurrentTssID = 0xACEE

type CurrentTssSQLType struct {
	gorm.Model
	ID        int64
	TssPubkey string
}

func ToCurrentTssSQLType(tssPubkey string) *CurrentTssSQLType {
	return &CurrentTssSQLType{
		ID:        CurrentTssID,
		TssPubkey: tssPubkey,
	}
}

type FirstNonceToScanSQLType struct {
	gorm.Model
	ID         int64
	FirstNonce uint64
}

func ToFirstNonceToScanSQLType(chainID int64, firstNonceToScan uint64) *FirstNonceToScanSQLType {
	return &FirstNonceToScanSQLType{
		ID:         chainID,
		FirstNonce: firstNonceToScan,
	}
}
