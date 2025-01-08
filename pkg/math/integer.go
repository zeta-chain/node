package math

import "math"

// IncreaseIntByPercent is a function that increases integer by a percentage.
// Example1: IncreaseIntByPercent(10, 15, true) = 10 * 1.15 = 12
// Example2: IncreaseIntByPercent(10, 15, false) = 10 + 10 * 0.15 = 11
func IncreaseIntByPercent(value int64, percent uint64, round bool) int64 {
	switch {
	case percent == 0:
		return value
	case percent%100 == 0:
		// optimization: a simple multiplication
		increase := value * int64(percent/100)
		return value + increase
	default:
		increase := float64(value) * float64(percent) / 100
		if round {
			return value + int64(math.Round(increase))
		}
		return value + int64(increase)
	}
}
