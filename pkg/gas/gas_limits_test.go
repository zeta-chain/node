package gas

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestMultiplyGasPrice(t *testing.T) {
	testCases := []struct {
		name             string
		medianGasPrice   string
		multiplierString string
		expectedGasPrice string
		wantErr          bool
	}{
		{
			name:             "valid multiplication",
			medianGasPrice:   "100",
			multiplierString: "1.5",
			expectedGasPrice: "150", // 100 * 1.5
			wantErr:          false,
		},
		{
			name:             "invalid multiplier format",
			medianGasPrice:   "100",
			multiplierString: "abc",
			expectedGasPrice: "",
			wantErr:          true,
		},
		{
			name:             "zero median price",
			medianGasPrice:   "0",
			multiplierString: "1.5",
			expectedGasPrice: "0",
			wantErr:          false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			medianGasPriceUint := math.NewUintFromString(tc.medianGasPrice)

			result, err := MultiplyGasPrice(medianGasPriceUint, tc.multiplierString)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				expectedGasPriceUint := math.NewUintFromString(tc.expectedGasPrice)
				require.True(t, result.Equal(expectedGasPriceUint))
			}
		})
	}
}
