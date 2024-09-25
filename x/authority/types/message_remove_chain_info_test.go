package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/authority/types"
)

func TestMsgRemoveChainInfo_ValidateBasic(t *testing.T) {
	tests := []struct {
		name        string
		msg         *types.MsgRemoveChainInfo
		errContains string
	}{
		{
			name: "valid message",
			msg:  types.NewMsgRemoveChainInfo(sample.AccAddress(), 42),
		},
		{
			name:        "invalid creator address",
			msg:         types.NewMsgRemoveChainInfo("invalid", 42),
			errContains: "invalid creator address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.errContains != "" {
				require.ErrorContains(t, err, tt.errContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgRemoveChainInfo_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    *types.MsgRemoveChainInfo
		panics bool
	}{
		{
			name:   "valid signer",
			msg:    types.NewMsgRemoveChainInfo(signer, 42),
			panics: false,
		},
		{
			name:   "invalid signer",
			msg:    types.NewMsgRemoveChainInfo("invalid", 42),
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

func TestMsgRemoveChainInfo_Type(t *testing.T) {
	msg := types.NewMsgRemoveChainInfo(sample.AccAddress(), 42)
	require.Equal(t, types.TypeMsgRemoveChainInfo, msg.Type())
}

func TestMsgRemoveChainInfo_Route(t *testing.T) {
	msg := types.NewMsgRemoveChainInfo(sample.AccAddress(), 42)
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgRemoveChainInfo_GetSignBytes(t *testing.T) {
	msg := types.NewMsgRemoveChainInfo(sample.AccAddress(), 42)
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
