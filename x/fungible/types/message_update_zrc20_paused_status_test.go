package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestMMsgUpdateZRC20PausedStatus_ValidateBasic(t *testing.T) {
	tt := []struct {
		name    string
		msg     *types.MsgUpdateZRC20PausedStatus
		wantErr bool
	}{
		{
			name: "valid pause message",
			msg: types.NewMsgUpdateZRC20PausedStatus(
				sample.AccAddress(),
				[]string{
					sample.EthAddress().String(),
					sample.EthAddress().String(),
					sample.EthAddress().String(),
				},
				types.UpdatePausedStatusAction_PAUSE,
			),
			wantErr: false,
		},
		{
			name: "valid unpause message",
			msg: types.NewMsgUpdateZRC20PausedStatus(
				sample.AccAddress(),
				[]string{
					sample.EthAddress().String(),
					sample.EthAddress().String(),
					sample.EthAddress().String(),
				},
				types.UpdatePausedStatusAction_UNPAUSE,
			),
			wantErr: false,
		},
		{
			name: "invalid creator address",
			msg: types.NewMsgUpdateZRC20PausedStatus(
				"invalid",
				[]string{
					sample.EthAddress().String(),
					sample.EthAddress().String(),
					sample.EthAddress().String(),
				},
				types.UpdatePausedStatusAction_PAUSE,
			),
			wantErr: true,
		},
		{
			name: "invalid empty zrc20 address",
			msg: types.NewMsgUpdateZRC20PausedStatus(
				sample.AccAddress(),
				[]string{},
				types.UpdatePausedStatusAction_PAUSE,
			),
			wantErr: true,
		},
		{
			name: "invalid zrc20 address",
			msg: types.NewMsgUpdateZRC20PausedStatus(
				sample.AccAddress(),
				[]string{
					sample.EthAddress().String(),
					"invalid",
					sample.EthAddress().String(),
				},
				types.UpdatePausedStatusAction_PAUSE,
			),
			wantErr: true,
		},
		{
			name: "invalid action",
			msg: types.NewMsgUpdateZRC20PausedStatus(
				sample.AccAddress(),
				[]string{
					sample.EthAddress().String(),
					sample.EthAddress().String(),
					sample.EthAddress().String(),
				},
				3,
			),
			wantErr: true,
		},
	}
	for _, tc := range tt {
		tc := tc
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

func TestMMsgUpdateZRC20PausedStatus_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgUpdateZRC20PausedStatus
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgUpdateZRC20PausedStatus{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgUpdateZRC20PausedStatus{
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

func TestMMsgUpdateZRC20PausedStatus_Type(t *testing.T) {
	msg := types.MsgUpdateZRC20PausedStatus{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgUpdateZRC20PausedStatus, msg.Type())
}

func TestMMsgUpdateZRC20PausedStatus_Route(t *testing.T) {
	msg := types.MsgUpdateZRC20PausedStatus{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMMsgUpdateZRC20PausedStatus_GetSignBytes(t *testing.T) {
	msg := types.MsgUpdateZRC20PausedStatus{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
