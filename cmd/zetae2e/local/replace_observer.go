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

// replaceObserverFlow performs the flow to replace an existing observer with a new validator
// that reuses the hotkey and TSS shard from zetaclient2.
// Unlike addNewObserver, we don't add a new node account - we reuse the existing one and just
// update the observer set using MsgUpdateObserver to swap the operator address.
func replaceObserverFlow(r *runner.E2ERunner) {
	fundHotkeyAccountForReplace(r)
	stakeToBecomeValidatorReplace(r)
	addGrantsReplace(r)
	updateObserverReplace(r)
	enableCctxAfterReplacement(r)
}

// stakeToBecomeValidatorReplace stakes tokens to create a new validator using a dedicated operator key.
// It expects a validator container reachable via SSH as "zetacore-new-validator".
// NOTE: Adjust the host string below if your docker-compose uses a different service name
// (e.g., "zetacored_new_validator").
func stakeToBecomeValidatorReplace(r *runner.E2ERunner) {
	stakeTokens, err := sdk.ParseCoinNormalized(StakeAmount)
	require.NoError(r, err, "failed to parse coin")

	validatorsKeyring := r.ZetaTxServer.GetValidatorsKeyring()
	pubkey, err := utils.FetchNodePubkey("zetacore-new-validator")
	if err != nil {
		// try alternative service name used in some compose files
		pubkey, err = utils.FetchNodePubkey("zetacored_new_validator")
	}
	require.NoError(r, err)
	var pk cryptotypes.PubKey
	cdc := r.ZetaTxServer.GetCodec()
	err = cdc.UnmarshalInterfaceJSON([]byte(pubkey), &pk)
	require.NoError(r, err, "failed to unmarshal pubkey")

	operator := getNewValidatorReplaceInfo(r)
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

// addGrantsReplace grants the necessary permissions between the operator and the zetaclient2 hotkey
// so that the zetaclient process can continue operating under the new operator address.
func addGrantsReplace(r *runner.E2ERunner) {
	// grantee (hotkey) info from zetaclient2
	oldInfo, err := utils.FetchHotkeyAddress("zetaclient2")
	require.NoError(r, err)
	// granter (observer) is the new operator from replacement container
	newInfo, err := utils.FetchHotkeyAddress("zetaclient-new-validator")
	if err != nil {
		newInfo, err = utils.FetchHotkeyAddress("zetaclient_new_validator")
	}
	require.NoError(r, err)
	txTypes := crosschaintypes.GetAllAuthzZetaclientTxTypes()
	validatorsKeyring := r.ZetaTxServer.GetValidatorsKeyring()

	zetaTxServer := r.ZetaTxServer
	validatorsTxServer := zetaTxServer.UpdateKeyring(validatorsKeyring)

	operator := getNewValidatorReplaceInfo(r)
	for _, txType := range txTypes {
		msg, err := authz.NewMsgGrant(
			sdk.MustAccAddressFromBech32(newInfo.ObserverAddress),
			sdk.MustAccAddressFromBech32(oldInfo.ZetaClientGranteeAddress),
			authz.NewGenericAuthorization(txType),
			nil,
		)
		require.NoError(r, err)
		_, err = validatorsTxServer.BroadcastTx(operator.Name, msg)
		require.NoError(r, err, "failed to broadcast transaction")
	}

	msg, err := authz.NewMsgGrant(
		sdk.MustAccAddressFromBech32(newInfo.ObserverAddress),
		sdk.MustAccAddressFromBech32(r.ZetaTxServer.MustGetAccountAddressFromName(utils.UserEmissionsWithdrawName)),
		authz.NewGenericAuthorization(sdk.MsgTypeURL(&emissionstypes.MsgWithdrawEmission{})),
		nil,
	)
	require.NoError(r, err)
	_, err = validatorsTxServer.BroadcastTx(operator.Name, msg)
	require.NoError(r, err, "failed to broadcast transaction")
}

// fundHotkeyAccountForReplace funds the zetaclient2 grantee account from the new operator to ensure
// it has sufficient balance to pay fees after the observer replacement.
func fundHotkeyAccountForReplace(r *runner.E2ERunner) {
	amount, err := sdk.ParseCoinNormalized("100000000000000000000azeta")
	require.NoError(r, err, "failed to parse coin")

	observerInfo, err := utils.FetchHotkeyAddress("zetaclient2")
	require.NoError(r, err)

	validatorsKeyring := r.ZetaTxServer.GetValidatorsKeyring()
	zetaTxServer := r.ZetaTxServer
	validatorsTxServer := zetaTxServer.UpdateKeyring(validatorsKeyring)

	operator := getNewValidatorReplaceInfo(r)
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

// updateObserverReplace updates the observer set to replace the old operator address with the new one.
// It sends MsgUpdateObserver with AdminUpdate reason using the operational policy account.
func updateObserverReplace(r *runner.E2ERunner) {
	// old observer operator address (current observer) from zetaclient2 os.json
	oldInfo, err := utils.FetchHotkeyAddress("zetaclient2")
	require.NoError(r, err)
	oldObserver := oldInfo.ObserverAddress

	// new observer operator address from the new-validator container os.json
	newInfo, err := utils.FetchHotkeyAddress("zetaclient-new-validator")
	require.NoError(r, err)
	newObserver := newInfo.ObserverAddress

	msg := observertypes.NewMsgUpdateObserver(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		oldObserver,
		newObserver,
		observertypes.ObserverUpdateReason_AdminUpdate,
	)

	_, err = r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, msg)
	require.NoError(r, err, "failed to broadcast MsgUpdateObserver")
}

// getNewValidatorReplaceInfo retrieves the keyring record for the replacement validator operator.
func getNewValidatorReplaceInfo(r *runner.E2ERunner) *keyring.Record {
	record, err := r.ZetaTxServer.GetValidatorsKeyring().Key("operator-new-validator")
	require.NoError(r, err, "failed to get operator-new-validator key")
	return record
}

// enableCctxAfterReplacement re-enables inbound and outbound processing after the replacement.
func enableCctxAfterReplacement(r *runner.E2ERunner) {
	msgEnable := observertypes.NewMsgEnableCCTX(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.OperationalPolicyName),
		true,
		true,
	)
	_, err := r.ZetaTxServer.BroadcastTx(utils.OperationalPolicyName, msgEnable)
	require.NoError(r, err)
	// Give it a block to propagate
	r.WaitForBlocks(1)
}
