package types_test

import (
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	"testing"
)

func TestGetGasStabilityPoolAddress(t *testing.T) {
	address := types.GetGasStabilityPoolAddress(42)
	require.False(t, address.Empty())

	address2 := types.GetGasStabilityPoolAddress(42)
	require.True(t, address.Equals(address2))

	address3 := types.GetGasStabilityPoolAddress(43)
	require.False(t, address3.Empty())
	require.False(t, address.Equals(address3))
}

func TestGetGasStabilityPoolAddressEVM(t *testing.T) {
	address := types.GetGasStabilityPoolAddressEVM(42)
	require.NotEmpty(t, address)

	address2 := types.GetGasStabilityPoolAddressEVM(42)
	require.True(t, address == address2)

	address3 := types.GetGasStabilityPoolAddressEVM(43)
	require.NotEmpty(t, address3)
	require.False(t, address == address3)
}
