package types

import (
	"crypto/sha256"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

// GasStabilityPoolAddress returns the address of the gas stability pool
func GasStabilityPoolAddress() sdk.AccAddress {
	hash := sha256.Sum256([]byte("gas_stability_pool"))
	return hash[:20]
}

// GasStabilityPoolAddressEVM returns the address of the gas stability pool in EVM format
func GasStabilityPoolAddressEVM() ethcommon.Address {
	return ethcommon.BytesToAddress(GasStabilityPoolAddress())
}
