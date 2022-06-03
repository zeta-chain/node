package common

import (
	"fmt"
	"strings"

	eth "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/common/cosmos"
)

type Address string

var NoAddress Address = Address("")

const ETHAddressLen = 42

// NewAddress create a new Address. Supports Ethereum, BSC, Polygon
func NewAddress(address string, chain Chain) (Address, error) {

	// Check is eth address
	if chain == ETHChain || chain == BSCChain || chain == POLYGONChain || chain == ROPSTENChain {
		if eth.IsHexAddress(address) {
			return Address(address), nil
		}
	}

	return NoAddress, fmt.Errorf("address format not supported: %s", address)
}

//func (addr Address) IsChain(chain Chain) bool {
//	switch chain {
//	case ETHChain:
//		return strings.HasPrefix(addr.String(), "0x")
//	default:
//		return false
//	}
//}

//func (addr Address) GetChain() Chain {
//	for _, chain := range []Chain{ETHChain} {
//		if addr.IsChain(chain) {
//			return chain
//		}
//	}
//	return EmptyChain
//}

func (addr Address) GetNetwork(chain Chain) ChainNetwork {
	switch chain {
	case ETHChain:
		return GetCurrentChainNetwork()
	}
	return MockNet
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
