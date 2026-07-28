package snapshot

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/types/bech32"
)

// addressLen is the byte length of a ZetaChain (eth-style) account address.
const addressLen = 20

// Canonical collapses any address form for the same underlying account to a
// single key: bech32 zeta1... (account), bech32 zetavaloper1... (validator
// operator) and 0x hex all reduce to the 20-byte account address as lowercase
// 0x-prefixed hex. This is what folds a validator's self-delegation and
// commission into its operator account without double-crediting.
//
// Canonical is strict: it errors on anything that is not a 20-byte account
// address, so it is used for inputs that must be 20-byte (validator operators,
// pins, module accounts). For bank holders and delegators, use Classify.
func Canonical(addr string) (string, error) {
	raw, err := decodeAddr(addr)
	if err != nil {
		return "", err
	}
	return canonicalFromBytes(raw)
}

// Classify decodes a holder address and returns its canonical 20-byte key.
// ok is false when the address is well-formed but not a 20-byte account (e.g. a
// 32-byte module-derived / group / interchain account): such holders have no
// eth-style keypair to claim with, so they are not attributed and fall through
// to the remainder. err is returned only for genuinely malformed input.
func Classify(addr string) (canon string, ok bool, err error) {
	raw, err := decodeAddr(addr)
	if err != nil {
		return "", false, err
	}
	if len(raw) != addressLen {
		return "", false, nil
	}
	return "0x" + hex.EncodeToString(raw), true, nil
}

// decodeAddr decodes a 0x-hex or bech32 address to its raw account bytes.
func decodeAddr(addr string) ([]byte, error) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return nil, fmt.Errorf("empty address")
	}
	if strings.HasPrefix(addr, "0x") || strings.HasPrefix(addr, "0X") {
		raw, err := hex.DecodeString(addr[2:])
		if err != nil {
			return nil, fmt.Errorf("decode hex address %q: %w", addr, err)
		}
		return raw, nil
	}
	_, raw, err := bech32.DecodeAndConvert(addr)
	if err != nil {
		return nil, fmt.Errorf("decode bech32 address %q: %w", addr, err)
	}
	return raw, nil
}

// canonicalFromBytes renders raw account-address bytes as lowercase 0x hex.
func canonicalFromBytes(raw []byte) (string, error) {
	if len(raw) != addressLen {
		return "", fmt.Errorf("expected %d-byte address, got %d bytes", addressLen, len(raw))
	}
	return "0x" + hex.EncodeToString(raw), nil
}
