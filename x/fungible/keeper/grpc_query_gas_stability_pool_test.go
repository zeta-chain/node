package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestKeeper_GasStabilityPoolAddress(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		res, err := k.GasStabilityPoolAddress(ctx, nil)
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return if req is not nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		res, err := k.GasStabilityPoolAddress(ctx, &types.QueryGetGasStabilityPoolAddress{})
		require.NoError(t, err)
		require.Equal(t, &types.QueryGetGasStabilityPoolAddressResponse{
			CosmosAddress: types.GasStabilityPoolAddress().String(),
			EvmAddress:    types.GasStabilityPoolAddressEVM().String(),
		}, res)
	})
}

func TestKeeper_GasStabilityPoolBalance(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		res, err := k.GasStabilityPoolBalance(ctx, nil)
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should error if system contracts not deployed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		chainID := 5

		res, err := k.GasStabilityPoolBalance(ctx, &types.QueryGetGasStabilityPoolBalance{
			ChainId: int64(chainID),
		})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return balance", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		chainID := 5
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		setupGasCoin(t, ctx, k, sdkk.EvmKeeper, int64(chainID), "foobar", "foobar")

		res, err := k.GasStabilityPoolBalance(ctx, &types.QueryGetGasStabilityPoolBalance{
			ChainId: int64(chainID),
		})
		require.NoError(t, err)
		require.Equal(t, "0", res.Balance)
	})

}

func TestKeeper_GasStabilityPoolBalanceAll(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		res, err := k.GasStabilityPoolBalanceAll(ctx, nil)
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return empty balances if no supported chains", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseObserverMock: true,
		})
		observerMock := keepertest.GetFungibleObserverMock(t, k)
		observerMock.On("GetSupportedChains", mock.Anything).Return([]chains.Chain{})

		res, err := k.GasStabilityPoolBalanceAll(ctx, &types.QueryAllGasStabilityPoolBalance{})
		require.NoError(t, err)
		require.Empty(t, res.Balances)
	})

	t.Run("should error if system contracts not deployed", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseObserverMock: true,
		})
		observerMock := keepertest.GetFungibleObserverMock(t, k)
		chainID := 5
		observerMock.On("GetSupportedChains", mock.Anything).Return([]chains.Chain{
			{
				ChainId: int64(chainID),
			},
		})

		res, err := k.GasStabilityPoolBalanceAll(ctx, &types.QueryAllGasStabilityPoolBalance{})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return balances", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseObserverMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		observerMock := keepertest.GetFungibleObserverMock(t, k)
		chainID := 5
		observerMock.On("GetSupportedChains", mock.Anything).Return([]chains.Chain{
			{
				ChainId: int64(chainID),
			},
		})

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		setupGasCoin(t, ctx, k, sdkk.EvmKeeper, int64(chainID), "foobar", "foobar")

		res, err := k.GasStabilityPoolBalanceAll(ctx, &types.QueryAllGasStabilityPoolBalance{})
		require.NoError(t, err)
		require.Len(t, res.Balances, 1)
		require.Equal(t, int64(chainID), res.Balances[0].ChainId)
		require.Equal(t, "0", res.Balances[0].Balance)
	})

	t.Run("should ignore ZetaChain chain ID in response", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseObserverMock: true,
		})
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		observerMock := keepertest.GetFungibleObserverMock(t, k)
		chainID := 5
		observerMock.On("GetSupportedChains", mock.Anything).Return([]chains.Chain{
			{
				ChainId: int64(chainID),
			},
			{
				ChainId: chains.ZetaChainMainnet.ChainId,
				Network: chains.Network_zeta,
			},
		})

		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		setupGasCoin(t, ctx, k, sdkk.EvmKeeper, int64(chainID), "foobar", "foobar")

		res, err := k.GasStabilityPoolBalanceAll(ctx, &types.QueryAllGasStabilityPoolBalance{})
		require.NoError(t, err)
		require.Len(t, res.Balances, 1)
		require.Equal(t, int64(chainID), res.Balances[0].ChainId)
		require.Equal(t, "0", res.Balances[0].Balance)
	})
}
