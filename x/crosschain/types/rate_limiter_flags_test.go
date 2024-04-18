package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestRateLimiterFlags_Validate(t *testing.T) {
	dec, err := sdk.NewDecFromStr("0.00042")
	require.NoError(t, err)
	duplicatedAddress := sample.EthAddress().String()

	tt := []struct {
		name  string
		flags types.RateLimiterFlags
		isErr bool
	}{
		{
			name: "valid flags",
			flags: types.RateLimiterFlags{
				Enabled: true,
				Window:  42,
				Rate:    sdk.NewUint(42),
				Conversions: []types.Conversion{
					{
						Zrc20: sample.EthAddress().String(),
						Rate:  sdk.NewDec(42),
					},
					{
						Zrc20: sample.EthAddress().String(),
						Rate:  dec,
					},
				},
			},
		},
		{
			name:  "empty is valid",
			flags: types.RateLimiterFlags{},
		},
		{
			name: "duplicated conversion",
			flags: types.RateLimiterFlags{
				Enabled: true,
				Window:  42,
				Rate:    sdk.NewUint(42),
				Conversions: []types.Conversion{
					{
						Zrc20: duplicatedAddress,
						Rate:  sdk.NewDec(42),
					},
					{
						Zrc20: duplicatedAddress,
						Rate:  dec,
					},
				},
			},
			isErr: true,
		},
		{
			name: "invalid conversion rate",
			flags: types.RateLimiterFlags{
				Enabled: true,
				Window:  42,
				Rate:    sdk.NewUint(42),
				Conversions: []types.Conversion{
					{
						Zrc20: sample.EthAddress().String(),
						Rate:  sdk.NewDec(42),
					},
					{
						Zrc20: sample.EthAddress().String(),
					},
				},
			},
			isErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.flags.Validate()
			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}

}

func TestRateLimiterFlags_GetConversion(t *testing.T) {
	dec, err := sdk.NewDecFromStr("0.00042")
	require.NoError(t, err)
	address := sample.EthAddress().String()

	tt := []struct {
		name       string
		flags      types.RateLimiterFlags
		zrc20      string
		expected   sdk.Dec
		shouldFind bool
	}{
		{
			name: "valid conversion",
			flags: types.RateLimiterFlags{
				Enabled: true,
				Window:  42,
				Rate:    sdk.NewUint(42),
				Conversions: []types.Conversion{
					{
						Zrc20: address,
						Rate:  sdk.NewDec(42),
					},
					{
						Zrc20: sample.EthAddress().String(),
						Rate:  dec,
					},
				},
			},
			zrc20:      address,
			expected:   sdk.NewDec(42),
			shouldFind: true,
		},
		{
			name: "not found",
			flags: types.RateLimiterFlags{
				Enabled: true,
				Window:  42,
				Rate:    sdk.NewUint(42),
				Conversions: []types.Conversion{
					{
						Zrc20: sample.EthAddress().String(),
						Rate:  sdk.NewDec(42),
					},
					{
						Zrc20: sample.EthAddress().String(),
						Rate:  dec,
					},
				},
			},
			zrc20:      address,
			expected:   sdk.NewDec(0),
			shouldFind: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual, found := tc.flags.GetConversion(tc.zrc20)
			require.Equal(t, tc.expected, actual)
			require.Equal(t, tc.shouldFind, found)
		})
	}
}
