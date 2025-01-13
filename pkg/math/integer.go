package math

import (
	"math"
	"math/big"
)

// IncreaseIntByPercent is a function that increases integer by a percentage.
func IncreaseIntByPercent(value int64, percent uint32) int64 {
	if percent == 0 {
		return value
	}

	if value < 0 {
		return -IncreaseIntByPercent(-value, percent)
	}

	bigValue := big.NewInt(value)
	bigPercent := big.NewInt(int64(percent))

	// product = value * percent
	product := new(big.Int).Mul(bigValue, bigPercent)

	// dividing product by 100
	product.Div(product, big.NewInt(100))

	// result = original value + product
	result := new(big.Int).Add(bigValue, product)

	// be mindful if result > MaxInt64
	if result.Cmp(big.NewInt(math.MaxInt64)) > 0 {
		return math.MaxInt64
	}
	return result.Int64()
}
