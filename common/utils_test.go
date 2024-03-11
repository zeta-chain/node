package common_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
)

func Test_GasPriceMultiplier(t *testing.T) {
	tt := []struct {
		name       string
		chainID    int64
		multiplier float64
	}{
		{
			name:       "get Ethereum multiplier",
			chainID:    1,
			multiplier: 1.2,
		},
		{
			name:       "get Goerli multiplier",
			chainID:    5,
			multiplier: 1.2,
		},
		{
			name:       "get BSC multiplier",
			chainID:    56,
			multiplier: 1.2,
		},
		{
			name:       "get BSC Testnet multiplier",
			chainID:    97,
			multiplier: 1.2,
		},
		{
			name:       "get Polygon multiplier",
			chainID:    137,
			multiplier: 1.2,
		},
		{
			name:       "get Mumbai Testnet multiplier",
			chainID:    80001,
			multiplier: 1.2,
		},
		{
			name:       "get Bitcoin multiplier",
			chainID:    8332,
			multiplier: 2.0,
		},
		{
			name:       "get Bitcoin Testnet multiplier",
			chainID:    18332,
			multiplier: 2.0,
		},
		{
			name:       "get unknown chain gas price multiplier",
			chainID:    1234,
			multiplier: 1.0,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			multiplier := common.GasPriceMultiplier(tc.chainID)
			require.Equal(t, tc.multiplier, multiplier)
		})
	}

}
