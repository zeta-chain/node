package types_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
)

func TestNewMsgMigrateConnectorFunds_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		msg     *types.MsgMigrateConnectorFunds
		wantErr bool
	}{
		{
			name: "valid message",
			msg: types.NewMsgMigrateConnectorFunds(
				sample.AccAddress(),
				1,
				"0x1234567890123456789012345678901234567890",
				sdkmath.NewUint(1000000),
			),
			wantErr: false,
		},
		{
			name: "invalid creator address",
			msg: types.NewMsgMigrateConnectorFunds(
				"invalid",
				1,
				"0x1234567890123456789012345678901234567890",
				sdkmath.NewUint(1000000),
			),
			wantErr: true,
		},
		{
			name: "empty new connector address",
			msg: types.NewMsgMigrateConnectorFunds(
				sample.AccAddress(),
				1,
				"",
				sdkmath.NewUint(1000000),
			),
			wantErr: true,
		},
		{
			name: "invalid new connector address",
			msg: types.NewMsgMigrateConnectorFunds(
				sample.AccAddress(),
				1,
				"invalid",
				sdkmath.NewUint(1000000),
			),
			wantErr: true,
		},
		{
			name: "zero amount",
			msg: types.NewMsgMigrateConnectorFunds(
				sample.AccAddress(),
				1,
				"0x1234567890123456789012345678901234567890",
				sdkmath.ZeroUint(),
			),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewMsgMigrateConnectorFunds_GetSigners(t *testing.T) {
	tests := []struct {
		name    string
		msg     types.MsgMigrateConnectorFunds
		wantErr bool
	}{
		{
			name: "valid signer",
			msg: types.MsgMigrateConnectorFunds{
				Creator: sample.AccAddress(),
			},
			wantErr: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgMigrateConnectorFunds{
				Creator: "invalid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				require.Panics(t, func() {
					tt.msg.GetSigners()
				})
			} else {
				signers := tt.msg.GetSigners()
				require.Len(t, signers, 1)
				require.Equal(t, tt.msg.Creator, signers[0].String())
			}
		})
	}
}

func TestNewMsgMigrateConnectorFunds_Type(t *testing.T) {
	msg := types.MsgMigrateConnectorFunds{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgMigrateConnectorFunds, msg.Type())
}

func TestNewMsgMigrateConnectorFunds_Route(t *testing.T) {
	msg := types.MsgMigrateConnectorFunds{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestNewMsgMigrateConnectorFunds_GetSignBytes(t *testing.T) {
	msg := types.MsgMigrateConnectorFunds{
		Creator: sample.AccAddress(),
	}
	require.NotNil(t, msg.GetSignBytes())
}
