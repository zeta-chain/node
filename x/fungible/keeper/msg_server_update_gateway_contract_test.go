package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/fungible/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestKeeper_UpdateGatewayContract(t *testing.T) {
	t.Run("can update the gateway contract address stored in the module", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		systemContractAddr := sample.EthAddress()
		connectorAddr := sample.EthAddress()
		k.SetSystemContract(ctx, types.SystemContract{
			SystemContract: systemContractAddr.Hex(),
			ConnectorZevm:  connectorAddr.Hex(),
			Gateway:        sample.EthAddress().Hex(),
		})

		newGatewayAddr := sample.EthAddress()

		msg := types.NewMsgUpdateGatewayContract(admin, newGatewayAddr.Hex())
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)

		// ACT
		_, err := msgServer.UpdateGatewayContract(ctx, msg)

		// ASSERT
		require.NoError(t, err)
		sc, found := k.GetSystemContract(ctx)
		require.True(t, found)

		// gateway is updated
		require.EqualValues(t, newGatewayAddr.Hex(), sc.Gateway)

		// system contract and connector remain the same
		require.EqualValues(t, systemContractAddr.Hex(), sc.SystemContract)
		require.EqualValues(t, connectorAddr.Hex(), sc.ConnectorZevm)
	})

	t.Run("can update and overwrite the gateway contract if system contract state variable not found", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		newGatewayAddr := sample.EthAddress()

		_, found := k.GetSystemContract(ctx)
		require.False(t, found)

		msg := types.NewMsgUpdateGatewayContract(admin, newGatewayAddr.Hex())
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)

		// ACT
		_, err := msgServer.UpdateGatewayContract(ctx, msg)

		// ASSERT
		require.NoError(t, err)
		sc, found := k.GetSystemContract(ctx)
		require.True(t, found)

		// gateway is updated
		require.EqualValues(t, newGatewayAddr.Hex(), sc.Gateway)

		// system contract and connector are not updated
		require.EqualValues(t, "", sc.SystemContract)
		require.EqualValues(t, "", sc.ConnectorZevm)
	})

	t.Run("should prevent update the gateway contract if not admin", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		msg := types.NewMsgUpdateSystemContract(admin, sample.EthAddress().Hex())
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, authoritytypes.ErrUnauthorized)

		// ACT
		_, err := msgServer.UpdateSystemContract(ctx, msg)

		// ASSERT
		require.Error(t, err)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

}
