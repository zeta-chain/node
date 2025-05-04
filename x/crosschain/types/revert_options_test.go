package types_test

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
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

func TestRevertOptions_GetSOLRevertAddress(t *testing.T) {
	t.Run("valid revert address", func(t *testing.T) {
		addr := sample.SolanaAddress(t)
		actualAddr, valid := types.RevertOptions{
			RevertAddress: addr,
		}.GetSOLRevertAddress()

		require.True(t, valid)
		require.Equal(t, addr, actualAddr.String())
	})

	t.Run("invalid revert address", func(t *testing.T) {
		_, valid := types.RevertOptions{
			RevertAddress: "invalid",
		}.GetSOLRevertAddress()

		require.False(t, valid)
	})

	t.Run("empty revert address", func(t *testing.T) {
		_, valid := types.RevertOptions{
			RevertAddress: "",
		}.GetSOLRevertAddress()

		require.False(t, valid)
	})
}

func TestRevertOptions_GetBTCRevertAddress(t *testing.T) {
	t.Run("valid Bitcoin revert address", func(t *testing.T) {
		r := sample.Rand()
		addr := sample.BTCAddressP2WPKH(t, r, &chaincfg.TestNet3Params).String()
		actualAddr, valid := types.RevertOptions{
			RevertAddress: addr,
		}.GetBTCRevertAddress(chains.BitcoinTestnet.ChainId)

		require.True(t, valid)
		require.Equal(t, addr, actualAddr)
	})

	t.Run("invalid Bitcoin revert address", func(t *testing.T) {
		actualAddr, valid := types.RevertOptions{
			// it's a regnet address, not testnet
			RevertAddress: "bcrt1qy9pqmk2pd9sv63g27jt8r657wy0d9uee4x2dt2",
		}.GetBTCRevertAddress(chains.BitcoinTestnet.ChainId)

		require.False(t, valid)
		require.Empty(t, actualAddr)
	})

	t.Run("empty revert address", func(t *testing.T) {
		actualAddr, valid := types.RevertOptions{
			RevertAddress: "",
		}.GetBTCRevertAddress(chains.BitcoinTestnet.ChainId)

		require.False(t, valid)
		require.Empty(t, actualAddr)
	})

	t.Run("unsupported Bitcoin revert address", func(t *testing.T) {
		actualAddr, valid := types.RevertOptions{
			// address not supported
			RevertAddress: "035e4ae279bd416b5da724972c9061ec6298dac020d1e3ca3f06eae715135cdbec",
		}.GetBTCRevertAddress(chains.BitcoinTestnet.ChainId)

		require.False(t, valid)
		require.Empty(t, actualAddr)
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
