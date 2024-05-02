package math

import (
	"math/big"
)

// Percentage calculates the percentage of A over B.
func Percentage(a, b *big.Int) *big.Float {
	// if a is zero, return nil to avoid division by zero
	if b.Cmp(big.NewInt(0)) == 0 {
		return nil
	}

	// convert a and a to big.Float
	floatA := new(big.Float).SetInt(a)
	floatB := new(big.Float).SetInt(b)

	// calculate the percentage of a over b
	percentage := new(big.Float).Quo(floatA, floatB)
	percentage.Mul(percentage, big.NewFloat(100))

	return percentage
}
