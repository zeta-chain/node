package math

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CantorPair_Unpair(t *testing.T) {
	tests := []struct {
		name string
		x    uint32
		y    uint32
		z    uint64
	}{
		{name: "test 1", x: 47, y: 32, z: 3192},
		{name: "test 2", x: 32767, y: 32767, z: 2147418112},
		{name: "test 3", x: 512628174, y: 6154648, z: 134567808466687901},
		{name: "test 4", x: 925478314, y: 91456237, z: 517077941108709313},
		{name: "test 5", x: MaxPairValue, y: MaxPairValue, z: 9223372032559808512},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// test CantorPair
			actualZ := CantorPair(tt.x, tt.y)
			require.Equal(t, tt.z, actualZ)

			// test CantorUnpair
			actualX, actualY := CantorUnpair(tt.z)
			require.Equal(t, tt.x, actualX)
			require.Equal(t, tt.y, actualY)
		})
	}
}
