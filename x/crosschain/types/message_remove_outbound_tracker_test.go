package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMsgRemoveOutboundTracker_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgRemoveOutboundTracker
		err  error
	}{
		{
			name: "valid message",
			msg:  types.NewMsgRemoveOutboundTracker(sample.AccAddress(), 1, 0),
		},
		{
			name: "invalid creator address",
			msg:  types.NewMsgRemoveOutboundTracker("invalid", 1, 0),
			err:  sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid chain id",
			msg:  types.NewMsgRemoveOutboundTracker(sample.AccAddress(), -1, 0),
			err:  sdkerrors.ErrInvalidChainID,
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

func TestMsgRemoveOutboundTracker_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    *types.MsgRemoveOutboundTracker
		panics bool
	}{
		{
			name:   "valid signer",
			msg:    types.NewMsgRemoveOutboundTracker(signer, 1, 0),
			panics: false,
		},
		{
			name:   "invalid signer",
			msg:    types.NewMsgRemoveOutboundTracker("invalid", 1, 0),
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

func TestMsgRemoveOutboundTracker_Type(t *testing.T) {
	msg := types.NewMsgRemoveOutboundTracker(sample.AccAddress(), 1, 0)
	require.Equal(t, types.TypeMsgRemoveOutboundTracker, msg.Type())
}

func TestMsgRemoveOutboundTracker_Route(t *testing.T) {
	msg := types.NewMsgRemoveOutboundTracker(sample.AccAddress(), 1, 0)
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgRemoveOutboundTracker_GetSignBytes(t *testing.T) {
	msg := types.NewMsgRemoveOutboundTracker(sample.AccAddress(), 1, 0)
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
