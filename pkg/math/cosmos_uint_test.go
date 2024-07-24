package math

import (
	"fmt"
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/assert"
)

func TestIncreaseUintByPercent(t *testing.T) {
	for i, tt := range []struct {
		in       math.Uint
		percent  uint64
		expected math.Uint
	}{
		{in: math.NewUint(444), percent: 0, expected: math.NewUint(444)},
		{in: math.NewUint(100), percent: 4, expected: math.NewUint(104)},
		{in: math.NewUint(100), percent: 100, expected: math.NewUint(200)},
		{in: math.NewUint(4000), percent: 50, expected: math.NewUint(6000)},
		{in: math.NewUint(2500), percent: 100, expected: math.NewUint(5000)},
		{in: math.NewUint(10000), percent: 33, expected: math.NewUint(13300)},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			actual, increase := IncreaseUintByPercent(tt.in, tt.percent)

			original := actual.Sub(increase).Uint64()

			assert.Equal(t, int(tt.expected.Uint64()), int(actual.Uint64()))
			assert.Equal(t, int(tt.in.Uint64()), int(original))

			t.Logf(
				"input: %d, percent: %d, expected: %d, actual: %d, increase: %d",
				tt.in.Uint64(),
				tt.percent,
				tt.expected.Uint64(),
				actual.Uint64(),
				increase.Uint64(),
			)
		})
	}
}
