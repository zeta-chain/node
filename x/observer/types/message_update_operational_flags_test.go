package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestMsgUpdateOperationalFlags_ValidateBasic(t *testing.T) {
	tt := []struct {
		name string
		msg  *types.MsgUpdateOperationalFlags
		err  require.ErrorAssertionFunc
	}{
		{
			name: "invalid creator address",
			msg:  types.NewMsgUpdateOperationalFlags("invalid", types.OperationalFlags{}),
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorContains(t, err, "invalid creator address")
			},
		},
		{
			name: "invalid operational flags",
			msg: types.NewMsgUpdateOperationalFlags(
				sample.AccAddress(),
				types.OperationalFlags{
					RestartHeight: -1,
				},
			),
			err: func(t require.TestingT, err error, i ...interface{}) {
				require.ErrorContains(t, err, "invalid request")
			},
		},
		{
			name: "valid",
			msg:  types.NewMsgUpdateOperationalFlags(sample.AccAddress(), sample.OperationalFlags()),
			err:  require.NoError,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.err(t, tc.msg.ValidateBasic())
		})
	}
}

func TestMsgUpdateOperationalFlags_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgUpdateOperationalFlags
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgUpdateOperationalFlags{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgUpdateOperationalFlags{
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

func TestMsgUpdateOperationalFlags_Type(t *testing.T) {
	msg := types.MsgUpdateOperationalFlags{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgUpdateOperationalFlags, msg.Type())
}

func TestMsgUpdateOperationalFlags_Route(t *testing.T) {
	msg := types.MsgUpdateOperationalFlags{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgUpdateOperationalFlags_GetSignBytes(t *testing.T) {
	msg := types.MsgUpdateOperationalFlags{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
