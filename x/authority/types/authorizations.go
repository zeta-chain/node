package types

import (
	"fmt"

	"cosmossdk.io/errors"
)

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

		{MsgUrl: "/zetachain.zetacore.observer.MsgUpdateChainParams",
			AuthorizedPolicy: PolicyType_groupOperational},
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
			MsgUrl:           "/zetachain.zetacore.lightclient.MsgDisableHeaderVerification",
			AuthorizedPolicy: PolicyType_groupEmergency,
		},
	}

	return AuthorizationList{
		Authorizations: authorizations,
	}
}

func (a *AuthorizationList) AddAuthorizations(authorizationList AuthorizationList) {
	a.Authorizations = append(a.Authorizations, authorizationList.Authorizations...)
}

func (a *AuthorizationList) RemoveAuthorizations(removeList AuthorizationList) {
	for _, removeAuth := range removeList.Authorizations {
		for i, auth := range a.Authorizations {
			if auth.MsgUrl == removeAuth.MsgUrl {
				a.Authorizations = append(a.Authorizations[:i], a.Authorizations[i+1:]...)
			}
		}
	}
}

func (a *AuthorizationList) Validate() error {
	if len(a.Authorizations) == 0 {
		return errors.Wrap(ErrInValidAuthorizationList, "empty authorization list")
	}
	checkMsgUrls := make(map[string]bool)
	for _, authorization := range a.Authorizations {
		if checkMsgUrls[authorization.MsgUrl] {
			return errors.Wrap(ErrInValidAuthorizationList, fmt.Sprintf("duplicate message url: %s", authorization.MsgUrl))
		}
		checkMsgUrls[authorization.MsgUrl] = true
	}
	return nil
}
