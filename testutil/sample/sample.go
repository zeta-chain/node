package sample

import (
	"errors"
	"hash/fnv"
	"math/rand"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

var ErrSample = errors.New("sample error")

func newRandFromSeed(s int64) *rand.Rand {
	return rand.New(rand.NewSource(s))
}

func newRandFromStringSeed(t *testing.T, s string) *rand.Rand {
	h := fnv.New64a()
	_, err := h.Write([]byte(s))
	require.NoError(t, err)
	return newRandFromSeed(int64(h.Sum64()))
}

// AccAddress returns a sample account address
func AccAddress() string {
	pk := ed25519.GenPrivKey().PubKey()
	addr := pk.Address()
	return sdk.AccAddress(addr).String()
}

// PrivKeyAddressPair returns a private key, address pair
func PrivKeyAddressPair() (*ed25519.PrivKey, sdk.AccAddress) {
	privKey := ed25519.GenPrivKey()
	addr := privKey.PubKey().Address()

	return privKey, sdk.AccAddress(addr)
}

// EthAddress returns a sample ethereum address
func EthAddress() ethcommon.Address {
	return ethcommon.HexToAddress(AccAddress())
}

// Bytes returns a sample byte array
func Bytes() []byte {
	return []byte("sample")
}

// String returns a sample string
func String() string {
	return "sample"
}

// StringRandom returns a sample string with random alphanumeric characters
func StringRandom(r *rand.Rand, length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}
