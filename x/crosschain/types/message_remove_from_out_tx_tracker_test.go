package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMsgRemoveFromOutTxTracker_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  *types.MsgRemoveFromOutTxTracker
		err  error
	}{
		{
			name: "valid message",
			msg:  types.NewMsgRemoveFromOutTxTracker(sample.AccAddress(), 1, 0),
		},
		{
			name: "invalid creator address",
			msg:  types.NewMsgRemoveFromOutTxTracker("invalid", 1, 0),
			err:  sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid chain id",
			msg:  types.NewMsgRemoveFromOutTxTracker(sample.AccAddress(), -1, 0),
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

func TestMsgRemoveFromOutTxTracker_GetSigners(t *testing.T) {
	signer := sample.AccAddress()
	tests := []struct {
		name   string
		msg    *types.MsgRemoveFromOutTxTracker
		panics bool
	}{
		{
			name:   "valid signer",
			msg:    types.NewMsgRemoveFromOutTxTracker(signer, 1, 0),
			panics: false,
		},
		{
			name:   "invalid signer",
			msg:    types.NewMsgRemoveFromOutTxTracker("invalid", 1, 0),
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

func TestMsgRemoveFromOutTxTracker_Type(t *testing.T) {
	msg := types.NewMsgRemoveFromOutTxTracker(sample.AccAddress(), 1, 0)
	require.Equal(t, types.TypeMsgRemoveFromOutTxTracker, msg.Type())
}

func TestMsgRemoveFromOutTxTracker_Route(t *testing.T) {
	msg := types.NewMsgRemoveFromOutTxTracker(sample.AccAddress(), 1, 0)
	require.Equal(t, types.RouterKey, msg.Route())
}

func TestMsgRemoveFromOutTxTracker_GetSignBytes(t *testing.T) {
	msg := types.NewMsgRemoveFromOutTxTracker(sample.AccAddress(), 1, 0)
	require.NotPanics(t, func() {
		msg.GetSignBytes()
	})
}
