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

	coinSet, err := CreateCoinSet(tokenAddr, amount)
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
		{"zrc20/0x0123456789abcdef", true},
		{"zrc20/0xabcdef0123456789", true},
		{"zrc200xabcdef", false},
		{"foo/0x0123456789", false},
	}

	for _, tt := range test {
		t.Run(tt.denom, func(t *testing.T) {
			result := CoinIsZRC20(tt.denom)
			if result != tt.expected {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}
