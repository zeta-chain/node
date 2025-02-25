package chains

import (
	sdkmath "cosmossdk.io/math"
)

const (
	// defaultInboundFastConfirmationLiquidityDivisor is the default ZRC20 liquidity cap divisor for inbound fast confirmation
	// For example: given a liquidity cap of 1M, the fast confirmation cap is 1M / 4000 = 250.
	defaultInboundFastConfirmationLiquidityDivisor = uint64(4000)
)

var (
	// customInboundFastConfirmationLiquidityDivisorMap maps chainID to custom ZRC20 liquidity cap divisor for inbound fast confirmation.
	// This map is used to override the default divisor for specific chains.
	customInboundFastConfirmationLiquidityDivisorMap = map[int64]uint64{}
)

// CalcInboundFastConfirmationAmountCap calculates the amount cap for inbound fast confirmation.
func CalcInboundFastConfirmationAmountCap(chainID int64, liquidityCap sdkmath.Uint) sdkmath.Uint {
	divisor := getInboundFastConfirmationLiquidityDivisor(chainID)
	return liquidityCap.QuoUint64(divisor)
}

// getInboundFastConfirmationLiquidityDivisor returns the ZRC20 liquidity cap divisor for inbound fast confirmation.
// Default divisor will be used if there is no custom divisor for given chainID.
func getInboundFastConfirmationLiquidityDivisor(chainID int64) uint64 {
	divisor, found := customInboundFastConfirmationLiquidityDivisorMap[chainID]
	if found && divisor > 0 {
		return divisor
	}
	return defaultInboundFastConfirmationLiquidityDivisor
}
