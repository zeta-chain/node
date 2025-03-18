package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestMsgUpdateOperationalChainParams_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgUpdateOperationalChainParams
		err  error
	}{
		{
			name: "valid message",
			msg: types.NewMsgUpdateOperationalChainParams(
				sample.AccAddress(),
				1,
				1,
				1,
				1,
				1,
				1,
				1,
				types.ConfirmationParams{
					SafeInboundCount:  1,
					FastInboundCount:  1,
					SafeOutboundCount: 1,
					FastOutboundCount: 1,
				},
			),
		},
		{
			name: "invalid address",
			msg: types.NewMsgUpdateOperationalChainParams(
				"invalid",
				1,
				1,
				1,
				1,
				1,
				1,
				1,
				types.ConfirmationParams{
					SafeInboundCount:  1,
					FastInboundCount:  1,
					SafeOutboundCount: 1,
					FastOutboundCount: 1,
				},
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid chain ID",
			msg: types.NewMsgUpdateOperationalChainParams(
				sample.AccAddress(),
				-1,
				1,
				1,
				1,
				1,
				1,
				1,
				types.ConfirmationParams{
					SafeInboundCount:  1,
					FastInboundCount:  1,
					SafeOutboundCount: 1,
					FastOutboundCount: 1,
				},
			),
			err: sdkerrors.ErrInvalidChainID,
		},
		{
			name: "invalid outbound schedule interval",
			msg: types.NewMsgUpdateOperationalChainParams(
				sample.AccAddress(),
				1,
				1,
				1,
				1,
				1,
				-1,
				1,
				types.ConfirmationParams{
					SafeInboundCount:  1,
					FastInboundCount:  1,
					SafeOutboundCount: 1,
					FastOutboundCount: 1,
				},
			),
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid outbound schedule look ahead",
			msg: types.NewMsgUpdateOperationalChainParams(
				sample.AccAddress(),
				1,
				1,
				1,
				1,
				1,
				1,
				-1,
				types.ConfirmationParams{
					SafeInboundCount:  1,
					FastInboundCount:  1,
					SafeOutboundCount: 1,
					FastOutboundCount: 1,
				},
			),
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid confirmation params",
			msg: types.NewMsgUpdateOperationalChainParams(
				sample.AccAddress(),
				1,
				1,
				1,
				1,
				1,
				1,
				1,
				types.ConfirmationParams{
					SafeInboundCount:  0,
					FastInboundCount:  1,
					SafeOutboundCount: 1,
					FastOutboundCount: 1,
				},
			),
			err: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgUpdateOperationalChainParams_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgUpdateOperationalChainParams
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgUpdateOperationalChainParams{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgUpdateOperationalChainParams{
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

func TestMsgUpdateOperationalChainParams_Type(t *testing.T) {
	msg := types.MsgUpdateOperationalChainParams{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgUpdateOperationalChainParams, msg.Type())
}

func TestMsgUpdateOperationalChainParams_Route(t *testing.T) {
	msg := types.MsgUpdateOperationalChainParams{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgUpdateOperationalChainParams_GetSignBytes(t *testing.T) {
	msg := types.MsgUpdateOperationalChainParams{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
