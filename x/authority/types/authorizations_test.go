package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/authority/types"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
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

func TestDefaultAuthorizationsList(t *testing.T) {
	t.Run("Default authorizations list", func(t *testing.T) {
		var OperationalPolicyMessageList = []string{
			sdk.MsgTypeURL(&crosschaintypes.MsgRefundAbortedCCTX{}),
			sdk.MsgTypeURL(&crosschaintypes.MsgAbortStuckCCTX{}),
			sdk.MsgTypeURL(&crosschaintypes.MsgUpdateRateLimiterFlags{}),
			sdk.MsgTypeURL(&crosschaintypes.MsgWhitelistERC20{}),
			sdk.MsgTypeURL(&fungibletypes.MsgDeployFungibleCoinZRC20{}),
			sdk.MsgTypeURL(&fungibletypes.MsgDeploySystemContracts{}),
			sdk.MsgTypeURL(&fungibletypes.MsgRemoveForeignCoin{}),
			sdk.MsgTypeURL(&fungibletypes.MsgUpdateZRC20LiquidityCap{}),
			sdk.MsgTypeURL(&fungibletypes.MsgUpdateZRC20WithdrawFee{}),
			sdk.MsgTypeURL(&fungibletypes.MsgUnpauseZRC20{}),
			sdk.MsgTypeURL(&observertypes.MsgAddObserver{}),
			sdk.MsgTypeURL(&observertypes.MsgRemoveChainParams{}),
			sdk.MsgTypeURL(&observertypes.MsgResetChainNonces{}),
			sdk.MsgTypeURL(&observertypes.MsgUpdateChainParams{}),
			//sdk.MsgTypeURL(&observertypes.MsgEnableCCTX{}),
			//sdk.MsgTypeURL(&observertypes.MsgUpdateGasPriceIncreaseFlags{}),
			sdk.MsgTypeURL(&lightclienttypes.MsgEnableHeaderVerification{}),
		}

		// EmergencyPolicyMessageList is a list of messages that can be authorized by the emergency policy
		var EmergencyPolicyMessageList = []string{
			sdk.MsgTypeURL(&crosschaintypes.MsgAddInboundTracker{}),
			sdk.MsgTypeURL(&crosschaintypes.MsgAddOutboundTracker{}),
			sdk.MsgTypeURL(&crosschaintypes.MsgRemoveOutboundTracker{}),
			sdk.MsgTypeURL(&fungibletypes.MsgPauseZRC20{}),
			sdk.MsgTypeURL(&observertypes.MsgUpdateKeygen{}),
			//sdk.MsgTypeURL(&observertypes.MsgDisableCCTX{}),
			sdk.MsgTypeURL(&lightclienttypes.MsgDisableHeaderVerification{}),
		}

		// AdminPolicyMessageList is a list of messages that can be authorized by the admin policy
		var AdminPolicyMessageList = []string{
			sdk.MsgTypeURL(&crosschaintypes.MsgMigrateTssFunds{}),
			sdk.MsgTypeURL(&crosschaintypes.MsgUpdateTssAddress{}),
			sdk.MsgTypeURL(&fungibletypes.MsgUpdateContractBytecode{}),
			sdk.MsgTypeURL(&fungibletypes.MsgUpdateSystemContract{}),
			sdk.MsgTypeURL(&observertypes.MsgUpdateObserver{}),
		}
		defaultList := types.DefaultAuthorizationsList()
		for _, msgUrl := range OperationalPolicyMessageList {
			require.True(t, defaultList.CheckAuthorizationExists(types.Authorization{MsgUrl: msgUrl}))
			policy, err := defaultList.GetAuthorizedPolicy(msgUrl)
			require.NoError(t, err)
			require.Equal(t, types.PolicyType_groupOperational, policy)
		}
		for _, msgUrl := range EmergencyPolicyMessageList {
			require.True(t, defaultList.CheckAuthorizationExists(types.Authorization{MsgUrl: msgUrl}))
			policy, err := defaultList.GetAuthorizedPolicy(msgUrl)
			require.NoError(t, err)
			require.Equal(t, types.PolicyType_groupEmergency, policy)
		}
		for _, msgUrl := range AdminPolicyMessageList {
			require.True(t, defaultList.CheckAuthorizationExists(types.Authorization{MsgUrl: msgUrl}))
			policy, err := defaultList.GetAuthorizedPolicy(msgUrl)
			require.NoError(t, err)
			require.Equal(t, types.PolicyType_groupAdmin, policy)
		}
		require.Len(t, defaultList.Authorizations, len(OperationalPolicyMessageList)+len(EmergencyPolicyMessageList)+len(AdminPolicyMessageList))
	})
}
