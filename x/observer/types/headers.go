package types

import (
	"bytes"
	"encoding/hex"
	"fmt"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// NewEthereumHeader returns a new HeaderData containing an Ethereum header
func NewEthereumHeader(header []byte) HeaderData {
	return HeaderData{
		Data: &HeaderData_EthereumHeader{
			EthereumHeader: header,
		},
	}
}

// ParentHash extracts the parent hash from the header
func (h HeaderData) ParentHash() ([]byte, error) {
	switch data := h.Data.(type) {
	case *HeaderData_EthereumHeader:
		var header ethtypes.Header
		if err := rlp.DecodeBytes(data.EthereumHeader, &header); err != nil {
			return nil, err
		}
		return header.ParentHash.Bytes(), nil
	default:
		return nil, ErrUnrecognizedBlockHeader
	}
}

// Validate performs a basic validation of the HeaderData
func (h HeaderData) Validate(blockHash []byte, height int64) error {
	switch data := h.Data.(type) {
	case *HeaderData_EthereumHeader:
		return validateEthereumHeader(data.EthereumHeader, blockHash, height)
	default:
		return ErrUnrecognizedBlockHeader
	}
}

// validateEthereumHeader performs a basic validation of the Ethereum header
func validateEthereumHeader(headerBytes []byte, blockHash []byte, height int64) error {
	// on ethereum the block header is ~538 bytes in RLP encoding
	if len(headerBytes) > 1024 {
		return fmt.Errorf("header too long (%d)", len(headerBytes))
	}

	// RLP encoded block header
	var header ethtypes.Header
	if err := rlp.DecodeBytes(headerBytes, &header); err != nil {
		return fmt.Errorf("cannot decode RLP (%s)")
	}
	if err := header.SanityCheck(); err != nil {
		return fmt.Errorf("sanity check failed (%s)")
	}
	if bytes.Compare(blockHash, header.Hash().Bytes()) != 0 {
		return fmt.Errorf("tx hash mismatch (%s) vs (%s)", hex.EncodeToString(blockHash), header.Hash().Hex())
	}
	if height != header.Number.Int64() {
		return fmt.Errorf("height mismatch (%d) vs (%d)", height, header.Number.Int64())
	}
	return nil
}
