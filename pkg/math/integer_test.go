package math

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IncreaseIntByPercent(t *testing.T) {
	for i, tt := range []struct {
		value    int64
		percent  uint32
		round    bool
		expected int64
	}{
		{value: 10, percent: 0, round: false, expected: 10},
		{value: 10, percent: 15, round: false, expected: 11},
		{value: 10, percent: 15, round: true, expected: 12},
		{value: 10, percent: 14, round: false, expected: 11},
		{value: 10, percent: 14, round: true, expected: 11},
		{value: 10, percent: 200, round: false, expected: 30},
		{value: 10, percent: 200, round: true, expected: 30},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			result := IncreaseIntByPercent(tt.value, tt.percent, tt.round)
			assert.Equal(t, tt.expected, result)
		})
	}
}
