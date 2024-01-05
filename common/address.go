package common

import (
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcutil"
	eth "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/common/cosmos"
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

func (addr Address) AccAddress() (cosmos.AccAddress, error) {
	return cosmos.AccAddressFromBech32(addr.String())
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
		return errors.New(fmt.Sprint(x))
	}
}

func DecodeBtcAddress(inputAddress string, chainId int64) (address btcutil.Address, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ConvertRecoverToError(r)
			err = fmt.Errorf("input address:%s,chainId:%d,err:%s", inputAddress, chainId, err.Error())
			return
		}
	}()
	chainParams, err := GetBTCChainParams(chainId)
	if err != nil {
		return nil, err
	}
	if chainParams == nil {
		return nil, fmt.Errorf("chain params not found")
	}
	oneIndex := strings.LastIndexByte(inputAddress, '1')
	if oneIndex > 1 {
		prefix := inputAddress[:oneIndex]
		ok := IsValidPrefix(prefix, chainId)
		if !ok {
			return nil, fmt.Errorf("invalid prefix:%s,chain-id:%d", prefix, chainId)
		}
		addressString := inputAddress[oneIndex+1:]
		if len(addressString) != 39 {
			return nil, fmt.Errorf("invalid address length:%d,inputaddress:%s", len(addressString), inputAddress)
		}
	}
	address, err = btcutil.DecodeAddress(inputAddress, chainParams)
	return
}
