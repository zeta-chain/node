// Package encoder contains functions that encode/decode transactions' logical time and hash
// ("lt:hash").
package encoder

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/tonkeeper/tongo/ton"
)

var (
	ErrInvalidString = errors.New("invalid string format")
	ErrInvalidLT     = errors.New("invalid logical time")
	ErrInvalidHash   = errors.New("invalid hash")
)

// EncodeTx encodes a transaction's logical time and hash (given the transaction).
func EncodeTx(tx ton.Transaction) string {
	return EncodeHash(tx.Lt, ton.Bits256(tx.Hash()))
}

// EncodeHash encodes a transaction's logical time and hash (given the LT and the hash).
func EncodeHash(lt uint64, hash ton.Bits256) string {
	return fmt.Sprintf("%d:%s", lt, hash.Hex())
}

// DecodeHash decodes an encoded transaction hash into logical time and hash.
func DecodeHash(encoded string) (lt uint64, hash ton.Bits256, err error) {
	parts := strings.Split(encoded, ":")
	if len(parts) != 2 {
		return lt, hash, fmt.Errorf("%w: %q", ErrInvalidString, encoded)
	}

	lt, err = strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return lt, hash, fmt.Errorf("%w: %w", ErrInvalidLT, err)
	}

	err = hash.FromHex(parts[1])
	if err != nil {
		return lt, hash, fmt.Errorf("%w: %w", ErrInvalidHash, err)
	}

	return lt, hash, nil
}
