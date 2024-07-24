package math

import "cosmossdk.io/math"

// IncreaseUintByPercent is a function that increases Uint by a percentage.
// Example: IncreaseUintByPercent(4000, 20) = 4000 * 1,2 = 4800
// Returns result and increase amount.
func IncreaseUintByPercent(amount math.Uint, percent uint64) (math.Uint, math.Uint) {
	switch {
	case percent == 0:
		// noop
		return amount, math.ZeroUint()
	case percent%100 == 0:
		// optimization: a simple multiplication
		increase := amount.MulUint64(percent / 100)
		return amount.Add(increase), increase
	default:
		increase := amount.MulUint64(percent).QuoUint64(100)
		return amount.Add(increase), increase
	}
}
