package math

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IncreaseIntByPercent(t *testing.T) {
	for i, tt := range []struct {
		value    int64
		percent  uint32
		expected int64
	}{
		{value: 10, percent: 0, expected: 10},
		{value: 10, percent: 15, expected: 11},
		{value: 10, percent: 225, expected: 32},
		{value: math.MaxInt64 / 2, percent: 101, expected: math.MaxInt64},
		{value: -10, percent: 0, expected: -10},
		{value: -10, percent: 15, expected: -11},
		{value: -10, percent: 225, expected: -32},
		{value: -math.MaxInt64 / 2, percent: 101, expected: -math.MaxInt64},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			result := IncreaseIntByPercent(tt.value, tt.percent)
			assert.Equal(t, tt.expected, result)
		})
	}
}
