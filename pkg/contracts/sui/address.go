package sui

import (
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

// EncodeAddress encodes Sui address to bytes
func EncodeAddress(addr string) ([]byte, error) {
	if !strings.HasPrefix(addr, "0x") {
		return nil, errors.New("address must start with 0x")
	}

	hexPart := addr[2:]

	if len(hexPart) == 0 {
		return nil, errors.New("address must not be empty")
	}

	if len(hexPart) > 64 {
		return nil, errors.New("address must be 64 characters or less")
	}

	return hex.DecodeString(fmt.Sprintf("%064s", hexPart))
}

// DecodeAddress converts a byte slice into a Sui address string (0x-prefixed, 64-char hex)
func DecodeAddress(data []byte) (string, error) {
	if len(data) > 32 {
		return "", errors.New("address must be 32 bytes or less")
	}

	// Left-pad with zeroes to make it 32 bytes (64 hex characters)
	padded := make([]byte, 32)
	copy(padded[32-len(data):], data)

	return "0x" + hex.EncodeToString(padded), nil
}

// ValidAddress checks whether the input string is a valid Sui address
func ValidAddress(addr string) error {
	if !strings.HasPrefix(addr, "0x") {
		return errors.New("address must start with 0x")
	}
	hexPart := addr[2:]

	if len(hexPart) == 0 {
		return errors.New("address must not be empty")
	}

	if len(hexPart) > 64 {
		return errors.New("address must be 64 characters or less")
	}

	_, err := hex.DecodeString(fmt.Sprintf("%064s", hexPart))
	return errors.Wrapf(err, "address %s is not valid hex", addr)
}
