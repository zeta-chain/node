package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestMsgUpdateContractBytecode_ValidateBasic(t *testing.T) {
	tt := []struct {
		name      string
		msg       *types.MsgUpdateContractBytecode
		wantError bool
	}{
		{
			name: "valid",
			msg: types.NewMsgUpdateContractBytecode(
				sample.AccAddress(),
				sample.EthAddress().Hex(),
				sample.Hash().Hex(),
			),
			wantError: false,
		},
		{
			name: "invalid creator",
			msg: types.NewMsgUpdateContractBytecode(
				"invalid",
				sample.EthAddress().Hex(),
				sample.Hash().Hex(),
			),
			wantError: true,
		},
		{
			name: "invalid contract address",
			msg: types.NewMsgUpdateContractBytecode(
				sample.AccAddress(),
				"invalid",
				sample.Hash().Hex(),
			),
			wantError: true,
		},
		{
			name: "invalid new code hash",
			msg: types.NewMsgUpdateContractBytecode(
				sample.AccAddress(),
				sample.EthAddress().Hex(),
				"invalid",
			),
			wantError: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgUpdateContractBytecode_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgUpdateContractBytecode
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgUpdateContractBytecode{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgUpdateContractBytecode{
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

func TestMsgUpdateContractBytecode_Type(t *testing.T) {
	msg := types.MsgUpdateContractBytecode{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgUpdateContractBytecode, msg.Type())
}

func TestMsgUpdateContractBytecode_Route(t *testing.T) {
	msg := types.MsgUpdateContractBytecode{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgUpdateContractBytecode_GetSignBytes(t *testing.T) {
	msg := types.MsgUpdateContractBytecode{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
