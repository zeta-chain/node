package sui

import (
	"testing"

	"github.com/pattonkan/sui-go/sui"
	"github.com/stretchr/testify/require"
)

func Test_TypeTagFromString(t *testing.T) {
	suiAddr, err := sui.AddressFromHex("0000000000000000000000000000000000000000000000000000000000000002")
	require.NoError(t, err)

	otherAddrStr := "0xae330284eefb31f37777c78007fc3bc0e88ca69b1be267e1b803d37c9ea52dc6"
	otherAddr, err := sui.AddressFromHex(otherAddrStr)
	require.NoError(t, err)

	tests := []struct {
		name     string
		coinType CoinType
		want     sui.StructTag
		errMsg   string
	}{
		{
			name:     "SUI coin type",
			coinType: SUI,
			want: sui.StructTag{
				Address: suiAddr,
				Module:  "sui",
				Name:    "SUI",
			},
		},
		{
			name:     "some other coin type",
			coinType: CoinType(otherAddrStr + "::other::TOKEN"),
			want: sui.StructTag{
				Address: otherAddr,
				Module:  "other",
				Name:    "TOKEN",
			},
		},
		{
			name:     "invalid type string",
			coinType: CoinType(otherAddrStr),
			want:     sui.StructTag{},
			errMsg:   "invalid type string",
		},
		{
			name:     "invalid address",
			coinType: "invalid::sui::SUI",
			want:     sui.StructTag{},
			errMsg:   "invalid address",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := TypeTagFromString(string(test.coinType))
			if test.errMsg != "" {
				require.Empty(t, got)
				require.ErrorContains(t, err, test.errMsg)
				return
			}

			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}
