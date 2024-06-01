package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/authority/types"
)

func TestMsgAddAuthorization_ValidateBasic(t *testing.T) {
	tests := []struct {
		name      string
		msg       *types.MsgAddAuthorization
		expectErr require.ErrorAssertionFunc
	}{
		{
			name: "invalid creator address",
			msg:  types.NewMsgAddAuthorization("invalid", "url", types.PolicyType_groupAdmin),
			expectErr: func(t require.TestingT, err error, msgAndArgs ...interface{}) {
				require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
				require.Contains(t, err.Error(), "invalid creator address")
			},
		},
		{
			name: "invalid authorized policy",
			msg:  types.NewMsgAddAuthorization(sample.AccAddress(), "url", types.PolicyType_groupEmpty),
			expectErr: func(t require.TestingT, err error, msgAndArgs ...interface{}) {
				require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
				require.Contains(t, err.Error(), "invalid authorized policy")
			},
		},
		{
			name: "invalid msg url",
			msg:  types.NewMsgAddAuthorization(sample.AccAddress(), "", types.PolicyType_groupAdmin),
			expectErr: func(t require.TestingT, err error, msgAndArgs ...interface{}) {
				require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
				require.Contains(t, err.Error(), "invalid msg url")
			},
		},
		{
			name:      "valid message",
			msg:       types.NewMsgAddAuthorization(sample.AccAddress(), "url", types.PolicyType_groupAdmin),
			expectErr: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expectErr(t, tt.msg.ValidateBasic())
		})
	}
}

func TestMsgAddAuthorization_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    *types.MsgAddAuthorization
		panics bool
	}{
		{
			name:   "valid signer",
			msg:    types.NewMsgAddAuthorization(signer, "url", types.PolicyType_groupAdmin),
			panics: false,
		},
		{
			name:   "invalid signer",
			msg:    types.NewMsgAddAuthorization("creator", "url", types.PolicyType_groupAdmin),
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

func TestMsgAddAuthorization_Type(t *testing.T) {
	msg := types.NewMsgAddAuthorization(sample.AccAddress(), "url", types.PolicyType_groupAdmin)
	require.Equal(t, types.TypeMsgAddAuthorization, msg.Type())
}

func TestMsgAddAuthorization_Route(t *testing.T) {
	msg := types.NewMsgAddAuthorization(sample.AccAddress(), "url", types.PolicyType_groupAdmin)
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgAddAuthorization_GetSignBytes(t *testing.T) {
	msg := types.NewMsgAddAuthorization(sample.AccAddress(), "url", types.PolicyType_groupAdmin)
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
