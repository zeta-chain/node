package common

import (
	"fmt"
	"strings"

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
func NewAddress(address string, chain Chain) (Address, error) {

	// Check is eth address
	if IsEVMChain(chain.ChainId) {
		if eth.IsHexAddress(address) {
			return Address(address), nil
		}
	} else if IsBitcoinChain(chain.ChainId) {

	}

	return NoAddress, fmt.Errorf("address format not supported: %s", address)
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
