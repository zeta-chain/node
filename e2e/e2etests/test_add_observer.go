package e2etests

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
	"github.com/zeta-chain/node/pkg/observer_info"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	emissionstypes "github.com/zeta-chain/node/x/emissions/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

const (
	StakeAmount = "1000000000000000000000azeta"
)

func AddNewObserver(r *runner.E2ERunner) {
	FundHotkeyAccountForNonValidatorNode(r)
	StakeToBecomeValidator(r)
	AddNodeAccount(r)
	AddObserverAccount(r)
	AddGrants(r)
}

func StakeToBecomeValidator(r *runner.E2ERunner) {
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

	zetaTxServer := r.ZetaTxServer
	validatorsTxServer := zetaTxServer.UpdateKeyring(validatorsKeyring)

	_, err = validatorsTxServer.BroadcastTx(operator.Name, msg)
	require.NoError(r, err, "failed to broadcast transaction")
}

func AddNodeAccount(r *runner.E2ERunner) {
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

func AddObserverAccount(r *runner.E2ERunner) {
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

func AddGrants(r *runner.E2ERunner) {
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

func FundHotkeyAccountForNonValidatorNode(r *runner.E2ERunner) {
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

func GetNewValidatorInfo(r *runner.E2ERunner) *keyring.Record {
	record, err := r.ZetaTxServer.GetValidatorsKeyring().Key("operator-new-validator")
	require.NoError(r, err, "failed to get operator-new-validator key")
	return record
}

func FetchHotkeyAddress(r *runner.E2ERunner) observer_info.ObserverInfoReader {
	cmd := exec.Command("ssh", "-q", "root@zetaclient-new-validator", "cat ~/.zetacored/os.json")
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	require.NoError(r, err)
	output := out.String()

	observerInfo := observer_info.ObserverInfoReader{}

	err = json.Unmarshal([]byte(output), &observerInfo)
	require.NoError(r, err)

	return observerInfo
}

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
