package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestMessageUpdateTssAddress_ValidateBasic(t *testing.T) {
	keeper.SetConfig(false)
	tests := []struct {
		name  string
		msg   *types.MsgUpdateTssAddress
		error bool
	}{
		{
			name: "invalid creator",
			msg: types.NewMsgUpdateTssAddress(
				"invalid_address",
				sample.PubKeyString(),
			),
			error: true,
		},
		{
			name: "invalid pubkey",
			msg: types.NewMsgUpdateTssAddress(
				sample.AccAddress(),
				"zetapub1addwnpepq28c57cvcs0a2htsem5zxr6qnlvq9mzhmm",
			),
			error: true,
		},
		{
			name: "valid msg",
			msg: types.NewMsgUpdateTssAddress(
				sample.AccAddress(),
				sample.PubKeyString(),
			),
			error: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.error {
				require.Error(t, err)
				return
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMessageUpdateTssAddress_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgUpdateTssAddress
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgUpdateTssAddress{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgUpdateTssAddress{
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

func TestMessageUpdateTssAddress_Type(t *testing.T) {
	msg := types.MsgUpdateTssAddress{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgUpdateTssAddress, msg.Type())
}

func TestMessageUpdateTssAddress_Route(t *testing.T) {
	msg := types.MsgUpdateTssAddress{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMessageUpdateTssAddress_GetSignBytes(t *testing.T) {
	msg := types.MsgUpdateTssAddress{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
