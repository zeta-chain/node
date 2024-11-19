package types

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func Test_ZRC20ToCosmosDenom(t *testing.T) {
	address := big.NewInt(12345) // 0x3039
	expected := "zrc20/0x0000000000000000000000000000000000003039"
	denom := ZRC20ToCosmosDenom(common.BigToAddress(address))
	require.Equal(t, expected, denom, "denom should be %s, got %s", expected, denom)
}

func Test_createCoinSet(t *testing.T) {
	tokenAddr := common.HexToAddress("0x0000000000000000000000000000000000003039")
	tokenDenom := ZRC20ToCosmosDenom(tokenAddr)
	amount := big.NewInt(100)

	coinSet, err := CreateZRC20CoinSet(tokenAddr, amount)
	require.NoError(t, err, "createCoinSet should not return an error")
	require.NotNil(t, coinSet, "coinSet should not be nil")

	coin := coinSet[0]
	require.Equal(t, tokenDenom, coin.Denom, "coin denom should be %s, got %s", tokenDenom, coin.Denom)
	require.Equal(t, amount, coin.Amount.BigInt(), "coin amount should be %s, got %s", amount, coin.Amount.BigInt())
}

func Test_CoinIsZRC20(t *testing.T) {
	test := []struct {
		denom    string
		expected bool
	}{
		{"", false}, // Empty string.
		{"zrc20/", false}, // Missing address.
		{"zrc20/0x514910771af9ca656af840dff83e8264ecf986ca", true}, // Valid ZRC20 address.
		{"zrc20/0xCa14007Eff0dB1f8135f4C25B34De49AB0d42766", true}, // Valid ZRC20 address.
		{"zrc200xabcdef", false}, // Malformed prefix.
		{"foo/0x0123456789", false}, // Invalid prefix.
		{"ZRC20/0x0123456789abcdef", false}, // Invalid prefix.
	}

	for _, tt := range test {
		t.Run(tt.denom, func(t *testing.T) {
			result := CoinIsZRC20(tt.denom)
			require.Equal(t, tt.expected, result, "got %v, want %v", result, tt.expected)
		})
	}
}
