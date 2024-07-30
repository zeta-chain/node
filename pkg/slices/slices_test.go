package slices

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		function func(int) int
		expected []int
	}{
		{
			name:     "double",
			input:    []int{1, 2, 3, 4},
			function: func(x int) int { return x * 2 },
			expected: []int{2, 4, 6, 8},
		},
		{
			name:     "square",
			input:    []int{1, 2, 3, 4},
			function: func(x int) int { return x * x },
			expected: []int{1, 4, 9, 16},
		},
		{
			name:     "increment",
			input:    []int{1, 2, 3, 4},
			function: func(x int) int { return x + 1 },
			expected: []int{2, 3, 4, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Map(tt.input, tt.function)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestElementsMatch(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []int
		expected bool
	}{
		{
			name:     "same elements in same order",
			a:        []int{1, 2, 3, 4},
			b:        []int{1, 2, 3, 4},
			expected: true,
		},
		{
			name:     "same elements in different order",
			a:        []int{4, 3, 2, 1},
			b:        []int{1, 2, 3, 4},
			expected: true,
		},
		{
			name:     "different elements",
			a:        []int{1, 2, 3, 5},
			b:        []int{1, 2, 3, 4},
			expected: false,
		},
		{
			name:     "different lengths",
			a:        []int{1, 2, 3},
			b:        []int{1, 2, 3, 4},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ElementsMatch(tt.a, tt.b)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestDiff(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []int
		expected []int
	}{
		{
			name:     "elements in a not in b",
			a:        []int{1, 2, 3, 4},
			b:        []int{3, 4, 5, 6},
			expected: []int{1, 2},
		},
		{
			name:     "no elements in a not in b",
			a:        []int{3, 4},
			b:        []int{3, 4, 5, 6},
			expected: nil,
		},
		{
			name:     "all elements in a not in b",
			a:        []int{1, 2},
			b:        []int{3, 4},
			expected: []int{1, 2},
		},
		{
			name:     "empty a",
			a:        []int{},
			b:        []int{1, 2, 3, 4},
			expected: nil,
		},
		{
			name:     "empty b",
			a:        []int{1, 2, 3, 4},
			b:        []int{},
			expected: []int{1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Diff(tt.a, tt.b)
			require.Equal(t, tt.expected, result)
		})
	}
}
