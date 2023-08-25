package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestGetGasStabilityPoolAddress(t *testing.T) {
	address := types.GasStabilityPoolAddress(42)
	require.False(t, address.Empty())

	address2 := types.GasStabilityPoolAddress(42)
	require.True(t, address.Equals(address2))

	address3 := types.GasStabilityPoolAddress(43)
	require.False(t, address3.Empty())
	require.False(t, address.Equals(address3))
}

func TestGetGasStabilityPoolAddressEVM(t *testing.T) {
	address := types.GasStabilityPoolAddressEVM(42)
	require.NotEmpty(t, address)

	address2 := types.GasStabilityPoolAddressEVM(42)
	require.True(t, address == address2)

	address3 := types.GasStabilityPoolAddressEVM(43)
	require.NotEmpty(t, address3)
	require.False(t, address == address3)
}
