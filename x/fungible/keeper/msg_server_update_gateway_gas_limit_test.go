package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/fungible/keeper"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestKeeper_UpdateGatewayGasLimit(t *testing.T) {
	t.Run(
		"can update the gateway gas limit stored in the system contract",
		func(t *testing.T) {
			// ARRANGE
			k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
				UseAuthorityMock: true,
			})

			msgServer := keeper.NewMsgServerImpl(*k)
			k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
			admin := sample.AccAddress()

			authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

			// new gas limit
			newGasLimit := sdkmath.NewInt(2_000_000)

			msg := types.NewMsgUpdateGatewayGasLimit(admin, newGasLimit)
			keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)

			// ACT
			_, err := msgServer.UpdateGatewayGasLimit(ctx, msg)

			// ASSERT
			require.NoError(t, err)
			sc, found := k.GetSystemContract(ctx)
			require.True(t, found)

			// gas limit is updated
			require.EqualValues(t, newGasLimit, sc.GatewayGasLimit)
		},
	)

	t.Run(
		"can update and overwrite the gateway gas limit if system contract state variable not found",
		func(t *testing.T) {
			// ARRANGE
			k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
				UseAuthorityMock: true,
			})

			msgServer := keeper.NewMsgServerImpl(*k)
			k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
			admin := sample.AccAddress()

			authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

			newGasLimit := sdkmath.NewInt(1_500_000)

			_, found := k.GetSystemContract(ctx)
			require.False(t, found)

			msg := types.NewMsgUpdateGatewayGasLimit(admin, newGasLimit)
			keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)

			// ACT
			_, err := msgServer.UpdateGatewayGasLimit(ctx, msg)

			// ASSERT
			require.NoError(t, err)
			sc, found := k.GetSystemContract(ctx)
			require.True(t, found)

			// gas limit is updated
			require.EqualValues(t, newGasLimit, sc.GatewayGasLimit)

			// other fields are not updated
			require.EqualValues(t, "", sc.SystemContract)
			require.EqualValues(t, "", sc.ConnectorZevm)
			require.EqualValues(t, "", sc.Gateway)
		},
	)

	t.Run("should prevent update the gateway gas limit if not admin", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		msg := types.NewMsgUpdateGatewayGasLimit(admin, sdkmath.NewInt(1000000))
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, authoritytypes.ErrUnauthorized)

		// ACT
		_, err := msgServer.UpdateGatewayGasLimit(ctx, msg)

		// ASSERT
		require.Error(t, err)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})
}
