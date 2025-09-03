package v5_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	v5 "github.com/zeta-chain/node/x/fungible/migrations/v5"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestMigrateStore(t *testing.T) {
	t.Run("set value if system contract not found", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.FungibleKeeper(t)

		// Act
		err := v5.MigrateStore(ctx, k)

		// Assert
		require.NoError(t, err)
		system, found := k.GetSystemContract(ctx)
		require.True(t, found)
		require.Equal(t, sdkmath.NewIntFromBigInt(types.GatewayGasLimit), system.GatewayGasLimit)
	})

	t.Run("set value if system contract is found", func(t *testing.T) {
		// Arrange
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		defaultSystemContract := *sample.SystemContract()
		defaultSystemContract.GatewayGasLimit = sdkmath.NewInt(999999)
		k.SetSystemContract(ctx, defaultSystemContract)

		// Act
		err := v5.MigrateStore(ctx, k)

		// Assert
		require.NoError(t, err)
		system, found := k.GetSystemContract(ctx)
		require.True(t, found)
		require.Equal(t, sdkmath.NewIntFromBigInt(types.GatewayGasLimit), system.GatewayGasLimit)
		require.Equal(t, defaultSystemContract.SystemContract, system.SystemContract)
		require.Equal(t, defaultSystemContract.ConnectorZevm, system.ConnectorZevm)
		require.Equal(t, defaultSystemContract.Gateway, system.Gateway)
	})
}
