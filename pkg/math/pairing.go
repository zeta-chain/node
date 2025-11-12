package math

import "math"

// CantorPair maps two numbers to a single number
// see: https://en.wikipedia.org/wiki/Pairing_function#Cantor_pairing_function
func CantorPair(x, y uint32) uint64 {
	sum := uint64(x) + uint64(y)
	return sum*(sum+1)/2 + uint64(y)
}

// CantorUnpair maps a single number to two numbers
func CantorUnpair(z uint64) (uint32, uint32) {
	w := uint64((math.Sqrt(float64(8*z+1)) - 1) / 2)
	t := w * (w + 1) / 2
	y := z - t
	x := uint32(w - y)

	return x, uint32(y)
}
