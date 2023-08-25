package common

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	secp256k1 "github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/cosmos/cosmos-sdk/crypto/codec"

	eth "github.com/ethereum/go-ethereum/crypto"
	"github.com/tendermint/tendermint/crypto"

	"github.com/zeta-chain/zetacore/common/cosmos"
)

// PubKey is bech32 encoded string

type (
	PubKey  string
	PubKeys []PubKey
)

var EmptyPubKey PubKey

// EmptyPubKeySet
//var EmptyPubKeySet PubKeySet

func NewPubKey(key string) (PubKey, error) {
	if len(key) == 0 {
		return EmptyPubKey, nil
	}
	_, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, key)
	if err != nil {
		return EmptyPubKey, fmt.Errorf("%s is not bech32 encoded pub key,err : %w", key, err)
	}
	return PubKey(key), nil
}

// NewPubKeyFromCrypto
func NewPubKeyFromCrypto(pk crypto.PubKey) (PubKey, error) {
	tmp, err := codec.FromTmPubKeyInterface(pk)
	if err != nil {
		return EmptyPubKey, fmt.Errorf("fail to create PubKey from crypto.PubKey,err:%w", err)
	}
	s, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, tmp)
	if err != nil {
		return EmptyPubKey, fmt.Errorf("fail to create PubKey from crypto.PubKey,err:%w", err)
	}
	return PubKey(s), nil
}

// Equals check whether two are the same
func (pubKey PubKey) Equals(pubKey1 PubKey) bool {
	return pubKey == pubKey1
}

// IsEmpty to check whether it is empty
func (pubKey PubKey) IsEmpty() bool {
	return len(pubKey) == 0
}

// String stringer implementation
func (pubKey PubKey) String() string {
	return string(pubKey)
}

// GetAddress will return an address for the given chain
func (pubKey PubKey) GetAddress(chain Chain) (Address, error) {
	if pubKey.IsEmpty() {
		return NoAddress, nil
	}

	if IsEVMChain(chain.ChainId) {
		// retrieve compressed pubkey bytes from bechh32 encoded str
		pk, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, string(pubKey))
		if err != nil {
			return NoAddress, err
		}
		// parse compressed bytes removing 5 first bytes (amino encoding) to get uncompressed
		pub, err := secp256k1.ParsePubKey(pk.Bytes())
		if err != nil {
			return NoAddress, err
		}
		str := strings.ToLower(eth.PubkeyToAddress(*pub.ToECDSA()).String())
		return NewAddress(str, chain)
	}
	return NoAddress, nil
}

func (pubKey PubKey) GetZetaAddress() (cosmos.AccAddress, error) {
	addr, err := pubKey.GetAddress(ZetaChain())
	if err != nil {
		return nil, err
	}
	return cosmos.AccAddressFromBech32(addr.String())
}

// MarshalJSON to Marshals to JSON using Bech32
func (pubKey PubKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(pubKey.String())
}

// UnmarshalJSON to Unmarshal from JSON assuming Bech32 encoding
func (pubKey *PubKey) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	pk, err := NewPubKey(s)
	if err != nil {
		return err
	}
	*pubKey = pk
	return nil
}

func (pks PubKeys) Valid() error {
	for _, pk := range pks {
		if _, err := NewPubKey(pk.String()); err != nil {
			return err
		}
	}
	return nil
}

func (pks PubKeys) Contains(pk PubKey) bool {
	for _, p := range pks {
		if p.Equals(pk) {
			return true
		}
	}
	return false
}

// Equals check whether two pub keys are identical
func (pks PubKeys) Equals(newPks PubKeys) bool {
	if len(pks) != len(newPks) {
		return false
	}
	source := append(pks[:0:0], pks...)
	dest := append(newPks[:0:0], newPks...)
	// sort both lists
	sort.Slice(source[:], func(i, j int) bool {
		return source[i].String() < source[j].String()
	})
	sort.Slice(dest[:], func(i, j int) bool {
		return dest[i].String() < dest[j].String()
	})
	for i := range source {
		if !source[i].Equals(dest[i]) {
			return false
		}
	}
	return true
}

// String implement stringer interface
func (pks PubKeys) String() string {
	strs := make([]string, len(pks))
	for i := range pks {
		strs[i] = pks[i].String()
	}
	return strings.Join(strs, ", ")
}

func (pks PubKeys) Strings() []string {
	allStrings := make([]string, len(pks))
	for i, pk := range pks {
		allStrings[i] = pk.String()
	}
	return allStrings
}

// ConvertAndEncode converts from a base64 encoded byte string to hex or base32 encoded byte string and then to bech32
func ConvertAndEncode(hrp string, data []byte) (string, error) {
	converted, err := bech32.ConvertBits(data, 8, 5, true)
	if err != nil {
		return "", fmt.Errorf("encoding bech32 failed,%w", err)
	}
	return bech32.Encode(hrp, converted)
}

// NewPubKeySet create a new instance of PubKeySet , which contains two keys
func NewPubKeySet(secp256k1, ed25519 PubKey) PubKeySet {
	return PubKeySet{
		Secp256k1: secp256k1,
		Ed25519:   ed25519,
	}
}
