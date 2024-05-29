package types

import (
	"fmt"

	"cosmossdk.io/errors"
)

// DefaultAuthorizationsList list is the list of authorizations that presently exist in the system.
// This is the minimum set of authorizations that are required to be set when the authorization table is deployed
func DefaultAuthorizationsList() AuthorizationList {
	var authorizations []Authorization

	authorizations = []Authorization{
		// OperationalPolicyMessageList
		{
			MsgUrl:           "/zetachain.zetacore.crosschain.MsgRefundAbortedCCTX",
			AuthorizedPolicy: PolicyType_groupOperational},
		{
			MsgUrl:           "/zetachain.zetacore.crosschain.MsgAbortStuckCCTX",
			AuthorizedPolicy: PolicyType_groupOperational,
		},
		{
			MsgUrl:           "/zetachain.zetacore.crosschain.MsgUpdateRateLimiterFlags",
			AuthorizedPolicy: PolicyType_groupOperational,
		},
		{
			MsgUrl:           "/zetachain.zetacore.crosschain.MsgWhitelistERC20",
			AuthorizedPolicy: PolicyType_groupOperational},
		{
			MsgUrl:           "/zetachain.zetacore.fungible.MsgDeployFungibleCoinZRC20",
			AuthorizedPolicy: PolicyType_groupOperational,
		},
		{
			MsgUrl:           "/zetachain.zetacore.fungible.MsgDeploySystemContracts",
			AuthorizedPolicy: PolicyType_groupOperational,
		},
		{
			MsgUrl:           "/zetachain.zetacore.fungible.MsgRemoveForeignCoin",
			AuthorizedPolicy: PolicyType_groupOperational,
		},
		{
			MsgUrl:           "/zetachain.zetacore.fungible.MsgUpdateZRC20LiquidityCap",
			AuthorizedPolicy: PolicyType_groupOperational,
		},
		{
			MsgUrl:           "/zetachain.zetacore.fungible.MsgUpdateZRC20WithdrawFee",
			AuthorizedPolicy: PolicyType_groupOperational,
		},
		{
			MsgUrl:           "/zetachain.zetacore.fungible.MsgUnpauseZRC20",
			AuthorizedPolicy: PolicyType_groupOperational,
		},
		{
			MsgUrl:           "/zetachain.zetacore.observer.MsgAddObserver",
			AuthorizedPolicy: PolicyType_groupOperational,
		},
		{
			MsgUrl:           "/zetachain.zetacore.observer.MsgRemoveChainParams",
			AuthorizedPolicy: PolicyType_groupOperational,
		},
		{
			MsgUrl:           "/zetachain.zetacore.observer.MsgResetChainNonces",
			AuthorizedPolicy: PolicyType_groupOperational,
		},

		{
			MsgUrl:           "/zetachain.zetacore.observer.MsgUpdateChainParams",
			AuthorizedPolicy: PolicyType_groupOperational,
		},
		{
			MsgUrl:           "/zetachain.zetacore.observer.MsgEnableCCTX",
			AuthorizedPolicy: PolicyType_groupOperational,
		},
		{
			MsgUrl:           "/zetachain.zetacore.observer.MsgUpdateGasPriceIncreaseFlags",
			AuthorizedPolicy: PolicyType_groupOperational,
		},
		{
			MsgUrl:           "/zetachain.zetacore.lightclient.MsgEnableHeaderVerification",
			AuthorizedPolicy: PolicyType_groupOperational,
		},
		// AdminPolicyMessageList
		{
			MsgUrl:           "/zetachain.zetacore.crosschain.MsgMigrateTssFunds",
			AuthorizedPolicy: PolicyType_groupAdmin,
		},
		{
			MsgUrl:           "/zetachain.zetacore.crosschain.MsgUpdateTssAddress",
			AuthorizedPolicy: PolicyType_groupAdmin,
		},
		{
			MsgUrl:           "/zetachain.zetacore.fungible.MsgUpdateContractBytecode",
			AuthorizedPolicy: PolicyType_groupAdmin,
		},
		{
			MsgUrl:           "/zetachain.zetacore.fungible.MsgUpdateSystemContract",
			AuthorizedPolicy: PolicyType_groupAdmin,
		},
		{
			MsgUrl:           "/zetachain.zetacore.observer.MsgUpdateObserver",
			AuthorizedPolicy: PolicyType_groupAdmin,
		},
		// EmergencyPolicyMessageList
		{
			MsgUrl:           "/zetachain.zetacore.crosschain.MsgAddInboundTracker",
			AuthorizedPolicy: PolicyType_groupEmergency,
		},
		{
			MsgUrl:           "/zetachain.zetacore.crosschain.MsgAddOutboundTracker",
			AuthorizedPolicy: PolicyType_groupEmergency,
		},
		{
			MsgUrl:           "/zetachain.zetacore.crosschain.MsgRemoveOutboundTracker",
			AuthorizedPolicy: PolicyType_groupEmergency,
		},
		{
			MsgUrl:           "/zetachain.zetacore.fungible.MsgPauseZRC20",
			AuthorizedPolicy: PolicyType_groupEmergency,
		},
		{
			MsgUrl:           "/zetachain.zetacore.observer.MsgUpdateKeygen",
			AuthorizedPolicy: PolicyType_groupEmergency,
		},
		{
			MsgUrl:           "/zetachain.zetacore.observer.MsgDisableCCTX",
			AuthorizedPolicy: PolicyType_groupEmergency,
		},
		{
			MsgUrl:           "/zetachain.zetacore.lightclient.MsgDisableHeaderVerification",
			AuthorizedPolicy: PolicyType_groupEmergency,
		},
	}

	return AuthorizationList{
		Authorizations: authorizations,
	}
}

// SetAuthorizations adds the authorization to the list. If the authorization already exists, it updates the policy.
func (a *AuthorizationList) SetAuthorizations(authorization Authorization) {
	for i, auth := range a.Authorizations {
		if auth.MsgUrl == authorization.MsgUrl {
			a.Authorizations[i].AuthorizedPolicy = authorization.AuthorizedPolicy
			return
		}
	}
	a.Authorizations = append(a.Authorizations, authorization)
}

// RemoveAuthorizations removes the authorization from the list. It does not check if the authorization exists or not.
func (a *AuthorizationList) RemoveAuthorizations(authorization Authorization) {
	for i, auth := range a.Authorizations {
		if auth.MsgUrl == authorization.MsgUrl {
			a.Authorizations = append(a.Authorizations[:i], a.Authorizations[i+1:]...)
		}
	}
}

// CheckAuthorizationExists checks if the authorization exists in the list.
func (a *AuthorizationList) CheckAuthorizationExists(authorization Authorization) bool {
	for _, auth := range a.Authorizations {
		if auth.MsgUrl == authorization.MsgUrl {
			return true
		}
	}
	return false
}

// GetAuthorizedPolicy returns the policy for the given message url.If the message url is not found,
// it returns an error and the first value of the enum.
func (a *AuthorizationList) GetAuthorizedPolicy(msgURL string) (PolicyType, error) {
	for _, auth := range a.Authorizations {
		if auth.MsgUrl == msgURL {
			return auth.AuthorizedPolicy, nil
		}
	}
	// Returning first value of enum, can consider adding a default value of `EmptyPolicy` in the enum.
	return PolicyType(0), ErrAuthorizationNotFound
}

// Validate checks if the authorization list is valid. It returns an error if the message url is duplicated with different policies.
// It does not check if the list is empty or not, as an empty list is also considered valid.
func (a *AuthorizationList) Validate() error {
	checkMsgUrls := make(map[string]bool)
	for _, authorization := range a.Authorizations {
		if checkMsgUrls[authorization.MsgUrl] {
			return errors.Wrap(
				ErrInValidAuthorizationList,
				fmt.Sprintf("duplicate message url: %s", authorization.MsgUrl),
			)
		}
		checkMsgUrls[authorization.MsgUrl] = true
	}
	return nil
}
