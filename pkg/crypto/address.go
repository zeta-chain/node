package crypto

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/zeta-chain/node/pkg/constant"
)

// IsEmptyAddress returns true if the address is empty
func IsEmptyAddress(address common.Address) bool {
	return address == (common.Address{}) || address.Hex() == constant.EVMZeroAddress
}
