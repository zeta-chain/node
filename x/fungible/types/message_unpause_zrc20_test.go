package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestMsgUnpauseZRC20_ValidateBasic(t *testing.T) {
	tt := []struct {
		name    string
		msg     *types.MsgUnpauseZRC20
		wantErr bool
	}{
		{
			name: "valid unpause message",
			msg: types.NewMsgUnpauseZRC20(
				sample.AccAddress(),
				[]string{
					sample.EthAddress().String(),
					sample.EthAddress().String(),
					sample.EthAddress().String(),
				},
			),
			wantErr: false,
		},
		{
			name: "valid unpause message",
			msg: types.NewMsgUnpauseZRC20(
				sample.AccAddress(),
				[]string{
					sample.EthAddress().String(),
					sample.EthAddress().String(),
					sample.EthAddress().String(),
				},
			),
			wantErr: false,
		},
		{
			name: "invalid creator address",
			msg: types.NewMsgUnpauseZRC20(
				"invalid",
				[]string{
					sample.EthAddress().String(),
					sample.EthAddress().String(),
					sample.EthAddress().String(),
				},
			),
			wantErr: true,
		},
		{
			name: "invalid empty zrc20 address",
			msg: types.NewMsgUnpauseZRC20(
				sample.AccAddress(),
				[]string{},
			),
			wantErr: true,
		},
		{
			name: "invalid zrc20 address",
			msg: types.NewMsgUnpauseZRC20(
				sample.AccAddress(),
				[]string{
					sample.EthAddress().String(),
					"invalid",
					sample.EthAddress().String(),
				},
			),
			wantErr: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgUnpauseZRC20_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgUnpauseZRC20
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgUnpauseZRC20{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgUnpauseZRC20{
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

func TestMsgUnpauseZRC20_Type(t *testing.T) {
	msg := types.MsgUnpauseZRC20{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgUnpauseZRC20, msg.Type())
}

func TestMsgUnpauseZRC20_Route(t *testing.T) {
	msg := types.MsgUnpauseZRC20{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgUnpauseZRC20_GetSignBytes(t *testing.T) {
	msg := types.MsgUnpauseZRC20{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
