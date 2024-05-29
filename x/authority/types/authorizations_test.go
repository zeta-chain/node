package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/authority/types"
)

func TestAuthorizationList_SetAuthorizations(t *testing.T) {
	t.Run("Set new authorization successfully", func(t *testing.T) {
		authorizationsList := types.DefaultAuthorizationsList()
		newAuthorization := sample.Authorization()
		require.False(t, authorizationsList.CheckAuthorizationExists(newAuthorization))
		authorizationsList.SetAuthorizations(newAuthorization)
		require.Len(t, authorizationsList.Authorizations, len(types.DefaultAuthorizationsList().Authorizations)+1)
		require.True(t, authorizationsList.CheckAuthorizationExists(newAuthorization))
	})

	t.Run("Update existing authorization successfully", func(t *testing.T) {
		authorizationsList := types.DefaultAuthorizationsList()
		newAuthorization := sample.Authorization()
		require.False(t, authorizationsList.CheckAuthorizationExists(newAuthorization))
		authorizationsList.SetAuthorizations(newAuthorization)
		require.Len(t, authorizationsList.Authorizations, len(types.DefaultAuthorizationsList().Authorizations)+1)
		require.True(t, authorizationsList.CheckAuthorizationExists(newAuthorization))

		newAuthorization.AuthorizedPolicy = types.PolicyType_groupEmergency
		authorizationsList.SetAuthorizations(newAuthorization)
		require.True(t, authorizationsList.CheckAuthorizationExists(newAuthorization))
		policy, err := authorizationsList.GetAuthorizedPolicy(newAuthorization.MsgUrl)
		require.NoError(t, err)
		require.Equal(t, newAuthorization.AuthorizedPolicy, policy)
	})
}

func TestAuthorizationList_GetAuthorizedPolicy(t *testing.T) {
	t.Run("Get authorized policy successfully", func(t *testing.T) {
		authorizationsList := types.DefaultAuthorizationsList()
		newAuthorization := sample.Authorization()
		authorizationsList.SetAuthorizations(newAuthorization)
		policy, err := authorizationsList.GetAuthorizedPolicy(newAuthorization.MsgUrl)
		require.NoError(t, err)
		require.Equal(t, newAuthorization.AuthorizedPolicy, policy)
	})

	t.Run("Get authorized policy failed with not found", func(t *testing.T) {
		authorizationsList := types.DefaultAuthorizationsList()
		policy, err := authorizationsList.GetAuthorizedPolicy("ABC")
		require.ErrorIs(t, err, types.ErrAuthorizationNotFound)
		require.Equal(t, types.PolicyType(0), policy)
	})
}

func TestAuthorizationList_CheckAuthorizationExists(t *testing.T) {
	t.Run("Check authorization exists successfully", func(t *testing.T) {
		authorizationsList := types.DefaultAuthorizationsList()
		newAuthorization := sample.Authorization()
		require.False(t, authorizationsList.CheckAuthorizationExists(newAuthorization))
		authorizationsList.SetAuthorizations(newAuthorization)
		require.True(t, authorizationsList.CheckAuthorizationExists(newAuthorization))
	})
}

func TestAuthorizationList_Validate(t *testing.T) {
	t.Run("Validate successfully", func(t *testing.T) {
		authorizationsList := types.DefaultAuthorizationsList()
		require.NoError(t, authorizationsList.Validate())
	})
	t.Run("Validate failed with duplicate msg url with different policies", func(t *testing.T) {
		authorizationsList := types.AuthorizationList{Authorizations: []types.Authorization{
			{
				MsgUrl:           "ABC",
				AuthorizedPolicy: types.PolicyType_groupOperational,
			},
			{
				MsgUrl:           "ABC",
				AuthorizedPolicy: types.PolicyType_groupEmergency,
			},
		}}

		require.ErrorIs(t, authorizationsList.Validate(), types.ErrInValidAuthorizationList)
	})

	t.Run("Validate failed with duplicate msg url with same policies", func(t *testing.T) {
		authorizationsList := types.AuthorizationList{Authorizations: []types.Authorization{
			{
				MsgUrl:           "ABC",
				AuthorizedPolicy: types.PolicyType_groupOperational,
			},
			{
				MsgUrl:           "ABC",
				AuthorizedPolicy: types.PolicyType_groupOperational,
			},
		}}

		require.ErrorIs(t, authorizationsList.Validate(), types.ErrInValidAuthorizationList)
	})

}

func TestAuthorizationList_RemoveAuthorizations(t *testing.T) {
	t.Run("Remove authorization successfully", func(t *testing.T) {
		authorizationsList := types.DefaultAuthorizationsList()
		newAuthorization := sample.Authorization()
		authorizationsList.SetAuthorizations(newAuthorization)
		require.True(t, authorizationsList.CheckAuthorizationExists(newAuthorization))
		authorizationsList.RemoveAuthorizations(newAuthorization)
		require.False(t, authorizationsList.CheckAuthorizationExists(newAuthorization))
		require.Len(t, authorizationsList.Authorizations, len(types.DefaultAuthorizationsList().Authorizations))
	})

	t.Run("do not remove anything if authorization not found", func(t *testing.T) {
		authorizationsList := types.DefaultAuthorizationsList()
		authorizationsList.RemoveAuthorizations(sample.Authorization())
		require.ElementsMatch(t, authorizationsList.Authorizations, types.DefaultAuthorizationsList().Authorizations)
	})
}
