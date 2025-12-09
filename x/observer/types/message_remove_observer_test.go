package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestMsgRemoveObserver_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgRemoveObserver
		err  error
	}{
		{
			name: "invalid creator",
			msg: types.NewMsgRemoveObserver(
				"invalid_address",
				sample.AccAddress(),
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid observer address",
			msg: types.NewMsgRemoveObserver(
				sample.AccAddress(),
				"invalid_address",
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid addresses",
			msg: types.NewMsgRemoveObserver(
				sample.AccAddress(),
				sample.AccAddress(),
			),
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

func TestMsgRemoveObserver_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgRemoveObserver
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgRemoveObserver{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgRemoveObserver{
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

func TestMsgRemoveObserver_Type(t *testing.T) {
	msg := types.MsgRemoveObserver{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgRemoveObserver, msg.Type())
}

func TestMsgRemoveObserver_Route(t *testing.T) {
	msg := types.MsgRemoveObserver{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgRemoveObserver_GetSignBytes(t *testing.T) {
	msg := types.MsgRemoveObserver{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
