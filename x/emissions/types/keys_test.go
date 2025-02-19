package types_test

import (
	"testing"

	"github.com/zeta-chain/node/x/emissions/types"
)

func Test_SecondsToBlocks(t *testing.T) {
	tt := []struct {
		name           string
		seconds        int64
		expectedBlocks int64
	}{
		{
			name:           "non-fractional result",
			seconds:        6,
			expectedBlocks: 1,
		},
		{
			name:           "default pending ballots buffer",
			seconds:        60 * 60 * 24 * 10,
			expectedBlocks: 144000,
		},
		{
			name:           "fractional result rounded of to lower int",
			seconds:        10,
			expectedBlocks: 1,
		},
		{
			name:           "negative value(not expected , but adding a test case for completeness)",
			seconds:        -10,
			expectedBlocks: -1,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			blocks := types.SecondsToBlocks(tc.seconds)
			if blocks != tc.expectedBlocks {
				t.Fatalf("expected %d, got %d", tc.expectedBlocks, blocks)
			}
		})
	}
}
