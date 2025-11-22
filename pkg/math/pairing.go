package math

import (
	"math"
	"math/big"
)

const (
	// MaxPairValue is the maximum supported value for a pair of numbers
	MaxPairValue = math.MaxUint32 / 2 // 2147483647
)

// CantorPair maps two numbers to a single number
// see: https://en.wikipedia.org/wiki/Pairing_function#Cantor_pairing_function
func CantorPair(x, y uint32) uint64 {
	sum := uint64(x) + uint64(y)
	return sum*(sum+1)/2 + uint64(y)
}

// CantorUnpair maps a single number to two numbers
// Note: this is currently only used for unit tests
func CantorUnpair(z uint64) (uint32, uint32) {
	zBig := new(big.Int).SetUint64(z)

	// w = (sqrt(8*z + 1) - 1) / 2
	eightZ := new(big.Int).Mul(big.NewInt(8), zBig)
	eightZ1 := new(big.Int).Add(eightZ, big.NewInt(1))

	// Calculate integer square root
	w := new(big.Int).Sqrt(eightZ1)
	w.Sub(w, big.NewInt(1))
	w.Div(w, big.NewInt(2))

	// t = w * (w + 1) / 2
	w1 := new(big.Int).Add(w, big.NewInt(1))
	t := new(big.Int).Mul(w, w1)
	t.Div(t, big.NewInt(2))

	// y = z - t
	y := new(big.Int).Sub(zBig, t)

	// x = w - y
	x := new(big.Int).Sub(w, y)

	// #nosec G115 e2e - always in range
	return uint32(x.Uint64()), uint32(y.Uint64())
}
