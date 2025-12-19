package local

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	emissionstypes "github.com/zeta-chain/node/x/emissions/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

const (
	StakeAmount = "1000000000000000000000azeta"
)

// addNewObserver is a test function that adds a new observer to the network.
func addNewObserver(r *runner.E2ERunner) {
	fundHotkeyAccountForNewNode(r)
	stakeToBecomeValidator(r)
	addNodeAccount(r)
	addObserverAccount(r)
	addGrants(r)
}

// stakeToBecomeValidator is a helper function that stakes tokens to become a validator.
func stakeToBecomeValidator(r *runner.E2ERunner) {
	stakeTokens, err := sdk.ParseCoinNormalized(StakeAmount)
	require.NoError(r, err, "failed to parse coin")

	validatorsKeyring := r.ZetaTxServer.GetValidatorsKeyring()
	pubkey, err := utils.FetchNodePubkey("zetacore-new-validator")
	require.NoError(r, err)
	var pk cryptotypes.PubKey
	cdc := r.ZetaTxServer.GetCodec()
	err = cdc.UnmarshalInterfaceJSON([]byte(pubkey), &pk)
	require.NoError(r, err, "failed to unmarshal pubkey")

	operator := getNewValidatorInfo(r)
	address, err := operator.GetAddress()
	require.NoError(r, err, "failed to get address")

	msg, err := stakingtypes.NewMsgCreateValidator(
		sdk.ValAddress(address).String(),
		pk,
		stakeTokens,
		stakingtypes.NewDescription("zetacore-new-validator", "", "", "", ""),
		stakingtypes.NewCommissionRates(
			sdkmath.LegacyMustNewDecFromStr("0.10"),
			sdkmath.LegacyMustNewDecFromStr("0.20"),
			sdkmath.LegacyMustNewDecFromStr("0.010"),
		),
		sdkmath.OneInt(),
	)
	require.NoError(r, err, "failed to create MsgCreateValidator")

	zetaTxServer := r.ZetaTxServer
	validatorsTxServer := zetaTxServer.UpdateKeyring(validatorsKeyring)

	_, err = validatorsTxServer.BroadcastTx(operator.Name, msg)
	require.NoError(r, err, "failed to broadcast transaction")
}

// addNodeAccount adds the node account of a new validator to the network.
func addNodeAccount(r *runner.E2ERunner) {
	observerInfo, err := utils.FetchHotkeyAddress("zetaclient-new-validator")
	require.NoError(r, err)
	msg := observertypes.MsgAddObserver{
		Creator:                 r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		ObserverAddress:         observerInfo.ObserverAddress,
		ZetaclientGranteePubkey: observerInfo.ZetaClientGranteePubKey,
		AddNodeAccountOnly:      true,
	}
	_, err = r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, &msg)
	require.NoError(r, err)
}

// addObserverAccount adds a validator account to the observer set.
func addObserverAccount(r *runner.E2ERunner) {
	observerInfo, err := utils.FetchHotkeyAddress("zetaclient-new-validator")
	require.NoError(r, err)
	msg := observertypes.MsgAddObserver{
		Creator:                 r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		ObserverAddress:         observerInfo.ObserverAddress,
		ZetaclientGranteePubkey: observerInfo.ZetaClientGranteePubKey,
		AddNodeAccountOnly:      false,
	}
	_, err = r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, &msg)
	require.NoError(r, err)
}

// addGrants adds the necessary grants between operator and hotkey accounts.
func addGrants(r *runner.E2ERunner) {
	addGrantsWithHotkey(r, "zetaclient-new-validator", "zetaclient-new-validator")
}

// addGrantsWithHotkey grants the necessary permissions from granter to grantee
func addGrantsWithHotkey(r *runner.E2ERunner, granter, grantee string) {
	granterInfo, err := utils.FetchHotkeyAddress(granter)
	require.NoError(r, err)

	granteeInfo, err := utils.FetchHotkeyAddress(grantee)
	require.NoError(r, err)

	txTypes := crosschaintypes.GetAllAuthzZetaclientTxTypes()
	validatorsKeyring := r.ZetaTxServer.GetValidatorsKeyring()

	zetaTxServer := r.ZetaTxServer
	validatorsTxServer := zetaTxServer.UpdateKeyring(validatorsKeyring)

	operator := getNewValidatorInfo(r)
	for _, txType := range txTypes {
		msg, err := authz.NewMsgGrant(
			sdk.MustAccAddressFromBech32(granterInfo.ObserverAddress),
			sdk.MustAccAddressFromBech32(granteeInfo.ZetaClientGranteeAddress),
			authz.NewGenericAuthorization(txType),
			nil,
		)
		require.NoError(r, err)
		_, err = validatorsTxServer.BroadcastTx(operator.Name, msg)
		require.NoError(r, err, "failed to broadcast transaction")
	}

	msg, err := authz.NewMsgGrant(
		sdk.MustAccAddressFromBech32(granterInfo.ObserverAddress),
		sdk.MustAccAddressFromBech32(r.ZetaTxServer.MustGetAccountAddressFromName(utils.UserEmissionsWithdrawName)),
		authz.NewGenericAuthorization(sdk.MsgTypeURL(&emissionstypes.MsgWithdrawEmission{})),
		nil,
	)
	require.NoError(r, err)
	_, err = validatorsTxServer.BroadcastTx(operator.Name, msg)
	require.NoError(r, err, "failed to broadcast transaction")
}

// fundHotkeyAccountForNewNode funds the hotkey address of a new validator.
func fundHotkeyAccountForNewNode(r *runner.E2ERunner) {
	fundHotkeyAccount(r, "zetaclient-new-validator")
}

// fundHotkeyAccount funds the hotkey address of the specified zetaclient.
func fundHotkeyAccount(r *runner.E2ERunner, zetaclientName string) {
	amount, err := sdk.ParseCoinNormalized("100000000000000000000azeta")
	require.NoError(r, err, "failed to parse coin")

	observerInfo, err := utils.FetchHotkeyAddress(zetaclientName)
	require.NoError(r, err)

	validatorsKeyring := r.ZetaTxServer.GetValidatorsKeyring()
	zetaTxServer := r.ZetaTxServer
	validatorsTxServer := zetaTxServer.UpdateKeyring(validatorsKeyring)

	operator := getNewValidatorInfo(r)
	operatorAddress, err := operator.GetAddress()
	require.NoError(r, err, "failed to get address for operator")

	msg := banktypes.MsgSend{
		FromAddress: operatorAddress.String(),
		ToAddress:   observerInfo.ZetaClientGranteeAddress,
		Amount:      sdk.NewCoins(amount),
	}

	_, err = validatorsTxServer.BroadcastTx(operator.Name, &msg)
	require.NoError(r, err, "failed to broadcast transaction")
}

// getNewValidatorInfo retrieves the keyring record for the new validator from the keyring.
// operator-new-validator is created when copying the keys into the container
func getNewValidatorInfo(r *runner.E2ERunner) *keyring.Record {
	record, err := r.ZetaTxServer.GetValidatorsKeyring().Key("operator-new-validator")
	require.NoError(r, err, "failed to get operator-new-validator key")
	return record
}
