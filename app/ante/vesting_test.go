package ante_test

import (
	"math/rand"
	"testing"
	"time"

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	vesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/app"
	"github.com/zeta-chain/node/app/ante"
	serverconfig "github.com/zeta-chain/node/server/config"
	"github.com/zeta-chain/node/testutil/sample"
)

func TestVesting_AnteHandle(t *testing.T) {
	txConfig := app.MakeEncodingConfig(serverconfig.DefaultEVMChainID).TxConfig

	testPrivKey, testAddress := sample.PrivKeyAddressPair()
	_, testAddress2 := sample.PrivKeyAddressPair()

	decorator := ante.NewVestingAccountDecorator()

	tests := []struct {
		name       string
		msg        sdk.Msg
		wantHasErr bool
		wantErr    string
	}{
		{
			"MsgCreateVestingAccount",
			vesting.NewMsgCreateVestingAccount(
				testAddress, testAddress2,
				sdk.NewCoins(sdk.NewInt64Coin("azeta", 100_000_000)),
				time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				false,
			),
			true,
			"MsgTypeURL /cosmos.vesting.v1beta1.MsgCreateVestingAccount not supported",
		},
		{
			"MsgCreatePermanentLockedAccount",
			vesting.NewMsgCreatePermanentLockedAccount(
				testAddress, testAddress2,
				sdk.NewCoins(sdk.NewInt64Coin("azeta", 100_000_000)),
			),
			true,
			"MsgTypeURL /cosmos.vesting.v1beta1.MsgCreatePermanentLockedAccount not supported",
		},
		{
			"MsgCreatePeriodicVestingAccount",
			vesting.NewMsgCreatePeriodicVestingAccount(
				testAddress, testAddress2,
				time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
				nil,
			),
			true,
			"MsgTypeURL /cosmos.vesting.v1beta1.MsgCreatePeriodicVestingAccount not supported",
		},
		{
			"Non blocked message",
			banktypes.NewMsgSend(
				testAddress, testAddress2,
				sdk.NewCoins(sdk.NewInt64Coin("azeta", 100_000_000)),
			),
			false,
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := simtestutil.GenSignedMockTx(
				rand.New(rand.NewSource(time.Now().UnixNano())),
				txConfig,
				[]sdk.Msg{
					tt.msg,
				},
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

			_, err = decorator.AnteHandle(ctx, tx, false, mmd.AnteHandle)

			if tt.wantHasErr {
				require.ErrorContains(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
