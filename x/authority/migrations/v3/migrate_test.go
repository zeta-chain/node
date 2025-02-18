package v3_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	v3 "github.com/zeta-chain/node/x/authority/migrations/v3"
	"github.com/zeta-chain/node/x/authority/types"
)

func TestMigrateStore(t *testing.T) {
	var (
		updateZRC20NameAuthorization = types.Authorization{
			MsgUrl:           "/zetachain.zetacore.fungible.MsgUpdateZRC20Name",
			AuthorizedPolicy: types.PolicyType_groupAdmin,
		}
		removeInboundAuthorization = types.Authorization{
			MsgUrl:           "/zetachain.zetacore.crosschain.MsgRemoveInboundTracker",
			AuthorizedPolicy: types.PolicyType_groupEmergency,
		}
		updateOperationalChainParamsAuthorization = types.Authorization{
			MsgUrl:           "/zetachain.zetacore.observer.MsgUpdateOperationalChainParams",
			AuthorizedPolicy: types.PolicyType_groupOperational,
		}
		updateChainParamsAuthorization = types.Authorization{
			MsgUrl:           "/zetachain.zetacore.observer.MsgUpdateChainParams",
			AuthorizedPolicy: types.PolicyType_groupAdmin,
		}
		disableFastConfirmationAuthorization = types.Authorization{
			MsgUrl:           "/zetachain.zetacore.observer.MsgDisableFastConfirmation",
			AuthorizedPolicy: types.PolicyType_groupEmergency,
		}
	)

	t.Run("update authorization list", func(t *testing.T) {
		// Arrange
		k, ctx := keepertest.AuthorityKeeper(t)

		list := types.DefaultAuthorizationsList()
		list.RemoveAuthorization("/zetachain.zetacore.fungible.MsgUpdateZRC20Name")
		list.RemoveAuthorization("/zetachain.zetacore.crosschain.MsgRemoveInboundTracker")
		list.RemoveAuthorization("/zetachain.zetacore.observer.MsgUpdateOperationalChainParams")
		list.RemoveAuthorization("/zetachain.zetacore.observer.MsgUpdateChainParams")
		list.RemoveAuthorization("/zetachain.zetacore.observer.MsgDisableFastConfirmation")
		k.SetAuthorizationList(ctx, list)

		// Act
		err := v3.MigrateStore(ctx, *k)

		// Assert
		require.NoError(t, err)
		newList, found := k.GetAuthorizationList(ctx)
		require.True(t, found)

		// two lists should be equal if adds the removed authorizations back
		list.SetAuthorization(updateZRC20NameAuthorization)
		list.SetAuthorization(removeInboundAuthorization)
		list.SetAuthorization(updateOperationalChainParamsAuthorization)
		list.SetAuthorization(updateChainParamsAuthorization)
		list.SetAuthorization(disableFastConfirmationAuthorization)
		require.Equal(t, list, newList)
	})

	t.Run("set default authorization list if list is not found", func(t *testing.T) {
		// Arrange
		k, ctx := keepertest.AuthorityKeeper(t)

		// Act
		err := v3.MigrateStore(ctx, *k)

		// Assert
		require.NoError(t, err)
		list, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		require.Equal(t, types.DefaultAuthorizationsList(), list)
	})

	t.Run("return error if authorization list is invalid", func(t *testing.T) {
		// Arrange
		k, ctx := keepertest.AuthorityKeeper(t)

		k.SetAuthorizationList(ctx, types.AuthorizationList{Authorizations: []types.Authorization{
			{
				MsgUrl:           "ABC",
				AuthorizedPolicy: types.PolicyType_groupEmergency,
			},
			{
				MsgUrl:           "ABC",
				AuthorizedPolicy: types.PolicyType_groupEmergency,
			},
		}})

		// Act
		err := v3.MigrateStore(ctx, *k)

		// Assert
		require.Error(t, err)
	})
}
