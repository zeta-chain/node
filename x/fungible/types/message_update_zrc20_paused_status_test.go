package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestMMsgUpdateZRC20PausedStatus_ValidateBasic(t *testing.T) {
	sampleAddress := sample.EthAddress().String()
	tt := []struct {
		name    string
		msg     types.MsgUpdateZRC20PausedStatus
		wantErr bool
	}{
		{
			name: "valid pause message",
			msg: types.MsgUpdateZRC20PausedStatus{
				Creator: sample.AccAddress(),
				Zrc20Addresses: []string{
					sample.EthAddress().String(),
					sample.EthAddress().String(),
					sample.EthAddress().String(),
				},
				Action: types.UpdatePausedStatusAction_PAUSE,
			},
			wantErr: false,
		},
		{
			name: "valid unpause message",
			msg: types.MsgUpdateZRC20PausedStatus{
				Creator: sample.AccAddress(),
				Zrc20Addresses: []string{
					sample.EthAddress().String(),
					sample.EthAddress().String(),
					sample.EthAddress().String(),
				},
				Action: types.UpdatePausedStatusAction_UNPAUSE,
			},
			wantErr: false,
		},
		{
			name: "invalid creator address",
			msg: types.MsgUpdateZRC20PausedStatus{
				Creator: "invalid",
				Zrc20Addresses: []string{
					sample.EthAddress().String(),
					sample.EthAddress().String(),
					sample.EthAddress().String(),
				},
				Action: types.UpdatePausedStatusAction_PAUSE,
			},
			wantErr: true,
		},
		{
			name: "invalid empty zrc20 address",
			msg: types.MsgUpdateZRC20PausedStatus{
				Creator:        sample.AccAddress(),
				Zrc20Addresses: []string{},
				Action:         types.UpdatePausedStatusAction_PAUSE,
			},
			wantErr: true,
		},
		{
			name: "invalid zrc20 address",
			msg: types.MsgUpdateZRC20PausedStatus{
				Creator: sample.AccAddress(),
				Zrc20Addresses: []string{
					sample.EthAddress().String(),
					"invalid",
					sample.EthAddress().String(),
				},
				Action: types.UpdatePausedStatusAction_PAUSE,
			},
			wantErr: true,
		},
		{
			name: "duplicate zrc20 address",
			msg: types.MsgUpdateZRC20PausedStatus{
				Creator: sample.AccAddress(),
				Zrc20Addresses: []string{
					sampleAddress,
					sampleAddress,
				},
				Action: types.UpdatePausedStatusAction_PAUSE,
			},
			wantErr: true,
		},
		{
			name: "invalid action",
			msg: types.MsgUpdateZRC20PausedStatus{
				Creator: sample.AccAddress(),
				Zrc20Addresses: []string{
					sample.EthAddress().String(),
					sample.EthAddress().String(),
					sample.EthAddress().String(),
				},
				Action: 3,
			},
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
