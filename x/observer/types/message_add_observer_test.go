package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestMsgAddObserver_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgAddObserver
		err  error
	}{
		{
			name: "invalid creator",
			msg: types.NewMsgAddObserver(
				"invalid_address",
				sample.AccAddress(),
				sample.PubKeyString(),
				true,
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid pubkey",
			msg: types.NewMsgAddObserver(
				sample.AccAddress(),
				sample.AccAddress(),
				"sample.PubKey()",
				true,
			),
			err: sdkerrors.ErrInvalidPubKey,
		},
		{
			name: "invalid observer address",
			msg: types.NewMsgAddObserver(
				sample.AccAddress(),
				"invalid_address",
				sample.PubKeyString(),
				true,
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "valid address",
			msg: types.NewMsgAddObserver(
				sample.AccAddress(),
				sample.AccAddress(),
				sample.PubKeyString(),
				true,
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
				require.Equal(t, []sdk.AccAddress{sdk.MustAccAddressFromBech32(signer)}, signers)
			} else {
				require.Panics(t, func() {
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
	require.Equal(t, types.TypeMsgAddObserver, msg.Type())
}

func TestMsgAddObserver_Route(t *testing.T) {
	msg := types.MsgAddObserver{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgAddObserver_GetSignBytes(t *testing.T) {
	msg := types.MsgAddObserver{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
