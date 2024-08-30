package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestNewMsgUpdateObserver_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgUpdateObserver
		err  error
	}{
		{
			name: "invalid creator",
			msg: types.NewMsgUpdateObserver(
				"invalid_address",
				sample.AccAddress(),
				sample.AccAddress(),
				types.ObserverUpdateReason_AdminUpdate,
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid old observer address",
			msg: types.NewMsgUpdateObserver(
				sample.AccAddress(),
				"invalid_address",
				sample.AccAddress(),
				types.ObserverUpdateReason_AdminUpdate,
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid new observer address",
			msg: types.NewMsgUpdateObserver(
				sample.AccAddress(),
				sample.AccAddress(),
				"invalid_address",
				types.ObserverUpdateReason_AdminUpdate,
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "old observer address is not creator",
			msg: types.NewMsgUpdateObserver(
				sample.AccAddress(),
				sample.AccAddress(),
				sample.AccAddress(),
				types.ObserverUpdateReason_Tombstoned,
			),
			err: types.ErrUpdateObserver,
		},
		{
			name: "invalid Update Reason",
			msg: types.NewMsgUpdateObserver(
				sample.AccAddress(),
				sample.AccAddress(),
				sample.AccAddress(),
				types.ObserverUpdateReason(100),
			),
			err: types.ErrUpdateObserver,
		},
		{
			name: "valid message",
			msg: types.NewMsgUpdateObserver(
				sample.AccAddress(),
				sample.AccAddress(),
				sample.AccAddress(),
				types.ObserverUpdateReason_AdminUpdate,
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

func TestNewMsgUpdateObserver_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgUpdateObserver
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgUpdateObserver{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgUpdateObserver{
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

func TestNewMsgUpdateObserver_Type(t *testing.T) {
	msg := types.MsgUpdateObserver{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgUpdateObserver, msg.Type())
}

func TestNewMsgUpdateObserver_Route(t *testing.T) {
	msg := types.MsgUpdateObserver{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestNewMsgUpdateObserver_GetSignBytes(t *testing.T) {
	msg := types.MsgUpdateObserver{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
