package types_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/x/authority/types"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	lightclienttypes "github.com/zeta-chain/node/x/lightclient/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
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
			name:    "set new authorization successfully with empty list",
			oldList: types.AuthorizationList{Authorizations: []types.Authorization{}},
			addAuthorization: types.Authorization{
				MsgUrl:           "XYZ",
				AuthorizedPolicy: types.PolicyType_groupOperational,
			},
			expectedList: types.AuthorizationList{Authorizations: []types.Authorization{
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
		{
			name: "update existing authorization successfully in the middle of the list",
			oldList: types.AuthorizationList{Authorizations: []types.Authorization{
				{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
				{
					MsgUrl:           "XYZ",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
				{
					MsgUrl:           "DEF",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
			}},
			addAuthorization: types.Authorization{
				MsgUrl:           "XYZ",
				AuthorizedPolicy: types.PolicyType_groupEmergency,
			},
			expectedList: types.AuthorizationList{Authorizations: []types.Authorization{
				{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
				{
					MsgUrl:           "XYZ",
					AuthorizedPolicy: types.PolicyType_groupEmergency,
				},
				{
					MsgUrl:           "DEF",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
			}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.oldList.SetAuthorization(tc.addAuthorization)
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
			name: "get authorizations successfully for admin policy",
			authorizations: types.AuthorizationList{Authorizations: []types.Authorization{
				{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupAdmin,
				},
			}},
			getPolicyMsgUrl: "ABC",
			expectedPolicy:  types.PolicyType_groupAdmin,
			error:           nil,
		},
		{
			name: "get authorizations successfully for emergency policy",
			authorizations: types.AuthorizationList{Authorizations: []types.Authorization{
				{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupEmergency,
				},
			}},
			getPolicyMsgUrl: "ABC",
			expectedPolicy:  types.PolicyType_groupEmergency,
			error:           nil,
		},
		{
			name: "get authorizations fails when msg not found in list",
			authorizations: types.AuthorizationList{Authorizations: []types.Authorization{
				{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
			}},
			getPolicyMsgUrl: "XYZ",
			expectedPolicy:  types.PolicyType_groupEmpty,
			error:           types.ErrAuthorizationNotFound,
		},
		{
			name:            "get authorizations fails when msg not found in list",
			authorizations:  types.AuthorizationList{Authorizations: []types.Authorization{}},
			getPolicyMsgUrl: "ABC",
			expectedPolicy:  types.PolicyType_groupEmpty,
			error:           types.ErrAuthorizationNotFound,
		},
		{
			name:            "get authorizations fails when when queried for empty string",
			authorizations:  types.AuthorizationList{Authorizations: []types.Authorization{}},
			getPolicyMsgUrl: "",
			expectedPolicy:  types.PolicyType_groupEmpty,
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
			name:           "validate default authorizations list",
			authorizations: types.DefaultAuthorizationsList(),
			expectedError:  nil,
		},
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
			name:           "validate successfully with empty list",
			authorizations: types.AuthorizationList{Authorizations: []types.Authorization{}},
			expectedError:  nil,
		},
		{
			name:           "validate successfully for default list",
			authorizations: types.DefaultAuthorizationsList(),
			expectedError:  nil,
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
		name         string
		oldList      types.AuthorizationList
		removeMsgUrl string
		expectedList types.AuthorizationList
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
			removeMsgUrl: "ABC",
			expectedList: types.AuthorizationList{Authorizations: []types.Authorization{
				{
					MsgUrl:           "XYZ",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
			}},
		},
		{
			name: "remove authorization successfully in the middle of the list",
			oldList: types.AuthorizationList{Authorizations: []types.Authorization{
				{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
				{
					MsgUrl:           "XYZ",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
				{
					MsgUrl:           "DEF",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
			}},
			removeMsgUrl: "XYZ",
			expectedList: types.AuthorizationList{Authorizations: []types.Authorization{
				{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
				{
					MsgUrl:           "DEF",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
			}},
		},
		{
			name:         "do not remove anything when trying to remove from an empty list",
			oldList:      types.AuthorizationList{Authorizations: []types.Authorization{}},
			removeMsgUrl: "XYZ",
			expectedList: types.AuthorizationList{Authorizations: []types.Authorization{}},
		},
		{
			name: "do not remove anything if authorization not found",
			oldList: types.AuthorizationList{Authorizations: []types.Authorization{
				{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
			}},
			removeMsgUrl: "XYZ",
			expectedList: types.AuthorizationList{Authorizations: []types.Authorization{
				{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
			}},
		},
		// The list is invalid, but this test case tries to assert the expected functionality
		{
			name: "return after removing first occurrence",
			oldList: types.AuthorizationList{Authorizations: []types.Authorization{
				{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
				{
					MsgUrl:           "ABC",
					AuthorizedPolicy: types.PolicyType_groupOperational,
				},
			}},
			removeMsgUrl: "ABC",
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
			tc.oldList.RemoveAuthorization(tc.removeMsgUrl)
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
			sdk.MsgTypeURL(&fungibletypes.MsgDeploySystemContracts{}),
			sdk.MsgTypeURL(&fungibletypes.MsgUpdateZRC20LiquidityCap{}),
			sdk.MsgTypeURL(&fungibletypes.MsgUpdateZRC20WithdrawFee{}),
			sdk.MsgTypeURL(&fungibletypes.MsgUnpauseZRC20{}),
			sdk.MsgTypeURL(&fungibletypes.MsgUpdateGatewayGasLimit{}),
			sdk.MsgTypeURL(&observertypes.MsgResetChainNonces{}),
			sdk.MsgTypeURL(&observertypes.MsgEnableCCTX{}),
			sdk.MsgTypeURL(&observertypes.MsgUpdateGasPriceIncreaseFlags{}),
			sdk.MsgTypeURL(&observertypes.MsgUpdateOperationalFlags{}),
			sdk.MsgTypeURL(&observertypes.MsgUpdateOperationalChainParams{}),
		}

		// EmergencyPolicyMessageList is a list of messages that can be authorized by the emergency policy
		var EmergencyPolicyMessageList = []string{
			sdk.MsgTypeURL(&crosschaintypes.MsgAddInboundTracker{}),
			sdk.MsgTypeURL(&crosschaintypes.MsgRemoveInboundTracker{}),
			sdk.MsgTypeURL(&crosschaintypes.MsgAddOutboundTracker{}),
			sdk.MsgTypeURL(&crosschaintypes.MsgRemoveOutboundTracker{}),
			sdk.MsgTypeURL(&fungibletypes.MsgPauseZRC20{}),
			sdk.MsgTypeURL(&observertypes.MsgUpdateKeygen{}),
			sdk.MsgTypeURL(&observertypes.MsgDisableCCTX{}),
			sdk.MsgTypeURL(&observertypes.MsgDisableFastConfirmation{}),
			sdk.MsgTypeURL(&lightclienttypes.MsgDisableHeaderVerification{}),
		}

		// AdminPolicyMessageList is a list of messages that can be authorized by the admin policy
		var AdminPolicyMessageList = []string{
			sdk.MsgTypeURL(&crosschaintypes.MsgMigrateTssFunds{}),
			sdk.MsgTypeURL(&crosschaintypes.MsgUpdateTssAddress{}),
			sdk.MsgTypeURL(&crosschaintypes.MsgWhitelistERC20{}),
			sdk.MsgTypeURL(&fungibletypes.MsgDeployFungibleCoinZRC20{}),
			sdk.MsgTypeURL(&fungibletypes.MsgUpdateContractBytecode{}),
			sdk.MsgTypeURL(&fungibletypes.MsgUpdateSystemContract{}),
			sdk.MsgTypeURL(&fungibletypes.MsgUpdateGatewayContract{}),
			sdk.MsgTypeURL(&fungibletypes.MsgRemoveForeignCoin{}),
			sdk.MsgTypeURL(&fungibletypes.MsgUpdateZRC20Name{}),
			sdk.MsgTypeURL(&observertypes.MsgUpdateObserver{}),
			sdk.MsgTypeURL(&observertypes.MsgAddObserver{}),
			sdk.MsgTypeURL(&observertypes.MsgRemoveChainParams{}),
			sdk.MsgTypeURL(&types.MsgAddAuthorization{}),
			sdk.MsgTypeURL(&types.MsgRemoveAuthorization{}),
			sdk.MsgTypeURL(&types.MsgUpdateChainInfo{}),
			sdk.MsgTypeURL(&types.MsgRemoveChainInfo{}),
			sdk.MsgTypeURL(&lightclienttypes.MsgEnableHeaderVerification{}),
			sdk.MsgTypeURL(&observertypes.MsgUpdateChainParams{}),
			sdk.MsgTypeURL(&fungibletypes.MsgBurnFungibleModuleAsset{}),
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
