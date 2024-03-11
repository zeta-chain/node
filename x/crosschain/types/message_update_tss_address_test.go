package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMessageUpdateTssAddress_ValidateBasic(t *testing.T) {
	tests := []struct {
		name  string
		msg   crosschaintypes.MsgUpdateTssAddress
		error bool
	}{
		{
			name: "invalid creator",
			msg: crosschaintypes.MsgUpdateTssAddress{
				Creator:   "invalid_address",
				TssPubkey: "zetapub1addwnpepq28c57cvcs0a2htsem5zxr6qnlvq9mzhmm76z3jncsnzz32rclangr2g35p",
			},
			error: true,
		},
		{
			name: "invalid pubkey",
			msg: crosschaintypes.MsgUpdateTssAddress{
				Creator:   "zeta15ruj2tc76pnj9xtw64utktee7cc7w6vzaes73z",
				TssPubkey: "zetapub1addwnpepq28c57cvcs0a2htsem5zxr6qnlvq9mzhmm",
			},
			error: true,
		},
		{
			name: "valid msg",
			msg: crosschaintypes.MsgUpdateTssAddress{
				Creator:   "zeta15ruj2tc76pnj9xtw64utktee7cc7w6vzaes73z",
				TssPubkey: "zetapub1addwnpepq28c57cvcs0a2htsem5zxr6qnlvq9mzhmm76z3jncsnzz32rclangr2g35p",
			},
			error: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keeper.SetConfig(false)
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
		msg    crosschaintypes.MsgUpdateTssAddress
		panics bool
	}{
		{
			name: "valid signer",
			msg: crosschaintypes.MsgUpdateTssAddress{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: crosschaintypes.MsgUpdateTssAddress{
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

func TestMessageUpdateTssAddress_Type(t *testing.T) {
	msg := crosschaintypes.MsgUpdateTssAddress{
		Creator: sample.AccAddress(),
	}
	assert.Equal(t, types.TypeMsgUpdateTssAddress, msg.Type())
}

func TestMessageUpdateTssAddress_Route(t *testing.T) {
	msg := crosschaintypes.MsgUpdateTssAddress{
		Creator: sample.AccAddress(),
	}
	assert.Equal(t, types.RouterKey, msg.Route())
}

func TestMessageUpdateTssAddress_GetSignBytes(t *testing.T) {
	msg := crosschaintypes.MsgUpdateTssAddress{
		Creator: sample.AccAddress(),
	}
	assert.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
