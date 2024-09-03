package types_test

import (
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
	"testing"
)

func TestRevertOptions_GetEVMRevertAddress(t *testing.T) {
	t.Run("valid revert address", func(t *testing.T) {
		addr := sample.EthAddress()
		actualAddr, valid := types.RevertOptions{
			RevertAddress: addr.Hex(),
		}.GetEVMRevertAddress()

		require.True(t, valid)
		require.Equal(t, addr.Hex(), actualAddr.Hex())
	})

	t.Run("invalid revert address", func(t *testing.T) {
		_, valid := types.RevertOptions{
			RevertAddress: "invalid",
		}.GetEVMRevertAddress()

		require.False(t, valid)
	})

	t.Run("empty revert address", func(t *testing.T) {
		_, valid := types.RevertOptions{
			RevertAddress: "",
		}.GetEVMRevertAddress()

		require.False(t, valid)
	})

	t.Run("zero revert address", func(t *testing.T) {
		_, valid := types.RevertOptions{
			RevertAddress: constant.EVMZeroAddress,
		}.GetEVMRevertAddress()

		require.False(t, valid)
	})
}

func TestRevertOptions_GetEVMAbortAddress(t *testing.T) {
	t.Run("valid abort address", func(t *testing.T) {
		addr := sample.EthAddress()
		actualAddr, valid := types.RevertOptions{
			AbortAddress: addr.Hex(),
		}.GetEVMAbortAddress()

		require.True(t, valid)
		require.Equal(t, addr.Hex(), actualAddr.Hex())
	})

	t.Run("invalid abort address", func(t *testing.T) {
		_, valid := types.RevertOptions{
			AbortAddress: "invalid",
		}.GetEVMAbortAddress()

		require.False(t, valid)
	})

	t.Run("empty abort address", func(t *testing.T) {
		_, valid := types.RevertOptions{
			AbortAddress: "",
		}.GetEVMAbortAddress()

		require.False(t, valid)
	})

	t.Run("zero abort address", func(t *testing.T) {
		_, valid := types.RevertOptions{
			AbortAddress: constant.EVMZeroAddress,
		}.GetEVMAbortAddress()

		require.False(t, valid)
	})
}
