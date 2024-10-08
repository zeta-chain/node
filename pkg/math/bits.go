package math

import (
	"math/bits"
)

// SetBit sets the bit at the given position (0-7) in the byte to 1
func SetBit(b *byte, position uint8) {
	if position > 7 {
		return
	}
	*b |= 1 << position
}

// IsBitSet returns true if the bit at the given position (0-7) is set in the byte, false otherwise
func IsBitSet(b byte, position uint8) bool {
	if position > 7 {
		return false
	}
	bitMask := byte(1 << position)
	return b&bitMask != 0
}

// GetBits extracts the value of bits for a given mask
//
// Example: given b = 0b11011001 and mask = 0b11100000, the function returns 0b110
func GetBits(b byte, mask byte) byte {
	extracted := b & mask

	// get the number of trailing zero bits
	trailingZeros := bits.TrailingZeros8(mask)

	// remove trailing zeros
	return extracted >> trailingZeros
}

// SetBits sets the value to the bits specified in the mask
//
// Example: given b = 0b00100001 and mask = 0b11100000, and value = 0b110, the function returns 0b11000001
func SetBits(b byte, mask byte, value byte) byte {
	// get the number of trailing zero bits in the mask
	trailingZeros := bits.TrailingZeros8(mask)

	// shift the value left by the number of trailing zeros
	valueShifted := value << trailingZeros

	// clear the bits in 'b' that correspond to the mask
	bCleared := b &^ mask

	// Set the bits by ORing the cleared 'b' with the shifted value
	return bCleared | valueShifted
}
