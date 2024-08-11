package types

import (
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/zeta-chain/zetacore/pkg/crypto"
)

// GetEVMRevertAddress returns the EVM revert address
// if the revert address is not a valid address, it returns false
func (r RevertOptions) GetEVMRevertAddress() (ethcommon.Address, bool) {
	addr := ethcommon.HexToAddress(r.RevertAddress)
	return addr, !crypto.IsEmptyAddress(addr)
}

// GetEVMAbortAddress returns the EVM revert address
// if the revert address is not a valid address, it returns false
func (r RevertOptions) GetEVMAbortAddress() (ethcommon.Address, bool) {
	addr := ethcommon.HexToAddress(r.AbortAddress)
	return addr, !crypto.IsEmptyAddress(addr)
}
