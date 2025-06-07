package sui

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/pkg/errors"
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

// DecodeAddress converts a byte slice into a Sui address string (0x-prefixed)
func DecodeAddress(data []byte) string {
	return "0x" + hex.EncodeToString(data)
}

// ValidateAddress checks whether the input string is a valid Sui address
// For WithdrawAndCall, the receiver is the target package ID. It follows same format, so we use same validation for both
func ValidateAddress(addr string) error {
	addr, hasPrefix := strings.CutPrefix(addr, "0x")
	switch {
	case !hasPrefix:
		return errors.New("address must start with 0x")
	case addr != strings.ToLower(addr):
		return errors.New("address must be lowercase")
	case len(addr) != 64:
		// accept full Sui address format only to make the validation easier
		return errors.New("address must be 64 characters")
	}

	if _, err := hex.DecodeString(addr); err != nil {
		return errors.Wrapf(err, "address %s is not valid hex", addr)
	}

	return nil
}
