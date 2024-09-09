package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/x/fungible/types"
)

func TestNewZRC20Data(t *testing.T) {
	zrc20 := types.NewZRC20Data("name", "symbol", 8)
	require.Equal(t, "name", zrc20.Name)
	require.Equal(t, "symbol", zrc20.Symbol)
	require.Equal(t, uint8(8), zrc20.Decimals)
}
