package types_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/x/authority/types"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestAuthorizationList_SetAuthorizations(t *testing.T) {
	tt := []struct {
		name             string
		oldList          types.AuthorizationList
		addAuthorization types.Authorization
		expectedList     types.AuthorizationList
	}{
		{
			name: "set new authorization successfully",
			oldList: types.AuthorizationList{Authorizations: []types.Authorization{
				{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
			}},
			addAuthorization: types.Authorization{
				MsgUrl:           "XYZ",
				AuthorizedPolicy: types.PolicyType_groupOperational,
			},
			expectedList: types.AuthorizationList{Authorizations: []types.Authorization{
				{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
				{
					MsgUrl:           "XYZ",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
			}},
		},
		{
			name: "update existing authorization successfully",
			oldList: types.AuthorizationList{Authorizations: []types.Authorization{
				{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
			}},
			addAuthorization: types.Authorization{
				MsgUrl:           "ABC",
				AuthorizedPolicy: types.PolicyType_groupEmergency,
			},
			expectedList: types.AuthorizationList{Authorizations: []types.Authorization{
				{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupEmergency,
				},
			}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.oldList.SetAuthorizations(tc.addAuthorization)
			require.Equal(t, tc.expectedList, tc.oldList)
		})
	}
}

func TestAuthorizationList_GetAuthorizations(t *testing.T) {
	tt := []struct {
		name            string
		authorizations  types.AuthorizationList
		getPolicyMsgUrl string
		expectedPolicy  types.PolicyType
		error           error
	}{
		{
			name: "get authorizations successfully",
			authorizations: types.AuthorizationList{Authorizations: []types.Authorization{
				{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
			}},
			getPolicyMsgUrl: "ABC",
			expectedPolicy:  types.PolicyType_groupOperational,
			error:           nil,
		},
		{
			name:            "get authorizations fails when msg not found in list",
			authorizations:  types.AuthorizationList{Authorizations: []types.Authorization{}},
			getPolicyMsgUrl: "ABC",
			expectedPolicy:  types.PolicyType(0),
			error:           types.ErrAuthorizationNotFound,
		},
		{
			name:            "get authorizations fails when when queried for empty string",
			authorizations:  types.AuthorizationList{Authorizations: []types.Authorization{}},
			getPolicyMsgUrl: "",
			expectedPolicy:  types.PolicyType(0),
			error:           types.ErrAuthorizationNotFound,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			policy, err := tc.authorizations.GetAuthorizedPolicy(tc.getPolicyMsgUrl)
			require.ErrorIs(t, err, tc.error)
			require.Equal(t, tc.expectedPolicy, policy)
		})
	}
}

func TestAuthorizationList_Validate(t *testing.T) {
	tt := []struct {
		name           string
		authorizations types.AuthorizationList
		expectedError  error
	}{
		{
			name: "validate successfully",
			authorizations: types.AuthorizationList{Authorizations: []types.Authorization{
				{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
				{
					MsgUrl:           "XYZ",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
			}},
			expectedError: nil,
		},
		{
			name: "validate failed with duplicate msg url with different policies",
			authorizations: types.AuthorizationList{Authorizations: []types.Authorization{
				{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
				{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupEmergency,
				},
			}},
			expectedError: errors.Wrap(
				types.ErrInvalidAuthorizationList,
				fmt.Sprintf("duplicate message url: %s", "ABC")),
		},
		{
			name: "validate failed with duplicate msg url with same policies",
			authorizations: types.AuthorizationList{Authorizations: []types.Authorization{
				{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
				{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
			}},
			expectedError: errors.Wrap(
				types.ErrInvalidAuthorizationList,
				fmt.Sprintf("duplicate message url: %s", "ABC")),
		}}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.authorizations.Validate()
			require.ErrorIs(t, err, tc.expectedError)
		})
	}
}

func TestAuthorizationList_RemoveAuthorizations(t *testing.T) {
	tt := []struct {
		name                string
		oldList             types.AuthorizationList
		removeAuthorization types.Authorization
		expectedList        types.AuthorizationList
	}{
		{
			name: "remove authorization successfully",
			oldList: types.AuthorizationList{Authorizations: []types.Authorization{
				{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
				{
					MsgUrl:           "XYZ",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
			}},
			removeAuthorization: types.Authorization{
				MsgUrl: "ABC",
			},
			expectedList: types.AuthorizationList{Authorizations: []types.Authorization{
				{
					MsgUrl:           "XYZ",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
			}},
		},
		{
			name: "do not remove anything if authorization not found",
			oldList: types.AuthorizationList{Authorizations: []types.Authorization{
				{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
			}},
			removeAuthorization: types.Authorization{
				MsgUrl: "XYZ",
			},
			expectedList: types.AuthorizationList{Authorizations: []types.Authorization{
				{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
			}},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.oldList.RemoveAuthorizations(tc.removeAuthorization)
			require.Equal(t, tc.expectedList, tc.oldList)
		})
	}
}

func TestDefaultAuthorizationsList(t *testing.T) {
	t.Run("default authorizations list", func(t *testing.T) {
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
			sdk.MsgTypeURL(&observertypes.MsgEnableCCTX{}),
			sdk.MsgTypeURL(&observertypes.MsgUpdateGasPriceIncreaseFlags{}),
			sdk.MsgTypeURL(&lightclienttypes.MsgEnableHeaderVerification{}),
		}

		// EmergencyPolicyMessageList is a list of messages that can be authorized by the emergency policy
		var EmergencyPolicyMessageList = []string{
			sdk.MsgTypeURL(&crosschaintypes.MsgAddInboundTracker{}),
			sdk.MsgTypeURL(&crosschaintypes.MsgAddOutboundTracker{}),
			sdk.MsgTypeURL(&crosschaintypes.MsgRemoveOutboundTracker{}),
			sdk.MsgTypeURL(&fungibletypes.MsgPauseZRC20{}),
			sdk.MsgTypeURL(&observertypes.MsgUpdateKeygen{}),
			sdk.MsgTypeURL(&observertypes.MsgDisableCCTX{}),
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
			_, err := defaultList.GetAuthorizedPolicy(msgUrl)
			require.NoError(t, err)
			policy, err := defaultList.GetAuthorizedPolicy(msgUrl)
			require.NoError(t, err)
			require.Equal(t, types.PolicyType_groupOperational, policy)
		}
		for _, msgUrl := range EmergencyPolicyMessageList {
			_, err := defaultList.GetAuthorizedPolicy(msgUrl)
			require.NoError(t, err)
			policy, err := defaultList.GetAuthorizedPolicy(msgUrl)
			require.NoError(t, err)
			require.Equal(t, types.PolicyType_groupEmergency, policy)
		}
		for _, msgUrl := range AdminPolicyMessageList {
			_, err := defaultList.GetAuthorizedPolicy(msgUrl)
			require.NoError(t, err)
			policy, err := defaultList.GetAuthorizedPolicy(msgUrl)
			require.NoError(t, err)
			require.Equal(t, types.PolicyType_groupAdmin, policy)
		}
		require.Len(
			t,
			defaultList.Authorizations,
			len(OperationalPolicyMessageList)+len(EmergencyPolicyMessageList)+len(AdminPolicyMessageList),
		)
	})
}
