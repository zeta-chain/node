package types

import sdkmath "cosmossdk.io/math"

type GasFee struct {
	ShouldPayGas bool
	GasFeePaid   sdkmath.Int
}
