package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceMedianValue(t *testing.T) {
	for _, tt := range []struct {
		name     string
		input    []int
		expected int
		inPlace  bool
	}{
		{
			name:     "empty",
			input:    nil,
			expected: 0,
			inPlace:  false,
		},
		{
			name:     "single",
			input:    []int{10},
			expected: 10,
		},
		{
			name:     "two",
			input:    []int{10, 20},
			expected: 15,
		},
		{
			name:     "even",
			input:    []int{30, 20, 10, 20},
			expected: 20,
		},
		{
			name:     "even in-place",
			input:    []int{30, 20, 10, 20},
			expected: 20,
			inPlace:  true,
		},
		{
			name:     "odd",
			input:    []int{5, 5, 6, 1, 1, 1, 4},
			expected: 4,
		},
		{
			name:     "odd in-place",
			input:    []int{1, 1, 1, 1, 7, 7, 7, 7},
			expected: 4,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// ASSERT
			// Given a copy of the input slice
			var snapshot []int
			for _, v := range tt.input {
				snapshot = append(snapshot, v)
			}

			// ACT
			out := SliceMedianValue(tt.input, tt.inPlace)

			// ASSERT
			assert.Equal(t, tt.expected, out)

			// Check that elements of the input slice are unchanged
			if !tt.inPlace {
				assert.Equal(t, snapshot, tt.input)
			}
		})
	}

}
