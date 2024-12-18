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

func TestMsgServer_UpdateOperationalFlags(t *testing.T) {
	t.Run("can update operational flags", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)

		// set admin
		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		// test initial set
		restartHeight := int64(100)
		msg := types.MsgUpdateOperationalFlags{
			Creator: admin,
			OperationalFlags: types.OperationalFlags{
				RestartHeight: restartHeight,
			},
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := srv.UpdateOperationalFlags(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)

		operationalFlags, found := k.GetOperationalFlags(ctx)
		require.True(t, found)
		require.Equal(t, restartHeight, operationalFlags.RestartHeight)

		// verify that we can set it again
		restartHeight = 101
		msg = types.MsgUpdateOperationalFlags{
			Creator: admin,
			OperationalFlags: types.OperationalFlags{
				RestartHeight: restartHeight,
			},
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err = srv.UpdateOperationalFlags(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)

		operationalFlags, found = k.GetOperationalFlags(ctx)
		require.True(t, found)
		require.Equal(t, restartHeight, operationalFlags.RestartHeight)
	})

	t.Run("cannot update operational flags if not authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		msg := types.MsgUpdateOperationalFlags{
			Creator:          admin,
			OperationalFlags: sample.OperationalFlags(),
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := srv.UpdateOperationalFlags(sdk.WrapSDKContext(ctx), &msg)

		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})
}
