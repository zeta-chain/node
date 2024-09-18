package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestMsgVoteTSS_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgVoteTSS
		err  error
	}{
		{
			name: "valid message",
			msg:  types.NewMsgVoteTSS(sample.AccAddress(), "pubkey", 1, chains.ReceiveStatus_success),
		},
		{
			name: "valid message with receive status failed",
			msg:  types.NewMsgVoteTSS(sample.AccAddress(), "pubkey", 1, chains.ReceiveStatus_failed),
		},
		{
			name: "invalid creator address",
			msg:  types.NewMsgVoteTSS("invalid", "pubkey", 1, chains.ReceiveStatus_success),
			err:  sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid observation status",
			msg:  types.NewMsgVoteTSS(sample.AccAddress(), "pubkey", 1, chains.ReceiveStatus_created),
			err:  sdkerrors.ErrInvalidRequest,
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

func TestMsgVoteTSS_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    *types.MsgVoteTSS
		panics bool
	}{
		{
			name:   "valid signer",
			msg:    types.NewMsgVoteTSS(signer, "pubkey", 1, chains.ReceiveStatus_success),
			panics: false,
		},
		{
			name:   "invalid signer",
			msg:    types.NewMsgVoteTSS("invalid", "pubkey", 1, chains.ReceiveStatus_success),
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

func TestMsgVoteTSS_Type(t *testing.T) {
	msg := types.NewMsgVoteTSS(sample.AccAddress(), "pubkey", 1, chains.ReceiveStatus_success)
	require.Equal(t, types.TypeMsgVoteTSS, msg.Type())
}

func TestMsgVoteTSS_Route(t *testing.T) {
	msg := types.NewMsgVoteTSS(sample.AccAddress(), "pubkey", 1, chains.ReceiveStatus_success)
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgVoteTSS_GetSignBytes(t *testing.T) {
	msg := types.NewMsgVoteTSS(sample.AccAddress(), "pubkey", 1, chains.ReceiveStatus_success)
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}

func TestMsgVoteTSS_Digest(t *testing.T) {
	vote1 := types.NewMsgVoteTSS(sample.AccAddress(), "pubkey", 1, chains.ReceiveStatus_success)
	require.Equal(t, "1-pubkey-tss-keygen", vote1.Digest())

	vote2 := types.NewMsgVoteTSS(sample.AccAddress(), "pubkey", 1, chains.ReceiveStatus_success)
	require.Equal(t, "1-pubkey-tss-keygen", vote2.Digest())

	require.Equal(t, vote1.Digest(), vote2.Digest())

	vote3 := types.NewMsgVoteTSS(sample.AccAddress(), "pubkeyNew", 2, chains.ReceiveStatus_success)
	require.Equal(t, "2-pubkeyNew-tss-keygen", vote3.Digest())
	// Different pubkey changes digest
	require.NotEqual(t, vote1.Digest(), vote3.Digest())

	vote4 := types.NewMsgVoteTSS(sample.AccAddress(), "pubkey", 3, chains.ReceiveStatus_success)
	require.Equal(t, "3-pubkey-tss-keygen", vote4.Digest())
	// Different height changes digest
	require.NotEqual(t, vote1.Digest(), vote4.Digest())
}
