package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/observer/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestMsgServer_UpdateV2ZetaFlows(t *testing.T) {
	t.Run("can enable V2 ZETA flows", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		msg := types.MsgUpdateV2ZetaFlows{
			Creator:         admin,
			IsV2ZetaEnabled: true,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// ACT
		_, err := srv.UpdateV2ZetaFlows(sdk.WrapSDKContext(ctx), &msg)

		// ASSERT
		require.NoError(t, err)

		flags, found := k.GetCrosschainFlags(ctx)
		require.True(t, found)
		require.True(t, flags.IsV2ZetaEnabled)
	})

	t.Run("can disable V2 ZETA flows", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)

		// set initial state with V2 ZETA flows enabled
		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled:  true,
			IsOutboundEnabled: true,
			IsV2ZetaEnabled:   true,
		})

		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		msg := types.MsgUpdateV2ZetaFlows{
			Creator:         admin,
			IsV2ZetaEnabled: false,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// ACT
		_, err := srv.UpdateV2ZetaFlows(sdk.WrapSDKContext(ctx), &msg)

		// ASSERT
		require.NoError(t, err)

		flags, found := k.GetCrosschainFlags(ctx)
		require.True(t, found)
		require.False(t, flags.IsV2ZetaEnabled)
		// verify other flags are preserved
		require.True(t, flags.IsInboundEnabled)
		require.True(t, flags.IsOutboundEnabled)
	})

	t.Run("sets default flags if not found", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		msg := types.MsgUpdateV2ZetaFlows{
			Creator:         admin,
			IsV2ZetaEnabled: true,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// ACT
		_, err := srv.UpdateV2ZetaFlows(sdk.WrapSDKContext(ctx), &msg)

		// ASSERT
		require.NoError(t, err)

		flags, found := k.GetCrosschainFlags(ctx)
		require.True(t, found)
		require.True(t, flags.IsV2ZetaEnabled)
		// verify default flags are applied
		defaultFlags := types.DefaultCrosschainFlags()
		require.Equal(t, defaultFlags.IsInboundEnabled, flags.IsInboundEnabled)
		require.Equal(t, defaultFlags.IsOutboundEnabled, flags.IsOutboundEnabled)
	})

	t.Run("cannot update V2 ZETA flows if not authorized", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		msg := types.MsgUpdateV2ZetaFlows{
			Creator:         admin,
			IsV2ZetaEnabled: true,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)

		// ACT
		_, err := srv.UpdateV2ZetaFlows(sdk.WrapSDKContext(ctx), &msg)

		// ASSERT
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})
}
