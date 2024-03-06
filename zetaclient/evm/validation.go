package evm

import (
	"encoding/hex"
	"fmt"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/onrik/ethrpc"
	"github.com/pkg/errors"
)

// ValidateEvmTxLog checks the basics of an EVM tx log
func ValidateEvmTxLog(vLog *ethtypes.Log, wantAddress ethcommon.Address, wantHash string, wantTopics int) error {
	if vLog.Removed {
		return fmt.Errorf("log is removed, chain reorg?")
	}
	if vLog.Address != wantAddress {
		return fmt.Errorf("log emitter address mismatch: want %s got %s", wantAddress.Hex(), vLog.Address.Hex())
	}
	if wantHash != "" && vLog.TxHash.Hex() != wantHash {
		return fmt.Errorf("log tx hash mismatch: want %s got %s", wantHash, vLog.TxHash.Hex())
	}
	if len(vLog.Topics) != wantTopics {
		return fmt.Errorf("number of topics mismatch: want %d got %d", wantTopics, len(vLog.Topics))
	}
	return nil
}

// ValidateEvmTransaction checks the basics of an EVM transaction
// Note: these checks are to ensure the transaction is well-formed
// and can be safely used for further processing by zetaclient
func ValidateEvmTransaction(tx *ethrpc.Transaction) error {
	if tx == nil {
		return fmt.Errorf("transaction is nil")
	}
	if tx.Hash == "" {
		return fmt.Errorf("transaction hash is empty")
	}
	if tx.Nonce < 0 {
		return fmt.Errorf("transaction nonce %d is negative", tx.Nonce)
	}
	if !ethcommon.IsHexAddress(tx.From) {
		return fmt.Errorf("transaction from %s is not a valid hex address", tx.From)
	}
	if tx.To != "" && !ethcommon.IsHexAddress(tx.To) {
		// To address can be empty for contract creation
		return fmt.Errorf("transaction to %s is not a valid hex address", tx.To)
	}
	if tx.Value.Sign() < 0 {
		return fmt.Errorf("transaction value %s is negative", tx.Value.String())
	}
	if tx.Gas < 0 {
		return fmt.Errorf("transaction gas %d is negative", tx.Gas)
	}
	if tx.GasPrice.Sign() < 0 {
		return fmt.Errorf("transaction gas price %s is negative", tx.GasPrice.String())
	}
	// remove '0x' prefix from input data to be consistent with ethclient
	tx.Input = strings.TrimPrefix(tx.Input, "0x")

	// tx input data should be hex encoded
	if _, err := hex.DecodeString(tx.Input); err != nil {
		return errors.Wrapf(err, "transaction input data is not hex encoded: %s", tx.Input)
	}

	// inclusion checks
	if tx.BlockNumber != nil {
		if *tx.BlockNumber <= 0 {
			return fmt.Errorf("transaction block number %d is not positive", *tx.BlockNumber)
		}
		if tx.BlockHash == "" {
			return fmt.Errorf("transaction block hash is empty")
		}
		if tx.TransactionIndex == nil {
			return fmt.Errorf("transaction index is nil")
		}
		if *tx.TransactionIndex < 0 {
			return fmt.Errorf("transaction index %d is negative", *tx.TransactionIndex)
		}
	}
	return nil
}
