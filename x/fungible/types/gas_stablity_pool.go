package types

import (
	"crypto/sha256"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

// GasStabilityPoolAddress returns the address of the gas stability pool
func GasStabilityPoolAddress(chainID int64) sdk.AccAddress {
	hash := sha256.Sum256([]byte(fmt.Sprintf("gas_stability_pool/%d", chainID)))
	return hash[:20]
}

// GasStabilityPoolAddressEVM returns the address of the gas stability pool in EVM format
func GasStabilityPoolAddressEVM(chainID int64) ethcommon.Address {
	return ethcommon.BytesToAddress(GasStabilityPoolAddress(chainID))
}
