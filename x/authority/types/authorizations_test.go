package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/authority/types"
)

func TestAuthorizationList_AddAuthorizations(t *testing.T) {
	t.Run("AddAuthorizations", func(t *testing.T) {
		authorizationsList := types.DefaultAuthorizationsList()
		newList := sample.AuthorizationList("sample")
		authorizationsList.AddAuthorizations(newList)
		require.ElementsMatch(
			t,
			append(types.DefaultAuthorizationsList().Authorizations, newList.Authorizations...),
			authorizationsList.Authorizations,
		)
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
