package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestNewMsgUpdateZRC20Name_ValidateBasics(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgUpdateZRC20Name
		err  error
	}{
		{
			name: "valid message",
			msg: types.NewMsgUpdateZRC20Name(
				sample.AccAddress(),
				sample.EthAddress().String(),
				"foo",
				"bar",
			),
		},
		{
			name: "valid message with empty name",
			msg: types.NewMsgUpdateZRC20Name(
				sample.AccAddress(),
				sample.EthAddress().String(),
				"",
				"bar",
			),
		},
		{
			name: "valid message with empty symbol",
			msg: types.NewMsgUpdateZRC20Name(
				sample.AccAddress(),
				sample.EthAddress().String(),
				"foo",
				"",
			),
		},

		{
			name: "invalid address",
			msg: types.NewMsgUpdateZRC20Name(
				"invalid_address",
				sample.EthAddress().String(),
				"foo",
				"bar",
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid contract address",
			msg: types.NewMsgUpdateZRC20Name(
				sample.AccAddress(),
				"invalid_address",
				"foo",
				"bar",
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "nothing to update",
			msg: types.NewMsgUpdateZRC20Name(
				sample.AccAddress(),
				sample.EthAddress().String(),
				"",
				"",
			),
			err: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestNewMsgUpdateZRC20Name_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgUpdateZRC20Name
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgUpdateZRC20Name{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgUpdateZRC20Name{
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

func TestNewMsgUpdateZRC20Name_Type(t *testing.T) {
	msg := types.MsgUpdateZRC20Name{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgUpdateZRC20Name, msg.Type())
}

func TestNewMsgUpdateZRC20Name_Route(t *testing.T) {
	msg := types.MsgUpdateZRC20Name{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestNewMsgUpdateZRC20Name_GetSignBytes(t *testing.T) {
	msg := types.MsgUpdateZRC20Name{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
