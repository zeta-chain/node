package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgUpdateGasPriceIncreaseFlags_ValidateBasic(t *testing.T) {
	tt := []struct {
		name string
		msg  *types.MsgUpdateGasPriceIncreaseFlags
		err  require.ErrorAssertionFunc
	}{
		{
			name: "invalid creator address",
			msg:  types.NewMsgUpdateGasPriceIncreaseFlags("invalid", types.DefaultGasPriceIncreaseFlags),
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.Contains(t, err.Error(), "invalid creator address")
			},
		},
		{
			name: "invalid gas price increase flags",
			msg: types.NewMsgUpdateGasPriceIncreaseFlags(
				sample.AccAddress(),
				types.GasPriceIncreaseFlags{
					EpochLength:             -1,
					RetryInterval:           1,
					GasPriceIncreasePercent: 1,
				},
			),
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.Contains(t, err.Error(), "invalid request")
			},
		},
		{
			name: "valid",
			msg:  types.NewMsgUpdateGasPriceIncreaseFlags(sample.AccAddress(), sample.GasPriceIncreaseFlags()),
			err:  require.NoError,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.err(t, tc.msg.ValidateBasic())
		})
	}
}

func TestMsgUpdateGasPriceIncreaseFlags_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgUpdateGasPriceIncreaseFlags
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgUpdateGasPriceIncreaseFlags{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgUpdateGasPriceIncreaseFlags{
				Creator: "invalid",
			},
			panics: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.panics {
				signers := tt.msg.GetSigners()
				require.Equal(t, []sdk.AccAddress{sdk.MustAccAddressFromBech32(signer)}, signers)
			} else {
				require.Panics(t, func() {
					tt.msg.GetSigners()
				})
			}
		})
	}
}

func TestMsgUpdateGasPriceIncreaseFlags_Type(t *testing.T) {
	msg := types.MsgUpdateGasPriceIncreaseFlags{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgUpdateGasPriceIncreaseFlags, msg.Type())
}

func TestMsgUpdateGasPriceIncreaseFlags_Route(t *testing.T) {
	msg := types.MsgUpdateGasPriceIncreaseFlags{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgUpdateGasPriceIncreaseFlags_GetSignBytes(t *testing.T) {
	msg := types.MsgUpdateGasPriceIncreaseFlags{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}

func TestGasPriceIncreaseFlags_Validate(t *testing.T) {
	tests := []struct {
		name        string
		gpf         types.GasPriceIncreaseFlags
		errContains string
	}{
		{
			name: "invalid epoch length",
			gpf: types.GasPriceIncreaseFlags{
				EpochLength:             -1,
				RetryInterval:           1,
				GasPriceIncreasePercent: 1,
			},
			errContains: "epoch length must be positive",
		},
		{
			name: "invalid retry interval",
			gpf: types.GasPriceIncreaseFlags{
				EpochLength:             1,
				RetryInterval:           -1,
				GasPriceIncreasePercent: 1,
			},
			errContains: "retry interval must be positive",
		},
		{
			name: "valid",
			gpf: types.GasPriceIncreaseFlags{
				EpochLength:             1,
				RetryInterval:           1,
				GasPriceIncreasePercent: 1,
			},
		},
		{
			name: "percent can be 0",
			gpf: types.GasPriceIncreaseFlags{
				EpochLength:             1,
				RetryInterval:           1,
				GasPriceIncreasePercent: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.gpf.Validate()
			if tt.errContains != "" {
				require.ErrorContains(t, err, tt.errContains)
				return
			}
			require.NoError(t, err)
		})
	}
}
