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
	"github.com/cosmos/cosmos-sdk/x/group"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/app"
	"github.com/zeta-chain/node/app/ante"
	serverconfig "github.com/zeta-chain/node/server/config"
	"github.com/zeta-chain/node/testutil/sample"
)

// disabledMsgs mirrors app.go's DisabledAuthzMsgs for the vesting entries: the three
// vesting-account creation messages must not be executable indirectly.
func disabledVestingMsgs() []string {
	return []string{
		sdk.MsgTypeURL(&vesting.MsgCreateVestingAccount{}),
		sdk.MsgTypeURL(&vesting.MsgCreatePermanentLockedAccount{}),
		sdk.MsgTypeURL(&vesting.MsgCreatePeriodicVestingAccount{}),
	}
}

// TestAuthzLimiter_AnteHandle verifies the decorator blocks disabled vesting msgs when they
// are wrapped in authz.MsgExec OR group.MsgSubmitProposal (including nested), and lets
// non-disabled messages through.
func TestAuthzLimiter_AnteHandle(t *testing.T) {
	// ARRANGE
	txConfig := app.MakeEncodingConfig(serverconfig.DefaultEVMChainID).TxConfig

	testPrivKey, testAddress := sample.PrivKeyAddressPair()
	_, testAddress2 := sample.PrivKeyAddressPair()
	_, policyAddress := sample.PrivKeyAddressPair()

	decorator := ante.NewAuthzLimiterDecorator(disabledVestingMsgs()...)

	createVestingMsg := vesting.NewMsgCreateVestingAccount(
		testAddress, testAddress2,
		sdk.NewCoins(sdk.NewInt64Coin("azeta", 100_000_000)),
		time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		false,
	)
	lockedAcctMsg := vesting.NewMsgCreatePermanentLockedAccount(
		testAddress, testAddress2,
		sdk.NewCoins(sdk.NewInt64Coin("azeta", 100_000_000)),
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
			groupProposal(t, lockedAcctMsg),
			true,
			"found disabled msg type",
		},
		{
			"authz exec wrapping a vesting msg is blocked",
			authzExec(lockedAcctMsg),
			true,
			"found disabled msg type",
		},
		{
			"authz exec wrapping a group proposal wrapping a vesting msg is blocked (nested)",
			authzExec(groupProposal(t, lockedAcctMsg)),
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
			lockedAcctMsg,
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
