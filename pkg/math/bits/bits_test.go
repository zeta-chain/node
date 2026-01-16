package math_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	zetabits "github.com/zeta-chain/node/pkg/math/bits"
)

func TestSetBit(t *testing.T) {
	tests := []struct {
		name     string
		initial  byte
		position uint8
		expected byte
	}{
		{
			name:     "set bit at position 0",
			initial:  0b00001000,
			position: 0,
			expected: 0b00001001,
		},
		{
			name:     "set bit at position 7",
			initial:  0b00001000,
			position: 7,
			expected: 0b10001000,
		},
		{
			name:     "out of range bit position (no effect)",
			initial:  0b00000000,
			position: 8, // Out of range
			expected: 0b00000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.initial
			zetabits.SetBit(&b, tt.position)
			require.Equal(t, tt.expected, b)
		})
	}
}

func TestIsBitSet(t *testing.T) {
	tests := []struct {
		name     string
		b        byte
		position uint8
		expected bool
	}{
		{
			name:     "bit 0 set",
			b:        0b00000001,
			position: 0,
			expected: true,
		},
		{
			name:     "bit 7 set",
			b:        0b10000000,
			position: 7,
			expected: true,
		},
		{
			name:     "bit 2 not set",
			b:        0b00000001,
			position: 2,
			expected: false,
		},
		{
			name:     "bit out of range",
			b:        0b00000001,
			position: 8, // Position out of range
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := zetabits.IsBitSet(tt.b, tt.position)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestGetBits(t *testing.T) {
	tests := []struct {
		name     string
		b        byte
		mask     byte
		expected byte
	}{
		{
			name:     "extract upper 3 bits",
			b:        0b11011001,
			mask:     0b11100000,
			expected: 0b110,
		},
		{
			name:     "extract middle 3 bits",
			b:        0b11011001,
			mask:     0b00011100,
			expected: 0b110,
		},
		{
			name:     "extract lower 3 bits",
			b:        0b11011001,
			mask:     0b00000111,
			expected: 0b001,
		},
		{
			name:     "extract no bits",
			b:        0b11011001,
			mask:     0b00000000,
			expected: 0b000,
		},
		{
			name:     "extract all bits",
			b:        0b11111111,
			mask:     0b11111111,
			expected: 0b11111111,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := zetabits.GetBits(tt.b, tt.mask)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestSetBits(t *testing.T) {
	tests := []struct {
		name     string
		b        byte
		mask     byte
		value    byte
		expected byte
	}{
		{
			name:     "set upper 3 bits",
			b:        0b00100001,
			mask:     0b11100000,
			value:    0b110,
			expected: 0b11000001,
		},
		{
			name:     "set middle 3 bits",
			b:        0b00100001,
			mask:     0b00011100,
			value:    0b101,
			expected: 0b00110101,
		},
		{
			name:     "set lower 3 bits",
			b:        0b11111100,
			mask:     0b00000111,
			value:    0b101,
			expected: 0b11111101,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := zetabits.SetBits(tt.b, tt.mask, tt.value)
			require.Equal(t, tt.expected, result)
		})
	}
}
