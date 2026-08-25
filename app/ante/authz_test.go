package ante_test

import (
	"math/rand"
	"testing"
	"time"

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	vesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/cosmos/cosmos-sdk/x/group"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/app"
	"github.com/zeta-chain/node/app/ante"
	serverconfig "github.com/zeta-chain/node/server/config"
	"github.com/zeta-chain/node/testutil/sample"
)

// TestDisabledAuthzMsgs_WiresVestingCreation asserts the production disabled-msg list actually
// contains the three vesting-account creation messages. This binds the ante test to app.go's real
// wiring: dropping any of them from app.DisabledAuthzMsgs would fail here.
func TestDisabledAuthzMsgs_WiresVestingCreation(t *testing.T) {
	disabled := app.DisabledAuthzMsgs()
	require.Contains(t, disabled, sdk.MsgTypeURL(&vesting.MsgCreateVestingAccount{}))
	require.Contains(t, disabled, sdk.MsgTypeURL(&vesting.MsgCreatePermanentLockedAccount{}))
	require.Contains(t, disabled, sdk.MsgTypeURL(&vesting.MsgCreatePeriodicVestingAccount{}))
}

// TestAuthzLimiter_AnteHandle verifies the decorator blocks the disabled vesting msgs when they are
// wrapped in authz.MsgExec OR group.MsgSubmitProposal (including nested), and lets non-disabled
// messages through. The decorator is built from app.DisabledAuthzMsgs() (the real production list)
// so the test cannot pass while the chain wiring regresses.
func TestAuthzLimiter_AnteHandle(t *testing.T) {
	// ARRANGE
	txConfig := app.MakeEncodingConfig(serverconfig.DefaultEVMChainID).TxConfig

	testPrivKey, testAddress := sample.PrivKeyAddressPair()
	_, testAddress2 := sample.PrivKeyAddressPair()
	_, policyAddress := sample.PrivKeyAddressPair()

	decorator := ante.NewAuthzLimiterDecorator(app.DisabledAuthzMsgs()...)

	createVestingMsg := vesting.NewMsgCreateVestingAccount(
		testAddress, testAddress2,
		sdk.NewCoins(sdk.NewInt64Coin("azeta", 100_000_000)),
		time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		false,
	)
	permanentLockedMsg := vesting.NewMsgCreatePermanentLockedAccount(
		testAddress, testAddress2,
		sdk.NewCoins(sdk.NewInt64Coin("azeta", 100_000_000)),
	)
	periodicVestingMsg := vesting.NewMsgCreatePeriodicVestingAccount(
		testAddress, testAddress2,
		time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		nil,
	)
	bankSend := banktypes.NewMsgSend(
		testAddress, testAddress2,
		sdk.NewCoins(sdk.NewInt64Coin("azeta", 100_000_000)),
	)

	groupProposal := func(t *testing.T, inner sdk.Msg) sdk.Msg {
		msg, err := group.NewMsgSubmitProposal(
			policyAddress.String(),
			[]string{testAddress.String()},
			[]sdk.Msg{inner},
			"",
			group.Exec_EXEC_UNSPECIFIED,
			"title",
			"summary",
		)
		require.NoError(t, err)
		return msg
	}
	authzExec := func(inner sdk.Msg) sdk.Msg {
		msg := authz.NewMsgExec(testAddress, []sdk.Msg{inner})
		return &msg
	}
	authzGrant := func(t *testing.T, msgTypeURL string) sdk.Msg {
		msg, err := authz.NewMsgGrant(
			testAddress, testAddress2,
			authz.NewGenericAuthorization(msgTypeURL),
			nil,
		)
		require.NoError(t, err)
		return msg
	}
	govProposal := func(t *testing.T, inner sdk.Msg) sdk.Msg {
		msg, err := govv1.NewMsgSubmitProposal(
			[]sdk.Msg{inner},
			sdk.NewCoins(sdk.NewInt64Coin("azeta", 100_000_000)),
			testAddress.String(),
			"",
			"title",
			"summary",
			false,
		)
		require.NoError(t, err)
		return msg
	}

	tests := []struct {
		name       string
		msg        sdk.Msg
		wantHasErr bool
		wantErr    string
	}{
		{
			"group proposal wrapping MsgCreateVestingAccount is blocked",
			groupProposal(t, createVestingMsg),
			true,
			"found disabled msg type",
		},
		{
			"group proposal wrapping MsgCreatePermanentLockedAccount is blocked",
			groupProposal(t, permanentLockedMsg),
			true,
			"found disabled msg type",
		},
		{
			"group proposal wrapping MsgCreatePeriodicVestingAccount is blocked",
			groupProposal(t, periodicVestingMsg),
			true,
			"found disabled msg type",
		},
		{
			"authz exec wrapping a vesting msg is blocked",
			authzExec(createVestingMsg),
			true,
			"found disabled msg type",
		},
		{
			"authz grant of a vesting-create authorization is blocked",
			authzGrant(t, sdk.MsgTypeURL(&vesting.MsgCreateVestingAccount{})),
			true,
			"found disabled msg type",
		},
		{
			"authz exec wrapping a group proposal wrapping a vesting msg is blocked (nested)",
			authzExec(groupProposal(t, createVestingMsg)),
			true,
			"found disabled msg type",
		},
		{
			"group proposal wrapping a non-disabled msg is allowed",
			groupProposal(t, bankSend),
			false,
			"",
		},
		{
			"top-level vesting msg is not blocked here (VestingAccountDecorator owns that)",
			createVestingMsg,
			false,
			"",
		},
		{
			// Documents a known gap: AuthzLimiterDecorator has no gov.MsgSubmitProposal case, so a
			// governance proposal can still carry a vesting-create. That path is privileged (a passed
			// gov vote, self-funded from the gov account), unlike the permissionless group path.
			"gov proposal wrapping a vesting msg is NOT blocked by this decorator (gov is out of scope)",
			govProposal(t, createVestingMsg),
			false,
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := simtestutil.GenSignedMockTx(
				rand.New(rand.NewSource(time.Now().UnixNano())),
				txConfig,
				[]sdk.Msg{tt.msg},
				sdk.NewCoins(),
				simtestutil.DefaultGenTxGas,
				"testing-chain-id",
				[]uint64{0},
				[]uint64{0},
				testPrivKey,
			)
			require.NoError(t, err)

			mmd := MockAnteHandler{}
			ctx := sdk.Context{}.WithIsCheckTx(true)

			// ACT
			_, err = decorator.AnteHandle(ctx, tx, false, mmd.AnteHandle)

			// ASSERT
			if tt.wantHasErr {
				require.ErrorContains(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
