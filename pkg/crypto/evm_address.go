package crypto

import (
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/zeta-chain/node/pkg/constant"
)

// IsEmptyAddress returns true if the address is empty
func IsEmptyAddress(address common.Address) bool {
	return address == (common.Address{}) || address.Hex() == constant.EVMZeroAddress
}

// IsEVMAddress returns true if the string is an EVM address
// independently of the checksum format
func IsEVMAddress(address string) bool {
	return len(address) == 42 && strings.HasPrefix(address, "0x") && common.IsHexAddress(address)
}

// IsChecksumAddress returns true if the EVM address string is a valid checksum address
// See https://eips.ethereum.org/EIPS/eip-55
func IsChecksumAddress(address string) bool {
	return address == common.HexToAddress(address).Hex()
}

// ToChecksumAddress returns the checksum address of the given EVM address
func ToChecksumAddress(address string) string {
	return common.HexToAddress(address).Hex()
}
