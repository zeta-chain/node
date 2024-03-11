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

func TestMsgAddObserver_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgAddObserver
		err  error
	}{
		{
			name: "invalid msg",
			msg: types.MsgAddObserver{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid creator",
			msg: types.MsgAddObserver{
				Creator:                 "invalid_address",
				ObserverAddress:         sample.AccAddress(),
				ZetaclientGranteePubkey: sample.PubKeyString(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid pubkey",
			msg: types.MsgAddObserver{
				Creator:                 sample.AccAddress(),
				ObserverAddress:         sample.AccAddress(),
				ZetaclientGranteePubkey: "sample.PubKey()",
			},
			err: sdkerrors.ErrInvalidPubKey,
		},
		{
			name: "invalid observer address",
			msg: types.MsgAddObserver{
				Creator:                 sample.AccAddress(),
				ObserverAddress:         "invalid_address",
				ZetaclientGranteePubkey: sample.PubKeyString(),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg: types.MsgAddObserver{
				Creator:                 sample.AccAddress(),
				ObserverAddress:         sample.AccAddress(),
				ZetaclientGranteePubkey: sample.PubKeyString(),
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

func TestMsgAddObserver_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgAddObserver
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgAddObserver{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgAddObserver{
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

func TestMsgAddObserver_Type(t *testing.T) {
	msg := types.MsgAddObserver{
		Creator: sample.AccAddress(),
	}
	assert.Equal(t, types.TypeMsgAddObserver, msg.Type())
}

func TestMsgAddObserver_Route(t *testing.T) {
	msg := types.MsgAddObserver{
		Creator: sample.AccAddress(),
	}
	assert.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgAddObserver_GetSignBytes(t *testing.T) {
	msg := types.MsgAddObserver{
		Creator: sample.AccAddress(),
	}
	assert.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
