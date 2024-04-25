package pkg

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	crosschainTypes "github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

func AuthorizationTable() map[string]authoritytypes.PolicyType {
	return map[string]authoritytypes.PolicyType{
		// Crosschain admin messages
		sdk.MsgTypeURL(&crosschainTypes.MsgRefundAbortedCCTX{}):      authoritytypes.PolicyType_groupOperational,
		sdk.MsgTypeURL(&crosschainTypes.MsgAbortStuckCCTX{}):         authoritytypes.PolicyType_groupOperational,
		sdk.MsgTypeURL(&crosschainTypes.MsgAddToInTxTracker{}):       authoritytypes.PolicyType_groupEmergency,
		sdk.MsgTypeURL(&crosschainTypes.MsgAddToOutTxTracker{}):      authoritytypes.PolicyType_groupEmergency,
		sdk.MsgTypeURL(&crosschainTypes.MsgMigrateTssFunds{}):        authoritytypes.PolicyType_groupAdmin,
		sdk.MsgTypeURL(&crosschainTypes.MsgRefundAbortedCCTX{}):      authoritytypes.PolicyType_groupOperational,
		sdk.MsgTypeURL(&crosschainTypes.MsgRemoveFromOutTxTracker{}): authoritytypes.PolicyType_groupEmergency,
		sdk.MsgTypeURL(&crosschainTypes.MsgUpdateRateLimiterFlags{}): authoritytypes.PolicyType_groupOperational,
		sdk.MsgTypeURL(&crosschainTypes.MsgUpdateTssAddress{}):       authoritytypes.PolicyType_groupAdmin,
		sdk.MsgTypeURL(&crosschainTypes.MsgWhitelistERC20{}):         authoritytypes.PolicyType_groupOperational,

		// Fungible admin messages
		sdk.MsgTypeURL(&fungibletypes.MsgDeployFungibleCoinZRC20{}): authoritytypes.PolicyType_groupOperational,
		sdk.MsgTypeURL(&fungibletypes.MsgDeploySystemContracts{}):   authoritytypes.PolicyType_groupOperational,
		sdk.MsgTypeURL(&fungibletypes.MsgRemoveForeignCoin{}):       authoritytypes.PolicyType_groupOperational,
		sdk.MsgTypeURL(&fungibletypes.MsgUpdateContractBytecode{}):  authoritytypes.PolicyType_groupAdmin,
		sdk.MsgTypeURL(&fungibletypes.MsgUpdateSystemContract{}):    authoritytypes.PolicyType_groupAdmin,
		sdk.MsgTypeURL(&fungibletypes.MsgUpdateZRC20LiquidityCap{}): authoritytypes.PolicyType_groupOperational,
		sdk.MsgTypeURL(&fungibletypes.MsgUpdateZRC20WithdrawFee{}):  authoritytypes.PolicyType_groupOperational,
	}
}
