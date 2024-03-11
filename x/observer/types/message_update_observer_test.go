package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestNewMsgUpdateObserver_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgUpdateObserver
		err  error
	}{
		{
			name: "invalid creator",
			msg: types.MsgUpdateObserver{
				Creator:            "invalid_address",
				OldObserverAddress: sample.AccAddress(),
				NewObserverAddress: sample.AccAddress(),
				UpdateReason:       types.ObserverUpdateReason_AdminUpdate,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid old observer address",
			msg: types.MsgUpdateObserver{
				Creator:            sample.AccAddress(),
				OldObserverAddress: "invalid_address",
				NewObserverAddress: sample.AccAddress(),
				UpdateReason:       types.ObserverUpdateReason_AdminUpdate,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid new observer address",
			msg: types.MsgUpdateObserver{
				Creator:            sample.AccAddress(),
				OldObserverAddress: sample.AccAddress(),
				NewObserverAddress: "invalid_address",
				UpdateReason:       types.ObserverUpdateReason_AdminUpdate,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "old observer address is not creator",
			msg: types.MsgUpdateObserver{
				Creator:            sample.AccAddress(),
				OldObserverAddress: sample.AccAddress(),
				NewObserverAddress: sample.AccAddress(),
				UpdateReason:       types.ObserverUpdateReason_Tombstoned,
			},
			err: types.ErrUpdateObserver,
		},
		{
			name: "invalid Update Reason",
			msg: types.MsgUpdateObserver{
				Creator:            sample.AccAddress(),
				OldObserverAddress: sample.AccAddress(),
				NewObserverAddress: sample.AccAddress(),
				UpdateReason:       types.ObserverUpdateReason(100),
			},
			err: types.ErrUpdateObserver,
		},
		{
			name: "valid message",
			msg: types.MsgUpdateObserver{
				Creator:            sample.AccAddress(),
				OldObserverAddress: sample.AccAddress(),
				NewObserverAddress: sample.AccAddress(),
				UpdateReason:       types.ObserverUpdateReason_AdminUpdate,
			},
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
				assert.Equal(t, []sdk.AccAddress{sdk.MustAccAddressFromBech32(signer)}, signers)
			} else {
				assert.Panics(t, func() {
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
	assert.Equal(t, types.TypeMsgUpdateObserver, msg.Type())
}

func TestNewMsgUpdateObserver_Route(t *testing.T) {
	msg := types.MsgUpdateObserver{
		Creator: sample.AccAddress(),
	}
	assert.Equal(t, types.RouterKey, msg.Route())
}

func TestNewMsgUpdateObserver_GetSignBytes(t *testing.T) {
	msg := types.MsgUpdateObserver{
		Creator: sample.AccAddress(),
	}
	assert.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
