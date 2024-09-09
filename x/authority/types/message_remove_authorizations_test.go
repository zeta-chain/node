package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/authority/types"
)

func TestMsgRemoveAuthorization_ValidateBasic(t *testing.T) {
	tests := []struct {
		name      string
		msg       *types.MsgRemoveAuthorization
		expectErr require.ErrorAssertionFunc
	}{
		{
			name: "invalid creator address",
			msg:  types.NewMsgRemoveAuthorization("invalid", "url"),
			expectErr: func(t require.TestingT, err error, msgAndArgs ...interface{}) {
				require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
				require.ErrorContains(t, err, "invalid creator address")
			},
		},
		{
			name: "invalid msg url",
			msg:  types.NewMsgRemoveAuthorization(sample.AccAddress(), ""),
			expectErr: func(t require.TestingT, err error, msgAndArgs ...interface{}) {
				require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
				require.ErrorContains(t, err, "invalid msg url")
			},
		},
		{
			name:      "valid message",
			msg:       types.NewMsgRemoveAuthorization(sample.AccAddress(), "url"),
			expectErr: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expectErr(t, tt.msg.ValidateBasic())
		})
	}
}

func TestMsgRemoveAuthorization_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    *types.MsgRemoveAuthorization
		panics bool
	}{
		{
			name:   "valid signer",
			msg:    types.NewMsgRemoveAuthorization(signer, "url"),
			panics: false,
		},
		{
			name:   "invalid signer",
			msg:    types.NewMsgRemoveAuthorization("creator", "url"),
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

func TestMsgRemoveAuthorization_Type(t *testing.T) {
	msg := types.NewMsgRemoveAuthorization(sample.AccAddress(), "url")
	require.Equal(t, types.TypeRemoveAuthorization, msg.Type())
}

func TestMsgRemoveAuthorization_Route(t *testing.T) {
	msg := types.NewMsgRemoveAuthorization(sample.AccAddress(), "url")
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgRemoveAuthorization_GetSignBytes(t *testing.T) {
	msg := types.NewMsgRemoveAuthorization(sample.AccAddress(), "url")
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
