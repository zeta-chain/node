package chains_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
)

func Test_CalcInboundFastConfirmationAmountCap(t *testing.T) {
	tests := []struct {
		name         string
		chainID      int64
		liquidityCap sdkmath.Uint
		divisor      sdkmath.LegacyDec
		expected     sdkmath.Uint
	}{
		{
			name:         "1000000 / 4000",
			chainID:      1,
			liquidityCap: sdkmath.NewUintFromString("1000000"),
			expected:     sdkmath.NewUint(250),
		},
		{
			name:         "700000 / 4000",
			chainID:      1,
			liquidityCap: sdkmath.NewUintFromString("700000"),
			expected:     sdkmath.NewUint(175),
		},
		{
			name:         "70000 / 4000",
			chainID:      1,
			liquidityCap: sdkmath.NewUintFromString("70000"),
			expected:     sdkmath.NewUint(17), // truncate 17.5 to 17
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := chains.CalcInboundFastConfirmationAmountCap(tt.chainID, tt.liquidityCap)
			require.Equal(t, tt.expected, actual)
		})
	}
}
