package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMsgAddOutboundTracker_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgAddOutboundTracker
		err  error
	}{
		{
			name: "invalid address",
			msg: types.NewMsgAddOutboundTracker(
				"invalid",
				1,
				1,
				"",
				nil,
				"",
				1,
			),
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid chain id",
			msg: types.NewMsgAddOutboundTracker(
				sample.AccAddress(),
				-1,
				1,
				"",
				nil,
				"",
				1,
			),
			err: sdkerrors.ErrInvalidChainID,
		},
		{
			name: "valid address",
			msg: types.NewMsgAddOutboundTracker(
				sample.AccAddress(),
				1,
				1,
				"",
				nil,
				"",
				1,
			),
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

func TestMsgAddOutboundTracker_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    types.MsgAddOutboundTracker
		panics bool
	}{
		{
			name: "valid signer",
			msg: types.MsgAddOutboundTracker{
				Creator: signer,
			},
			panics: false,
		},
		{
			name: "invalid signer",
			msg: types.MsgAddOutboundTracker{
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

func TestMsgAddOutboundTracker_Type(t *testing.T) {
	msg := types.MsgAddOutboundTracker{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.TypeMsgAddOutboundTracker, msg.Type())
}

func TestMsgAddOutboundTracker_Route(t *testing.T) {
	msg := types.MsgAddOutboundTracker{
		Creator: sample.AccAddress(),
	}
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgAddOutboundTracker_GetSignBytes(t *testing.T) {
	msg := types.MsgAddOutboundTracker{
		Creator: sample.AccAddress(),
	}
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
