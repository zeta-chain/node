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
	tokenDenom := "zrc20/0x0000000000000000000000000000000000003039"
	amount := big.NewInt(100)

	coinSet, err := CreateCoinSet(tokenDenom, amount)
	require.NoError(t, err, "createCoinSet should not return an error")
	require.NotNil(t, coinSet, "coinSet should not be nil")

	coin := coinSet[0]
	require.Equal(t, tokenDenom, coin.Denom, "coin denom should be %s, got %s", tokenDenom, coin.Denom)
	require.Equal(t, amount, coin.Amount.BigInt(), "coin amount should be %s, got %s", amount, coin.Amount.BigInt())
}
