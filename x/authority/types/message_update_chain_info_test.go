package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/authority/types"
)

func TestMsgUpdateChainInfo_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgUpdateChainInfo
		err  error
	}{
		{
			name: "valid message",
			msg:  types.NewMsgUpdateChainInfo(sample.AccAddress()),
		},
		{
			name: "invalid creator address",
			msg:  types.NewMsgUpdateChainInfo("invalid"),
			err:  sdkerrors.ErrInvalidAddress,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgUpdateChainInfo_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    *types.MsgUpdateChainInfo
		panics bool
	}{
		{
			name:   "valid signer",
			msg:    types.NewMsgUpdateChainInfo(signer),
			panics: false,
		},
		{
			name:   "invalid signer",
			msg:    types.NewMsgUpdateChainInfo("invalid"),
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

func TestMsgUpdateChainInfo_Type(t *testing.T) {
	msg := types.NewMsgUpdateChainInfo(sample.AccAddress())
	require.Equal(t, types.TypeMsgUpdateChainInfo, msg.Type())
}

func TestMsgUpdateChainInfo_Route(t *testing.T) {
	msg := types.NewMsgUpdateChainInfo(sample.AccAddress())
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgUpdateChainInfo_GetSignBytes(t *testing.T) {
	msg := types.NewMsgUpdateChainInfo(sample.AccAddress())
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
