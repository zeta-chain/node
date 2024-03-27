package common

import (
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcutil"
	eth "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/common/bitcoin"
)

type Address string

var NoAddress Address

var (
	DeadAddress = eth.HexToAddress("0xdEAD000000000000000042069420694206942069")
)

const ETHAddressLen = 42

// NewAddress create a new Address. Supports Ethereum, BSC, Polygon
func NewAddress(address string) Address {
	// Check is eth address
	if eth.IsHexAddress(address) {
		return Address(address)
	}
	return NoAddress
}

func (addr Address) Equals(addr2 Address) bool {
	return strings.EqualFold(addr.String(), addr2.String())
}

func (addr Address) IsEmpty() bool {
	return strings.TrimSpace(addr.String()) == ""
}

func (addr Address) String() string {
	return string(addr)
}

func ConvertRecoverToError(r interface{}) error {
	switch x := r.(type) {
	case string:
		return errors.New(x)
	case error:
		return x
	default:
		return fmt.Errorf("%v", x)
	}
}

// DecodeBtcAddress decodes a BTC address from a given string and chainID
func DecodeBtcAddress(inputAddress string, chainID int64) (address btcutil.Address, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ConvertRecoverToError(r)
			err = fmt.Errorf("input address:%s, chainId:%d, err:%s", inputAddress, chainID, err.Error())
			return
		}
	}()
	chainParams, err := GetBTCChainParams(chainID)
	if err != nil {
		return nil, err
	}
	if chainParams == nil {
		return nil, fmt.Errorf("chain params not found")
	}
	// test taproot address type
	address, err = bitcoin.DecodeTaprootAddress(inputAddress)
	if err == nil {
		if address.IsForNet(chainParams) {
			return address, nil
		}
		return nil, fmt.Errorf("address %s is not for network %s", inputAddress, chainParams.Name)
	}
	// test taproot address failed; continue testing other types: P2WSH, P2WPKH, P2SH, P2PKH
	address, err = btcutil.DecodeAddress(inputAddress, chainParams)
	if err != nil {
		return nil, fmt.Errorf("decode address failed: %s, for input address %s", err.Error(), inputAddress)
	}
	ok := address.IsForNet(chainParams)
	if !ok {
		return nil, fmt.Errorf("address %s is not for network %s", inputAddress, chainParams.Name)
	}
	return
}

// IsBtcAddressSupported returns true if the given BTC address is supported
func IsBtcAddressSupported(addr btcutil.Address) bool {
	switch addr.(type) {
	// P2TR address
	case *bitcoin.AddressTaproot,
		// P2WSH address
		*btcutil.AddressWitnessScriptHash,
		// P2WPKH address
		*btcutil.AddressWitnessPubKeyHash,
		// P2SH address
		*btcutil.AddressScriptHash,
		// P2PKH address
		*btcutil.AddressPubKeyHash:
		return true
	}
	return false
}
