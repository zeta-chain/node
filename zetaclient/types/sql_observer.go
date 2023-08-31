package types

import (
	"gorm.io/gorm"
)

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
