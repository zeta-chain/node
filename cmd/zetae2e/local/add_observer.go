package local

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strings"

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
	"github.com/zeta-chain/node/pkg/parsers"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	emissionstypes "github.com/zeta-chain/node/x/emissions/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

const (
	StakeAmount = "1000000000000000000000azeta"
)

// addNewObserver is a test function that adds a new observer to the network.
func addNewObserver(r *runner.E2ERunner) {
	fundHotkeyAccountForNonValidatorNode(r)
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
	pubkey := FetchNodePubkey(r)
	var pk cryptotypes.PubKey
	cdc := r.ZetaTxServer.GetCodec()
	err = cdc.UnmarshalInterfaceJSON([]byte(pubkey), &pk)
	require.NoError(r, err, "failed to unmarshal pubkey")

	operator := GetNewValidatorInfo(r)
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
	observerInfo := FetchHotkeyAddress(r)
	msg := observertypes.MsgAddObserver{
		Creator:                 r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		ObserverAddress:         observerInfo.ObserverAddress,
		ZetaclientGranteePubkey: observerInfo.ZetaClientGranteePubKey,
		AddNodeAccountOnly:      true,
	}
	_, err := r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, &msg)
	require.NoError(r, err)
}

// addObserverAccount adds a validator account to the observer set.
func addObserverAccount(r *runner.E2ERunner) {
	observerInfo := FetchHotkeyAddress(r)
	msg := observertypes.MsgAddObserver{
		Creator:                 r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		ObserverAddress:         observerInfo.ObserverAddress,
		ZetaclientGranteePubkey: observerInfo.ZetaClientGranteePubKey,
		AddNodeAccountOnly:      false,
	}
	_, err := r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, &msg)
	require.NoError(r, err)
}

// addGrants adds the necessary grants between operator and hotkey accounts.
func addGrants(r *runner.E2ERunner) {
	observerInfo := FetchHotkeyAddress(r)
	txTypes := crosschaintypes.GetAllAuthzZetaclientTxTypes()
	validatorsKeyring := r.ZetaTxServer.GetValidatorsKeyring()

	zetaTxServer := r.ZetaTxServer
	validatorsTxServer := zetaTxServer.UpdateKeyring(validatorsKeyring)

	operator := GetNewValidatorInfo(r)
	for _, txType := range txTypes {
		msg, err := authz.NewMsgGrant(
			sdk.MustAccAddressFromBech32(observerInfo.ObserverAddress),
			sdk.MustAccAddressFromBech32(observerInfo.ZetaClientGranteeAddress),
			authz.NewGenericAuthorization(txType),
			nil,
		)
		require.NoError(r, err)
		_, err = validatorsTxServer.BroadcastTx(operator.Name, msg)
		require.NoError(r, err, "failed to broadcast transaction")
	}

	msg, err := authz.NewMsgGrant(
		sdk.MustAccAddressFromBech32(observerInfo.ObserverAddress),
		sdk.MustAccAddressFromBech32(r.ZetaTxServer.MustGetAccountAddressFromName(utils.UserEmissionsWithdrawName)),
		authz.NewGenericAuthorization(sdk.MsgTypeURL(&emissionstypes.MsgWithdrawEmission{})),
		nil,
	)
	require.NoError(r, err)
	_, err = validatorsTxServer.BroadcastTx(operator.Name, msg)
	require.NoError(r, err, "failed to broadcast transaction")
}

// fundHotkeyAccountForNonValidatorNode funds the hotkey address of a new validator.
func fundHotkeyAccountForNonValidatorNode(r *runner.E2ERunner) {
	amount, err := sdk.ParseCoinNormalized("100000000000000000000azeta")
	require.NoError(r, err, "failed to parse coin")

	observerInfo := FetchHotkeyAddress(r)

	validatorsKeyring := r.ZetaTxServer.GetValidatorsKeyring()
	zetaTxServer := r.ZetaTxServer
	validatorsTxServer := zetaTxServer.UpdateKeyring(validatorsKeyring)

	operator := GetNewValidatorInfo(r)
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

// GetNewValidatorInfo retrieves the keyring record for the new validator from the keyring.
func GetNewValidatorInfo(r *runner.E2ERunner) *keyring.Record {
	record, err := r.ZetaTxServer.GetValidatorsKeyring().Key("operator-new-validator")
	require.NoError(r, err, "failed to get operator-new-validator key")
	return record
}

// FetchHotkeyAddress retrieves the hotkey address of a new validator.
func FetchHotkeyAddress(r *runner.E2ERunner) parsers.ObserverInfoReader {
	cmd := exec.Command("ssh", "-q", "root@zetaclient-new-validator", "cat ~/.zetacored/os.json")
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	require.NoError(r, err)
	output := out.String()

	observerInfo := parsers.ObserverInfoReader{}

	err = json.Unmarshal([]byte(output), &observerInfo)
	require.NoError(r, err)

	return observerInfo
}

// FetchNodePubkey retrieves the public key of the new validator node.
func FetchNodePubkey(r *runner.E2ERunner) string {
	cmd := exec.Command("ssh", "-q", "root@zetacore-new-validator", "zetacored tendermint show-validator")
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	require.NoError(r, err)
	output := out.String()
	output = strings.TrimSpace(output)
	return output
}
