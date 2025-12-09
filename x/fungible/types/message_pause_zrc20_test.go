package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestMsgPauseZRC20_ValidateBasic(t *testing.T) {
	tt := []struct {
		name    string
		msg     *types.MsgPauseZRC20
		wantErr bool
	}{
		{
			name: "valid pause message",
			msg: types.NewMsgPauseZRC20(
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
			msg: types.NewMsgPauseZRC20(
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
			msg: types.NewMsgPauseZRC20(
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
			msg: types.NewMsgPauseZRC20(
				sample.AccAddress(),
				[]string{},
			),
			wantErr: true,
		},
		{
			name: "invalid zrc20 address",
			msg: types.NewMsgPauseZRC20(
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

func TestMsgPauseZRC20_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgPauseZRC20
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgPauseZRC20{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgPauseZRC20{
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

func TestMsgPauseZRC20_Type(t *testing.T) {
	msg := types.MsgPauseZRC20{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgPauseZrc20, msg.Type())
}

func TestMsgPauseZRC20_Route(t *testing.T) {
	msg := types.MsgPauseZRC20{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgPauseZRC20_GetSignBytes(t *testing.T) {
	msg := types.MsgPauseZRC20{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
