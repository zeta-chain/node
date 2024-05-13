package authorizations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func GetRequiredPolicy(msgURL string) authoritytypes.PolicyType {
	if CheckPolicyList(msgURL, OperationalPolicyMessageList) {
		return authoritytypes.PolicyType_groupOperational
	}
	if CheckPolicyList(msgURL, EmergencyPolicyMessageList) {
		return authoritytypes.PolicyType_groupEmergency
	}
	if CheckPolicyList(msgURL, AdminPolicyMessageList) {
		return authoritytypes.PolicyType_groupAdmin
	}
	return authoritytypes.PolicyType_emptyPolicyType
}

func CheckPolicyList(msgURl string, msgList []string) bool {
	for _, policy := range msgList {
		if policy == msgURl {
			return true
		}
	}
	return false
}

var OperationalPolicyMessageList = []string{
	// Crosschain admin messages
	sdk.MsgTypeURL(&crosschaintypes.MsgRefundAbortedCCTX{}),
	sdk.MsgTypeURL(&crosschaintypes.MsgAbortStuckCCTX{}),
	sdk.MsgTypeURL(&crosschaintypes.MsgUpdateRateLimiterFlags{}),
	sdk.MsgTypeURL(&crosschaintypes.MsgWhitelistERC20{}),
	// Fungible admin messages
	sdk.MsgTypeURL(&fungibletypes.MsgDeployFungibleCoinZRC20{}),
	sdk.MsgTypeURL(&fungibletypes.MsgDeploySystemContracts{}),
	sdk.MsgTypeURL(&fungibletypes.MsgRemoveForeignCoin{}),
	sdk.MsgTypeURL(&fungibletypes.MsgUpdateZRC20LiquidityCap{}),
	sdk.MsgTypeURL(&fungibletypes.MsgUpdateZRC20WithdrawFee{}),
	sdk.MsgTypeURL(&fungibletypes.MsgUnpauseZRC20{}),
	// Observer admin messages
	sdk.MsgTypeURL(&observertypes.MsgAddObserver{}),
	sdk.MsgTypeURL(&observertypes.MsgRemoveChainParams{}),
	sdk.MsgTypeURL(&observertypes.MsgResetChainNonces{}),
	sdk.MsgTypeURL(&observertypes.MsgUpdateChainParams{}),
	sdk.MsgTypeURL(&observertypes.MsgEnableCCTXFlags{}),
	sdk.MsgTypeURL(&observertypes.MsgUpdateGasPriceIncreaseFlags{}),
	// Lightclient admin messages
	sdk.MsgTypeURL(&lightclienttypes.MsgEnableHeaderVerification{}),
}

var EmergencyPolicyMessageList = []string{
	// Crosschain admin messages
	sdk.MsgTypeURL(&crosschaintypes.MsgAddToInTxTracker{}),
	sdk.MsgTypeURL(&crosschaintypes.MsgAddToOutTxTracker{}),
	sdk.MsgTypeURL(&crosschaintypes.MsgRemoveFromOutTxTracker{}),
	// Fungible admin messages
	sdk.MsgTypeURL(&fungibletypes.MsgPauseZRC20{}),
	// Observer admin messages
	sdk.MsgTypeURL(&observertypes.MsgUpdateKeygen{}),
	sdk.MsgTypeURL(&observertypes.MsgDisableCCTXFlags{}),
	// Lightclient admin messages
	sdk.MsgTypeURL(&lightclienttypes.MsgDisableHeaderVerification{}),
}

var AdminPolicyMessageList = []string{
	// Crosschain admin messages
	sdk.MsgTypeURL(&crosschaintypes.MsgMigrateTssFunds{}),
	sdk.MsgTypeURL(&crosschaintypes.MsgUpdateTssAddress{}),
	// Fungible admin messages
	sdk.MsgTypeURL(&fungibletypes.MsgUpdateContractBytecode{}),
	sdk.MsgTypeURL(&fungibletypes.MsgUpdateSystemContract{}),
	// Observer admin messages
	sdk.MsgTypeURL(&observertypes.MsgUpdateObserver{}),
}
