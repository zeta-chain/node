package sample

import (
	"errors"
	"hash/fnv"
	"math/rand"
	"strconv"
	"testing"

	"github.com/zeta-chain/zetacore/cmd/zetacored/config"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/cosmos"
)

var ErrSample = errors.New("sample error")

func newRandFromSeed(s int64) *rand.Rand {
	// #nosec G404 test purpose - weak randomness is not an issue here
	return rand.New(rand.NewSource(s))
}

func newRandFromStringSeed(t *testing.T, s string) *rand.Rand {
	h := fnv.New64a()
	_, err := h.Write([]byte(s))
	require.NoError(t, err)
	return newRandFromSeed(int64(h.Sum64()))
}

// PubKey returns a sample account PubKey
func PubKey(r *rand.Rand) cryptotypes.PubKey {
	seed := []byte(strconv.Itoa(r.Int()))
	return ed25519.GenPrivKeyFromSecret(seed).PubKey()
}

// Bech32AccAddress returns a sample account address
func Bech32AccAddress() sdk.AccAddress {
	pk := ed25519.GenPrivKey().PubKey()
	addr := pk.Address()
	return sdk.AccAddress(addr)
}

// AccAddress returns a sample account address in string
func AccAddress() string {
	pk := ed25519.GenPrivKey().PubKey()
	addr := pk.Address()
	return sdk.AccAddress(addr).String()
}

func ConsAddress() sdk.ConsAddress {
	return sdk.ConsAddress(PubKey(newRandFromSeed(1)).Address())
}

// ValAddress returns a sample validator operator address
func ValAddress(r *rand.Rand) sdk.ValAddress {
	return sdk.ValAddress(PubKey(r).Address())
}

// Validator returns a sample staking validator
func Validator(t testing.TB, r *rand.Rand) stakingtypes.Validator {
	seed := []byte(strconv.Itoa(r.Int()))
	val, err := stakingtypes.NewValidator(
		ValAddress(r),
		ed25519.GenPrivKeyFromSecret(seed).PubKey(),
		stakingtypes.Description{})
	require.NoError(t, err)
	return val
}

// PubKeyString returns a sample public key string
func PubKeyString() string {
	priKey := ed25519.GenPrivKey()
	s, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, priKey.PubKey())
	if err != nil {
		panic(err)
	}
	pubkey, err := common.NewPubKey(s)
	if err != nil {
		panic(err)
	}
	return pubkey.String()
}

// PrivKeyAddressPair returns a private key, address pair
func PrivKeyAddressPair() (*ed25519.PrivKey, sdk.AccAddress) {
	privKey := ed25519.GenPrivKey()
	addr := privKey.PubKey().Address()

	return privKey, sdk.AccAddress(addr)
}

// EthAddress returns a sample ethereum address
func EthAddress() ethcommon.Address {
	return ethcommon.BytesToAddress(sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()).Bytes())
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

// Coins returns a sample sdk.Coins
func Coins() sdk.Coins {
	return sdk.NewCoins(sdk.NewCoin(config.BaseDenom, sdk.NewInt(42)))
}
