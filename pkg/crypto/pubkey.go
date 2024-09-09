package crypto

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	secp256k1 "github.com/btcsuite/btcd/btcec/v2"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	eth "github.com/ethereum/go-ethereum/crypto"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/cosmos"
)

// PubKey is bech32 encoded string

type (
	PubKey  string
	PubKeys []PubKey
)

var EmptyPubKey PubKey

// EmptyPubKeySet
//var EmptyPubKeySet PubKeySet

// NewPubKey create a new instance of PubKey
// key is bech32 encoded string

func GetAddressFromPubkeyString(pubkey string) (sdk.AccAddress, error) {
	cryptopub, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, pubkey)
	if err != nil {
		return nil, err
	}
	addr, err := sdk.AccAddressFromHexUnsafe(cryptopub.Address().String())
	if err != nil {
		return nil, err
	}
	return addr, nil
}
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
func (pubKey PubKey) GetAddress(chain chains.Chain) (chains.Address, error) {
	if chain.IsEVMChain() {
		return pubKey.GetEVMAddress()
	}
	return chains.NoAddress, nil
}

// GetEVMAddress will return the evm address
func (pubKey PubKey) GetEVMAddress() (chains.Address, error) {
	if pubKey.IsEmpty() {
		return chains.NoAddress, nil
	}

	// retrieve compressed pubkey bytes from bechh32 encoded str
	pk, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, string(pubKey))
	if err != nil {
		return chains.NoAddress, err
	}
	// parse compressed bytes removing 5 first bytes (amino encoding) to get uncompressed
	pub, err := secp256k1.ParsePubKey(pk.Bytes())
	if err != nil {
		return chains.NoAddress, err
	}
	str := strings.ToLower(eth.PubkeyToAddress(*pub.ToECDSA()).String())
	return chains.NewAddress(str), nil
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

func GetPubkeyBech32FromRecord(record *keyring.Record) (string, error) {
	pk, ok := record.PubKey.GetCachedValue().(cryptotypes.PubKey)
	if !ok {
		return "", errors.New("unable to cast any to cryptotypes.PubKey")
	}

	s, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pk)
	if err != nil {
		return "", err
	}
	return s, nil
}
